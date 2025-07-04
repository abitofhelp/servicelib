// Copyright (c) 2024 A Bit of Help, Inc.

// Package db provides utilities for working with databases.
package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	dberrors "github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/retry"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
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
		return nil, dberrors.NewDatabaseError("failed to connect to MongoDB", "connect", "MongoDB", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(connectCtx, nil); err != nil {
		return nil, dberrors.NewDatabaseError("failed to ping MongoDB", "ping", "MongoDB", err)
	}

	return client, nil
}

// PostgresConfig holds configuration for a PostgreSQL connection pool.
type PostgresConfig struct {
	// URI is the PostgreSQL connection URI
	URI string

	// Timeout is the timeout for connection operations
	Timeout time.Duration

	// MaxOpenConns is the maximum number of open connections
	MaxConns int32

	// MinOpenConns is the minimum number of open connections
	MinConns int32

	// MaxConnLifetime is the maximum lifetime of a connection
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the maximum idle time of a connection
	MaxConnIdleTime time.Duration

	// HealthCheckPeriod is the period between health checks
	HealthCheckPeriod time.Duration
}

// InitPostgresPool initializes a PostgreSQL connection pool.
// Parameters:
//   - ctx: The context for the operation
//   - config: The PostgreSQL connection configuration
//
// Returns:
//   - *pgxpool.Pool: The initialized PostgreSQL connection pool
//   - error: An error if initialization fails
func InitPostgresPool(ctx context.Context, config PostgresConfig) (*pgxpool.Pool, error) {
	// Create a context with timeout for the connection
	connectCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Validate configuration
	if config.URI == "" {
		return nil, dberrors.NewConfigurationError("invalid PostgreSQL URI: cannot be empty", "PostgresURI", "", nil)
	}

	// Create pool configuration
	poolConfig, err := pgxpool.ParseConfig(config.URI)
	if err != nil {
		return nil, dberrors.NewDatabaseError("failed to parse PostgreSQL URI", "parse", "PostgreSQL", err)
	}

	// Set pool settings if provided
	if config.MaxConns > 0 {
		poolConfig.MaxConns = config.MaxConns
	}
	if config.MinConns > 0 {
		poolConfig.MinConns = config.MinConns
	}
	if config.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = config.MaxConnLifetime
	}
	if config.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	}
	if config.HealthCheckPeriod > 0 {
		poolConfig.HealthCheckPeriod = config.HealthCheckPeriod
	}

	// Connect to PostgreSQL
	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, dberrors.NewDatabaseError("failed to connect to PostgreSQL", "connect", "PostgreSQL", err)
	}

	// Ping the database to verify connection
	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, dberrors.NewDatabaseError("failed to ping PostgreSQL", "ping", "PostgreSQL", err)
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
		return nil, dberrors.NewDatabaseError("failed to open SQLite database", "open", "SQLite", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Ping the database to verify connection
	if err := db.PingContext(connectCtx); err != nil {
		return nil, dberrors.NewDatabaseError("failed to ping SQLite database", "ping", "SQLite", err)
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
		return dberrors.NewDatabaseError("failed to ping PostgreSQL during health check", "health_check", "PostgreSQL", err)
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
		return dberrors.NewDatabaseError("failed to ping MongoDB during health check", "health_check", "MongoDB", err)
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
		return dberrors.NewDatabaseError("failed to ping SQLite during health check", "health_check", "SQLite", err)
	}

	return nil
}

// IsTransientError checks if an error is a transient database error that can be retried.
func IsTransientError(err error) bool {
	if err == nil {
		return false
	}

	// Check for PostgreSQL-specific transient errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Connection-related errors
		switch pgErr.Code {
		case "08000", // connection_exception
			"08003", // connection_does_not_exist
			"08006", // connection_failure
			"08001", // sqlclient_unable_to_establish_sqlconnection
			"08004", // sqlserver_rejected_establishment_of_sqlconnection
			"57P01", // admin_shutdown
			"57P02", // crash_shutdown
			"57P03", // cannot_connect_now
			"53300", // too_many_connections
			"55P03", // lock_not_available
			"55006", // object_in_use
			"40001", // serialization_failure
			"40P01", // deadlock_detected
			"XX000": // internal_error (sometimes transient)
			return true
		}
	}

	// Check for common transient error messages
	errMsg := err.Error()
	return strings.Contains(errMsg, "connection reset by peer") ||
		strings.Contains(errMsg, "broken pipe") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "connection timed out") ||
		strings.Contains(errMsg, "no connection") ||
		strings.Contains(errMsg, "connection closed") ||
		strings.Contains(errMsg, "EOF") ||
		strings.Contains(errMsg, "write: broken pipe") ||
		strings.Contains(errMsg, "i/o timeout") ||
		strings.Contains(errMsg, "too many connections") ||
		strings.Contains(errMsg, "connection terminated") ||
		strings.Contains(errMsg, "connection terminated unexpectedly") ||
		strings.Contains(errMsg, "server closed the connection unexpectedly")
}

// RetryConfig holds configuration for retry operations.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// InitialBackoff is the initial backoff duration
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration
	MaxBackoff time.Duration

	// BackoffFactor is the factor by which the backoff increases
	BackoffFactor float64

	// Logger is the logger to use for logging retries
	Logger *logging.ContextLogger
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
		BackoffFactor:  2.0,
		Logger:         logging.NewContextLogger(zap.NewNop()),
	}
}

// ExecutePostgresTransaction executes a function within a PostgreSQL transaction.
// Parameters:
//   - ctx: The context for the operation
//   - pool: The PostgreSQL connection pool
//   - fn: The function to execute within the transaction
//   - retryConfig: Optional retry configuration for transient errors
//
// Returns:
//   - error: An error if the transaction fails
func ExecutePostgresTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error, retryConfig ...RetryConfig) error {
	// Use default retry config if not provided
	config := DefaultRetryConfig()
	if len(retryConfig) > 0 {
		config = retryConfig[0]
	}

	// Convert db RetryConfig to retry.Config
	retryOpts := retry.DefaultConfig().
		WithMaxRetries(config.MaxRetries).
		WithInitialBackoff(config.InitialBackoff).
		WithMaxBackoff(config.MaxBackoff).
		WithBackoffFactor(config.BackoffFactor)

	// Create retry options with the logger
	options := retry.Options{
		Logger: config.Logger,
		Tracer: telemetry.NewNoopTracer(), // Use no-op tracer by default
	}

	// Create a function that will be retried
	retryFn := func(ctx context.Context) error {
		// Create a context with timeout for the transaction
		txCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()

		// Begin transaction
		tx, err := pool.Begin(txCtx)
		if err != nil {
			return dberrors.NewDatabaseError("failed to begin transaction", "begin", "PostgreSQL", err)
		}

		// Execute function
		err = fn(tx)
		if err != nil {
			// Rollback transaction
			rollbackCtx, rollbackCancel := context.WithTimeout(ctx, 5*time.Second)
			rollbackErr := tx.Rollback(rollbackCtx)
			rollbackCancel()

			if rollbackErr != nil {
				config.Logger.Warn(ctx, "Failed to rollback transaction", zap.Error(rollbackErr))
			}

			return dberrors.NewDatabaseError("transaction function failed", "execute", "PostgreSQL", err)
		}

		// Commit transaction
		err = tx.Commit(txCtx)
		if err != nil {
			return dberrors.NewDatabaseError("failed to commit transaction", "commit", "PostgreSQL", err)
		}

		// Transaction succeeded
		return nil
	}

	// Use the retry package to execute the function with retries
	err := retry.DoWithOptions(ctx, retryFn, retryOpts, IsTransientError, options)
	if err != nil {
		// Check if it's already a DatabaseError
		if dberrors.IsDatabaseError(err) {
			return err
		}
		// Wrap other errors
		return dberrors.NewDatabaseError("transaction failed", "retry", "PostgreSQL", err)
	}

	return nil
}

// ExecuteSQLTransaction executes a function within a SQL transaction.
// Parameters:
//   - ctx: The context for the operation
//   - db: The SQL database connection
//   - fn: The function to execute within the transaction
//   - retryConfig: Optional retry configuration for transient errors
//
// Returns:
//   - error: An error if the transaction fails
func ExecuteSQLTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error, retryConfig ...RetryConfig) error {
	// Use default retry config if not provided
	config := DefaultRetryConfig()
	if len(retryConfig) > 0 {
		config = retryConfig[0]
	}

	// Convert db RetryConfig to retry.Config
	retryOpts := retry.DefaultConfig().
		WithMaxRetries(config.MaxRetries).
		WithInitialBackoff(config.InitialBackoff).
		WithMaxBackoff(config.MaxBackoff).
		WithBackoffFactor(config.BackoffFactor)

	// Create retry options with the logger
	options := retry.Options{
		Logger: config.Logger,
		Tracer: telemetry.NewNoopTracer(), // Use no-op tracer by default
	}

	// Create a function that will be retried
	retryFn := func(ctx context.Context) error {
		// Create a context with timeout for the transaction
		txCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()

		// Begin transaction
		tx, err := db.BeginTx(txCtx, nil)
		if err != nil {
			return dberrors.NewDatabaseError("failed to begin transaction", "begin", "SQL", err)
		}

		// Execute function
		err = fn(tx)
		if err != nil {
			// Rollback transaction
			rollbackErr := tx.Rollback()

			if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
				config.Logger.Warn(ctx, "Failed to rollback transaction", zap.Error(rollbackErr))
			}

			return dberrors.NewDatabaseError("transaction function failed", "execute", "SQL", err)
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			return dberrors.NewDatabaseError("failed to commit transaction", "commit", "SQL", err)
		}

		// Transaction succeeded
		return nil
	}

	// Use the retry package to execute the function with retries
	err := retry.DoWithOptions(ctx, retryFn, retryOpts, IsTransientError, options)
	if err != nil {
		// Check if it's already a DatabaseError
		if dberrors.IsDatabaseError(err) {
			return err
		}
		// Wrap other errors
		return dberrors.NewDatabaseError("transaction failed", "retry", "SQL", err)
	}

	return nil
}
