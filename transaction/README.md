# Transaction Package

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

## Usage

### Basic Saga Transaction

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/abitofhelp/servicelib/transaction/saga"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    
    // Create a context
    ctx := context.Background()
    
    // Execute operations within a transaction
    err := saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
        // Add operations with their corresponding rollback operations
        tx.AddOperation(
            // Operation to create a user
            func(ctx context.Context) error {
                fmt.Println("Creating user...")
                // In a real application, this would call a user service
                return nil
            },
            // Rollback operation to delete the user if a later operation fails
            func(ctx context.Context) error {
                fmt.Println("Rolling back user creation...")
                // In a real application, this would call the user service to delete the user
                return nil
            },
        )
        
        tx.AddOperation(
            // Operation to create an order
            func(ctx context.Context) error {
                fmt.Println("Creating order...")
                // In a real application, this would call an order service
                return nil
            },
            // Rollback operation to delete the order if a later operation fails
            func(ctx context.Context) error {
                fmt.Println("Rolling back order creation...")
                // In a real application, this would call the order service to delete the order
                return nil
            },
        )
        
        tx.AddOperation(
            // Operation to process payment
            func(ctx context.Context) error {
                fmt.Println("Processing payment...")
                // Simulate a failure in payment processing
                return fmt.Errorf("payment processing failed")
            },
            // Rollback operation to refund the payment
            func(ctx context.Context) error {
                fmt.Println("Rolling back payment processing...")
                // In a real application, this would call the payment service to issue a refund
                return nil
            },
        )
        
        return nil
    })
    
    if err != nil {
        fmt.Printf("Transaction failed: %v\n", err)
    } else {
        fmt.Println("Transaction completed successfully")
    }
    
    // Output:
    // Creating user...
    // Creating order...
    // Processing payment...
    // Rolling back order creation...
    // Rolling back user creation...
    // Transaction failed: transaction execution failed: transaction operation failed: payment processing failed
}
```

### Custom Transaction Execution

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/abitofhelp/servicelib/transaction/saga"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    
    // Create a context
    ctx := context.Background()
    
    // Create a transaction
    tx := saga.NewTransaction(logger)
    
    // Add operations with their corresponding rollback operations
    tx.AddOperation(
        // Operation to reserve inventory
        func(ctx context.Context) error {
            fmt.Println("Reserving inventory...")
            return nil
        },
        // Rollback operation to release inventory
        func(ctx context.Context) error {
            fmt.Println("Releasing inventory...")
            return nil
        },
    )
    
    tx.AddOperation(
        // Operation to create shipment
        func(ctx context.Context) error {
            fmt.Println("Creating shipment...")
            return nil
        },
        // Rollback operation to cancel shipment
        func(ctx context.Context) error {
            fmt.Println("Cancelling shipment...")
            return nil
        },
    )
    
    tx.AddOperation(
        // Operation to send notification
        func(ctx context.Context) error {
            fmt.Println("Sending notification...")
            return nil
        },
        // No-op rollback since notifications can't be "unsent"
        saga.NoopRollback(),
    )
    
    // Execute the transaction
    if err := tx.Execute(ctx); err != nil {
        fmt.Printf("Transaction failed: %v\n", err)
    } else {
        fmt.Println("Transaction completed successfully")
    }
    
    // Output:
    // Reserving inventory...
    // Creating shipment...
    // Sending notification...
    // Transaction completed successfully
}
```

### Error Handling with Checked Rollbacks

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/abitofhelp/servicelib/transaction/saga"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    
    // Create a context
    ctx := context.Background()
    
    // Create a transaction
    tx := saga.NewTransaction(logger)
    
    // Add operations with checked rollbacks
    tx.AddOperation(
        // Operation to update account balance
        func(ctx context.Context) error {
            fmt.Println("Updating account balance...")
            return nil
        },
        // Checked rollback with operation name and error message
        saga.CheckedRollback(
            func(ctx context.Context) error {
                fmt.Println("Restoring account balance...")
                return nil
            },
            "RestoreBalance",
            "Failed to restore account balance",
        ),
    )
    
    tx.AddOperation(
        // Operation to record transaction history
        func(ctx context.Context) error {
            fmt.Println("Recording transaction history...")
            return nil
        },
        // Checked rollback with additional details
        saga.CheckedRollbackWithDetails(
            func(ctx context.Context) error {
                fmt.Println("Marking transaction as reversed...")
                return nil
            },
            "ReverseTransaction",
            "Failed to mark transaction as reversed",
            map[string]interface{}{
                "transaction_type": "account_update",
                "importance": "high",
            },
        ),
    )
    
    // Execute the transaction
    if err := tx.Execute(ctx); err != nil {
        fmt.Printf("Transaction failed: %v\n", err)
    } else {
        fmt.Println("Transaction completed successfully")
    }
    
    // Output:
    // Updating account balance...
    // Recording transaction history...
    // Transaction completed successfully
}
```

## Best Practices

1. **Define Clear Compensation Actions**: For each operation, define a clear rollback operation that undoes its effects.

   ```go
   tx.AddOperation(
       // Operation
       func(ctx context.Context) error {
           return createResource(ctx)
       },
       // Compensation action
       func(ctx context.Context) error {
           return deleteResource(ctx)
       },
   )
   ```

2. **Handle Rollback Failures**: Be prepared for rollback operations to fail and handle those failures appropriately.

   ```go
   // Use CheckedRollback to add error context
   saga.CheckedRollback(
       func(ctx context.Context) error {
           return deleteResource(ctx)
       },
       "DeleteResource",
       "Failed to delete resource during rollback",
   )
   ```

3. **Use Context Properly**: Pass context through all operations to ensure proper cancellation and timeout handling.

   ```go
   // Create a context with timeout
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   
   // Use the context in the transaction
   err := saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
       // ...
   })
   ```

4. **Log Transaction Events**: Enable detailed logging to track transaction progress and diagnose issues.

   ```go
   // Create a logger with appropriate level
   logger, _ := zap.NewProduction()
   
   // Pass the logger to the transaction
   tx := saga.NewTransaction(logger)
   ```

5. **Idempotent Operations**: Design operations and rollbacks to be idempotent when possible.

   ```go
   // Idempotent operation example
   func createResourceIdempotent(ctx context.Context, id string) error {
       // Check if resource already exists
       if resourceExists(ctx, id) {
           return nil
       }
       // Create the resource
       return createResource(ctx, id)
   }
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.