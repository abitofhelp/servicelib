// Copyright (c) 2025 A Bit of Help, Inc.

// Package temporal provides value objects related to temporal information.
package temporal

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Duration represents a time duration value object
type Duration struct {
	base.BaseStructValueObject
	duration time.Duration
}

// Regular expression for parsing duration strings in format "1h30m45s"
var durationRegex = regexp.MustCompile(`^(\d+h)?(\d+m)?(\d+s)?$`)

// NewDuration creates a new Duration with validation
func NewDuration(duration time.Duration) (Duration, error) {
	vo := Duration{
		duration: duration,
	}

	// Validate the value object
	if err := vo.Validate(); err != nil {
		return Duration{}, err
	}

	return vo, nil
}

// ParseDuration creates a new Duration from a string
func ParseDuration(s string) (Duration, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return Duration{}, errors.New("duration string cannot be empty")
	}

	// Try to parse using Go's built-in parser
	duration, err := time.ParseDuration(trimmed)
	if err != nil {
		return Duration{}, errors.New("invalid duration format, expected format like '1h30m45s'")
	}

	return NewDuration(duration)
}

// String returns the string representation of the Duration
func (v Duration) String() string {
	return v.duration.String()
}

// Equals checks if two Durations are equal
func (v Duration) Equals(other Duration) bool {
	return v.duration == other.duration
}

// IsEmpty checks if the Duration is empty (zero value)
func (v Duration) IsEmpty() bool {
	return v.duration == 0
}

// Validate checks if the Duration is valid
func (v Duration) Validate() error {
	if v.duration < 0 {
		return errors.New("duration cannot be negative")
	}
	return nil
}

// ToMap converts the Duration to a map[string]interface{}
func (v Duration) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"duration": v.duration,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (v Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

// Value returns the underlying time.Duration value
func (v Duration) Value() time.Duration {
	return v.duration
}

// Hours returns the duration as a floating point number of hours
func (v Duration) Hours() float64 {
	return v.duration.Hours()
}

// Minutes returns the duration as a floating point number of minutes
func (v Duration) Minutes() float64 {
	return v.duration.Minutes()
}

// Seconds returns the duration as a floating point number of seconds
func (v Duration) Seconds() float64 {
	return v.duration.Seconds()
}

// Milliseconds returns the duration as an integer number of milliseconds
func (v Duration) Milliseconds() int64 {
	return v.duration.Milliseconds()
}

// Add adds another Duration and returns a new Duration
func (v Duration) Add(other Duration) (Duration, error) {
	return NewDuration(v.duration + other.duration)
}

// Subtract subtracts another Duration and returns a new Duration
// If the result would be negative, it returns zero duration
func (v Duration) Subtract(other Duration) (Duration, error) {
	result := v.duration - other.duration
	if result < 0 {
		result = 0
	}
	return NewDuration(result)
}

// Format returns a formatted string representation of the duration
// Format options:
// - "short": "1h 30m 45s"
// - "long": "1 hour 30 minutes 45 seconds"
// - "compact": "1:30:45"
func (v Duration) Format(format string) string {
	hours := int(v.duration.Hours())
	minutes := int(v.duration.Minutes()) % 60
	seconds := int(v.duration.Seconds()) % 60

	switch format {
	case "short":
		parts := []string{}
		if hours > 0 {
			parts = append(parts, fmt.Sprintf("%dh", hours))
		}
		if minutes > 0 || (hours > 0 && seconds > 0) {
			parts = append(parts, fmt.Sprintf("%dm", minutes))
		}
		if seconds > 0 || len(parts) == 0 {
			parts = append(parts, fmt.Sprintf("%ds", seconds))
		}
		return strings.Join(parts, " ")

	case "long":
		parts := []string{}
		if hours > 0 {
			if hours == 1 {
				parts = append(parts, "1 hour")
			} else {
				parts = append(parts, fmt.Sprintf("%d hours", hours))
			}
		}
		if minutes > 0 {
			if minutes == 1 {
				parts = append(parts, "1 minute")
			} else {
				parts = append(parts, fmt.Sprintf("%d minutes", minutes))
			}
		}
		if seconds > 0 || len(parts) == 0 {
			if seconds == 1 {
				parts = append(parts, "1 second")
			} else {
				parts = append(parts, fmt.Sprintf("%d seconds", seconds))
			}
		}
		return strings.Join(parts, " ")

	case "compact":
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)

	default:
		return v.String()
	}
}
