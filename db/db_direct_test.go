// Copyright (c) 2025 A Bit of Help, Inc.

//go:build integration
// +build integration

package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestCheckPostgresHealthDirect tests the CheckPostgresHealth function directly
// This test requires a real PostgreSQL connection
func TestCheckPostgresHealthDirect(t *testing.T) {
	// Skip this test by default
	t.Skip("Skipping TestCheckPostgresHealthDirect as it requires a real PostgreSQL connection")

	// Initialize a PostgreSQL connection
	pool, err := pgxpool.New(context.Background(), "postgres://localhost:5432")
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Test case: Healthy connection
	err = CheckPostgresHealth(context.Background(), pool)
	assert.NoError(t, err)
}

// TestCheckMongoHealthDirect tests the CheckMongoHealth function directly
// This test requires a real MongoDB connection
func TestCheckMongoHealthDirect(t *testing.T) {
	// Skip this test by default
	t.Skip("Skipping TestCheckMongoHealthDirect as it requires a real MongoDB connection")

	// Initialize a MongoDB connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Test case: Healthy connection
	err = CheckMongoHealth(context.Background(), client)
	assert.NoError(t, err)
}

// TestCheckSQLiteHealthDirect tests the CheckSQLiteHealth function directly
// This test requires a real SQLite connection
func TestCheckSQLiteHealthDirect(t *testing.T) {
	// Skip this test by default
	t.Skip("Skipping TestCheckSQLiteHealthDirect as it requires a real SQLite connection")

	// Initialize a SQLite connection
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer db.Close()

	// Test case: Healthy connection
	err = CheckSQLiteHealth(context.Background(), db)
	assert.NoError(t, err)
}

// TestExecutePostgresTransactionDirect tests the ExecutePostgresTransaction function directly
// This test requires a real PostgreSQL connection
func TestExecutePostgresTransactionDirect(t *testing.T) {
	// Skip this test by default
	t.Skip("Skipping TestExecutePostgresTransactionDirect as it requires a real PostgreSQL connection")

	// Initialize a PostgreSQL connection
	pool, err := pgxpool.New(context.Background(), "postgres://localhost:5432")
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Test case: Successful transaction
	err = ExecutePostgresTransaction(context.Background(), pool, func(tx pgx.Tx) error {
		// Execute a simple query
		_, err := tx.Exec(context.Background(), "SELECT 1")
		return err
	})
	assert.NoError(t, err)
}

// TestExecuteSQLTransactionDirect tests the ExecuteSQLTransaction function directly
// This test requires a real SQLite connection
func TestExecuteSQLTransactionDirect(t *testing.T) {
	// Skip this test by default
	t.Skip("Skipping TestExecuteSQLTransactionDirect as it requires a real SQLite connection")

	// Initialize a SQLite connection
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer db.Close()

	// Test case: Successful transaction
	err = ExecuteSQLTransaction(context.Background(), db, func(tx *sql.Tx) error {
		// Execute a simple query
		_, err := tx.Exec("SELECT 1")
		return err
	})
	assert.NoError(t, err)
}
