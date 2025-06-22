// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// ID represents a unique identifier value object
// This is a wrapper around identification.ID for backward compatibility.
// New code should use identification.ID directly.
type ID = identification.ID

// NewID creates a new ID with validation
// This function is provided for backward compatibility.
// New code should use identification.NewID directly.
func NewID(id string) (ID, error) {
	return identification.NewID(id)
}

// GenerateID creates a new random ID
// This function is provided for backward compatibility.
// New code should use identification.GenerateID directly.
func GenerateID() ID {
	return identification.GenerateID()
}
