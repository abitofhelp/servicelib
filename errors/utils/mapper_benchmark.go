// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
)

func BenchmarkErrorCodeMapper_ToHTTPStatus(b *testing.B) {
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorCodeMapper.ToHTTPStatus(code)
	}
}

func BenchmarkErrorCodeMapper_ToCategory(b *testing.B) {
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorCodeMapper.ToCategory(code)
	}
}

func BenchmarkErrorCodeMapper_ToMessage(b *testing.B) {
	code := core.NotFoundCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorCodeMapper.ToMessage(code)
	}
}

func BenchmarkErrorCodeMapper_ToHTTPStatusInvalidCode(b *testing.B) {
	code := core.ErrorCode("INVALID_CODE")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorCodeMapper.ToHTTPStatus(code)
	}
}

func BenchmarkErrorCodeMapper_ToCategoryInvalidCode(b *testing.B) {
	code := core.ErrorCode("INVALID_CODE")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorCodeMapper.ToCategory(code)
	}
}

func BenchmarkErrorCodeMapper_ToMessageInvalidCode(b *testing.B) {
	code := core.ErrorCode("INVALID_CODE")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorCodeMapper.ToMessage(code)
	}
}
