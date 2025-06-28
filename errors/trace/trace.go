// Copyright (c) 2025 A Bit of Help, Inc.

// Package trace provides tracing integration for the error handling system.
package trace

import (
	"context"
	"fmt"

	"github.com/abitofhelp/servicelib/errors/core"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// AddErrorToSpan adds error information to the current span.
// It extracts error details like code, operation, and details if available.
func AddErrorToSpan(ctx context.Context, err error) {
	if err == nil {
		return
	}

	// Get the current span from context
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	// Set span status to error
	span.SetStatus(codes.Error, err.Error())

	// Create base attributes
	attrs := []attribute.KeyValue{
		attribute.String("error.message", err.Error()),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		attrs = append(attrs, attribute.String("error.code", string(e.GetCode())))
	}

	// Add error type
	attrs = append(attrs, attribute.String("error.type", getErrorType(err)))

	// Add operation if available
	if e, ok := err.(interface{ GetOperation() string }); ok && e.GetOperation() != "" {
		attrs = append(attrs, attribute.String("error.operation", e.GetOperation()))
	}

	// Add source location if available
	if e, ok := err.(interface {
		GetSource() string
		GetLine() int
	}); ok && e.GetSource() != "" {
		attrs = append(attrs, attribute.String("error.source", e.GetSource()))
		if e.GetLine() > 0 {
			attrs = append(attrs, attribute.Int("error.line", e.GetLine()))
		}
	}

	// Add details if available
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		details := e.GetDetails()
		if details != nil {
			for k, v := range details {
				// Convert details to string to avoid type issues
				attrs = append(attrs, attribute.String("error.detail."+k, fmt.Sprintf("%v", v)))
			}
		}
	}

	// Add HTTP status if available
	if e, ok := err.(interface{ GetHTTPStatus() int }); ok && e.GetHTTPStatus() != 0 {
		attrs = append(attrs, attribute.Int("error.http_status", e.GetHTTPStatus()))
	}

	// Add attributes to span
	span.SetAttributes(attrs...)

	// Record the error event
	span.RecordError(err, trace.WithAttributes(attrs...))
}

// getErrorType returns the type of the error as a string
func getErrorType(err error) string {
	// Check for domain errors
	if e, ok := err.(interface{ IsDomainError() bool }); ok && e.IsDomainError() {
		if e, ok := err.(interface{ IsValidationError() bool }); ok && e.IsValidationError() {
			return "validation_error"
		}
		if e, ok := err.(interface{ IsBusinessRuleError() bool }); ok && e.IsBusinessRuleError() {
			return "business_rule_error"
		}
		if e, ok := err.(interface{ IsNotFoundError() bool }); ok && e.IsNotFoundError() {
			return "not_found_error"
		}
		return "domain_error"
	}

	// Check for infrastructure errors
	if e, ok := err.(interface{ IsInfrastructureError() bool }); ok && e.IsInfrastructureError() {
		if e, ok := err.(interface{ IsDatabaseError() bool }); ok && e.IsDatabaseError() {
			return "database_error"
		}
		if e, ok := err.(interface{ IsNetworkError() bool }); ok && e.IsNetworkError() {
			return "network_error"
		}
		if e, ok := err.(interface{ IsExternalServiceError() bool }); ok && e.IsExternalServiceError() {
			return "external_service_error"
		}
		return "infrastructure_error"
	}

	// Check for application errors
	if e, ok := err.(interface{ IsApplicationError() bool }); ok && e.IsApplicationError() {
		if e, ok := err.(interface{ IsConfigurationError() bool }); ok && e.IsConfigurationError() {
			return "configuration_error"
		}
		if e, ok := err.(interface{ IsAuthenticationError() bool }); ok && e.IsAuthenticationError() {
			return "authentication_error"
		}
		if e, ok := err.(interface{ IsAuthorizationError() bool }); ok && e.IsAuthorizationError() {
			return "authorization_error"
		}
		return "application_error"
	}

	// Default to generic error
	return "error"
}

// AddErrorToSpanWithAttributes adds error information to the current span with additional attributes.
func AddErrorToSpanWithAttributes(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	if err == nil {
		return
	}

	// Get the current span from context
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	// Set span status to error
	span.SetStatus(codes.Error, err.Error())

	// Create base attributes
	baseAttrs := []attribute.KeyValue{
		attribute.String("error.message", err.Error()),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		baseAttrs = append(baseAttrs, attribute.String("error.code", string(e.GetCode())))
	}

	// Add error type
	baseAttrs = append(baseAttrs, attribute.String("error.type", getErrorType(err)))

	// Add operation if available
	if e, ok := err.(interface{ GetOperation() string }); ok && e.GetOperation() != "" {
		baseAttrs = append(baseAttrs, attribute.String("error.operation", e.GetOperation()))
	}

	// Add custom attributes
	baseAttrs = append(baseAttrs, attrs...)

	// Add attributes to span
	span.SetAttributes(baseAttrs...)

	// Record the error event
	span.RecordError(err, trace.WithAttributes(baseAttrs...))
}
