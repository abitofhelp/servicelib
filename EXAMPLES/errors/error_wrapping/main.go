// Copyright (c) 2025 A Bit of Help, Inc.

// This example demonstrates how to wrap errors to add context as they propagate up the call stack.
package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/abitofhelp/servicelib/errors"
)

// simulateDBQuery simulates a database query that might fail
func simulateDBQuery(id string) (string, error) {
	// Simulate a database error
	if id == "invalid" {
		return "", sql.ErrNoRows
	}
	return "User data", nil
}

// getUserFromDB attempts to get a user from the database
func getUserFromDB(id string) (string, error) {
	data, err := simulateDBQuery(id)
	if err != nil {
		// Wrap the database error with more context
		return "", errors.NewDatabaseError("failed to get user", "SELECT", "users", err)
	}
	return data, nil
}

// processUserData processes user data from the database
func processUserData(id string) (string, error) {
	data, err := getUserFromDB(id)
	if err != nil {
		// Wrap the error with additional context
		return "", errors.Wrap(err, errors.InternalErrorCode, "error processing user data")
	}
	return data, nil
}

// handleUserRequest handles a user request
func handleUserRequest(id string) (string, error) {
	data, err := processUserData(id)
	if err != nil {
		// Wrap the error with operation information
		return "", errors.WrapWithOperation(err, errors.InternalErrorCode, "request handling failed", "handleUserRequest")
	}
	return data, nil
}

func main() {
	fmt.Println("Error Wrapping Example")
	fmt.Println("======================")

	// Example 1: Basic error wrapping
	err1 := errors.New(errors.InvalidInputCode, "invalid input")
	wrappedErr1 := errors.Wrap(err1, errors.InternalErrorCode, "operation failed")
	fmt.Printf("Original error: %v\n", err1)
	fmt.Printf("Wrapped error: %v\n", wrappedErr1)
	fmt.Println()

	// Example 2: Wrapping with operation
	err2 := errors.New(errors.TimeoutCode, "database timeout")
	wrappedErr2 := errors.WrapWithOperation(err2, errors.DatabaseErrorCode, "query failed", "getUserData")
	fmt.Printf("Original error: %v\n", err2)
	fmt.Printf("Wrapped with operation: %v\n", wrappedErr2)
	fmt.Println()

	// Example 3: Wrapping with details
	err3 := errors.New(errors.NetworkErrorCode, "connection refused")
	details := map[string]interface{}{
		"host":    "example.com",
		"port":    8080,
		"timeout": "30s",
	}
	wrappedErr3 := errors.WrapWithDetails(err3, errors.ExternalServiceErrorCode, "API call failed", details)
	fmt.Printf("Original error: %v\n", err3)
	fmt.Printf("Wrapped with details: %v\n", wrappedErr3)
	fmt.Println()

	// Example 4: Wrapping standard library errors
	fileErr := os.ErrNotExist
	wrappedFileErr := errors.Wrap(fileErr, errors.NotFoundCode, "config file not found")
	fmt.Printf("Original error: %v\n", fileErr)
	fmt.Printf("Wrapped error: %v\n", wrappedFileErr)
	fmt.Println()

	// Example 5: Error chain with multiple wrappings
	fmt.Println("Error chain example:")
	result, err := handleUserRequest("invalid")
	if err != nil {
		fmt.Printf("Final error: %v\n", err)

		// Unwrap the error chain
		fmt.Println("\nUnwrapping error chain:")
		currentErr := err
		level := 1
		for currentErr != nil {
			fmt.Printf("Level %d: %v\n", level, currentErr)
			currentErr = errors.Unwrap(currentErr)
			level++
		}
	} else {
		fmt.Printf("Result: %s\n", result)
	}
	fmt.Println()

	// Example 6: Using errors.Is to check for specific errors
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("The root cause was sql.ErrNoRows")
	}

	// Example 7: Using errors.As to get a specific error type
	var dbErr *errors.DatabaseError
	if errors.As(err, &dbErr) {
		fmt.Printf("Found DatabaseError in the chain: %v\n", dbErr)
		fmt.Printf("Operation: %s, Table: %s\n", dbErr.Operation, dbErr.Table)
	}
}

// To run this example:
// go run examples/errors/error_wrapping_example.go
