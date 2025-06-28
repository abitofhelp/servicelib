// Copyright (c) 2025 A Bit of Help, Inc.

// This example demonstrates how to include and retrieve caller information from errors.
package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/errors/utils"
)

// Function that creates an error at a specific call site
func createError() error {
	return errors.New(errors.InternalErrorCode, "something went wrong")
}

// Function that wraps an error at a specific call site
func wrapError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, errors.InternalErrorCode, "error in wrapError function")
}

// Function that creates a chain of errors
func createErrorChain() error {
	// Create an error
	err := createError()

	// Wrap the error at different call sites
	err = wrapError(err)
	err = errors.Wrap(err, errors.InternalErrorCode, "error in createErrorChain function")

	return err
}

// Function to extract source information from an error
func extractSourceInfo(err error) (string, int) {
	// Check if the error has source information
	if e, ok := err.(interface {
		GetSource() string
		GetLine() int
	}); ok {
		return e.GetSource(), e.GetLine()
	}
	return "unknown", 0
}

// Function to get the current function name
func getCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}

// Function to demonstrate how to use the utils.GetCallerPackage function
func demonstrateCallerPackage() {
	// Get the package name of the caller
	pkg := utils.GetCallerPackage(0)
	fmt.Printf("Current package: %s\n", pkg)

	// Call a function that gets its caller's package
	getCallerPackageWrapper()
}

// Helper function for demonstrateCallerPackage
func getCallerPackageWrapper() {
	// Get the package name of the caller (demonstrateCallerPackage)
	pkg := utils.GetCallerPackage(1)
	fmt.Printf("Caller package: %s\n", pkg)
}

// Function to print error details including source information
func printErrorDetails(err error) {
	fmt.Printf("Error: %v\n", err)

	// Extract source information
	source, line := extractSourceInfo(err)
	fmt.Printf("Source: %s:%d\n", source, line)

	// Convert error to JSON to see all details
	jsonBytes, _ := json.MarshalIndent(struct {
		Message string `json:"message"`
		Source  string `json:"source"`
		Line    int    `json:"line"`
	}{
		Message: err.Error(),
		Source:  source,
		Line:    line,
	}, "", "  ")

	fmt.Printf("JSON: %s\n", string(jsonBytes))
}

func main() {
	fmt.Println("Stack Traces Example")
	fmt.Println("===================")

	// Example 1: Creating an error with automatic source information
	fmt.Println("\nExample 1: Creating an error with automatic source information")
	err1 := errors.New(errors.InternalErrorCode, "simple error")
	printErrorDetails(err1)

	// Example 2: Creating an error in a different function
	fmt.Println("\nExample 2: Creating an error in a different function")
	err2 := createError()
	printErrorDetails(err2)

	// Example 3: Creating an error chain
	fmt.Println("\nExample 3: Creating an error chain")
	err3 := createErrorChain()
	printErrorDetails(err3)

	// Example 4: Unwrapping errors to see the chain
	fmt.Println("\nExample 4: Unwrapping errors to see the chain")
	currentErr := err3
	level := 1
	for currentErr != nil {
		source, line := extractSourceInfo(currentErr)
		fmt.Printf("Level %d: %v (Source: %s:%d)\n", level, currentErr, source, line)
		currentErr = errors.Unwrap(currentErr)
		level++
	}

	// Example 5: Getting the current function name
	fmt.Println("\nExample 5: Getting the current function name")
	funcName := getCurrentFunctionName()
	fmt.Printf("Current function: %s\n", funcName)

	// Example 6: Demonstrating caller package
	fmt.Println("\nExample 6: Demonstrating caller package")
	demonstrateCallerPackage()
}

// To run this example:
// go run examples/errors/stack_traces_example.go
