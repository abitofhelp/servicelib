// Copyright (c) 2025 A Bit of Help, Inc.

package wrappers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRetryableError_Error tests the Error method of RetryableError
func TestRetryableError_Error(t *testing.T) {
	// Create a retryable error
	originalErr := errors.New("original error")
	err := &RetryableError{
		Original:   originalErr,
		Attempts:   1,
		MaxAttempts: 3,
		RetryAfter: 1000,
	}

	// Test the Error method
	assert.Equal(t, "retryable error (attempt 1/3): original error", err.Error())
}

// TestRetryableError_Unwrap tests the Unwrap method of RetryableError
func TestRetryableError_Unwrap(t *testing.T) {
	// Create a retryable error
	originalErr := errors.New("original error")
	err := &RetryableError{
		Original:   originalErr,
		Attempts:   1,
		MaxAttempts: 3,
		RetryAfter: 1000,
	}

	// Test the Unwrap method
	assert.Equal(t, originalErr, err.Unwrap())
}

// TestRetryableError_IsRetryable tests the IsRetryable method of RetryableError
func TestRetryableError_IsRetryable(t *testing.T) {
	// Create a retryable error with attempts < maxAttempts
	err1 := &RetryableError{
		Original:   errors.New("original error"),
		Attempts:   1,
		MaxAttempts: 3,
		RetryAfter: 1000,
	}

	// Test the IsRetryable method with attempts < maxAttempts
	assert.True(t, err1.IsRetryable())

	// Create a retryable error with attempts = maxAttempts
	err2 := &RetryableError{
		Original:   errors.New("original error"),
		Attempts:   3,
		MaxAttempts: 3,
		RetryAfter: 1000,
	}

	// Test the IsRetryable method with attempts = maxAttempts
	assert.False(t, err2.IsRetryable())

	// Create a retryable error with attempts > maxAttempts
	err3 := &RetryableError{
		Original:   errors.New("original error"),
		Attempts:   4,
		MaxAttempts: 3,
		RetryAfter: 1000,
	}

	// Test the IsRetryable method with attempts > maxAttempts
	assert.False(t, err3.IsRetryable())
}

// TestNewRetryableError tests the NewRetryableError function
func TestNewRetryableError(t *testing.T) {
	// Create a retryable error using NewRetryableError
	originalErr := errors.New("original error")
	err := NewRetryableError(originalErr, 1, 3, 1000)

	// Test the error properties
	assert.Equal(t, originalErr, err.Original)
	assert.Equal(t, 1, err.Attempts)
	assert.Equal(t, 3, err.MaxAttempts)
	assert.Equal(t, int64(1000), err.RetryAfter)
}

// TestWrapAsRetryable tests the WrapAsRetryable function
func TestWrapAsRetryable(t *testing.T) {
	// Test with nil error
	err := WrapAsRetryable(nil, 1, 3, 1000)
	assert.Nil(t, err)

	// Test with standard error
	originalErr := errors.New("original error")
	err = WrapAsRetryable(originalErr, 1, 3, 1000)
	
	// Check that the error is a RetryableError
	re, ok := err.(*RetryableError)
	assert.True(t, ok)
	
	// Check the error properties
	assert.Equal(t, originalErr, re.Original)
	assert.Equal(t, 1, re.Attempts)
	assert.Equal(t, 3, re.MaxAttempts)
	assert.Equal(t, int64(1000), re.RetryAfter)
}

// TestIsRetryableError tests the IsRetryableError function
func TestIsRetryableError(t *testing.T) {
	// Test with nil error
	assert.False(t, IsRetryableError(nil))

	// Test with standard error
	assert.False(t, IsRetryableError(errors.New("standard error")))

	// Test with retryable error
	re := &RetryableError{
		Original:   errors.New("original error"),
		Attempts:   1,
		MaxAttempts: 3,
		RetryAfter: 1000,
	}
	assert.True(t, IsRetryableError(re))
}