# Retry Package

The retry package provides functionality for retrying operations with configurable backoff and jitter. It is designed to work seamlessly with the servicelib error framework.

## Features

- **Configurable Retry Parameters**: Customize maximum retries, initial backoff, maximum backoff, backoff factor, and jitter factor.
- **Exponential Backoff**: Automatically increases wait time between retry attempts.
- **Jitter**: Adds randomness to backoff times to prevent thundering herd problems.
- **Context Integration**: Respects context cancellation and timeouts.
- **Error Framework Integration**: Uses the servicelib error framework for consistent error handling.
- **Retryable Error Detection**: Provides helper functions to determine if an error is retryable.
- **Structured Logging**: Logs retry attempts, backoff times, and errors using the ContextLogger.
- **Distributed Tracing**: Integrates with OpenTelemetry for tracing retry operations.
- **Metrics**: Records metrics for retry attempts, success/failure, and backoff times.

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

### Logging and Tracing

The retry package supports structured logging and distributed tracing. You can provide a logger and tracer using the `DoWithOptions` function:

```go
import (
    "context"
    "time"

    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/retry"
    "go.opentelemetry.io/otel"
    "go.uber.org/zap"
)

func executeWithLogging(ctx context.Context, logger *zap.Logger) error {
    // Create a context logger
    contextLogger := logging.NewContextLogger(logger)

    // Define a function to retry
    fn := func(ctx context.Context) error {
        // Your operation here
        return nil
    }

    // Use default retry configuration
    config := retry.DefaultConfig()

    // Create options with logger
    options := retry.DefaultOptions()
    options.Logger = contextLogger

    // Execute with retry and logging
    return retry.DoWithOptions(ctx, fn, config, nil, options)
}
```

#### Custom Tracer

The retry package uses a `Tracer` interface for distributed tracing, which allows for easier testing and mocking. By default, it uses OpenTelemetry, but you can provide your own implementation or use the no-op implementation for testing or when tracing is not needed:

```go
import (
    "context"
    "time"

    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/retry"
    "go.opentelemetry.io/otel"
    "go.uber.org/zap"
)

func executeWithCustomTracing(ctx context.Context, logger *zap.Logger) error {
    // Create a context logger
    contextLogger := logging.NewContextLogger(logger)

    // Define a function to retry
    fn := func(ctx context.Context) error {
        // Your operation here
        return nil
    }

    // Use default retry configuration
    config := retry.DefaultConfig()

    // Create options with logger and custom tracer
    options := retry.Options{
        Logger: contextLogger,
        Tracer: retry.NewOtelTracer(otel.Tracer("my-custom-tracer")),
    }

    // Or use a no-op tracer when tracing is not needed
    // options := retry.Options{
    //     Logger: contextLogger,
    //     Tracer: retry.NewNoopTracer(),
    // }

    // Execute with retry, logging, and tracing
    return retry.DoWithOptions(ctx, fn, config, nil, options)
}
```

This will log detailed information about each retry attempt, including:
- When an attempt starts and finishes
- The result of each attempt (success or failure)
- Backoff durations between attempts
- Error details when attempts fail
- Final outcome of the retry operation

The logs include trace IDs and span IDs when tracing is enabled, making it easy to correlate logs with traces.

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
