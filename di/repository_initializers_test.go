// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestGenericMongoInitializer tests the GenericMongoInitializer function
func TestGenericMongoInitializer(t *testing.T) {
	// Skip this test as it requires a real MongoDB connection
	// In a real-world scenario, we would use a mock MongoDB client
	t.Skip("Skipping TestGenericMongoInitializer as it requires a real MongoDB connection")
}

// TestGenericPostgresInitializer tests the GenericPostgresInitializer function
func TestGenericPostgresInitializer(t *testing.T) {
	// Skip this test as it requires a real PostgreSQL connection
	// In a real-world scenario, we would use a mock PostgreSQL pool
	t.Skip("Skipping TestGenericPostgresInitializer as it requires a real PostgreSQL connection")
}

// TestGenericSQLiteInitializer tests the GenericSQLiteInitializer function
func TestGenericSQLiteInitializer(t *testing.T) {
	// Create a context and logger
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	// Test with in-memory SQLite database (this should work without a real connection)
	t.Run("Success with in-memory database", func(t *testing.T) {
		// Call the function being tested
		result, err := GenericSQLiteInitializer(ctx, ":memory:", logger)
		
		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		
		// Clean up
		if db, ok := result.(*interface{}); ok && db != nil {
			// This is a bit tricky since we're dealing with interface{}
			// In a real test, we would type assert to sql.DB and close it
		}
	})
	
	// Test with invalid URI
	t.Run("Error with invalid URI", func(t *testing.T) {
		// Call the function being tested with an invalid URI
		result, err := GenericSQLiteInitializer(ctx, "file:nonexistent?mode=ro", logger)
		
		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestGenericInitializersWithNilLogger tests the initializers with a nil logger
func TestGenericInitializersWithNilLogger(t *testing.T) {
	// Create a context
	ctx := context.Background()
	
	// Test GenericMongoInitializer with nil logger
	t.Run("GenericMongoInitializer with nil logger", func(t *testing.T) {
		_, err := GenericMongoInitializer(ctx, "mongodb://localhost:27017", "testdb", "testcollection", nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || 
			err.Error() == "failed to initialize MongoDB client: context canceled" || 
			err.Error() == "failed to initialize MongoDB client: context deadline exceeded")
	})
	
	// Test GenericPostgresInitializer with nil logger
	t.Run("GenericPostgresInitializer with nil logger", func(t *testing.T) {
		_, err := GenericPostgresInitializer(ctx, "postgres://localhost:5432", nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || 
			err.Error() == "failed to initialize PostgreSQL connection pool: context canceled" || 
			err.Error() == "failed to initialize PostgreSQL connection pool: context deadline exceeded")
	})
	
	// Test GenericSQLiteInitializer with nil logger
	t.Run("GenericSQLiteInitializer with nil logger", func(t *testing.T) {
		_, err := GenericSQLiteInitializer(ctx, ":memory:", nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || 
			err.Error() == "failed to initialize SQLite database connection: context canceled" || 
			err.Error() == "failed to initialize SQLite database connection: context deadline exceeded")
	})
}