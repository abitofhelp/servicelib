# Dependency Injection Examples

This directory contains examples demonstrating how to use the `di` package, which provides a container-based dependency injection system for Go applications.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates basic usage of the Container type, including:
- Creating a container with context, logger, and configuration
- Getting the context, logger, and configuration from the container

### 2. Service Container Example

[service_container_example.go](service_container_example.go)

Demonstrates how to use the ServiceContainer type, including:
- Implementing the Repository, DomainService, and ApplicationService interfaces
- Implementing the Config interface with AppConfig and DatabaseConfig
- Creating a service container with dependencies
- Getting repositories, domain services, and application services from the container

### 3. Generic Container Example

[generic_container_example.go](generic_container_example.go)

Demonstrates how to use the GenericAppContainer type, including:
- Implementing repository, domain service, and application service initializers
- Creating a generic container with dependencies
- Getting repositories, domain services, and application services from the container

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run basic_usage_example.go
go run service_container_example.go
go run generic_container_example.go
```

## Additional Resources

For more information about the di package, see the [di package documentation](../../di/README.md).