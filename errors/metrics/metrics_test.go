// Copyright (c) 2025 A Bit of Help, Inc.

package metrics

import (
	"context"
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/errors/app"
	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/abitofhelp/servicelib/errors/domain"
	"github.com/abitofhelp/servicelib/errors/infra"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

// TestGetErrorType tests the getErrorType function
func TestGetErrorType(t *testing.T) {
	// Test with standard error
	stdErr := errors.New("standard error")
	assert.Equal(t, "error", getErrorType(stdErr))

	// Test with domain error
	domainErr := domain.NewDomainError(core.ValidationErrorCode, "domain error", nil)
	assert.Equal(t, "domain_error", getErrorType(domainErr))

	// Test with validation error
	validationErr := domain.NewValidationError("invalid input", "field", nil)
	assert.Equal(t, "validation_error", getErrorType(validationErr))

	// Test with business rule error
	businessRuleErr := domain.NewBusinessRuleError("rule violated", "rule", nil)
	assert.Equal(t, "business_rule_error", getErrorType(businessRuleErr))

	// Test with not found error
	notFoundErr := domain.NewNotFoundError("User", "123", nil)
	assert.Equal(t, "not_found_error", getErrorType(notFoundErr))

	// Test with infrastructure error
	infraErr := infra.NewInfrastructureError(core.DatabaseErrorCode, "infra error", nil)
	assert.Equal(t, "infrastructure_error", getErrorType(infraErr))

	// Test with database error
	dbErr := infra.NewDatabaseError("db error", "SELECT", "users", nil)
	assert.Equal(t, "database_error", getErrorType(dbErr))

	// Test with network error
	networkErr := infra.NewNetworkError("network error", "localhost", "8080", nil)
	assert.Equal(t, "network_error", getErrorType(networkErr))

	// Test with external service error
	externalErr := infra.NewExternalServiceError("external error", "service", "/api", nil)
	assert.Equal(t, "external_service_error", getErrorType(externalErr))

	// Test with application error
	appErr := app.NewApplicationError(core.ConfigurationErrorCode, "app error", nil)
	assert.Equal(t, "application_error", getErrorType(appErr))

	// Test with configuration error
	configErr := app.NewConfigurationError("config error", "key", "value", nil)
	assert.Equal(t, "configuration_error", getErrorType(configErr))

	// Test with authentication error
	authErr := app.NewAuthenticationError("auth error", "user", nil)
	assert.Equal(t, "authentication_error", getErrorType(authErr))

	// Test with authorization error
	authzErr := app.NewAuthorizationError("authz error", "user", "resource", "action", nil)
	assert.Equal(t, "authorization_error", getErrorType(authzErr))
}

// TestRecordError_NilError tests that RecordError handles nil errors gracefully
func TestRecordError_NilError(t *testing.T) {
	// Save the original counter
	originalCounter := errorCounter

	// Set the counter to nil to simulate uninitialized state
	errorCounter = nil

	// This should not panic
	RecordError(context.Background(), nil)

	// Restore the original counter
	errorCounter = originalCounter
}

// TestRecordError_NilCounter tests that RecordError handles nil counter gracefully
func TestRecordError_NilCounter(t *testing.T) {
	// Save the original counter
	originalCounter := errorCounter

	// Set the counter to nil to simulate uninitialized state
	errorCounter = nil

	// This should not panic
	RecordError(context.Background(), errors.New("test error"))

	// Restore the original counter
	errorCounter = originalCounter
}

// TestRecordErrorWithAttributes_NilError tests that RecordErrorWithAttributes handles nil errors gracefully
func TestRecordErrorWithAttributes_NilError(t *testing.T) {
	// Save the original counter
	originalCounter := errorCounter

	// Set the counter to nil to simulate uninitialized state
	errorCounter = nil

	// This should not panic
	RecordErrorWithAttributes(context.Background(), nil)

	// Restore the original counter
	errorCounter = originalCounter
}

// TestRecordErrorWithAttributes_NilCounter tests that RecordErrorWithAttributes handles nil counter gracefully
func TestRecordErrorWithAttributes_NilCounter(t *testing.T) {
	// Save the original counter
	originalCounter := errorCounter

	// Set the counter to nil to simulate uninitialized state
	errorCounter = nil

	// This should not panic
	attrs := []attribute.KeyValue{
		attribute.String("test", "value"),
	}
	RecordErrorWithAttributes(context.Background(), errors.New("test error"), attrs...)

	// Restore the original counter
	errorCounter = originalCounter
}
