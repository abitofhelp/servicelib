# Transaction Package

## Overview

The `transaction` package provides utilities for managing distributed transactions in Go applications. It implements the Saga pattern, which is a way to maintain data consistency across multiple services, each with its own database, where traditional ACID transactions are not possible.

## Features

- **Saga Pattern Implementation**: Manage distributed transactions with compensating actions
- **Transaction Coordination**: Coordinate multiple operations as a single transaction
- **Automatic Rollback**: Roll back completed operations if any operation fails
- **Error Handling**: Detailed error reporting with operation context
- **Logging Integration**: Comprehensive logging of transaction events
- **Context Support**: Respect context cancellation and deadlines

## Installation

```bash
go get github.com/abitofhelp/servicelib/transaction
```

## Quick Start

See the [Basic Saga Example](../examples/transaction/basic_saga_example.go) for a complete, runnable example of how to use the transaction package.

## Configuration

See the [Custom Transaction Example](../examples/transaction/custom_transaction_example.go) for a complete, runnable example of how to configure the transaction package.

## API Documentation

### Core Types

#### Transaction

The `Transaction` struct is the main entry point for the transaction package. It provides methods for adding operations and executing transactions.

See the [Basic Usage Example](../examples/transaction/basic_usage_example.go) for a complete, runnable example of how to use the Transaction struct.

#### Operation

The `Operation` type represents a function that performs a specific operation within a transaction.

#### Rollback

The `Rollback` type represents a function that undoes the effects of an operation if the transaction fails.

### Key Methods

#### WithTransaction

The `WithTransaction` function provides a convenient way to execute operations within a transaction.

See the [Basic Saga Example](../examples/transaction/basic_saga_example.go) for a complete, runnable example of how to use the WithTransaction function.

#### AddOperation

The `AddOperation` method adds an operation and its corresponding rollback operation to a transaction.

#### Execute

The `Execute` method executes all operations in a transaction and rolls back completed operations if any operation fails.

#### CheckedRollback

The `CheckedRollback` function creates a rollback operation with error checking and logging.

See the [Error Handling Example](../examples/transaction/error_handling_example.go) for a complete, runnable example of how to use the CheckedRollback function.

## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Saga Example](../examples/transaction/basic_saga_example.go) - Shows how to use the WithTransaction function
- [Basic Usage Example](../examples/transaction/basic_usage_example.go) - Shows how to create and use a Transaction
- [Context Timeout Example](../examples/transaction/context_timeout_example.go) - Shows how to use context timeouts with transactions
- [Custom Transaction Example](../examples/transaction/custom_transaction_example.go) - Shows how to create a custom transaction
- [Error Handling Example](../examples/transaction/error_handling_example.go) - Shows how to handle errors in transactions
- [Idempotent Operations Example](../examples/transaction/idempotent_operations_example.go) - Shows how to create idempotent operations

## Best Practices

1. **Define Clear Compensation Actions**: For each operation, define a clear rollback operation that undoes its effects.

2. **Handle Rollback Failures**: Be prepared for rollback operations to fail and handle those failures appropriately.

3. **Use Context Properly**: Pass context through all operations to ensure proper cancellation and timeout handling.

4. **Log Transaction Events**: Enable detailed logging to track transaction progress and diagnose issues.

5. **Idempotent Operations**: Design operations and rollbacks to be idempotent when possible.

See the [Idempotent Operations Example](../examples/transaction/idempotent_operations_example.go) for a complete, runnable example of how to create idempotent operations.

## Troubleshooting

### Common Issues

#### Rollback Failures

**Issue**: Rollback operations fail, leaving the system in an inconsistent state.

**Solution**: Use `CheckedRollback` to add error context and logging to rollback operations. Consider implementing a recovery mechanism for failed rollbacks.

#### Context Cancellation

**Issue**: Transactions are not respecting context cancellation or timeouts.

**Solution**: Ensure that all operations check the context for cancellation. See the [Context Timeout Example](../examples/transaction/context_timeout_example.go) for a complete, runnable example.

#### Transaction Coordination

**Issue**: Transactions are not being coordinated properly across multiple services.

**Solution**: Ensure that each service has a well-defined API for both operations and their compensating actions. Consider using a transaction coordinator service for complex scenarios.

## Related Components

- [Logging](../logging/README.md) - The logging component is used by the transaction package for logging transaction events.
- [Context](../context/README.md) - The context component is used by the transaction package for context propagation.
- [Errors](../errors/README.md) - The errors component is used by the transaction package for error handling.

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
