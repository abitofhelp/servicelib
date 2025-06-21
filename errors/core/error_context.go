// Copyright (c) 2024 A Bit of Help, Inc.

package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
)

// ErrorContext holds additional context for an error.
type ErrorContext struct {
	Operation string `json:"operation,omitempty"`
	Source    string `json:"source,omitempty"`
	Line      int    `json:"line,omitempty"`
	Code      ErrorCode `json:"code,omitempty"`
	HTTPStatus int    `json:"http_status,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// ContextualError is an error with additional context.
type ContextualError struct {
	Original error
	Context ErrorContext
}

// Error returns the error message with contextual information.
func (e *ContextualError) Error() string {
	if e.Context.Operation != "" {
		return fmt.Sprintf("%s: %v", e.Context.Operation, e.Original)
	}
	return e.Original.Error()
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
	return json.Marshal(e.Context)
}

// getCallerInfo returns the file name and line number of the caller.
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown", 0
	}
	return filepath.Base(file), line
}
