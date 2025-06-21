// Copyright (c) 2024 A Bit of Help, Inc.

package wrappers

import (
	"fmt"
)

// RetryableError wraps an error with retryable context
type RetryableError struct {
	Original error
	Attempts int
	MaxAttempts int
	RetryAfter int64 // milliseconds
}

// Error returns the error message
func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d/%d): %v", e.Attempts, e.MaxAttempts, e.Original)
}

// Unwrap returns the original error
func (e *RetryableError) Unwrap() error {
	return e.Original
}

// IsRetryable returns true if the error is retryable
func (e *RetryableError) IsRetryable() bool {
	return e.Attempts < e.MaxAttempts
}

// NewRetryableError creates a new retryable error
func NewRetryableError(err error, attempts, maxAttempts int, retryAfter int64) *RetryableError {
	return &RetryableError{
		Original:   err,
		Attempts:   attempts,
		MaxAttempts: maxAttempts,
		RetryAfter: retryAfter,
	}
}

// WrapAsRetryable wraps an error as retryable
func WrapAsRetryable(err error, attempts, maxAttempts int, retryAfter int64) error {
	if err == nil {
		return nil
	}
	return NewRetryableError(err, attempts, maxAttempts, retryAfter)
}

// IsRetryableError checks if an error is a retryable error
func IsRetryableError(err error) bool {
	if _, ok := err.(*RetryableError); ok {
		return true
	}
	return false
}
