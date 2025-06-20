// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"strings"
)

// Address represents a postal address value object
type Address string

// NewAddress creates a new Address with validation
func NewAddress(address string) (Address, error) {
	// Trim whitespace
	trimmedAddress := strings.TrimSpace(address)

	// Empty address is allowed (optional field)
	if trimmedAddress == "" {
		return "", nil
	}

	// Basic validation - ensure minimum length
	if len(trimmedAddress) < 5 {
		return "", errors.New("address is too short")
	}

	// Maximum length validation to prevent abuse
	if len(trimmedAddress) > 200 {
		return "", errors.New("address is too long")
	}

	return Address(trimmedAddress), nil
}

// String returns the string representation of the Address
func (a Address) String() string {
	return string(a)
}

// Equals checks if two Addresses are equal
func (a Address) Equals(other Address) bool {
	return strings.EqualFold(string(a), string(other))
}

// IsEmpty checks if the Address is empty
func (a Address) IsEmpty() bool {
	return a == ""
}
