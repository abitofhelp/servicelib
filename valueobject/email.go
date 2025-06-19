// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

// Email represents an email address value object
type Email string

// Regular expression for validating email format
// Only allows alphanumeric characters, dots, underscores, and hyphens in the local part
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error) {
	// Trim whitespace
	trimmedEmail := strings.TrimSpace(email)

	// Empty email is allowed (optional field)
	if trimmedEmail == "" {
		return "", nil
	}

	// Validate email format
	_, err := mail.ParseAddress(trimmedEmail)
	if err != nil {
		return "", errors.New("invalid email format")
	}

	// Additional validation with regex to restrict special characters
	// Skip regex validation for emails with names in the format "Name <email@example.com>"
	if !strings.Contains(trimmedEmail, "<") && !strings.Contains(trimmedEmail, ">") {
		if !emailRegex.MatchString(trimmedEmail) {
			return "", errors.New("invalid email format: contains invalid characters")
		}
	}

	return Email(trimmedEmail), nil
}

// String returns the string representation of the Email
func (e Email) String() string {
	return string(e)
}

// Equals checks if two Emails are equal
func (e Email) Equals(other Email) bool {
	return strings.EqualFold(string(e), string(other))
}

// IsEmpty checks if the Email is empty
func (e Email) IsEmpty() bool {
	return e == ""
}
