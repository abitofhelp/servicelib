// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// DateOfBirth represents a date of birth value object
// This is a wrapper around identification.DateOfBirth for backward compatibility.
// New code should use identification.DateOfBirth directly.
type DateOfBirth = identification.DateOfBirth

// NewDateOfBirth creates a new DateOfBirth with validation
// This function is provided for backward compatibility.
// New code should use identification.NewDateOfBirth directly.
func NewDateOfBirth(year, month, day int) (DateOfBirth, error) {
	return identification.NewDateOfBirth(year, month, day)
}

// ParseDateOfBirth creates a new DateOfBirth from a string in format "YYYY-MM-DD"
// This function is provided for backward compatibility.
// New code should use identification.ParseDateOfBirth directly.
func ParseDateOfBirth(dateStr string) (DateOfBirth, error) {
	return identification.ParseDateOfBirth(dateStr)
}
