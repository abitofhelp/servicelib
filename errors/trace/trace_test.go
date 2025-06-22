// Copyright (c) 2025 A Bit of Help, Inc.

package trace

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

// TestAddErrorToSpan_NilError tests that AddErrorToSpan handles nil errors gracefully
func TestAddErrorToSpan_NilError(t *testing.T) {
	// Create a context with no span
	ctx := context.Background()

	// This should not panic
	AddErrorToSpan(ctx, nil)
}

// TestAddErrorToSpan_NoSpan tests that AddErrorToSpan handles contexts with no span gracefully
func TestAddErrorToSpan_NoSpan(t *testing.T) {
	// Create a context with no span
	ctx := context.Background()

	// This should not panic
	AddErrorToSpan(ctx, errors.New("test error"))
}

// TestAddErrorToSpanWithAttributes_NilError tests that AddErrorToSpanWithAttributes handles nil errors gracefully
func TestAddErrorToSpanWithAttributes_NilError(t *testing.T) {
	// Create a context with no span
	ctx := context.Background()

	// This should not panic
	AddErrorToSpanWithAttributes(ctx, nil)
}

// TestAddErrorToSpanWithAttributes_NoSpan tests that AddErrorToSpanWithAttributes handles contexts with no span gracefully
func TestAddErrorToSpanWithAttributes_NoSpan(t *testing.T) {
	// Create a context with no span
	ctx := context.Background()

	// This should not panic
	AddErrorToSpanWithAttributes(ctx, errors.New("test error"))
}

// TestAddErrorToSpanWithAttributes_WithAttributes tests that AddErrorToSpanWithAttributes handles attributes correctly
func TestAddErrorToSpanWithAttributes_WithAttributes(t *testing.T) {
	// Create a context with no span
	ctx := context.Background()

	// Create custom attributes
	customAttrs := []attribute.KeyValue{
		attribute.String("custom.key", "custom.value"),
		attribute.Int("custom.count", 42),
	}

	// This should not panic
	AddErrorToSpanWithAttributes(ctx, errors.New("test error"), customAttrs...)
}
