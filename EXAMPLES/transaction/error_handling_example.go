// Copyright (c) 2025 A Bit of Help, Inc.

// Example of error handling with checked rollbacks
package example_transaction

import (
	"context"
	"errors"
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

	// Add operations with checked rollbacks
	tx.AddOperation(
		// Operation to update account balance
		func(ctx context.Context) error {
			fmt.Println("Updating account balance...")
			return nil
		},
		// Rollback operation with error checking
		saga.CheckedRollback(
			func(ctx context.Context) error {
				fmt.Println("Rolling back account balance update...")
				// Simulate a rollback error
				return errors.New("failed to rollback account balance update")
			},
			"update_balance",
			"Failed to rollback account balance update",
		),
	)

	tx.AddOperation(
		// Operation to create an order that will fail
		func(ctx context.Context) error {
			fmt.Println("Creating order...")
			return errors.New("order creation failed")
		},
		saga.CheckedRollback(
			func(ctx context.Context) error {
				fmt.Println("Rolling back order creation...")
				return nil
			},
			"create_order",
			"Failed to rollback order creation",
		),
	)

	// Execute the transaction
	err = tx.Execute(ctx)
	if err != nil {
		logger.Error("Transaction failed",
			zap.Error(err))
		return
	}

	logger.Info("Transaction completed successfully")
}
