// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the date package
package example_date

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/date"
)

func main() {
	fmt.Println("Date Package - Basic Usage Example")
	fmt.Println("==================================")

	// Parse a date string
	dateStr := "2023-01-15T14:30:00Z"
	fmt.Printf("\nParsing date string: %s\n", dateStr)

	parsedDate, err := date.ParseDate(dateStr)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return
	}

	fmt.Printf("Parsed date: %v\n", parsedDate)
	fmt.Printf("Year: %d\n", parsedDate.Year())
	fmt.Printf("Month: %s\n", parsedDate.Month())
	fmt.Printf("Day: %d\n", parsedDate.Day())
	fmt.Printf("Hour: %d\n", parsedDate.Hour())
	fmt.Printf("Minute: %d\n", parsedDate.Minute())
	fmt.Printf("Second: %d\n", parsedDate.Second())

	// Format a date
	now := time.Now()
	fmt.Printf("\nFormatting current time: %v\n", now)

	formatted := date.FormatDate(now)
	fmt.Printf("Formatted current time: %s\n", formatted)

	// Parse the formatted date back
	reparsed, err := date.ParseDate(formatted)
	if err != nil {
		fmt.Printf("Error reparsing date: %v\n", err)
		return
	}

	fmt.Printf("\nReparsed date: %v\n", reparsed)

	// Expected output:
	// Date Package - Basic Usage Example
	// ==================================
	//
	// Parsing date string: 2023-01-15T14:30:00Z
	// Parsed date: 2023-01-15 14:30:00 +0000 UTC
	// Year: 2023
	// Month: January
	// Day: 15
	// Hour: 14
	// Minute: 30
	// Second: 0
	//
	// Formatting current time: 2023-05-10 15:45:30.123456789 -0700 PDT m=+0.000000001
	// Formatted current time: 2023-05-10T15:45:30-07:00
	//
	// Reparsed date: 2023-05-10 15:45:30 -0700 PDT
}
