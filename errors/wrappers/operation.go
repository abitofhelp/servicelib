// Copyright (c) 2024 A Bit of Help, Inc.

package wrappers

import (
	"fmt"
	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/abitofhelp/servicelib/errors/utils"
)

// OperationError wraps an error with operation-specific context
func OperationError(err error, operation string, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	
	msg := fmt.Sprintf(format, args...)
	return withContext(err, operation, "", 0, map[string]interface{}{
		"message": msg,
		"operation": operation,
		"package": utils.GetCallerPackage(1),
	})
}

// OperationNotFound creates a not found error for an operation
func OperationNotFound(operation string, resourceType string, id string) error {
	return OperationError(
		core.ErrNotFound,
		operation,
		"%s not found: %s",
		resourceType,
		id,
	)
}

// OperationInvalidInput creates an invalid input error for an operation
func OperationInvalidInput(operation string, format string, args ...interface{}) error {
	return OperationError(
		core.ErrInvalidInput,
		operation,
		format,
		args...,
	)
}

// OperationInternal creates an internal error for an operation
func OperationInternal(operation string, err error, format string, args ...interface{}) error {
	return OperationError(
		err,
		operation,
		format,
		args...,
	)
}
