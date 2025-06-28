// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
)

func BenchmarkErrorCodeCategorizer_Categorize(b *testing.B) {
	c := NewErrorCodeCategorizer()
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Categorize(code)
	}
}

func BenchmarkErrorCodeCategorizer_IsClientError(b *testing.B) {
	c := NewErrorCodeCategorizer()
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.IsClientError(code)
	}
}

func BenchmarkErrorCodeCategorizer_IsServerError(b *testing.B) {
	c := NewErrorCodeCategorizer()
	code := core.InternalErrorCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.IsServerError(code)
	}
}

func BenchmarkErrorCodeCategorizer_IsSystemError(b *testing.B) {
	c := NewErrorCodeCategorizer()
	code := core.TimeoutCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.IsSystemError(code)
	}
}

func BenchmarkErrorCodeCategorizer_IsExternalError(b *testing.B) {
	c := NewErrorCodeCategorizer()
	code := core.ExternalServiceErrorCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.IsExternalError(code)
	}
}

func BenchmarkErrorCodeCategorizer_IsSecurityError(b *testing.B) {
	c := NewErrorCodeCategorizer()
	code := core.UnauthorizedCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.IsSecurityError(code)
	}
}
