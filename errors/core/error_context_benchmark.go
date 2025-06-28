// Copyright (c) 2024 A Bit of Help, Inc.

package core

import (
	"errors"
	"testing"
)

func BenchmarkContextualErrorCreation(b *testing.B) {
	originalErr := errors.New("original error")
	ctx := ErrorContext{
		Operation:  "create_user",
		Code:       "INVALID_INPUT",
		HTTPStatus: 400,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &ContextualError{
			Original: originalErr,
			Context:  ctx,
		}
	}
}

func BenchmarkErrorWrapping(b *testing.B) {
	originalErr := errors.New("original error")
	ctx := ErrorContext{
		Operation:  "create_user",
		Code:       "INVALID_INPUT",
		HTTPStatus: 400,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ce := &ContextualError{
			Original: originalErr,
			Context:  ctx,
		}
		_ = ce.Error()
	}
}

func BenchmarkErrorUnwrapping(b *testing.B) {
	originalErr := errors.New("original error")
	ctx := ErrorContext{
		Operation:  "create_user",
		Code:       "INVALID_INPUT",
		HTTPStatus: 400,
	}
	err := &ContextualError{
		Original: originalErr,
		Context:  ctx,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Unwrap()
	}
}

func BenchmarkErrorJSONMarshaling(b *testing.B) {
	originalErr := errors.New("original error")
	ctx := ErrorContext{
		Operation:  "create_user",
		Code:       "INVALID_INPUT",
		HTTPStatus: 400,
		Details: map[string]interface{}{
			"field":  "email",
			"reason": "invalid format",
		},
	}
	err := &ContextualError{
		Original: originalErr,
		Context:  ctx,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = err.MarshalJSON()
	}
}

func BenchmarkGetCallerInfo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetCallerInfo(1)
	}
}
