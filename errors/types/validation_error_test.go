// Copyright (c) 2025 A Bit of Help, Inc.

package types

import (
	"net/http"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

// TestValidationError_Error tests the Error method of ValidationError
func TestValidationError_Error(t *testing.T) {
	// Create a validation error with a field
	err1 := &ValidationError{
		Msg:   "test message",
		Field: "testField",
	}

	// Test the Error method with a field
	assert.Equal(t, "validation error: test message (field: testField)", err1.Error())

	// Create a validation error without a field
	err2 := &ValidationError{
		Msg: "test message",
	}

	// Test the Error method without a field
	assert.Equal(t, "validation error: test message", err2.Error())
}

// TestValidationError_Code tests the Code method of ValidationError
func TestValidationError_Code(t *testing.T) {
	// Create a validation error
	err := &ValidationError{
		Msg:   "test message",
		Field: "testField",
	}

	// Test the Code method
	assert.Equal(t, string(core.ValidationErrorCode), err.Code())
}

// TestValidationError_HTTPStatus tests the HTTPStatus method of ValidationError
func TestValidationError_HTTPStatus(t *testing.T) {
	// Create a validation error
	err := &ValidationError{
		Msg:   "test message",
		Field: "testField",
	}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus())
}

// TestValidationError_IsValidationError tests the IsValidationError method of ValidationError
func TestValidationError_IsValidationError(t *testing.T) {
	// Create a validation error
	err := &ValidationError{
		Msg:   "test message",
		Field: "testField",
	}

	// Test the IsValidationError method
	assert.True(t, err.IsValidationError())
}

// TestNewValidationError tests the NewValidationError function
func TestNewValidationError(t *testing.T) {
	// Create a validation error using NewValidationError
	err := NewValidationError("test message")

	// Test the error properties
	assert.Equal(t, "test message", err.Msg)
	assert.Equal(t, "", err.Field)
}

// TestNewFieldValidationError tests the NewFieldValidationError function
func TestNewFieldValidationError(t *testing.T) {
	// Create a validation error using NewFieldValidationError
	err := NewFieldValidationError("test message", "testField")

	// Test the error properties
	assert.Equal(t, "test message", err.Msg)
	assert.Equal(t, "testField", err.Field)
}

// TestValidationErrors_Error tests the Error method of ValidationErrors
func TestValidationErrors_Error(t *testing.T) {
	// Create a ValidationErrors with no errors
	errs1 := &ValidationErrors{}

	// Test the Error method with no errors
	assert.Equal(t, "", errs1.Error())

	// Create a ValidationErrors with one error
	errs2 := &ValidationErrors{
		Errors: []*ValidationError{
			{
				Msg:   "test message",
				Field: "testField",
			},
		},
	}

	// Test the Error method with one error
	assert.Equal(t, "validation error: test message (field: testField)", errs2.Error())

	// Create a ValidationErrors with multiple errors
	errs3 := &ValidationErrors{
		Errors: []*ValidationError{
			{
				Msg:   "test message 1",
				Field: "testField1",
			},
			{
				Msg:   "test message 2",
				Field: "testField2",
			},
		},
	}

	// Test the Error method with multiple errors
	assert.Equal(t, "2 validation errors: validation error: test message 1 (field: testField1); validation error: test message 2 (field: testField2)", errs3.Error())
}

// TestValidationErrors_Code tests the Code method of ValidationErrors
func TestValidationErrors_Code(t *testing.T) {
	// Create a ValidationErrors
	errs := &ValidationErrors{}

	// Test the Code method
	assert.Equal(t, string(core.ValidationErrorCode), errs.Code())
}

// TestValidationErrors_HTTPStatus tests the HTTPStatus method of ValidationErrors
func TestValidationErrors_HTTPStatus(t *testing.T) {
	// Create a ValidationErrors
	errs := &ValidationErrors{}

	// Test the HTTPStatus method
	assert.Equal(t, http.StatusBadRequest, errs.HTTPStatus())
}

// TestValidationErrors_IsValidationError tests the IsValidationError method of ValidationErrors
func TestValidationErrors_IsValidationError(t *testing.T) {
	// Create a ValidationErrors
	errs := &ValidationErrors{}

	// Test the IsValidationError method
	assert.True(t, errs.IsValidationError())
}

// TestNewValidationErrors tests the NewValidationErrors function
func TestNewValidationErrors(t *testing.T) {
	// Create validation errors
	err1 := &ValidationError{
		Msg:   "test message 1",
		Field: "testField1",
	}
	err2 := &ValidationError{
		Msg:   "test message 2",
		Field: "testField2",
	}

	// Create a ValidationErrors using NewValidationErrors
	errs := NewValidationErrors(err1, err2)

	// Test the error properties
	assert.Equal(t, 2, len(errs.Errors))
	assert.Equal(t, err1, errs.Errors[0])
	assert.Equal(t, err2, errs.Errors[1])
}

// TestValidationErrors_AddError tests the AddError method of ValidationErrors
func TestValidationErrors_AddError(t *testing.T) {
	// Create a ValidationErrors
	errs := &ValidationErrors{}

	// Create validation errors
	err1 := &ValidationError{
		Msg:   "test message 1",
		Field: "testField1",
	}
	err2 := &ValidationError{
		Msg:   "test message 2",
		Field: "testField2",
	}

	// Add errors
	errs.AddError(err1)
	errs.AddError(err2)

	// Test the error properties
	assert.Equal(t, 2, len(errs.Errors))
	assert.Equal(t, err1, errs.Errors[0])
	assert.Equal(t, err2, errs.Errors[1])
}

// TestValidationErrors_HasErrors tests the HasErrors method of ValidationErrors
func TestValidationErrors_HasErrors(t *testing.T) {
	// Create a ValidationErrors with no errors
	errs1 := &ValidationErrors{}

	// Test the HasErrors method with no errors
	assert.False(t, errs1.HasErrors())

	// Create a ValidationErrors with errors
	errs2 := &ValidationErrors{
		Errors: []*ValidationError{
			{
				Msg:   "test message",
				Field: "testField",
			},
		},
	}

	// Test the HasErrors method with errors
	assert.True(t, errs2.HasErrors())
}