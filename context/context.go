// Copyright (c) 2025 A Bit of Help, Inc.

// Package context provides utilities for working with Go's context package.
// It includes functions for creating contexts with various timeouts, adding and retrieving
// values from contexts, and checking context status.
package context

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/google/uuid"
)

// Key represents a key for context values
type Key string

// Context keys
const (
	RequestIDKey     Key = "request_id"
	TraceIDKey       Key = "trace_id"
	UserIDKey        Key = "user_id"
	TenantIDKey      Key = "tenant_id"
	OperationKey     Key = "operation"
	CorrelationIDKey Key = "correlation_id"
	ServiceNameKey   Key = "service_name"
)

// Default timeout values
const (
	DefaultTimeout         = 30 * time.Second
	ShortTimeout           = 5 * time.Second
	LongTimeout            = 60 * time.Second
	DatabaseTimeout        = 10 * time.Second
	NetworkTimeout         = 15 * time.Second
	ExternalServiceTimeout = 20 * time.Second
)

// ContextOptions contains options for creating a context
type ContextOptions struct {
	// Timeout is the duration after which the context will be canceled
	Timeout time.Duration

	// RequestID is a unique identifier for the request
	RequestID string

	// TraceID is a unique identifier for tracing
	TraceID string

	// UserID is the ID of the user making the request
	UserID string

	// TenantID is the ID of the tenant
	TenantID string

	// Operation is the name of the operation being performed
	Operation string

	// CorrelationID is a unique identifier for correlating related operations
	CorrelationID string

	// ServiceName is the name of the service
	ServiceName string

	// Parent is the parent context
	Parent context.Context
}

// NewContext creates a new context with the specified options
func NewContext(opts ContextOptions) (context.Context, context.CancelFunc) {
	if opts.Parent == nil {
		opts.Parent = context.Background()
	}

	if opts.Timeout > 0 {
		return WithTimeout(opts.Parent, opts.Timeout, opts)
	}

	// If no timeout is specified, create a context with cancel
	ctx, cancel := context.WithCancel(opts.Parent)
	return enrichContext(ctx, opts), cancel
}

// WithTimeout creates a new context with the specified timeout and options
func WithTimeout(ctx context.Context, timeout time.Duration, opts ContextOptions) (context.Context, context.CancelFunc) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	return enrichContext(timeoutCtx, opts), cancel
}

// enrichContext adds the specified options to the context
func enrichContext(ctx context.Context, opts ContextOptions) context.Context {
	if opts.RequestID != "" {
		ctx = context.WithValue(ctx, RequestIDKey, opts.RequestID)
	} else if GetRequestID(ctx) == "" {
		ctx = context.WithValue(ctx, RequestIDKey, uuid.New().String())
	}

	if opts.TraceID != "" {
		ctx = context.WithValue(ctx, TraceIDKey, opts.TraceID)
	} else if GetTraceID(ctx) == "" {
		ctx = context.WithValue(ctx, TraceIDKey, uuid.New().String())
	}

	if opts.UserID != "" {
		ctx = context.WithValue(ctx, UserIDKey, opts.UserID)
	}

	if opts.TenantID != "" {
		ctx = context.WithValue(ctx, TenantIDKey, opts.TenantID)
	}

	if opts.Operation != "" {
		ctx = context.WithValue(ctx, OperationKey, opts.Operation)
	}

	if opts.CorrelationID != "" {
		ctx = context.WithValue(ctx, CorrelationIDKey, opts.CorrelationID)
	} else if GetCorrelationID(ctx) == "" {
		ctx = context.WithValue(ctx, CorrelationIDKey, uuid.New().String())
	}

	if opts.ServiceName != "" {
		ctx = context.WithValue(ctx, ServiceNameKey, opts.ServiceName)
	}

	return ctx
}

// WithDefaultTimeout creates a new context with the default timeout
func WithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, DefaultTimeout, ContextOptions{Parent: ctx})
}

// WithShortTimeout creates a new context with a short timeout
func WithShortTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, ShortTimeout, ContextOptions{Parent: ctx})
}

// WithLongTimeout creates a new context with a long timeout
func WithLongTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, LongTimeout, ContextOptions{Parent: ctx})
}

// WithDatabaseTimeout creates a new context with a database timeout
func WithDatabaseTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, DatabaseTimeout, ContextOptions{
		Parent:    ctx,
		Operation: "database",
	})
}

// WithNetworkTimeout creates a new context with a network timeout
func WithNetworkTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, NetworkTimeout, ContextOptions{
		Parent:    ctx,
		Operation: "network",
	})
}

// WithExternalServiceTimeout creates a new context with an external service timeout
func WithExternalServiceTimeout(ctx context.Context, serviceName string) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, ExternalServiceTimeout, ContextOptions{
		Parent:      ctx,
		Operation:   "external_service",
		ServiceName: serviceName,
	})
}

// WithOperation creates a new context with an operation name
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// GetOperation gets the operation name from the context
func GetOperation(ctx context.Context) string {
	if op, ok := ctx.Value(OperationKey).(string); ok {
		return op
	}
	return ""
}

// WithCorrelationID adds a correlation ID to the context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

// GetCorrelationID gets the correlation ID from the context
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}

// WithServiceName adds a service name to the context
func WithServiceName(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, ServiceNameKey, serviceName)
}

// GetServiceName gets the service name from the context
func GetServiceName(ctx context.Context) string {
	if name, ok := ctx.Value(ServiceNameKey).(string); ok {
		return name
	}
	return ""
}

// CheckContext checks if the context is done and returns an appropriate error
func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == context.DeadlineExceeded {
			operation := GetOperation(ctx)
			service := GetServiceName(ctx)
			if operation != "" {
				if service != "" {
					// Wrap ErrTimeout with additional context
					return fmt.Errorf("operation '%s' for service '%s' timed out: %w", operation, service, errors.ErrTimeout)
				}
				return fmt.Errorf("operation '%s' timed out: %w", operation, errors.ErrTimeout)
			}
			return errors.ErrTimeout
		}
		if err == context.Canceled {
			operation := GetOperation(ctx)
			if operation != "" {
				// Wrap ErrCanceled with additional context
				return fmt.Errorf("operation '%s' was canceled: %w", operation, errors.ErrCanceled)
			}
			return errors.ErrCanceled
		}
		return errors.Wrap(err, errors.InternalErrorCode, "context error")
	default:
		return nil
	}
}

// MustCheck checks if the context is done and panics with an appropriate error if it is
func MustCheck(ctx context.Context) {
	if err := CheckContext(ctx); err != nil {
		panic(err)
	}
}

// WithValues creates a new context with the specified key-value pairs
func WithValues(ctx context.Context, keyValues ...interface{}) context.Context {
	if len(keyValues)%2 != 0 {
		// If odd number of arguments, ignore the last one
		keyValues = keyValues[:len(keyValues)-1]
	}

	for i := 0; i < len(keyValues); i += 2 {
		ctx = context.WithValue(ctx, keyValues[i], keyValues[i+1])
	}

	return ctx
}

// Background returns a non-nil, empty Context. It is never canceled, has no
// values, and has no deadline.
func Background() context.Context {
	return context.Background()
}

// TODO returns a non-nil, empty Context. Code should use context.TODO when
// it's unclear which Context to use or it is not yet available (because the
// surrounding function has not yet been extended to accept a Context parameter).
func TODO() context.Context {
	return context.TODO()
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestIDKey, uuid.New().String())
}

// GetRequestID gets the request ID from the context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceIDKey, uuid.New().String())
}

// GetTraceID gets the trace ID from the context
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return ""
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID gets the user ID from the context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// WithTenantID adds a tenant ID to the context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// GetTenantID gets the tenant ID from the context
func GetTenantID(ctx context.Context) string {
	if id, ok := ctx.Value(TenantIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextInfo returns a string with all the context information
func ContextInfo(ctx context.Context) string {
	info := fmt.Sprintf("RequestID: %s, TraceID: %s", GetRequestID(ctx), GetTraceID(ctx))

	if userID := GetUserID(ctx); userID != "" {
		info += fmt.Sprintf(", UserID: %s", userID)
	}

	if tenantID := GetTenantID(ctx); tenantID != "" {
		info += fmt.Sprintf(", TenantID: %s", tenantID)
	}

	if operation := GetOperation(ctx); operation != "" {
		info += fmt.Sprintf(", Operation: %s", operation)
	}

	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		info += fmt.Sprintf(", CorrelationID: %s", correlationID)
	}

	if serviceName := GetServiceName(ctx); serviceName != "" {
		info += fmt.Sprintf(", ServiceName: %s", serviceName)
	}

	deadline, ok := ctx.Deadline()
	if ok {
		info += fmt.Sprintf(", Deadline: %s", deadline.Format(time.RFC3339))
	}

	return info
}
