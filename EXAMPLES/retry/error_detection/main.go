// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating custom error detection with the retry package
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	serviceErrors "github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/retry"
)

// simulateOperation simulates an operation that might fail with different types of errors
func simulateOperation(ctx context.Context, attempt int) error {
	fmt.Printf("Attempt %d: Executing operation...\n", attempt+1)

	// Simulate different errors based on the attempt number
	switch attempt {
	case 0:
		// First attempt: network error (retryable)
		return &net.OpError{
			Op:  "read",
			Net: "tcp",
			Err: fmt.Errorf("connection reset by peer"),
		}
	case 1:
		// Second attempt: timeout error (retryable)
		return context.DeadlineExceeded
	case 2:
		// Third attempt: transient error (retryable)
		return serviceErrors.NewDatabaseError("connection refused", "connect", "users", nil)
	case 3:
		// Fourth attempt: success
		return nil
	default:
		// Any further attempts: non-retryable error
		return errors.New("permanent error")
	}
}

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Track the current attempt
	attempt := 0

	// Define a function to retry that uses our simulateOperation function
	fn := func(ctx context.Context) error {
		err := simulateOperation(ctx, attempt)
		attempt++
		return err
	}

	// Define a custom function to determine if an error is retryable
	isRetryable := func(err error) bool {
		// Check if the error is a network error
		if serviceErrors.IsNetworkError(err) {
			fmt.Println("Detected network error (retryable)")
			return true
		}

		// Check if the error is a timeout error
		if serviceErrors.IsTimeout(err) {
			fmt.Println("Detected timeout error (retryable)")
			return true
		}

		// Check if the error is a transient error
		if serviceErrors.IsTransientError(err) {
			fmt.Println("Detected transient error (retryable)")
			return true
		}

		// Check for specific error types
		if errors.Is(err, io.EOF) {
			fmt.Println("Detected EOF error (retryable)")
			return true
		}

		// Add your custom logic here
		// For example, check for specific error messages or types
		if err != nil && err.Error() == "connection refused" {
			fmt.Println("Detected connection refused error (retryable)")
			return true
		}

		fmt.Println("Error is not retryable:", err)
		return false
	}

	// Use a configuration with more retries to demonstrate multiple error types
	config := retry.DefaultConfig().
		WithMaxRetries(5).
		WithInitialBackoff(100 * time.Millisecond)

	// Execute with retry and custom error detection
	err := retry.Do(ctx, fn, config, isRetryable)
	if err != nil {
		fmt.Printf("\nOperation failed after retries: %v\n", err)
	} else {
		fmt.Println("\nOperation succeeded!")
	}
}