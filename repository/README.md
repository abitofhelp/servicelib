# Repository

## Overview

The Repository component provides a robust implementation of the Repository pattern for data access. It abstracts the details of data storage and retrieval, making it easier to switch between different data sources and to test your code.

## Features

- **Generic Implementation**: Works with any data type
- **CRUD Operations**: Standard Create, Read, Update, Delete operations
- **Transaction Support**: Support for transactions across multiple operations
- **Query Building**: Fluent API for building queries
- **Pagination**: Support for paginated results

## Installation

```bash
go get github.com/abitofhelp/servicelib/repository
```

## Quick Start

See the [Basic Repository example](../EXAMPLES/repository/basic_repository_example/README.md) for a complete, runnable example of how to use the repository component.

## Configuration

See the [Repository Factory example](../EXAMPLES/repository/repository_factory_example/README.md) for a complete, runnable example of how to configure the repository component.

## API Documentation

### Core Types

The repository component provides several core types for implementing the Repository pattern.

#### Repository

The main interface that defines the Repository pattern.

```
type Repository[T any, ID any] interface {
    // Methods
}
```

#### BaseRepository

A base implementation of the Repository interface.

```
type BaseRepository[T any, ID any] struct {
    // Fields
}
```

### Key Methods

The repository component provides several key methods for implementing the Repository pattern.

#### Create

Creates a new entity.

```
func (r *BaseRepository[T, ID]) Create(ctx context.Context, entity T) (T, error)
```

#### FindByID

Finds an entity by its ID.

```
func (r *BaseRepository[T, ID]) FindByID(ctx context.Context, id ID) (T, error)
```

#### Update

Updates an existing entity.

```
func (r *BaseRepository[T, ID]) Update(ctx context.Context, entity T) (T, error)
```

#### Delete

Deletes an entity.

```
func (r *BaseRepository[T, ID]) Delete(ctx context.Context, entity T) error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Repository](../EXAMPLES/repository/basic_repository_example/README.md) - Basic repository pattern
- [Dependency Injection](../EXAMPLES/repository/dependency_injection_example/README.md) - DI with repositories
- [Repository Factory](../EXAMPLES/repository/repository_factory_example/README.md) - Using repository factories

## Best Practices

1. **Use Interfaces**: Define repository interfaces for better testability
2. **Use Transactions**: Use transactions for operations that need to be atomic
3. **Handle Errors**: Properly handle repository errors
4. **Use Context**: Pass context to repository methods for cancellation and timeouts
5. **Use Pagination**: Use pagination for large result sets

## Troubleshooting

### Common Issues

#### Entity Not Found

If an entity is not found, check that you're using the correct ID and that the entity exists in the data source.

#### Transaction Failures

If transactions are failing, check that your data source supports transactions and that you're using them correctly.

## Related Components

- [DB](../db/README.md) - Database access for repositories
- [Errors](../errors/README.md) - Error handling for repositories
- [DI](../di/README.md) - Dependency injection for repositories

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.