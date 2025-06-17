// Copyright (c) 2025 A Bit of Help, Inc.

// Package date provides utilities for working with dates and times.
package date

import (
	"time"

	"github.com/abitofhelp/servicelib/errors"
)

const (
	// StandardDateFormat is the standard date format used throughout the application
	StandardDateFormat = time.RFC3339
)

// ParseDate parses a date string in the standard format.
// Parameters:
//   - dateStr: The date string to parse
//
// Returns:
//   - time.Time: The parsed time
//   - error: An error if parsing fails
func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(StandardDateFormat, dateStr)
	if err != nil {
		return time.Time{}, errors.NewValidationError("invalid date format")
	}
	return parsedDate, nil
}

// ParseOptionalDate parses an optional date string in the standard format.
// Parameters:
//   - dateStr: Pointer to the date string to parse, can be nil
//
// Returns:
//   - *time.Time: Pointer to the parsed time, nil if input is nil
//   - error: An error if parsing fails
func ParseOptionalDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}

	parsedDate, err := time.Parse(StandardDateFormat, *dateStr)
	if err != nil {
		return nil, errors.NewValidationError("invalid date format")
	}

	return &parsedDate, nil
}

// FormatDate formats a time.Time as a string in the standard format.
// Parameters:
//   - date: The time to format
//
// Returns:
//   - string: The formatted date string
func FormatDate(date time.Time) string {
	return date.Format(StandardDateFormat)
}

// FormatOptionalDate formats an optional time.Time as a string in the standard format.
// Parameters:
//   - date: Pointer to the time to format, can be nil
//
// Returns:
//   - *string: Pointer to the formatted date string, nil if input is nil
func FormatOptionalDate(date *time.Time) *string {
	if date == nil {
		return nil
	}

	formatted := date.Format(StandardDateFormat)
	return &formatted
}
