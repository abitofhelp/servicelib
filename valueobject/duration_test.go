// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
	"time"
)

func TestNewDuration(t *testing.T) {
	tests := []struct {
		name        string
		duration    time.Duration
		expectError bool
	}{
		{"Valid Duration", 1 * time.Hour, false},
		{"Zero Duration", 0, false},
		{"Negative Duration", -1 * time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDuration(tt.duration)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if d.Value() != tt.duration {
					t.Errorf("Expected duration %v, got %v", tt.duration, d.Value())
				}
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{"Valid Duration", "1h30m45s", 1*time.Hour + 30*time.Minute + 45*time.Second, false},
		{"Valid Duration with Spaces", " 1h30m45s ", 1*time.Hour + 30*time.Minute + 45*time.Second, false},
		{"Hours Only", "2h", 2 * time.Hour, false},
		{"Minutes Only", "30m", 30 * time.Minute, false},
		{"Seconds Only", "45s", 45 * time.Second, false},
		{"Zero Duration", "0s", 0, false},
		{"Empty String", "", 0, true},
		{"Invalid Format", "1hour", 0, true},
		{"Negative Duration", "-1h", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ParseDuration(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if d.Value() != tt.expected {
					t.Errorf("Expected duration %v, got %v", tt.expected, d.Value())
				}
			}
		})
	}
}

func TestDuration_String(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"1 Hour", 1 * time.Hour, "1h0m0s"},
		{"1 Hour 30 Minutes", 1*time.Hour + 30*time.Minute, "1h30m0s"},
		{"1 Hour 30 Minutes 45 Seconds", 1*time.Hour + 30*time.Minute + 45*time.Second, "1h30m45s"},
		{"Zero Duration", 0, "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := NewDuration(tt.duration)
			if d.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, d.String())
			}
		})
	}
}

func TestDuration_Equals(t *testing.T) {
	d1, _ := NewDuration(1 * time.Hour)
	d2, _ := NewDuration(1 * time.Hour)
	d3, _ := NewDuration(2 * time.Hour)

	if !d1.Equals(d2) {
		t.Errorf("Expected durations to be equal")
	}

	if d1.Equals(d3) {
		t.Errorf("Expected different durations to not be equal")
	}
}

func TestDuration_IsEmpty(t *testing.T) {
	emptyDuration, _ := NewDuration(0)
	nonEmptyDuration, _ := NewDuration(1 * time.Hour)

	if !emptyDuration.IsEmpty() {
		t.Errorf("Expected empty duration to be empty")
	}

	if nonEmptyDuration.IsEmpty() {
		t.Errorf("Expected non-empty duration to not be empty")
	}
}

func TestDuration_Value(t *testing.T) {
	expected := 1*time.Hour + 30*time.Minute + 45*time.Second
	d, _ := NewDuration(expected)

	if d.Value() != expected {
		t.Errorf("Expected value %v, got %v", expected, d.Value())
	}
}

func TestDuration_Conversions(t *testing.T) {
	// 1 hour 30 minutes 45 seconds
	duration := 1*time.Hour + 30*time.Minute + 45*time.Second
	d, _ := NewDuration(duration)

	// Expected values
	expectedHours := 1.5125
	expectedMinutes := 90.75
	expectedSeconds := 5445.0
	expectedMilliseconds := int64(5445000)

	// Test Hours
	if d.Hours() != expectedHours {
		t.Errorf("Expected hours %f, got %f", expectedHours, d.Hours())
	}

	// Test Minutes
	if d.Minutes() != expectedMinutes {
		t.Errorf("Expected minutes %f, got %f", expectedMinutes, d.Minutes())
	}

	// Test Seconds
	if d.Seconds() != expectedSeconds {
		t.Errorf("Expected seconds %f, got %f", expectedSeconds, d.Seconds())
	}

	// Test Milliseconds
	if d.Milliseconds() != expectedMilliseconds {
		t.Errorf("Expected milliseconds %d, got %d", expectedMilliseconds, d.Milliseconds())
	}
}

func TestDuration_Add(t *testing.T) {
	d1, _ := NewDuration(1 * time.Hour)
	d2, _ := NewDuration(30 * time.Minute)
	expected := 1*time.Hour + 30*time.Minute

	result, err := d1.Add(d2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.Value() != expected {
		t.Errorf("Expected %v, got %v", expected, result.Value())
	}
}

func TestDuration_Subtract(t *testing.T) {
	tests := []struct {
		name     string
		d1       time.Duration
		d2       time.Duration
		expected time.Duration
	}{
		{"Normal Subtraction", 1 * time.Hour, 30 * time.Minute, 30 * time.Minute},
		{"Zero Result", 1 * time.Hour, 1 * time.Hour, 0},
		{"Would Be Negative", 30 * time.Minute, 1 * time.Hour, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1, _ := NewDuration(tt.d1)
			d2, _ := NewDuration(tt.d2)

			result, err := d1.Subtract(d2)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result.Value() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result.Value())
			}
		})
	}
}

func TestDuration_Format(t *testing.T) {
	// 1 hour 30 minutes 45 seconds
	duration := 1*time.Hour + 30*time.Minute + 45*time.Second
	d, _ := NewDuration(duration)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"Default Format", "", "1h30m45s"},
		{"Short Format", "short", "1h 30m 45s"},
		{"Long Format", "long", "1 hour 30 minutes 45 seconds"},
		{"Compact Format", "compact", "1:30:45"},
		{"Invalid Format", "invalid", "1h30m45s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Format(tt.format)
			if result != tt.expected {
				t.Errorf("Expected format %s to be %s, got %s", tt.format, tt.expected, result)
			}
		})
	}

	// Test special cases for short format
	t.Run("Short Format - Hours Only", func(t *testing.T) {
		hoursOnly, _ := NewDuration(2 * time.Hour)
		expected := "2h"
		if hoursOnly.Format("short") != expected {
			t.Errorf("Expected %s, got %s", expected, hoursOnly.Format("short"))
		}
	})

	t.Run("Short Format - Minutes Only", func(t *testing.T) {
		minutesOnly, _ := NewDuration(30 * time.Minute)
		expected := "30m"
		if minutesOnly.Format("short") != expected {
			t.Errorf("Expected %s, got %s", expected, minutesOnly.Format("short"))
		}
	})

	t.Run("Short Format - Seconds Only", func(t *testing.T) {
		secondsOnly, _ := NewDuration(45 * time.Second)
		expected := "45s"
		if secondsOnly.Format("short") != expected {
			t.Errorf("Expected %s, got %s", expected, secondsOnly.Format("short"))
		}
	})

	t.Run("Short Format - Zero Duration", func(t *testing.T) {
		zeroDuration, _ := NewDuration(0)
		expected := "0s"
		if zeroDuration.Format("short") != expected {
			t.Errorf("Expected %s, got %s", expected, zeroDuration.Format("short"))
		}
	})

	// Test special cases for long format
	t.Run("Long Format - Singular Units", func(t *testing.T) {
		singularUnits, _ := NewDuration(1*time.Hour + 1*time.Minute + 1*time.Second)
		expected := "1 hour 1 minute 1 second"
		if singularUnits.Format("long") != expected {
			t.Errorf("Expected %s, got %s", expected, singularUnits.Format("long"))
		}
	})

	t.Run("Long Format - Zero Duration", func(t *testing.T) {
		zeroDuration, _ := NewDuration(0)
		expected := "0 seconds"
		if zeroDuration.Format("long") != expected {
			t.Errorf("Expected %s, got %s", expected, zeroDuration.Format("long"))
		}
	})
}
