// Copyright (c) 2025 A Bit of Help, Inc.

// Package date provides utilities for working with dates and times in a consistent manner.
//
// This package standardizes date and time handling across the application by providing
// a set of functions for parsing, formatting, and manipulating dates. It uses a standard
// date format (RFC3339) to ensure consistency in date representations.
//
// The package addresses common challenges in date handling:
//   - Consistent date formatting across the application
//   - Handling of optional (nullable) dates
//   - Proper error handling for date parsing failures
//
// Key components:
//   - StandardDateFormat: The standard date format constant (RFC3339)
//   - Parse functions: For converting strings to time.Time values
//   - Format functions: For converting time.Time values to strings
//   - Support for both required and optional (pointer-based) dates
//
// Example usage:
//
//	// Parsing a date string
//	dateStr := "2023-04-15T14:30:00Z"
//	parsedDate, err := date.ParseDate(dateStr)
//	if err != nil {
//	    log.Fatalf("Failed to parse date: %v", err)
//	}
//
//	// Working with the parsed date
//	if parsedDate.After(time.Now()) {
//	    fmt.Println("The date is in the future")
//	}
//
//	// Formatting a date
//	formattedDate := date.FormatDate(parsedDate)
//	fmt.Println("Formatted date:", formattedDate)
//
//	// Working with optional dates
//	var optionalDateStr *string
//	// optionalDateStr is nil
//	optionalDate, _ := date.ParseOptionalDate(optionalDateStr)
//	// optionalDate will be nil
//
//	// Setting a value
//	dateValue := "2023-05-20T10:15:00Z"
//	optionalDateStr = &dateValue
//	optionalDate, _ = date.ParseOptionalDate(optionalDateStr)
//	// optionalDate will contain the parsed date
//
//	// Formatting an optional date
//	formattedOptional := date.FormatOptionalDate(optionalDate)
//	// formattedOptional will be a pointer to the formatted string
package date
