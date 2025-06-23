# Retry Package

The retry package provides functionality for retrying operations with configurable backoff and jitter. It is designed to work seamlessly with the servicelib error framework.

## Features

- **Configurable Retry Parameters**: Customize maximum retries, initial backoff, maximum backoff, backoff factor, and jitter factor.
- **Exponential Backoff**: Automatically increases wait time between retry attempts.
- **Jitter**: Adds randomness to backoff times to prevent thundering herd problems.
- **Context Integration**: Respects context cancellation and timeouts.
- **Error Framework Integration**: Uses the servicelib error framework for consistent error handling.
- **Retryable Error Detection**: Provides helper functions to determine if an error is retryable.

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/abitofhelp/servicelib/retry"
)

func main() {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Define a function to retry
    fn := func(ctx context.Context) error {
        // Your operation here
        // Return nil if successful, or an error if it should be retried
        return nil
    }

    // Use default retry configuration
    config := retry.DefaultConfig()

    // Execute with retry
    err := retry.Do(ctx, fn, config, nil)
    if err != nil {
        fmt.Printf("Operation failed after retries: %v\n", err)
    }
}
```

### Custom Configuration

```go
// Create a custom retry configuration
config := retry.DefaultConfig().
    WithMaxRetries(5).
    WithInitialBackoff(200 * time.Millisecond).
    WithMaxBackoff(5 * time.Second).
    WithBackoffFactor(2.5).
    WithJitterFactor(0.3)
```

### Custom Retryable Error Detection

```go
// Define a custom function to determine if an error is retryable
isRetryable := func(err error) bool {
    // Check if the error is a network error
    if retry.IsNetworkError(err) {
        return true
    }
    
    // Check if the error is a timeout error
    if retry.IsTimeoutError(err) {
        return true
    }
    
    // Check if the error is a transient error
    if retry.IsTransientError(err) {
        return true
    }
    
    // Add your custom logic here
    
    return false
}

// Execute with retry and custom retryable error detection
err := retry.Do(ctx, fn, config, isRetryable)
```

## Error Handling

The retry package integrates with the servicelib error framework. It uses the following error types:

- **RetryError**: Represents an error that occurred during a retry operation.
- **ContextError**: Represents an error that occurred due to a context cancellation or timeout.

These error types can be checked using the following functions from the errors package:

```go
import "github.com/abitofhelp/servicelib/errors"

// Check if an error is a retry error
if errors.IsRetryError(err) {
    // Handle retry error
}

// Check if an error is a context error
if errors.IsContextError(err) {
    // Handle context error
}

// Check if an error is a timeout error
if errors.IsTimeout(err) {
    // Handle timeout error
}
```

## Helper Functions

The package provides several helper functions to determine if an error is retryable:

- **IsNetworkError**: Returns true if the error is a network-related error.
- **IsTimeoutError**: Returns true if the error is a timeout error.
- **IsTransientError**: Returns true if the error is a transient error.

## Integration with Other Packages

The retry package is designed to work with other servicelib packages. For example, it can be used with the database package to retry database operations:

```go
import (
    "context"
    "database/sql"
    
    "github.com/abitofhelp/servicelib/db"
    "github.com/abitofhelp/servicelib/retry"
)

func executeWithRetry(ctx context.Context, db *sql.DB) error {
    // Define a function to retry
    fn := func(ctx context.Context) error {
        // Your database operation here
        return nil
    }
    
    // Use default retry configuration
    config := retry.DefaultConfig()
    
    // Execute with retry
    return retry.Do(ctx, fn, config, db.IsTransientError)
}
```

## License

This package is part of the servicelib project and is licensed under the same terms.