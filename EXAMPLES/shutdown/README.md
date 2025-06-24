# Shutdown Package Examples

This directory contains examples demonstrating how to use the `shutdown` package, which provides utilities for implementing graceful shutdown in Go applications. Graceful shutdown is essential for properly releasing resources, completing in-flight requests, and ensuring data integrity when an application is being terminated.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates the basic implementation of graceful shutdown for an HTTP server.

Key concepts:
- Creating an HTTP server
- Defining a shutdown function
- Using GracefulShutdown to wait for shutdown signals
- Handling shutdown timeouts
- Logging shutdown events

### 2. Multiple Resource Example

[multiple_resource_example.go](multiple_resource_example.go)

Shows how to gracefully shut down multiple resources in a specific order.

Key concepts:
- Managing shutdown of multiple resources
- Controlling shutdown order
- Handling errors from multiple shutdown operations
- Coordinating dependent resource shutdown
- Setting appropriate timeouts for different resources

### 3. Programmatic Shutdown Example

[programmatic_shutdown_example.go](programmatic_shutdown_example.go)

Demonstrates how to trigger shutdown programmatically rather than waiting for OS signals.

Key concepts:
- Triggering shutdown from application code
- Combining signal-based and programmatic shutdown
- Implementing health-based shutdown triggers
- Handling shutdown requests from management endpoints
- Coordinating shutdown across multiple services

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/shutdown/basic_usage_example.go
```

## Additional Resources

For more information about the shutdown package, see the [shutdown package documentation](../../shutdown/README.md).