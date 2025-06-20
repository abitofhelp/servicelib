# Repository Package Examples

This directory contains examples demonstrating how to use the `repository` package, which provides utilities for implementing the repository pattern in Go applications. The repository pattern abstracts data access logic, making it easier to switch between different data storage implementations and improving testability.

## Examples

### 1. Basic Repository Example

[basic_repository_example.go](basic_repository_example.go)

Demonstrates the basic implementation and usage of a repository.

Key concepts:
- Creating a domain entity
- Implementing a repository interface
- Using an in-memory repository implementation
- Performing basic CRUD operations
- Using generics for type-safe repository operations

### 2. Dependency Injection Example

[dependency_injection_example.go](dependency_injection_example.go)

Shows how to use dependency injection with repositories.

Key concepts:
- Injecting repository dependencies
- Decoupling business logic from data access
- Using interfaces for loose coupling
- Testing with mock repositories
- Configuring repositories at runtime

### 3. Repository Factory Example

[repository_factory_example.go](repository_factory_example.go)

Demonstrates how to use a repository factory to create repositories.

Key concepts:
- Creating a repository factory
- Managing multiple repository types
- Centralizing repository creation logic
- Configuring repositories based on application settings
- Handling repository initialization

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/repository/basic_repository_example.go
```

## Additional Resources

For more information about the repository package, see the [repository package documentation](../../repository/README.md).