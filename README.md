# ServiceLib

[![codecov](https://codecov.io/gh/abitofhelp/servicelib/graph/badge.svg)](https://codecov.io/gh/abitofhelp/servicelib)
[![Go Report Card](https://goreportcard.com/badge/github.com/abitofhelp/servicelib)](https://goreportcard.com/report/github.com/abitofhelp/servicelib)
[![GoDoc](https://godoc.org/github.com/abitofhelp/servicelib?status.svg)](https://godoc.org/github.com/abitofhelp/servicelib)

## Coverage (click for details)
[![codecov](https://codecov.io/gh/abitofhelp/servicelib/graphs/sunburst.svg)](https://codecov.io/gh/abitofhelp/servicelib)

## Overview
### Note: The Family-Service repository implements a GraphQL service using ServiceLib. It serves as an excellent illustration of ServiceLib's capabilities and potential.

ServiceLib is a comprehensive Go library designed to accelerate the development of robust, production-ready microservices. It provides a collection of reusable components and utilities that address common challenges in service development, allowing developers to focus on business logic rather than infrastructure concerns.

The library follows modern Go practices and design patterns, with a focus on:

- **Modularity**: Each component can be used independently or together with others
- **Testability**: All components are designed with testing in mind
- **Performance**: Optimized for high-throughput microservices
- **Reliability**: Built-in error handling and recovery mechanisms
- **Observability**: Integrated logging, metrics, and tracing

## Documentation

For comprehensive documentation, please see the following resources:

- **[Developer Guide](docs/ServiceLib_Developer_Guide.md)** - Detailed component descriptions, architecture diagrams, usage examples, best practices, and troubleshooting guidance
- **[Integration Tests](docs/Integration_Tests.md)** - Instructions for running integration tests, test environment setup, and troubleshooting test failures
- **[UML Diagrams](docs/diagrams/README.md)** - Architectural and component diagrams showing package dependencies and component relationships
- **[API Reference](https://pkg.go.dev/github.com/abitofhelp/servicelib)** - Generated API documentation from Go doc comments

## Features

ServiceLib provides a comprehensive set of components organized by functionality:

### Core Infrastructure

- **[Configuration](config/README.md)** - Flexible configuration management with adapters for various sources (files, environment variables, etc.)
- **[ConfigScan](configscan/README.md)** - Tool for scanning packages to determine whether any have configuration requirements and don't have default values set
- **[Context](context/README.md)** - Context utilities for request handling, cancellation, and value propagation
- **[Dependency Injection](di/README.md)** - Container-based DI system for managing service dependencies
- **[Environment Variables](env/README.md)** - Utilities for working with environment variables with fallback values
- **[Error Handling](errors/README.md)** - Structured error types and handling with rich context information
- **[Logging](logging/README.md)** - Structured logging with Zap for high-performance logging
- **[Retry](retry/README.md)** - Retry operations with configurable backoff and jitter

### Security

- **[Authentication](auth/README.md)** - JWT, OAuth2, and OIDC implementations for secure service-to-service and user authentication
- **[Validation](validation/README.md)** - Request and data validation using go-playground/validator

### Data Access

- **[Database](db/README.md)** - Database connection and transaction management with support for PostgreSQL, SQLite, and MongoDB
- **[Repository Pattern](repository/README.md)** - Generic repository implementations for data access abstraction
- **[Transaction](transaction/README.md)** - Distributed transaction management with Saga pattern support

### API & Communication

- **[GraphQL](graphql/README.md)** - Utilities for building GraphQL services with gqlgen integration
- **[Middleware](middleware/README.md)** - HTTP middleware components for common cross-cutting concerns

### Operations & Observability

- **[Health Checks](health/README.md)** - Health check endpoints and handlers for Kubernetes readiness and liveness probes
- **[Shutdown](shutdown/README.md)** - Graceful shutdown utilities for clean service termination
- **[Signal Handling](signal/README.md)** - OS signal handling for responding to system events
- **[Telemetry](telemetry/README.md)** - Metrics, tracing, and monitoring with Prometheus and OpenTelemetry

### Domain Modeling

- **[Value Objects](valueobject/README.md)** - Immutable objects that represent domain concepts

## Compatibility and Versioning

ServiceLib follows [Semantic Versioning](https://semver.org/). The current version is available in the [go.mod](go.mod) file.

- **Go Version**: Requires Go 1.24 or later
- **Dependencies**: All dependencies are managed through Go modules
- **API Stability**: APIs marked as stable will not have breaking changes within the same major version
- **Backward Compatibility**: We strive to maintain backward compatibility within the same major version

## Installation

```bash
go get github.com/abitofhelp/servicelib
```

## Getting Started

> **Important**: In order to ensure that developers can build and work with this package within different IDEs and environments, please use the Makefile to build, test, etc.

To get started with ServiceLib, you can create a simple HTTP server with logging and middleware:

1. Create a new Go module:
   ```bash
   mkdir myservice
   cd myservice
   go mod init myservice
   ```

2. Install ServiceLib:
   ```bash
   go get github.com/abitofhelp/servicelib
   ```

3. Create a main.go file with a basic HTTP server using ServiceLib components.

For a complete example, see the [Quickstart Example](examples/quickstart_example.go).

## Examples

Each package in ServiceLib includes its own README.md with detailed documentation and examples:

- [Authentication](auth/README.md) - JWT, OAuth2, and OIDC examples
- [Configuration](config/README.md) - Configuration management examples
- [ConfigScan](configscan/README.md) - Package scanning for configuration requirements
- [Database](db/README.md) - Database connection and transaction examples
- [Dependency Injection](di/README.md) - DI container usage examples
- [Health Checks](health/README.md) - Health check implementation examples
- [Logging](logging/README.md) - Structured logging examples
- [Telemetry](telemetry/README.md) - Metrics and tracing examples
- [Transaction](transaction/README.md) - Distributed transaction examples

For complete example applications, see the [Examples Directory](examples/README.md).

For comprehensive documentation, please see the [Developer Guide](docs/ServiceLib_Developer_Guide.md).

## Component Documentation

### Authentication

The `auth` package provides implementations for JWT, OAuth2, and OIDC authentication methods:

- **JWT**: JSON Web Token implementation for stateless authentication
  - Token generation and validation
  - Support for standard claims and custom claims
  - Configurable signing methods (HMAC, RSA, ECDSA)

- **OAuth2**: OAuth 2.0 client implementation
  - Authorization Code, Client Credentials, and Password grant types
  - Token refresh and validation
  - State management for CSRF protection

- **OIDC**: OpenID Connect implementation
  - Discovery of provider configuration
  - ID token validation
  - User info retrieval

### Configuration

The `config` package provides a flexible configuration system that supports multiple sources and formats:

- **Multiple Sources**:
  - YAML and JSON files
  - Environment variables
  - Command-line flags
  - In-memory values

- **Features**:
  - Hierarchical configuration with dot notation
  - Default values
  - Type conversion
  - Configuration reloading
  - Validation

- **Adapters**: Easily create custom adapters for different configuration sources

### ConfigScan

The `configscan` package provides tools for scanning packages to determine whether any have configuration requirements and don't have default values set:

- **Scanner**: Scans packages for configuration requirements
  - Identifies structs with names containing "Config"
  - Checks if there's a default function for each config struct
  - Reports config structs that don't have default functions

- **Command-Line Tool**: Easily scan your project from the command line
  - Scan the current directory or a specific directory
  - Get a report of packages with missing default values

- **Best Practices**: Encourages providing default values for all configuration structs
  - Ensures users don't have to specify every configuration parameter
  - Promotes consistent configuration patterns across packages

### Context

The `context` package extends Go's standard context package with additional utilities:

- **Value Management**: Strongly typed context values
- **Timeout Management**: Utilities for working with context deadlines
- **Cancellation**: Simplified cancellation patterns
- **Propagation**: Utilities for propagating context values across service boundaries

### Database

The `db` package provides utilities for database connection management and operations:

- **Connection Management**:
  - Connection pooling
  - Automatic reconnection
  - Health checks

- **Supported Databases**:
  - PostgreSQL (via pgx)
  - SQLite
  - MongoDB

- **Features**:
  - Transaction management
  - Query execution with retries
  - Result mapping
  - Migrations

### Date

The `date` package provides utilities for working with dates and times:

- **Formatting**: Consistent date/time formatting
- **Parsing**: Robust date/time parsing with error handling
- **Comparison**: Date comparison utilities
- **Timezone**: Timezone conversion and management

### Dependency Injection

The `di` package provides a container-based dependency injection system:

- **Container Types**:
  - Base container
  - Service container
  - Repository container
  - Generic container

- **Features**:
  - Constructor injection
  - Singleton instances
  - Lazy initialization
  - Scoped instances
  - Circular dependency detection

### Error Handling

The `errors` package provides structured error types and handling:

- **Error Types**:
  - Domain errors
  - Infrastructure errors
  - Application errors
  - Validation errors

- **Features**:
  - Error wrapping with context
  - Error codes
  - Localized error messages
  - Stack traces
  - Error categorization

### GraphQL

The `graphql` package provides utilities for building GraphQL services:

- **Integration**: Integration with gqlgen
- **Error Handling**: Structured error handling for GraphQL
- **Middleware**: GraphQL-specific middleware
- **Resolvers**: Utilities for implementing resolvers

### Health Checks

The `health` package provides components for implementing health check endpoints:

- **Check Types**:
  - Liveness checks
  - Readiness checks
  - Startup checks

- **Features**:
  - Configurable check intervals
  - Automatic registration with HTTP server
  - Detailed health status reporting
  - Integration with Kubernetes probes

### Logging

The `logging` package provides structured logging with Zap:

- **Log Levels**: Debug, Info, Warn, Error, Fatal
- **Structured Logging**: Key-value pairs for better searchability
- **Output Formats**: JSON, console
- **Integration**: Context-aware logging
- **Performance**: High-performance logging with minimal allocations

### Middleware

The `middleware` package provides HTTP middleware components:

- **Authentication**: JWT authentication middleware
- **Logging**: Request/response logging
- **Metrics**: Request metrics collection
- **Tracing**: Distributed tracing
- **Recovery**: Panic recovery
- **CORS**: Cross-Origin Resource Sharing
- **Rate Limiting**: Request rate limiting

### Repository Pattern

The `repository` package provides generic repository implementations:

- **Generic Repository**: Type-safe repository implementation
- **CRUD Operations**: Create, Read, Update, Delete
- **Query Building**: Fluent query building
- **Pagination**: Offset and cursor-based pagination
- **Sorting**: Multi-field sorting

### Shutdown

The `shutdown` package provides graceful shutdown utilities:

- **Graceful Shutdown**: Orderly shutdown of services
- **Timeout Management**: Configurable shutdown timeouts
- **Dependency Order**: Shutdown in the correct order
- **Resource Cleanup**: Ensure all resources are properly released

### Signal Handling

The `signal` package provides OS signal handling:

- **Signal Types**: SIGINT, SIGTERM, SIGHUP
- **Custom Handlers**: Register custom signal handlers
- **Graceful Shutdown**: Integration with shutdown package

### String Utilities

The `stringutil` package provides string manipulation utilities:

- **Formatting**: String formatting utilities
- **Validation**: String validation
- **Transformation**: Case conversion, trimming, etc.
- **Generation**: Random string generation

### Telemetry

The `telemetry` package provides utilities for metrics, tracing, and monitoring:

- **Metrics**:
  - Prometheus integration
  - Counter, gauge, histogram, and summary metrics
  - Default metrics for HTTP, gRPC, and database

- **Tracing**:
  - OpenTelemetry integration
  - Span creation and management
  - Context propagation
  - Automatic instrumentation for HTTP and gRPC

- **Monitoring**:
  - Health check integration
  - Alerting utilities
  - Dashboard templates

### Transaction

The `transaction` package provides utilities for managing distributed transactions:

- **Saga Pattern**: Implementation of the Saga pattern for distributed transactions
- **Compensation**: Transaction compensation for rollback
- **Coordination**: Transaction coordination across services

### Validation

The `validation` package provides request and data validation:

- **Integration**: Integration with go-playground/validator
- **Custom Validators**: Define custom validation rules
- **Validation Middleware**: HTTP request validation
- **Error Handling**: Structured validation errors

## Building and Testing

ServiceLib uses Go modules for dependency management and Make for build automation.

### Prerequisites

- Go 1.24 or higher
- Make (optional, for using the Makefile)

### Build Commands

```bash
# Build the library
make build

# Run tests
make test

# Run tests with coverage
make coverage

# Run linter
make lint

# Format code
make fmt

# Check for security vulnerabilities
make security
```

## Troubleshooting

### Common Issues

#### Connection Pooling

**Issue**: Database connections are not being properly released, leading to connection pool exhaustion.

**Solution**: Ensure that all database operations properly close their resources, especially in error cases. Use the `defer` statement to ensure connections are returned to the pool.

#### Memory Leaks

**Issue**: Memory usage grows over time, indicating potential memory leaks.

**Solution**: Use the telemetry package to monitor memory usage and identify leaks. Common causes include:

- Forgetting to close response bodies
- Goroutines that never terminate
- Large objects stored in context values

#### Circular Dependencies

**Issue**: Dependency injection container fails with circular dependency errors.

**Solution**: Restructure your dependencies to break the cycle. Consider:

- Using interfaces to break direct dependencies
- Introducing a mediator or facade
- Using lazy initialization for some dependencies

### Debugging

#### Enabling Debug Logging

To enable debug logging for troubleshooting, configure the logger with debug level and console format for better readability during development.

#### Tracing Requests

For detailed request tracing:

1. Enable the tracing middleware
2. Set the sampling rate to 1.0 (100%)
3. Use the OpenTelemetry UI or Jaeger to view traces

## Best Practices

### Service Structure

- **Layered Architecture**: Organize your service with clear separation between:
  - API/Transport layer (HTTP, gRPC)
  - Service layer (business logic)
  - Repository layer (data access)

- **Dependency Injection**: Use the DI container to manage dependencies and make testing easier

- **Configuration**: Externalize all configuration and use environment variables for deployment-specific settings

### Error Handling

- **Structured Errors**: Use the errors package to create structured errors with context information for better error handling and debugging.

- **Error Categorization**: Categorize errors to handle them appropriately at the API boundary, mapping different error types to appropriate HTTP status codes.

### Performance Optimization

- **Connection Pooling**: Configure database connection pools based on expected load by setting appropriate values for max open connections, max idle connections, and connection lifetime.

- **Caching**: Use caching for frequently accessed, rarely changed data

- **Pagination**: Always implement pagination for endpoints that return collections

### Testing

- **Unit Tests**: Test each component in isolation using mocks

- **[Integration Tests](docs/Integration_Tests.md)**: Test the integration between components

- **End-to-End Tests**: Test the complete service flow

- **Load Tests**: Test performance under load to identify bottlenecks

## Architecture and Design

ServiceLib is designed with the following architectural principles:

### Modularity

Each package in ServiceLib is designed to be used independently or together with other packages. This allows you to use only the components you need without bringing in unnecessary dependencies.

### Design Patterns

ServiceLib implements several design patterns:

- **Repository Pattern**: Abstracts data access behind interfaces
- **Dependency Injection**: Manages dependencies and facilitates testing
- **Factory Pattern**: Creates complex objects with consistent configuration
- **Adapter Pattern**: Converts interfaces to work with different systems
- **Observer Pattern**: Implements event-based communication

## Compatibility and Versioning

### Version Compatibility

ServiceLib follows semantic versioning (SemVer):

- **Major version** (X.y.z): Incompatible API changes
- **Minor version** (x.Y.z): Backwards-compatible functionality
- **Patch version** (x.y.Z): Backwards-compatible bug fixes

### Go Version Compatibility

- **Minimum Go version**: 1.24
- **Tested Go versions**: 1.24, 1.25

### Dependencies

ServiceLib has the following major dependencies:

- **zap**: Structured logging
- **prometheus**: Metrics collection
- **opentelemetry**: Distributed tracing
- **validator**: Request validation
- **pgx**: PostgreSQL driver
- **gqlgen**: GraphQL implementation

### Backward Compatibility Guarantees

- No breaking changes will be introduced in minor or patch releases
- Deprecated features will be marked with `Deprecated` in the documentation
- Deprecated features will be removed only in major version releases
- Migration guides will be provided for major version upgrades

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. See the [CONTRIBUTING.md](CONTRIBUTING.md) file for detailed guidelines and instructions.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

### Coding Standards

- Follow Go best practices and style guidelines
- Write tests for new functionality
- Document public APIs
- Keep backward compatibility in mind

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
