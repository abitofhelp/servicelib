// Copyright (c) 2025 A Bit of Help, Inc.

// Example of custom transaction execution
package main

import (
	"context"
	"fmt"

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

	// Create a context
	ctx := context.Background()

	// Create a transaction
	tx := saga.NewTransaction(logger)

	// Add operations with their corresponding rollback operations
	tx.AddOperation(
		// Operation to reserve inventory
		func(ctx context.Context) error {
			fmt.Println("Reserving inventory...")
			return nil
		},
		// Rollback operation to release inventory
		func(ctx context.Context) error {
			fmt.Println("Releasing inventory...")
			return nil
		},
	)

	tx.AddOperation(
		// Operation to create shipment
		func(ctx context.Context) error {
			fmt.Println("Creating shipment...")
			return nil
		},
		// Rollback operation to cancel shipment
		func(ctx context.Context) error {
			fmt.Println("Cancelling shipment...")
			return nil
		},
	)

	tx.AddOperation(
		// Operation to send notification
		func(ctx context.Context) error {
			fmt.Println("Sending notification...")
			return nil
		},
		// No-op rollback since notifications can't be "unsent"
		saga.NoopRollback(),
	)

	// Execute the transaction
	if err := tx.Execute(ctx); err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}

	// Expected output:
	// Reserving inventory...
	// Creating shipment...
	// Sending notification...
	// Transaction completed successfully
}
