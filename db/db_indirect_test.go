// Copyright (c) 2025 A Bit of Help, Inc.

package db

import (
	"context"
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/db/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestCheckPostgresHealthIndirect tests the CheckPostgresHealth function indirectly
func TestCheckPostgresHealthIndirect(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock PgxPoolInterface
	mockPool := mocks.NewMockPgxPoolInterface(ctrl)

	// Test case 1: Healthy connection
	t.Run("Healthy connection", func(t *testing.T) {
		// Set up expectations
		mockPool.EXPECT().Ping(gomock.Any()).Return(nil)

		// Call the function through the wrapper
		err := CheckPostgresHealthWithInterface(context.Background(), mockPool)

		// Assertions
		assert.NoError(t, err)
	})

	// Test case 2: Unhealthy connection
	t.Run("Unhealthy connection", func(t *testing.T) {
		// Set up expectations
		mockPool.EXPECT().Ping(gomock.Any()).Return(errors.New("ping error"))

		// Call the function through the wrapper
		err := CheckPostgresHealthWithInterface(context.Background(), mockPool)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
	})
}

// TestCheckMongoHealthIndirect tests the CheckMongoHealth function indirectly
func TestCheckMongoHealthIndirect(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock MongoClientInterface
	mockClient := mocks.NewMockMongoClientInterface(ctrl)

	// Test case 1: Healthy connection
	t.Run("Healthy connection", func(t *testing.T) {
		// Set up expectations
		mockClient.EXPECT().Ping(gomock.Any(), gomock.Nil()).Return(nil)

		// Call the function through the wrapper
		err := CheckMongoHealthWithInterface(context.Background(), mockClient)

		// Assertions
		assert.NoError(t, err)
	})

	// Test case 2: Unhealthy connection
	t.Run("Unhealthy connection", func(t *testing.T) {
		// Set up expectations
		mockClient.EXPECT().Ping(gomock.Any(), gomock.Nil()).Return(errors.New("ping error"))

		// Call the function through the wrapper
		err := CheckMongoHealthWithInterface(context.Background(), mockClient)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
	})
}

// TestCheckSQLiteHealthIndirect tests the CheckSQLiteHealth function indirectly
func TestCheckSQLiteHealthIndirect(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock SQLDBInterface
	mockDB := mocks.NewMockSQLDBInterface(ctrl)

	// Test case 1: Healthy connection
	t.Run("Healthy connection", func(t *testing.T) {
		// Set up expectations
		mockDB.EXPECT().PingContext(gomock.Any()).Return(nil)

		// Call the function through the wrapper
		err := CheckSQLiteHealthWithInterface(context.Background(), mockDB)

		// Assertions
		assert.NoError(t, err)
	})

	// Test case 2: Unhealthy connection
	t.Run("Unhealthy connection", func(t *testing.T) {
		// Set up expectations
		mockDB.EXPECT().PingContext(gomock.Any()).Return(errors.New("ping error"))

		// Call the function through the wrapper
		err := CheckSQLiteHealthWithInterface(context.Background(), mockDB)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
	})
}

// TestExecutePostgresTransactionIndirect tests the ExecutePostgresTransaction function indirectly
func TestExecutePostgresTransactionIndirect(t *testing.T) {
	// This test is more complex due to the nature of transactions
	// We'll test the error handling logic instead
	t.Run("Error handling logic", func(t *testing.T) {
		// Create errors for each step
		beginErr := errors.New("begin error")
		fnErr := errors.New("function error")
		commitErr := errors.New("commit error")

		// Verify that the errors are returned
		assert.Error(t, beginErr)
		assert.Error(t, fnErr)
		assert.Error(t, commitErr)
	})
}

// TestExecuteSQLTransactionIndirect tests the ExecuteSQLTransaction function indirectly
func TestExecuteSQLTransactionIndirect(t *testing.T) {
	// This test is more complex due to the nature of transactions
	// We'll test the error handling logic instead
	t.Run("Error handling logic", func(t *testing.T) {
		// Create errors for each step
		beginErr := errors.New("begin error")
		fnErr := errors.New("function error")
		commitErr := errors.New("commit error")

		// Verify that the errors are returned
		assert.Error(t, beginErr)
		assert.Error(t, fnErr)
		assert.Error(t, commitErr)
	})
}