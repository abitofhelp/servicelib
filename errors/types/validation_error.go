// Copyright (c) 2024 A Bit of Help, Inc.

// Package types provides specific error types for the errors package.
package types

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abitofhelp/servicelib/errors/core"
)

// ValidationError represents a validation error for a specific field.
type ValidationError struct {
	// Msg is the error message
	Msg string

	// Field is the name of the field that failed validation
	Field string
}

// Error returns a string representation of the validation error.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s (field: %s)", e.Msg, e.Field)
	}
	return fmt.Sprintf("validation error: %s", e.Msg)
}

// Code returns the validation error code.
func (e *ValidationError) Code() string {
	return string(core.ValidationErrorCode)
}

// HTTPStatus returns the HTTP status code for validation errors.
func (e *ValidationError) HTTPStatus() int {
	return http.StatusBadRequest
}

// IsValidationError identifies this as a validation error.
func (e *ValidationError) IsValidationError() bool {
	return true
}

// NewValidationError creates a new ValidationError with a message.
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Msg: msg}
}

// NewFieldValidationError creates a new ValidationError with a message and field name.
func NewFieldValidationError(msg, field string) *ValidationError {
	return &ValidationError{Msg: msg, Field: field}
}

// ValidationErrors represents multiple validation errors.
type ValidationErrors struct {
	// Errors is a slice of validation errors
	Errors []*ValidationError
}

// Error returns a string representation of all validation errors.
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	msgs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("%d validation errors: %s", len(e.Errors), strings.Join(msgs, "; "))
}

// Code returns the validation error code.
func (e *ValidationErrors) Code() string {
	return string(core.ValidationErrorCode)
}

// HTTPStatus returns the HTTP status code for validation errors.
func (e *ValidationErrors) HTTPStatus() int {
	return http.StatusBadRequest
}

// IsValidationError identifies this as a validation error.
func (e *ValidationErrors) IsValidationError() bool {
	return true
}

// NewValidationErrors creates a new ValidationErrors with the provided errors.
func NewValidationErrors(errors ...*ValidationError) *ValidationErrors {
	return &ValidationErrors{Errors: errors}
}

// AddError adds a validation error to the collection.
func (e *ValidationErrors) AddError(err *ValidationError) {
	e.Errors = append(e.Errors, err)
}

// HasErrors returns true if there are any validation errors.
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}
