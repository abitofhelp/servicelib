//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of error handling with the context package
package example_context

import (
	"context" // Standard context package
	"fmt"
	"time"

	svcctx "github.com/abitofhelp/servicelib/context" // Aliased to avoid name collision
	"github.com/abitofhelp/servicelib/errors"
)

// errorHandlingMain is the main function for this example.
// It's not named "main" to avoid conflicts with other example files in the same package.
// To run this example, rename this function to "main" and comment out the main function in other example files.
func errorHandlingMain() {
	fmt.Println("Context Error Handling Examples")
	fmt.Println("==============================")

	// 1. Handling context timeout
	fmt.Println("\n1. Handling Context Timeout:")
	handleTimeout()

	// 2. Handling context cancellation
	fmt.Println("\n2. Handling Context Cancellation:")
	handleCancellation()

	// 3. Using CheckContext
	fmt.Println("\n3. Using CheckContext:")
	useCheckContext()

	// 4. Using MustCheck (with panic recovery)
	fmt.Println("\n4. Using MustCheck (with panic recovery):")
	useMustCheck()

	// 5. Error context with operation and service name
	fmt.Println("\n5. Error Context with Operation and Service Name:")
	errorWithContext()
}

// handleTimeout demonstrates handling context timeout
func handleTimeout() {
	// Create a context with a short timeout
	ctx, cancel := svcctx.WithTimeout(svcctx.Background(), 100*time.Millisecond, svcctx.ContextOptions{
		Operation: "database-query",
	})
	defer cancel()

	// Simulate a long-running operation
	err := simulateLongOperation(ctx, 200*time.Millisecond)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Check if it's a timeout error
		if errors.Is(err, errors.ErrTimeout) {
			fmt.Println("Detected timeout error")
		}
	}
}

// handleCancellation demonstrates handling context cancellation
func handleCancellation() {
	// Create a context
	ctx, cancel := svcctx.NewContext(svcctx.ContextOptions{
		Operation: "user-request",
	})

	// Simulate a cancellation after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Cancelling context...")
		cancel()
	}()

	// Simulate a long-running operation
	err := simulateLongOperation(ctx, 200*time.Millisecond)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Check if it's a cancellation error
		if errors.Is(err, errors.ErrCanceled) {
			fmt.Println("Detected cancellation error")
		}
	}
}

// useCheckContext demonstrates using the CheckContext function
func useCheckContext() {
	// Create contexts with different states
	activeCtx, activeCancel := svcctx.NewContext(svcctx.ContextOptions{})
	defer activeCancel()

	timeoutCtx, timeoutCancel := svcctx.WithTimeout(svcctx.Background(), 10*time.Millisecond, svcctx.ContextOptions{
		Operation: "timeout-operation",
	})
	defer timeoutCancel()

	cancelledCtx, cancelledCancel := svcctx.NewContext(svcctx.ContextOptions{
		Operation: "cancelled-operation",
	})
	cancelledCancel() // Cancel immediately

	// Wait for timeout to occur
	time.Sleep(20 * time.Millisecond)

	// Check active context
	if err := svcctx.CheckContext(activeCtx); err != nil {
		fmt.Printf("Active context error: %v\n", err)
	} else {
		fmt.Println("Active context is valid")
	}

	// Check timeout context
	if err := svcctx.CheckContext(timeoutCtx); err != nil {
		fmt.Printf("Timeout context error: %v\n", err)
	} else {
		fmt.Println("Timeout context is valid (unexpected)")
	}

	// Check cancelled context
	if err := svcctx.CheckContext(cancelledCtx); err != nil {
		fmt.Printf("Cancelled context error: %v\n", err)
	} else {
		fmt.Println("Cancelled context is valid (unexpected)")
	}
}

// useMustCheck demonstrates using the MustCheck function with panic recovery
func useMustCheck() {
	// Create a context with a short timeout
	ctx, cancel := svcctx.WithTimeout(svcctx.Background(), 10*time.Millisecond, svcctx.ContextOptions{
		Operation: "critical-operation",
	})
	defer cancel()

	// Wait for timeout to occur
	time.Sleep(20 * time.Millisecond)

	// Use MustCheck with panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	fmt.Println("Calling MustCheck on expired context...")
	svcctx.MustCheck(ctx)
	fmt.Println("This line will not be executed")
}

// errorWithContext demonstrates error context with operation and service name
func errorWithContext() {
	// Create a context with operation and service name
	ctx, cancel := svcctx.WithTimeout(svcctx.Background(), 10*time.Millisecond, svcctx.ContextOptions{
		Operation:   "data-processing",
		ServiceName: "analytics-service",
	})
	defer cancel()

	// Wait for timeout to occur
	time.Sleep(20 * time.Millisecond)

	// Check context and print error
	if err := svcctx.CheckContext(ctx); err != nil {
		fmt.Printf("Error with context: %v\n", err)

		// The error message will include the operation and service name
	}
}

// simulateLongOperation simulates a long-running operation that respects context
func simulateLongOperation(ctx context.Context, duration time.Duration) error {
	done := make(chan struct{})

	// Start the operation in a goroutine
	go func() {
		// Simulate work
		time.Sleep(duration)
		close(done)
	}()

	// Wait for either the operation to complete or the context to be done
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return svcctx.CheckContext(ctx)
	}
}
