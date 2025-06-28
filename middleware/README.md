# Middleware

## Overview

The Middleware component provides a collection of HTTP middleware for Go web applications. It offers ready-to-use middleware for common tasks like request logging, error handling, timeout management, and CORS support. The component also includes utilities for middleware chaining and context management.

## Features

- **Request Context**: Add request IDs and timing information to request contexts
- **Logging**: Log request details including method, path, status code, and duration
- **Error Handling**: Map errors to appropriate HTTP responses with status codes
- **Panic Recovery**: Catch and handle panics to prevent application crashes
- **Timeout Management**: Add request timeouts with proper cancellation handling
- **CORS Support**: Add Cross-Origin Resource Sharing headers to responses
- **Context Cancellation**: Detect and handle client disconnections
- **Middleware Chaining**: Apply multiple middleware in a specific order
- **Thread Safety**: Thread-safe response writing for concurrent operations

## Installation

```bash
go get github.com/abitofhelp/servicelib/middleware
```

## Quick Start

```go
package main

import (
    "net/http"
    
    "github.com/abitofhelp/servicelib/middleware"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    
    // Create a simple handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Apply middleware
    wrappedHandler := middleware.ApplyMiddleware(handler, logger)
    
    // Start the server
    http.ListenAndServe(":8080", wrappedHandler)
}
```

## API Documentation

### Core Types

#### Middleware

Represents a middleware function that wraps an http.Handler.

```go
type Middleware func(http.Handler) http.Handler
```

### Key Functions

#### Chain

Applies multiple middleware to an http.Handler in the order they are provided.

```go
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler
```

#### ApplyMiddleware

Applies all middleware to a handler in the recommended order.

```go
func ApplyMiddleware(handler http.Handler, logger *zap.Logger) http.Handler
```

### Individual Middleware

#### WithRequestContext

Adds request context information to the request.

```go
func WithRequestContext(next http.Handler) http.Handler
```

#### WithLogging

Adds request logging.

```go
func WithLogging(logger *logging.ContextLogger, next http.Handler) http.Handler
```

#### WithRecovery

Adds panic recovery to the request.

```go
func WithRecovery(logger *logging.ContextLogger, next http.Handler) http.Handler
```

#### WithTimeout

Adds a timeout to the request context.

```go
func WithTimeout(timeout time.Duration) func(http.Handler) http.Handler
```

#### WithErrorHandling

Adds centralized error handling.

```go
func WithErrorHandling(next http.Handler) http.Handler
```

#### WithCORS

Adds CORS headers to allow cross-origin requests.

```go
func WithCORS(next http.Handler) http.Handler
```

#### WithContextCancellation

Checks for context cancellation and detects when a client disconnects.

```go
func WithContextCancellation(logger *logging.ContextLogger) func(next http.Handler) http.Handler
```

### Context Utilities

#### RequestID

Returns the request ID from the context.

```go
func RequestID(ctx context.Context) string
```

#### StartTime

Returns the request start time from the context.

```go
func StartTime(ctx context.Context) time.Time
```

#### RequestDuration

Returns the duration since the request started.

```go
func RequestDuration(ctx context.Context) time.Duration
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/middleware/basic_usage/README.md) - Shows how to use basic middleware
- [Custom Middleware](../EXAMPLES/middleware/custom_middleware/README.md) - Shows how to create custom middleware
- [Middleware Chain](../EXAMPLES/middleware/middleware_chain/README.md) - Shows how to chain multiple middleware
- [Error Handling](../EXAMPLES/middleware/error_handling/README.md) - Shows how to use error handling middleware

## Best Practices

1. **Apply Middleware in the Right Order**: The order of middleware matters; typically, request context should be first, followed by recovery, logging, etc.
2. **Use Chain for Custom Middleware Combinations**: Use the Chain function when you need a specific combination of middleware
3. **Handle Context Cancellation**: Always handle context cancellation to avoid resource leaks
4. **Set Appropriate Timeouts**: Configure timeouts based on expected operation durations
5. **Log Request Details**: Use the logging middleware to track request details for debugging and monitoring

## Troubleshooting

### Common Issues

#### Middleware Order Problems

If middleware isn't working as expected:
- Check the order in which middleware is applied
- Remember that middleware is executed in reverse order of application (last applied is first executed)
- Use the ApplyMiddleware function for a standard ordering

#### Response Already Written

If you're seeing "http: multiple response.WriteHeader calls":
- Ensure only one middleware is writing the response
- Use the error handling middleware to centralize error responses
- Check for race conditions in concurrent handlers

## Related Components

- [Logging](../logging/README.md) - Logging integration for middleware
- [Errors](../errors/README.md) - Error handling used by middleware
- [Context](../context/README.md) - Context utilities used by middleware
- [Health](../health/README.md) - Health checks that can use middleware

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.