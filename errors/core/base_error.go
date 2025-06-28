// Copyright (c) 2025 A Bit of Help, Inc.

// Package core provides the core error handling functionality for the application.
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// BaseError is the foundation for all error types in the system.
// It provides common functionality for all errors, including error code,
// message, operation, details, cause, and stack trace.
type BaseError struct {
	// Code is a unique error code for categorizing errors
	Code ErrorCode `json:"code,omitempty"`

	// Message is a human-readable error message
	Message string `json:"message,omitempty"`

	// Operation is the name of the operation that failed
	Operation string `json:"operation,omitempty"`

	// Details contains additional information about the error
	Details map[string]interface{} `json:"details,omitempty"`

	// Cause is the underlying error that caused this error
	Cause error `json:"-"`

	// Source is the file and line where the error occurred
	Source string `json:"source,omitempty"`

	// Line is the line number where the error occurred
	Line int `json:"line,omitempty"`
}

// Error returns a string representation of the error.
func (e *BaseError) Error() string {
	var builder strings.Builder

	// Add operation if available
	if e.Operation != "" {
		builder.WriteString(fmt.Sprintf("%s: ", e.Operation))
	}

	// Add message
	builder.WriteString(e.Message)

	// Add cause if available
	if e.Cause != nil {
		builder.WriteString(fmt.Sprintf(": %v", e.Cause))
	}

	// Add source location if available
	if e.Source != "" {
		builder.WriteString(fmt.Sprintf(" (source: %s", e.Source))
		if e.Line > 0 {
			builder.WriteString(fmt.Sprintf(":%d", e.Line))
		}
		builder.WriteString(")")
	}

	return builder.String()
}

// Unwrap returns the underlying error.
func (e *BaseError) Unwrap() error {
	return e.Cause
}

// GetCode returns the error code.
func (e *BaseError) GetCode() ErrorCode {
	return e.Code
}

// GetMessage returns the error message.
func (e *BaseError) GetMessage() string {
	return e.Message
}

// GetOperation returns the operation that failed.
func (e *BaseError) GetOperation() string {
	return e.Operation
}

// GetDetails returns additional details about the error.
func (e *BaseError) GetDetails() map[string]interface{} {
	return e.Details
}

// GetCause returns the underlying error.
func (e *BaseError) GetCause() error {
	return e.Cause
}

// GetSource returns the source location of the error.
func (e *BaseError) GetSource() string {
	return e.Source
}

// GetLine returns the line number where the error occurred.
func (e *BaseError) GetLine() int {
	return e.Line
}

// GetHTTPStatus returns the HTTP status code for the error.
func (e *BaseError) GetHTTPStatus() int {
	return GetHTTPStatus(e.Code)
}

// MarshalJSON implements the json.Marshaler interface.
func (e *BaseError) MarshalJSON() ([]byte, error) {
	type Alias BaseError

	// Create a custom message that includes the message and cause, but not operation or source
	message := e.Message
	if e.Cause != nil {
		message = message + ": " + e.Cause.Error()
	}

	return json.Marshal(&struct {
		*Alias
		Message string `json:"message"`
		Cause   string `json:"cause,omitempty"`
	}{
		Alias:   (*Alias)(e),
		Message: message,
		Cause: func() string {
			if e.Cause != nil {
				return e.Cause.Error()
			}
			return ""
		}(),
	})
}

// NewBaseError creates a new BaseError with caller information.
func NewBaseError(code ErrorCode, message string, cause error) *BaseError {
	// Get caller information
	source, line := getCallerInfo(2)

	return &BaseError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Source:  source,
		Line:    line,
	}
}

// WithOperation adds an operation to the error.
func (e *BaseError) WithOperation(operation string) *BaseError {
	e.Operation = operation
	return e
}

// WithDetails adds details to the error.
func (e *BaseError) WithDetails(details map[string]interface{}) *BaseError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// getCallerInfo returns the file name and line number of the caller.
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0
	}
	return filepath.Base(file), line
}

// Is reports whether the error is of the given target type.
func (e *BaseError) Is(target error) bool {
	if t, ok := target.(*BaseError); ok {
		return e.Code == t.Code
	}
	return errors.Is(e.Cause, target)
}

// As finds the first error in err's chain that matches target.
func (e *BaseError) As(target interface{}) bool {
	if t, ok := target.(*BaseError); ok {
		*t = *e
		return true
	}
	return errors.As(e.Cause, target)
}
