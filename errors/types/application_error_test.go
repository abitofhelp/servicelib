// Copyright (c) 2025 A Bit of Help, Inc.

package types

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApplicationError_Error tests the Error method of ApplicationError
func TestApplicationError_Error(t *testing.T) {
	// Create an application error
	err := &ApplicationError{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the Error method
	assert.Equal(t, "Application error (TEST_CODE): test message", err.Error())
}

// TestApplicationError_Unwrap tests the Unwrap method of ApplicationError
func TestApplicationError_Unwrap(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create an application error
	err := &ApplicationError{
		Err:     originalErr,
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the Unwrap method
	assert.Equal(t, originalErr, err.Unwrap())
}

// TestApplicationError_HTTPStatus tests the HTTPStatus method of ApplicationError
func TestApplicationError_HTTPStatus(t *testing.T) {
	// Create an application error
	err := &ApplicationError{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
}

// TestApplicationError_IsApplicationError tests the IsApplicationError method of ApplicationError
func TestApplicationError_IsApplicationError(t *testing.T) {
	// Create an application error
	err := &ApplicationError{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the IsApplicationError method
	assert.True(t, err.IsApplicationError())
}

// TestNewApplicationError tests the NewApplicationError function
func TestNewApplicationError(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create an application error using NewApplicationError
	err := NewApplicationError(originalErr, "test message", "TEST_CODE")

	// Test the error properties
	assert.Equal(t, originalErr, err.Err)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, "TEST_CODE", err.Code)
}

// TestAppError_Error tests the Error method of AppError
func TestAppError_Error(t *testing.T) {
	// Create an app error
	err := &AppError[string]{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
		Type:    "TEST_TYPE",
	}

	// Test the Error method
	assert.Equal(t, "TEST_TYPE (TEST_CODE): test message", err.Error())
}

// TestAppError_Unwrap tests the Unwrap method of AppError
func TestAppError_Unwrap(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create an app error
	err := &AppError[string]{
		Err:     originalErr,
		Message: "test message",
		Code:    "TEST_CODE",
		Type:    "TEST_TYPE",
	}

	// Test the Unwrap method
	assert.Equal(t, originalErr, err.Unwrap())
}

// TestAppError_HTTPStatus tests the HTTPStatus method of AppError
func TestAppError_HTTPStatus(t *testing.T) {
	// Create an app error
	err := &AppError[string]{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
		Type:    "TEST_TYPE",
	}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
}

// TestAppError_IsApplicationError tests the IsApplicationError method of AppError
func TestAppError_IsApplicationError(t *testing.T) {
	// Create an app error
	err := &AppError[string]{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
		Type:    "TEST_TYPE",
	}

	// Test the IsApplicationError method
	assert.True(t, err.IsApplicationError())
}

// TestAppError_ErrorType tests the ErrorType method of AppError
func TestAppError_ErrorType(t *testing.T) {
	// Create an app error
	err := &AppError[string]{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
		Type:    "TEST_TYPE",
	}

	// Test the ErrorType method
	assert.Equal(t, "TEST_TYPE", err.ErrorType())
}

// TestNewAppError tests the NewAppError function
func TestNewAppError(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create an app error using NewAppError
	err := NewAppError(originalErr, "test message", "TEST_CODE", "TEST_TYPE")

	// Test the error properties
	assert.Equal(t, originalErr, err.Err)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, "TEST_CODE", err.Code)
	assert.Equal(t, "TEST_TYPE", err.Type)
}