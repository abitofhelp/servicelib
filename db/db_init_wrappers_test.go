// Copyright (c) 2025 A Bit of Help, Inc.

package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/db/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// MockInitMongoClient is a mock implementation for testing MongoDB client initialization
func MockInitMongoClient(ctx context.Context, uri string, timeout time.Duration, mockClient *mocks.MockMongoClientInterface) (*mocks.MockMongoClientInterface, error) {
	// Simulate connection error if URI is invalid
	if uri == "invalid-uri" {
		return nil, errors.New("failed to parse URI")
	}

	// Simulate connection timeout if timeout is too short
	if timeout < time.Millisecond {
		return nil, errors.New("connection timeout")
	}

	// Return the mock client for successful initialization
	return mockClient, nil
}

// TestMockInitMongoClient tests the MockInitMongoClient function
func TestMockInitMongoClient(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock MongoClientInterface
	mockClient := mocks.NewMockMongoClientInterface(ctrl)

	// Test case 1: Successful initialization
	t.Run("Successful initialization", func(t *testing.T) {
		// Call the function
		client, err := MockInitMongoClient(context.Background(), "mongodb://localhost:27017", 1*time.Second, mockClient)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, mockClient, client)
	})

	// Test case 2: Invalid URI
	t.Run("Invalid URI", func(t *testing.T) {
		// Call the function
		client, err := MockInitMongoClient(context.Background(), "invalid-uri", 1*time.Second, mockClient)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to parse URI")
	})

	// Test case 3: Connection timeout
	t.Run("Connection timeout", func(t *testing.T) {
		// Call the function
		client, err := MockInitMongoClient(context.Background(), "mongodb://localhost:27017", 0, mockClient)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "connection timeout")
	})
}

// MockInitPostgresPool is a mock implementation for testing PostgreSQL pool initialization
func MockInitPostgresPool(ctx context.Context, uri string, timeout time.Duration, mockPool *mocks.MockPgxPoolInterface) (*mocks.MockPgxPoolInterface, error) {
	// Simulate connection error if URI is invalid
	if uri == "invalid-uri" {
		return nil, errors.New("failed to parse URI")
	}

	// Simulate connection timeout if timeout is too short
	if timeout < time.Millisecond {
		return nil, errors.New("connection timeout")
	}

	// Return the mock pool for successful initialization
	return mockPool, nil
}

// TestMockInitPostgresPool tests the MockInitPostgresPool function
func TestMockInitPostgresPool(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock PgxPoolInterface
	mockPool := mocks.NewMockPgxPoolInterface(ctrl)

	// Test case 1: Successful initialization
	t.Run("Successful initialization", func(t *testing.T) {
		// Call the function
		pool, err := MockInitPostgresPool(context.Background(), "postgres://localhost:5432", 1*time.Second, mockPool)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, mockPool, pool)
	})

	// Test case 2: Invalid URI
	t.Run("Invalid URI", func(t *testing.T) {
		// Call the function
		pool, err := MockInitPostgresPool(context.Background(), "invalid-uri", 1*time.Second, mockPool)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, pool)
		assert.Contains(t, err.Error(), "failed to parse URI")
	})

	// Test case 3: Connection timeout
	t.Run("Connection timeout", func(t *testing.T) {
		// Call the function
		pool, err := MockInitPostgresPool(context.Background(), "postgres://localhost:5432", 0, mockPool)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, pool)
		assert.Contains(t, err.Error(), "connection timeout")
	})
}

// MockInitSQLiteDB is a mock implementation for testing SQLite DB initialization
func MockInitSQLiteDB(ctx context.Context, uri string, timeout time.Duration, connMaxLifetime time.Duration, maxOpenConns int, maxIdleConns int, mockDB *mocks.MockSQLDBInterface) (*mocks.MockSQLDBInterface, error) {
	// Simulate connection error if URI is invalid
	if uri == "invalid-uri" {
		return nil, errors.New("failed to open database")
	}

	// Simulate connection timeout if timeout is too short
	if timeout < time.Millisecond {
		return nil, errors.New("connection timeout")
	}

	// Set connection parameters
	mockDB.SetConnMaxLifetime(connMaxLifetime)
	mockDB.SetMaxOpenConns(maxOpenConns)
	mockDB.SetMaxIdleConns(maxIdleConns)

	// Return the mock DB for successful initialization
	return mockDB, nil
}

// TestMockInitSQLiteDB tests the MockInitSQLiteDB function
func TestMockInitSQLiteDB(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock SQLDBInterface
	mockDB := mocks.NewMockSQLDBInterface(ctrl)

	// Test case 1: Successful initialization
	t.Run("Successful initialization", func(t *testing.T) {
		// Set up expectations
		mockDB.EXPECT().SetConnMaxLifetime(1 * time.Hour).Times(1)
		mockDB.EXPECT().SetMaxOpenConns(10).Times(1)
		mockDB.EXPECT().SetMaxIdleConns(5).Times(1)

		// Call the function
		db, err := MockInitSQLiteDB(context.Background(), ":memory:", 1*time.Second, 1*time.Hour, 10, 5, mockDB)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, mockDB, db)
	})

	// Test case 2: Invalid URI
	t.Run("Invalid URI", func(t *testing.T) {
		// Call the function
		db, err := MockInitSQLiteDB(context.Background(), "invalid-uri", 1*time.Second, 1*time.Hour, 10, 5, mockDB)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "failed to open database")
	})

	// Test case 3: Connection timeout
	t.Run("Connection timeout", func(t *testing.T) {
		// Call the function
		db, err := MockInitSQLiteDB(context.Background(), ":memory:", 0, 1*time.Hour, 10, 5, mockDB)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "connection timeout")
	})
}
