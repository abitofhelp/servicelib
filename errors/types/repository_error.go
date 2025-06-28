// Copyright (c) 2024 A Bit of Help, Inc.

// Package types provides specific error types for the errors package.
package types

import (
	"fmt"
	"net/http"

	"github.com/abitofhelp/servicelib/errors/core"
)

// RepositoryError represents an error that occurred in the repository layer.
type RepositoryError struct {
	// Err is the original error
	Err error

	// Message is a human-readable error message
	Message string

	// Code is a machine-readable error code
	Code string
}

// Error returns a string representation of the error.
func (e *RepositoryError) Error() string {
	return fmt.Sprintf("Repository error (%s): %s", e.Code, e.Message)
}

// Unwrap returns the original error.
func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the HTTP status code for repository errors.
func (e *RepositoryError) HTTPStatus() int {
	return http.StatusInternalServerError
}

// IsRepositoryError identifies this as a repository error.
func (e *RepositoryError) IsRepositoryError() bool {
	return true
}

// NewRepositoryError creates a new RepositoryError.
func NewRepositoryError(err error, message, code string) *RepositoryError {
	return &RepositoryError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	// ResourceType is the type of resource that was not found
	ResourceType string

	// ID is the identifier of the resource that was not found
	ID string
}

// Error returns a string representation of the error.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.ResourceType, e.ID)
}

// Code returns the not found error code.
func (e *NotFoundError) Code() string {
	return string(core.NotFoundCode)
}

// HTTPStatus returns the HTTP status code for not found errors.
func (e *NotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

// IsNotFoundError identifies this as a not found error.
func (e *NotFoundError) IsNotFoundError() bool {
	return true
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(resourceType, id string) *NotFoundError {
	return &NotFoundError{
		ResourceType: resourceType,
		ID:           id,
	}
}

// Is implements the errors.Is interface for NotFoundError.
func (e *NotFoundError) Is(target error) bool {
	if t, ok := target.(*NotFoundError); ok {
		return e.ResourceType == t.ResourceType && e.ID == t.ID
	}
	return false
}
