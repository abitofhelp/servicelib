// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Duration represents a time duration value object
type Duration struct {
	duration time.Duration
}

// Regular expression for parsing duration strings in format "1h30m45s"
var durationRegex = regexp.MustCompile(`^(\d+h)?(\d+m)?(\d+s)?$`)

// NewDuration creates a new Duration with validation
func NewDuration(duration time.Duration) (Duration, error) {
	// Validate duration is not negative
	if duration < 0 {
		return Duration{}, errors.New("duration cannot be negative")
	}

	return Duration{duration: duration}, nil
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
func (d Duration) String() string {
	return d.duration.String()
}

// Equals checks if two Durations are equal
func (d Duration) Equals(other Duration) bool {
	return d.duration == other.duration
}

// IsEmpty checks if the Duration is zero
func (d Duration) IsEmpty() bool {
	return d.duration == 0
}

// Value returns the underlying time.Duration value
func (d Duration) Value() time.Duration {
	return d.duration
}

// Hours returns the duration as a floating point number of hours
func (d Duration) Hours() float64 {
	return d.duration.Hours()
}

// Minutes returns the duration as a floating point number of minutes
func (d Duration) Minutes() float64 {
	return d.duration.Minutes()
}

// Seconds returns the duration as a floating point number of seconds
func (d Duration) Seconds() float64 {
	return d.duration.Seconds()
}

// Milliseconds returns the duration as an integer number of milliseconds
func (d Duration) Milliseconds() int64 {
	return d.duration.Milliseconds()
}

// Add adds another Duration and returns a new Duration
func (d Duration) Add(other Duration) (Duration, error) {
	return NewDuration(d.duration + other.duration)
}

// Subtract subtracts another Duration and returns a new Duration
// If the result would be negative, it returns zero duration
func (d Duration) Subtract(other Duration) (Duration, error) {
	result := d.duration - other.duration
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
func (d Duration) Format(format string) string {
	hours := int(d.duration.Hours())
	minutes := int(d.duration.Minutes()) % 60
	seconds := int(d.duration.Seconds()) % 60

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
		return d.String()
	}
}
