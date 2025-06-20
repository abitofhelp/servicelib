# Middleware Package Examples

This directory contains examples demonstrating how to use the `middleware` package, which provides HTTP middleware components for Go web applications. The package offers various middleware functions that can be chained together to add functionality such as logging, error handling, CORS support, and more to HTTP handlers.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates the basic setup and usage of the middleware package.

Key concepts:
- Creating a simple HTTP handler
- Applying multiple middleware components
- Chaining middleware together
- Registering the handler with the HTTP server
- Understanding middleware execution order

### 2. CORS Example

[cors_example.go](cors_example.go)

Shows how to use Cross-Origin Resource Sharing (CORS) middleware.

Key concepts:
- Setting up CORS middleware
- Configuring allowed origins, methods, and headers
- Handling preflight requests
- Securing API endpoints for cross-origin requests
- Testing CORS functionality

### 3. Error Handling Example

[error_handling_example.go](error_handling_example.go)

Demonstrates how to use error handling middleware.

Key concepts:
- Consistent error handling across endpoints
- Converting errors to appropriate HTTP responses
- Handling different error types
- Customizing error responses
- Maintaining clean error handling logic

### 4. Logging Example

[logging_example.go](logging_example.go)

Shows how to use logging middleware for HTTP requests.

Key concepts:
- Logging HTTP request details
- Recording response status and timing
- Structured logging for requests
- Correlating logs with request IDs
- Configuring log levels for requests

### 5. Recovery Example

[recovery_example.go](recovery_example.go)

Demonstrates how to use recovery middleware to handle panics.

Key concepts:
- Recovering from panics in HTTP handlers
- Preventing server crashes
- Logging panic information
- Returning appropriate error responses
- Maintaining application stability

### 6. Timeout Example

[timeout_example.go](timeout_example.go)

Shows how to use timeout middleware to limit request processing time.

Key concepts:
- Setting request timeouts
- Handling long-running requests
- Canceling context on timeout
- Preventing resource exhaustion
- Improving service reliability

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/middleware/basic_usage_example.go
```

## Additional Resources

For more information about the middleware package, see the [middleware package documentation](../../middleware/README.md).