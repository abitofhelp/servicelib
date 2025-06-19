// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
)

func TestNewPhone(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		expected    string
		expectError bool
	}{
		{"Valid US Phone", "123-456-7890", "123-456-7890", false},
		{"Valid US Phone with Parentheses", "(123) 456-7890", "(123) 456-7890", false},
		{"Valid US Phone with Dots", "123.456.7890", "123.456.7890", false},
		{"Valid US Phone with Spaces", "123 456 7890", "123 456 7890", false},
		{"Valid International Phone", "+1234567890", "+1234567890", false},
		{"Valid International Phone with Format", "+1 (123) 456-7890", "+1 (123) 456-7890", false},
		{"Empty Phone", "", "", false}, // Empty is allowed
		{"Invalid Phone - Letters", "abc-def-ghij", "", true},
		{"Invalid Phone - Too Short", "123-456", "", true},
		{"Invalid Phone - Wrong Format", "1234-567-890", "", true},
		{"Phone with Whitespace", " 123-456-7890 ", "123-456-7890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.phone)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if phone.String() != tt.expected {
					t.Errorf("Expected phone %s, got %s", tt.expected, phone.String())
				}
			}
		})
	}
}

func TestPhone_String(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{"Regular Phone", "123-456-7890", "123-456-7890"},
		{"Empty Phone", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, _ := NewPhone(tt.phone)
			if phone.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, phone.String())
			}
		})
	}
}

func TestPhone_Equals(t *testing.T) {
	tests := []struct {
		name     string
		phone1   string
		phone2   string
		expected bool
	}{
		{"Same Phone", "123-456-7890", "123-456-7890", true},
		{"Different Format - Dashes vs Dots", "123-456-7890", "123.456.7890", true},
		{"Different Format - Dashes vs Spaces", "123-456-7890", "123 456 7890", true},
		{"Different Format - Dashes vs Parentheses", "123-456-7890", "(123) 456-7890", true},
		{"Different Format - With/Without Country Code", "+1-123-456-7890", "123-456-7890", false}, // Different numbers
		{"Different Phone", "123-456-7890", "987-654-3210", false},
		{"Empty Phones", "", "", true},
		{"One Empty Phone", "123-456-7890", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone1, _ := NewPhone(tt.phone1)
			phone2, _ := NewPhone(tt.phone2)

			if phone1.Equals(phone2) != tt.expected {
				t.Errorf("Expected Equals to return %v for %s and %s", tt.expected, tt.phone1, tt.phone2)
			}
		})
	}
}

func TestPhone_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"Empty Phone", "", true},
		{"Non-Empty Phone", "123-456-7890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, _ := NewPhone(tt.phone)
			if phone.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %s", tt.expected, tt.phone)
			}
		})
	}
}