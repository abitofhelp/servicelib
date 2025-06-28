# transaction basic_usage_example Example

## Overview

This example demonstrates the basic usage of the ServiceLib transaction package with the WithTransaction function. It shows how to create a transaction, add operations with their corresponding rollback operations, and execute the transaction.

## Features

- **WithTransaction Function**: Using the WithTransaction helper function to create and execute a transaction
- **Multiple Operations**: Adding multiple operations to a transaction
- **Rollback Operations**: Defining rollback operations for each operation
- **Success Handling**: Handling successful transaction completion

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
3. **Payment Processing**: Processes a payment and defines a rollback operation to refund the payment if a later operation fails.

### Error Handling

The example handles the result of the transaction execution, printing a success message if the transaction completes successfully or an error message if it fails.

## Expected Output

```
Creating user...
Creating order...
Processing payment...
Transaction completed successfully
```

## Related Examples


- [basic_saga_example](../basic_saga_example/README.md) - Related example for basic_saga_example
- [context_timeout_example](../context_timeout_example/README.md) - Related example for context_timeout_example
- [custom_transaction_example](../custom_transaction_example/README.md) - Related example for custom_transaction_example

## Related Components

- [transaction Package](../../../transaction/README.md) - The transaction package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
