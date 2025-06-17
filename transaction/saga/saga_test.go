// Copyright (c) 2025 A Bit of Help, Inc.

package saga

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "github.com/abitofhelp/servicelib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// TestNewTransaction_WithLogger tests creating a new transaction with a provided logger
func TestNewTransaction_WithLogger(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a new transaction
	tx := NewTransaction(testLogger)

	// Assert that the transaction was created correctly
	assert.NotNil(t, tx)
	assert.NotNil(t, tx.logger)
	assert.Empty(t, tx.operations)
	assert.Empty(t, tx.rollbacks)
}

// TestNewTransaction_WithoutLogger tests creating a new transaction without a logger
func TestNewTransaction_WithoutLogger(t *testing.T) {
	// Create a new transaction with nil logger
	tx := NewTransaction(nil)

	// Assert that the transaction was created correctly with a default logger
	assert.NotNil(t, tx)
	assert.NotNil(t, tx.logger)
	assert.Empty(t, tx.operations)
	assert.Empty(t, tx.rollbacks)
}

// TestAddOperation tests adding operations to a transaction
func TestAddOperation(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a new transaction
	tx := NewTransaction(testLogger)

	// Create test operations and rollbacks
	op1 := func(ctx context.Context) error { return nil }
	rollback1 := func(ctx context.Context) error { return nil }

	op2 := func(ctx context.Context) error { return nil }
	rollback2 := func(ctx context.Context) error { return nil }

	// Add operations to the transaction
	tx.AddOperation(op1, rollback1)
	tx.AddOperation(op2, rollback2)

	// Assert that the operations and rollbacks were added correctly
	assert.Len(t, tx.operations, 2)
	assert.Len(t, tx.rollbacks, 2)
}

// TestExecute_Success tests executing a transaction successfully
func TestExecute_Success(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a new transaction
	tx := NewTransaction(testLogger)

	// Create test operations and rollbacks with tracking variables
	op1Called := false
	op1 := func(ctx context.Context) error {
		op1Called = true
		return nil
	}

	rollback1Called := false
	rollback1 := func(ctx context.Context) error {
		rollback1Called = true
		return nil
	}

	op2Called := false
	op2 := func(ctx context.Context) error {
		op2Called = true
		return nil
	}

	rollback2Called := false
	rollback2 := func(ctx context.Context) error {
		rollback2Called = true
		return nil
	}

	// Add operations to the transaction
	tx.AddOperation(op1, rollback1)
	tx.AddOperation(op2, rollback2)

	// Execute the transaction
	err := tx.Execute(context.Background())

	// Assert that the transaction executed successfully
	assert.NoError(t, err)
	assert.True(t, op1Called, "First operation should have been called")
	assert.True(t, op2Called, "Second operation should have been called")
	assert.False(t, rollback1Called, "First rollback should not have been called")
	assert.False(t, rollback2Called, "Second rollback should not have been called")
}

// TestExecute_OperationError tests executing a transaction with an operation that fails
func TestExecute_OperationError(t *testing.T) {
	// Create a test logger
	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	// Create a new transaction
	tx := NewTransaction(observedLogger)

	// Create test operations and rollbacks with tracking variables
	op1Called := false
	op1 := func(ctx context.Context) error {
		op1Called = true
		return nil
	}

	rollback1Called := false
	rollback1 := func(ctx context.Context) error {
		rollback1Called = true
		return nil
	}

	op2Called := false
	expectedErr := errors.New("operation 2 failed")
	op2 := func(ctx context.Context) error {
		op2Called = true
		return expectedErr
	}

	rollback2Called := false
	rollback2 := func(ctx context.Context) error {
		rollback2Called = true
		return nil
	}

	// Add operations to the transaction
	tx.AddOperation(op1, rollback1)
	tx.AddOperation(op2, rollback2)

	// Execute the transaction
	err := tx.Execute(context.Background())

	// Assert that the transaction failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation 2 failed")

	// Assert that operations and rollbacks were called correctly
	assert.True(t, op1Called, "First operation should have been called")
	assert.True(t, op2Called, "Second operation should have been called")
	assert.True(t, rollback1Called, "First rollback should have been called")
	assert.False(t, rollback2Called, "Second rollback should not have been called")

	// Check that the error contains the expected details
	var contextualErr *apperrors.ContextualError
	assert.True(t, apperrors.As(err, &contextualErr))
	assert.NotNil(t, contextualErr.Context.Details)
	assert.Equal(t, 1, contextualErr.Context.Details["operation_index"])
	assert.Equal(t, 1, contextualErr.Context.Details["failed_operation"])
}

// TestExecute_RollbackError tests executing a transaction with a rollback that fails
func TestExecute_RollbackError(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	// Create a new transaction
	tx := NewTransaction(observedLogger)

	// Create test operations and rollbacks with tracking variables
	op1Called := false
	op1 := func(ctx context.Context) error {
		op1Called = true
		return nil
	}

	rollback1Called := false
	rollbackErr := errors.New("rollback 1 failed")
	rollback1 := func(ctx context.Context) error {
		rollback1Called = true
		return rollbackErr
	}

	op2Called := false
	expectedErr := errors.New("operation 2 failed")
	op2 := func(ctx context.Context) error {
		op2Called = true
		return expectedErr
	}

	rollback2Called := false
	rollback2 := func(ctx context.Context) error {
		rollback2Called = true
		return nil
	}

	// Add operations to the transaction
	tx.AddOperation(op1, rollback1)
	tx.AddOperation(op2, rollback2)

	// Execute the transaction
	err := tx.Execute(context.Background())

	// Assert that the transaction failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation 2 failed")

	// Assert that operations and rollbacks were called correctly
	assert.True(t, op1Called, "First operation should have been called")
	assert.True(t, op2Called, "Second operation should have been called")
	assert.True(t, rollback1Called, "First rollback should have been called")
	assert.False(t, rollback2Called, "Second rollback should not have been called")

	// Check that the error contains the expected details
	var contextualErr *apperrors.ContextualError
	assert.True(t, apperrors.As(err, &contextualErr))
	assert.NotNil(t, contextualErr.Context.Details)
	assert.Equal(t, 1, contextualErr.Context.Details["operation_index"])
	assert.Equal(t, 1, contextualErr.Context.Details["failed_operation"])

	// Check that rollback errors are included in the details
	rollbackErrors, ok := contextualErr.Context.Details["rollback_errors"].([]error)
	assert.True(t, ok)
	assert.Len(t, rollbackErrors, 1)
	assert.Equal(t, rollbackErr.Error(), rollbackErrors[0].Error())

	// Verify logs
	logs := observedLogs.All()
	var foundRollbackError bool
	for _, log := range logs {
		if log.Message == "Failed to rollback operation" {
			foundRollbackError = true
			break
		}
	}
	assert.True(t, foundRollbackError, "Expected 'Failed to rollback operation' log message")
}

// TestExecute_ContextCancellation tests executing a transaction with a cancelled context
func TestExecute_ContextCancellation(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a new transaction
	tx := NewTransaction(testLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Create test operations and rollbacks
	opCalled := false
	op := func(ctx context.Context) error {
		opCalled = true
		return nil
	}

	rollbackCalled := false
	rollback := func(ctx context.Context) error {
		rollbackCalled = true
		return nil
	}

	// Add operations to the transaction
	tx.AddOperation(op, rollback)

	// Execute the transaction with the cancelled context
	err := tx.Execute(ctx)

	// Assert that the transaction failed with a context error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation cancelled")
	assert.False(t, opCalled, "Operation should not have been called")
	assert.False(t, rollbackCalled, "Rollback should not have been called")
}

// TestExecute_ContextCancellationDuringOperation tests context cancellation during operation execution
func TestExecute_ContextCancellationDuringOperation(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a new transaction
	tx := NewTransaction(testLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create test operations and rollbacks with tracking variables
	op1Called := false
	op1 := func(ctx context.Context) error {
		op1Called = true
		return nil
	}

	rollback1Called := false
	rollback1 := func(ctx context.Context) error {
		rollback1Called = true
		return nil
	}

	op2Called := false
	op2 := func(ctx context.Context) error {
		op2Called = true
		// Cancel the context during the second operation
		cancel()
		// Sleep briefly to ensure the context cancellation is processed
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	rollback2Called := false
	rollback2 := func(ctx context.Context) error {
		rollback2Called = true
		return nil
	}

	op3Called := false
	op3 := func(ctx context.Context) error {
		op3Called = true
		return nil
	}

	rollback3Called := false
	rollback3 := func(ctx context.Context) error {
		rollback3Called = true
		return nil
	}

	// Add operations to the transaction
	tx.AddOperation(op1, rollback1)
	tx.AddOperation(op2, rollback2)
	tx.AddOperation(op3, rollback3)

	// Execute the transaction
	err := tx.Execute(ctx)

	// Assert that the transaction failed with a context error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation cancelled")

	// Assert that operations and rollbacks were called correctly
	assert.True(t, op1Called, "First operation should have been called")
	assert.True(t, op2Called, "Second operation should have been called")
	assert.False(t, op3Called, "Third operation should not have been called")
	assert.True(t, rollback1Called, "First rollback should have been called")
	assert.True(t, rollback2Called, "Second rollback should have been called")
	assert.False(t, rollback3Called, "Third rollback should not have been called")
}

// TestExecute_OperationsMismatch tests executing a transaction with mismatched operations and rollbacks
func TestExecute_OperationsMismatch(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a new transaction
	tx := NewTransaction(testLogger)

	// Add an operation without a rollback (manually to bypass AddOperation)
	tx.operations = append(tx.operations, func(ctx context.Context) error { return nil })

	// Execute the transaction
	err := tx.Execute(context.Background())

	// Assert that the transaction failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction operations and rollbacks count mismatch")
}

// TestWithTransaction_Success tests WithTransaction with a successful transaction
func TestWithTransaction_Success(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a context
	ctx := context.Background()

	// Create tracking variables
	setupCalled := false
	op1Called := false
	op2Called := false

	// Call WithTransaction
	err := WithTransaction(ctx, testLogger, func(tx *Transaction) error {
		setupCalled = true

		// Add operations to the transaction
		tx.AddOperation(
			func(ctx context.Context) error {
				op1Called = true
				return nil
			},
			func(ctx context.Context) error { return nil },
		)

		tx.AddOperation(
			func(ctx context.Context) error {
				op2Called = true
				return nil
			},
			func(ctx context.Context) error { return nil },
		)

		return nil
	})

	// Assert that the transaction executed successfully
	assert.NoError(t, err)
	assert.True(t, setupCalled, "Setup function should have been called")
	assert.True(t, op1Called, "First operation should have been called")
	assert.True(t, op2Called, "Second operation should have been called")
}

// TestWithTransaction_SetupError tests WithTransaction with a setup error
func TestWithTransaction_SetupError(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a context
	ctx := context.Background()

	// Expected error
	expectedErr := errors.New("setup error")

	// Call WithTransaction with a setup error
	err := WithTransaction(ctx, testLogger, func(tx *Transaction) error {
		return expectedErr
	})

	// Assert that the transaction failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "setup error")
}

// TestWithTransaction_ExecuteError tests WithTransaction with an execution error
func TestWithTransaction_ExecuteError(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a context
	ctx := context.Background()

	// Expected error
	expectedErr := errors.New("operation error")

	// Call WithTransaction with an operation that fails
	err := WithTransaction(ctx, testLogger, func(tx *Transaction) error {
		tx.AddOperation(
			func(ctx context.Context) error {
				return expectedErr
			},
			func(ctx context.Context) error { return nil },
		)
		return nil
	})

	// Assert that the transaction failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation error")
}

// TestWithTransaction_ContextCancellation tests WithTransaction with a cancelled context
func TestWithTransaction_ContextCancellation(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Call WithTransaction with a cancelled context
	err := WithTransaction(ctx, testLogger, func(tx *Transaction) error {
		// This should not be executed
		require.Fail(t, "Setup function should not be called with cancelled context")
		return nil
	})

	// Assert that the transaction failed with a context error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation cancelled")
}

// TestNoopRollback tests the NoopRollback function
func TestNoopRollback(t *testing.T) {
	// Get the noop rollback function
	rollback := NoopRollback()

	// Execute the rollback
	err := rollback(context.Background())

	// Assert that the rollback succeeded
	assert.NoError(t, err)
}

// TestCheckedRollback_Success tests CheckedRollback with a successful rollback
func TestCheckedRollback_Success(t *testing.T) {
	// Create a rollback that succeeds
	baseRollback := func(ctx context.Context) error {
		return nil
	}

	// Wrap the rollback with CheckedRollback
	rollback := CheckedRollback(baseRollback, "TestOperation", "Rollback failed")

	// Execute the rollback
	err := rollback(context.Background())

	// Assert that the rollback succeeded
	assert.NoError(t, err)
}

// TestCheckedRollback_Error tests CheckedRollback with a failing rollback
func TestCheckedRollback_Error(t *testing.T) {
	// Expected error
	expectedErr := errors.New("rollback error")

	// Create a rollback that fails
	baseRollback := func(ctx context.Context) error {
		return expectedErr
	}

	// Wrap the rollback with CheckedRollback
	rollback := CheckedRollback(baseRollback, "TestOperation", "Rollback failed")

	// Execute the rollback
	err := rollback(context.Background())

	// Assert that the rollback failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rollback error")
	assert.Contains(t, err.Error(), "operation TestOperation")

	// Check that the error is wrapped correctly
	var contextualErr *apperrors.ContextualError
	assert.True(t, apperrors.As(err, &contextualErr))
	assert.Equal(t, "TestOperation", contextualErr.Context.Operation)
}

// TestCheckedRollbackWithDetails_Success tests CheckedRollbackWithDetails with a successful rollback
func TestCheckedRollbackWithDetails_Success(t *testing.T) {
	// Create a rollback that succeeds
	baseRollback := func(ctx context.Context) error {
		return nil
	}

	// Create details
	details := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	// Wrap the rollback with CheckedRollbackWithDetails
	rollback := CheckedRollbackWithDetails(baseRollback, "TestOperation", "Rollback failed", details)

	// Execute the rollback
	err := rollback(context.Background())

	// Assert that the rollback succeeded
	assert.NoError(t, err)
}

// TestCheckedRollbackWithDetails_Error tests CheckedRollbackWithDetails with a failing rollback
func TestCheckedRollbackWithDetails_Error(t *testing.T) {
	// Expected error
	expectedErr := errors.New("rollback error")

	// Create a rollback that fails
	baseRollback := func(ctx context.Context) error {
		return expectedErr
	}

	// Create details
	details := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	// Wrap the rollback with CheckedRollbackWithDetails
	rollback := CheckedRollbackWithDetails(baseRollback, "TestOperation", "Rollback failed", details)

	// Execute the rollback
	err := rollback(context.Background())

	// Assert that the rollback failed with the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rollback error")
	assert.Contains(t, err.Error(), "operation TestOperation")

	// Check that the error is wrapped correctly
	var contextualErr *apperrors.ContextualError
	assert.True(t, apperrors.As(err, &contextualErr))
	assert.Equal(t, "TestOperation", contextualErr.Context.Operation)

	// Check that the details are included
	assert.Equal(t, "value1", contextualErr.Context.Details["key1"])
	assert.Equal(t, 42, contextualErr.Context.Details["key2"])
}
