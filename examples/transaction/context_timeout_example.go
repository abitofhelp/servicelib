// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using context with timeout in transactions
package example_transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/transaction/saga"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute operations within a transaction
	err = saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
		// Add operations with their corresponding rollback operations
		tx.AddOperation(
			// Operation to process data
			func(ctx context.Context) error {
				fmt.Println("Processing data...")

				// Check if context is done (timeout or cancellation)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(1 * time.Second): // Simulate work
					fmt.Println("Data processing completed")
					return nil
				}
			},
			// Rollback operation
			func(ctx context.Context) error {
				fmt.Println("Rolling back data processing...")
				return nil
			},
		)

		tx.AddOperation(
			// Operation that takes longer than the timeout
			func(ctx context.Context) error {
				fmt.Println("Starting long operation...")

				// Check if context is done (timeout or cancellation)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(10 * time.Second): // This will exceed the timeout
					fmt.Println("Long operation completed (this won't be printed)")
					return nil
				}
			},
			// Rollback operation
			func(ctx context.Context) error {
				fmt.Println("Rolling back long operation...")
				return nil
			},
		)

		return nil
	})

	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}

	// Expected output:
	// Processing data...
	// Data processing completed
	// Starting long operation...
	// Rolling back data processing...
	// Transaction failed: transaction execution failed: transaction operation failed: context deadline exceeded
}
