// Copyright (c) 2025 A Bit of Help, Inc.

// Example of custom validation using the validation package
package main

import (
	"fmt"
	
	"github.com/abitofhelp/servicelib/validation"
)

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Email    string
	Age      int
	Role     string
}

// ValidateUser validates a user
func ValidateUser(user User) error {
	result := validation.NewValidationResult()
	
	// Validate ID
	validation.ValidateID(user.ID, "id", result)
	
	// Validate username
	validation.Required(user.Username, "username", result)
	validation.MinLength(user.Username, 3, "username", result)
	validation.MaxLength(user.Username, 50, "username", result)
	
	// Validate email
	validation.Required(user.Email, "email", result)
	validation.Pattern(user.Email, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "email", result)
	
	// Validate age
	if user.Age < 18 {
		result.AddError("must be at least 18 years old", "age")
	}
	
	// Validate role
	if user.Role != "admin" && user.Role != "user" && user.Role != "guest" {
		result.AddError("must be one of: admin, user, guest", "role")
	}
	
	return result.Error()
}

func main() {
	// Create an invalid user
	invalidUser := User{
		ID:       "",
		Username: "jo", // Too short
		Email:    "invalid-email", // Invalid format
		Age:      16, // Too young
		Role:     "superuser", // Invalid role
	}
	
	fmt.Println("Validating invalid user:", invalidUser)
	
	// Validate the invalid user
	if err := ValidateUser(invalidUser); err != nil {
		fmt.Printf("User validation failed: %v\n", err)
	} else {
		fmt.Println("User validation passed!")
	}
	
	// Create a valid user
	validUser := User{
		ID:       "123",
		Username: "johndoe",
		Email:    "john.doe@example.com",
		Age:      25,
		Role:     "admin",
	}
	
	fmt.Println("\nValidating valid user:", validUser)
	
	// Validate the valid user
	if err := ValidateUser(validUser); err != nil {
		fmt.Printf("User validation failed: %v\n", err)
	} else {
		fmt.Println("User validation passed!")
	}
	
	// Create a user with some validation errors
	partialUser := User{
		ID:       "456",
		Username: "jane",
		Email:    "jane.doe@example.com",
		Age:      30,
		Role:     "manager", // Invalid role
	}
	
	fmt.Println("\nValidating user with partial errors:", partialUser)
	
	// Validate the user with partial errors
	if err := ValidateUser(partialUser); err != nil {
		fmt.Printf("User validation failed: %v\n", err)
	} else {
		fmt.Println("User validation passed!")
	}
	
	// Expected output:
	// Validating invalid user: { jo invalid-email 16 superuser}
	// User validation failed: validation errors: id: is required, username: must be at least 3 characters long, email: has an invalid format, age: must be at least 18 years old, role: must be one of: admin, user, guest
	//
	// Validating valid user: {123 johndoe john.doe@example.com 25 admin}
	// User validation passed!
	//
	// Validating user with partial errors: {456 jane jane.doe@example.com 30 manager}
	// User validation failed: validation errors: role: must be one of: admin, user, guest
}