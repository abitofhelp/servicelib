// Copyright (c) 2025 A Bit of Help, Inc.

// Package base provides common interfaces and utilities for value objects.
package base

import (
	"strings"
)

// StringValueObject is a base type for string-based value objects.
// It implements the ValueObject and Equatable interfaces.
type StringValueObject struct {
	value string
}

// NewStringValueObject creates a new StringValueObject with the given value.
func NewStringValueObject(value string) StringValueObject {
	return StringValueObject{
		value: strings.TrimSpace(value),
	}
}

// String returns the string representation of the value object.
func (vo StringValueObject) String() string {
	return vo.value
}

// IsEmpty checks if the value object is empty (zero value).
func (vo StringValueObject) IsEmpty() bool {
	return vo.value == ""
}

// Equals checks if two StringValueObjects are equal.
// This is a case-sensitive comparison. For case-insensitive comparison,
// use EqualsIgnoreCase.
func (vo StringValueObject) Equals(other StringValueObject) bool {
	return vo.value == other.value
}

// EqualsIgnoreCase checks if two StringValueObjects are equal, ignoring case.
func (vo StringValueObject) EqualsIgnoreCase(other StringValueObject) bool {
	return StringsEqualFold(vo.value, other.value)
}

// Value returns the underlying string value.
func (vo StringValueObject) Value() string {
	return vo.value
}

// WithValue returns a new StringValueObject with the given value.
func (vo StringValueObject) WithValue(value string) StringValueObject {
	return NewStringValueObject(value)
}
