// Copyright (c) 2025 A Bit of Help, Inc.

// Package contact provides value objects related to contact information.
package contact

import (
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Email represents an email address value object
type Email string

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error) {
	// Validate the email using the common validation function
	if err := base.ValidateEmail(email); err != nil {
		return "", err
	}

	// Trim whitespace
	trimmedEmail := strings.TrimSpace(email)

	return Email(trimmedEmail), nil
}

// String returns the string representation of the Email
func (e Email) String() string {
	return string(e)
}

// Equals checks if two Emails are equal
func (e Email) Equals(other Email) bool {
	return base.StringsEqualFold(string(e), string(other))
}

// IsEmpty checks if the Email is empty
func (e Email) IsEmpty() bool {
	return e == ""
}

// Validate checks if the Email is valid
func (e Email) Validate() error {
	return base.ValidateEmail(string(e))
}