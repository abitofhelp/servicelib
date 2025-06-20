# Middleware Module

The Middleware Module provides HTTP middleware components for common cross-cutting concerns in Go applications. It includes middleware for request context, timeout handling, panic recovery, logging, error handling, and CORS support.

## Features

- **Request Context**: Add request ID and other context information to requests
- **Timeout Handling**: Set timeouts for request processing
- **Recovery**: Catch and handle panics in HTTP handlers
- **Logging**: Log request and response details
- **Error Handling**: Centralized error handling for HTTP responses
- **CORS**: Cross-Origin Resource Sharing support

## Installation

```bash
go get github.com/abitofhelp/servicelib/middleware
```

## Quick Start

See the [Basic Usage example](../examples/middleware/basic_usage_example.go) for a complete, runnable example of how to use the Middleware module.

## API Documentation

### Basic Usage

See the [Basic Usage example](../examples/middleware/basic_usage_example.go) for a complete, runnable example of how to use the middleware components together.

### Request Context Middleware

The `WithRequestContext` middleware adds request context information, including request ID and start time.

See the [Basic Usage example](../examples/middleware/basic_usage_example.go) for a complete, runnable example of how to use the request context middleware.

### Timeout Middleware

The `WithTimeout` middleware sets a timeout for request processing.

See the [Timeout example](../examples/middleware/timeout_example.go) for a complete, runnable example of how to use the timeout middleware.

### Recovery Middleware

The `WithRecovery` middleware catches panics in HTTP handlers and converts them to error responses.

See the [Recovery example](../examples/middleware/recovery_example.go) for a complete, runnable example of how to use the recovery middleware.

### Logging Middleware

The `WithLogging` middleware logs request and response details.

See the [Logging example](../examples/middleware/logging_example.go) for a complete, runnable example of how to use the logging middleware.

### Error Handling Middleware

The `WithErrorHandling` middleware provides centralized error handling for HTTP responses.

See the [Error Handling example](../examples/middleware/error_handling_example.go) for a complete, runnable example of how to use the error handling middleware.

### CORS Middleware

The `WithCORS` middleware adds Cross-Origin Resource Sharing support.

See the [CORS example](../examples/middleware/cors_example.go) for a complete, runnable example of how to use the CORS middleware.

### Middleware Composition

```go
// Example of middleware composition
package example

import (
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/middleware"
	"go.uber.org/zap"
)

func composeMiddleware(logger *zap.Logger) http.Handler {
	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Create a simple HTTP handler
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
	})

	// Apply middleware in the correct order
	handler := middleware.WithRequestContext(helloHandler)
	handler = middleware.WithRecovery(contextLogger, handler)
	handler = middleware.WithLogging(contextLogger, handler)
	handler = middleware.WithErrorHandling(handler)
	handler = middleware.WithCORS(handler)

	return handler
}
```

### Creating Custom Middleware

```go
// Example of creating a custom middleware
package example

import (
	"net/http"
	"time"
)

// TimingMiddleware is a custom middleware that measures request duration
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record start time
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Calculate request duration
		duration := time.Since(start)

		// Add a custom header with the request duration
		w.Header().Set("X-Request-Duration", duration.String())
	})
}

// Usage example
func useTimingMiddleware() {
	// Create a simple HTTP handler
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Apply your custom middleware
	handler := TimingMiddleware(helloHandler)

	// Register the handler
	http.Handle("/", handler)
}
```

## Middleware Best Practices

1. **Order Matters**: The order of middleware in the chain is important. For example, recovery middleware should be first to catch panics in other middleware.

   ```go
   // Example of middleware ordering
   package example

   import (
       "net/http"

       "github.com/abitofhelp/servicelib/logging"
       "github.com/abitofhelp/servicelib/middleware"
   )

   func applyMiddlewareInOrder(logger *logging.ContextLogger, handler http.Handler) http.Handler {
       // Apply middleware in the correct order
       handler = middleware.WithRecovery(logger, handler)     // First to catch panics
       handler = middleware.WithRequestContext(handler)       // Second to add request context
       handler = middleware.WithLogging(logger, handler)      // Third to log with request context
       handler = middleware.WithErrorHandling(handler)        // Fourth to handle errors
       handler = middleware.WithCORS(handler)                 // Fifth to add CORS headers

       return handler
   }
   ```

2. **Performance**: Be mindful of middleware performance, especially for high-traffic services.

3. **Error Handling**: Ensure middleware properly handles errors and doesn't swallow them.

4. **Context Values**: Use context to pass values between middleware and handlers.

5. **Middleware Configuration**: Use configuration options to customize middleware behavior.

6. **Testing**: Test middleware in isolation and as part of the chain.

7. **Logging**: Include relevant information in logs, but be careful not to log sensitive data.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
