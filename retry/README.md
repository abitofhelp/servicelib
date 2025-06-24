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
- **Distributed Tracing**: Optional integration with OpenTelemetry for tracing retry operations.

## Usage

See the following examples for how to use the retry package:

- [Basic Usage](../examples/retry/basic_usage_example.go) - A simple example of how to use the retry package with default configuration.
- [Custom Configuration](../examples/retry/custom_configuration_example.go) - How to customize retry parameters like max retries, backoff, and jitter.
- [Error Detection](../examples/retry/error_detection_example.go) - How to use custom error detection to determine if an error is retryable.
- [Logging and Tracing](../examples/retry/logging_tracing_example.go) - How to use the retry package with logging and tracing.

> **Note**: The `retry.IsNetworkError`, `retry.IsTimeoutError`, and `retry.IsTransientError` functions are deprecated and will be removed in a future version. Use the corresponding functions from the `errors` package instead.

### Logging and Tracing

The retry package supports structured logging and distributed tracing. You can provide a logger and tracer using the `DoWithOptions` function.

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
package main

import "github.com/abitofhelp/servicelib/errors"

func checkErrorTypes(err error) {
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
}
```

## Helper Functions

The package provides several helper functions to determine if an error is retryable:

- **IsNetworkError** (Deprecated): Returns true if the error is a network-related error. Use `errors.IsNetworkError` instead.
- **IsTimeoutError** (Deprecated): Returns true if the error is a timeout error. Use `errors.IsTimeout` instead.
- **IsTransientError** (Deprecated): Returns true if the error is a transient error. Use `errors.IsTransientError` instead.

These functions are deprecated and will be removed in a future version. Use the corresponding functions from the `errors` package instead, as they provide more comprehensive error detection and are maintained as part of the central error handling framework.

## Integration with Other Packages

The retry package is designed to work with other servicelib packages. For example, it can be used with the database package to retry database operations:

```go
package main

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
