// Copyright (c) 2025 A Bit of Help, Inc.

// Package identification provides value objects related to identification.
package identification

import (
	"errors"
	"time"
)

// DateOfBirth represents a date of birth value object
type DateOfBirth struct {
	date time.Time
}

// NewDateOfBirth creates a new DateOfBirth with validation
func NewDateOfBirth(year, month, day int) (DateOfBirth, error) {
	// Create date
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate date is not in the future
	if date.After(time.Now()) {
		return DateOfBirth{}, errors.New("date of birth cannot be in the future")
	}

	// Validate reasonable age (e.g., not more than 150 years old)
	maxAge := time.Now().AddDate(-150, 0, 0)
	if date.Before(maxAge) {
		return DateOfBirth{}, errors.New("date of birth indicates an unreasonable age")
	}

	return DateOfBirth{date: date}, nil
}

// ParseDateOfBirth creates a new DateOfBirth from a string in format "YYYY-MM-DD"
func ParseDateOfBirth(dateStr string) (DateOfBirth, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return DateOfBirth{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Extract year, month, day and create using the constructor for validation
	year, month, day := date.Date()
	return NewDateOfBirth(year, int(month), day)
}

// String returns the string representation of the DateOfBirth
func (dob DateOfBirth) String() string {
	return dob.date.Format("2006-01-02")
}

// Equals checks if two DateOfBirth values are equal
func (dob DateOfBirth) Equals(other DateOfBirth) bool {
	return dob.date.Equal(other.date)
}

// IsEmpty checks if the DateOfBirth is empty (zero value)
func (dob DateOfBirth) IsEmpty() bool {
	return dob.date.IsZero()
}

// Age calculates the age based on the date of birth and current date
func (dob DateOfBirth) Age() int {
	if dob.IsEmpty() {
		return 0
	}

	now := time.Now()
	years := now.Year() - dob.date.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < dob.date.Month() || (now.Month() == dob.date.Month() && now.Day() < dob.date.Day()) {
		years--
	}

	return years
}

// IsAdult checks if the person is an adult (18 years or older)
func (dob DateOfBirth) IsAdult() bool {
	return dob.Age() >= 18
}

// Date returns the underlying time.Time value
func (dob DateOfBirth) Date() time.Time {
	return dob.date
}

// Format returns the date formatted according to the given layout
func (dob DateOfBirth) Format(layout string) string {
	if dob.IsEmpty() {
		return ""
	}
	return dob.date.Format(layout)
}