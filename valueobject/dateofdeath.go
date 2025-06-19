// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"time"
)

// DateOfDeath represents a date of death value object
type DateOfDeath struct {
	date time.Time
}

// NewDateOfDeath creates a new DateOfDeath with validation
func NewDateOfDeath(year, month, day int) (DateOfDeath, error) {
	// Create date
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate date is not in the future
	if date.After(time.Now()) {
		return DateOfDeath{}, errors.New("date of death cannot be in the future")
	}

	// Validate reasonable age (e.g., not more than 150 years old)
	maxAge := time.Now().AddDate(-150, 0, 0)
	if date.Before(maxAge) {
		return DateOfDeath{}, errors.New("date of death indicates an unreasonable age")
	}

	return DateOfDeath{date: date}, nil
}

// ParseDateOfDeath creates a new DateOfDeath from a string in format "YYYY-MM-DD"
func ParseDateOfDeath(dateStr string) (DateOfDeath, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return DateOfDeath{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Extract year, month, day and create using the constructor for validation
	year, month, day := date.Date()
	return NewDateOfDeath(year, int(month), day)
}

// String returns the string representation of the DateOfDeath
func (dod DateOfDeath) String() string {
	return dod.date.Format("2006-01-02")
}

// Equals checks if two DateOfDeath values are equal
func (dod DateOfDeath) Equals(other DateOfDeath) bool {
	return dod.date.Equal(other.date)
}

// IsEmpty checks if the DateOfDeath is empty (zero value)
func (dod DateOfDeath) IsEmpty() bool {
	return dod.date.IsZero()
}

// TimeSinceDeath calculates the time elapsed since death in years
func (dod DateOfDeath) TimeSinceDeath() int {
	if dod.IsEmpty() {
		return 0
	}

	now := time.Now()
	years := now.Year() - dod.date.Year()

	// Adjust years if death anniversary hasn't occurred yet this year
	if now.Month() < dod.date.Month() || (now.Month() == dod.date.Month() && now.Day() < dod.date.Day()) {
		years--
	}

	return years
}

// IsRecent checks if the death occurred within the last year
func (dod DateOfDeath) IsRecent() bool {
	if dod.IsEmpty() {
		return false
	}
	return dod.TimeSinceDeath() < 1
}

// Date returns the underlying time.Time value
func (dod DateOfDeath) Date() time.Time {
	return dod.date
}

// Format returns the date formatted according to the given layout
func (dod DateOfDeath) Format(layout string) string {
	if dod.IsEmpty() {
		return ""
	}
	return dod.date.Format(layout)
}
