// Copyright (c) 2024 A Bit of Help, Inc.

package wrappers

import (
	"fmt"
	"github.com/abitofhelp/servicelib/errors/core"
)

// withContext wraps an error with contextual information.
func withContext(err error, operation string, code core.ErrorCode, httpStatus int, details map[string]interface{}) error {
	if err == nil {
		return nil
	}

	if ce, ok := err.(*core.ContextualError); ok {
		if operation != "" {
			ce.Context.Operation = operation
		}
		if code != "" {
			ce.Context.Code = code
		}
		if httpStatus != 0 {
			ce.Context.HTTPStatus = httpStatus
		}
		if details != nil {
			if ce.Context.Details == nil {
				ce.Context.Details = make(map[string]interface{})
			}
			for k, v := range details {
				ce.Context.Details[k] = v
			}
		}
		return ce
	}

	return &core.ContextualError{
		Original: err,
		Context: core.ErrorContext{
			Operation:  operation,
			Code:       code,
			HTTPStatus: httpStatus,
			Details:    details,
		},
	}
}

// WrapWithOperation wraps an error with an operation name and message.
func WrapWithOperation(err error, operation string, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	return withContext(err, operation, "", 0, map[string]interface{}{"message": msg})
}

// WithDetails adds details to an error.
func WithDetails(err error, details map[string]interface{}) error {
	if err == nil {
		return nil
	}

	return withContext(err, "", "", 0, details)
}
