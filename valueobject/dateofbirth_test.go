// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
	"time"
)

func TestNewDateOfBirth(t *testing.T) {
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
			dob, err := NewDateOfBirth(tt.year, tt.month, tt.day)

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
				if !dob.date.Equal(expected) {
					t.Errorf("Expected date %v, got %v", expected, dob.date)
				}
			}
		})
	}
}

func TestParseDateOfBirth(t *testing.T) {
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
			dob, err := ParseDateOfBirth(tt.dateStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify the string representation matches the input
				if dob.String() != tt.dateStr {
					t.Errorf("Expected string %s, got %s", tt.dateStr, dob.String())
				}
			}
		})
	}
}

func TestDateOfBirth_String(t *testing.T) {
	dob, _ := NewDateOfBirth(1990, 1, 1)
	expected := "1990-01-01"

	if dob.String() != expected {
		t.Errorf("Expected string %s, got %s", expected, dob.String())
	}
}

func TestDateOfBirth_Equals(t *testing.T) {
	dob1, _ := NewDateOfBirth(1990, 1, 1)
	dob2, _ := NewDateOfBirth(1990, 1, 1)
	dob3, _ := NewDateOfBirth(1991, 1, 1)

	if !dob1.Equals(dob2) {
		t.Errorf("Expected dates to be equal")
	}

	if dob1.Equals(dob3) {
		t.Errorf("Expected different dates to not be equal")
	}
}

func TestDateOfBirth_IsEmpty(t *testing.T) {
	dob, _ := NewDateOfBirth(1990, 1, 1)
	emptyDob := DateOfBirth{}

	if dob.IsEmpty() {
		t.Errorf("Expected non-empty date of birth")
	}

	if !emptyDob.IsEmpty() {
		t.Errorf("Expected empty date of birth")
	}
}

func TestDateOfBirth_Age(t *testing.T) {
	// Current time for testing
	now := time.Now()

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected int
	}{
		{"Age - Birthday Passed This Year", now.Year() - 30, int(now.Month()) - 1, now.Day(), 30},
		{"Age - Birthday Not Yet This Year", now.Year() - 30, int(now.Month()) + 1, now.Day(), 29},
		{"Age - Birthday Today", now.Year() - 30, int(now.Month()), now.Day(), 30},
		{"Age - Born This Year", now.Year(), int(now.Month()) - 1, now.Day(), 0},
		{"Age - Empty DOB", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dob DateOfBirth
			var err error

			if tt.year == 0 && tt.month == 0 && tt.day == 0 {
				dob = DateOfBirth{} // Empty DOB
			} else {
				dob, err = NewDateOfBirth(tt.year, tt.month, tt.day)
				if err != nil {
					t.Fatalf("Failed to create DateOfBirth: %v", err)
				}
			}

			age := dob.Age()
			if age != tt.expected {
				t.Errorf("Expected age %d, got %d", tt.expected, age)
			}
		})
	}
}

func TestDateOfBirth_IsAdult(t *testing.T) {
	// Current time for testing
	now := time.Now()

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected bool
	}{
		{"Adult - Over 18", now.Year() - 19, int(now.Month()), now.Day(), true},
		{"Adult - Exactly 18", now.Year() - 18, int(now.Month()), now.Day(), true},
		{"Not Adult - Under 18", now.Year() - 17, int(now.Month()), now.Day(), false},
		{"Not Adult - Empty DOB", 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dob DateOfBirth
			var err error

			if tt.year == 0 && tt.month == 0 && tt.day == 0 {
				dob = DateOfBirth{} // Empty DOB
			} else {
				dob, err = NewDateOfBirth(tt.year, tt.month, tt.day)
				if err != nil {
					t.Fatalf("Failed to create DateOfBirth: %v", err)
				}
			}

			isAdult := dob.IsAdult()
			if isAdult != tt.expected {
				t.Errorf("Expected IsAdult to return %v, got %v", tt.expected, isAdult)
			}
		})
	}
}

func TestDateOfBirth_Date(t *testing.T) {
	year, month, day := 1990, 1, 1
	dob, _ := NewDateOfBirth(year, month, day)

	expected := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if !dob.Date().Equal(expected) {
		t.Errorf("Expected date %v, got %v", expected, dob.Date())
	}
}

func TestDateOfBirth_Format(t *testing.T) {
	dob, _ := NewDateOfBirth(1990, 1, 1)
	emptyDob := DateOfBirth{}

	tests := []struct {
		name     string
		dob      DateOfBirth
		layout   string
		expected string
	}{
		{"ISO Format", dob, "2006-01-02", "1990-01-01"},
		{"US Format", dob, "01/02/2006", "01/01/1990"},
		{"Custom Format", dob, "Jan 2, 2006", "Jan 1, 1990"},
		{"Empty DOB", emptyDob, "2006-01-02", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := tt.dob.Format(tt.layout)
			if formatted != tt.expected {
				t.Errorf("Expected format %s, got %s", tt.expected, formatted)
			}
		})
	}
}
