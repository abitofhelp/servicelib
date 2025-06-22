// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Username represents a username value object
// This is a wrapper around identification.Username for backward compatibility.
// New code should use identification.Username directly.
type Username = identification.Username

// NewUsername creates a new Username with validation
// This function is provided for backward compatibility.
// New code should use identification.NewUsername directly.
func NewUsername(username string) (Username, error) {
	return identification.NewUsername(username)
}
