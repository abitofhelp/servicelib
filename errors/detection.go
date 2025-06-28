// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides error handling utilities for the servicelib package.
package errors

import (
	"context"
	"fmt"
)

// DetectErrorType returns a string describing the type of error
func DetectErrorType(err error) string {
	if err == nil {
		return "nil"
	}

	if IsValidationError(err) {
		return "ValidationError"
	}

	if IsDomainError(err) {
		return "DomainError"
	}

	if IsNotFoundError(err) {
		return "NotFoundError"
	}

	if IsDatabaseError(err) {
		return "DatabaseError"
	}

	return fmt.Sprintf("%T", err)
}

// FormatErrorWithSource formats an error with its source information
func FormatErrorWithSource(err error, source string) string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("[%s] %s", source, err.Error())
}

// DetectErrorFromContext extracts an error from a context if present
func DetectErrorFromContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

// This function is used to ensure the context package is used
var _ = context.Background
