//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// This example demonstrates how to add and retrieve context information from errors.
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/errors/utils"
)

// UserService is a simple service that demonstrates error context
type UserService struct{}

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Email    string
	Age      int
}

// ValidateUser validates a user and returns an error with context if validation fails
func (s *UserService) ValidateUser(user *User) error {
	if user == nil {
		// Create an error with context
		return errors.WrapWithDetails(
			errors.ErrInvalidInput,
			errors.InvalidInputCode,
			"user cannot be nil",
			map[string]interface{}{
				"operation": "ValidateUser",
				"source":    "UserService",
			},
		)
	}

	if user.Username == "" {
		// Create a validation error with field information
		return errors.NewValidationError("username cannot be empty", "username", nil)
	}

	if user.Email == "" {
		// Create a validation error with field information
		return errors.NewValidationError("email cannot be empty", "email", nil)
	}

	if user.Age < 18 {
		// Create a validation error with field information and additional context
		validationErr := errors.NewValidationError("user must be at least 18 years old", "age", nil)

		// Add additional context to the error
		return errors.WrapWithDetails(
			validationErr,
			errors.ValidationErrorCode,
			"age validation failed",
			map[string]interface{}{
				"provided_age": user.Age,
				"minimum_age":  18,
				"user_id":      user.ID,
			},
		)
	}

	return nil
}

// CreateUser creates a user and demonstrates error context propagation
func (s *UserService) CreateUser(user *User) error {
	// Validate the user
	if err := s.ValidateUser(user); err != nil {
		// Add more context to the error
		return errors.WrapWithDetails(
			err,
			errors.InvalidInputCode,
			"user validation failed during creation",
			map[string]interface{}{
				"operation": "CreateUser",
				"timestamp": "2025-01-01T12:00:00Z",
			},
		)
	}

	// Simulate a database error with context
	return errors.WrapWithDetails(
		errors.NewDatabaseError("failed to insert user", "INSERT", "users", nil),
		errors.DatabaseErrorCode,
		"database operation failed",
		map[string]interface{}{
			"user_id":   user.ID,
			"username":  user.Username,
			"operation": "CreateUser",
			"timestamp": "2025-01-01T12:00:00Z",
		},
	)
}

func main() {
	fmt.Println("Error Context Example")
	fmt.Println("====================")

	// Example 1: Creating an error with context
	fmt.Println("Example 1: Creating an error with context")
	err1 := errors.WrapWithDetails(
		errors.ErrNotFound,
		errors.NotFoundCode,
		"user not found",
		map[string]interface{}{
			"user_id":   "123",
			"operation": "GetUser",
			"timestamp": "2025-01-01T12:00:00Z",
		},
	)
	fmt.Printf("Error: %v\n", err1)

	// Get details from the error
	if e, ok := err1.(interface{ GetDetails() map[string]interface{} }); ok {
		details := e.GetDetails()
		fmt.Println("Details:")
		for k, v := range details {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
	fmt.Println()

	// Example 2: Using the utils.GetDetails function
	fmt.Println("Example 2: Using the utils.GetDetails function")
	err2 := errors.WrapWithDetails(
		errors.ErrTimeout,
		errors.TimeoutCode,
		"operation timed out",
		map[string]interface{}{
			"operation": "FetchData",
			"timeout":   "30s",
			"retries":   3,
		},
	)
	fmt.Printf("Error: %v\n", err2)

	// Get details using the utils package
	details := utils.GetDetails(err2)
	fmt.Println("Details:")
	for k, v := range details {
		fmt.Printf("  %s: %v\n", k, v)
	}
	fmt.Println()

	// Example 3: Error context propagation
	fmt.Println("Example 3: Error context propagation")
	service := &UserService{}

	// Create an invalid user to trigger validation errors
	invalidUser := &User{
		ID:       "456",
		Username: "johndoe",
		Email:    "john@example.com",
		Age:      16, // Too young
	}

	err3 := service.CreateUser(invalidUser)
	fmt.Printf("Error: %v\n", err3)

	// Get details from the propagated error
	details = utils.GetDetails(err3)
	fmt.Println("Details:")
	for k, v := range details {
		fmt.Printf("  %s: %v\n", k, v)
	}
	fmt.Println()

	// Example 4: Getting specific context values
	fmt.Println("Example 4: Getting specific context values")
	if age, exists := utils.GetDetails(err3)["provided_age"]; exists {
		fmt.Printf("Provided age: %v\n", age)
	}

	if minAge, exists := utils.GetDetails(err3)["minimum_age"]; exists {
		fmt.Printf("Minimum age: %v\n", minAge)
	}

	if operation, exists := utils.GetDetails(err3)["operation"]; exists {
		fmt.Printf("Operation: %v\n", operation)
	}
	fmt.Println()

	// Example 5: Adding context to existing errors
	fmt.Println("Example 5: Adding context to existing errors")
	baseErr := errors.New(errors.NetworkErrorCode, "connection failed")
	fmt.Printf("Original error: %v\n", baseErr)

	// Add context to the error
	contextErr := errors.WrapWithDetails(
		baseErr,
		errors.NetworkErrorCode,
		"network operation failed",
		map[string]interface{}{
			"host":    "example.com",
			"port":    8080,
			"timeout": "5s",
		},
	)
	fmt.Printf("Error with context: %v\n", contextErr)

	// Get details from the error
	details = utils.GetDetails(contextErr)
	fmt.Println("Details:")
	for k, v := range details {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

// To run this example:
// go run examples/errors/error_context_example.go
