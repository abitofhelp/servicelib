// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/contact"
)

// Email represents an email address value object
// This is a wrapper around contact.Email for backward compatibility.
// New code should use contact.Email directly.
type Email = contact.Email

// NewEmail creates a new Email with validation
// This function is provided for backward compatibility.
// New code should use contact.NewEmail directly.
func NewEmail(email string) (Email, error) {
	return contact.NewEmail(email)
}
