// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"errors"
	"strconv"
	"testing"
)

func BenchmarkAddDetail(b *testing.B) {
	originalErr := errors.New("original error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddDetail(originalErr, "field", "value")
	}
}

func BenchmarkAddDetails(b *testing.B) {
	originalErr := errors.New("original error")
	details := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddDetails(originalErr, details)
	}
}

func BenchmarkGetDetail(b *testing.B) {
	originalErr := errors.New("original error")
	err := AddDetail(originalErr, "field", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetDetail(err, "field")
	}
}

func BenchmarkGetDetails(b *testing.B) {
	originalErr := errors.New("original error")
	details := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	err := AddDetails(originalErr, details)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetDetails(err)
	}
}

func BenchmarkAddDetailMultiple(b *testing.B) {
	originalErr := errors.New("original error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			err := AddDetail(originalErr, "field_"+strconv.Itoa(j), "value_"+strconv.Itoa(j))
			originalErr = err
		}
	}
}
