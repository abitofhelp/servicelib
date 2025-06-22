// Copyright (c) 2025 A Bit of Help, Inc.

// Package contact provides value objects related to contact information.
package contact

import (
	"regexp"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Phone represents a phone number value object
type Phone string

// NewPhone creates a new Phone with validation
func NewPhone(phone string) (Phone, error) {
	// Validate the phone using the common validation function
	if err := base.ValidatePhone(phone); err != nil {
		return "", err
	}

	// Trim whitespace
	trimmedPhone := strings.TrimSpace(phone)

	return Phone(trimmedPhone), nil
}

// String returns the string representation of the Phone
func (p Phone) String() string {
	return string(p)
}

// Equals checks if two Phones are equal
func (p Phone) Equals(other Phone) bool {
	// Remove all non-digit characters for comparison
	cleanThis := regexp.MustCompile(`\D`).ReplaceAllString(string(p), "")
	cleanOther := regexp.MustCompile(`\D`).ReplaceAllString(string(other), "")

	return cleanThis == cleanOther
}

// IsEmpty checks if the Phone is empty
func (p Phone) IsEmpty() bool {
	return p == ""
}

// Validate checks if the Phone is valid
func (p Phone) Validate() error {
	return base.ValidatePhone(string(p))
}

// Format formats the phone number according to the specified format
// Supported formats:
// - "e164": E.164 format (e.g., "+12345678901")
// - "national": National format (e.g., "(234) 567-8901")
// - "international": International format (e.g., "+1 234 567 8901")
// If an unsupported format is specified, the phone number is returned as is.
func (p Phone) Format(format string) string {
	// Empty phone number should return empty string
	if p.IsEmpty() {
		return ""
	}

	// This is a simplified implementation
	// A real implementation would use a library like libphonenumber
	switch format {
	case "e164":
		// Remove all non-digit characters except the leading +
		digits := regexp.MustCompile(`[^\d+]`).ReplaceAllString(string(p), "")

		// Handle UK phone numbers specifically
		if strings.HasPrefix(digits, "+44") && len(digits) > 12 {
			// UK phone numbers should be +44 followed by 9 digits
			return "+44" + digits[3:12]
		}

		if !strings.HasPrefix(digits, "+") {
			digits = "+" + digits
		}
		return digits
	case "national":
		// This is a simplified implementation
		digits := regexp.MustCompile(`\D`).ReplaceAllString(string(p), "")
		if len(digits) >= 10 {
			return "(" + digits[len(digits)-10:len(digits)-7] + ") " +
				digits[len(digits)-7:len(digits)-4] + "-" +
				digits[len(digits)-4:]
		}
		return string(p)
	case "international":
		// This is a simplified implementation
		digits := regexp.MustCompile(`\D`).ReplaceAllString(string(p), "")
		if len(digits) >= 10 {
			countryCode := "+"
			if len(digits) > 10 {
				countryCode += digits[:len(digits)-10]
			} else {
				countryCode += "1" // Default to US/Canada
			}
			return countryCode + " " +
				digits[len(digits)-10:len(digits)-7] + " " +
				digits[len(digits)-7:len(digits)-4] + " " +
				digits[len(digits)-4:]
		}
		return string(p)
	default:
		return string(p)
	}
}

// CountryCode returns the country code of the phone number
func (p Phone) CountryCode() string {
	// This is a simplified implementation
	// A real implementation would use a library like libphonenumber
	digits := regexp.MustCompile(`\D`).ReplaceAllString(string(p), "")
	if len(digits) > 10 {
		return digits[:len(digits)-10]
	}
	return "1" // Default to US/Canada
}

// Normalized returns the phone number in E.164 format
func (p Phone) Normalized() string {
	return p.Format("e164")
}

// IsValidForCountry checks if the phone number is valid for the specified country
func (p Phone) IsValidForCountry(countryCode string) bool {
	// This is a simplified implementation
	// A real implementation would use a library like libphonenumber
	return p.CountryCode() == countryCode
}
