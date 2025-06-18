# ServiceLib Documentation

Welcome to the ServiceLib documentation. This documentation provides comprehensive information about the ServiceLib library, its components, and how to use them effectively.

## Overview

ServiceLib is a comprehensive Go library designed to accelerate the development of robust, production-ready microservices. It provides a collection of reusable components and utilities that address common challenges in service development, allowing developers to focus on business logic rather than infrastructure concerns.

## Documentation Structure

- [Developer Guide](ServiceLib_Developer_Guide.md) - Comprehensive guide for developers using ServiceLib
- [API Reference](https://pkg.go.dev/github.com/abitofhelp/servicelib) - Generated API documentation
- [Examples](../examples/) - Example applications using ServiceLib

## Key Components

ServiceLib includes the following key components:

- **Authentication** - JWT, OAuth2, and OIDC implementations for secure service-to-service and user authentication
- **Configuration** - Flexible configuration management with adapters for various sources
- **Context** - Context utilities for request handling, cancellation, and value propagation
- **Database** - Database connection and transaction management
- **Dependency Injection** - Container-based DI system for managing service dependencies
- **Error Handling** - Structured error types and handling with rich context information
- **GraphQL** - Utilities for building GraphQL services
- **Health Checks** - Health check endpoints and handlers for Kubernetes readiness and liveness probes
- **Logging** - Structured logging with Zap
- **Middleware** - HTTP middleware components for common cross-cutting concerns
- **Repository Pattern** - Generic repository implementations for data access abstraction
- **Shutdown** - Graceful shutdown utilities for clean service termination
- **Signal Handling** - OS signal handling for responding to system events
- **Telemetry** - Metrics, tracing, and monitoring with Prometheus and OpenTelemetry
- **Validation** - Request and data validation

## Getting Started

To get started with ServiceLib, follow these steps:

1. Install the library:
   ```bash
   go get github.com/abitofhelp/servicelib
   ```

2. Import the packages you need:
   ```go
   import (
       "github.com/abitofhelp/servicelib/auth"
       "github.com/abitofhelp/servicelib/config"
       "github.com/abitofhelp/servicelib/logging"
       // Import other packages as needed
   )
   ```

3. Check the [Developer Guide](ServiceLib_Developer_Guide.md) for detailed usage instructions and examples.

## Contributing

Contributions to ServiceLib are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

ServiceLib is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.