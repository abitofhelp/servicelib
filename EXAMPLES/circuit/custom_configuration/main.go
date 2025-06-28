// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating custom configuration of the circuit package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/circuit"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a logger
	logger, _ := zap.NewDevelopment()
	contextLogger := logging.NewContextLogger(logger)

	// Create a circuit breaker with custom configuration
	cfg := circuit.DefaultConfig().
		WithEnabled(true).                // Explicitly enable the circuit breaker
		WithTimeout(2 * time.Second).     // Set timeout to 2 seconds
		WithMaxConcurrent(50).            // Set max concurrent requests to 50
		WithErrorThreshold(0.3).          // Trip the circuit at 30% error rate
		WithVolumeThreshold(5).           // Require at least 5 requests before tripping
		WithSleepWindow(5 * time.Second)  // Wait 5 seconds before trying again

	options := circuit.DefaultOptions().
		WithName("custom-example").       // Set a custom name
		WithLogger(contextLogger)         // Use a custom logger

	cb := circuit.NewCircuitBreaker(cfg, options)

	// Define a function to execute with circuit breaking
	fn := func(ctx context.Context) (string, error) {
		// Your operation here
		// For demonstration, we'll simulate a successful operation
		fmt.Println("Executing operation with custom configuration...")
		return "success", nil
	}

	// Execute the function with circuit breaking
	fmt.Println("Executing function with custom circuit breaking...")
	result, err := circuit.Execute(ctx, cb, "custom-operation", fn)
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
	fmt.Println("With error threshold of 0.3 and volume threshold of 5,")
	fmt.Println("the circuit should trip after 2 failures (30% of 5 requests).")

	for i := 0; i < 5; i++ {
		result, err = circuit.Execute(ctx, cb, "failing-operation", failFn)
		if err != nil {
			fmt.Printf("Attempt %d: Operation failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Attempt %d: Operation succeeded with result: %s\n", i+1, result)
		}
	}

	// Check the circuit state
	fmt.Printf("\nCircuit state: %s\n", cb.GetState())

	// Demonstrate fallback functionality
	fmt.Println("\nDemonstrating fallback functionality...")
	
	// Define a fallback function
	fallbackFn := func(ctx context.Context, err error) (string, error) {
		fmt.Printf("Fallback called with error: %v\n", err)
		return "fallback result", nil
	}
	
	// Execute with fallback
	result, err = circuit.ExecuteWithFallback(ctx, cb, "failing-operation", failFn, fallbackFn)
	if err != nil {
		fmt.Printf("Operation with fallback failed: %v\n", err)
	} else {
		fmt.Printf("Operation with fallback succeeded with result: %s\n", result)
	}

	// Demonstrate individual parameter configuration
	fmt.Println("\nDemonstrating individual parameter configuration:")
	
	// Start with default config
	cfg = circuit.DefaultConfig()
	
	// Configure timeout only
	cfg = cfg.WithTimeout(3 * time.Second)
	fmt.Printf("Timeout: %v\n", cfg.Timeout)
	
	// Configure error threshold only
	cfg = cfg.WithErrorThreshold(0.2)
	fmt.Printf("Error threshold: %.1f\n", cfg.ErrorThreshold)
	
	// Configure volume threshold only
	cfg = cfg.WithVolumeThreshold(10)
	fmt.Printf("Volume threshold: %d\n", cfg.VolumeThreshold)
	
	// Configure sleep window only
	cfg = cfg.WithSleepWindow(10 * time.Second)
	fmt.Printf("Sleep window: %v\n", cfg.SleepWindow)
}