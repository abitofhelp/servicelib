# Retry

## Overview

The Retry component provides functionality for retrying operations with configurable backoff and jitter. It implements an exponential backoff algorithm with jitter to help distribute retry attempts and prevent thundering herd problems in distributed systems.

## Features

- **Exponential Backoff**: Automatically increase wait time between retry attempts
- **Configurable Jitter**: Add randomness to retry intervals to prevent synchronized retries
- **Flexible Error Handling**: Customize which errors should trigger retries
- **Context Awareness**: Respect context cancellation and timeouts
- **Telemetry Integration**: Built-in OpenTelemetry tracing for monitoring retry operations
- **Logging Integration**: Comprehensive logging of retry attempts and outcomes
- **Type-Safe API**: Generic functions for type-safe operation with Go generics

## Installation

```bash
go get github.com/abitofhelp/servicelib/retry
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/abitofhelp/servicelib/retry"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create a retry configuration
    config := retry.DefaultConfig().
        WithMaxRetries(5).
        WithInitialBackoff(100 * time.Millisecond).
        WithMaxBackoff(2 * time.Second)
    
    // Define a function to retry
    operation := func(ctx context.Context) error {
        // Your operation that might fail
        return callExternalService(ctx)
    }
    
    // Define which errors should be retried
    isRetryable := func(err error) bool {
        return errors.IsNetworkError(err) || errors.IsTimeout(err)
    }
    
    // Execute with retry
    ctx := context.Background()
    err := retry.Do(ctx, operation, config, isRetryable)
    if err != nil {
        if errors.IsRetryError(err) {
            fmt.Println("Operation failed after maximum retry attempts")
        } else {
            fmt.Println("Operation failed with non-retryable error:", err)
        }
        return
    }
    
    fmt.Println("Operation succeeded")
}

func callExternalService(ctx context.Context) error {
    // Simulate an external service call
    return nil
}
```

## API Documentation

### Core Types

#### Config

Configuration for retry operations.

```go
type Config struct {
    MaxRetries      int           // Maximum number of retry attempts
    InitialBackoff  time.Duration // Initial backoff duration
    MaxBackoff      time.Duration // Maximum backoff duration
    BackoffFactor   float64       // Factor by which the backoff increases
    JitterFactor    float64       // Factor for random jitter (0-1)
    RetryableErrors []error       // Errors that are considered retryable
}
```

#### Options

Additional options for retry operations.

```go
type Options struct {
    // Logger is used for logging retry operations
    Logger *logging.ContextLogger
    // Tracer is used for tracing retry operations
    Tracer telemetry.Tracer
}
```

#### RetryableFunc

A function that can be retried.

```go
type RetryableFunc func(ctx context.Context) error
```

#### IsRetryableError

A function that determines if an error is retryable.

```go
type IsRetryableError func(err error) bool
```

### Key Methods

#### DefaultConfig

Returns a default retry configuration.

```go
func DefaultConfig() Config
```

#### DefaultOptions

Returns default options for retry operations.

```go
func DefaultOptions() Options
```

#### Do

Executes the given function with retry logic using default options.

```go
func Do(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError) error
```

#### DoWithOptions

Executes the given function with retry logic and custom options.

```go
func DoWithOptions(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError, options Options) error
```

## Examples

Currently, there are no dedicated examples for the retry package in the EXAMPLES directory. The following code snippets demonstrate common usage patterns:

### Basic Retry with Default Configuration

```go
// Define a function to retry
operation := func(ctx context.Context) error {
    // Your operation that might fail
    return callExternalService(ctx)
}

// Execute with default configuration
err := retry.Do(ctx, operation, retry.DefaultConfig(), nil)
if err != nil {
    // Handle error
}
```

### Custom Retry Configuration

```go
// Create a custom retry configuration
config := retry.DefaultConfig().
    WithMaxRetries(5).                      // Maximum 5 retry attempts
    WithInitialBackoff(200 * time.Millisecond). // Start with 200ms backoff
    WithMaxBackoff(5 * time.Second).        // Cap backoff at 5 seconds
    WithBackoffFactor(2.5).                 // Increase backoff by factor of 2.5
    WithJitterFactor(0.3)                   // Add 30% jitter

// Execute with custom configuration
err := retry.Do(ctx, operation, config, nil)
```

### Custom Error Handling

```go
// Define which errors should be retried
isRetryable := func(err error) bool {
    // Retry network errors, timeouts, and rate limit errors
    return errors.IsNetworkError(err) || 
           errors.IsTimeout(err) || 
           errors.Is(err, errors.New(errors.ResourceExhaustedCode, ""))
}

// Execute with custom error handling
err := retry.Do(ctx, operation, retry.DefaultConfig(), isRetryable)
```

### Retry with Telemetry and Logging

```go
// Create a logger
logger := logging.NewContextLogger(zapLogger)

// Create a tracer
tracer := telemetry.NewOtelTracer(otelTracer)

// Create options with logger and tracer
options := retry.DefaultOptions().
    WithLogger(logger).
    WithOtelTracer(tracer)

// Execute with custom options
err := retry.DoWithOptions(ctx, operation, retry.DefaultConfig(), isRetryable, options)
```

## Best Practices

1. **Choose Appropriate Retry Limits**: Set MaxRetries based on your operation's importance and expected recovery time
2. **Use Exponential Backoff**: The default BackoffFactor of 2.0 works well for most cases
3. **Always Add Jitter**: Use JitterFactor to prevent synchronized retries in distributed systems
4. **Be Selective About Retrying**: Only retry errors that are likely to be transient
5. **Respect Context Cancellation**: Always pass a proper context that can be cancelled or timed out

## Troubleshooting

### Common Issues

#### Too Many Retries

If your application is performing too many retries, consider:
- Reducing MaxRetries
- Increasing InitialBackoff
- Being more selective in your isRetryable function

#### Long Recovery Times

If recovery takes too long:
- Decrease BackoffFactor
- Decrease MaxBackoff
- Increase MaxRetries for critical operations

#### High Latency Spikes

If you're seeing latency spikes:
- Increase JitterFactor to spread out retry attempts
- Ensure your isRetryable function is correctly identifying transient errors

## Related Components

- [Errors](../errors/README.md) - Error handling and classification for retry operations
- [Circuit](../circuit/README.md) - Circuit breaker pattern for preventing cascading failures
- [Telemetry](../telemetry/README.md) - Telemetry integration for monitoring retry operations

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
