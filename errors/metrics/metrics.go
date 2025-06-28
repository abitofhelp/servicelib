// Copyright (c) 2025 A Bit of Help, Inc.

// Package metrics provides metrics integration for the error handling system.
package metrics

import (
	"context"

	"github.com/abitofhelp/servicelib/errors/core"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	// errorCounter is a counter for tracking errors by code and type
	errorCounter metric.Int64Counter
)

// Initialize sets up the error metrics
func Initialize(meter metric.Meter) error {
	var err error

	// Create error counter
	errorCounter, err = meter.Int64Counter(
		"errors.count",
		metric.WithDescription("Number of errors by code and type"),
	)
	if err != nil {
		return err
	}

	return nil
}

// RecordError records an error metric with the error code and type
func RecordError(ctx context.Context, err error) {
	if err == nil || errorCounter == nil {
		return
	}

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

	// Record the error
	errorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
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

// RecordErrorWithAttributes records an error metric with additional attributes
func RecordErrorWithAttributes(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	if err == nil || errorCounter == nil {
		return
	}

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

	// Record the error
	errorCounter.Add(ctx, 1, metric.WithAttributes(baseAttrs...))
}
