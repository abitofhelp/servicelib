# Health

## Overview

The Health component provides functionality for health checking applications. It offers a robust HTTP handler for exposing health endpoints, checking the status of dependencies like databases, and returning standardized health status responses. This component is essential for implementing health checks in microservices and containerized applications.

## Features

- **HTTP Health Endpoint**: Ready-to-use HTTP handler for health check endpoints
- **Dependency Status Checking**: Check the health of dependencies like databases
- **Standardized Responses**: Consistent JSON response format for health status
- **Configurable Timeouts**: Set timeouts for health check operations
- **Status Reporting**: Report overall system status as Healthy or Degraded
- **Service-Level Reporting**: Individual status reporting for each dependency
- **Logging Integration**: Automatic logging of health check results

## Installation

```bash
go get github.com/abitofhelp/servicelib/health
```

## Quick Start

```go
package main

import (
    "net/http"
    
    "github.com/abitofhelp/servicelib/config"
    "github.com/abitofhelp/servicelib/health"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    
    // Create or get your application configuration
    appConfig := config.NewAppConfig()
    
    // Create a health check provider (typically your DI container)
    healthProvider := YourHealthCheckProvider()
    
    // Create a health check handler
    healthHandler := health.NewHandler(healthProvider, logger, appConfig)
    
    // Register the health check endpoint
    http.HandleFunc("/health", healthHandler)
    
    // Start the server
    http.ListenAndServe(":8080", nil)
}
```

## API Documentation

### Core Types

#### HealthStatus

Represents the health status response.

```go
type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp string            `json:"timestamp"`
    Version   string            `json:"version"`
    Services  map[string]string `json:"services,omitempty"`
}
```

### Key Interfaces

#### HealthCheckProvider

Combines the interfaces needed for health checking.

```go
type HealthCheckProvider interface {
    RepositoryFactoryProvider
}
```

#### RepositoryFactoryProvider

Defines the interface for getting a repository factory.

```go
type RepositoryFactoryProvider interface {
    GetRepositoryFactory() any
}
```

#### HealthConfig

Defines the interface for health check configuration.

```go
type HealthConfig interface {
    GetVersion() string
    GetName() string
    GetEnvironment() string
    GetTimeout() int
}
```

### Key Methods

#### NewHandler

Creates a new health check HTTP handler.

```go
func NewHandler(provider HealthCheckProvider, logger *zap.Logger, cfg config.Config) http.HandlerFunc
```

### Constants

#### Status Constants

```go
const (
    StatusHealthy = "Healthy"
    StatusDegraded = "Degraded"
)
```

#### Service Status Constants

```go
const (
    ServiceUp = "Up"
    ServiceDown = "Down"
)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Health Check](../EXAMPLES/health/basic_health_check/README.md) - Shows how to set up a basic health check endpoint
- [Custom Health Checks](../EXAMPLES/health/custom_health_checks/README.md) - Shows how to implement custom health checks
- [Kubernetes Integration](../EXAMPLES/health/kubernetes_integration/README.md) - Shows how to integrate health checks with Kubernetes

## Best Practices

1. **Separate Liveness and Readiness**: Consider implementing separate endpoints for liveness and readiness checks
2. **Keep Checks Lightweight**: Health checks should be fast and not consume significant resources
3. **Include Version Information**: Always include application version in health responses
4. **Set Appropriate Timeouts**: Configure timeouts based on expected check durations
5. **Log Health Status Changes**: Log when health status changes from healthy to degraded or vice versa

## Troubleshooting

### Common Issues

#### Health Check Timeouts

If health checks are timing out:
- Check if your database or other dependencies are responding slowly
- Consider increasing the timeout for health checks
- Ensure your health check logic is efficient

#### False Negatives

If health checks are reporting services as down when they're actually up:
- Verify your repository factory is properly initialized
- Check connection strings and credentials
- Ensure your health check logic is correctly implemented

## Related Components

- [Config](../config/README.md) - Configuration management used by health checks
- [Logging](../logging/README.md) - Logging integration for health check results
- [DB](../db/README.md) - Database utilities that can be health-checked
- [Middleware](../middleware/README.md) - HTTP middleware that can be used with health endpoints

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.