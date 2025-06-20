// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic validation using the validation package
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/validation"
)

func main() {
	// Create a validation result
	result := validation.NewValidationResult()

	// Validate a username
	username := "john"
	fmt.Println("Validating username:", username)

	// Check if username is required
	validation.Required(username, "username", result)

	// Check if username meets minimum length requirement
	validation.MinLength(username, 5, "username", result)

	// Check if username meets maximum length requirement
	validation.MaxLength(username, 50, "username", result)

	// Validate an email
	email := "john@example.com"
	fmt.Println("Validating email:", email)

	// Check if email is required
	validation.Required(email, "email", result)

	// Check if email matches a pattern
	validation.Pattern(email, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "email", result)

	// Check if validation passed
	if !result.IsValid() {
		fmt.Printf("Validation failed: %v\n", result.Error())
	} else {
		fmt.Println("Validation passed!")
	}

	// Create a new validation result for a valid example
	validResult := validation.NewValidationResult()

	// Validate a valid username
	validUsername := "johndoe"
	fmt.Println("\nValidating valid username:", validUsername)

	// Check if username is required
	validation.Required(validUsername, "username", validResult)

	// Check if username meets minimum length requirement
	validation.MinLength(validUsername, 5, "username", validResult)

	// Check if username meets maximum length requirement
	validation.MaxLength(validUsername, 50, "username", validResult)

	// Validate a valid email
	validEmail := "john.doe@example.com"
	fmt.Println("Validating valid email:", validEmail)

	// Check if email is required
	validation.Required(validEmail, "email", validResult)

	// Check if email matches a pattern
	validation.Pattern(validEmail, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "email", validResult)

	// Check if validation passed
	if !validResult.IsValid() {
		fmt.Printf("Validation failed: %v\n", validResult.Error())
	} else {
		fmt.Println("Validation passed!")
	}

	// Expected output:
	// Validating username: john
	// Validating email: john@example.com
	// Validation failed: validation errors: username: must be at least 5 characters long
	//
	// Validating valid username: johndoe
	// Validating valid email: john.doe@example.com
	// Validation passed!
}
