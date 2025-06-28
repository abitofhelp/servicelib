// Copyright (c) 2025 A Bit of Help, Inc.

// Package base provides common interfaces and utilities for value objects.
package base

import (
	"encoding/json"
)

// StructValueObject is a base interface for struct-based value objects.
// It extends the ValueObject interface with methods specific to struct-based value objects.
type StructValueObject interface {
	ValueObject
	// ToMap converts the value object to a map[string]interface{}.
	ToMap() map[string]interface{}
	// MarshalJSON implements the json.Marshaler interface.
	MarshalJSON() ([]byte, error)
}

// BaseStructValueObject is a base implementation of StructValueObject.
// It provides common functionality for struct-based value objects.
// Embed this struct in your value object to inherit its functionality.
type BaseStructValueObject struct{}

// IsEmpty checks if the value object is empty (zero value).
// This is a default implementation that always returns false.
// Override this method in your value object to provide a proper implementation.
func (vo BaseStructValueObject) IsEmpty() bool {
	return false
}

// String returns the string representation of the value object.
// This is a default implementation that returns a generic string.
// Override this method in your value object to provide a proper implementation.
func (vo BaseStructValueObject) String() string {
	return "BaseStructValueObject"
}

// ToMap converts the value object to a map[string]interface{}.
// This is a default implementation that returns an empty map.
// Override this method in your value object to provide a proper implementation.
func (vo BaseStructValueObject) ToMap() map[string]interface{} {
	// Return an empty map to avoid infinite recursion
	return map[string]interface{}{}
}

// MarshalJSON implements the json.Marshaler interface.
// This is a default implementation that uses the ToMap method.
// Override this method in your value object to provide a more efficient implementation.
func (vo BaseStructValueObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.ToMap())
}

// WithValidation is a helper function that validates a value object after creation.
// It's useful for constructor functions that need to validate the created object.
// Example usage:
//
//	func NewMyValueObject(value string) (MyValueObject, error) {
//	    vo := MyValueObject{value: value}
//	    return WithValidation(vo, vo.Validate())
//	}
func WithValidation[T any](vo T, err error) (T, error) {
	if err != nil {
		var zero T
		return zero, err
	}
	return vo, nil
}
