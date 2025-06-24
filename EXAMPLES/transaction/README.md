# Transaction Package Examples

This directory contains examples demonstrating how to use the `transaction` package, which provides utilities for managing distributed transactions in Go applications using the Saga pattern.

## Examples

### 1. Basic Saga Example

[basic_saga_example.go](basic_saga_example.go)

Demonstrates a basic saga transaction where one of the operations fails, triggering rollbacks for the previously completed operations.

Key concepts:
- Using `WithTransaction` to execute operations within a transaction
- Adding operations with their corresponding rollback operations
- Handling operation failures and automatic rollbacks

### 2. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Shows a successful transaction using the `WithTransaction` function.

Key concepts:
- Using `WithTransaction` for a simple successful transaction
- Adding multiple operations to a transaction
- Handling the transaction result

### 3. Context Timeout Example

[context_timeout_example.go](context_timeout_example.go)

Demonstrates using context with timeout in transactions.

Key concepts:
- Creating a context with timeout
- Respecting context cancellation in operations
- Handling context deadline exceeded errors
- Rollback behavior when context times out

### 4. Custom Transaction Example

[custom_transaction_example.go](custom_transaction_example.go)

Shows how to create a transaction explicitly using `NewTransaction` and then execute it.

Key concepts:
- Creating a transaction explicitly with `NewTransaction`
- Adding operations to the transaction
- Using `NoopRollback` for operations that don't need a rollback
- Executing the transaction with `Execute`

### 5. Error Handling Example

[error_handling_example.go](error_handling_example.go)

Demonstrates error handling with checked rollbacks.

Key concepts:
- Using `CheckedRollback` for better error reporting
- Handling rollback errors
- Extracting detailed error information from a `RollbackError`

### 6. Idempotent Operations Example

[idempotent_operations_example.go](idempotent_operations_example.go)

Shows how to implement idempotent operations in transactions.

Key concepts:
- Designing operations to be idempotent
- Implementing idempotent Create and Delete operations
- Running multiple transactions with the same operations
- Handling rollbacks with idempotent operations

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/transaction/basic_saga_example.go
```

## Additional Resources

For more information about the transaction package, see the [transaction package documentation](../../transaction/README.md).