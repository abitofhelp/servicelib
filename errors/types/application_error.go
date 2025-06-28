// Copyright (c) 2024 A Bit of Help, Inc.

// Package types provides specific error types for the errors package.
package types

import (
	"fmt"
	"net/http"
)

// ApplicationError represents an error that occurred in the application layer.
type ApplicationError struct {
	// Err is the original error
	Err error

	// Message is a human-readable error message
	Message string

	// Code is a machine-readable error code
	Code string
}

// Error returns a string representation of the error.
func (e *ApplicationError) Error() string {
	return fmt.Sprintf("Application error (%s): %s", e.Code, e.Message)
}

// Unwrap returns the original error.
func (e *ApplicationError) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the HTTP status code for application errors.
func (e *ApplicationError) HTTPStatus() int {
	return http.StatusInternalServerError
}

// IsApplicationError identifies this as an application error.
func (e *ApplicationError) IsApplicationError() bool {
	return true
}

// NewApplicationError creates a new ApplicationError.
func NewApplicationError(err error, message, code string) *ApplicationError {
	return &ApplicationError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// AppError is a generic error type that can be used for different error categories.
type AppError[T any] struct {
	// Err is the original error
	Err error

	// Message is a human-readable error message
	Message string

	// Code is a machine-readable error code
	Code string

	// Type is the error type
	Type T
}

// Error returns a string representation of the error.
func (e *AppError[T]) Error() string {
	return fmt.Sprintf("%v (%s): %s", e.Type, e.Code, e.Message)
}

// Unwrap returns the original error.
func (e *AppError[T]) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the HTTP status code for application errors.
func (e *AppError[T]) HTTPStatus() int {
	return http.StatusInternalServerError
}

// IsApplicationError identifies this as an application error.
func (e *AppError[T]) IsApplicationError() bool {
	return true
}

// ErrorType returns the error type.
func (e *AppError[T]) ErrorType() T {
	return e.Type
}

// NewAppError creates a new AppError.
func NewAppError[T any](err error, message, code string, errorType T) *AppError[T] {
	return &AppError[T]{
		Err:     err,
		Message: message,
		Code:    code,
		Type:    errorType,
	}
}
