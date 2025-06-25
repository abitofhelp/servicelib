# Retry Package
The retry package provides functionality for retrying operations with configurable backoff and jitter. It is designed to work seamlessly with the servicelib error framework.


## Overview

The Retry package is a core component of the ServiceLib library that provides robust functionality for executing operations that may fail transiently and need to be retried. It implements industry-standard retry patterns including exponential backoff with jitter to help distributed systems recover from temporary failures while avoiding the "thundering herd" problem.

This package is designed to work seamlessly with the ServiceLib error framework, providing consistent error handling and detailed logging of retry attempts. It's particularly useful for operations that interact with external systems like databases, APIs, or network services where transient failures are common.

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

- [Basic Usage](../EXAMPLES/retry/basic_usage_example.go) - A simple example of how to use the retry package with default configuration.
- [Custom Configuration](../EXAMPLES/retry/custom_configuration_example.go) - How to customize retry parameters like max retries, backoff, and jitter.
- [Error Detection](../EXAMPLES/retry/error_detection_example.go) - How to use custom error detection to determine if an error is retryable.
- [Logging and Tracing](../EXAMPLES/retry/logging_tracing_example.go) - How to use the retry package with logging and tracing.

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


## Installation

```bash
go get github.com/abitofhelp/servicelib/retry
```

## Quick Start

See the [Quick Start example](../EXAMPLES/retry/quickstart_example.go) for a complete, runnable example of how to use the retry.

## Configuration

See the [Configuration example](../EXAMPLES/retry/configuration_example.go) for a complete, runnable example of how to configure the retry.

## API Documentation


### Core Types

The retry package provides several core types for configuring and executing retryable operations:

#### Config

The `Config` struct holds the configuration for retry operations, including:

- **MaxRetries**: Maximum number of retry attempts
- **InitialBackoff**: Initial backoff duration before the first retry
- **MaxBackoff**: Maximum backoff duration for any retry attempt
- **BackoffFactor**: Factor by which the backoff increases with each retry
- **JitterFactor**: Factor for random jitter to prevent thundering herd problems
- **RetryableErrors**: Specific errors that should be considered retryable

The `Config` struct provides fluent methods for configuration:

- `WithMaxRetries(max int)`: Sets the maximum number of retry attempts
- `WithInitialBackoff(duration time.Duration)`: Sets the initial backoff duration
- `WithMaxBackoff(duration time.Duration)`: Sets the maximum backoff duration
- `WithBackoffFactor(factor float64)`: Sets the backoff factor
- `WithJitterFactor(factor float64)`: Sets the jitter factor
- `WithRetryableErrors(errors []error)`: Sets specific errors that should be retried

#### RetryableFunc

The `RetryableFunc` type represents a function that can be retried:

```
type RetryableFunc func(ctx context.Context) error
```

This is the function you want to execute with retry logic. It should return an error if the operation fails and should be retried.

#### IsRetryableError

The `IsRetryableError` type is a function that determines if an error is retryable:

```
type IsRetryableError func(err error) bool
```

This function is called after each failed attempt to determine if another retry should be attempted.

See the [Configuration example](../EXAMPLES/retry/configuration_example.go) for a complete, runnable example of how to configure retry operations.

### Key Methods

The retry package provides several key methods for executing retryable operations:

#### Do

The `Do` method executes a function with retry logic:

```
func Do(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError) error
```

This is the main method for executing a function with retry logic. It takes a context, the function to retry, a retry configuration, and a function to determine if an error is retryable.

#### DoWithOptions

The `DoWithOptions` method executes a function with retry logic and additional options:

```
func DoWithOptions(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError, options Options) error
```

This method is similar to `Do`, but it allows you to specify additional options like a logger and tracer for more detailed logging and tracing of retry attempts.

#### DefaultConfig

The `DefaultConfig` function returns a default retry configuration:

```
func DefaultConfig() Config
```

This function returns a Config struct with sensible default values for retry operations.

See the [Basic Usage example](../EXAMPLES/retry/basic_usage_example.go) for a complete, runnable example of how to use these methods.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/retry/basic_usage_example.go) - Shows basic usage of the retry
- [Advanced Configuration](../EXAMPLES/retry/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/retry/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Use Exponential Backoff**: Always use exponential backoff for retries to avoid overwhelming the system during recovery. The default configuration provides a good starting point.

2. **Add Jitter**: Include jitter in your retry strategy to prevent synchronized retry attempts from multiple clients, which can cause "thundering herd" problems.

3. **Set Appropriate Timeouts**: Configure context timeouts that are appropriate for your operation, considering the maximum possible retry time.

4. **Limit Retry Attempts**: Set a reasonable maximum number of retry attempts to prevent infinite retry loops. Consider the nature of the operation and the likelihood of recovery.

5. **Use Custom Error Detection**: Implement custom error detection logic to determine which errors are truly retryable. Not all errors should trigger a retry.

6. **Log Retry Attempts**: Enable logging to track retry attempts, which can help diagnose issues with external systems.

7. **Consider Circuit Breakers**: For critical systems, combine retries with circuit breakers to prevent cascading failures when a system is consistently unavailable.

## Troubleshooting

### Common Issues

#### Excessive Retry Attempts

**Issue**: Operations are being retried too many times, causing increased latency or resource consumption.

**Solution**: Adjust the `MaxRetries` parameter in the Config struct to limit the number of retry attempts. Also, ensure your error detection function (`IsRetryableError`) is correctly identifying which errors should be retried.

#### Insufficient Backoff

**Issue**: Retry attempts are happening too quickly, overwhelming the target system during recovery.

**Solution**: Increase the `InitialBackoff` and `BackoffFactor` parameters to create more space between retry attempts. Consider setting a higher `MaxBackoff` value for operations that might require longer recovery times.

#### Context Timeouts

**Issue**: Operations are timing out before retry attempts can complete.

**Solution**: Ensure that your context timeout is sufficient to accommodate the maximum possible retry time. You can calculate this as approximately: `initialBackoff * (1 - backoffFactor^maxRetries) / (1 - backoffFactor)` plus some buffer for the actual operation time.

#### Synchronized Retry Storms

**Issue**: Multiple clients are retrying operations at the same time, causing load spikes.

**Solution**: Ensure you're using jitter by setting the `JitterFactor` parameter to a value greater than 0 (typically between 0.1 and 0.3). This adds randomness to the backoff times, preventing synchronized retry attempts.

## Related Components

- [Errors](../errors/README.md) - The errors package provides error types and utilities used by the retry package, including RetryError and functions for checking error types.
- [Context](../context/README.md) - The context package provides utilities for working with context, which is used extensively in the retry package for cancellation and timeout handling.
- [Logging](../logging/README.md) - The logging package provides structured logging, which is used by the retry package for logging retry attempts and errors.
- [Telemetry](../telemetry/README.md) - The telemetry package provides tracing, which can be integrated with the retry package for tracing retry operations.
- [Database](../db/README.md) - The database package can use the retry package to implement retry logic for database operations.

## Contributing

Contributions to this retry are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This package is part of the servicelib project and is licensed under the same terms.
