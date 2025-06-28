// Copyright (c) 2025 A Bit of Help, Inc.

// Package measurement provides value objects related to measurement information.
package measurement

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abitofhelp/servicelib/valueobject/base"
	"regexp"
	"strconv"
	"strings"
)

// MemSizeUnit represents a memory size unit (using base 2 binary values)
type MemSizeUnit string

const (
	// Bytes unit (B)
	Bytes MemSizeUnit = "B"
	// Kibibytes unit (KiB) - 2^10 bytes
	Kibibytes MemSizeUnit = "KiB"
	// Mebibytes unit (MiB) - 2^20 bytes
	Mebibytes MemSizeUnit = "MiB"
	// Gibibytes unit (GiB) - 2^30 bytes
	Gibibytes MemSizeUnit = "GiB"
	// Tebibytes unit (TiB) - 2^40 bytes
	Tebibytes MemSizeUnit = "TiB"
	// Pebibytes unit (PiB) - 2^50 bytes
	Pebibytes MemSizeUnit = "PiB"
)

// MemSize represents a memory size value object using base 2 binary values
// This is specifically for computer memory measurements, not file sizes.
type MemSize struct {
	base.BaseStructValueObject
	bytes uint64
}

// Regular expression for parsing memory size strings
var memSizeRegex = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([KMGTP]i?B|B)$`)

// NewMemSize creates a new MemSize with validation
func NewMemSize(value float64, unit MemSizeUnit) (MemSize, error) {
	// Validate value
	if value < 0 {
		return MemSize{}, errors.New("memory size cannot be negative")
	}

	// Convert to bytes using base 2 (binary) values
	var bytes uint64
	switch unit {
	case Bytes:
		bytes = uint64(value)
	case Kibibytes:
		bytes = uint64(value * 1024)
	case Mebibytes:
		bytes = uint64(value * 1024 * 1024)
	case Gibibytes:
		bytes = uint64(value * 1024 * 1024 * 1024)
	case Tebibytes:
		bytes = uint64(value * 1024 * 1024 * 1024 * 1024)
	case Pebibytes:
		bytes = uint64(value * 1024 * 1024 * 1024 * 1024 * 1024)
	default:
		return MemSize{}, errors.New("invalid memory size unit")
	}

	return MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 bytes,
	}, nil
}

// ParseMemSize creates a new MemSize from a string
func ParseMemSize(s string) (MemSize, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return MemSize{
			BaseStructValueObject: base.BaseStructValueObject{},
		}, errors.New("memory size string cannot be empty")
	}

	// Match against regex
	matches := memSizeRegex.FindStringSubmatch(trimmed)
	if matches == nil {
		return MemSize{
			BaseStructValueObject: base.BaseStructValueObject{},
		}, errors.New("invalid memory size format, expected '1024B', '1.5KiB', '2.5MiB', etc.")
	}

	// Parse value
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return MemSize{
			BaseStructValueObject: base.BaseStructValueObject{},
		}, errors.New("invalid memory size value")
	}

	// Parse unit
	unitStr := matches[2]
	var unit MemSizeUnit

	// Convert legacy units to binary units if needed
	switch unitStr {
	case "B":
		unit = Bytes
	case "KB", "KiB":
		unit = Kibibytes
	case "MB", "MiB":
		unit = Mebibytes
	case "GB", "GiB":
		unit = Gibibytes
	case "TB", "TiB":
		unit = Tebibytes
	case "PB", "PiB":
		unit = Pebibytes
	default:
		return MemSize{
			BaseStructValueObject: base.BaseStructValueObject{},
		}, errors.New("invalid memory size unit")
	}

	return NewMemSize(value, unit)
}

// String returns the string representation of the MemSize
// Uses the most appropriate unit for readability with binary (base 2) units
func (ms MemSize) String() string {
	return ms.Format("auto")
}

// Equals checks if two MemSizes are equal
func (ms MemSize) Equals(other MemSize) bool {
	return ms.bytes == other.bytes
}

// IsEmpty checks if the MemSize is empty (zero value)
func (ms MemSize) IsEmpty() bool {
	return ms.bytes == 0
}

// Validate checks if the MemSize is valid
func (ms MemSize) Validate() error {
	return nil
}

// ToMap converts the MemSize to a map[string]interface{}
func (ms MemSize) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"bytes": ms.bytes,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (ms MemSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(ms.ToMap())
}

// Bytes returns the size in bytes
func (ms MemSize) Bytes() uint64 {
	return ms.bytes
}

// Kibibytes returns the size in kibibytes (KiB)
func (ms MemSize) Kibibytes() float64 {
	return float64(ms.bytes) / 1024
}

// Mebibytes returns the size in mebibytes (MiB)
func (ms MemSize) Mebibytes() float64 {
	return float64(ms.bytes) / (1024 * 1024)
}

// Gibibytes returns the size in gibibytes (GiB)
func (ms MemSize) Gibibytes() float64 {
	return float64(ms.bytes) / (1024 * 1024 * 1024)
}

// Tebibytes returns the size in tebibytes (TiB)
func (ms MemSize) Tebibytes() float64 {
	return float64(ms.bytes) / (1024 * 1024 * 1024 * 1024)
}

// Pebibytes returns the size in pebibytes (PiB)
func (ms MemSize) Pebibytes() float64 {
	return float64(ms.bytes) / (1024 * 1024 * 1024 * 1024 * 1024)
}

// Add adds another MemSize and returns a new MemSize
func (ms MemSize) Add(other MemSize) MemSize {
	return MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 ms.bytes + other.bytes,
	}
}

// Subtract subtracts another MemSize and returns a new MemSize
// If the result would be negative, it returns zero
func (ms MemSize) Subtract(other MemSize) MemSize {
	if other.bytes > ms.bytes {
		return MemSize{
			BaseStructValueObject: base.BaseStructValueObject{},
			bytes:                 0,
		}
	}
	return MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 ms.bytes - other.bytes,
	}
}

// IsLargerThan checks if this MemSize is larger than another MemSize
func (ms MemSize) IsLargerThan(other MemSize) bool {
	return ms.bytes > other.bytes
}

// IsSmallerThan checks if this MemSize is smaller than another MemSize
func (ms MemSize) IsSmallerThan(other MemSize) bool {
	return ms.bytes < other.bytes
}

// Format returns the memory size in the specified format
// Format options:
// - "B": Always in bytes
// - "KiB": Always in kibibytes
// - "MiB": Always in mebibytes
// - "GiB": Always in gibibytes
// - "TiB": Always in tebibytes
// - "PiB": Always in pebibytes
// - "auto": Uses the most appropriate unit for readability
func (ms MemSize) Format(format string) string {
	switch format {
	case "B":
		return fmt.Sprintf("%d B", ms.bytes)
	case "KiB":
		return fmt.Sprintf("%.2f KiB", ms.Kibibytes())
	case "MiB":
		return fmt.Sprintf("%.2f MiB", ms.Mebibytes())
	case "GiB":
		return fmt.Sprintf("%.2f GiB", ms.Gibibytes())
	case "TiB":
		return fmt.Sprintf("%.2f TiB", ms.Tebibytes())
	case "PiB":
		return fmt.Sprintf("%.2f PiB", ms.Pebibytes())
	case "auto":
		return ms.formatAuto()
	default:
		return ms.formatAuto()
	}
}

// formatAuto formats the memory size using the most appropriate unit for readability
// Always uses binary (base 2) units with the "i" prefix (KiB, MiB, etc.)
func (ms MemSize) formatAuto() string {
	bytes := float64(ms.bytes)

	if bytes < 1024 {
		return fmt.Sprintf("%d B", ms.bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KiB", ms.Kibibytes())
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MiB", ms.Mebibytes())
	}
	if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GiB", ms.Gibibytes())
	}
	if bytes < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f TiB", ms.Tebibytes())
	}
	return fmt.Sprintf("%.2f PiB", ms.Pebibytes())
}
