// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Percentage represents a percentage value object
type Percentage float64

// NewPercentage creates a new Percentage with validation
func NewPercentage(value float64) (Percentage, error) {
	// Validate percentage range
	if value < 0 {
		return 0, errors.New("percentage cannot be negative")
	}
	
	if value > 100 {
		return 0, errors.New("percentage cannot exceed 100")
	}
	
	// Round to 2 decimal places for precision
	roundedValue := float64(int64(value*100+0.5)) / 100
	
	return Percentage(roundedValue), nil
}

// ParsePercentage creates a new Percentage from a string
func ParsePercentage(s string) (Percentage, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)
	
	// Empty string is not allowed
	if trimmed == "" {
		return 0, errors.New("percentage string cannot be empty")
	}
	
	// Remove percentage sign if present
	if strings.HasSuffix(trimmed, "%") {
		trimmed = strings.TrimSuffix(trimmed, "%")
	}
	
	// Parse the value
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, errors.New("invalid percentage format")
	}
	
	return NewPercentage(value)
}

// String returns the string representation of the Percentage
func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", float64(p))
}

// Equals checks if two Percentages are equal
func (p Percentage) Equals(other Percentage) bool {
	// Compare with small epsilon to handle floating point precision
	epsilon := 0.0001
	diff := float64(p) - float64(other)
	return diff > -epsilon && diff < epsilon
}

// IsEmpty checks if the Percentage is zero
func (p Percentage) IsEmpty() bool {
	return p == 0
}

// Value returns the raw float64 value
func (p Percentage) Value() float64 {
	return float64(p)
}

// AsDecimal returns the percentage as a decimal (e.g., 75% -> 0.75)
func (p Percentage) AsDecimal() float64 {
	return float64(p) / 100
}

// Of calculates the percentage of a given value
// For example: 25% of 200 = 50
func (p Percentage) Of(value float64) float64 {
	return p.AsDecimal() * value
}

// Inverse returns the inverse percentage (100% - p)
func (p Percentage) Inverse() (Percentage, error) {
	return NewPercentage(100 - float64(p))
}

// Add adds another percentage and returns a new Percentage
// The result is capped at 100%
func (p Percentage) Add(other Percentage) (Percentage, error) {
	sum := float64(p) + float64(other)
	return NewPercentage(sum)
}

// Subtract subtracts another percentage and returns a new Percentage
// The result is floored at 0%
func (p Percentage) Subtract(other Percentage) (Percentage, error) {
	diff := float64(p) - float64(other)
	if diff < 0 {
		diff = 0
	}
	return NewPercentage(diff)
}