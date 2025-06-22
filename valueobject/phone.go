// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/contact"
)

// Phone represents a phone number value object
// This is a wrapper around contact.Phone for backward compatibility.
// New code should use contact.Phone directly.
type Phone = contact.Phone

// NewPhone creates a new Phone with validation
// This function is provided for backward compatibility.
// New code should use contact.NewPhone directly.
func NewPhone(phone string) (Phone, error) {
	return contact.NewPhone(phone)
}
