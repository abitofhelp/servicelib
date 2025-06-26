# Context

## Overview

The Context component provides utilities for working with Go's standard context package. It extends the functionality with additional features for managing request-specific data, timeouts, and error handling in distributed systems.

## Features

- **Enhanced Context Creation**: Create contexts with various timeouts and built-in metadata
- **Automatic ID Generation**: Generate and manage request IDs, trace IDs, and correlation IDs
- **Timeout Management**: Predefined timeouts for common operations (database, network, external services)
- **Context Value Management**: Add and retrieve values from contexts with type safety
- **Error Handling**: Utilities for checking context status and handling context errors
- **Context Information**: Generate comprehensive context information for logging and debugging

## Installation

```bash
go get github.com/abitofhelp/servicelib/context
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/context/basic_usage/README.md) for a complete, runnable example of how to use the context component.

## API Documentation

### Core Types

#### Key

A type for context keys to avoid collisions.

```go
type Key string
```

#### ContextOptions

Options for creating a new context.

```go
type ContextOptions struct {
    // Timeout is the duration after which the context will be canceled
    Timeout time.Duration

    // RequestID is a unique identifier for the request
    RequestID string

    // TraceID is a unique identifier for tracing
    TraceID string

    // UserID is the ID of the user making the request
    UserID string

    // TenantID is the ID of the tenant
    TenantID string

    // Operation is the name of the operation being performed
    Operation string

    // CorrelationID is a unique identifier for correlating related operations
    CorrelationID string

    // ServiceName is the name of the service
    ServiceName string

    // Parent is the parent context
    Parent context.Context
}
```

### Key Methods

#### NewContext

Creates a new context with the specified options.

```go
func NewContext(opts ContextOptions) (context.Context, context.CancelFunc)
```

#### WithTimeout

Creates a new context with the specified timeout and options.

```go
func WithTimeout(ctx context.Context, timeout time.Duration, opts ContextOptions) (context.Context, context.CancelFunc)
```

#### CheckContext

Checks if the context is done and returns an appropriate error.

```go
func CheckContext(ctx context.Context) error
```

#### WithValues

Creates a new context with the specified key-value pairs.

```go
func WithValues(ctx context.Context, keyValues ...interface{}) context.Context
```

#### ContextInfo

Returns a string with all the context information.

```go
func ContextInfo(ctx context.Context) string
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/context/basic_usage/README.md) - Shows basic usage of the context component
- [Timeout Management](../EXAMPLES/context/timeout/README.md) - Shows how to use timeouts with contexts
- [Value Propagation](../EXAMPLES/context/value_propagation/README.md) - Shows how to propagate values through contexts
- [Error Handling](../EXAMPLES/context/error_handling_example.go) - Shows how to handle context errors

## Best Practices

1. **Always Pass Context First**: Make context the first parameter in functions that need it
2. **Use Appropriate Timeouts**: Choose the right timeout function for your operation type
3. **Check Context Status**: Regularly check if the context is done to allow for early cancellation
4. **Propagate Context Values**: Pass relevant context values down the call stack
5. **Don't Store Contexts**: Contexts should be passed as parameters, not stored in structs

## Troubleshooting

### Common Issues

#### Context Deadline Exceeded

If you're seeing "context deadline exceeded" errors, your operation is taking longer than the timeout allows. Consider:
- Increasing the timeout duration
- Optimizing the operation to complete faster
- Breaking the operation into smaller parts

#### Context Canceled

If you're seeing "context canceled" errors, the context was explicitly canceled. This is often intentional, but if not, check:
- If a parent context is being canceled unexpectedly
- If cancel() is being called too early
- If the wrong context is being passed to functions

## Related Components

- [Errors](../errors/README.md) - Error handling for context-related errors
- [Logging](../logging/README.md) - Logging with context information
- [Telemetry](../telemetry/README.md) - Distributed tracing with context propagation

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
