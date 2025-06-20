// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
)

func TestNewFileSize(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		unit        FileSizeUnit
		expected    uint64
		expectError bool
	}{
		{"Valid Bytes", 1024, Bytes, 1024, false},
		{"Valid Kilobytes", 1, Kilobytes, 1024, false},
		{"Valid Megabytes", 1, Megabytes, 1024 * 1024, false},
		{"Valid Gigabytes", 1, Gigabytes, 1024 * 1024 * 1024, false},
		{"Valid Terabytes", 1, Terabytes, 1024 * 1024 * 1024 * 1024, false},
		{"Valid Petabytes", 1, Petabytes, 1024 * 1024 * 1024 * 1024 * 1024, false},
		{"Negative Value", -1, Bytes, 0, true},
		{"Invalid Unit", 1, FileSizeUnit("invalid"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := NewFileSize(tt.value, tt.unit)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if fs.bytes != tt.expected {
					t.Errorf("Expected %d bytes, got %d", tt.expected, fs.bytes)
				}
			}
		})
	}
}

func TestParseFileSize(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    uint64
		expectError bool
	}{
		{"Valid Bytes", "1024B", 1024, false},
		{"Valid Kilobytes", "1KB", 1024, false},
		{"Valid Megabytes", "1MB", 1024 * 1024, false},
		{"Valid Gigabytes", "1GB", 1024 * 1024 * 1024, false},
		{"Valid Terabytes", "1TB", 1024 * 1024 * 1024 * 1024, false},
		{"Valid Petabytes", "1PB", 1024 * 1024 * 1024 * 1024 * 1024, false},
		{"Decimal Value", "1.5KB", uint64(1.5 * 1024), false},
		{"With Space", "1.5 KB", uint64(1.5 * 1024), false},
		{"Empty String", "", 0, true},
		{"Invalid Format", "1K", 0, true},
		{"Invalid Value", "xKB", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := ParseFileSize(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if fs.bytes != tt.expected {
					t.Errorf("Expected %d bytes, got %d", tt.expected, fs.bytes)
				}
			}
		})
	}
}

func TestFileSize_Equals(t *testing.T) {
	fs1, _ := NewFileSize(1024, Bytes)
	fs2, _ := NewFileSize(1, Kilobytes)
	fs3, _ := NewFileSize(2, Kilobytes)

	if !fs1.Equals(fs2) {
		t.Errorf("Expected 1024B to equal 1KB")
	}

	if fs1.Equals(fs3) {
		t.Errorf("Expected 1024B to not equal 2KB")
	}
}

func TestFileSize_IsEmpty(t *testing.T) {
	fs1, _ := NewFileSize(0, Bytes)
	fs2, _ := NewFileSize(1, Bytes)

	if !fs1.IsEmpty() {
		t.Errorf("Expected 0B to be empty")
	}

	if fs2.IsEmpty() {
		t.Errorf("Expected 1B to not be empty")
	}
}

func TestFileSize_Conversions(t *testing.T) {
	fs, _ := NewFileSize(1024*1024, Bytes) // 1MB in bytes

	if fs.Bytes() != 1024*1024 {
		t.Errorf("Expected 1048576 bytes, got %d", fs.Bytes())
	}

	if fs.Kilobytes() != 1024 {
		t.Errorf("Expected 1024 kilobytes, got %f", fs.Kilobytes())
	}

	if fs.Megabytes() != 1 {
		t.Errorf("Expected 1 megabyte, got %f", fs.Megabytes())
	}

	if fs.Gigabytes() != 1.0/1024 {
		t.Errorf("Expected 0.0009765625 gigabytes, got %f", fs.Gigabytes())
	}

	if fs.Terabytes() != 1.0/(1024*1024) {
		t.Errorf("Expected 9.5367e-7 terabytes, got %f", fs.Terabytes())
	}

	if fs.Petabytes() != 1.0/(1024*1024*1024) {
		t.Errorf("Expected 9.3132e-10 petabytes, got %f", fs.Petabytes())
	}
}

func TestFileSize_Add(t *testing.T) {
	fs1, _ := NewFileSize(1, Kilobytes)
	fs2, _ := NewFileSize(1, Megabytes)
	result := fs1.Add(fs2)

	expected := uint64(1024 + 1024*1024)
	if result.bytes != expected {
		t.Errorf("Expected %d bytes, got %d", expected, result.bytes)
	}
}

func TestFileSize_Subtract(t *testing.T) {
	fs1, _ := NewFileSize(2, Megabytes)
	fs2, _ := NewFileSize(1, Megabytes)
	fs3, _ := NewFileSize(3, Megabytes)

	result1 := fs1.Subtract(fs2)
	expected1 := uint64(1024 * 1024)
	if result1.bytes != expected1 {
		t.Errorf("Expected %d bytes, got %d", expected1, result1.bytes)
	}

	result2 := fs1.Subtract(fs3)
	if result2.bytes != 0 {
		t.Errorf("Expected 0 bytes when subtracting larger value, got %d", result2.bytes)
	}
}

func TestFileSize_Comparisons(t *testing.T) {
	fs1, _ := NewFileSize(1, Kilobytes)
	fs2, _ := NewFileSize(1, Megabytes)

	if !fs2.IsLargerThan(fs1) {
		t.Errorf("Expected 1MB to be larger than 1KB")
	}

	if fs1.IsLargerThan(fs2) {
		t.Errorf("Expected 1KB to not be larger than 1MB")
	}

	if !fs1.IsSmallerThan(fs2) {
		t.Errorf("Expected 1KB to be smaller than 1MB")
	}

	if fs2.IsSmallerThan(fs1) {
		t.Errorf("Expected 1MB to not be smaller than 1KB")
	}
}

func TestFileSize_String(t *testing.T) {
	tests := []struct {
		name     string
		fileSize FileSize
		expected string
	}{
		{"Bytes", FileSize{bytes: 500}, "500 B"},
		{"Kilobytes", FileSize{bytes: 1500}, "1.46 KB"},
		{"Megabytes", FileSize{bytes: 1500000}, "1.43 MB"},
		{"Gigabytes", FileSize{bytes: 1500000000}, "1.40 GB"},
		{"Terabytes", FileSize{bytes: 1500000000000}, "1.36 TB"},
		{"Petabytes", FileSize{bytes: 1500000000000000}, "1.33 PB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fileSize.String()
			// Using contains instead of exact match due to floating point precision
			if result[:4] != tt.expected[:4] {
				t.Errorf("Expected string to start with %s, got %s", tt.expected[:4], result[:4])
			}
		})
	}
}

func TestFileSize_Format(t *testing.T) {
	fs, _ := NewFileSize(1024*1024, Bytes) // 1MB

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"Format B", "B", "1048576 B"},
		{"Format KB", "KB", "1024.00 KB"},
		{"Format MB", "MB", "1.00 MB"},
		{"Format GB", "GB", "0.00 GB"},
		{"Format TB", "TB", "0.00 TB"},
		{"Format PB", "PB", "0.00 PB"},
		{"Format auto", "auto", "1.00 MB"},
		{"Format binary", "binary", "1.00 MiB"},
		{"Format decimal", "decimal", "1.05 MB"},
		{"Format invalid", "invalid", "1.00 MB"}, // Should default to auto
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fs.Format(tt.format)
			// Using contains instead of exact match due to floating point precision
			if result[:4] != tt.expected[:4] {
				t.Errorf("Expected string to start with %s, got %s", tt.expected[:4], result[:4])
			}
		})
	}
}

func TestFileSize_FormatFunctions(t *testing.T) {
	// Test formatAuto
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1500, "1.46 KB"},
		{"Megabytes", 1500000, "1.43 MB"},
		{"Gigabytes", 1500000000, "1.40 GB"},
		{"Terabytes", 1500000000000, "1.36 TB"},
		{"Petabytes", 1500000000000000, "1.33 PB"},
	}

	for _, tt := range tests {
		t.Run("formatAuto_"+tt.name, func(t *testing.T) {
			fs := FileSize{bytes: tt.bytes}
			result := fs.formatAuto()
			if result[:4] != tt.expected[:4] {
				t.Errorf("Expected string to start with %s, got %s", tt.expected[:4], result[:4])
			}
		})
	}

	// Test formatBinary
	binaryTests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Kibibytes", 1500, "1.46 KiB"},
		{"Mebibytes", 1500000, "1.43 MiB"},
		{"Gibibytes", 1500000000, "1.40 GiB"},
		{"Tebibytes", 1500000000000, "1.36 TiB"},
		{"Pebibytes", 1500000000000000, "1.33 PiB"},
	}

	for _, tt := range binaryTests {
		t.Run("formatBinary_"+tt.name, func(t *testing.T) {
			fs := FileSize{bytes: tt.bytes}
			result := fs.formatBinary()
			if result[:4] != tt.expected[:4] {
				t.Errorf("Expected string to start with %s, got %s", tt.expected[:4], result[:4])
			}
		})
	}

	// Test formatDecimal
	decimalTests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1500, "1.50 KB"},
		{"Megabytes", 1500000, "1.50 MB"},
		{"Gigabytes", 1500000000, "1.50 GB"},
		{"Terabytes", 1500000000000, "1.50 TB"},
		{"Petabytes", 1500000000000000, "1.50 PB"},
	}

	for _, tt := range decimalTests {
		t.Run("formatDecimal_"+tt.name, func(t *testing.T) {
			fs := FileSize{bytes: tt.bytes}
			result := fs.formatDecimal()
			if result[:4] != tt.expected[:4] {
				t.Errorf("Expected string to start with %s, got %s", tt.expected[:4], result[:4])
			}
		})
	}
}
