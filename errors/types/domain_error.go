// Copyright (c) 2024 A Bit of Help, Inc.

package types

import "fmt"

// DomainErrorType represents the type of domain error
type DomainErrorType string

const (
	DomainErrorGeneral DomainErrorType = "DOMAIN_ERROR"
)

// DomainError represents a domain error with operation context
type DomainError struct {
	Original error
	Code    string
	Message string
	Op      string
	Param   string
}

// Error returns a string representation of the error
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Unwrap returns the original error
func (e *DomainError) Unwrap() error {
	return e.Original
}

// Is reports whether the error is of the given target type
func (e *DomainError) Is(target error) bool {
	if t, ok := target.(*DomainError); ok {
		return e.Code == t.Code
	}
	return false
}

// New creates a new DomainError
func New(op, code, message string, original error) *DomainError {
	return &DomainError{
		Original: original,
		Code:    code,
		Message: message,
		Op:      op,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, op, message string) *DomainError {
	return New(op, "", message, err)
}
