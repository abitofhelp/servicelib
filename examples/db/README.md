# Database Package Examples

This directory contains examples demonstrating how to use the `db` package, which provides utilities for working with databases in Go applications. The package offers support for different database types, connection management, health checks, querying, transactions, and more.

## Examples

### 1. Connection Example

[connection_example.go](connection_example.go)

Demonstrates how to connect to different types of databases.

Key concepts:
- Connecting to PostgreSQL using `InitPostgresPool`
- Connecting to SQLite using `InitSQLiteDB` with connection pool settings
- Connecting to MongoDB using `InitMongoClient`
- Using custom timeouts for database connections

### 2. Health Check Example

[health_check_example.go](health_check_example.go)

Shows how to perform health checks on database connections.

Key concepts:
- Checking PostgreSQL health using `CheckPostgresHealth`
- Checking SQLite health using `CheckSQLiteHealth`
- Checking MongoDB health using `CheckMongoHealth`
- Setting up periodic health checks using a ticker
- Performing health checks with a custom timeout

### 3. Query Example

[query_example.go](query_example.go)

Demonstrates how to execute queries with different database types.

Key concepts:
- Executing basic queries with PostgreSQL
- Executing queries with SQLite, showing the difference in parameter placeholders
- Querying a single row using `QueryRowContext`
- Using parameterized queries to prevent SQL injection
- Handling NULL values from the database using `sql.Null*` types

### 4. Repository Example

[repository_example.go](repository_example.go)

Shows how to implement the repository pattern with the db package.

Key concepts:
- Defining a repository interface
- Implementing the repository for SQL databases
- Implementing the repository for MongoDB
- Using the repository to perform CRUD operations
- Benefits of using the repository pattern

### 5. Retry Example

[retry_example.go](retry_example.go)

Demonstrates how to execute database queries with retries.

Key concepts:
- Executing a basic query with retries using `QueryWithRetries`
- Customizing retry behavior with different retry options
- Retrying transactions using `TransactionWithRetries`
- Handling specific errors with custom retry conditions
- Using exponential backoff with jitter to prevent thundering herd problems

### 6. Transaction Example

[transaction_example.go](transaction_example.go)

Shows how to use transactions with different database types.

Key concepts:
- Executing a transaction with PostgreSQL using `ExecutePostgresTransaction`
- Executing a transaction with SQLite using `ExecuteSQLTransaction`
- Handling errors within transactions
- Information about nested transactions and savepoints
- Information about transaction isolation levels

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/db/connection_example.go
```

## Additional Resources

For more information about the db package, see the [db package documentation](../../db/README.md).