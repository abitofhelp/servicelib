# Health Package Examples

This directory contains examples demonstrating how to use the `health` package, which provides utilities for implementing health check endpoints in Go applications. Health checks are essential for monitoring service health, readiness, and liveness in microservice architectures.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates the basic setup and usage of the health package.

Key concepts:
- Implementing the health.HealthCheckProvider interface
- Setting up configuration objects
- Creating a health check handler
- Registering the health check handler with an HTTP server
- Exposing a health check endpoint

### 2. Custom Health Status Example

[custom_health_status_example.go](custom_health_status_example.go)

Shows how to implement custom health status checks.

Key concepts:
- Creating custom health status indicators
- Implementing custom health check logic
- Reporting detailed health information
- Handling different health status scenarios
- Customizing health check responses

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/health/basic_usage_example.go
```

## Additional Resources

For more information about the health package, see the [health package documentation](../../health/README.md).