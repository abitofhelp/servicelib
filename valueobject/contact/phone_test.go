// Copyright (c) 2025 A Bit of Help, Inc.

package contact

import (
	"testing"

	"github.com/abitofhelp/servicelib/valueobject/base"
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

func TestPhone_Validate(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		expectError bool
	}{
		{"Valid Phone", "123-456-7890", false},
		{"Empty Phone", "", false},
		{"Invalid Phone", "abc-def-ghij", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For invalid phones, we need to directly call the validation function
			// because NewPhone would already return an error
			var err error
			if tt.expectError {
				err = base.ValidatePhone(tt.phone)
			} else {
				phone, _ := NewPhone(tt.phone)
				err = phone.Validate()
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

func TestPhone_Format(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		format   string
		expected string
	}{
		{"E164 Format", "123-456-7890", "e164", "+1234567890"},
		{"E164 Format with Country Code", "+44 123-456-7890", "e164", "+44123456789"},
		{"National Format", "123-456-7890", "national", "(123) 456-7890"},
		{"International Format", "123-456-7890", "international", "+1 123 456 7890"},
		{"International Format with Country Code", "+44 123-456-7890", "international", "+44 123 456 7890"},
		{"Invalid Format", "123-456-7890", "invalid", "123-456-7890"},
		{"Empty Phone", "", "e164", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, _ := NewPhone(tt.phone)
			formatted := phone.Format(tt.format)

			if formatted != tt.expected {
				t.Errorf("Expected format %s to return %s, got %s", tt.format, tt.expected, formatted)
			}
		})
	}
}

func TestPhone_CountryCode(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{"US Phone", "123-456-7890", "1"},
		{"UK Phone", "+44 123-456-7890", "44"},
		{"Empty Phone", "", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, _ := NewPhone(tt.phone)
			countryCode := phone.CountryCode()

			if countryCode != tt.expected {
				t.Errorf("Expected country code %s, got %s", tt.expected, countryCode)
			}
		})
	}
}

func TestPhone_Normalized(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{"US Phone", "123-456-7890", "+1234567890"},
		{"UK Phone", "+44 123-456-7890", "+44123456789"},
		{"Empty Phone", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, _ := NewPhone(tt.phone)
			normalized := phone.Normalized()

			if normalized != tt.expected {
				t.Errorf("Expected normalized phone %s, got %s", tt.expected, normalized)
			}
		})
	}
}

func TestPhone_IsValidForCountry(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		countryCode string
		expected    bool
	}{
		{"US Phone - Valid", "123-456-7890", "1", true},
		{"US Phone - Invalid", "123-456-7890", "44", false},
		{"UK Phone - Valid", "+44 123-456-7890", "44", true},
		{"UK Phone - Invalid", "+44 123-456-7890", "1", false},
		{"Empty Phone", "", "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, _ := NewPhone(tt.phone)
			isValid := phone.IsValidForCountry(tt.countryCode)

			if isValid != tt.expected {
				t.Errorf("Expected IsValidForCountry to return %v for %s and %s", tt.expected, tt.phone, tt.countryCode)
			}
		})
	}
}
