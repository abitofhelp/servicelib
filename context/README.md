# Context Package

The `context` package extends Go's standard context package with additional utilities for request handling, cancellation, and value propagation. It provides a set of helper functions and types to make working with contexts easier and more type-safe.

## Features

- **Value Management**: Strongly typed context values
- **Timeout Management**: Utilities for working with context deadlines
- **Cancellation**: Simplified cancellation patterns
- **Propagation**: Utilities for propagating context values across service boundaries

## Installation

```bash
go get github.com/abitofhelp/servicelib/context
```

## Usage

### Value Management

```go
package main

import (
    "fmt"
    "context"
    "github.com/abitofhelp/servicelib/context"
)

// Define a type-safe key
type userIDKey struct{}

func main() {
    // Create a context with a value
    ctx := context.WithValue(context.Background(), userIDKey{}, "user-123")
    
    // Get a value from the context
    userID, ok := context.GetValue[string](ctx, userIDKey{})
    if ok {
        fmt.Println("User ID:", userID)
    }
    
    // Using helper functions
    ctx = context.WithUserID(ctx, "user-456")
    userID = context.GetUserID(ctx)
    fmt.Println("User ID:", userID)
}
```

### Timeout Management

```go
package main

import (
    "fmt"
    "time"
    "context"
    "github.com/abitofhelp/servicelib/context"
)

func main() {
    // Create a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Check if the context has a deadline
    if deadline, ok := context.GetDeadline(ctx); ok {
        fmt.Printf("Context will expire at: %v\n", deadline)
        fmt.Printf("Time remaining: %v\n", context.GetRemainingTime(ctx))
    }
    
    // Execute a function with a timeout
    result, err := context.ExecuteWithTimeout(ctx, func(ctx context.Context) (string, error) {
        // Simulate work
        time.Sleep(2 * time.Second)
        return "Success", nil
    })
    
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }
}
```

### Cancellation

```go
package main

import (
    "fmt"
    "time"
    "context"
    "github.com/abitofhelp/servicelib/context"
)

func main() {
    // Create a cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start a goroutine that will be cancelled
    go func() {
        select {
        case <-ctx.Done():
            fmt.Println("Goroutine cancelled")
            return
        case <-time.After(10 * time.Second):
            fmt.Println("Goroutine completed")
        }
    }()
    
    // Cancel the context after 2 seconds
    time.Sleep(2 * time.Second)
    cancel()
    
    // Wait to see the output
    time.Sleep(1 * time.Second)
}
```

### Propagation

```go
package main

import (
    "fmt"
    "context"
    "net/http"
    "github.com/abitofhelp/servicelib/context"
)

func main() {
    // Create an HTTP handler that propagates context values
    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        // Extract values from the request context
        ctx := r.Context()
        
        // Add a request ID to the context
        ctx = context.WithRequestID(ctx, "req-123")
        
        // Call a service function with the context
        result := serviceFunction(ctx)
        
        fmt.Fprintf(w, "Result: %s", result)
    })
    
    http.ListenAndServe(":8080", nil)
}

func serviceFunction(ctx context.Context) string {
    // Get values from the context
    requestID := context.GetRequestID(ctx)
    
    // Use the values
    return fmt.Sprintf("Processing request %s", requestID)
}
```

## Common Context Values

The package provides helpers for common context values:

```go
// Request ID
ctx = context.WithRequestID(ctx, "req-123")
requestID := context.GetRequestID(ctx)

// User ID
ctx = context.WithUserID(ctx, "user-456")
userID := context.GetUserID(ctx)

// Tenant ID
ctx = context.WithTenantID(ctx, "tenant-789")
tenantID := context.GetTenantID(ctx)

// Correlation ID
ctx = context.WithCorrelationID(ctx, "corr-abc")
correlationID := context.GetCorrelationID(ctx)

// Transaction ID
ctx = context.WithTransactionID(ctx, "tx-def")
transactionID := context.GetTransactionID(ctx)
```

## Best Practices

1. **Type Safety**: Use strongly typed keys for context values to avoid runtime type errors.

2. **Context Propagation**: Always propagate context through your application, especially across API boundaries.

3. **Cancellation**: Use context cancellation to properly clean up resources and prevent goroutine leaks.

4. **Timeout Management**: Set appropriate timeouts for operations to prevent hanging requests.

5. **Value Scope**: Context values should be request-scoped and not used for passing optional parameters to functions.

## License

This project is licensed under the MIT License - see the LICENSE file for details.