// Copyright (c) 2025 A Bit of Help, Inc.

package wrappers

import (
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

// TestWrapWithOperation tests the WrapWithOperation function
func TestWrapWithOperation(t *testing.T) {
	// Test with nil error
	err := WrapWithOperation(nil, "TestOperation", "test message")
	assert.Nil(t, err)

	// Test with standard error
	originalErr := errors.New("original error")
	err = WrapWithOperation(originalErr, "TestOperation", "test message")

	// Check that the error is a ContextualError
	ce, ok := err.(*core.ContextualError)
	assert.True(t, ok)

	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "test message", ce.Context.Details["message"])

	// Test with format string
	err = WrapWithOperation(originalErr, "TestOperation", "test message with %s", "format")

	// Check that the error is a ContextualError
	ce, ok = err.(*core.ContextualError)
	assert.True(t, ok)

	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "test message with format", ce.Context.Details["message"])

	// Test with existing ContextualError
	existingCE := &core.ContextualError{
		Original: originalErr,
		Context: core.ErrorContext{
			Operation:  "ExistingOperation",
			Code:       "EXISTING_CODE",
			HTTPStatus: 500,
			Details:    map[string]interface{}{"existing": "value"},
		},
	}

	err = WrapWithOperation(existingCE, "NewOperation", "new message")

	// Check that the error is a ContextualError
	ce, ok = err.(*core.ContextualError)
	assert.True(t, ok)

	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "NewOperation", ce.Context.Operation)
	assert.Equal(t, core.ErrorCode("EXISTING_CODE"), ce.Context.Code)
	assert.Equal(t, 500, ce.Context.HTTPStatus)
	assert.Equal(t, "value", ce.Context.Details["existing"])
	assert.Equal(t, "new message", ce.Context.Details["message"])
}

// TestWithDetails tests the WithDetails function
func TestWithDetails(t *testing.T) {
	// Test with nil error
	err := WithDetails(nil, map[string]interface{}{"key": "value"})
	assert.Nil(t, err)

	// Test with standard error
	originalErr := errors.New("original error")
	details := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	err = WithDetails(originalErr, details)

	// Check that the error is a ContextualError
	ce, ok := err.(*core.ContextualError)
	assert.True(t, ok)

	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "", ce.Context.Operation)
	assert.Equal(t, core.ErrorCode(""), ce.Context.Code)
	assert.Equal(t, 0, ce.Context.HTTPStatus)
	assert.Equal(t, "value1", ce.Context.Details["key1"])
	assert.Equal(t, 42, ce.Context.Details["key2"])

	// Test with existing ContextualError
	existingCE := &core.ContextualError{
		Original: originalErr,
		Context: core.ErrorContext{
			Operation:  "ExistingOperation",
			Code:       "EXISTING_CODE",
			HTTPStatus: 500,
			Details:    map[string]interface{}{"existing": "value"},
		},
	}

	newDetails := map[string]interface{}{
		"key1": "new value",
		"key3": true,
	}

	err = WithDetails(existingCE, newDetails)

	// Check that the error is a ContextualError
	ce, ok = err.(*core.ContextualError)
	assert.True(t, ok)

	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "ExistingOperation", ce.Context.Operation)
	assert.Equal(t, core.ErrorCode("EXISTING_CODE"), ce.Context.Code)
	assert.Equal(t, 500, ce.Context.HTTPStatus)
	assert.Equal(t, "value", ce.Context.Details["existing"])
	assert.Equal(t, "new value", ce.Context.Details["key1"])
	assert.Equal(t, true, ce.Context.Details["key3"])
}
