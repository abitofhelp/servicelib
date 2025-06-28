// Copyright (c) 2025 A Bit of Help, Inc.

package wrappers

import (
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

// TestOperationError tests the OperationError function
func TestOperationError(t *testing.T) {
	// Test with nil error
	err := OperationError(nil, "TestOperation", "test message")
	assert.Nil(t, err)

	// Test with standard error
	originalErr := errors.New("original error")
	err = OperationError(originalErr, "TestOperation", "test message")
	
	// Check that the error is a ContextualError
	ce, ok := err.(*core.ContextualError)
	assert.True(t, ok)
	
	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "test message", ce.Context.Details["message"])
	assert.Equal(t, "TestOperation", ce.Context.Details["operation"])
	assert.NotEmpty(t, ce.Context.Details["package"])

	// Test with format string
	err = OperationError(originalErr, "TestOperation", "test message with %s", "format")
	
	// Check that the error is a ContextualError
	ce, ok = err.(*core.ContextualError)
	assert.True(t, ok)
	
	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "test message with format", ce.Context.Details["message"])
}

// TestOperationNotFound tests the OperationNotFound function
func TestOperationNotFound(t *testing.T) {
	// Test with standard parameters
	err := OperationNotFound("TestOperation", "User", "123")
	
	// Check that the error is a ContextualError
	ce, ok := err.(*core.ContextualError)
	assert.True(t, ok)
	
	// Check the error properties
	assert.Equal(t, core.ErrNotFound, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "User not found: 123", ce.Context.Details["message"])
	assert.Equal(t, "TestOperation", ce.Context.Details["operation"])
	assert.NotEmpty(t, ce.Context.Details["package"])
}

// TestOperationInvalidInput tests the OperationInvalidInput function
func TestOperationInvalidInput(t *testing.T) {
	// Test with standard parameters
	err := OperationInvalidInput("TestOperation", "Invalid input: %s", "test")
	
	// Check that the error is a ContextualError
	ce, ok := err.(*core.ContextualError)
	assert.True(t, ok)
	
	// Check the error properties
	assert.Equal(t, core.ErrInvalidInput, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "Invalid input: test", ce.Context.Details["message"])
	assert.Equal(t, "TestOperation", ce.Context.Details["operation"])
	assert.NotEmpty(t, ce.Context.Details["package"])
}

// TestOperationInternal tests the OperationInternal function
func TestOperationInternal(t *testing.T) {
	// Test with standard parameters
	originalErr := errors.New("original error")
	err := OperationInternal("TestOperation", originalErr, "Internal error: %s", "test")
	
	// Check that the error is a ContextualError
	ce, ok := err.(*core.ContextualError)
	assert.True(t, ok)
	
	// Check the error properties
	assert.Equal(t, originalErr, ce.Original)
	assert.Equal(t, "TestOperation", ce.Context.Operation)
	assert.Equal(t, "Internal error: test", ce.Context.Details["message"])
	assert.Equal(t, "TestOperation", ce.Context.Details["operation"])
	assert.NotEmpty(t, ce.Context.Details["package"])
}