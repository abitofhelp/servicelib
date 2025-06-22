// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic saga transaction
package example_transaction

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

	// Execute operations within a transaction
	err = saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
		// Add operations with their corresponding rollback operations
		tx.AddOperation(
			// Operation to create a user
			func(ctx context.Context) error {
				fmt.Println("Creating user...")
				// In a real application, this would call a user service
				return nil
			},
			// Rollback operation to delete the user if a later operation fails
			func(ctx context.Context) error {
				fmt.Println("Rolling back user creation...")
				// In a real application, this would call the user service to delete the user
				return nil
			},
		)

		tx.AddOperation(
			// Operation to create an order
			func(ctx context.Context) error {
				fmt.Println("Creating order...")
				// In a real application, this would call an order service
				return nil
			},
			// Rollback operation to delete the order if a later operation fails
			func(ctx context.Context) error {
				fmt.Println("Rolling back order creation...")
				// In a real application, this would call the order service to delete the order
				return nil
			},
		)

		tx.AddOperation(
			// Operation to process payment
			func(ctx context.Context) error {
				fmt.Println("Processing payment...")
				// Simulate a failure in payment processing
				return fmt.Errorf("payment processing failed")
			},
			// Rollback operation to refund the payment
			func(ctx context.Context) error {
				fmt.Println("Rolling back payment processing...")
				// In a real application, this would call the payment service to issue a refund
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
	// Creating user...
	// Creating order...
	// Processing payment...
	// Rolling back order creation...
	// Rolling back user creation...
	// Transaction failed: transaction execution failed: transaction operation failed: payment processing failed
}
