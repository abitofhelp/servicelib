// Copyright (c) 2025 A Bit of Help, Inc.

// Example of error handling with the date package
package example_date

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/date"
	"github.com/abitofhelp/servicelib/errors"
)

func main() {
	fmt.Println("Date Package - Error Handling Example")
	fmt.Println("=====================================")

	// 1. Handling invalid date formats
	fmt.Println("\n1. Handling invalid date formats:")

	invalidDateStr := "not-a-date"
	fmt.Printf("Trying to parse invalid date: %s\n", invalidDateStr)

	_, err := date.ParseDate(invalidDateStr)
	if err != nil {
		fmt.Printf("Error (as expected): %v\n", err)

		// Check if it's a validation error
		var validationErr *errors.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Println("Detected validation error type")
			fmt.Printf("Error message: %s\n", validationErr.Error())
		}
	} else {
		fmt.Println("No error returned (unexpected)")
	}

	// 2. Handling invalid optional date formats
	fmt.Println("\n2. Handling invalid optional date formats:")

	invalidOptStr := "also-not-a-date"
	invalidOptStrPtr := &invalidOptStr
	fmt.Printf("Trying to parse invalid optional date: %s\n", invalidOptStr)

	_, err = date.ParseOptionalDate(invalidOptStrPtr)
	if err != nil {
		fmt.Printf("Error (as expected): %v\n", err)
	} else {
		fmt.Println("No error returned (unexpected)")
	}

	// 3. Graceful error handling in a function
	fmt.Println("\n3. Graceful error handling in a function:")

	// Try with valid date
	validResult := parseAndDisplayDate("2023-01-15T14:30:00Z")
	fmt.Printf("Valid date result: %v\n", validResult)

	// Try with invalid date
	invalidResult := parseAndDisplayDate("invalid-date-format")
	fmt.Printf("Invalid date result: %v\n", invalidResult)

	// 4. Handling multiple date parsing attempts
	fmt.Println("\n4. Handling multiple date parsing attempts:")

	dateStrings := []string{
		"2023-01-15T14:30:00Z",      // Valid
		"not-a-date",                // Invalid
		"2023-05-20T10:15:30+02:00", // Valid
		"2023/05/20",                // Invalid
	}

	for i, ds := range dateStrings {
		fmt.Printf("\nAttempt %d - Parsing: %s\n", i+1, ds)
		if parsedDate, err := date.ParseDate(ds); err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Success: %v\n", parsedDate)
		}
	}

	// Expected output:
	// Date Package - Error Handling Example
	// =====================================
	//
	// 1. Handling invalid date formats:
	// Trying to parse invalid date: not-a-date
	// Error (as expected): invalid date format
	// Detected validation error type
	// Error message: invalid date format
	//
	// 2. Handling invalid optional date formats:
	// Trying to parse invalid optional date: also-not-a-date
	// Error (as expected): invalid date format
	//
	// 3. Graceful error handling in a function:
	// Valid date result: 2023-01-15 14:30:00 +0000 UTC
	// Invalid date result: <zero date>
	//
	// 4. Handling multiple date parsing attempts:
	//
	// Attempt 1 - Parsing: 2023-01-15T14:30:00Z
	//   Success: 2023-01-15 14:30:00 +0000 UTC
	//
	// Attempt 2 - Parsing: not-a-date
	//   Error: invalid date format
	//
	// Attempt 3 - Parsing: 2023-05-20T10:15:30+02:00
	//   Success: 2023-05-20 10:15:30 +0200 CEST
	//
	// Attempt 4 - Parsing: 2023/05/20
	//   Error: invalid date format
}

// parseAndDisplayDate attempts to parse a date string and returns a default value on error
func parseAndDisplayDate(dateStr string) time.Time {
	parsedDate, err := date.ParseDate(dateStr)
	if err != nil {
		fmt.Printf("Error parsing date '%s': %v\n", dateStr, err)
		return time.Time{} // Return zero value for time.Time
	}
	return parsedDate
}
