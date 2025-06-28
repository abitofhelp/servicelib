// Copyright (c) 2025 A Bit of Help, Inc.

package core

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	// Create a base error
	baseErr := errors.New("base error")

	// Create a wrapped error
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)

	// Test Is function
	assert.True(t, Is(wrappedErr, baseErr), "Is should return true for wrapped error")
	assert.False(t, Is(wrappedErr, errors.New("different error")), "Is should return false for different error")
}

// Custom error type for testing
type customError struct {
	msg string
}

// Implement the error interface
func (e *customError) Error() string {
	return e.msg
}

func TestAs(t *testing.T) {
	// Create an instance of the custom error
	custErr := &customError{msg: "custom error"}

	// Wrap the custom error
	wrappedErr := fmt.Errorf("wrapped: %w", custErr)

	// Test As function
	var target *customError
	assert.True(t, As(wrappedErr, &target), "As should return true for wrapped custom error")
	assert.Equal(t, "custom error", target.Error(), "Target should be set to the custom error")

	var wrongTarget *BaseError
	assert.False(t, As(wrappedErr, &wrongTarget), "As should return false for wrong target type")
}

func TestUnwrap(t *testing.T) {
	// Create a base error
	baseErr := errors.New("base error")

	// Create a wrapped error
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)

	// Test Unwrap function
	unwrappedErr := Unwrap(wrappedErr)
	assert.Equal(t, baseErr, unwrappedErr, "Unwrap should return the base error")

	// Test Unwrap with non-wrapped error
	nonWrappedErr := errors.New("non-wrapped error")
	assert.Nil(t, Unwrap(nonWrappedErr), "Unwrap should return nil for non-wrapped error")
}

func TestGetCode(t *testing.T) {
	// Create a BaseError with a code
	baseErr := &BaseError{
		Code:    InvalidInputCode,
		Message: "invalid input",
	}

	// Test GetCode function
	assert.Equal(t, InvalidInputCode, GetCode(baseErr), "GetCode should return the error code")

	// Test GetCode with standard error
	stdErr := errors.New("standard error")
	assert.Equal(t, InternalErrorCode, GetCode(stdErr), "GetCode should return InternalErrorCode for standard error")
}

func TestGetHTTPStatusFromError(t *testing.T) {
	// Create errors with different codes
	notFoundErr := &BaseError{
		Code:    NotFoundCode,
		Message: "not found",
	}

	badRequestErr := &BaseError{
		Code:    InvalidInputCode,
		Message: "invalid input",
	}

	internalErr := &BaseError{
		Code:    InternalErrorCode,
		Message: "internal error",
	}

	// Test GetHTTPStatusFromError function
	assert.Equal(t, http.StatusNotFound, GetHTTPStatusFromError(notFoundErr), "Should return 404 for NotFoundCode")
	assert.Equal(t, http.StatusBadRequest, GetHTTPStatusFromError(badRequestErr), "Should return 400 for InvalidInputCode")
	assert.Equal(t, http.StatusInternalServerError, GetHTTPStatusFromError(internalErr), "Should return 500 for InternalErrorCode")

	// Test with standard error
	stdErr := errors.New("standard error")
	assert.Equal(t, http.StatusInternalServerError, GetHTTPStatusFromError(stdErr), "Should return 500 for standard error")
}

func TestToJSON(t *testing.T) {
	// Create a BaseError with details
	baseErr := &BaseError{
		Code:    InvalidInputCode,
		Message: "invalid input",
		Details: map[string]interface{}{
			"field": "username",
			"value": "invalid",
		},
	}

	// Test ToJSON function
	jsonStr := ToJSON(baseErr)
	assert.Contains(t, jsonStr, `"code":"INVALID_INPUT"`, "JSON should contain the error code")
	assert.Contains(t, jsonStr, `"message":"invalid input"`, "JSON should contain the error message")
	assert.Contains(t, jsonStr, `"details":`, "JSON should contain the details")
	assert.Contains(t, jsonStr, `"field":"username"`, "JSON should contain the field detail")
	assert.Contains(t, jsonStr, `"value":"invalid"`, "JSON should contain the value detail")

	// Test with standard error
	stdErr := errors.New("standard error")
	stdJSON := ToJSON(stdErr)
	assert.Contains(t, stdJSON, `"code":"INTERNAL_ERROR"`, "JSON should contain the default error code")
	assert.Contains(t, stdJSON, `"message":"standard error"`, "JSON should contain the error message")
}

func TestWrapWithOperation(t *testing.T) {
	// Create a base error
	baseErr := errors.New("base error")

	// Wrap with operation
	wrappedErr := WrapWithOperation(baseErr, "TestOperation")

	// Test that the wrapped error is a ContextualError
	var contextualErr *ContextualError
	assert.True(t, As(wrappedErr, &contextualErr), "Wrapped error should be a ContextualError")

	// Test that the operation is set
	assert.Equal(t, "TestOperation", contextualErr.Context.Operation, "Operation should be set")

	// Test that the original error is set
	assert.Equal(t, baseErr, contextualErr.Original, "Original error should be set to the original error")

	// Test wrapping a BaseError
	originalBaseErr := &BaseError{
		Code:      InvalidInputCode,
		Message:   "invalid input",
		Operation: "OriginalOperation",
	}

	wrappedBaseErr := WrapWithOperation(originalBaseErr, "NewOperation")
	var wrappedContextualErr *ContextualError
	assert.True(t, As(wrappedBaseErr, &wrappedContextualErr), "Wrapped error should be a ContextualError")
	assert.Equal(t, "NewOperation", wrappedContextualErr.Context.Operation, "Operation should be updated")
	assert.Equal(t, InvalidInputCode, wrappedContextualErr.Context.Code, "Code should be preserved")
}

func TestWithDetails(t *testing.T) {
	// Create a base error
	baseErr := errors.New("base error")

	// Add details
	details := map[string]interface{}{
		"field": "username",
		"value": "invalid",
	}

	detailedErr := WithDetails(baseErr, details)

	// Test that the detailed error is a ContextualError
	var contextualErr *ContextualError
	assert.True(t, As(detailedErr, &contextualErr), "Detailed error should be a ContextualError")

	// Test that the details are set
	assert.Equal(t, details, contextualErr.Context.Details, "Details should be set")

	// Test that the original error is set
	assert.Equal(t, baseErr, contextualErr.Original, "Original error should be set to the original error")

	// Test adding details to a BaseError
	originalBaseErr := &BaseError{
		Code:    InvalidInputCode,
		Message: "invalid input",
		Details: map[string]interface{}{
			"original": "detail",
		},
	}

	newDetails := map[string]interface{}{
		"new": "detail",
	}

	detailedBaseErr := WithDetails(originalBaseErr, newDetails)
	var detailedContextualErr *ContextualError
	assert.True(t, As(detailedBaseErr, &detailedContextualErr), "Detailed error should be a ContextualError")

	// Test that the new details are merged with the original details
	assert.Equal(t, "detail", detailedContextualErr.Context.Details["original"], "Original details should be preserved")
	assert.Equal(t, "detail", detailedContextualErr.Context.Details["new"], "New details should be added")
}
