// Copyright (c) 2025 A Bit of Help, Inc.

// Package measurement provides value objects related to measurement information.
package measurement

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Percentage represents a percentage value object
type Percentage struct {
	base.BaseStructValueObject

	value float64
}

// NewPercentage creates a new Percentage with validation
func NewPercentage(value float64) (Percentage, error) {
	vo := Percentage{

		value: value,
	}

	// Validate the value object
	if err := vo.Validate(); err != nil {
		return Percentage{}, err
	}

	return vo, nil
}

// String returns the string representation of the Percentage
func (v Percentage) String() string {
	return fmt.Sprintf("%.2f%%", v.value)
}

// Equals checks if two Percentages are equal
func (v Percentage) Equals(other Percentage) bool {

	if !base.FloatsEqual(v.value, other.value) {
		return false
	}

	return true
}

// IsEmpty checks if the Percentage is empty (zero value)
func (v Percentage) IsEmpty() bool {
	return v.value == float64(0)
}

// Validate checks if the Percentage is valid
func (v Percentage) Validate() error {

	// Validate percentage range
	if v.value < 0 {
		return errors.New("percentage cannot be negative")
	}

	if v.value > 100 {
		return errors.New("percentage cannot exceed 100")
	}

	return nil
}

// ToMap converts the Percentage to a map[string]interface{}
func (v Percentage) ToMap() map[string]interface{} {
	return map[string]interface{}{

		"value": v.value,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (v Percentage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

// Value returns the raw float64 value
func (v Percentage) Value() float64 {
	return v.value
}

// AsDecimal returns the percentage as a decimal (e.g., 75% -> 0.75)
func (v Percentage) AsDecimal() float64 {
	return v.value / 100
}

// Of calculates the percentage of a given value
// For example: 25% of 200 = 50
func (v Percentage) Of(value float64) float64 {
	return v.AsDecimal() * value
}

// Inverse returns the inverse percentage (100% - p)
func (v Percentage) Inverse() (Percentage, error) {
	return NewPercentage(100 - v.value)
}

// Add adds another percentage and returns a new Percentage
// The result is capped at 100%
func (v Percentage) Add(other Percentage) (Percentage, error) {
	sum := v.value + other.value
	return NewPercentage(sum)
}

// Subtract subtracts another percentage and returns a new Percentage
// The result is floored at 0%
func (v Percentage) Subtract(other Percentage) (Percentage, error) {
	diff := v.value - other.value
	if diff < 0 {
		diff = 0
	}
	return NewPercentage(diff)
}

// ParsePercentage creates a new Percentage from a string
func ParsePercentage(s string) (Percentage, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return Percentage{}, errors.New("percentage string cannot be empty")
	}

	// Remove percentage sign if present
	if strings.HasSuffix(trimmed, "%") {
		trimmed = strings.TrimSuffix(trimmed, "%")
	}

	// Parse the value
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return Percentage{}, errors.New("invalid percentage format")
	}

	return NewPercentage(value)
}
