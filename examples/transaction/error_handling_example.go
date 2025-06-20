// Copyright (c) 2025 A Bit of Help, Inc.

// Example of error handling with checked rollbacks
package main

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
        saga.CheckedRollback(func(ctx context.Context) error {
            fmt.Println("Rolling back account balance update...")
            // Simulate a rollback error
            return errors.New("failed to rollback account balance update")
        }),
    )
    
    tx.AddOperation(
        // Operation to create an order that will fail
        func(ctx context.Context) error {
            fmt.Println("Creating order...")
            return errors.New("order creation failed")
        },
        // Rollback operation that will not be called because the operation fails
        func(ctx context.Context) error {
            fmt.Println("This rollback should not be called")
            return nil
        },
    )
    
    // Execute the transaction
    err = tx.Execute(ctx)
    if err != nil {
        fmt.Printf("Transaction failed as expected: %v\n", err)
        
        // Check if it's a rollback error
        var rollbackErr *saga.RollbackError
        if errors.As(err, &rollbackErr) {
            fmt.Printf("Rollback error details: %v\n", rollbackErr.Error())
            fmt.Printf("Original error: %v\n", rollbackErr.OriginalError())
            fmt.Printf("Rollback errors: %v\n", rollbackErr.RollbackErrors())
        }
    } else {
        fmt.Println("Transaction completed successfully (unexpected)")
    }
    
    // Expected output:
    // Updating account balance...
    // Creating order...
    // Rolling back account balance update...
    // Transaction failed as expected: transaction execution failed: transaction operation failed: order creation failed
    // Rollback error details: transaction rollback failed: failed to rollback account balance update
    // Original error: order creation failed
    // Rollback errors: [failed to rollback account balance update]
}