// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Name represents a person's name value object
// This is a wrapper around identification.Name for backward compatibility.
// New code should use identification.Name directly.
type Name = identification.Name

// NewName creates a new Name with validation
// This function is provided for backward compatibility.
// New code should use identification.NewName directly.
func NewName(name string) (Name, error) {
	return identification.NewName(name)
}
