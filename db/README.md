# Database

## Overview

The Database component provides utilities for working with various database systems, including PostgreSQL, SQLite, and MongoDB. It offers connection management, health checks, transaction handling, and retry mechanisms for transient errors.

## Features

- **Multi-Database Support**: Work with PostgreSQL, SQLite, and MongoDB using a consistent API
- **Connection Management**: Initialize and configure database connections with sensible defaults
- **Health Checks**: Verify database connectivity with built-in health check functions
- **Transaction Handling**: Execute operations within transactions with automatic rollback on errors
- **Retry Mechanisms**: Automatically retry operations that fail due to transient errors
- **Error Handling**: Structured error types for better error handling and reporting

## Installation

```bash
go get github.com/abitofhelp/servicelib/db
```

## Quick Start

See the [Connection example](../EXAMPLES/db/connection/README.md) for a complete, runnable example of how to use the database component.

## API Documentation

### Core Types

#### PostgresConfig

Configuration for a PostgreSQL connection pool.

```go
type PostgresConfig struct {
    // URI is the PostgreSQL connection URI
    URI string

    // Timeout is the timeout for connection operations
    Timeout time.Duration

    // MaxConns is the maximum number of open connections
    MaxConns int32

    // MinConns is the minimum number of open connections
    MinConns int32

    // MaxConnLifetime is the maximum lifetime of a connection
    MaxConnLifetime time.Duration

    // MaxConnIdleTime is the maximum idle time of a connection
    MaxConnIdleTime time.Duration

    // HealthCheckPeriod is the period between health checks
    HealthCheckPeriod time.Duration
}
```

#### RetryConfig

Configuration for retry operations.

```go
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
```

### Key Methods

#### InitPostgresPool

Initializes a PostgreSQL connection pool.

```go
func InitPostgresPool(ctx context.Context, config PostgresConfig) (*pgxpool.Pool, error)
```

#### InitSQLiteDB

Initializes a SQLite database connection.

```go
func InitSQLiteDB(ctx context.Context, uri string, timeout, connMaxLifetime time.Duration, maxOpenConns, maxIdleConns int) (*sql.DB, error)
```

#### InitMongoClient

Initializes a MongoDB client.

```go
func InitMongoClient(ctx context.Context, uri string, timeout time.Duration) (*mongo.Client, error)
```

#### ExecutePostgresTransaction

Executes a function within a PostgreSQL transaction.

```go
func ExecutePostgresTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error, retryConfig ...RetryConfig) error
```

#### ExecuteSQLTransaction

Executes a function within a SQL transaction.

```go
func ExecuteSQLTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error, retryConfig ...RetryConfig) error
```

#### CheckPostgresHealth / CheckSQLiteHealth / CheckMongoHealth

Check if a database connection is healthy.

```go
func CheckPostgresHealth(ctx context.Context, pool *pgxpool.Pool) error
func CheckSQLiteHealth(ctx context.Context, db *sql.DB) error
func CheckMongoHealth(ctx context.Context, client *mongo.Client) error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Connection](../EXAMPLES/db/connection/README.md) - Shows how to connect to different databases
- [Health Check](../EXAMPLES/db/health_check/README.md) - Shows how to perform health checks on database connections
- [Query](../EXAMPLES/db/query/README.md) - Shows how to execute queries against databases
- [Transaction](../EXAMPLES/db/transaction/README.md) - Shows how to use transactions with automatic retries
- [Repository](../EXAMPLES/db/repository/README.md) - Shows how to implement the repository pattern
- [Retry](../EXAMPLES/db/retry/README.md) - Shows how to use retry mechanisms for database operations

## Best Practices

1. **Use Connection Pooling**: Always use connection pooling for better performance and resource management
2. **Set Appropriate Timeouts**: Configure timeouts based on expected operation durations
3. **Use Transactions**: Wrap related operations in transactions to ensure data consistency
4. **Handle Transient Errors**: Use retry mechanisms for operations that may fail due to transient errors
5. **Close Connections**: Always close database connections when they are no longer needed

## Troubleshooting

### Common Issues

#### Connection Failures

If you're experiencing connection failures, check:
- Database server is running and accessible
- Connection string is correct
- Network connectivity between your application and the database
- Firewall rules allow the connection

#### Transaction Deadlocks

If you're experiencing transaction deadlocks:
- Ensure consistent order of resource access
- Keep transactions short and focused
- Consider using optimistic concurrency control
- Monitor and analyze deadlock logs

## Related Components

- [Retry](../retry/README.md) - Used for retrying operations that fail due to transient errors
- [Errors](../errors/README.md) - Provides structured error types for database operations
- [Logging](../logging/README.md) - Used for logging database operations
- [Telemetry](../telemetry/README.md) - Used for tracing database operations

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
