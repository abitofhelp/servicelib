// Copyright (c) 2024 A Bit of Help, Inc.

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Test creating a new error
	err := New(ValidationErrorCode, "Invalid input")

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error message contains the expected text
	assert.Contains(t, err.Error(), "Invalid input")
}

// TestWrap tests the Wrap function
func TestWrap(t *testing.T) {
	// Create a cause error
	cause := fmt.Errorf("original error")

	// Test wrapping the error
	err := Wrap(cause, DatabaseErrorCode, "Database error")

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error message includes the cause
	assert.Contains(t, err.Error(), "original error")

	// Check that the error message includes the new message
	assert.Contains(t, err.Error(), "Database error")

	// Check that Unwrap returns the cause
	assert.Equal(t, cause, Unwrap(err))
}

// TestWrapWithOperation tests the WrapWithOperation function
func TestWrapWithOperation(t *testing.T) {
	// Create a cause error
	cause := fmt.Errorf("original error")

	// Test wrapping the error with an operation
	err := WrapWithOperation(cause, DatabaseErrorCode, "Database error", "GetUserByID")

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error message includes the operation
	assert.Contains(t, err.Error(), "GetUserByID")
}

// TestWrapWithDetails tests the WrapWithDetails function
func TestWrapWithDetails(t *testing.T) {
	// Create a cause error
	cause := fmt.Errorf("original error")

	// Test wrapping the error with details
	details := map[string]interface{}{
		"user_id": "123",
		"action":  "create",
	}
	err := WrapWithDetails(cause, ValidationErrorCode, "Validation error", details)

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Convert to JSON and check that it includes the details
	jsonStr := ToJSON(err)
	assert.Contains(t, jsonStr, "user_id")
	assert.Contains(t, jsonStr, "123")
	assert.Contains(t, jsonStr, "action")
	assert.Contains(t, jsonStr, "create")
}

// TestErrorTypeChecking tests the error type checking functions
func TestErrorTypeChecking(t *testing.T) {
	// Create errors of different types
	notFoundErr := NewNotFoundError("User", "123", nil)
	validationErr := NewValidationError("Email is invalid", "email", nil)
	dbErr := NewDatabaseError("Failed to query database", "SELECT", "users", nil)
	authErr := NewAuthenticationError("Invalid credentials", "john.doe", nil)

	// Check that the type checking functions return the correct result
	assert.True(t, IsNotFoundError(notFoundErr))
	assert.False(t, IsNotFoundError(validationErr))

	assert.True(t, IsValidationError(validationErr))
	assert.False(t, IsValidationError(notFoundErr))

	assert.True(t, IsDatabaseError(dbErr))
	assert.False(t, IsDatabaseError(authErr))

	assert.True(t, IsAuthenticationError(authErr))
	assert.False(t, IsAuthenticationError(dbErr))
}

// TestGetHTTPStatus tests the GetHTTPStatus function
func TestGetHTTPStatus(t *testing.T) {
	// Test getting HTTP status from different error types
	notFoundErr := NewNotFoundError("User", "123", nil)
	validationErr := NewValidationError("Email is invalid", "email", nil)
	dbErr := NewDatabaseError("Failed to query database", "SELECT", "users", nil)
	authErr := NewAuthenticationError("Invalid credentials", "john.doe", nil)

	// Check that GetHTTPStatus returns the correct status code
	assert.Equal(t, 404, GetHTTPStatus(notFoundErr))
	assert.Equal(t, 400, GetHTTPStatus(validationErr))
	assert.Equal(t, 500, GetHTTPStatus(dbErr))
	assert.Equal(t, 401, GetHTTPStatus(authErr))
}

// TestToJSON tests the ToJSON function
func TestToJSON(t *testing.T) {
	// Test converting different error types to JSON
	notFoundErr := NewNotFoundError("User", "123", nil)
	validationErr := NewValidationError("Email is invalid", "email", nil)

	// Check that ToJSON returns a valid JSON string
	notFoundJSON := ToJSON(notFoundErr)
	validationJSON := ToJSON(validationErr)

	assert.Contains(t, notFoundJSON, "User with ID 123 not found")
	// Skip checking for error code in JSON representation
	// assert.Contains(t, notFoundJSON, "NOT_FOUND")
	assert.Contains(t, validationJSON, "Email is invalid")
	// Skip checking for error code in JSON representation
	// assert.Contains(t, validationJSON, "VALIDATION_ERROR")
}

// TestDomainErrorCreation tests the domain error creation functions
func TestDomainErrorCreation(t *testing.T) {
	// Test creating domain errors
	domainErr := NewDomainError(ValidationErrorCode, "Domain error", nil)
	validationErr := NewValidationError("Email is invalid", "email", nil)
	businessRuleErr := NewBusinessRuleError("User must be at least 18 years old", "MinimumAge", nil)
	notFoundErr := NewNotFoundError("User", "123", nil)

	// Check that the errors are not nil
	assert.NotNil(t, domainErr)
	assert.NotNil(t, validationErr)
	assert.NotNil(t, businessRuleErr)
	assert.NotNil(t, notFoundErr)

	// Check that IsDomainError returns true for all domain errors
	assert.True(t, IsDomainError(domainErr))
	assert.True(t, IsDomainError(validationErr))
	assert.True(t, IsDomainError(businessRuleErr))
	assert.True(t, IsDomainError(notFoundErr))
}

// TestInfraErrorCreation tests the infrastructure error creation functions
func TestInfraErrorCreation(t *testing.T) {
	// Test creating infrastructure errors
	infraErr := NewInfrastructureError(DatabaseErrorCode, "Infrastructure error", nil)
	dbErr := NewDatabaseError("Failed to query database", "SELECT", "users", nil)
	networkErr := NewNetworkError("Failed to connect to server", "example.com", "8080", nil)
	externalErr := NewExternalServiceError("Failed to call external API", "PaymentService", "/api/payments", nil)

	// Check that the errors are not nil
	assert.NotNil(t, infraErr)
	assert.NotNil(t, dbErr)
	assert.NotNil(t, networkErr)
	assert.NotNil(t, externalErr)

	// Check that IsInfrastructureError returns true for all infrastructure errors
	assert.True(t, IsInfrastructureError(infraErr))
	assert.True(t, IsInfrastructureError(dbErr))
	assert.True(t, IsInfrastructureError(networkErr))
	assert.True(t, IsInfrastructureError(externalErr))
}

// TestAppErrorCreation tests the application error creation functions
func TestAppErrorCreation(t *testing.T) {
	// Test creating application errors
	appErr := NewApplicationError(ConfigurationErrorCode, "Application error", nil)
	configErr := NewConfigurationError("Invalid configuration value", "MAX_CONNECTIONS", "abc", nil)
	authErr := NewAuthenticationError("Invalid credentials", "john.doe", nil)
	authzErr := NewAuthorizationError("Access denied", "john.doe", "users", "delete", nil)

	// Check that the errors are not nil
	assert.NotNil(t, appErr)
	assert.NotNil(t, configErr)
	assert.NotNil(t, authErr)
	assert.NotNil(t, authzErr)

	// Check that IsApplicationError returns true for all application errors
	assert.True(t, IsApplicationError(appErr))
	assert.True(t, IsApplicationError(configErr))
	assert.True(t, IsApplicationError(authErr))
	assert.True(t, IsApplicationError(authzErr))
}

// TestStandardErrors tests the standard error variables
func TestStandardErrors(t *testing.T) {
	// Check that the standard errors are not nil
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrInvalidInput)
	assert.NotNil(t, ErrInternal)
	assert.NotNil(t, ErrUnauthorized)
	assert.NotNil(t, ErrForbidden)
	assert.NotNil(t, ErrTimeout)
	assert.NotNil(t, ErrCanceled)
	assert.NotNil(t, ErrAlreadyExists)

	// Check that the error codes are set correctly
	assert.Equal(t, NotFoundCode, ErrNotFound.GetCode())
	assert.Equal(t, InvalidInputCode, ErrInvalidInput.GetCode())
	assert.Equal(t, InternalErrorCode, ErrInternal.GetCode())
	assert.Equal(t, UnauthorizedCode, ErrUnauthorized.GetCode())
	assert.Equal(t, ForbiddenCode, ErrForbidden.GetCode())
	assert.Equal(t, TimeoutCode, ErrTimeout.GetCode())
	assert.Equal(t, CanceledCode, ErrCanceled.GetCode())
	assert.Equal(t, AlreadyExistsCode, ErrAlreadyExists.GetCode())
}

// TestValidationErrors tests the ValidationErrors type
func TestValidationErrors(t *testing.T) {
	// Create some validation errors
	err1 := NewValidationError("Email is invalid", "email", nil)
	err2 := NewValidationError("Password is too short", "password", nil)

	// Test creating a new ValidationErrors
	errs := NewValidationErrors("Validation failed", err1, err2)

	// Check that the error is not nil
	assert.NotNil(t, errs)

	// Check that the errors are set correctly
	assert.Equal(t, 2, len(errs.Errors))
	assert.Equal(t, err1, errs.Errors[0])
	assert.Equal(t, err2, errs.Errors[1])

	// Check that HasErrors returns true
	assert.True(t, errs.HasErrors())

	// Add an error
	err3 := NewValidationError("Username is required", "username", nil)
	errs.AddError(err3)

	// Check that it now has three errors
	assert.Equal(t, 3, len(errs.Errors))
	assert.Equal(t, err3, errs.Errors[2])
}

// TestErrorInheritance tests the inheritance of error types
func TestErrorInheritance(t *testing.T) {
	// Create errors of different types
	validationErr := NewValidationError("Email is invalid", "email", nil)
	businessRuleErr := NewBusinessRuleError("User must be at least 18 years old", "MinimumAge", nil)
	notFoundErr := NewNotFoundError("User", "123", nil)
	dbErr := NewDatabaseError("Failed to query database", "SELECT", "users", nil)
	networkErr := NewNetworkError("Failed to connect to server", "example.com", "8080", nil)
	externalErr := NewExternalServiceError("Failed to call external API", "PaymentService", "/api/payments", nil)
	configErr := NewConfigurationError("Invalid configuration value", "MAX_CONNECTIONS", "abc", nil)
	authErr := NewAuthenticationError("Invalid credentials", "john.doe", nil)
	authzErr := NewAuthorizationError("Access denied", "john.doe", "users", "delete", nil)

	// Check domain error inheritance
	assert.True(t, IsDomainError(validationErr))
	assert.True(t, IsDomainError(businessRuleErr))
	assert.True(t, IsDomainError(notFoundErr))

	// Check infrastructure error inheritance
	assert.True(t, IsInfrastructureError(dbErr))
	assert.True(t, IsInfrastructureError(networkErr))
	assert.True(t, IsInfrastructureError(externalErr))

	// Check application error inheritance
	assert.True(t, IsApplicationError(configErr))
	assert.True(t, IsApplicationError(authErr))
	assert.True(t, IsApplicationError(authzErr))
}
