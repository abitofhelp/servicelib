// Copyright (c) 2025 A Bit of Help, Inc.

// Package core provides the core error handling functionality for the errors package.
package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Is checks if an error matches a target error.
// This is a wrapper around errors.Is that adds support for ContextualError.
// Parameters:
//   - err: The error to check
//   - target: The target error to match against
//
// Returns:
//   - bool: True if the error matches the target
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
// This is a wrapper around errors.As that adds support for ContextualError.
// Parameters:
//   - err: The error to check
//   - target: A pointer to the error type to match against
//
// Returns:
//   - bool: True if a match was found
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the underlying error.
// This is a wrapper around errors.Unwrap that adds support for ContextualError.
// Parameters:
//   - err: The error to unwrap
//
// Returns:
//   - error: The underlying error, or nil if there is none
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// GetCode returns the error code from an error.
// If the error is a ContextualError, it returns the code from the context.
// Otherwise, it returns an empty string.
// Parameters:
//   - err: The error to get the code from
//
// Returns:
//   - ErrorCode: The error code, or an empty string if not available
func GetCode(err error) ErrorCode {
	if err == nil {
		return ""
	}

	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		return contextualErr.Context.Code
	}

	return ""
}

// GetHTTPStatusFromError returns the HTTP status code from an error.
// If the error is a ContextualError, it returns the HTTP status from the context.
// Otherwise, it returns 0.
// Parameters:
//   - err: The error to get the HTTP status from
//
// Returns:
//   - int: The HTTP status code, or 0 if not available
func GetHTTPStatusFromError(err error) int {
	if err == nil {
		return 0
	}

	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		return contextualErr.Context.HTTPStatus
	}

	return 0
}

// ToJSON converts an error to a JSON string.
// If the error is a BaseError, it uses the MarshalJSON method.
// Otherwise, it creates a simple JSON object with the error message.
// Parameters:
//   - err: The error to convert to JSON
//
// Returns:
//   - string: The JSON representation of the error
func ToJSON(err error) string {
	if err == nil {
		return "{}"
	}

	var baseErr *BaseError
	if errors.As(err, &baseErr) {
		// Create a JSON object with the error details
		jsonObj := map[string]interface{}{
			"message": baseErr.Error(),
			"code":    baseErr.Code,
		}

		// Add other fields if available
		if baseErr.Operation != "" {
			jsonObj["operation"] = baseErr.Operation
		}
		if baseErr.Source != "" {
			jsonObj["source"] = baseErr.Source
		}
		if baseErr.Line > 0 {
			jsonObj["line"] = baseErr.Line
		}
		if baseErr.Details != nil && len(baseErr.Details) > 0 {
			jsonObj["details"] = baseErr.Details
		}

		jsonBytes, jsonErr := json.Marshal(jsonObj)
		if jsonErr != nil {
			return fmt.Sprintf(`{"message":"Error marshaling error to JSON: %s"}`, jsonErr.Error())
		}
		return string(jsonBytes)
	}

	// For non-BaseError errors, create a simple JSON object
	jsonBytes, jsonErr := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})
	if jsonErr != nil {
		return fmt.Sprintf(`{"message":"Error marshaling error to JSON: %s"}`, jsonErr.Error())
	}
	return string(jsonBytes)
}

// WrapWithOperation wraps an error with an operation name.
// It preserves the error chain and adds source location information.
// Parameters:
//   - err: The error to wrap
//   - operation: The name of the operation that failed
//
// Returns:
//   - error: A new error that wraps the original error
func WrapWithOperation(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Get caller information
	source, line := GetCallerInfo(1)

	// If the error is already a ContextualError, update its context
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		// Create a new ContextualError with the updated operation
		return &ContextualError{
			Original: contextualErr.Original,
			Context: ErrorContext{
				Operation:  operation,
				Source:     source,
				Line:       line,
				Code:       contextualErr.Context.Code,
				HTTPStatus: contextualErr.Context.HTTPStatus,
				Details:    contextualErr.Context.Details,
			},
		}
	}

	// Create a new ContextualError
	return &ContextualError{
		Original: err,
		Context: ErrorContext{
			Operation: operation,
			Source:    source,
			Line:      line,
		},
	}
}

// WithDetails adds details to an error.
// It preserves the error chain and adds additional context information.
// Parameters:
//   - err: The error to enhance
//   - details: A map of additional details to add to the error
//
// Returns:
//   - error: A new error with the added details
func WithDetails(err error, details map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// Get caller information
	source, line := GetCallerInfo(1)

	// If the error is already a ContextualError, update its context
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		// Merge the details
		newDetails := make(map[string]interface{})
		if contextualErr.Context.Details != nil {
			for k, v := range contextualErr.Context.Details {
				newDetails[k] = v
			}
		}
		for k, v := range details {
			newDetails[k] = v
		}

		// Create a new ContextualError with the merged details
		return &ContextualError{
			Original: contextualErr.Original,
			Context: ErrorContext{
				Operation:  contextualErr.Context.Operation,
				Source:     source,
				Line:       line,
				Code:       contextualErr.Context.Code,
				HTTPStatus: contextualErr.Context.HTTPStatus,
				Details:    newDetails,
			},
		}
	}

	// Create a new ContextualError
	return &ContextualError{
		Original: err,
		Context: ErrorContext{
			Source:  source,
			Line:    line,
			Details: details,
		},
	}
}
