// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestMongoDBInitializer tests the MongoDBInitializer function
func TestMongoDBInitializer(t *testing.T) {
	// Skip this test as it requires a real MongoDB connection
	// In a real-world scenario, we would use a mock MongoDB client
	t.Skip("Skipping TestMongoDBInitializer as it requires a real MongoDB connection")
}

// TestPostgresInitializer tests the PostgresInitializer function
func TestPostgresInitializer(t *testing.T) {
	// Skip this test as it requires a real PostgreSQL connection
	// In a real-world scenario, we would use a mock PostgreSQL pool
	t.Skip("Skipping TestPostgresInitializer as it requires a real PostgreSQL connection")
}

// TestSQLiteInitializer tests the SQLiteInitializer function
func TestSQLiteInitializer(t *testing.T) {
	// Create a context and logger
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	// Test with in-memory SQLite database (this should work without a real connection)
	t.Run("Success with in-memory database", func(t *testing.T) {
		// Call the function being tested
		db, err := SQLiteInitializer(ctx, ":memory:", logger)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, db)

		// Clean up
		if db != nil {
			db.Close()
		}
	})

	// Test with invalid URI
	t.Run("Error with invalid URI", func(t *testing.T) {
		// Call the function being tested with an invalid URI
		db, err := SQLiteInitializer(ctx, "file:nonexistent?mode=ro", logger)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

// TestDBInitializersWithNilLogger tests the initializers with a nil logger
func TestDBInitializersWithNilLogger(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Test MongoDBInitializer with nil logger
	t.Run("MongoDBInitializer with nil logger", func(t *testing.T) {
		_, err := MongoDBInitializer(ctx, "mongodb://localhost:27017", "testdb", "testcollection", nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || 
			err.Error() == "failed to initialize MongoDB client: context canceled" || 
			err.Error() == "failed to initialize MongoDB client: context deadline exceeded")
	})

	// Test PostgresInitializer with nil logger
	t.Run("PostgresInitializer with nil logger", func(t *testing.T) {
		_, err := PostgresInitializer(ctx, "postgres://localhost:5432", nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || 
			err.Error() == "failed to initialize PostgreSQL connection pool: context canceled" || 
			err.Error() == "failed to initialize PostgreSQL connection pool: context deadline exceeded")
	})

	// Test SQLiteInitializer with nil logger
	t.Run("SQLiteInitializer with nil logger", func(t *testing.T) {
		_, err := SQLiteInitializer(ctx, ":memory:", nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || 
			err.Error() == "failed to initialize SQLite database connection: context canceled" || 
			err.Error() == "failed to initialize SQLite database connection: context deadline exceeded")
	})
}
