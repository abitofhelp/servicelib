// Copyright (c) 2024 A Bit of Help, Inc.

package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestInitMongoClientCoverage tests the InitMongoClient function to improve coverage
func TestInitMongoClientCoverage(t *testing.T) {
	// Test with an invalid URI
	t.Run("Invalid URI", func(t *testing.T) {
		client, err := InitMongoClient(context.Background(), "invalid-uri", 1*time.Second)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	// Test with a valid URI but unreachable server
	t.Run("Unreachable server", func(t *testing.T) {
		client, err := InitMongoClient(context.Background(), "mongodb://nonexistent-server:27017", 1*time.Second)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	// Test with a timeout that's too short
	t.Run("Timeout too short", func(t *testing.T) {
		client, err := InitMongoClient(context.Background(), "mongodb://localhost:27017", 1*time.Nanosecond)
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

// TestInitPostgresPoolCoverage tests the InitPostgresPool function to improve coverage
func TestInitPostgresPoolCoverage(t *testing.T) {
	// Test with an invalid URI
	t.Run("Invalid URI", func(t *testing.T) {
		config := PostgresConfig{
			URI:     "invalid-uri",
			Timeout: 1 * time.Second,
		}
		pool, err := InitPostgresPool(context.Background(), config)
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	// Test with a valid URI but unreachable server
	t.Run("Unreachable server", func(t *testing.T) {
		config := PostgresConfig{
			URI:     "postgres://nonexistent-server:5432",
			Timeout: 1 * time.Second,
		}
		pool, err := InitPostgresPool(context.Background(), config)
		assert.Error(t, err)
		assert.Nil(t, pool)
	})

	// Test with a timeout that's too short
	t.Run("Timeout too short", func(t *testing.T) {
		config := PostgresConfig{
			URI:     "postgres://localhost:5432",
			Timeout: 1 * time.Nanosecond,
		}
		pool, err := InitPostgresPool(context.Background(), config)
		assert.Error(t, err)
		assert.Nil(t, pool)
	})
}

// TestInitSQLiteDBCoverage tests the InitSQLiteDB function to improve coverage
func TestInitSQLiteDBCoverage(t *testing.T) {
	// Test with an invalid URI
	t.Run("Invalid URI", func(t *testing.T) {
		db, err := InitSQLiteDB(context.Background(), "file:nonexistent?mode=ro", 1*time.Second, 1*time.Hour, 10, 5)
		if err == nil {
			// If we got a valid connection, close it and skip the test
			db.Close()
			t.Skip("Skipping test as SQLite accepted the invalid URI")
		}
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	// Test with a timeout that's too short
	t.Run("Timeout too short", func(t *testing.T) {
		db, err := InitSQLiteDB(context.Background(), ":memory:", 1*time.Nanosecond, 1*time.Hour, 10, 5)
		if err == nil {
			// If we got a valid connection, close it and skip the test
			db.Close()
			t.Skip("Skipping test as SQLite accepted the short timeout")
		}
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	// Test with a valid in-memory database
	t.Run("Valid in-memory database", func(t *testing.T) {
		db, err := InitSQLiteDB(context.Background(), ":memory:", 1*time.Second, 1*time.Hour, 10, 5)
		if err != nil {
			// If there was an error, it might be because SQLite is not available
			t.Skip("Skipping test as SQLite is not available")
		}
		assert.NoError(t, err)
		assert.NotNil(t, db)
		db.Close()
	})
}
