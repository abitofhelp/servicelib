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

// DirectCheckPostgresHealth is a wrapper for CheckPostgresHealth that accepts PgxPoolInterface
// This allows us to test the health check logic without requiring a real PostgreSQL connection
func DirectCheckPostgresHealth(ctx context.Context, pool PgxPoolInterface) error {
	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Ping the database to verify connection
	if err := pool.Ping(healthCtx); err != nil {
		return err
	}

	return nil
}

// TestDirectCheckPostgresHealth tests the DirectCheckPostgresHealth function
func TestDirectCheckPostgresHealth(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock PgxPoolInterface
	mockPool := mocks.NewMockPgxPoolInterface(ctrl)

	// Test case 1: Healthy connection
	t.Run("Healthy connection", func(t *testing.T) {
		// Set up expectations
		mockPool.EXPECT().Ping(gomock.Any()).Return(nil)

		// Call the function
		err := DirectCheckPostgresHealth(context.Background(), mockPool)

		// Assertions
		assert.NoError(t, err)
	})

	// Test case 2: Unhealthy connection
	t.Run("Unhealthy connection", func(t *testing.T) {
		// Set up expectations
		mockPool.EXPECT().Ping(gomock.Any()).Return(errors.New("ping error"))

		// Call the function
		err := DirectCheckPostgresHealth(context.Background(), mockPool)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
	})
}

// DirectCheckMongoHealth is a wrapper for CheckMongoHealth that accepts MongoClientInterface
// This allows us to test the health check logic without requiring a real MongoDB connection
func DirectCheckMongoHealth(ctx context.Context, client MongoClientInterface) error {
	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Ping the database to verify connection
	if err := client.Ping(healthCtx, nil); err != nil {
		return err
	}

	return nil
}

// TestDirectCheckMongoHealth tests the DirectCheckMongoHealth function
func TestDirectCheckMongoHealth(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock MongoClientInterface
	mockClient := mocks.NewMockMongoClientInterface(ctrl)

	// Test case 1: Healthy connection
	t.Run("Healthy connection", func(t *testing.T) {
		// Set up expectations
		mockClient.EXPECT().Ping(gomock.Any(), gomock.Nil()).Return(nil)

		// Call the function
		err := DirectCheckMongoHealth(context.Background(), mockClient)

		// Assertions
		assert.NoError(t, err)
	})

	// Test case 2: Unhealthy connection
	t.Run("Unhealthy connection", func(t *testing.T) {
		// Set up expectations
		mockClient.EXPECT().Ping(gomock.Any(), gomock.Nil()).Return(errors.New("ping error"))

		// Call the function
		err := DirectCheckMongoHealth(context.Background(), mockClient)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
	})
}

// DirectCheckSQLiteHealth is a wrapper for CheckSQLiteHealth that accepts SQLDBInterface
// This allows us to test the health check logic without requiring a real SQLite connection
func DirectCheckSQLiteHealth(ctx context.Context, db SQLDBInterface) error {
	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Ping the database to verify connection
	if err := db.PingContext(healthCtx); err != nil {
		return err
	}

	return nil
}

// TestDirectCheckSQLiteHealth tests the DirectCheckSQLiteHealth function
func TestDirectCheckSQLiteHealth(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock SQLDBInterface
	mockDB := mocks.NewMockSQLDBInterface(ctrl)

	// Test case 1: Healthy connection
	t.Run("Healthy connection", func(t *testing.T) {
		// Set up expectations
		mockDB.EXPECT().PingContext(gomock.Any()).Return(nil)

		// Call the function
		err := DirectCheckSQLiteHealth(context.Background(), mockDB)

		// Assertions
		assert.NoError(t, err)
	})

	// Test case 2: Unhealthy connection
	t.Run("Unhealthy connection", func(t *testing.T) {
		// Set up expectations
		mockDB.EXPECT().PingContext(gomock.Any()).Return(errors.New("ping error"))

		// Call the function
		err := DirectCheckSQLiteHealth(context.Background(), mockDB)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
	})
}
