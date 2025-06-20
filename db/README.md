# Database Package

## Overview

The `db` package provides utilities for database connection management and operations in Go applications. It supports multiple database types and provides features for connection pooling, transaction management, and query execution.

## Features

- **Connection Management**:
  - Connection pooling
  - Automatic reconnection
  - Health checks

- **Supported Databases**:
  - PostgreSQL (via pgx)
  - SQLite
  - MongoDB

- **Features**:
  - Transaction management
  - Query execution with retries
  - Result mapping
  - Migrations

## Installation

```bash
go get github.com/abitofhelp/servicelib/db
```

## Quick Start

See the [Connection Example](../examples/db/connection_example.go) for a complete, runnable example of how to connect to a database and execute basic operations.

## Configuration

See the [Repository Example](../examples/db/repository_example.go) for a complete, runnable example of how to configure the database package.

## API Documentation

### Core Types

#### Database

The `Database` interface provides methods for database operations.

See the [Connection Example](../examples/db/connection_example.go) for a complete, runnable example of how to use the Database interface.

#### Transaction

The `Transaction` interface provides methods for transaction management.

See the [Transaction Example](../examples/db/transaction_example.go) for a complete, runnable example of how to use transactions.

### Key Methods

#### New

The `New` function creates a new database connection.

See the [Connection Example](../examples/db/connection_example.go) for a complete, runnable example.

#### Transaction

The `Transaction` method executes a function within a transaction.

See the [Transaction Example](../examples/db/transaction_example.go) for a complete, runnable example of how to use the Transaction method.

#### QueryWithRetries

The `QueryWithRetries` function executes a query with automatic retries.

See the [Retry Example](../examples/db/retry_example.go) for a complete, runnable example of how to use the QueryWithRetries function.

## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Connection Example](../examples/db/connection_example.go) - Shows how to connect to a database
- [Query Example](../examples/db/query_example.go) - Shows how to execute queries
- [Transaction Example](../examples/db/transaction_example.go) - Shows how to use transactions
- [Retry Example](../examples/db/retry_example.go) - Shows how to use query retries
- [Repository Example](../examples/db/repository_example.go) - Shows how to implement the repository pattern
- [Health Check Example](../examples/db/health_check_example.go) - Shows how to implement database health checks

## Best Practices

1. **Connection Pooling**: Configure connection pools based on your application's needs and the database's capacity.

2. **Transactions**: Use transactions for operations that need to be atomic.

3. **Prepared Statements**: Use prepared statements to prevent SQL injection and improve performance.

4. **Context**: Always pass a context to database operations to enable cancellation and timeouts.

5. **Error Handling**: Handle database errors appropriately, distinguishing between different types of errors.

6. **Repository Pattern**: Use the repository pattern to abstract database access and make testing easier.

7. **Migrations**: Use database migrations to manage schema changes.

## Troubleshooting

### Common Issues

#### Connection Pool Exhaustion

**Issue**: The application runs out of database connections.

**Solution**: Configure the connection pool appropriately by setting max open connections, max idle connections, and connection lifetime. Monitor connection usage and adjust as needed.

#### Query Timeouts

**Issue**: Database queries are timing out.

**Solution**: Use context with timeouts for long-running queries. Optimize queries that are consistently slow. Consider using indexes or query optimization.

#### Transaction Deadlocks

**Issue**: Transactions are deadlocking.

**Solution**: Keep transactions short and focused. Avoid holding transactions open for long periods. Consider the order of operations to minimize lock contention.

## Related Components

- [Config](../config/README.md) - The config component is used to configure database connections.
- [Logging](../logging/README.md) - The logging component can be used to log database operations.
- [Telemetry](../telemetry/README.md) - The telemetry component can be used to monitor database performance.
- [Transaction](../transaction/README.md) - The transaction component provides utilities for distributed transactions.

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
