// Copyright (c) 2025 A Bit of Help, Inc.

package identification

import (
	"testing"
	"time"
)

func TestNewDateOfDeath(t *testing.T) {
	// Get current time for testing
	now := time.Now()

	tests := []struct {
		name        string
		year        int
		month       int
		day         int
		expectError bool
	}{
		{"Valid Date", 1990, 1, 1, false},
		{"Valid Recent Date", now.Year() - 1, int(now.Month()), now.Day(), false},
		{"Future Date", now.Year() + 1, int(now.Month()), now.Day(), true},
		{"Too Old Date", 1800, 1, 1, true},
		{"Edge Case - Max Age", now.Year() - 150, int(now.Month()), now.Day(), true},
		{"Edge Case - Just Beyond Max Age", now.Year() - 151, int(now.Month()), now.Day() + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dod, err := NewDateOfDeath(tt.year, tt.month, tt.day)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify the date was set correctly
				expected := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
				if !dod.date.Equal(expected) {
					t.Errorf("Expected date %v, got %v", expected, dod.date)
				}
			}
		})
	}
}

func TestParseDateOfDeath(t *testing.T) {
	// Get current time for testing
	now := time.Now()

	tests := []struct {
		name        string
		dateStr     string
		expectError bool
	}{
		{"Valid Date", "1990-01-01", false},
		{"Valid Recent Date", now.AddDate(-1, 0, 0).Format("2006-01-02"), false},
		{"Future Date", now.AddDate(1, 0, 0).Format("2006-01-02"), true},
		{"Too Old Date", "1800-01-01", true},
		{"Invalid Format", "01/01/1990", true},
		{"Invalid Date", "1990-13-01", true},
		{"Invalid Date - Non-existent", "1990-02-30", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dod, err := ParseDateOfDeath(tt.dateStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify the string representation matches the input
				if dod.String() != tt.dateStr {
					t.Errorf("Expected string %s, got %s", tt.dateStr, dod.String())
				}
			}
		})
	}
}

func TestDateOfDeath_String(t *testing.T) {
	dod, _ := NewDateOfDeath(1990, 1, 1)
	expected := "1990-01-01"

	if dod.String() != expected {
		t.Errorf("Expected string %s, got %s", expected, dod.String())
	}
}

func TestDateOfDeath_Equals(t *testing.T) {
	dod1, _ := NewDateOfDeath(1990, 1, 1)
	dod2, _ := NewDateOfDeath(1990, 1, 1)
	dod3, _ := NewDateOfDeath(1991, 1, 1)

	if !dod1.Equals(dod2) {
		t.Errorf("Expected dates to be equal")
	}

	if dod1.Equals(dod3) {
		t.Errorf("Expected different dates to not be equal")
	}
}

func TestDateOfDeath_IsEmpty(t *testing.T) {
	dod, _ := NewDateOfDeath(1990, 1, 1)
	emptyDod := DateOfDeath{}

	if dod.IsEmpty() {
		t.Errorf("Expected non-empty date of death")
	}

	if !emptyDod.IsEmpty() {
		t.Errorf("Expected empty date of death")
	}
}

func TestDateOfDeath_TimeSinceDeath(t *testing.T) {
	// Current time for testing
	now := time.Now()

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected int
	}{
		{"Time Since Death - Anniversary Passed This Year", now.Year() - 30, int(now.Month()) - 1, now.Day(), 30},
		{"Time Since Death - Anniversary Not Yet This Year", now.Year() - 30, int(now.Month()) + 1, now.Day(), 29},
		{"Time Since Death - Anniversary Today", now.Year() - 30, int(now.Month()), now.Day(), 30},
		{"Time Since Death - This Year", now.Year(), int(now.Month()) - 1, now.Day(), 0},
		{"Time Since Death - Empty DOD", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dod DateOfDeath
			var err error

			if tt.year == 0 && tt.month == 0 && tt.day == 0 {
				dod = DateOfDeath{} // Empty DOD
			} else {
				dod, err = NewDateOfDeath(tt.year, tt.month, tt.day)
				if err != nil {
					t.Fatalf("Failed to create DateOfDeath: %v", err)
				}
			}

			timeSince := dod.TimeSinceDeath()
			if timeSince != tt.expected {
				t.Errorf("Expected time since death %d, got %d", tt.expected, timeSince)
			}
		})
	}
}

func TestDateOfDeath_IsRecent(t *testing.T) {
	// Current time for testing
	now := time.Now()

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected bool
	}{
		{"Recent Death - Less Than 1 Year", now.Year(), int(now.Month()) - 1, now.Day(), true},
		{"Recent Death - Today", now.Year(), int(now.Month()), now.Day(), true},
		{"Not Recent Death - More Than 1 Year", now.Year() - 2, int(now.Month()), now.Day(), false},
		{"Not Recent Death - Empty DOD", 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dod DateOfDeath
			var err error

			if tt.year == 0 && tt.month == 0 && tt.day == 0 {
				dod = DateOfDeath{} // Empty DOD
			} else {
				dod, err = NewDateOfDeath(tt.year, tt.month, tt.day)
				if err != nil {
					t.Fatalf("Failed to create DateOfDeath: %v", err)
				}
			}

			isRecent := dod.IsRecent()
			if isRecent != tt.expected {
				t.Errorf("Expected IsRecent to return %v, got %v", tt.expected, isRecent)
			}
		})
	}
}

func TestDateOfDeath_Date(t *testing.T) {
	year, month, day := 1990, 1, 1
	dod, _ := NewDateOfDeath(year, month, day)

	expected := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if !dod.Date().Equal(expected) {
		t.Errorf("Expected date %v, got %v", expected, dod.Date())
	}
}

func TestDateOfDeath_Format(t *testing.T) {
	dod, _ := NewDateOfDeath(1990, 1, 1)
	emptyDod := DateOfDeath{}

	tests := []struct {
		name     string
		dod      DateOfDeath
		layout   string
		expected string
	}{
		{"ISO Format", dod, "2006-01-02", "1990-01-01"},
		{"US Format", dod, "01/02/2006", "01/01/1990"},
		{"Custom Format", dod, "Jan 2, 2006", "Jan 1, 1990"},
		{"Empty DOD", emptyDod, "2006-01-02", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := tt.dod.Format(tt.layout)
			if formatted != tt.expected {
				t.Errorf("Expected format %s, got %s", tt.expected, formatted)
			}
		})
	}
}
