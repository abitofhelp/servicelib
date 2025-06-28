//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating custom configuration of the retry package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/retry"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Define a function to retry
	fn := func(ctx context.Context) error {
		// Your operation here
		// For demonstration, we'll simulate a successful operation
		fmt.Println("Executing operation with custom configuration...")
		return nil
	}

	// Create a custom retry configuration
	config := retry.DefaultConfig().
		WithMaxRetries(5).                    // Increase max retries to 5
		WithInitialBackoff(200 * time.Millisecond). // Start with 200ms backoff
		WithMaxBackoff(5 * time.Second).      // Cap backoff at 5 seconds
		WithBackoffFactor(2.5).               // Use 2.5x backoff factor
		WithJitterFactor(0.3)                 // Add 30% jitter

	// Execute with retry and custom configuration
	err := retry.Do(ctx, fn, config, nil)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("Operation succeeded with custom configuration!")
	}

	// You can also configure individual parameters
	fmt.Println("\nDemonstrating individual parameter configuration:")
	
	// Start with default config
	config = retry.DefaultConfig()
	
	// Configure max retries only
	config = config.WithMaxRetries(10)
	fmt.Printf("Max retries: %d\n", config.MaxRetries)
	
	// Configure initial backoff only
	config = config.WithInitialBackoff(500 * time.Millisecond)
	fmt.Printf("Initial backoff: %v\n", config.InitialBackoff)
	
	// Configure max backoff only
	config = config.WithMaxBackoff(10 * time.Second)
	fmt.Printf("Max backoff: %v\n", config.MaxBackoff)
	
	// Configure backoff factor only
	config = config.WithBackoffFactor(3.0)
	fmt.Printf("Backoff factor: %.1f\n", config.BackoffFactor)
	
	// Configure jitter factor only
	config = config.WithJitterFactor(0.5)
	fmt.Printf("Jitter factor: %.1f\n", config.JitterFactor)
}