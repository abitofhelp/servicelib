// Copyright (c) 2025 A Bit of Help, Inc.

package auth_test

import (
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/auth"
	serviceErrors "github.com/abitofhelp/servicelib/errors"
	"github.com/stretchr/testify/assert"
)

// TestErrorConstants tests the error constants defined in errors_compat.go
func TestErrorConstants(t *testing.T) {
	// Test that error constants are not nil
	assert.NotNil(t, auth.ErrInvalidToken)
	assert.NotNil(t, auth.ErrExpiredToken)
	assert.NotNil(t, auth.ErrMissingToken)
	assert.NotNil(t, auth.ErrInvalidSignature)
	assert.NotNil(t, auth.ErrInvalidClaims)
	assert.NotNil(t, auth.ErrUnauthorized)
	assert.NotNil(t, auth.ErrForbidden)
	assert.NotNil(t, auth.ErrInvalidConfig)
	assert.NotNil(t, auth.ErrInternal)
	assert.NotNil(t, auth.ErrNotImplemented)
}

// TestWithContext tests the WithContext function
func TestWithContext(t *testing.T) {
	// Test with a nil error
	err := auth.WithContext(nil, "key", "value")
	assert.Nil(t, err)

	// Test with a standard error
	originalErr := errors.New("test error")
	err = auth.WithContext(originalErr, "key", "value")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "test error")

	// Test with a service error
	serviceErr := serviceErrors.New(serviceErrors.InvalidInputCode, "service error")
	err = auth.WithContext(serviceErr, "key", "value")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service error")

	// Test getting the context
	value, ok := auth.GetContext(err, "key")
	assert.True(t, ok)
	assert.Equal(t, "value", value)

	// Test getting a non-existent context
	value, ok = auth.GetContext(err, "nonexistent")
	assert.False(t, ok)
	assert.Nil(t, value)

	// Test getting context from a standard error
	value, ok = auth.GetContext(originalErr, "key")
	assert.False(t, ok)
	assert.Nil(t, value)
}

// TestWithOp tests the WithOp function
func TestWithOp(t *testing.T) {
	// Test with a nil error
	err := auth.WithOp(nil, "operation")
	assert.Nil(t, err)

	// Test with a standard error
	originalErr := errors.New("test error")
	err = auth.WithOp(originalErr, "operation")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "test error")

	// Test with a service error
	serviceErr := serviceErrors.New(serviceErrors.InvalidInputCode, "service error")
	err = auth.WithOp(serviceErr, "operation")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service error")

	// Test getting the operation
	op, ok := auth.GetOp(err)
	assert.True(t, ok)
	assert.Equal(t, "operation", op)

	// Test getting operation from a standard error
	op, ok = auth.GetOp(originalErr)
	assert.False(t, ok)
	assert.Equal(t, "", op)
}

// TestWithMessage tests the WithMessage function
func TestWithMessage(t *testing.T) {
	// Test with a nil error
	err := auth.WithMessage(nil, "message")
	assert.Nil(t, err)

	// Test with a standard error
	originalErr := errors.New("test error")
	err = auth.WithMessage(originalErr, "message")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "test error")
	assert.Contains(t, err.Error(), "message")

	// Test with a service error
	serviceErr := serviceErrors.New(serviceErrors.InvalidInputCode, "service error")
	err = auth.WithMessage(serviceErr, "message")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service error")
	assert.Contains(t, err.Error(), "message")

	// Test getting the message
	msg, ok := auth.GetMessage(err)
	assert.True(t, ok)
	assert.Contains(t, msg, "message")

	// Test getting message from a standard error
	msg, ok = auth.GetMessage(originalErr)
	assert.True(t, ok)
	assert.Equal(t, "test error", msg)
}

// TestWrap tests the Wrap function
func TestWrap(t *testing.T) {
	// Test with a nil error
	err := auth.Wrap(nil, "message")
	assert.Nil(t, err)

	// Test with a standard error
	originalErr := errors.New("test error")
	err = auth.Wrap(originalErr, "message")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "test error")
	assert.Contains(t, err.Error(), "message")

	// Test with a service error
	serviceErr := serviceErrors.New(serviceErrors.InvalidInputCode, "service error")
	err = auth.Wrap(serviceErr, "message")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service error")
	assert.Contains(t, err.Error(), "message")
}

// TestGetErrorCode tests the getErrorCode function indirectly through other functions
func TestGetErrorCode(t *testing.T) {
	// Test with different error types that have known HTTP status codes
	testCases := []struct {
		err  error
		code serviceErrors.ErrorCode
	}{
		{serviceErrors.New(serviceErrors.InvalidInputCode, "invalid input"), serviceErrors.InvalidInputCode},
		{serviceErrors.New(serviceErrors.UnauthorizedCode, "unauthorized"), serviceErrors.UnauthorizedCode},
		{serviceErrors.New(serviceErrors.ForbiddenCode, "forbidden"), serviceErrors.ForbiddenCode},
		{serviceErrors.New(serviceErrors.NotFoundCode, "not found"), serviceErrors.NotFoundCode},
		{serviceErrors.New(serviceErrors.AlreadyExistsCode, "already exists"), serviceErrors.AlreadyExistsCode},
		{serviceErrors.New(serviceErrors.ValidationErrorCode, "validation error"), serviceErrors.ValidationErrorCode},
		{serviceErrors.New(serviceErrors.ResourceExhaustedCode, "resource exhausted"), serviceErrors.ResourceExhaustedCode},
		{serviceErrors.New(serviceErrors.InternalErrorCode, "internal error"), serviceErrors.InternalErrorCode},
		{serviceErrors.New(serviceErrors.ExternalServiceErrorCode, "external service error"), serviceErrors.ExternalServiceErrorCode},
		{errors.New("standard error"), serviceErrors.InternalErrorCode}, // Default case for standard errors
	}

	for _, tc := range testCases {
		// Use WithOp to indirectly test getErrorCode
		wrappedErr := auth.WithOp(tc.err, "operation")

		// Check that the error code is preserved
		if e, ok := wrappedErr.(interface{ GetCode() serviceErrors.ErrorCode }); ok {
			assert.Equal(t, tc.code, e.GetCode())
		} else if tc.err.Error() == "standard error" {
			// For standard errors, we can't check the code directly
			// but we can check that the error message is preserved
			assert.Contains(t, wrappedErr.Error(), "standard error")
		} else {
			t.Errorf("Expected wrapped error to implement GetCode, but it doesn't")
		}
	}
}
