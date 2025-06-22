// Copyright (c) 2025 A Bit of Help, Inc.

package domain

import (
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

func TestNewDomainError(t *testing.T) {
	// Test creating a new DomainError
	err := NewDomainError(core.ValidationErrorCode, "Invalid input", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.ValidationErrorCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Invalid input", err.BaseError.GetMessage())
	
	// Check that IsDomainError returns true
	assert.True(t, err.IsDomainError())
}

func TestNewValidationError(t *testing.T) {
	// Test creating a new ValidationError
	err := NewValidationError("Email is invalid", "email", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.ValidationErrorCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Email is invalid", err.BaseError.GetMessage())
	
	// Check that the field is set correctly
	assert.Equal(t, "email", err.Field)
	
	// Check that IsValidationError returns true
	assert.True(t, err.IsValidationError())
	
	// Check that IsDomainError returns true (inheritance)
	assert.True(t, err.IsDomainError())
}

func TestNewValidationErrors(t *testing.T) {
	// Create some validation errors
	err1 := NewValidationError("Email is invalid", "email", nil)
	err2 := NewValidationError("Password is too short", "password", nil)
	
	// Test creating a new ValidationErrors
	errs := NewValidationErrors("Validation failed", err1, err2)
	
	// Check that the error is not nil
	assert.NotNil(t, errs)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.ValidationErrorCode, errs.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Validation failed", errs.BaseError.GetMessage())
	
	// Check that the errors are set correctly
	assert.Equal(t, 2, len(errs.Errors))
	assert.Equal(t, err1, errs.Errors[0])
	assert.Equal(t, err2, errs.Errors[1])
	
	// Check that HasErrors returns true
	assert.True(t, errs.HasErrors())
}

func TestValidationErrorsAddError(t *testing.T) {
	// Create a ValidationErrors
	errs := NewValidationErrors("Validation failed")
	
	// Check that it has no errors initially
	assert.Equal(t, 0, len(errs.Errors))
	assert.False(t, errs.HasErrors())
	
	// Add an error
	err := NewValidationError("Email is invalid", "email", nil)
	errs.AddError(err)
	
	// Check that it now has one error
	assert.Equal(t, 1, len(errs.Errors))
	assert.Equal(t, err, errs.Errors[0])
	assert.True(t, errs.HasErrors())
}

func TestNewBusinessRuleError(t *testing.T) {
	// Test creating a new BusinessRuleError
	err := NewBusinessRuleError("User must be at least 18 years old", "MinimumAge", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.BusinessRuleViolationCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "User must be at least 18 years old", err.BaseError.GetMessage())
	
	// Check that the rule is set correctly
	assert.Equal(t, "MinimumAge", err.Rule)
	
	// Check that IsBusinessRuleError returns true
	assert.True(t, err.IsBusinessRuleError())
	
	// Check that IsDomainError returns true (inheritance)
	assert.True(t, err.IsDomainError())
}

func TestNewNotFoundError(t *testing.T) {
	// Test creating a new NotFoundError
	err := NewNotFoundError("User", "123", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.NotFoundCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "User with ID 123 not found", err.BaseError.GetMessage())
	
	// Check that the resource type and ID are set correctly
	assert.Equal(t, "User", err.ResourceType)
	assert.Equal(t, "123", err.ResourceID)
	
	// Check that IsNotFoundError returns true
	assert.True(t, err.IsNotFoundError())
	
	// Check that IsDomainError returns true (inheritance)
	assert.True(t, err.IsDomainError())
}

func TestErrorsAs(t *testing.T) {
	// Create errors of different types
	domainErr := NewDomainError(core.ValidationErrorCode, "Domain error", nil)
	validationErr := NewValidationError("Email is invalid", "email", nil)
	businessRuleErr := NewBusinessRuleError("User must be at least 18 years old", "MinimumAge", nil)
	notFoundErr := NewNotFoundError("User", "123", nil)
	
	// Test errors.As with DomainError
	var de *DomainError
	assert.True(t, errors.As(domainErr, &de))
	
	// Test errors.As with ValidationError
	var ve *ValidationError
	assert.True(t, errors.As(validationErr, &ve))
	
	// Test errors.As with BusinessRuleError
	var bre *BusinessRuleError
	assert.True(t, errors.As(businessRuleErr, &bre))
	
	// Test errors.As with NotFoundError
	var nfe *NotFoundError
	assert.True(t, errors.As(notFoundErr, &nfe))
	
	// Test inheritance (ValidationError is also a DomainError)
	assert.True(t, errors.As(validationErr, &de))
}