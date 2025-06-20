// Copyright (c) 2025 A Bit of Help, Inc.

// Example of working with optional dates using the date package
package main

import (
	"fmt"
	"time"
	
	"github.com/abitofhelp/servicelib/date"
)

func main() {
	fmt.Println("Date Package - Optional Date Example")
	fmt.Println("====================================")
	
	// Working with nil date strings
	fmt.Println("\n1. Parsing nil date strings:")
	var nilDateStr *string
	
	parsedNil, err := date.ParseOptionalDate(nilDateStr)
	if err != nil {
		fmt.Printf("Error parsing nil date: %v\n", err)
		return
	}
	
	if parsedNil == nil {
		fmt.Println("Parsed nil date is nil (as expected)")
	} else {
		fmt.Printf("Parsed nil date: %v (unexpected)\n", *parsedNil)
	}
	
	// Working with non-nil date strings
	fmt.Println("\n2. Parsing non-nil date strings:")
	validDateStr := "2023-01-15T14:30:00Z"
	dateStrPtr := &validDateStr
	
	parsed, err := date.ParseOptionalDate(dateStrPtr)
	if err != nil {
		fmt.Printf("Error parsing optional date: %v\n", err)
		return
	}
	
	if parsed == nil {
		fmt.Println("Parsed optional date is nil (unexpected)")
	} else {
		fmt.Printf("Parsed optional date: %v\n", *parsed)
		fmt.Printf("Year: %d\n", parsed.Year())
		fmt.Printf("Month: %s\n", parsed.Month())
		fmt.Printf("Day: %d\n", parsed.Day())
	}
	
	// Formatting nil dates
	fmt.Println("\n3. Formatting nil dates:")
	var nilDate *time.Time
	
	formattedNil := date.FormatOptionalDate(nilDate)
	if formattedNil == nil {
		fmt.Println("Formatted nil date is nil (as expected)")
	} else {
		fmt.Printf("Formatted nil date: %s (unexpected)\n", *formattedNil)
	}
	
	// Formatting non-nil dates
	fmt.Println("\n4. Formatting non-nil dates:")
	now := time.Now()
	nowPtr := &now
	
	formatted := date.FormatOptionalDate(nowPtr)
	if formatted == nil {
		fmt.Println("Formatted optional date is nil (unexpected)")
	} else {
		fmt.Printf("Formatted optional date: %s\n", *formatted)
	}
	
	// Using optional dates in a struct
	fmt.Println("\n5. Using optional dates in a struct:")
	
	type Event struct {
		Name      string
		StartDate time.Time
		EndDate   *time.Time // Optional end date
	}
	
	// Event with no end date
	ongoingEvent := Event{
		Name:      "Ongoing Conference",
		StartDate: time.Now(),
		EndDate:   nil,
	}
	
	// Event with end date
	endTime := time.Now().AddDate(0, 0, 3) // 3 days from now
	completedEvent := Event{
		Name:      "Completed Workshop",
		StartDate: time.Now().AddDate(0, 0, -5), // 5 days ago
		EndDate:   &endTime,
	}
	
	// Display events
	fmt.Printf("Ongoing Event: %s (Started: %s, Ends: %v)\n", 
		ongoingEvent.Name, 
		date.FormatDate(ongoingEvent.StartDate),
		formatOptionalEndDate(ongoingEvent.EndDate))
	
	fmt.Printf("Completed Event: %s (Started: %s, Ends: %v)\n", 
		completedEvent.Name, 
		date.FormatDate(completedEvent.StartDate),
		formatOptionalEndDate(completedEvent.EndDate))
	
	// Expected output:
	// Date Package - Optional Date Example
	// ====================================
	//
	// 1. Parsing nil date strings:
	// Parsed nil date is nil (as expected)
	//
	// 2. Parsing non-nil date strings:
	// Parsed optional date: 2023-01-15 14:30:00 +0000 UTC
	// Year: 2023
	// Month: January
	// Day: 15
	//
	// 3. Formatting nil dates:
	// Formatted nil date is nil (as expected)
	//
	// 4. Formatting non-nil dates:
	// Formatted optional date: 2023-05-10T15:45:30-07:00
	//
	// 5. Using optional dates in a struct:
	// Ongoing Event: Ongoing Conference (Started: 2023-05-10T15:45:30-07:00, Ends: Not scheduled)
	// Completed Event: Completed Workshop (Started: 2023-05-05T15:45:30-07:00, Ends: 2023-05-13T15:45:30-07:00)
}

// Helper function to format an optional end date
func formatOptionalEndDate(endDate *time.Time) string {
	if endDate == nil {
		return "Not scheduled"
	}
	return date.FormatDate(*endDate)
}