# Health Module

The Health Module provides functionality for implementing health check endpoints in Go applications. It helps you create health check handlers for monitoring the health of your services and their dependencies.

## Features

- **Health Status Reporting**: Report the health status of your application and its dependencies
- **Service Dependency Checks**: Check the health of service dependencies like databases
- **Configurable Responses**: Customize health check responses to include additional information
- **HTTP Integration**: Easy integration with HTTP servers
- **Kubernetes Compatibility**: Compatible with Kubernetes liveness and readiness probes
- **Tracing Integration**: OpenTelemetry tracing for all operations
- **Logging Integration**: Structured logging with zap

## Installation

```bash
go get github.com/abitofhelp/servicelib/health
```

## Quick Start

See the [Basic Usage example](../examples/health/basic_usage_example.go) for a complete, runnable example of how to use the Health module.

## API Documentation

### Health Handler

The `NewHandler` function creates a new health check HTTP handler that can be registered with an HTTP server.

#### Creating a Health Handler

See the [Basic Usage example](../examples/health/basic_usage_example.go) for a complete, runnable example of how to create a health handler.

#### Custom Health Status

See the [Custom Health Status example](../examples/health/custom_health_status_example.go) for a complete, runnable example of how to create a custom health status response.

### Health Status

The `HealthStatus` struct represents the health status response returned by the health handler.

```go
// Example of the HealthStatus struct
package example

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services,omitempty"`
}
```

### Status Constants

The health package provides constants for health status:

```go
// Example of the status constants
package example

// Status constants for health check responses
const (
	// StatusHealthy represents a healthy system status
	StatusHealthy = "Healthy"

	// StatusDegraded represents a degraded system status
	StatusDegraded = "Degraded"
)

// Service status constants for health check responses
const (
	// ServiceUp represents a service that is up and running
	ServiceUp = "Up"

	// ServiceDown represents a service that is down or unavailable
	ServiceDown = "Down"
)
```

## Health Check Best Practices

1. **Keep Checks Lightweight**: Health checks should be fast and not consume significant resources.

2. **Appropriate Timeouts**: Set appropriate timeouts for health checks to prevent hanging.

3. **Include Version Information**: Include version information in health status to help with debugging.

4. **Detailed Status**: Provide detailed status information for each dependency to aid in troubleshooting.

5. **Monitoring Integration**: Integrate health checks with your monitoring system.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
