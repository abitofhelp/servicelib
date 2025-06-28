// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a comprehensive framework for implementing domain value objects.
//
// Value objects are immutable objects that represent concepts or entities in the domain model.
// They are defined by their attributes rather than their identity, meaning two value objects
// with the same attributes are considered equal regardless of their memory location.
//
// This package implements the Value Object pattern from Domain-Driven Design (DDD) and provides:
//   - Base interfaces and implementations for common value object types
//   - Validation mechanisms to ensure value objects maintain invariants
//   - Serialization support for persistence and data transfer
//   - Type-specific implementations for common domain concepts
//
// The package is organized into several subpackages:
//   - base: Core interfaces and base implementations for value objects
//   - appearance: Value objects related to visual appearance (color, style, etc.)
//   - contact: Value objects for contact information (email, phone, address, etc.)
//   - identification: Value objects for identifiers (UUID, custom IDs, etc.)
//   - location: Value objects for geographical locations (coordinates, addresses, etc.)
//   - measurement: Value objects for measurements (length, weight, temperature, etc.)
//   - network: Value objects for network concepts (IP address, URL, etc.)
//   - temporal: Value objects for time-related concepts (date range, time period, etc.)
//   - generator: Utilities for generating value objects
//
// Example usage:
//
//	// Create an email address value object
//	email, err := contact.NewEmail("user@example.com")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Create a UUID identifier
//	id, err := identification.NewUUID()
//	if err != nil {
//	    // Handle error
//	}
//
//	// Create a date range
//	dateRange, err := temporal.NewDateRange(startDate, endDate)
//	if err != nil {
//	    // Handle validation error
//	}
//
// Value objects are designed to be immutable, so they cannot be changed after creation.
// To modify a value object, create a new instance with the desired values.
package valueobject
