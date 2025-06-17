// Copyright (c) 2025 A Bit of Help, Inc.

// Package db provides utilities for working with databases.
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DefaultTimeout is the default timeout for database operations
const DefaultTimeout = 30 * time.Second

// InitMongoClient initializes a MongoDB client.
// Parameters:
//   - ctx: The context for the operation
//   - uri: The MongoDB connection URI
//   - timeout: The timeout for the connection operation
//
// Returns:
//   - *mongo.Client: The initialized MongoDB client
//   - error: An error if initialization fails
func InitMongoClient(ctx context.Context, uri string, timeout time.Duration) (*mongo.Client, error) {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a context with timeout for the connection
	connectCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(connectCtx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// InitPostgresPool initializes a PostgreSQL connection pool.
// Parameters:
//   - ctx: The context for the operation
//   - uri: The PostgreSQL connection URI
//   - timeout: The timeout for the connection operation
//
// Returns:
//   - *pgxpool.Pool: The initialized PostgreSQL connection pool
//   - error: An error if initialization fails
func InitPostgresPool(ctx context.Context, uri string, timeout time.Duration) (*pgxpool.Pool, error) {
	// Create a context with timeout for the connection
	connectCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect to PostgreSQL
	pool, err := pgxpool.New(connectCtx, uri)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := pool.Ping(connectCtx); err != nil {
		return nil, err
	}

	return pool, nil
}

// InitSQLiteDB initializes a SQLite database connection.
// Parameters:
//   - ctx: The context for the operation
//   - uri: The SQLite connection URI
//   - timeout: The timeout for the connection operation
//   - maxOpenConns: The maximum number of open connections
//   - maxIdleConns: The maximum number of idle connections
//   - connMaxLifetime: The maximum lifetime of a connection
//
// Returns:
//   - *sql.DB: The initialized SQLite database connection
//   - error: An error if initialization fails
func InitSQLiteDB(ctx context.Context, uri string, timeout, connMaxLifetime time.Duration, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	// Create a context with timeout for the connection
	connectCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect to SQLite
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Ping the database to verify connection
	if err := db.PingContext(connectCtx); err != nil {
		return nil, err
	}

	return db, nil
}

// LogDatabaseConnection logs a successful database connection.
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use
//   - dbType: The type of database (e.g., "MongoDB", "PostgreSQL", "SQLite")
func LogDatabaseConnection(ctx context.Context, logger *logging.ContextLogger, dbType string) {
	logger.Info(ctx, "Connected to "+dbType)
}

// CheckPostgresHealth checks if a PostgreSQL connection is healthy.
// Parameters:
//   - ctx: The context for the operation
//   - pool: The PostgreSQL connection pool
//
// Returns:
//   - error: An error if the health check fails
func CheckPostgresHealth(ctx context.Context, pool *pgxpool.Pool) error {
	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Ping the database to verify connection
	if err := pool.Ping(healthCtx); err != nil {
		return err
	}

	return nil
}

// CheckMongoHealth checks if a MongoDB connection is healthy.
// Parameters:
//   - ctx: The context for the operation
//   - client: The MongoDB client
//
// Returns:
//   - error: An error if the health check fails
func CheckMongoHealth(ctx context.Context, client *mongo.Client) error {
	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Ping the database to verify connection
	if err := client.Ping(healthCtx, nil); err != nil {
		return err
	}

	return nil
}

// CheckSQLiteHealth checks if a SQLite connection is healthy.
// Parameters:
//   - ctx: The context for the operation
//   - db: The SQLite database connection
//
// Returns:
//   - error: An error if the health check fails
func CheckSQLiteHealth(ctx context.Context, db *sql.DB) error {
	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Ping the database to verify connection
	if err := db.PingContext(healthCtx); err != nil {
		return err
	}

	return nil
}

// ExecutePostgresTransaction executes a function within a PostgreSQL transaction.
// Parameters:
//   - ctx: The context for the operation
//   - pool: The PostgreSQL connection pool
//   - fn: The function to execute within the transaction
//
// Returns:
//   - error: An error if the transaction fails
func ExecutePostgresTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	// Create a context with timeout for the transaction
	txCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Begin transaction
	tx, err := pool.Begin(txCtx)
	if err != nil {
		return err
	}

	// Ensure transaction is rolled back if not committed
	defer func() {
		if tx != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Execute function
	if err := fn(tx); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(txCtx); err != nil {
		return err
	}

	return nil
}

// ExecuteSQLTransaction executes a function within a SQL transaction.
// Parameters:
//   - ctx: The context for the operation
//   - db: The SQL database connection
//   - fn: The function to execute within the transaction
//
// Returns:
//   - error: An error if the transaction fails
func ExecuteSQLTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	// Create a context with timeout for the transaction
	txCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Begin transaction
	tx, err := db.BeginTx(txCtx, nil)
	if err != nil {
		return err
	}

	// Ensure transaction is rolled back if not committed
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// Execute function
	if err := fn(tx); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
