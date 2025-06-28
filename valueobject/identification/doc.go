// Copyright (c) 2025 A Bit of Help, Inc.

// Package identification provides value objects related to identification.
//
// This package contains value objects that represent different types of identifiers
// and identifying information such as IDs, names, usernames, and other personal
// identifiers. These value objects are immutable and follow the Value Object pattern
// from Domain-Driven Design.
//
// Key value objects in this package:
//   - ID: Represents a unique identifier (UUID)
//   - Name: Represents a person's name
//   - Username: Represents a username with validation rules
//   - DateOfBirth: Represents a person's date of birth
//   - DateOfDeath: Represents a person's date of death
//   - Gender: Represents a person's gender
//   - Password: Represents a password with security requirements
//
// Each value object provides methods for:
//   - Creating and validating instances
//   - String representation
//   - Equality comparison
//   - Emptiness checking
//
// Some value objects provide additional functionality specific to their domain:
//   - ID: Generation of new random UUIDs
//   - Username: Case conversion, length calculation, substring checking
//   - Password: Strength checking, hashing
//
// Example usage:
//
//	// Create a new ID
//	id, err := identification.NewID("550e8400-e29b-41d4-a716-446655440000")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Generate a random ID
//	randomID := identification.GenerateID()
//
//	// Create a new name
//	name, err := identification.NewName("John Doe")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Create a new username
//	username, err := identification.NewUsername("john_doe")
//	if err != nil {
//	    // Handle validation error
//	}
//
// All value objects in this package are designed to be immutable, so they cannot
// be changed after creation. To modify a value object, create a new instance with
// the desired values.
package identification