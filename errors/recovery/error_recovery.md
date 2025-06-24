# Error Recovery and Rate Limiting

## Overview
The error recovery package provides robust error handling mechanisms including panic recovery, error rate limiting, and circuit breaker patterns.

## Features

### Panic Recovery
The `RecoveryHandler` provides a safe way to execute operations with automatic panic recovery:

```go
handler := recovery.NewRecoveryHandler(logger, maxErrorRate)
err := handler.WithRecovery(ctx, "operation-name", func() error {
    // Your operation here
    return nil
})
```

### Error Rate Limiting
Built-in error rate limiting prevents log flooding and system overload:

```go
// Limit to 10 errors per second
handler := recovery.NewRecoveryHandler(logger, 10)
```

### Circuit Breaker
The circuit breaker pattern protects services from cascading failures:

```go
// Opens after 5 failures, resets after 1 minute
cb := recovery.NewCircuitBreaker(5, time.Minute)

err := cb.Execute(ctx, func() error {
    return someOperation()
})
```

States:
- Closed: Normal operation
- Open: All requests fail fast
- Half-Open: Testing if service has recovered

## Context Handling

The recovery package properly handles context cancellation and timeouts:

```go
// Context cancellation
ctx, cancel := context.WithCancel(context.Background())
cancel()
err := cb.Execute(ctx, func() error {
    // This won't execute if context is canceled
    return nil
})
// err will be ErrContextCanceled

// Context timeout
ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()
time.Sleep(200 * time.Millisecond) // Simulate timeout
err = cb.Execute(ctx, func() error {
    // This won't execute if context deadline exceeded
    return nil
})
// err will be ErrContextDeadlineExceeded
```

## Best Practices

1. Always wrap critical operations with panic recovery
2. Set appropriate error rate limits based on operation type
3. Use circuit breakers for external service calls
4. Include meaningful operation names for better tracking
5. Configure appropriate timeouts for circuit breaker reset
6. Always pass context as the first parameter to functions
7. Use ContextLogger for logging to include trace information

## Example

```go
func main() {
    ctx := context.Background()
    logger := logging.NewContextLogger(zapLogger)
    handler := recovery.NewRecoveryHandler(logger, 10)
    cb := recovery.NewCircuitBreaker(5, time.Minute)

    err := handler.WithRecovery(ctx, "external-api-call", func() error {
        return cb.Execute(ctx, func() error {
            return callExternalAPI()
        })
    })
}
```

For a complete example, see [Circuit Breaker Example](../examples/recovery/circuit_breaker_example.go).