// Copyright (c) 2025 A Bit of Help, Inc.

package types

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDomainError_Error tests the Error method of DomainError
func TestDomainError_Error(t *testing.T) {
	// Create a domain error
	err := &DomainError{
		Original: errors.New("original error"),
		Code:     "TEST_CODE",
		Message:  "test message",
		Op:       "TestOperation",
		Param:    "testParam",
	}

	// Test the Error method
	assert.Equal(t, "TestOperation: test message", err.Error())
}

// TestDomainError_Unwrap tests the Unwrap method of DomainError
func TestDomainError_Unwrap(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create a domain error
	err := &DomainError{
		Original: originalErr,
		Code:     "TEST_CODE",
		Message:  "test message",
		Op:       "TestOperation",
		Param:    "testParam",
	}

	// Test the Unwrap method
	assert.Equal(t, originalErr, err.Unwrap())
}

// TestDomainError_HTTPStatus tests the HTTPStatus method of DomainError
func TestDomainError_HTTPStatus(t *testing.T) {
	// Create a domain error
	err := &DomainError{
		Original: errors.New("original error"),
		Code:     "TEST_CODE",
		Message:  "test message",
		Op:       "TestOperation",
		Param:    "testParam",
	}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus())
}

// TestDomainError_IsDomainError tests the IsDomainError method of DomainError
func TestDomainError_IsDomainError(t *testing.T) {
	// Create a domain error
	err := &DomainError{
		Original: errors.New("original error"),
		Code:     "TEST_CODE",
		Message:  "test message",
		Op:       "TestOperation",
		Param:    "testParam",
	}

	// Test the IsDomainError method
	assert.True(t, err.IsDomainError())
}

// TestDomainError_Is tests the Is method of DomainError
func TestDomainError_Is(t *testing.T) {
	// Create a domain error
	err1 := &DomainError{
		Original: errors.New("original error"),
		Code:     "TEST_CODE",
		Message:  "test message",
		Op:       "TestOperation",
		Param:    "testParam",
	}

	// Create another domain error with the same code
	err2 := &DomainError{
		Original: errors.New("another error"),
		Code:     "TEST_CODE",
		Message:  "another message",
		Op:       "AnotherOperation",
		Param:    "anotherParam",
	}

	// Create another domain error with a different code
	err3 := &DomainError{
		Original: errors.New("third error"),
		Code:     "DIFFERENT_CODE",
		Message:  "third message",
		Op:       "ThirdOperation",
		Param:    "thirdParam",
	}

	// Test the Is method
	assert.True(t, err1.Is(err1))
	assert.True(t, err1.Is(err2))
	assert.False(t, err1.Is(err3))
	assert.False(t, err1.Is(errors.New("not a domain error")))
}

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create a domain error using New
	err := New("TestOperation", "TEST_CODE", "test message", originalErr)

	// Test the error properties
	assert.Equal(t, originalErr, err.Original)
	assert.Equal(t, "TEST_CODE", err.Code)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, "TestOperation", err.Op)
	assert.Equal(t, "", err.Param)
}

// TestWrap tests the Wrap function
func TestWrap(t *testing.T) {
	// Create an original error
	originalErr := errors.New("original error")

	// Create a domain error using Wrap
	err := Wrap(originalErr, "TestOperation", "test message")

	// Test the error properties
	assert.Equal(t, originalErr, err.Original)
	assert.Equal(t, "", err.Code)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, "TestOperation", err.Op)
	assert.Equal(t, "", err.Param)
}
