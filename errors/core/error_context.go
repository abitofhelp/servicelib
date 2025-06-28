// Copyright (c) 2024 A Bit of Help, Inc.

// Package core provides the core error handling functionality for the errors package.
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrorContext holds additional context for an error.
// It includes information about the operation that failed, the source location,
// and any additional details that might be useful for debugging or error reporting.
type ErrorContext struct {
	// Operation is the name of the operation that failed
	Operation string `json:"operation,omitempty"`

	// Source is the file and line where the error occurred
	Source string `json:"source,omitempty"`

	// Line is the line number where the error occurred
	Line int `json:"line,omitempty"`

	// Code is the error code
	Code ErrorCode `json:"code,omitempty"`

	// HTTPStatus is the HTTP status code to return for this error
	HTTPStatus int `json:"http_status,omitempty"`

	// Details contains additional information about the error
	Details map[string]interface{} `json:"details,omitempty"`
}

// ContextualError is an error with additional context.
// It wraps another error and adds contextual information like operation name,
// error code, HTTP status, and source location.
type ContextualError struct {
	// Original is the original error that was wrapped
	Original error

	// Context contains additional information about the error
	Context ErrorContext
}

// Error returns the error message with contextual information.
func (e *ContextualError) Error() string {
	var builder strings.Builder

	// Add operation if available
	if e.Context.Operation != "" {
		builder.WriteString(fmt.Sprintf("%s: ", e.Context.Operation))
	}

	// Add original error message
	if e.Original != nil {
		builder.WriteString(e.Original.Error())
	} else {
		builder.WriteString("an error occurred")
	}

	// Add source location if available
	if e.Context.Source != "" {
		builder.WriteString(fmt.Sprintf(" (source: %s", e.Context.Source))
		if e.Context.Line > 0 {
			builder.WriteString(fmt.Sprintf(":%d", e.Context.Line))
		}
		builder.WriteString(")")
	}

	return builder.String()
}

// Unwrap returns the original error.
func (e *ContextualError) Unwrap() error {
	return e.Original
}

// Code returns the error code.
func (e *ContextualError) Code() ErrorCode {
	return e.Context.Code
}

// HTTPStatus returns the HTTP status code.
func (e *ContextualError) HTTPStatus() int {
	return e.Context.HTTPStatus
}

// MarshalJSON implements the json.Marshaler interface.
func (e *ContextualError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message    string                 `json:"message"`
		Operation  string                 `json:"operation,omitempty"`
		Source     string                 `json:"source,omitempty"`
		Line       int                    `json:"line,omitempty"`
		Code       ErrorCode              `json:"code,omitempty"`
		HTTPStatus int                    `json:"http_status,omitempty"`
		Details    map[string]interface{} `json:"details,omitempty"`
	}{
		Message:    e.Error(),
		Operation:  e.Context.Operation,
		Source:     e.Context.Source,
		Line:       e.Context.Line,
		Code:       e.Context.Code,
		HTTPStatus: e.Context.HTTPStatus,
		Details:    e.Context.Details,
	})
}

// GetCallerInfo returns the file name and line number of the caller.
func GetCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0
	}
	return filepath.Base(file), line
}

// WithContext wraps an error with contextual information.
// It adds operation name, error code, HTTP status, and source location to the error.
// If the error is already a ContextualError, it updates the context with the new information.
func WithContext(err error, operation string, code ErrorCode, httpStatus int, details map[string]interface{}) error {
	// Get caller information
	source, line := GetCallerInfo(2)

	// If the error is already a ContextualError, update its context
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		// Only update operation if it's not already set
		if operation != "" && contextualErr.Context.Operation == "" {
			contextualErr.Context.Operation = operation
		}

		// Only update code if it's not already set
		if code != "" && contextualErr.Context.Code == "" {
			contextualErr.Context.Code = code
		}

		// Only update HTTP status if it's not already set
		if httpStatus != 0 && contextualErr.Context.HTTPStatus == 0 {
			contextualErr.Context.HTTPStatus = httpStatus
		}

		// Merge details if provided
		if details != nil {
			if contextualErr.Context.Details == nil {
				contextualErr.Context.Details = make(map[string]interface{})
			}
			for k, v := range details {
				contextualErr.Context.Details[k] = v
			}
		}

		return contextualErr
	}

	// Create a new ContextualError
	return &ContextualError{
		Original: err,
		Context: ErrorContext{
			Operation:  operation,
			Source:     source,
			Line:       line,
			Code:       code,
			HTTPStatus: httpStatus,
			Details:    details,
		},
	}
}
