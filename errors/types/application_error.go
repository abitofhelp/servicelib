// Copyright (c) 2024 A Bit of Help, Inc.

package types

import "fmt"

// ApplicationError represents an error that occurred in the application layer
type ApplicationError struct {
	Err     error
	Message string
	Code    string
}

// Error returns a string representation of the error
func (e *ApplicationError) Error() string {
	return fmt.Sprintf("Application error (%s): %s", e.Code, e.Message)
}

// Unwrap returns the original error
func (e *ApplicationError) Unwrap() error {
	return e.Err
}

// NewApplicationError creates a new ApplicationError
func NewApplicationError(err error, message, code string) *ApplicationError {
	return &ApplicationError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// AppError is a generic error type that can be used for different error categories
type AppError[T any] struct {
	Err     error
	Message string
	Code    string
	Type    T
}

// Error returns the error message
func (e *AppError[T]) Error() string {
	return fmt.Sprintf("%v (%s): %s", e.Type, e.Code, e.Message)
}

// Unwrap returns the original error
func (e *AppError[T]) Unwrap() error {
	return e.Err
}

// ErrorType returns the error type
func (e *AppError[T]) ErrorType() T {
	return e.Type
}

// NewAppError creates a new AppError
func NewAppError[T any](err error, message, code string, errorType T) *AppError[T] {
	return &AppError[T]{
		Err:     err,
		Message: message,
		Code:    code,
		Type:    errorType,
	}
}
