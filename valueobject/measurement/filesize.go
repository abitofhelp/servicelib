// Copyright (c) 2025 A Bit of Help, Inc.

// Package measurement provides value objects related to measurement information.
package measurement

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// FileSizeUnit represents a file size unit (using base 10 decimal values per ISO standards)
type FileSizeUnit string

const (
	// FileSizeBytes unit (B)
	FileSizeBytes FileSizeUnit = "B"
	// Kilobytes unit (KB) - 10^3 bytes
	Kilobytes FileSizeUnit = "KB"
	// Megabytes unit (MB) - 10^6 bytes
	Megabytes FileSizeUnit = "MB"
	// Gigabytes unit (GB) - 10^9 bytes
	Gigabytes FileSizeUnit = "GB"
	// Terabytes unit (TB) - 10^12 bytes
	Terabytes FileSizeUnit = "TB"
	// Petabytes unit (PB) - 10^15 bytes
	Petabytes FileSizeUnit = "PB"
)

// FileSize represents a file size value object that uses base 10 decimal values
// in accordance with the International System of Units (SI Units).
// This is the standard for file sizes in most contexts, where 1 KB = 1000 bytes,
// 1 MB = 1000 KB, etc., as opposed to the binary-based units used for memory sizes.
type FileSize struct {
	base.BaseStructValueObject
	bytes uint64
}

// Regular expression for parsing file size strings
var fileSizeRegex = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([KMGTP]?B)$`)

// NewFileSize creates a new FileSize with validation
// Uses base 10 decimal values in accordance with the International System of Units (SI Units),
// where 1 KB = 1000 bytes, 1 MB = 1000 KB, etc.
func NewFileSize(value float64, unit FileSizeUnit) (FileSize, error) {
	// Validate value
	if value < 0 {
		return FileSize{}, errors.New("file size cannot be negative")
	}

	// Convert to bytes
	var bytes uint64
	switch unit {
	case FileSizeBytes:
		bytes = uint64(value)
	case Kilobytes:
		bytes = uint64(value * 1000)
	case Megabytes:
		bytes = uint64(value * 1000 * 1000)
	case Gigabytes:
		bytes = uint64(value * 1000 * 1000 * 1000)
	case Terabytes:
		bytes = uint64(value * 1000 * 1000 * 1000 * 1000)
	case Petabytes:
		bytes = uint64(value * 1000 * 1000 * 1000 * 1000 * 1000)
	default:
		return FileSize{}, errors.New("invalid file size unit")
	}

	return FileSize{bytes: bytes}, nil
}

// NewFileSizeFromBytes creates a new FileSize directly from bytes
// The resulting FileSize will use base 10 decimal values in accordance with the
// International System of Units (SI Units) for all conversions and formatting.
func NewFileSizeFromBytes(bytes uint64) FileSize {
	return FileSize{bytes: bytes}
}

// ParseFileSize creates a new FileSize from a string
// The parsed value will use base 10 decimal values in accordance with the
// International System of Units (SI Units), where 1 KB = 1000 bytes, 1 MB = 1000 KB, etc.
func ParseFileSize(s string) (FileSize, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return FileSize{}, errors.New("file size string cannot be empty")
	}

	// Match against regex
	matches := fileSizeRegex.FindStringSubmatch(trimmed)
	if matches == nil {
		return FileSize{}, errors.New("invalid file size format, expected '1000B', '1.5KB', '2.5MB', etc.")
	}

	// Parse value
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return FileSize{}, errors.New("invalid file size value")
	}

	// Parse unit
	unit := FileSizeUnit(matches[2])

	return NewFileSize(value, unit)
}

// String returns the string representation of the FileSize
// Uses the most appropriate unit for readability
func (v FileSize) String() string {
	return v.Format("auto")
}

// Equals checks if two FileSizes are equal
func (v FileSize) Equals(other FileSize) bool {
	return v.bytes == other.bytes
}

// IsEmpty checks if the FileSize is empty (zero value)
func (v FileSize) IsEmpty() bool {
	return v.bytes == 0
}

// Validate checks if the FileSize is valid
func (v FileSize) Validate() error {
	return nil
}

// ToMap converts the FileSize to a map[string]interface{}
func (v FileSize) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"bytes": v.bytes,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (v FileSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

// Bytes returns the size in bytes
func (v FileSize) Bytes() uint64 {
	return v.bytes
}

// Kilobytes returns the size in kilobytes (base 10, ISO standard)
func (v FileSize) Kilobytes() float64 {
	return float64(v.bytes) / 1000
}

// Megabytes returns the size in megabytes (base 10, ISO standard)
func (v FileSize) Megabytes() float64 {
	return float64(v.bytes) / (1000 * 1000)
}

// Gigabytes returns the size in gigabytes (base 10, ISO standard)
func (v FileSize) Gigabytes() float64 {
	return float64(v.bytes) / (1000 * 1000 * 1000)
}

// Terabytes returns the size in terabytes (base 10, ISO standard)
func (v FileSize) Terabytes() float64 {
	return float64(v.bytes) / (1000 * 1000 * 1000 * 1000)
}

// Petabytes returns the size in petabytes (base 10, ISO standard)
func (v FileSize) Petabytes() float64 {
	return float64(v.bytes) / (1000 * 1000 * 1000 * 1000 * 1000)
}

// Add adds another FileSize and returns a new FileSize
func (v FileSize) Add(other FileSize) FileSize {
	return FileSize{bytes: v.bytes + other.bytes}
}

// Subtract subtracts another FileSize and returns a new FileSize
// If the result would be negative, it returns zero
func (v FileSize) Subtract(other FileSize) FileSize {
	if other.bytes > v.bytes {
		return FileSize{bytes: 0}
	}
	return FileSize{bytes: v.bytes - other.bytes}
}

// IsLargerThan checks if this FileSize is larger than another FileSize
func (v FileSize) IsLargerThan(other FileSize) bool {
	return v.bytes > other.bytes
}

// IsSmallerThan checks if this FileSize is smaller than another FileSize
func (v FileSize) IsSmallerThan(other FileSize) bool {
	return v.bytes < other.bytes
}

// Format returns the file size in the specified format
// Format options:
// - "B": Always in bytes
// - "KB": Always in kilobytes (base 10, ISO standard)
// - "MB": Always in megabytes (base 10, ISO standard)
// - "GB": Always in gigabytes (base 10, ISO standard)
// - "TB": Always in terabytes (base 10, ISO standard)
// - "PB": Always in petabytes (base 10, ISO standard)
// - "auto": Uses the most appropriate unit for readability (base 10, ISO standard)
// - "binary": Uses binary prefixes (KiB, MiB, GiB, etc.) - for compatibility only
func (v FileSize) Format(format string) string {
	switch format {
	case "B":
		return fmt.Sprintf("%d B", v.bytes)
	case "KB":
		return fmt.Sprintf("%.2f KB", v.Kilobytes())
	case "MB":
		return fmt.Sprintf("%.2f MB", v.Megabytes())
	case "GB":
		return fmt.Sprintf("%.2f GB", v.Gigabytes())
	case "TB":
		return fmt.Sprintf("%.2f TB", v.Terabytes())
	case "PB":
		return fmt.Sprintf("%.2f PB", v.Petabytes())
	case "binary":
		return v.formatBinary()
	case "auto":
		return v.formatDecimal() // Use decimal (base 10) as default per ISO standards
	default:
		return v.formatDecimal() // Use decimal (base 10) as default per ISO standards
	}
}

// formatAuto formats the file size using the most appropriate unit for readability
// This method is deprecated and maintained for backward compatibility only.
// It uses base 2 (1024) thresholds with base 10 (1000) conversion, which is inconsistent.
// Use formatDecimal instead for ISO-compliant file sizes.
func (v FileSize) formatAuto() string {
	return v.formatDecimal() // Redirect to formatDecimal for ISO compliance
}

// formatBinary formats the file size using binary prefixes (KiB, MiB, GiB, etc.)
func (v FileSize) formatBinary() string {
	bytes := float64(v.bytes)

	if bytes < 1024 {
		return fmt.Sprintf("%d B", v.bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KiB", v.Kilobytes())
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MiB", v.Megabytes())
	}
	if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GiB", v.Gigabytes())
	}
	if bytes < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f TiB", v.Terabytes())
	}
	return fmt.Sprintf("%.2f PiB", v.Petabytes())
}

// formatDecimal formats the file size using decimal prefixes (KB, MB, GB, etc.)
// in accordance with the International System of Units (SI Units).
// This uses base 10 values where 1 KB = 1000 bytes, 1 MB = 1000 KB, etc.,
// which is the standard for file sizes in most operating systems and file managers.
func (v FileSize) formatDecimal() string {
	bytes := float64(v.bytes)

	if bytes < 1000 {
		return fmt.Sprintf("%d B", v.bytes)
	}
	if bytes < 1000*1000 {
		kb := float64(v.bytes) / 1000
		return fmt.Sprintf("%.2f KB", kb)
	}
	if bytes < 1000*1000*1000 {
		mb := float64(v.bytes) / (1000 * 1000)
		return fmt.Sprintf("%.2f MB", mb)
	}
	if bytes < 1000*1000*1000*1000 {
		gb := float64(v.bytes) / (1000 * 1000 * 1000)
		return fmt.Sprintf("%.2f GB", gb)
	}
	if bytes < 1000*1000*1000*1000*1000 {
		tb := float64(v.bytes) / (1000 * 1000 * 1000 * 1000)
		return fmt.Sprintf("%.2f TB", tb)
	}
	pb := float64(v.bytes) / (1000 * 1000 * 1000 * 1000 * 1000)
	return fmt.Sprintf("%.2f PB", pb)
}
