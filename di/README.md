# Dependency Injection

## Overview

The Dependency Injection (DI) component provides a flexible and type-safe dependency injection container for Go applications. It helps manage application dependencies, promotes loose coupling, and simplifies testing by providing a centralized way to create and access services.

## Features

- **Generic Containers**: Type-safe dependency injection using Go generics
- **Layered Architecture Support**: Built-in support for repositories, domain services, and application services
- **Configuration Management**: Integrated configuration handling
- **Validation**: Built-in validation using go-playground/validator
- **Logging Integration**: Seamless integration with the logging component
- **Resource Management**: Automatic resource cleanup with proper Close methods

## Installation

```bash
go get github.com/abitofhelp/servicelib/di
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/di/basic_usage/README.md) for a complete, runnable example of how to use the dependency injection component.

## API Documentation

### Core Types

#### BaseContainer

The foundation container that provides common functionality.

```go
type BaseContainer[C any] struct {
    // Internal fields
}
```

#### GenericAppContainer

A generic container for applications following a layered architecture.

```go
type GenericAppContainer[R any, D any, A any, C any] struct {
    *BaseContainer[C]
    // Internal fields
}
```

#### Container

A non-generic container for backward compatibility.

```go
type Container struct {
    *BaseContainer[interface{}]
}
```

### Key Methods

#### NewBaseContainer

Creates a new base dependency injection container.

```go
func NewBaseContainer[C any](ctx context.Context, logger *zap.Logger, cfg C) (*BaseContainer[C], error)
```

#### NewGenericAppContainer

Creates a new generic application container.

```go
func NewGenericAppContainer[R any, D any, A any, C any](
    ctx context.Context,
    logger *zap.Logger,
    cfg C,
    connectionString string,
    initRepo AppRepositoryInitializer[R],
    initDomainService GenericDomainServiceInitializer[R, D],
    initAppService GenericApplicationServiceInitializer[R, D, A],
) (*GenericAppContainer[R, D, A, C], error)
```

#### NewContainer

Creates a new generic dependency injection container (for backward compatibility).

```go
func NewContainer(ctx context.Context, logger *zap.Logger, cfg interface{}) (*Container, error)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/di/basic_usage/README.md) - Shows basic usage of the dependency injection container
- [Generic Container](../EXAMPLES/di/generic_container/README.md) - Shows how to use the generic container with type safety
- [Service Container](../EXAMPLES/di/service_container/README.md) - Shows how to use the service container for more complex applications

## Best Practices

1. **Use Type Parameters**: Leverage Go generics for type safety when possible
2. **Follow Layered Architecture**: Organize code into repositories, domain services, and application services
3. **Inject Dependencies**: Pass dependencies through constructors rather than creating them inside components
4. **Close Resources**: Always call Close() on containers when they're no longer needed
5. **Use Context**: Pass context through the container to enable proper cancellation and timeouts

## Troubleshooting

### Common Issues

#### Circular Dependencies

If you encounter circular dependencies (A depends on B, which depends on A):
- Refactor your code to break the cycle
- Consider introducing an interface to decouple components
- Use a mediator pattern to handle communication between components

#### Memory Leaks

If resources aren't being cleaned up properly:
- Ensure you're calling Close() on the container when it's no longer needed
- Implement Close() methods on your custom types that need resource cleanup
- Use defer to ensure Close() is called even if errors occur

## Related Components

- [Logging](../logging/README.md) - Logging component used by the DI container
- [Config](../config/README.md) - Configuration management that works well with DI
- [DB](../db/README.md) - Database utilities that can be integrated with DI containers

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
