// Copyright (c) 2025 A Bit of Help, Inc.

// Package contact provides value objects related to contact information.
//
// This package contains value objects that represent different types of contact
// information such as email addresses, phone numbers, and postal addresses. These
// value objects are immutable and follow the Value Object pattern from Domain-Driven
// Design.
//
// Key value objects in this package:
//   - Email: Represents an email address
//   - Phone: Represents a phone number with formatting capabilities
//   - Address: Represents a postal address
//
// Each value object provides methods for:
//   - Creating and validating instances
//   - String representation
//   - Equality comparison
//   - Emptiness checking
//
// The Phone value object additionally provides methods for:
//   - Formatting in different standards (E.164, national, international)
//   - Extracting country codes
//   - Normalizing to standard formats
//   - Validating for specific countries
//
// Example usage:
//
//	// Create a new email address
//	email, err := contact.NewEmail("user@example.com")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Create a new phone number
//	phone, err := contact.NewPhone("+1 (234) 567-8901")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Format phone number in different ways
//	e164Format := phone.Format("e164")        // "+12345678901"
//	nationalFormat := phone.Format("national") // "(234) 567-8901"
//
//	// Create a new address
//	address, err := contact.NewAddress("123 Main St, Anytown, USA")
//	if err != nil {
//	    // Handle validation error
//	}
//
// All value objects in this package are designed to be immutable, so they cannot
// be changed after creation. To modify a value object, create a new instance with
// the desired values.
package contact
