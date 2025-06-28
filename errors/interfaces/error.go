// Copyright (c) 2024 A Bit of Help, Inc.

package interfaces

// Error is a generic error interface that can be used for different error categories
type Error interface {
	error
	Code() string
	HTTPStatus() int
	Is(target error) bool
	Unwrap() error
}

// Wrapper is an interface for error wrappers
type Wrapper interface {
	Wrap(err error, operation string, format string, args ...interface{}) error
	WithDetails(err error, details map[string]interface{}) error
}

// ErrorType is a type for error categories
type ErrorType string

// Standard error types
const (
	DomainErrorType      ErrorType = "DOMAIN_ERROR"
	ValidationErrorType  ErrorType = "VALIDATION_ERROR"
	RepositoryErrorType  ErrorType = "REPOSITORY_ERROR"
	ApplicationErrorType ErrorType = "APPLICATION_ERROR"
)
