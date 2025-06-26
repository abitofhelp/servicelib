# Rate

## Overview

The Rate component provides functionality for rate limiting to protect resources from being overwhelmed by too many requests. It implements a token bucket rate limiter that allows a configurable number of requests per second with a configurable burst size.

## Features

- **Token Bucket Algorithm**: Implements the token bucket algorithm for efficient rate limiting
- **Configurable Rate Limits**: Set requests per second and burst size to match your application needs
- **Immediate Rejection**: Reject requests immediately when rate limit is exceeded
- **Wait Mode**: Optionally wait for tokens to become available instead of rejecting requests
- **Telemetry Integration**: Built-in OpenTelemetry tracing for monitoring rate limiter operations
- **Context Awareness**: Respects context cancellation for graceful shutdowns
- **Generic Functions**: Type-safe execution with Go generics

## Installation

```bash
go get github.com/abitofhelp/servicelib/rate
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/abitofhelp/servicelib/rate"
)

func main() {
    // Create a rate limiter with default configuration
    cfg := rate.DefaultConfig()
    options := rate.DefaultOptions().WithName("api-limiter")
    limiter := rate.NewRateLimiter(cfg, options)
    
    // Execute a function with rate limiting
    ctx := context.Background()
    result, err := rate.Execute(ctx, limiter, "get-user", func(ctx context.Context) (string, error) {
        // Your rate-limited operation here
        return "user data", nil
    })
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Result: %s\n", result)
}
```

## API Documentation

### Core Types

#### Config

Configuration for the rate limiter.

```go
type Config struct {
    // Enabled determines if the rate limiter is enabled
    Enabled bool
    // RequestsPerSecond is the number of requests allowed per second
    RequestsPerSecond int
    // BurstSize is the maximum number of requests allowed in a burst
    BurstSize int
}
```

#### Options

Additional options for the rate limiter.

```go
type Options struct {
    // Logger is used for logging rate limiter operations
    Logger *logging.ContextLogger
    // Tracer is used for tracing rate limiter operations
    Tracer telemetry.Tracer
    // Name is the name of the rate limiter
    Name string
}
```

#### RateLimiter

The main rate limiter struct.

```go
type RateLimiter struct {
    // Internal fields
}
```

### Key Methods

#### NewRateLimiter

Creates a new rate limiter.

```go
func NewRateLimiter(config Config, options Options) *RateLimiter
```

#### Allow

Checks if a request should be allowed based on the rate limit.

```go
func (rl *RateLimiter) Allow() bool
```

#### Execute

Executes a function with rate limiting, returning an error if the rate limit is exceeded.

```go
func Execute[T any](ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (T, error)) (T, error)
```

#### ExecuteWithWait

Executes a function with rate limiting, waiting for a token to become available if necessary.

```go
func ExecuteWithWait[T any](ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (T, error)) (T, error)
```

#### Reset

Resets the rate limiter to its initial state.

```go
func (rl *RateLimiter) Reset()
```

## Examples

Currently, there are no dedicated examples for the rate package in the EXAMPLES directory. The following code snippets demonstrate common usage patterns:

### Basic Rate Limiting

```go
// Create a rate limiter that allows 100 requests per second with a burst of 50
cfg := rate.DefaultConfig().
    WithRequestsPerSecond(100).
    WithBurstSize(50)
options := rate.DefaultOptions().WithName("api-limiter")
limiter := rate.NewRateLimiter(cfg, options)

// Check if a request is allowed
if limiter.Allow() {
    // Process the request
} else {
    // Return rate limit exceeded error
}
```

### Rate Limited Function Execution

```go
// Create a rate limiter
limiter := rate.NewRateLimiter(rate.DefaultConfig(), rate.DefaultOptions())

// Execute a function with rate limiting
ctx := context.Background()
result, err := rate.Execute(ctx, limiter, "get-data", func(ctx context.Context) ([]byte, error) {
    // Make an API call or perform other rate-limited operation
    return fetchData(ctx)
})

if err != nil {
    // Handle error (including rate limit exceeded)
    return err
}

// Use the result
processData(result)
```

### Waiting for Rate Limit

```go
// Create a rate limiter
limiter := rate.NewRateLimiter(rate.DefaultConfig(), rate.DefaultOptions())

// Execute a function with waiting for rate limit
ctx := context.Background()
result, err := rate.ExecuteWithWait(ctx, limiter, "process-item", func(ctx context.Context) (bool, error) {
    // This function will be executed once a token is available
    return processItem(ctx, item)
})

if err != nil {
    // Handle error (context cancellation, etc.)
    return err
}

// Use the result
if result {
    fmt.Println("Item processed successfully")
}
```

## Best Practices

1. **Choose Appropriate Limits**: Set rate limits based on your resource capacity and expected load
2. **Use Descriptive Operation Names**: Provide meaningful operation names for better telemetry and logging
3. **Handle Rate Limit Errors**: Properly handle rate limit exceeded errors in your application
4. **Consider Wait vs. Reject**: Use ExecuteWithWait when appropriate to smooth out traffic spikes
5. **Monitor Rate Limiter Metrics**: Use the integrated telemetry to monitor rate limiter performance

## Troubleshooting

### Common Issues

#### Rate Limiter Too Restrictive

If your rate limiter is rejecting too many requests, consider increasing the RequestsPerSecond or BurstSize parameters.

```go
cfg := rate.DefaultConfig().
    WithRequestsPerSecond(200).  // Increase from default 100
    WithBurstSize(100)           // Increase from default 50
```

#### High Latency with ExecuteWithWait

If you're experiencing high latency with ExecuteWithWait, it might indicate that your rate limits are too low for your traffic. Consider increasing the limits or adding more capacity.

## Related Components

- [Errors](../errors/README.md) - Error handling for rate limiting errors
- [Telemetry](../telemetry/README.md) - Telemetry integration for rate limiter monitoring
- [Logging](../logging/README.md) - Logging for rate limiter events

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
