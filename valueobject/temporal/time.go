// Copyright (c) 2025 A Bit of Help, Inc.

// Package temporal provides value objects related to time and duration.
// It includes implementations for time-based value objects that follow
// the domain-driven design principles.
package temporal

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Time represents a time value object based on RFC 3339 format.
// It encapsulates a time.Time value and provides methods for
// comparison, arithmetic, and conversion operations.
type Time struct {
	base.BaseStructValueObject
	value time.Time
}

// NewTime creates a new Time with validation.
// It returns an error if the validation fails.
func NewTime(value time.Time) (Time, error) {
	vo := Time{
		BaseStructValueObject: base.BaseStructValueObject{},
		value:                 value,
	}

	// No specific validation needed for time.Time as it's always valid
	// We could add domain-specific validation here if needed in the future

	return vo, nil
}

// NewTimeFromString creates a new Time from an RFC 3339 formatted string.
// It returns an error if the string is empty or not in valid RFC 3339 format.
func NewTimeFromString(value string) (Time, error) {
	if value == "" {
		return Time{}, errors.New("time string cannot be empty")
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return Time{}, errors.New("invalid RFC 3339 time format")
	}

	return NewTime(t)
}

// String returns the RFC 3339 string representation of the time.
func (t Time) String() string {
	return t.value.Format(time.RFC3339)
}

// IsEmpty checks if the time is empty (zero value).
func (t Time) IsEmpty() bool {
	return t.value.IsZero()
}

// Value returns the underlying time.Time value.
func (t Time) Value() time.Time {
	return t.value
}

// Equals checks if two times are equal.
// This method follows the naming convention used in other value objects.
func (t Time) Equals(other Time) bool {
	return t.value.Equal(other.value)
}

// Equal is maintained for backward compatibility.
// New code should use Equals() instead.
func (t Time) Equal(other Time) bool {
	return t.Equals(other)
}

// Before checks if this time is before another time.
func (t Time) Before(other Time) bool {
	return t.value.Before(other.value)
}

// After checks if this time is after another time.
func (t Time) After(other Time) bool {
	return t.value.After(other.value)
}

// Add returns a new Time with the given duration added.
// It properly handles any errors from the NewTime constructor.
func (t Time) Add(d time.Duration) Time {
	newTime, err := NewTime(t.value.Add(d))
	if err != nil {
		// This should never happen with a valid time.Time
		// But we handle it anyway for robustness
		return t // Return the original time if there's an error
	}
	return newTime
}

// Sub returns the duration between two times.
func (t Time) Sub(other Time) time.Duration {
	return t.value.Sub(other.value)
}

// ToMap converts the time to a map representation.
func (t Time) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"value": t.String(),
	}
}

// MarshalJSON implements json.Marshaler interface.
// It marshals the time as an RFC 3339 formatted string.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements json.Unmarshaler interface.
// It unmarshals an RFC 3339 formatted string into a Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	newTime, err := NewTimeFromString(timeStr)
	if err != nil {
		return err
	}

	*t = newTime
	return nil
}

// Validate checks if the Time is valid.
// Currently, all time.Time values are considered valid.
// This method is included for consistency with other value objects.
func (t Time) Validate() error {
	// No specific validation needed for time.Time
	// We could add domain-specific validation here if needed
	return nil
}
