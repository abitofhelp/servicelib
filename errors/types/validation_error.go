// Copyright (c) 2024 A Bit of Help, Inc.

package types

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Msg   string
	Field string
}

// Error returns the error message
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Msg)
	}
	return e.Msg
}

// NewValidationError creates a new ValidationError
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{
		Msg: msg,
	}
}

// NewFieldValidationError creates a new ValidationError with a field name
func NewFieldValidationError(msg, field string) *ValidationError {
	return &ValidationError{
		Msg:   msg,
		Field: field,
	}
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []*ValidationError
}

// Error returns a combined error message
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	msgs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("%d validation errors: %s", len(e.Errors), strings.Join(msgs, "; "))
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors(errors ...*ValidationError) *ValidationErrors {
	return &ValidationErrors{
		Errors: errors,
	}
}

// AddError adds a validation error to the collection
func (e *ValidationErrors) AddError(err *ValidationError) {
	e.Errors = append(e.Errors, err)
}

// HasErrors returns true if there are any validation errors
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}
