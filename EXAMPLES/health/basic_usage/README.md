# health basic_usage Example

## Overview

This example demonstrates how to set up a basic health check endpoint using the ServiceLib health package. It shows how to create a health check provider, configure it, and register it with an HTTP server.

## Features

- **Health Check Provider**: Implement a custom health check provider
- **Configuration Integration**: Integrate with the config package for application and database information
- **HTTP Endpoint**: Expose a health check endpoint via HTTP
- **Structured Health Information**: Return structured health information in the response

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Health Check Provider Implementation

The example implements a custom health check provider that satisfies the health.HealthCheckProvider interface:

```go
// MyHealthProvider implements the health.HealthCheckProvider interface
type MyHealthProvider struct{}

// GetRepositoryFactory returns a mock repository factory
func (p *MyHealthProvider) GetRepositoryFactory() any {
    // In a real application, this would return an actual repository factory
    // For this example, we'll just return a non-nil value
    return &struct{}{}
}
```

### Configuration Setup

The example sets up configuration objects that implement the required interfaces:

```go
// MyConfig implements the config.Config interface
type MyConfig struct{}

// GetApp returns the application configuration
func (c *MyConfig) GetApp() config.AppConfig {
    return &MyAppConfig{}
}

// GetDatabase returns the database configuration
func (c *MyConfig) GetDatabase() config.DatabaseConfig {
    return &MyDatabaseConfig{}
}
```

### Health Check Handler Creation and Registration

The example creates a health check handler and registers it with an HTTP server:

```go
// Create a health check provider
provider := &MyHealthProvider{}

// Create a configuration
cfg := &MyConfig{}

// Create a health check handler
healthHandler := health.NewHandler(provider, logger, cfg)

// Register the health check handler
http.HandleFunc("/health", healthHandler)
```

## Expected Output

When you run the example and access the health check endpoint at http://localhost:8080/health, you should see a JSON response similar to:

```json
{
  "status": "UP",
  "version": "1.0.0",
  "name": "my-service",
  "environment": "development",
  "database": {
    "status": "UP",
    "type": "postgres",
    "name": "mydb"
  }
}
```

Note: The example doesn't actually start the HTTP server, but it shows how to set up the health check endpoint.

## Related Examples


- [custom_health_status](../custom_health_status/README.md) - Related example for custom_health_status

## Related Components

- [health Package](../../../health/README.md) - The health package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
