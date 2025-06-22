// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Gender represents a person's gender value object
// This is a wrapper around identification.Gender for backward compatibility.
// New code should use identification.Gender directly.
type Gender = identification.Gender

// Gender constants
const (
	GenderMale   = identification.GenderMale
	GenderFemale = identification.GenderFemale
	GenderOther  = identification.GenderOther
)

// NewGender creates a new Gender with validation
// This function is provided for backward compatibility.
// New code should use identification.NewGender directly.
func NewGender(gender string) (Gender, error) {
	return identification.NewGender(gender)
}
