// Copyright (c) 2025 A Bit of Help, Inc.

package network

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIPAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Valid Value", "valid", "valid", false},
		{"Valid Value with Spaces", " valid ", "valid", false},
		{"Empty Value", "", "", false}, // Empty is allowed
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewIPAddress(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, value.String())
			}
		})
	}
}

func TestIPAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Regular Value", "test", "test"},
		{"Empty Value", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			assert.Equal(t, tt.expected, value.String())
		})
	}
}

func TestIPAddress_Equals(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{"Same Values", "test", "test", true},
		{"Different Case", "test", "TEST", true},
		{"Different Values", "test", "other", false},
		{"Empty Values", "", "", true},
		{"One Empty Value", "test", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1, _ := NewIPAddress(tt.value1)
			value2, _ := NewIPAddress(tt.value2)

			assert.Equal(t, tt.expected, value1.Equals(value2))
		})
	}
}

func TestIPAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"Empty Value", "", true},
		{"Non-Empty Value", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.value)
			assert.Equal(t, tt.expected, value.IsEmpty())
		})
	}
}

func TestIPAddress_Validate(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"Valid Value", "valid", false},
		{"Empty Value", "", false}, // Empty is allowed
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.value)
			err := value.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
