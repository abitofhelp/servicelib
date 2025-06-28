# transaction basic_saga_example Example

## Overview

This example demonstrates the Saga pattern for distributed transactions using the ServiceLib transaction package. It shows how to create a transaction that fails during execution, triggering the rollback of previously executed operations.

## Features

- **Saga Pattern**: Implementation of the Saga pattern for distributed transactions
- **Automatic Rollback**: Automatic rollback of operations when a transaction fails
- **Error Handling**: Detailed error information when transactions fail
- **Multiple Operations**: Adding multiple operations to a transaction

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Transaction Creation

The example starts by creating a logger and a context, then uses the WithTransaction function to create and execute a transaction. The WithTransaction function takes a context, a logger, and a function that receives a transaction object.

### Adding Operations

The example adds three operations to the transaction, each with a corresponding rollback operation:

1. **User Creation**: Creates a user and defines a rollback operation to delete the user if a later operation fails.
2. **Order Creation**: Creates an order and defines a rollback operation to delete the order if a later operation fails.
3. **Payment Processing**: Processes a payment but intentionally fails with an error, triggering the rollback of the previous operations.

### Error Handling

The example handles the result of the transaction execution, printing an error message that includes details about the failed operation and the rollback operations.

## Expected Output

```
Creating user...
Creating order...
Processing payment...
Rolling back order creation...
Rolling back user creation...
Transaction failed: transaction execution failed: transaction operation failed: payment processing failed
```

## Related Examples


- [basic_usage_example](../basic_usage_example/README.md) - Related example for basic_usage_example
- [context_timeout_example](../context_timeout_example/README.md) - Related example for context_timeout_example
- [custom_transaction_example](../custom_transaction_example/README.md) - Related example for custom_transaction_example

## Related Components

- [transaction Package](../../../transaction/README.md) - The transaction package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
