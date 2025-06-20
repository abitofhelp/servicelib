// Copyright (c) 2025 A Bit of Help, Inc.

// Package saga provides utilities for implementing the saga pattern for distributed transactions.
package saga

import (
	"context"
	"go.uber.org/zap"

	appctx "github.com/abitofhelp/servicelib/context"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
)

// Operation represents a function that performs a database operation.
type Operation func(ctx context.Context) error

// RollbackOperation represents a function that rolls back a database operation.
type RollbackOperation func(ctx context.Context) error

// Transaction represents a transaction that can be executed and rolled back.
type Transaction struct {
	operations []Operation
	rollbacks  []RollbackOperation
	logger     *logging.ContextLogger
}

// NewTransaction creates a new transaction.
//
// Parameters:
//   - logger: The logger to use for logging transaction events
//
// Returns:
//   - *Transaction: A new transaction instance
func NewTransaction(logger *zap.Logger) *Transaction {
	var contextLogger *logging.ContextLogger
	if logger == nil {
		// Create a default logger with info level in development mode
		defaultLogger, _ := logging.NewLogger("info", true)
		contextLogger = logging.NewContextLogger(defaultLogger)
	} else {
		contextLogger = logging.NewContextLogger(logger)
	}
	return &Transaction{
		operations: make([]Operation, 0),
		rollbacks:  make([]RollbackOperation, 0),
		logger:     contextLogger,
	}
}

// AddOperation adds an operation to the transaction with its corresponding rollback operation.
//
// Parameters:
//   - op: The operation to execute
//   - rollback: The rollback operation to execute if the transaction fails
func (t *Transaction) AddOperation(op Operation, rollback RollbackOperation) {
	t.operations = append(t.operations, op)
	t.rollbacks = append(t.rollbacks, rollback)
}

// Execute executes all operations in the transaction. If any operation fails,
// it rolls back all previously executed operations in reverse order.
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - error: An error if the transaction fails
func (t *Transaction) Execute(ctx context.Context) error {
	// Check if context is done
	if err := appctx.CheckContext(ctx); err != nil {
		return errors.WrapWithOperation(err, "Transaction.Execute", "transaction execution aborted")
	}

	// Validate that operations and rollbacks have the same length
	if len(t.operations) != len(t.rollbacks) {
		return errors.Internal(nil, "transaction operations and rollbacks count mismatch")
	}

	for i, op := range t.operations {
		// Check context before each operation
		if err := appctx.CheckContext(ctx); err != nil {
			rollbackErrors := t.rollback(ctx, i-1)

			// Add rollback errors as details to the context error
			if len(rollbackErrors) > 0 {
				details := map[string]interface{}{
					"operation_index": i,
					"rollback_errors": rollbackErrors,
				}
				return errors.WithDetails(
					errors.WrapWithOperation(err, "Transaction.Execute", "transaction operation aborted"),
					details,
				)
			}

			return errors.WrapWithOperation(err, "Transaction.Execute", "transaction operation aborted")
		}

		if err := op(ctx); err != nil {
			// Roll back all previously executed operations in reverse order
			rollbackErrors := t.rollback(ctx, i-1)

			// Add operation index and rollback errors as details
			details := map[string]interface{}{
				"operation_index":  i,
				"failed_operation": i,
			}

			if len(rollbackErrors) > 0 {
				details["rollback_errors"] = rollbackErrors
			}

			return errors.WithDetails(
				errors.WrapWithOperation(err, "Transaction.Execute", "transaction operation failed"),
				details,
			)
		}
	}
	return nil
}

// rollback rolls back operations from index i down to 0 and returns any errors.
func (t *Transaction) rollback(ctx context.Context, i int) []error {
	var errors []error
	for j := i; j >= 0; j-- {
		if err := t.rollbacks[j](ctx); err != nil {
			t.logger.Error(ctx, "Failed to rollback operation", zap.Error(err))
			errors = append(errors, err)
		}
	}
	return errors
}

// WithTransaction executes a function within a transaction and handles rollback if needed.
//
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use for logging transaction events
//   - fn: The function to execute within the transaction
//
// Returns:
//   - error: An error if the transaction fails
func WithTransaction(ctx context.Context, logger *zap.Logger, fn func(tx *Transaction) error) error {
	// Check if context is done
	if err := appctx.CheckContext(ctx); err != nil {
		return errors.WrapWithOperation(err, "WithTransaction", "transaction creation aborted")
	}

	tx := NewTransaction(logger)
	if err := fn(tx); err != nil {
		// Add operation information to the error
		return errors.WrapWithOperation(err, "WithTransaction", "transaction setup failed")
	}

	// Execute the transaction
	if err := tx.Execute(ctx); err != nil {
		return errors.WrapWithOperation(err, "WithTransaction", "transaction execution failed")
	}

	return nil
}

// NoopRollback returns a rollback operation that does nothing.
//
// Returns:
//   - RollbackOperation: A rollback operation that does nothing
func NoopRollback() RollbackOperation {
	return func(ctx context.Context) error {
		return nil
	}
}

// CheckedRollback wraps a rollback operation with error checking.
//
// Parameters:
//   - rollback: The rollback operation to wrap
//   - operation: The name of the operation
//   - errorMsg: The error message to use if the rollback fails
//
// Returns:
//   - RollbackOperation: A wrapped rollback operation
func CheckedRollback(rollback RollbackOperation, operation string, errorMsg string) RollbackOperation {
	return func(ctx context.Context) error {
		if err := rollback(ctx); err != nil {
			return errors.WrapWithOperation(err, operation, errorMsg)
		}
		return nil
	}
}

// CheckedRollbackWithDetails wraps a rollback operation with error checking and additional details.
//
// Parameters:
//   - rollback: The rollback operation to wrap
//   - operation: The name of the operation
//   - errorMsg: The error message to use if the rollback fails
//   - details: Additional details to add to the error
//
// Returns:
//   - RollbackOperation: A wrapped rollback operation
func CheckedRollbackWithDetails(rollback RollbackOperation, operation string, errorMsg string, details map[string]interface{}) RollbackOperation {
	return func(ctx context.Context) error {
		if err := rollback(ctx); err != nil {
			wrappedErr := errors.WrapWithOperation(err, operation, errorMsg)
			return errors.WithDetails(wrappedErr, details)
		}
		return nil
	}
}
