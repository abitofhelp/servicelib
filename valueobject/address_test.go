// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"strings"
	"testing"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Valid Address", "123 Main St, City, State 12345", "123 Main St, City, State 12345", false},
		{"Valid Address with Spaces", " 123 Main St, City, State 12345 ", "123 Main St, City, State 12345", false},
		{"Empty Address", "", "", false}, // Empty is allowed
		{"Too Short Address", "123", "", true},
		{"Minimum Length Address", "12345", "12345", false},
		{"Too Long Address", strings.Repeat("a", 201), "", true},
		{"Maximum Length Address", strings.Repeat("a", 200), strings.Repeat("a", 200), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := NewAddress(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if address.String() != tt.expected {
					t.Errorf("Expected address %s, got %s", tt.expected, address.String())
				}
			}
		})
	}
}

func TestAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Regular Address", "123 Main St, City, State 12345", "123 Main St, City, State 12345"},
		{"Empty Address", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, _ := NewAddress(tt.input)
			if address.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, address.String())
			}
		})
	}
}

func TestAddress_Equals(t *testing.T) {
	tests := []struct {
		name     string
		address1 string
		address2 string
		expected bool
	}{
		{"Same Address", "123 Main St", "123 Main St", true},
		{"Different Case", "123 main st", "123 MAIN ST", true},
		{"Different Address", "123 Main St", "456 Oak Ave", false},
		{"Empty Addresses", "", "", true},
		{"One Empty Address", "123 Main St", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address1, _ := NewAddress(tt.address1)
			address2, _ := NewAddress(tt.address2)

			if address1.Equals(address2) != tt.expected {
				t.Errorf("Expected Equals to return %v for %s and %s", tt.expected, tt.address1, tt.address2)
			}
		})
	}
}

func TestAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected bool
	}{
		{"Empty Address", "", true},
		{"Non-Empty Address", "123 Main St", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, _ := NewAddress(tt.address)
			if address.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %s", tt.expected, tt.address)
			}
		})
	}
}
