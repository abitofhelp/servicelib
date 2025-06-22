// Copyright (c) 2024 A Bit of Help, Inc.

// Package types provides specific error types for the errors package.
package types

import (
	"fmt"
	"net/http"
)

// DomainErrorType represents the type of domain error.
type DomainErrorType string

// Domain error type constants.
const (
	DomainErrorGeneral DomainErrorType = "DOMAIN_ERROR"
)

// DomainError represents a domain error with operation context.
type DomainError struct {
	// Original is the original error
	Original error

	// Code is a machine-readable error code
	Code string

	// Message is a human-readable error message
	Message string

	// Op is the operation that caused the error
	Op string

	// Param is the parameter that caused the error
	Param string
}

// Error returns a string representation of the error.
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Unwrap returns the original error.
func (e *DomainError) Unwrap() error {
	return e.Original
}

// HTTPStatus returns the HTTP status code for domain errors.
func (e *DomainError) HTTPStatus() int {
	return http.StatusBadRequest
}

// IsDomainError identifies this as a domain error.
func (e *DomainError) IsDomainError() bool {
	return true
}

// Is reports whether the error is of the given target type.
func (e *DomainError) Is(target error) bool {
	if t, ok := target.(*DomainError); ok {
		return e.Code == t.Code
	}
	return false
}

// New creates a new DomainError.
func New(op, code, message string, original error) *DomainError {
	return &DomainError{
		Original: original,
		Code:    code,
		Message: message,
		Op:      op,
	}
}

// Wrap wraps an error with additional context.
func Wrap(err error, op, message string) *DomainError {
	return New(op, "", message, err)
}
