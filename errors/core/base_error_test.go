// Copyright (c) 2025 A Bit of Help, Inc.

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseError(t *testing.T) {
	// Test creating a new BaseError
	err := NewBaseError(ValidationErrorCode, "Invalid input", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, ValidationErrorCode, err.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Invalid input", err.GetMessage())
	
	// Check that the cause is nil
	assert.Nil(t, err.GetCause())
	
	// Check that the source and line are set
	assert.NotEmpty(t, err.GetSource())
	assert.Greater(t, err.GetLine(), 0)
}

func TestBaseErrorWithCause(t *testing.T) {
	// Create a cause error
	cause := fmt.Errorf("original error")
	
	// Create a BaseError with the cause
	err := NewBaseError(DatabaseErrorCode, "Database error", cause)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the cause is set correctly
	assert.Equal(t, cause, err.GetCause())
	
	// Check that the error message includes the cause
	assert.Contains(t, err.Error(), "original error")
}

func TestBaseErrorWithOperation(t *testing.T) {
	// Create a BaseError
	err := NewBaseError(NotFoundCode, "User not found", nil)
	
	// Add an operation
	err = err.WithOperation("GetUserByID")
	
	// Check that the operation is set correctly
	assert.Equal(t, "GetUserByID", err.GetOperation())
	
	// Check that the error message includes the operation
	assert.Contains(t, err.Error(), "GetUserByID")
}

func TestBaseErrorWithDetails(t *testing.T) {
	// Create a BaseError
	err := NewBaseError(ValidationErrorCode, "Invalid input", nil)
	
	// Add details
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid-email",
	}
	err = err.WithDetails(details)
	
	// Check that the details are set correctly
	assert.Equal(t, details, err.GetDetails())
}

func TestBaseErrorMarshalJSON(t *testing.T) {
	// Create a BaseError with all fields set
	cause := fmt.Errorf("original error")
	err := NewBaseError(ValidationErrorCode, "Invalid input", cause)
	err = err.WithOperation("ValidateUser")
	err = err.WithDetails(map[string]interface{}{
		"field": "email",
		"value": "invalid-email",
	})
	
	// Marshal to JSON
	jsonBytes, jsonErr := json.Marshal(err)
	assert.NoError(t, jsonErr)
	
	// Unmarshal back to a map
	var result map[string]interface{}
	jsonErr = json.Unmarshal(jsonBytes, &result)
	assert.NoError(t, jsonErr)
	
	// Check that all fields are present in the JSON
	assert.Equal(t, "Invalid input: original error", result["message"])
	assert.Equal(t, "VALIDATION_ERROR", result["code"])
	assert.Equal(t, "ValidateUser", result["operation"])
	assert.NotNil(t, result["details"])
	assert.NotEmpty(t, result["source"])
	assert.NotNil(t, result["line"])
	assert.Equal(t, "original error", result["cause"])
}

func TestBaseErrorIs(t *testing.T) {
	// Create a BaseError
	err1 := NewBaseError(NotFoundCode, "User not found", nil)
	
	// Create another BaseError with the same code
	err2 := NewBaseError(NotFoundCode, "Resource not found", nil)
	
	// Create a BaseError with a different code
	err3 := NewBaseError(ValidationErrorCode, "Invalid input", nil)
	
	// Check that err1 is err2 (same code)
	assert.True(t, err1.Is(err2))
	
	// Check that err1 is not err3 (different code)
	assert.False(t, err1.Is(err3))
	
	// Check with standard errors.Is
	assert.True(t, errors.Is(err1, err2))
	assert.False(t, errors.Is(err1, err3))
}

func TestBaseErrorAs(t *testing.T) {
	// Create a BaseError
	err := NewBaseError(NotFoundCode, "User not found", nil)
	
	// Check with standard errors.As
	var baseErr *BaseError
	assert.True(t, errors.As(err, &baseErr))
	assert.Equal(t, NotFoundCode, baseErr.GetCode())
}

func TestBaseErrorUnwrap(t *testing.T) {
	// Create a cause error
	cause := fmt.Errorf("original error")
	
	// Create a BaseError with the cause
	err := NewBaseError(DatabaseErrorCode, "Database error", cause)
	
	// Check that Unwrap returns the cause
	assert.Equal(t, cause, err.Unwrap())
	
	// Check with standard errors.Unwrap
	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestGetHTTPStatus(t *testing.T) {
	// Test all error codes
	testCases := []struct {
		code     ErrorCode
		expected int
	}{
		{NotFoundCode, 404},
		{InvalidInputCode, 400},
		{DatabaseErrorCode, 500},
		{InternalErrorCode, 500},
		{TimeoutCode, 504},
		{CanceledCode, 408},
		{AlreadyExistsCode, 409},
		{UnauthorizedCode, 401},
		{ForbiddenCode, 403},
		{ValidationErrorCode, 400},
		{BusinessRuleViolationCode, 422},
		{ExternalServiceErrorCode, 502},
		{NetworkErrorCode, 503},
		{ConfigurationErrorCode, 500},
		{ResourceExhaustedCode, 429},
		{DataCorruptionCode, 500},
		{ConcurrencyErrorCode, 409},
	}
	
	for _, tc := range testCases {
		t.Run(string(tc.code), func(t *testing.T) {
			assert.Equal(t, tc.expected, GetHTTPStatus(tc.code))
		})
	}
	
	// Test unknown error code
	assert.Equal(t, 500, GetHTTPStatus("UNKNOWN_CODE"))
}