// Copyright (c) 2025 A Bit of Help, Inc.

package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringValueObject_String(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Regular string", "test", "test"},
		{"String with whitespace", " test ", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo := NewStringValueObject(tt.value)
			assert.Equal(t, tt.expected, vo.String())
		})
	}
}

func TestStringValueObject_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Whitespace string", "   ", true},
		{"Non-empty string", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo := NewStringValueObject(tt.value)
			assert.Equal(t, tt.expected, vo.IsEmpty())
		})
	}
}

func TestStringValueObject_Equals(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{"Equal strings", "test", "test", true},
		{"Different strings", "test", "other", false},
		{"Case sensitive", "Test", "test", false},
		{"Empty strings", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo1 := NewStringValueObject(tt.value1)
			vo2 := NewStringValueObject(tt.value2)
			assert.Equal(t, tt.expected, vo1.Equals(vo2))
		})
	}
}

func TestStringValueObject_EqualsIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{"Equal strings", "test", "test", true},
		{"Different strings", "test", "other", false},
		{"Case insensitive", "Test", "test", true},
		{"Empty strings", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo1 := NewStringValueObject(tt.value1)
			vo2 := NewStringValueObject(tt.value2)
			assert.Equal(t, tt.expected, vo1.EqualsIgnoreCase(vo2))
		})
	}
}

func TestStringValueObject_Value(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Regular string", "test", "test"},
		{"String with whitespace", " test ", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo := NewStringValueObject(tt.value)
			assert.Equal(t, tt.expected, vo.Value())
		})
	}
}

func TestStringValueObject_WithValue(t *testing.T) {
	tests := []struct {
		name         string
		initialValue string
		newValue     string
		expected     string
	}{
		{"Change value", "test", "new", "new"},
		{"Change to empty", "test", "", ""},
		{"Change from empty", "", "new", "new"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo := NewStringValueObject(tt.initialValue)
			newVo := vo.WithValue(tt.newValue)

			// Original should be unchanged
			assert.Equal(t, tt.initialValue, vo.Value())

			// New should have new value
			assert.Equal(t, tt.expected, newVo.Value())
		})
	}
}
