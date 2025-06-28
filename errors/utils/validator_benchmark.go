// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
)

func BenchmarkValidator_Validate(b *testing.B) {
	v := NewErrorCodeValidator()
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Validate(code)
	}
}

func BenchmarkValidator_GetCodeName(b *testing.B) {
	v := NewErrorCodeValidator()
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.GetCodeName(code)
	}
}

func BenchmarkValidator_IsValid(b *testing.B) {
	v := NewErrorCodeValidator()
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.IsValid(code)
	}
}

func BenchmarkValidator_Validate_Invalid(b *testing.B) {
	v := NewErrorCodeValidator()
	code := core.ErrorCode("INVALID_CODE")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Validate(code)
	}
}

func BenchmarkValidator_GetCodeName_Invalid(b *testing.B) {
	v := NewErrorCodeValidator()
	code := core.ErrorCode("INVALID_CODE")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.GetCodeName(code)
	}
}
