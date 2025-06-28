// Copyright (c) 2025 A Bit of Help, Inc.

package contact

import (
	"testing"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expected    string
		expectError bool
	}{
		{"Valid Email", "test@example.com", "test@example.com", false},
		{"Valid Email with Spaces", " test@example.com ", "test@example.com", false},
		{"Valid Email with Uppercase", "TEST@EXAMPLE.COM", "TEST@EXAMPLE.COM", false},
		{"Valid Email with Name", "Test User <test@example.com>", "Test User <test@example.com>", false},
		{"Empty Email", "", "", false},
		{"Invalid Email - No @", "testexample.com", "", true},
		{"Invalid Email - No Domain", "test@", "", true},
		{"Invalid Email - Invalid Characters", "test!@example.com", "", true},
		{"Invalid Email - Multiple @", "test@example@com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if email.String() != tt.expected {
					t.Errorf("Expected email %s, got %s", tt.expected, email.String())
				}
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{"Regular Email", "test@example.com", "test@example.com"},
		{"Empty Email", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, _ := NewEmail(tt.email)
			if email.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, email.String())
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name     string
		email1   string
		email2   string
		expected bool
	}{
		{"Same Email", "test@example.com", "test@example.com", true},
		{"Different Case", "test@example.com", "TEST@EXAMPLE.COM", true},
		{"Different Email", "test@example.com", "other@example.com", false},
		{"Empty Emails", "", "", true},
		{"One Empty Email", "test@example.com", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email1, _ := NewEmail(tt.email1)
			email2, _ := NewEmail(tt.email2)

			if email1.Equals(email2) != tt.expected {
				t.Errorf("Expected Equals to return %v for %s and %s", tt.expected, tt.email1, tt.email2)
			}
		})
	}
}

func TestEmail_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Empty Email", "", true},
		{"Non-Empty Email", "test@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, _ := NewEmail(tt.email)
			if email.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %s", tt.expected, tt.email)
			}
		})
	}
}

func TestEmail_Validate(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{"Valid Email", "test@example.com", false},
		{"Empty Email", "", false},
		{"Invalid Email", "testexample.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For invalid emails, we need to directly call the validation function
			// because NewEmail would already return an error
			var err error
			if tt.expectError {
				err = base.ValidateEmail(tt.email)
			} else {
				email, _ := NewEmail(tt.email)
				err = email.Validate()
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
