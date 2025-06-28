// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
)

func BenchmarkErrorCodeGrouper_Group(b *testing.B) {
	code := core.NotFoundCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.Group(code)
	}
}

func BenchmarkErrorCodeGrouper_IsResourceError(b *testing.B) {
	code := core.NotFoundCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.IsResourceError(code)
	}
}

func BenchmarkErrorCodeGrouper_IsInputError(b *testing.B) {
	code := core.InvalidInputCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.IsInputError(code)
	}
}

func BenchmarkErrorCodeGrouper_IsSystemError(b *testing.B) {
	code := core.TimeoutCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.IsSystemError(code)
	}
}

func BenchmarkErrorCodeGrouper_IsSecurityError(b *testing.B) {
	code := core.UnauthorizedCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.IsSecurityError(code)
	}
}

func BenchmarkErrorCodeGrouper_IsExternalError(b *testing.B) {
	code := core.ExternalServiceErrorCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.IsExternalError(code)
	}
}

func BenchmarkErrorCodeGrouper_IsBusinessError(b *testing.B) {
	code := core.BusinessRuleViolationCode
	grouper := NewErrorCodeGrouper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.IsBusinessError(code)
	}
}
