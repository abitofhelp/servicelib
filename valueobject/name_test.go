// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"strings"
	"testing"
)

func TestNewName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Valid Name", "John Doe", "John Doe", false},
		{"Valid Name with Spaces", " Jane Smith ", "Jane Smith", false},
		{"Empty Name", "", "", true},
		{"Too Long Name", strings.Repeat("a", 101), "", true},
		{"Max Length Name", strings.Repeat("a", 100), strings.Repeat("a", 100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := NewName(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if name.String() != tt.expected {
					t.Errorf("Expected name %s, got %s", tt.expected, name.String())
				}
			}
		})
	}
}

func TestName_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Regular Name", "John Doe", "John Doe"},
		{"Single Word Name", "John", "John"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, _ := NewName(tt.input)
			if name.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, name.String())
			}
		})
	}
}

func TestName_Equals(t *testing.T) {
	tests := []struct {
		name     string
		name1    string
		name2    string
		expected bool
	}{
		{"Same Name", "John Doe", "John Doe", true},
		{"Different Case", "john doe", "JOHN DOE", true},
		{"Different Name", "John Doe", "Jane Smith", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name1, _ := NewName(tt.name1)
			name2, _ := NewName(tt.name2)

			if name1.Equals(name2) != tt.expected {
				t.Errorf("Expected Equals to return %v for %s and %s", tt.expected, tt.name1, tt.name2)
			}
		})
	}
}

func TestName_IsEmpty(t *testing.T) {
	// Create a valid name
	validName, _ := NewName("John Doe")

	// Create an empty name (this would normally fail validation, but we can test with the type directly)
	emptyName := Name("")

	if validName.IsEmpty() {
		t.Errorf("Expected valid name to not be empty")
	}

	if !emptyName.IsEmpty() {
		t.Errorf("Expected empty name to be empty")
	}
}
