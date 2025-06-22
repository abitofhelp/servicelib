// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Password represents a password value object
// This is a wrapper around identification.Password for backward compatibility.
// New code should use identification.Password directly.
type Password = identification.Password

// NewPassword creates a new Password with validation and hashing
// This function is provided for backward compatibility.
// New code should use identification.NewPassword directly.
func NewPassword(plainPassword string) (Password, error) {
	return identification.NewPassword(plainPassword)
}
