// Copyright (c) 2025 A Bit of Help, Inc.

// Package base provides common interfaces and utilities for value objects.
package base

import (
	"math"
	"strings"
)

// CompareStrings compares two strings case-insensitively.
// Returns:
// -1 if a < b
//
//	0 if a == b
//	1 if a > b
func CompareStrings(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// StringsEqualFold checks if two strings are equal, ignoring case.
func StringsEqualFold(a, b string) bool {
	return strings.EqualFold(a, b)
}

// CompareFloats compares two float64 values with a small epsilon to handle floating point precision.
// Returns:
// -1 if a < b
//
//	0 if a == b
//	1 if a > b
func CompareFloats(a, b float64) int {
	epsilon := 0.0000001
	diff := a - b
	switch {
	case math.Abs(diff) < epsilon:
		return 0
	case diff < 0:
		return -1
	default:
		return 1
	}
}

// FloatsEqual checks if two float64 values are equal, with a small epsilon to handle floating point precision.
func FloatsEqual(a, b float64) bool {
	return CompareFloats(a, b) == 0
}

// TrimAndNormalize trims whitespace and normalizes a string.
// This is useful for value objects that need to normalize their string representation.
func TrimAndNormalize(s string) string {
	return strings.TrimSpace(s)
}

// IsEmptyString checks if a string is empty after trimming whitespace.
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsZeroFloat checks if a float64 is zero, with a small epsilon to handle floating point precision.
func IsZeroFloat(f float64) bool {
	return math.Abs(f) < 0.0000001
}

// RoundFloat rounds a float64 to the specified number of decimal places.
func RoundFloat(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(f*shift) / shift
}
