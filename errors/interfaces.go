// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides generic error interfaces and implementations that can be used across different applications.
package errors

import (
	"net/http"
)

// ErrorWithCode is an interface for errors that have an error code
type ErrorWithCode interface {
	error
	// Code returns the error code
	Code() string
}

// ErrorWithHTTPStatus is an interface for errors that have an HTTP status code
type ErrorWithHTTPStatus interface {
	error
	// HTTPStatus returns the HTTP status code
	HTTPStatus() int
}

// ValidationErrorInterface is an interface for validation errors
type ValidationErrorInterface interface {
	ErrorWithCode
	ErrorWithHTTPStatus
	// IsValidationError identifies this as a validation error
	IsValidationError() bool
}

// NotFoundErrorInterface is an interface for not found errors
type NotFoundErrorInterface interface {
	ErrorWithCode
	ErrorWithHTTPStatus
	// IsNotFoundError identifies this as a not found error
	IsNotFoundError() bool
}

// ApplicationErrorInterface is an interface for application errors
type ApplicationErrorInterface interface {
	ErrorWithCode
	ErrorWithHTTPStatus
	// IsApplicationError identifies this as an application error
	IsApplicationError() bool
}

// RepositoryErrorInterface is an interface for repository errors
type RepositoryErrorInterface interface {
	ErrorWithCode
	ErrorWithHTTPStatus
	// IsRepositoryError identifies this as a repository error
	IsRepositoryError() bool
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return As(err, &validationErr)
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	var notFoundErr *NotFoundError
	return As(err, &notFoundErr)
}

// IsApplicationError checks if an error is an application error
func IsApplicationError(err error) bool {
	var appErr ApplicationErrorInterface
	return As(err, &appErr)
}

// IsRepositoryError checks if an error is a repository error
func IsRepositoryError(err error) bool {
	var repoErr RepositoryErrorInterface
	return As(err, &repoErr)
}

// GetHTTPStatusFromError returns the HTTP status code for an error
func GetHTTPStatusFromError(err error) int {
	// Check if the error implements ErrorWithHTTPStatus
	if httpErr, ok := err.(ErrorWithHTTPStatus); ok {
		return httpErr.HTTPStatus()
	}

	// Check specific error types
	if IsValidationError(err) {
		return http.StatusBadRequest
	}
	if IsNotFoundError(err) {
		return http.StatusNotFound
	}
	if IsApplicationError(err) || IsRepositoryError(err) {
		return http.StatusInternalServerError
	}

	// Default status code
	return http.StatusInternalServerError
}
