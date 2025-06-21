// Copyright (c) 2024 A Bit of Help, Inc.

package types

import "fmt"

// RepositoryError represents an error that occurred in the repository layer
type RepositoryError struct {
	Err     error
	Message string
	Code    string
}

// Error returns a string representation of the error
func (e *RepositoryError) Error() string {
	return fmt.Sprintf("Repository error (%s): %s", e.Code, e.Message)
}

// Unwrap returns the original error
func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError creates a new RepositoryError
func NewRepositoryError(err error, message, code string) *RepositoryError {
	return &RepositoryError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	ResourceType string
	ID           string
}

// Error returns a string representation of the error
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.ResourceType, e.ID)
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resourceType, id string) *NotFoundError {
	return &NotFoundError{
		ResourceType: resourceType,
		ID:          id,
	}
}

// Is implements the errors.Is interface for NotFoundError
func (e *NotFoundError) Is(target error) bool {
	if t, ok := target.(*NotFoundError); ok {
		return e.ResourceType == t.ResourceType && e.ID == t.ID
	}
	return false
}
