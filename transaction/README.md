# Transaction

## Overview

The Transaction package provides utilities for managing transactions in distributed systems. It offers tools and patterns for ensuring data consistency across multiple services and databases in a distributed environment, focusing on maintaining the ACID properties (Atomicity, Consistency, Isolation, Durability) as much as possible in distributed scenarios where traditional database transactions are not feasible.

## Features

- **Saga Pattern**: Implementation of the Saga pattern for distributed transactions
- **Automatic Rollback**: Automatic rollback of operations when a transaction fails
- **Error Handling**: Detailed error information when transactions fail
- **Context Support**: Support for context cancellation and timeouts
- **Logging**: Comprehensive logging of transaction events

## Installation

```bash
go get github.com/abitofhelp/servicelib/transaction
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/transaction/basic_usage_example/README.md) for a complete, runnable example of how to use the transaction package.

## Configuration

The transaction package does not require specific configuration. It is designed to be used with your existing logging and context infrastructure.

## API Documentation

### Core Types

The transaction package provides several core types for managing distributed transactions.

#### Transaction

The Transaction type is the main struct in the saga subpackage that manages a sequence of operations and their rollbacks. It provides methods for adding operations and executing the transaction.

See the [Basic Saga example](../EXAMPLES/transaction/basic_saga_example/README.md) for a complete, runnable example of how to use the Transaction type.

#### Operation

The Operation type is a function type that represents a local transaction. It takes a context.Context parameter and returns an error.

#### RollbackOperation

The RollbackOperation type is a function type that represents a compensating transaction. It takes a context.Context parameter and returns an error.

### Key Methods

The transaction package provides several key methods for managing distributed transactions.

#### WithTransaction

The WithTransaction function is a helper function for executing a function within a transaction. It creates a new transaction, executes the provided function, and then executes the transaction.

See the [Basic Usage example](../EXAMPLES/transaction/basic_usage_example/README.md) for a complete, runnable example of how to use the WithTransaction function.

#### AddOperation

The AddOperation method adds an operation and its corresponding rollback operation to a transaction.

#### Execute

The Execute method executes all operations in a transaction. If any operation fails, it rolls back all previously executed operations in reverse order.

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/transaction/basic_usage_example/README.md) - Shows basic usage of the transaction package
- [Basic Saga](../EXAMPLES/transaction/basic_saga_example/README.md) - Shows how to use the Saga pattern for distributed transactions
- [Error Handling](../EXAMPLES/transaction/error_handling_example/README.md) - Shows how to handle errors in transactions
- [Context Timeout](../EXAMPLES/transaction/context_timeout_example/README.md) - Shows how to handle context timeouts
- [Custom Transaction](../EXAMPLES/transaction/custom_transaction_example/README.md) - Shows how to create custom transaction types
- [Idempotent Operations](../EXAMPLES/transaction/idempotent_operations_example/README.md) - Shows how to create idempotent operations

## Best Practices

1. **Define Clear Boundaries**: Clearly define the boundaries of your transactions to minimize the scope and duration.
2. **Keep Operations Idempotent**: Design operations to be idempotent so they can be safely retried.
3. **Design Effective Rollbacks**: Ensure rollback operations effectively undo the effects of their corresponding operations.
4. **Handle Rollback Failures**: Be prepared to handle failures in rollback operations, possibly requiring manual intervention.
5. **Use Timeouts**: Set appropriate timeouts for operations to prevent transactions from hanging indefinitely.

## Troubleshooting

### Common Issues

#### Rollback Failures

If a rollback operation fails, the transaction will continue to execute the remaining rollback operations but will include information about the failed rollbacks in the returned error. You may need to implement additional error handling or manual intervention for these cases.

#### Context Cancellation

If the context is canceled during a transaction, the transaction will be aborted and any completed operations will be rolled back. Ensure your context management is appropriate for the transaction's expected duration.

## Related Components

- [Errors](../errors/README.md) - The transaction package uses the errors package for error handling and reporting.
- [Logging](../logging/README.md) - The transaction package uses the logging package for logging transaction events.
- [Context](../context/README.md) - The transaction package uses the context package for context management.

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
