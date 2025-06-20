# Repository Module

The Repository Module provides generic repository interfaces for entity persistence operations in Go applications. It implements the Repository Pattern, which is a key component of Domain-Driven Design (DDD) and Hexagonal Architecture.

## Features

- **Generic Interfaces**: Type-safe repository interfaces using Go generics
- **Hexagonal Architecture**: Supports ports and adapters pattern
- **Domain-Driven Design**: Facilitates separation of domain and infrastructure concerns
- **Repository Pattern**: Standardized approach to data access
- **Repository Factory**: Interface for creating repositories

## Installation

```bash
go get github.com/abitofhelp/servicelib/repository
```

## Quick Start

See the [Basic Repository example](../examples/repository/basic_repository_example.go) for a complete, runnable example of how to use the Repository module.

## API Documentation

### Basic Repository Implementation

The `Repository` interface provides a generic interface for entity persistence operations.

#### Basic Usage

See the [Basic Repository example](../examples/repository/basic_repository_example.go) for a complete, runnable example of how to implement and use a basic repository.

### Repository Factory

The `RepositoryFactory` interface provides a way to create and manage repositories for different entity types.

#### Using Repository Factory

See the [Repository Factory example](../examples/repository/repository_factory_example.go) for a complete, runnable example of how to implement and use a repository factory.

### Dependency Injection

The Repository pattern works well with dependency injection, allowing for better testability and flexibility.

#### Integration with Dependency Injection

See the [Dependency Injection example](../examples/repository/dependency_injection_example.go) for a complete, runnable example of how to use repositories with dependency injection.

## Best Practices

1. **Interface Segregation**: Keep repository interfaces focused on specific entity types.

2. **Dependency Inversion**: Depend on repository interfaces, not concrete implementations.

3. **Testability**: Use in-memory repository implementations for testing.

4. **Transaction Management**: Consider adding transaction support for operations that span multiple repositories.

5. **Error Handling**: Use domain-specific errors for repository operations.

6. **Context Usage**: Always pass a context to repository methods for cancellation and timeout support.

7. **Repository Factory**: Use a factory to create and manage repositories when dealing with multiple entity types.

8. **Concurrency**: Ensure thread safety in repository implementations.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
