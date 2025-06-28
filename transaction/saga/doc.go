// Copyright (c) 2025 A Bit of Help, Inc.

// Package saga provides utilities for implementing the saga pattern for distributed transactions.
//
// The saga pattern is a way to manage distributed transactions without using two-phase commit.
// It works by defining a sequence of local transactions, each with a corresponding compensating
// transaction that can undo its effects. If any local transaction fails, the saga executes
// the compensating transactions in reverse order to maintain data consistency.
//
// This package implements the saga pattern with the following features:
//   - Transaction management with automatic rollback on failure
//   - Support for defining operations and their corresponding rollback operations
//   - Error handling with detailed information about failures
//   - Logging of transaction events
//   - Context support for cancellation and timeouts
//
// Key components:
//   - Transaction: The main struct that manages a sequence of operations and their rollbacks
//   - Operation: A function that performs a local transaction
//   - RollbackOperation: A function that rolls back a local transaction
//   - WithTransaction: A helper function for executing a function within a transaction
//
// Example usage:
//
//	err := saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
//	    // Add operations and their rollbacks to the transaction
//	    tx.AddOperation(
//	        // Operation to create a user
//	        func(ctx context.Context) error {
//	            return userRepo.Create(ctx, user)
//	        },
//	        // Rollback operation to delete the user if a later operation fails
//	        func(ctx context.Context) error {
//	            return userRepo.Delete(ctx, user.ID)
//	        },
//	    )
//
//	    tx.AddOperation(
//	        // Operation to create an account for the user
//	        func(ctx context.Context) error {
//	            return accountRepo.Create(ctx, account)
//	        },
//	        // Rollback operation to delete the account if a later operation fails
//	        func(ctx context.Context) error {
//	            return accountRepo.Delete(ctx, account.ID)
//	        },
//	    )
//
//	    return nil
//	})
//
// If any operation fails, the transaction will automatically roll back all previously
// executed operations in reverse order to maintain data consistency.
package saga