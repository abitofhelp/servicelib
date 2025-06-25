# Dependency Injection Package

## Overview

The `di` package provides a container-based dependency injection system for Go applications. It helps manage dependencies between components, making your code more modular, testable, and maintainable.


## Features

- **Container Types**:
  - Base container
  - Service container
  - Repository container
  - Generic container

- **Features**:
  - Constructor injection
  - Type-safe dependency management with generics
  - Context-aware containers
  - Structured configuration support
  - Domain-driven design support


## Installation

```bash
go get github.com/abitofhelp/servicelib/di
```


## Quick Start

See the [Basic Usage Example](../EXAMPLES/di/basic_usage_example.go) for a complete, runnable example of how to use the di package.


## Configuration

See the [Service Container Example](../EXAMPLES/di/service_container_example.go) for a complete, runnable example of how to configure the di package.


## API Documentation


### Core Types

#### BaseContainer

The `BaseContainer` is the foundation for all container types. It provides access to common dependencies like context, logger, and configuration.

See the [Basic Usage Example](../EXAMPLES/di/basic_usage_example.go) for a complete, runnable example of how to use the BaseContainer.

#### ServiceContainer

The `ServiceContainer` is designed for domain-driven design applications. It manages repositories, domain services, and application services.

See the [Service Container Example](../EXAMPLES/di/service_container_example.go) for a complete, runnable example of how to use the ServiceContainer.

#### GenericAppContainer

The `GenericAppContainer` provides a flexible container for any application structure. It allows custom initialization of repositories, domain services, and application services.

See the [Generic Container Example](../EXAMPLES/di/generic_container_example.go) for a complete, runnable example of how to use the GenericAppContainer.


### Key Methods

#### NewBaseContainer

The `NewBaseContainer` function creates a new base container with context, logger, and configuration.

```go
baseContainer, err := di.NewBaseContainer(ctx, logger, cfg)
```

#### NewContainer

The `NewContainer` function creates a new generic container with context, logger, and configuration.

```go
container, err := di.NewContainer(ctx, logger, cfg)
```

#### NewServiceContainer

The `NewServiceContainer` function creates a new service container with context, logger, configuration, repository, and service initializers.

```go
container, err := di.NewServiceContainer(ctx, logger, cfg, repo, initDomainService, initAppService)
```

#### NewGenericAppContainer

The `NewGenericAppContainer` function creates a new generic application container with context, logger, configuration, connection string, and initializers.

```go
container, err := di.NewGenericAppContainer(ctx, logger, cfg, connectionString, initRepo, initDomainService, initAppService)
```


## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Usage Example](../EXAMPLES/di/basic_usage_example.go) - Shows basic usage of the Container type
- [Service Container Example](../EXAMPLES/di/service_container_example.go) - Shows how to use the ServiceContainer type
- [Generic Container Example](../EXAMPLES/di/generic_container_example.go) - Shows how to use the GenericAppContainer type


## Best Practices

1. **Use Interfaces**: Define clear interfaces for your repositories, domain services, and application services.

2. **Type Safety**: Use generics to ensure type safety in your dependency injection.

3. **Context Awareness**: Pass context through your application to enable proper cancellation and timeout handling.

4. **Configuration**: Use structured configuration objects rather than loose key-value pairs.

5. **Domain-Driven Design**: Organize your code according to domain-driven design principles with clear separation between repositories, domain services, and application services.

6. **Error Handling**: Always check for errors when creating containers and initializing dependencies.

7. **Testing**: Use dependency injection to make your code more testable by allowing mock implementations.


## Troubleshooting

### Common Issues

#### Container Creation Failures

**Issue**: Errors when creating a container.

**Solution**: Ensure that you're providing valid context, logger, and configuration objects. Check that your configuration implements the required interfaces.

#### Dependency Initialization Failures

**Issue**: Errors when initializing dependencies.

**Solution**: Check that your initializer functions are correctly implemented and handle all possible error cases. Ensure that dependencies are available before trying to use them.

#### Type Mismatches

**Issue**: Type assertion errors when using container methods.

**Solution**: Use the correct generic type parameters when creating containers. Ensure that your types implement the required interfaces.


## Related Components

- [Config](../config/README.md) - The config component is used to configure the di package.
- [Logging](../logging/README.md) - The logging component is used for logging in the di package.
- [Context](../context/README.md) - The context component is used for context management in the di package.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
