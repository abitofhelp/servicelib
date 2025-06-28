//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// This example demonstrates how to create different types of errors using the errors package.
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/errors"
)

func main() {
	// Example 1: Creating a simple error with errors.New
	simpleErr := errors.New(errors.InternalErrorCode, "something went wrong")
	fmt.Printf("Simple error: %v\n", simpleErr)
	fmt.Printf("Error code: %v\n", simpleErr.(interface{ GetCode() errors.ErrorCode }).GetCode())
	fmt.Println()

	// Example 2: Creating a domain error
	domainErr := errors.NewDomainError(errors.BusinessRuleViolationCode, "business rule violated", nil)
	fmt.Printf("Domain error: %v\n", domainErr)
	fmt.Printf("Error code: %v\n", domainErr.GetCode())
	fmt.Println()

	// Example 3: Creating a validation error
	validationErr := errors.NewValidationError("Email is invalid", "email", nil)
	fmt.Printf("Validation error: %v\n", validationErr)
	fmt.Printf("Field: %v\n", validationErr.Field)
	fmt.Println()

	// Example 4: Creating multiple validation errors
	emailErr := errors.NewValidationError("Email is required", "email", nil)
	passwordErr := errors.NewValidationError("Password must be at least 8 characters", "password", nil)
	validationErrors := errors.NewValidationErrors("Form validation failed", emailErr, passwordErr)
	fmt.Printf("Validation errors: %v\n", validationErrors)
	fmt.Printf("Number of errors: %d\n", len(validationErrors.Errors))
	for i, err := range validationErrors.Errors {
		fmt.Printf("Error %d: %v (Field: %s)\n", i+1, err.Error(), err.Field)
	}
	fmt.Println()

	// Example 5: Creating a not found error
	notFoundErr := errors.NewNotFoundError("User", "123", nil)
	fmt.Printf("Not found error: %v\n", notFoundErr)
	fmt.Printf("Resource type: %v\n", notFoundErr.ResourceType)
	fmt.Printf("Resource ID: %v\n", notFoundErr.ResourceID)
	fmt.Println()

	// Example 6: Creating an infrastructure error
	infraErr := errors.NewInfrastructureError(errors.NetworkErrorCode, "failed to connect to service", nil)
	fmt.Printf("Infrastructure error: %v\n", infraErr)
	fmt.Printf("Error code: %v\n", infraErr.GetCode())
	fmt.Println()

	// Example 7: Creating a database error
	dbErr := errors.NewDatabaseError("failed to insert record", "INSERT", "users", nil)
	fmt.Printf("Database error: %v\n", dbErr)
	fmt.Printf("Operation: %v\n", dbErr.Operation)
	fmt.Printf("Table: %v\n", dbErr.Table)
	fmt.Println()

	// Example 8: Creating an application error
	appErr := errors.NewApplicationError(errors.ConfigurationErrorCode, "invalid configuration", nil)
	fmt.Printf("Application error: %v\n", appErr)
	fmt.Printf("Error code: %v\n", appErr.GetCode())
	fmt.Println()

	// Example 9: Creating a configuration error
	configErr := errors.NewConfigurationError("invalid port number", "PORT", "abc", nil)
	fmt.Printf("Configuration error: %v\n", configErr)
	fmt.Printf("Config key: %v\n", configErr.ConfigKey)
	fmt.Printf("Config value: %v\n", configErr.ConfigValue)
	fmt.Println()

	// Example 10: Creating an authentication error
	authErr := errors.NewAuthenticationError("invalid credentials", "john.doe", nil)
	fmt.Printf("Authentication error: %v\n", authErr)
	fmt.Printf("Username: %v\n", authErr.Username)
	fmt.Println()

	// Example 11: Creating an authorization error
	authzErr := errors.NewAuthorizationError("insufficient permissions", "john.doe", "document", "edit", nil)
	fmt.Printf("Authorization error: %v\n", authzErr)
	fmt.Printf("Username: %v\n", authzErr.Username)
	fmt.Printf("Resource: %v\n", authzErr.Resource)
	fmt.Printf("Action: %v\n", authzErr.Action)
}

// To run this example:
// go run examples/errors/basic_error_creation_example.go
