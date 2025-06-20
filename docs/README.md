# ServiceLib Documentation

Welcome to the ServiceLib documentation. This documentation provides comprehensive information about the ServiceLib library, its components, and how to use them effectively.

## Overview

ServiceLib is a comprehensive Go library designed to accelerate the development of robust, production-ready microservices. It provides a collection of reusable components and utilities that address common challenges in service development, allowing developers to focus on business logic rather than infrastructure concerns.

## Documentation Structure

The ServiceLib documentation is organized into several key sections to help you find the information you need:

### Core Documentation

- **[Developer Guide](ServiceLib_Developer_Guide.md)** - Comprehensive guide for developers using ServiceLib, including detailed component descriptions, architecture overview, usage examples, best practices, and troubleshooting guidance.

- **[API Reference](https://pkg.go.dev/github.com/abitofhelp/servicelib)** - Generated API documentation from Go doc comments, providing function signatures, type definitions, package documentation, and usage examples.

- **[Examples](../examples/README.md)** - Example applications using ServiceLib, organized by component, including a quickstart example for new users, component-specific examples, advanced usage patterns, and integration examples.

### Development & Testing

- **[Integration Tests](Integration_Tests.md)** - Information about running integration tests, including test setup instructions, test environment requirements, running specific test suites, and troubleshooting test failures.

- **[Contributing Guide](../CONTRIBUTING.md)** - Guidelines for contributing to ServiceLib, including development workflow, coding standards, pull request process, and issue reporting.

### Architecture & Design

- **[UML Diagrams](diagrams/README.md)** - Architectural and component diagrams, including package dependencies, component relationships, sequence diagrams, and class diagrams.
  - [Package Dependencies Diagram](diagrams/svg/package_dependencies.svg) ([source](diagrams/source/package_dependencies.puml))
  - [Layered Architecture Diagram](diagrams/svg/layered_architecture.svg) ([source](diagrams/source/layered_architecture.puml))
  - [Authentication Component Diagram](diagrams/svg/auth_component.svg) ([source](diagrams/source/auth_component.puml))
  - [Dependency Injection Component Diagram](diagrams/svg/di_component.svg) ([source](diagrams/source/di_component.puml))
  - [Health Check Component Diagram](diagrams/svg/health_component.svg) ([source](diagrams/source/health_component.puml))

### Component Documentation

Each component in ServiceLib has its own README file with detailed documentation:

## Key Components

ServiceLib includes the following key components:

- **[Authentication](../auth/README.md)** - JWT, OAuth2, and OIDC implementations for secure service-to-service and user authentication
- **[Configuration](../config/README.md)** - Flexible configuration management with adapters for various sources
- **[Context](../context/README.md)** - Context utilities for request handling, cancellation, and value propagation
- **[Database](../db/README.md)** - Database connection and transaction management
- **[Dependency Injection](../di/README.md)** - Container-based DI system for managing service dependencies
- **[Error Handling](../errors/README.md)** - Structured error types and handling with rich context information
- **[GraphQL](../graphql/README.md)** - Utilities for building GraphQL services
- **[Health Checks](../health/README.md)** - Health check endpoints and handlers for Kubernetes readiness and liveness probes
- **[Logging](../logging/README.md)** - Structured logging with Zap
- **[Middleware](../middleware/README.md)** - HTTP middleware components for common cross-cutting concerns
- **[Repository Pattern](../repository/README.md)** - Generic repository implementations for data access abstraction
- **[Shutdown](../shutdown/README.md)** - Graceful shutdown utilities for clean service termination
- **[Signal Handling](../signal/README.md)** - OS signal handling for responding to system events
- **[Telemetry](../telemetry/README.md)** - Metrics, tracing, and monitoring with Prometheus and OpenTelemetry
- **[Validation](../validation/README.md)** - Request and data validation
- **[Value Objects](../valueobject/README.md)** - Immutable objects that represent domain concepts

## Getting Started

To get started with ServiceLib, follow these steps:

1. Install the library:
   ```bash
   go get github.com/abitofhelp/servicelib
   ```

2. Check out the [Quickstart Example](../examples/quickstart_example.go) to see a basic application using ServiceLib.

3. Review the [Developer Guide](ServiceLib_Developer_Guide.md) for detailed usage instructions and examples.

4. Explore the [Examples](../examples/) directory for component-specific examples.

5. For testing your implementation, refer to the [Integration Tests](Integration_Tests.md) documentation.

## Contributing

Contributions to ServiceLib are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

ServiceLib is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.
