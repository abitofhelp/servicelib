# Context Package Examples

This directory contains examples demonstrating how to use the `context` package, which provides utilities for working with Go's standard context package. The package offers enhanced functionality for context creation, value propagation, timeout management, and error handling.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates basic usage of the context package.

Key concepts:
- Creating contexts with default and custom options
- Automatically generating request IDs, trace IDs, and correlation IDs
- Retrieving context information
- Using Background and TODO contexts

### 2. Timeout Example

[timeout_example.go](timeout_example.go)

Shows how to use different timeout functions provided by the context package.

Key concepts:
- Using predefined timeout durations (default, short, long)
- Creating specialized timeout contexts (database, network, external service)
- Setting custom timeouts
- Handling context deadline exceeded errors

### 3. Value Propagation Example

[value_propagation_example.go](value_propagation_example.go)

Demonstrates how to add and retrieve values from contexts.

Key concepts:
- Adding individual values to contexts (user ID, tenant ID, operation, etc.)
- Adding multiple values at once
- Propagating context values through function calls
- Using WithValues for arbitrary key-value pairs
- Automatic generation of IDs

### 4. Error Handling Example

[error_handling_example.go](error_handling_example.go)

Shows how to handle context-related errors.

Key concepts:
- Handling context timeout errors
- Handling context cancellation
- Using CheckContext to check if a context is done
- Using MustCheck for critical operations
- Error context with operation and service name

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/context/basic_usage_example.go
```

## Key Features of the Context Package

1. **Enhanced Context Creation**: Create contexts with various options like timeout, request ID, trace ID, user ID, tenant ID, operation name, correlation ID, and service name.

2. **Predefined Timeouts**: Use predefined timeout durations for common scenarios:
   - Default Timeout: 30 seconds
   - Short Timeout: 5 seconds
   - Long Timeout: 60 seconds
   - Database Timeout: 10 seconds
   - Network Timeout: 15 seconds
   - External Service Timeout: 20 seconds

3. **Value Management**: Add and retrieve values from contexts with type-safe accessors.

4. **Error Handling**: Check context status and handle errors with detailed information.

5. **Automatic ID Generation**: Automatically generate request IDs, trace IDs, and correlation IDs when not provided.

## Additional Resources

For more information about the context package, see the [context package documentation](../../context/README.md).