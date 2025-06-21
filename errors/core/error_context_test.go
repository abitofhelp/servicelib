// Copyright (c) 2024 A Bit of Help, Inc.

package core

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextualError_Error(t *testing.T) {
	tests := []struct {
		name     string
		operation string
		original error
		want     string
	}{
		{
			name:     "with operation",
			operation: "create_user",
			original:  errors.New("failed to create user"),
			want:     "create_user: failed to create user",
		},
		{
			name:     "without operation",
			operation: "",
			original:  errors.New("failed to create user"),
			want:     "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ErrorContext{
				Operation: tt.operation,
			}
			err := &ContextualError{
				Original: tt.original,
				Context:  ctx,
			}
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestContextualError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := &ContextualError{
		Original: originalErr,
		Context:  ErrorContext{},
	}

	unwrapped := err.Unwrap()
	assert.Equal(t, originalErr, unwrapped)
}

func TestContextualError_Code(t *testing.T) {
	code := "INVALID_INPUT"
	err := &ContextualError{
		Context: ErrorContext{
			Code: ErrorCode(code),
		},
	}

	assert.Equal(t, ErrorCode(code), err.Code())
}

func TestContextualError_HTTPStatus(t *testing.T) {
	status := 400
	err := &ContextualError{
		Context: ErrorContext{
			HTTPStatus: status,
		},
	}

	assert.Equal(t, status, err.HTTPStatus())
}

func TestContextualError_MarshalJSON(t *testing.T) {
	ctx := ErrorContext{
		Operation: "create_user",
		Code:      "INVALID_INPUT",
		HTTPStatus: 400,
		Details: map[string]interface{}{
			"field": "email",
			"reason": "invalid format",
		},
	}
	err := &ContextualError{
		Context: ctx,
	}

	jsonBytes, marshalErr := err.MarshalJSON()
	require.NoError(t, marshalErr)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &result))

	assert.Equal(t, ctx.Operation, result["operation"])
	assert.Equal(t, string(ctx.Code), result["code"])
	assert.Equal(t, float64(ctx.HTTPStatus), result["http_status"])
	assert.Equal(t, ctx.Details, result["details"])
}

func TestGetCallerInfo(t *testing.T) {
	file, line := getCallerInfo(1)
	assert.NotEmpty(t, file)
	assert.NotZero(t, line)
}
