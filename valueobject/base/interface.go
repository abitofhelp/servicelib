// Copyright (c) 2025 A Bit of Help, Inc.

// Package base provides common interfaces and utilities for value objects.
package base

// ValueObject is the base interface for all value objects.
// It defines the common methods that all value objects should implement.
type ValueObject interface {
	// String returns the string representation of the value object.
	String() string

	// IsEmpty checks if the value object is empty (zero value).
	IsEmpty() bool
}

// Equatable is an interface for value objects that can be compared for equality.
type Equatable[T any] interface {
	// Equals checks if two value objects are equal.
	Equals(other T) bool
}

// Comparable is an interface for value objects that can be compared for ordering.
type Comparable[T any] interface {
	// CompareTo compares this value object with another and returns:
	// -1 if this < other
	//  0 if this == other
	//  1 if this > other
	CompareTo(other T) int
}

// Validatable is an interface for value objects that can be validated.
type Validatable interface {
	// Validate checks if the value object is valid.
	// Returns nil if valid, otherwise returns an error.
	Validate() error
}

// JSONMarshallable is an interface for value objects that can be marshalled to JSON.
type JSONMarshallable interface {
	// MarshalJSON implements the json.Marshaler interface.
	MarshalJSON() ([]byte, error)
}

// JSONUnmarshallable is an interface for value objects that can be unmarshalled from JSON.
type JSONUnmarshallable interface {
	// UnmarshalJSON implements the json.Unmarshaler interface.
	UnmarshalJSON(data []byte) error
}
