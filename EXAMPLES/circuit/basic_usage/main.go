//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating basic usage of the circuit package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/circuit"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a circuit breaker with default configuration
	cfg := circuit.DefaultConfig()
	options := circuit.DefaultOptions().WithName("example")
	cb := circuit.NewCircuitBreaker(cfg, options)

	// Define a function to execute with circuit breaking
	fn := func(ctx context.Context) (string, error) {
		// Your operation here
		// For demonstration, we'll simulate a successful operation
		fmt.Println("Executing operation...")
		return "success", nil
	}

	// Execute the function with circuit breaking
	fmt.Println("Executing function with circuit breaking...")
	result, err := circuit.Execute(ctx, cb, "example-operation", fn)
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Printf("Operation succeeded with result: %s\n", result)
	}

	// Define a function that will fail
	failFn := func(ctx context.Context) (string, error) {
		// Simulate a failure
		fmt.Println("Executing operation that will fail...")
		return "", fmt.Errorf("simulated failure")
	}

	// Execute the failing function with circuit breaking
	fmt.Println("\nExecuting failing function with circuit breaking...")
	result, err = circuit.Execute(ctx, cb, "failing-operation", failFn)
	if err != nil {
		fmt.Printf("Operation failed as expected: %v\n", err)
	} else {
		fmt.Printf("Operation unexpectedly succeeded with result: %s\n", result)
	}

	// Execute the failing function multiple times to trip the circuit
	fmt.Println("\nExecuting failing function multiple times to trip the circuit...")
	for i := 0; i < 10; i++ {
		result, err = circuit.Execute(ctx, cb, "failing-operation", failFn)
		if err != nil {
			fmt.Printf("Attempt %d: Operation failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Attempt %d: Operation succeeded with result: %s\n", i+1, result)
		}
	}

	// Check the circuit state
	fmt.Printf("\nCircuit state: %s\n", cb.GetState())

	// Try to execute a successful function after the circuit is open
	fmt.Println("\nTrying to execute a successful function after the circuit is open...")
	result, err = circuit.Execute(ctx, cb, "example-operation", fn)
	if err != nil {
		fmt.Printf("Operation failed because circuit is open: %v\n", err)
	} else {
		fmt.Printf("Operation succeeded with result: %s\n", result)
	}

	// Reset the circuit
	fmt.Println("\nResetting the circuit...")
	cb.Reset()
	fmt.Printf("Circuit state after reset: %s\n", cb.GetState())

	// Execute a successful function after reset
	fmt.Println("\nExecuting a successful function after reset...")
	result, err = circuit.Execute(ctx, cb, "example-operation", fn)
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Printf("Operation succeeded with result: %s\n", result)
	}
}