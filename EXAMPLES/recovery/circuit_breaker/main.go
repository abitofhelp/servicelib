// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/errors/recovery"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

func main() {
	// Create a context
	ctx := context.Background()

	// Create a logger
	logger, _ := zap.NewDevelopment()
	contextLogger := logging.NewContextLogger(logger)

	// Create a recovery handler
	handler := recovery.NewRecoveryHandler(contextLogger, 10)

	// Create a circuit breaker
	cb := recovery.NewCircuitBreaker(2, 5*time.Second)

	// Simulate some operations
	for i := 0; i < 5; i++ {
		err := handler.WithRecovery(ctx, fmt.Sprintf("operation-%d", i), func() error {
			return cb.Execute(ctx, func() error {
				// Simulate a failing operation
				return fmt.Errorf("operation failed")
			})
		})

		if err != nil {
			fmt.Printf("Attempt %d: %v\n", i+1, err)
		}

		// Wait a bit between attempts
		time.Sleep(time.Second)
	}

	// Wait for circuit breaker to reset
	fmt.Println("Waiting for circuit breaker to reset...")
	time.Sleep(6 * time.Second)

	// Try again after reset
	err := handler.WithRecovery(ctx, "final-operation", func() error {
		return cb.Execute(ctx, func() error {
			return nil // Successful operation
		})
	})

	if err != nil {
		fmt.Printf("After reset: %v\n", err)
	} else {
		fmt.Println("Operation successful after reset")
	}

	// Example with context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel the context immediately

	err = handler.WithRecovery(cancelCtx, "canceled-operation", func() error {
		return cb.Execute(cancelCtx, func() error {
			return nil // This should not be executed
		})
	})

	fmt.Printf("Canceled operation: %v\n", err)

	// Example with context timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // Ensure deadline is exceeded

	err = handler.WithRecovery(timeoutCtx, "timeout-operation", func() error {
		return cb.Execute(timeoutCtx, func() error {
			return nil // This should not be executed
		})
	})

	fmt.Printf("Timeout operation: %v\n", err)
}