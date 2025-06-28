// Copyright (c) 2025 A Bit of Help, Inc.

// Package base provides the core interfaces and implementations for value objects.
//
// This package serves as the foundation for the valueobject framework, defining the
// essential interfaces that all value objects should implement and providing base
// implementations for common value object types.
//
// Key interfaces defined in this package:
//   - ValueObject: The base interface for all value objects, defining methods for
//     string representation and emptiness checking.
//   - Equatable: For value objects that can be compared for equality.
//   - Comparable: For value objects that can be compared for ordering.
//   - Validatable: For value objects that can be validated.
//   - JSONMarshallable: For value objects that can be marshalled to JSON.
//   - JSONUnmarshallable: For value objects that can be unmarshalled from JSON.
//
// Base implementations provided:
//   - StringValueObject: A base implementation for string-based value objects.
//   - StructValueObject: A base implementation for struct-based value objects.
//
// The package also provides utility functions for common value object operations
// and validation helpers to ensure value objects maintain their invariants.
//
// Example usage:
//
//	// Define a custom string-based value object
//	type Email struct {
//	    base.StringValueObject
//	}
//
//	// Create a constructor with validation
//	func NewEmail(value string) (*Email, error) {
//	    // Create the value object
//	    email := &Email{}
//
//	    // Initialize with validation
//	    if err := email.Init(value); err != nil {
//	        return nil, err
//	    }
//
//	    // Additional validation
//	    if !strings.Contains(value, "@") {
//	        return nil, errors.New("invalid email format")
//	    }
//
//	    return email, nil
//	}
//
// The base package is designed to be extended by more specific value object types
// in other subpackages, providing a consistent foundation for all value objects
// in the application.
package base
