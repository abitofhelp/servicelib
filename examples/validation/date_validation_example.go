// Copyright (c) 2025 A Bit of Help, Inc.

// Example of date validation using the validation package
package main

import (
	"fmt"
	"time"
	
	"github.com/abitofhelp/servicelib/validation"
)

func main() {
	// Create a validation result
	result := validation.NewValidationResult()
	
	// Validate a birth date (must be in the past)
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("Validating birth date:", birthDate.Format("2006-01-02"))
	validation.PastDate(birthDate, "birthDate", result)
	
	// Validate a future date (will fail PastDate validation)
	futureDate := time.Now().AddDate(1, 0, 0) // 1 year in the future
	fmt.Println("Validating future date:", futureDate.Format("2006-01-02"))
	validation.PastDate(futureDate, "futureDate", result)
	
	// Validate a valid date range
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Validating date range: %s to %s\n", 
		startDate.Format("2006-01-02"), 
		endDate.Format("2006-01-02"))
	validation.ValidDateRange(startDate, endDate, "startDate", "endDate", result)
	
	// Validate an invalid date range (start date after end date)
	invalidStartDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	invalidEndDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Validating invalid date range: %s to %s\n", 
		invalidStartDate.Format("2006-01-02"), 
		invalidEndDate.Format("2006-01-02"))
	validation.ValidDateRange(invalidStartDate, invalidEndDate, "invalidStartDate", "invalidEndDate", result)
	
	// Check if validation passed
	if !result.IsValid() {
		fmt.Printf("Validation failed: %v\n", result.Error())
	} else {
		fmt.Println("Validation passed!")
	}
	
	// Create a new validation result for a valid example
	validResult := validation.NewValidationResult()
	
	// Validate a valid birth date (in the past)
	validBirthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("\nValidating valid birth date:", validBirthDate.Format("2006-01-02"))
	validation.PastDate(validBirthDate, "validBirthDate", validResult)
	
	// Validate a valid date range
	validStartDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	validEndDate := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Validating valid date range: %s to %s\n", 
		validStartDate.Format("2006-01-02"), 
		validEndDate.Format("2006-01-02"))
	validation.ValidDateRange(validStartDate, validEndDate, "validStartDate", "validEndDate", validResult)
	
	// Check if validation passed
	if !validResult.IsValid() {
		fmt.Printf("Validation failed: %v\n", validResult.Error())
	} else {
		fmt.Println("Validation passed!")
	}
	
	// Expected output (note: actual future date will vary):
	// Validating birth date: 2000-01-01
	// Validating future date: 2024-05-15
	// Validating date range: 2023-01-01 to 2023-12-31
	// Validating invalid date range: 2023-12-31 to 2023-01-01
	// Validation failed: validation errors: futureDate: must be in the past, invalidStartDate: must be before invalidEndDate
	//
	// Validating valid birth date: 1990-01-01
	// Validating valid date range: 2022-01-01 to 2022-12-31
	// Validation passed!
}