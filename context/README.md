# Context Package

## Overview

The `context` package extends Go's standard context package with additional utilities for request handling, cancellation, and value propagation. It provides a set of helper functions and types to make working with contexts easier and more type-safe.

## Features

- **Value Management**: Strongly typed context values
- **Timeout Management**: Utilities for working with context deadlines
- **Cancellation**: Simplified cancellation patterns
- **Propagation**: Utilities for propagating context values across service boundaries

## Installation

```bash
go get github.com/abitofhelp/servicelib/context
```

## Quick Start

See the [Basic Usage Example](../examples/context/basic_usage_example.go) for a complete, runnable example of how to use the context package.

## Configuration

The context package does not require any specific configuration. It can be used directly without any setup.

## API Documentation

### Core Types

#### Context Keys

The package provides strongly typed keys for context values to ensure type safety.

See the [Basic Usage Example](../examples/context/basic_usage_example.go) for a complete, runnable example of how to use context keys.

### Key Methods

#### Value Management

The package provides methods for setting and getting values from a context.

See the [Value Propagation Example](../examples/context/value_propagation_example.go) for a complete, runnable example of how to use value management methods.

#### Timeout Management

The package provides methods for working with context deadlines and timeouts.

See the [Timeout Example](../examples/context/timeout_example.go) for a complete, runnable example of how to use timeout management methods.

#### Error Handling

The package provides methods for handling errors in context operations.

See the [Error Handling Example](../examples/context/error_handling_example.go) for a complete, runnable example of how to handle errors.

## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Usage Example](../examples/context/basic_usage_example.go) - Shows basic usage of the context package
- [Value Propagation Example](../examples/context/value_propagation_example.go) - Shows how to propagate values through a context
- [Timeout Example](../examples/context/timeout_example.go) - Shows how to use timeouts with context
- [Error Handling Example](../examples/context/error_handling_example.go) - Shows how to handle errors in context operations

## Best Practices

1. **Type Safety**: Use strongly typed keys for context values to avoid runtime type errors.

2. **Context Propagation**: Always propagate context through your application, especially across API boundaries.

3. **Cancellation**: Use context cancellation to properly clean up resources and prevent goroutine leaks.

4. **Timeout Management**: Set appropriate timeouts for operations to prevent hanging requests.

5. **Value Scope**: Context values should be request-scoped and not used for passing optional parameters to functions.

## Troubleshooting

### Common Issues

#### Context Cancellation Not Propagating

**Issue**: Context cancellation is not being propagated to child goroutines.

**Solution**: Ensure that you're passing the context to all functions and goroutines that should respect cancellation. Check that you're properly selecting on `ctx.Done()` in long-running operations.

#### Context Values Not Found

**Issue**: Context values are not being found when expected.

**Solution**: Ensure that you're using the correct key type and that the value was actually set in the context. Remember that context values are immutable, so you need to capture the return value of `context.WithValue()`.

#### Deadlines Not Being Respected

**Issue**: Operations are not respecting context deadlines.

**Solution**: Ensure that you're checking `ctx.Err()` regularly in long-running operations. Use `select` statements with `ctx.Done()` to abort operations when the context is cancelled or times out.

## Related Components

- [Logging](../logging/README.md) - The logging component uses context for request-scoped logging.
- [Telemetry](../telemetry/README.md) - The telemetry component uses context for propagating trace information.
- [Middleware](../middleware/README.md) - The middleware component uses context for request processing.
- [Auth](../auth/README.md) - The auth component uses context for storing authentication information.

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
