// Copyright (c) 2025 A Bit of Help, Inc.

// Example of working with time zones using the date package
package main

import (
	"fmt"
	"time"
	
	"github.com/abitofhelp/servicelib/date"
)

func main() {
	fmt.Println("Date Package - Time Zone Example")
	fmt.Println("================================")
	
	// 1. Parsing dates with different time zones
	fmt.Println("\n1. Parsing dates with different time zones:")
	
	// UTC time
	utcDateStr := "2023-01-15T14:30:00Z"
	fmt.Printf("Parsing UTC date: %s\n", utcDateStr)
	
	utcDate, err := date.ParseDate(utcDateStr)
	if err != nil {
		fmt.Printf("Error parsing UTC date: %v\n", err)
		return
	}
	
	fmt.Printf("Parsed UTC date: %v\n", utcDate)
	fmt.Printf("Time zone: %s\n", utcDate.Location().String())
	
	// EST time
	estDateStr := "2023-01-15T09:30:00-05:00"
	fmt.Printf("\nParsing EST date: %s\n", estDateStr)
	
	estDate, err := date.ParseDate(estDateStr)
	if err != nil {
		fmt.Printf("Error parsing EST date: %v\n", err)
		return
	}
	
	fmt.Printf("Parsed EST date: %v\n", estDate)
	fmt.Printf("Time zone: %s\n", estDate.Location().String())
	
	// CEST time
	cestDateStr := "2023-01-15T20:30:00+02:00"
	fmt.Printf("\nParsing CEST date: %s\n", cestDateStr)
	
	cestDate, err := date.ParseDate(cestDateStr)
	if err != nil {
		fmt.Printf("Error parsing CEST date: %v\n", err)
		return
	}
	
	fmt.Printf("Parsed CEST date: %v\n", cestDate)
	fmt.Printf("Time zone: %s\n", cestDate.Location().String())
	
	// 2. Comparing dates from different time zones
	fmt.Println("\n2. Comparing dates from different time zones:")
	
	// These dates represent the same instant in time
	fmt.Printf("UTC date: %v\n", utcDate)
	fmt.Printf("EST date: %v\n", estDate)
	
	if utcDate.Equal(estDate) {
		fmt.Println("The UTC and EST dates are equal (represent the same instant)")
	} else {
		fmt.Println("The UTC and EST dates are not equal (unexpected)")
	}
	
	// 3. Converting between time zones
	fmt.Println("\n3. Converting between time zones:")
	
	// Load location for Tokyo
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Printf("Error loading Tokyo location: %v\n", err)
		return
	}
	
	// Convert UTC date to Tokyo time
	tokyoDate := utcDate.In(tokyo)
	fmt.Printf("UTC date: %v\n", utcDate)
	fmt.Printf("Tokyo date: %v\n", tokyoDate)
	fmt.Printf("Tokyo time zone: %s\n", tokyoDate.Location().String())
	
	// 4. Formatting dates with time zones
	fmt.Println("\n4. Formatting dates with time zones:")
	
	// Format the dates using the date package
	fmt.Printf("Formatted UTC date: %s\n", date.FormatDate(utcDate))
	fmt.Printf("Formatted EST date: %s\n", date.FormatDate(estDate))
	fmt.Printf("Formatted Tokyo date: %s\n", date.FormatDate(tokyoDate))
	
	// 5. Working with local time zone
	fmt.Println("\n5. Working with local time zone:")
	
	// Create a date in local time zone
	localDate := time.Now()
	fmt.Printf("Local date: %v\n", localDate)
	fmt.Printf("Local time zone: %s\n", localDate.Location().String())
	
	// Format the local date
	fmt.Printf("Formatted local date: %s\n", date.FormatDate(localDate))
	
	// Convert local date to UTC
	utcLocalDate := localDate.UTC()
	fmt.Printf("Local date in UTC: %v\n", utcLocalDate)
	fmt.Printf("Formatted local date in UTC: %s\n", date.FormatDate(utcLocalDate))
	
	// Expected output (time zones may vary based on your system):
	// Date Package - Time Zone Example
	// ================================
	//
	// 1. Parsing dates with different time zones:
	// Parsing UTC date: 2023-01-15T14:30:00Z
	// Parsed UTC date: 2023-01-15 14:30:00 +0000 UTC
	// Time zone: UTC
	//
	// Parsing EST date: 2023-01-15T09:30:00-05:00
	// Parsed EST date: 2023-01-15 09:30:00 -0500 -0500
	// Time zone: -0500
	//
	// Parsing CEST date: 2023-01-15T20:30:00+02:00
	// Parsed CEST date: 2023-01-15 20:30:00 +0200 +0200
	// Time zone: +0200
	//
	// 2. Comparing dates from different time zones:
	// UTC date: 2023-01-15 14:30:00 +0000 UTC
	// EST date: 2023-01-15 09:30:00 -0500 -0500
	// The UTC and EST dates are equal (represent the same instant)
	//
	// 3. Converting between time zones:
	// UTC date: 2023-01-15 14:30:00 +0000 UTC
	// Tokyo date: 2023-01-15 23:30:00 +0900 JST
	// Tokyo time zone: Asia/Tokyo
	//
	// 4. Formatting dates with time zones:
	// Formatted UTC date: 2023-01-15T14:30:00Z
	// Formatted EST date: 2023-01-15T09:30:00-05:00
	// Formatted Tokyo date: 2023-01-15T23:30:00+09:00
	//
	// 5. Working with local time zone:
	// Local date: 2023-05-10 15:45:30.123456789 -0700 PDT m=+0.000000001
	// Local time zone: Local
	// Formatted local date: 2023-05-10T15:45:30-07:00
	// Local date in UTC: 2023-05-10 22:45:30.123456789 +0000 UTC
	// Formatted local date in UTC: 2023-05-10T22:45:30Z
}