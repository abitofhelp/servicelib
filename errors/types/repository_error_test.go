// Copyright (c) 2025 A Bit of Help, Inc.

package types

import (
	"errors"
	"net/http"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

// TestRepositoryError_Error tests the Error method of RepositoryError
func TestRepositoryError_Error(t *testing.T) {
	// Create a repository error
	err := &RepositoryError{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the Error method
	assert.Equal(t, "Repository error (TEST_CODE): test message", err.Error())
}

// TestRepositoryError_Unwrap tests the Unwrap method of RepositoryError
func TestRepositoryError_Unwrap(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create a repository error
	err := &RepositoryError{
		Err:     originalErr,
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the Unwrap method
	assert.Equal(t, originalErr, err.Unwrap())
}

// TestRepositoryError_HTTPStatus tests the HTTPStatus method of RepositoryError
func TestRepositoryError_HTTPStatus(t *testing.T) {
	// Create a repository error
	err := &RepositoryError{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
}

// TestRepositoryError_IsRepositoryError tests the IsRepositoryError method of RepositoryError
func TestRepositoryError_IsRepositoryError(t *testing.T) {
	// Create a repository error
	err := &RepositoryError{
		Err:     errors.New("original error"),
		Message: "test message",
		Code:    "TEST_CODE",
	}

	// Test the IsRepositoryError method
	assert.True(t, err.IsRepositoryError())
}

// TestNewRepositoryError tests the NewRepositoryError function
func TestNewRepositoryError(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create a repository error using NewRepositoryError
	err := NewRepositoryError(originalErr, "test message", "TEST_CODE")

	// Test the error properties
	assert.Equal(t, originalErr, err.Err)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, "TEST_CODE", err.Code)
}

// TestNotFoundError_Error tests the Error method of NotFoundError
func TestNotFoundError_Error(t *testing.T) {
	// Create a not found error
	err := &NotFoundError{
		ResourceType: "User",
		ID:           "123",
	}

	// Test the Error method
	assert.Equal(t, "User not found: 123", err.Error())
}

// TestNotFoundError_Code tests the Code method of NotFoundError
func TestNotFoundError_Code(t *testing.T) {
	// Create a not found error
	err := &NotFoundError{
		ResourceType: "User",
		ID:           "123",
	}

	// Test the Code method
	assert.Equal(t, string(core.NotFoundCode), err.Code())
}

// TestNotFoundError_HTTPStatus tests the HTTPStatus method of NotFoundError
func TestNotFoundError_HTTPStatus(t *testing.T) {
	// Create a not found error
	err := &NotFoundError{
		ResourceType: "User",
		ID:           "123",
	}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusNotFound, err.HTTPStatus())
}

// TestNotFoundError_IsNotFoundError tests the IsNotFoundError method of NotFoundError
func TestNotFoundError_IsNotFoundError(t *testing.T) {
	// Create a not found error
	err := &NotFoundError{
		ResourceType: "User",
		ID:           "123",
	}

	// Test the IsNotFoundError method
	assert.True(t, err.IsNotFoundError())
}

// TestNewNotFoundError tests the NewNotFoundError function
func TestNewNotFoundError(t *testing.T) {
	// Create a not found error using NewNotFoundError
	err := NewNotFoundError("User", "123")

	// Test the error properties
	assert.Equal(t, "User", err.ResourceType)
	assert.Equal(t, "123", err.ID)
}

// TestNotFoundError_Is tests the Is method of NotFoundError
func TestNotFoundError_Is(t *testing.T) {
	// Create a not found error
	err1 := &NotFoundError{
		ResourceType: "User",
		ID:           "123",
	}

	// Create another not found error with the same resource type and ID
	err2 := &NotFoundError{
		ResourceType: "User",
		ID:           "123",
	}

	// Create another not found error with a different resource type
	err3 := &NotFoundError{
		ResourceType: "Product",
		ID:           "123",
	}

	// Create another not found error with a different ID
	err4 := &NotFoundError{
		ResourceType: "User",
		ID:           "456",
	}

	// Test the Is method
	assert.True(t, err1.Is(err1))
	assert.True(t, err1.Is(err2))
	assert.False(t, err1.Is(err3))
	assert.False(t, err1.Is(err4))
	assert.False(t, err1.Is(errors.New("not a not found error")))
}