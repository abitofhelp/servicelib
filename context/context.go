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

// Key represents a key for context values in the context package.
// It is a type alias for string that provides type safety for context keys.
type Key string

// Context keys define the standard keys used for storing and retrieving values from contexts.
// These keys should be used consistently throughout the application to ensure
// that context values can be reliably accessed.
const (
	// RequestIDKey is the key for storing and retrieving the request ID.
	RequestIDKey Key = "request_id"

	// TraceIDKey is the key for storing and retrieving the trace ID.
	TraceIDKey Key = "trace_id"

	// UserIDKey is the key for storing and retrieving the user ID.
	UserIDKey Key = "user_id"

	// TenantIDKey is the key for storing and retrieving the tenant ID.
	TenantIDKey Key = "tenant_id"

	// OperationKey is the key for storing and retrieving the operation name.
	OperationKey Key = "operation"

	// CorrelationIDKey is the key for storing and retrieving the correlation ID.
	CorrelationIDKey Key = "correlation_id"

	// ServiceNameKey is the key for storing and retrieving the service name.
	ServiceNameKey Key = "service_name"
)

// Default timeout values define standard durations for different types of operations.
// These constants should be used consistently throughout the application to ensure
// that timeouts are appropriate for the type of operation being performed.
const (
	// DefaultTimeout is the standard timeout for general operations (30 seconds).
	DefaultTimeout = 30 * time.Second

	// ShortTimeout is a brief timeout for quick operations (5 seconds).
	ShortTimeout = 5 * time.Second

	// LongTimeout is an extended timeout for operations that may take longer (60 seconds).
	LongTimeout = 60 * time.Second

	// DatabaseTimeout is the timeout for database operations (10 seconds).
	DatabaseTimeout = 10 * time.Second

	// NetworkTimeout is the timeout for network operations (15 seconds).
	NetworkTimeout = 15 * time.Second

	// ExternalServiceTimeout is the timeout for calls to external services (20 seconds).
	ExternalServiceTimeout = 20 * time.Second
)

// ContextOptions contains options for creating a new context.
// This struct allows for flexible context creation with various metadata
// and configuration options.
type ContextOptions struct {
	// Timeout is the duration after which the context will be canceled.
	// If set to 0, no timeout will be applied.
	Timeout time.Duration

	// RequestID is a unique identifier for the request.
	// If empty, a new UUID will be generated.
	RequestID string

	// TraceID is a unique identifier for distributed tracing.
	// If empty, a new UUID will be generated.
	TraceID string

	// UserID is the ID of the user making the request.
	// This is useful for authentication and authorization.
	UserID string

	// TenantID is the ID of the tenant in multi-tenant applications.
	// This is useful for data isolation and authorization.
	TenantID string

	// Operation is the name of the operation being performed.
	// This is useful for logging and error reporting.
	Operation string

	// CorrelationID is a unique identifier for correlating related operations.
	// If empty, a new UUID will be generated.
	CorrelationID string

	// ServiceName is the name of the service handling the request.
	// This is useful for distributed tracing and logging.
	ServiceName string

	// Parent is the parent context from which the new context will be derived.
	// If nil, context.Background() will be used.
	Parent context.Context
}

// NewContext creates a new context with the specified options.
// It enriches the context with metadata from the provided options and
// applies a timeout if specified.
//
// Parameters:
//   - opts: The options for creating the context, including timeout, metadata, and parent context.
//
// Returns:
//   - context.Context: The newly created context with all specified options applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
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

// WithTimeout creates a new context with the specified timeout and options.
// It applies the timeout to the context and enriches it with metadata from the provided options.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - timeout: The duration after which the context will be canceled.
//   - opts: The options for creating the context, including metadata.
//
// Returns:
//   - context.Context: The newly created context with timeout and options applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
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

// WithDefaultTimeout creates a new context with the default timeout (30 seconds).
// This is a convenience function for creating a context with a standard timeout
// for general operations.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with the default timeout applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
func WithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, DefaultTimeout, ContextOptions{Parent: ctx})
}

// WithShortTimeout creates a new context with a short timeout (5 seconds).
// This is a convenience function for creating a context with a brief timeout
// for operations that should complete quickly.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with the short timeout applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
func WithShortTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, ShortTimeout, ContextOptions{Parent: ctx})
}

// WithLongTimeout creates a new context with a long timeout (60 seconds).
// This is a convenience function for creating a context with an extended timeout
// for operations that may take longer to complete.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with the long timeout applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
func WithLongTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, LongTimeout, ContextOptions{Parent: ctx})
}

// WithDatabaseTimeout creates a new context with a database timeout (10 seconds).
// This is a convenience function for creating a context specifically for database operations.
// It also sets the operation name to "database" for better error reporting and logging.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with the database timeout applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
func WithDatabaseTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, DatabaseTimeout, ContextOptions{
		Parent:    ctx,
		Operation: "database",
	})
}

// WithNetworkTimeout creates a new context with a network timeout (15 seconds).
// This is a convenience function for creating a context specifically for network operations.
// It also sets the operation name to "network" for better error reporting and logging.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with the network timeout applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
func WithNetworkTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, NetworkTimeout, ContextOptions{
		Parent:    ctx,
		Operation: "network",
	})
}

// WithExternalServiceTimeout creates a new context with an external service timeout (20 seconds).
// This is a convenience function for creating a context specifically for external service calls.
// It sets the operation name to "external_service" and includes the service name for better
// error reporting and logging.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - serviceName: The name of the external service being called.
//
// Returns:
//   - context.Context: The newly created context with the external service timeout applied.
//   - context.CancelFunc: A cancel function that can be called to cancel the context.
func WithExternalServiceTimeout(ctx context.Context, serviceName string) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, ExternalServiceTimeout, ContextOptions{
		Parent:      ctx,
		Operation:   "external_service",
		ServiceName: serviceName,
	})
}

// WithOperation creates a new context with an operation name.
// The operation name is used for logging, error reporting, and tracing to identify
// the specific operation being performed.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - operation: The name of the operation being performed.
//
// Returns:
//   - context.Context: The newly created context with the operation name added.
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// GetOperation gets the operation name from the context.
// This function retrieves the operation name that was previously set with WithOperation.
//
// Parameters:
//   - ctx: The context from which to retrieve the operation name.
//
// Returns:
//   - string: The operation name, or an empty string if no operation name was set.
func GetOperation(ctx context.Context) string {
	if op, ok := ctx.Value(OperationKey).(string); ok {
		return op
	}
	return ""
}

// WithCorrelationID adds a correlation ID to the context.
// The correlation ID is used to track related operations across multiple services
// or components, especially in distributed systems.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - correlationID: The correlation ID to add to the context.
//
// Returns:
//   - context.Context: The newly created context with the correlation ID added.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

// GetCorrelationID gets the correlation ID from the context.
// This function retrieves the correlation ID that was previously set with WithCorrelationID.
//
// Parameters:
//   - ctx: The context from which to retrieve the correlation ID.
//
// Returns:
//   - string: The correlation ID, or an empty string if no correlation ID was set.
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}

// WithServiceName adds a service name to the context.
// The service name is used to identify the service handling the request,
// which is useful for logging, error reporting, and distributed tracing.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - serviceName: The name of the service to add to the context.
//
// Returns:
//   - context.Context: The newly created context with the service name added.
func WithServiceName(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, ServiceNameKey, serviceName)
}

// GetServiceName gets the service name from the context.
// This function retrieves the service name that was previously set with WithServiceName.
//
// Parameters:
//   - ctx: The context from which to retrieve the service name.
//
// Returns:
//   - string: The service name, or an empty string if no service name was set.
func GetServiceName(ctx context.Context) string {
	if name, ok := ctx.Value(ServiceNameKey).(string); ok {
		return name
	}
	return ""
}

// CheckContext checks if the context is done and returns an appropriate error.
// This function is useful for checking if a context has been canceled or has timed out
// before proceeding with an operation. It returns nil if the context is still valid,
// or an appropriate error if the context is done.
//
// Parameters:
//   - ctx: The context to check.
//
// Returns:
//   - error: nil if the context is still valid, or an error describing why the context is done.
//     The error will be errors.ErrTimeout if the context timed out, errors.ErrCanceled if the
//     context was canceled, or a wrapped error for other context errors.
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

// MustCheck checks if the context is done and panics with an appropriate error if it is.
// This function is similar to CheckContext, but instead of returning an error, it panics
// if the context is done. This is useful in situations where a context error should be
// treated as a fatal error.
//
// Parameters:
//   - ctx: The context to check.
//
// Panics:
//   - If the context is done, this function will panic with the error returned by CheckContext.
func MustCheck(ctx context.Context) {
	if err := CheckContext(ctx); err != nil {
		panic(err)
	}
}

// WithValues creates a new context with the specified key-value pairs.
// This function allows adding multiple key-value pairs to a context in a single call.
// The keys and values are provided as alternating arguments, where each key is followed
// by its corresponding value.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - keyValues: Alternating key-value pairs to add to the context. The number of arguments
//     should be even. If an odd number of arguments is provided, the last one is ignored.
//
// Returns:
//   - context.Context: The newly created context with all key-value pairs added.
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

// Background returns a non-nil, empty Context.
// It is never canceled, has no values, and has no deadline. This is a wrapper
// around context.Background() from the standard library.
//
// Returns:
//   - context.Context: A background context that is never canceled.
func Background() context.Context {
	return context.Background()
}

// TODO returns a non-nil, empty Context.
// This is a wrapper around context.TODO() from the standard library.
// Code should use TODO when it's unclear which Context to use or it is not yet
// available (because the surrounding function has not yet been extended to accept
// a Context parameter).
//
// Returns:
//   - context.Context: A TODO context that indicates the context should be provided
//     in the future.
func TODO() context.Context {
	return context.TODO()
}

// WithRequestID adds a request ID to the context.
// This function generates a new UUID and adds it to the context as the request ID.
// The request ID is useful for tracking requests through the system, especially
// in logs and error reports.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with a request ID added.
func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestIDKey, uuid.New().String())
}

// GetRequestID gets the request ID from the context.
// This function retrieves the request ID that was previously set with WithRequestID.
//
// Parameters:
//   - ctx: The context from which to retrieve the request ID.
//
// Returns:
//   - string: The request ID, or an empty string if no request ID was set.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithTraceID adds a trace ID to the context.
// This function generates a new UUID and adds it to the context as the trace ID.
// The trace ID is useful for distributed tracing, allowing requests to be tracked
// across multiple services.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//
// Returns:
//   - context.Context: The newly created context with a trace ID added.
func WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceIDKey, uuid.New().String())
}

// GetTraceID gets the trace ID from the context.
// This function retrieves the trace ID that was previously set with WithTraceID.
//
// Parameters:
//   - ctx: The context from which to retrieve the trace ID.
//
// Returns:
//   - string: The trace ID, or an empty string if no trace ID was set.
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return ""
}

// WithUserID adds a user ID to the context.
// The user ID is useful for authentication, authorization, and auditing purposes.
// It allows operations to be associated with a specific user.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - userID: The ID of the user to add to the context.
//
// Returns:
//   - context.Context: The newly created context with the user ID added.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID gets the user ID from the context.
// This function retrieves the user ID that was previously set with WithUserID.
//
// Parameters:
//   - ctx: The context from which to retrieve the user ID.
//
// Returns:
//   - string: The user ID, or an empty string if no user ID was set.
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// WithTenantID adds a tenant ID to the context.
// The tenant ID is useful for multi-tenant applications, allowing operations
// to be associated with a specific tenant for data isolation and authorization.
//
// Parameters:
//   - ctx: The parent context from which the new context will be derived.
//   - tenantID: The ID of the tenant to add to the context.
//
// Returns:
//   - context.Context: The newly created context with the tenant ID added.
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// GetTenantID gets the tenant ID from the context.
// This function retrieves the tenant ID that was previously set with WithTenantID.
//
// Parameters:
//   - ctx: The context from which to retrieve the tenant ID.
//
// Returns:
//   - string: The tenant ID, or an empty string if no tenant ID was set.
func GetTenantID(ctx context.Context) string {
	if id, ok := ctx.Value(TenantIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextInfo returns a string with all the context information.
// This function creates a formatted string containing all the metadata stored in the context,
// including request ID, trace ID, user ID, tenant ID, operation name, correlation ID,
// service name, and deadline (if set). This is useful for logging and debugging.
//
// Parameters:
//   - ctx: The context from which to retrieve the information.
//
// Returns:
//   - string: A formatted string containing all the context information.
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
