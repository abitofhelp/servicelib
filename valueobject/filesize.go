// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// FileSizeUnit represents a file size unit
type FileSizeUnit string

const (
	// Bytes unit (B)
	Bytes FileSizeUnit = "B"
	// Kilobytes unit (KB)
	Kilobytes FileSizeUnit = "KB"
	// Megabytes unit (MB)
	Megabytes FileSizeUnit = "MB"
	// Gigabytes unit (GB)
	Gigabytes FileSizeUnit = "GB"
	// Terabytes unit (TB)
	Terabytes FileSizeUnit = "TB"
	// Petabytes unit (PB)
	Petabytes FileSizeUnit = "PB"
)

// FileSize represents a file size value object
type FileSize struct {
	bytes uint64
}

// Regular expression for parsing file size strings
var fileSizeRegex = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([KMGTP]?B)$`)

// NewFileSize creates a new FileSize with validation
func NewFileSize(value float64, unit FileSizeUnit) (FileSize, error) {
	// Validate value
	if value < 0 {
		return FileSize{}, errors.New("file size cannot be negative")
	}

	// Convert to bytes
	var bytes uint64
	switch unit {
	case Bytes:
		bytes = uint64(value)
	case Kilobytes:
		bytes = uint64(value * 1024)
	case Megabytes:
		bytes = uint64(value * 1024 * 1024)
	case Gigabytes:
		bytes = uint64(value * 1024 * 1024 * 1024)
	case Terabytes:
		bytes = uint64(value * 1024 * 1024 * 1024 * 1024)
	case Petabytes:
		bytes = uint64(value * 1024 * 1024 * 1024 * 1024 * 1024)
	default:
		return FileSize{}, errors.New("invalid file size unit")
	}

	return FileSize{bytes: bytes}, nil
}

// ParseFileSize creates a new FileSize from a string
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
		return FileSize{}, errors.New("invalid file size format, expected '1024B', '1.5KB', '2.5MB', etc.")
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
func (fs FileSize) String() string {
	return fs.Format("auto")
}

// Equals checks if two FileSizes are equal
func (fs FileSize) Equals(other FileSize) bool {
	return fs.bytes == other.bytes
}

// IsEmpty checks if the FileSize is empty (zero value)
func (fs FileSize) IsEmpty() bool {
	return fs.bytes == 0
}

// Bytes returns the size in bytes
func (fs FileSize) Bytes() uint64 {
	return fs.bytes
}

// Kilobytes returns the size in kilobytes
func (fs FileSize) Kilobytes() float64 {
	return float64(fs.bytes) / 1024
}

// Megabytes returns the size in megabytes
func (fs FileSize) Megabytes() float64 {
	return float64(fs.bytes) / (1024 * 1024)
}

// Gigabytes returns the size in gigabytes
func (fs FileSize) Gigabytes() float64 {
	return float64(fs.bytes) / (1024 * 1024 * 1024)
}

// Terabytes returns the size in terabytes
func (fs FileSize) Terabytes() float64 {
	return float64(fs.bytes) / (1024 * 1024 * 1024 * 1024)
}

// Petabytes returns the size in petabytes
func (fs FileSize) Petabytes() float64 {
	return float64(fs.bytes) / (1024 * 1024 * 1024 * 1024 * 1024)
}

// Add adds another FileSize and returns a new FileSize
func (fs FileSize) Add(other FileSize) FileSize {
	return FileSize{bytes: fs.bytes + other.bytes}
}

// Subtract subtracts another FileSize and returns a new FileSize
// If the result would be negative, it returns zero
func (fs FileSize) Subtract(other FileSize) FileSize {
	if other.bytes > fs.bytes {
		return FileSize{bytes: 0}
	}
	return FileSize{bytes: fs.bytes - other.bytes}
}

// IsLargerThan checks if this FileSize is larger than another FileSize
func (fs FileSize) IsLargerThan(other FileSize) bool {
	return fs.bytes > other.bytes
}

// IsSmallerThan checks if this FileSize is smaller than another FileSize
func (fs FileSize) IsSmallerThan(other FileSize) bool {
	return fs.bytes < other.bytes
}

// Format returns the file size in the specified format
// Format options:
// - "B": Always in bytes
// - "KB": Always in kilobytes
// - "MB": Always in megabytes
// - "GB": Always in gigabytes
// - "TB": Always in terabytes
// - "PB": Always in petabytes
// - "auto": Uses the most appropriate unit for readability
// - "binary": Uses binary prefixes (KiB, MiB, GiB, etc.)
// - "decimal": Uses decimal prefixes (KB, MB, GB, etc.)
func (fs FileSize) Format(format string) string {
	switch format {
	case "B":
		return fmt.Sprintf("%d B", fs.bytes)
	case "KB":
		return fmt.Sprintf("%.2f KB", fs.Kilobytes())
	case "MB":
		return fmt.Sprintf("%.2f MB", fs.Megabytes())
	case "GB":
		return fmt.Sprintf("%.2f GB", fs.Gigabytes())
	case "TB":
		return fmt.Sprintf("%.2f TB", fs.Terabytes())
	case "PB":
		return fmt.Sprintf("%.2f PB", fs.Petabytes())
	case "binary":
		return fs.formatBinary()
	case "decimal":
		return fs.formatDecimal()
	case "auto":
		return fs.formatAuto()
	default:
		return fs.formatAuto()
	}
}

// formatAuto formats the file size using the most appropriate unit for readability
func (fs FileSize) formatAuto() string {
	bytes := float64(fs.bytes)

	if bytes < 1024 {
		return fmt.Sprintf("%d B", fs.bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", fs.Kilobytes())
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", fs.Megabytes())
	}
	if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", fs.Gigabytes())
	}
	if bytes < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f TB", fs.Terabytes())
	}
	return fmt.Sprintf("%.2f PB", fs.Petabytes())
}

// formatBinary formats the file size using binary prefixes (KiB, MiB, GiB, etc.)
func (fs FileSize) formatBinary() string {
	bytes := float64(fs.bytes)

	if bytes < 1024 {
		return fmt.Sprintf("%d B", fs.bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KiB", fs.Kilobytes())
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MiB", fs.Megabytes())
	}
	if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GiB", fs.Gigabytes())
	}
	if bytes < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f TiB", fs.Terabytes())
	}
	return fmt.Sprintf("%.2f PiB", fs.Petabytes())
}

// formatDecimal formats the file size using decimal prefixes (KB, MB, GB, etc.)
func (fs FileSize) formatDecimal() string {
	bytes := float64(fs.bytes)

	if bytes < 1000 {
		return fmt.Sprintf("%d B", fs.bytes)
	}
	if bytes < 1000*1000 {
		kb := float64(fs.bytes) / 1000
		return fmt.Sprintf("%.2f KB", kb)
	}
	if bytes < 1000*1000*1000 {
		mb := float64(fs.bytes) / (1000 * 1000)
		return fmt.Sprintf("%.2f MB", mb)
	}
	if bytes < 1000*1000*1000*1000 {
		gb := float64(fs.bytes) / (1000 * 1000 * 1000)
		return fmt.Sprintf("%.2f GB", gb)
	}
	if bytes < 1000*1000*1000*1000*1000 {
		tb := float64(fs.bytes) / (1000 * 1000 * 1000 * 1000)
		return fmt.Sprintf("%.2f TB", tb)
	}
	pb := float64(fs.bytes) / (1000 * 1000 * 1000 * 1000 * 1000)
	return fmt.Sprintf("%.2f PB", pb)
}
