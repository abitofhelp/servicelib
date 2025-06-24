# Rate Limiter Package

The `rate` package provides a generic implementation of a token bucket rate limiter to protect resources from being overwhelmed by too many requests.

## Overview

Rate limiting is a strategy to control the rate of requests sent or received by a system. It helps prevent resource exhaustion, maintain service quality, and ensure fair usage of shared resources.

This implementation provides:

- Generic support for any function that returns a value and an error
- Token bucket algorithm for precise rate limiting
- Configurable requests per second and burst size
- Support for both immediate rejection and waiting for available tokens
- Support for OpenTelemetry tracing
- Fluent interface for configuration
- Thread-safe implementation

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/abitofhelp/servicelib/rate"
    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
)

// Create a rate limiter
cfg := rate.DefaultConfig().
    WithEnabled(true).
    WithRequestsPerSecond(100).
    WithBurstSize(50)

logger := logging.NewContextLogger(zap.NewNop())
options := rate.DefaultOptions().
    WithName("my-service").
    WithLogger(logger)

rl := rate.NewRateLimiter(cfg, options)

// Execute a function with rate limiting
result, err := rate.Execute(ctx, rl, "GetUserProfile", func(ctx context.Context) (UserProfile, error) {
    // Call external service
    return userService.GetProfile(ctx, userID)
})

// If the rate limit is exceeded, an error will be returned
if err != nil {
    // Handle rate limit exceeded error
    return nil, err
}

// Use the result
return result, nil
```

### With Wait

If you want to wait for a token to become available instead of immediately rejecting the request:

```go
// Execute a function with rate limiting and wait for a token
result, err := rate.ExecuteWithWait(ctx, rl, "GetUserProfile", func(ctx context.Context) (UserProfile, error) {
    // Call external service
    return userService.GetProfile(ctx, userID)
})

// If the context is canceled while waiting, an error will be returned
if err != nil {
    // Handle context cancellation
    return nil, err
}

// Use the result
return result, nil
```

## Configuration

The rate limiter can be configured using the `Config` struct and the fluent interface:

```go
cfg := rate.DefaultConfig().
    WithEnabled(true).                  // Enable/disable the rate limiter
    WithRequestsPerSecond(100).         // Number of requests allowed per second
    WithBurstSize(50)                   // Maximum number of requests allowed in a burst
```

## OpenTelemetry Integration

The rate limiter supports OpenTelemetry tracing:

```go
import (
    "go.opentelemetry.io/otel/trace"
)

// Create a tracer
tracer := otelTracer // Your OpenTelemetry tracer

// Configure the rate limiter with the tracer
options := rate.DefaultOptions().
    WithName("my-service").
    WithLogger(logger).
    WithOtelTracer(tracer)

rl := rate.NewRateLimiter(cfg, options)
```

## Thread Safety

The rate limiter is thread-safe and can be used concurrently from multiple goroutines.

## Performance Considerations

- The token bucket algorithm provides precise rate limiting with minimal overhead
- Tokens are refilled based on the time elapsed since the last refill, ensuring accurate rate limiting
- The rate limiter uses a mutex to ensure thread safety, so concurrent operations are serialized
- The `ExecuteWithWait` method uses a polling approach with a small sleep interval, which may not be optimal for high-concurrency scenarios