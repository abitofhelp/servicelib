# Error Handling

## Overview

The Error Handling component provides a comprehensive solution for creating, categorizing, and managing errors in your services. It offers structured error types with context, stack traces, and categorization to make error handling more robust and informative.

## Features

- **Structured Errors**: Create errors with additional context and metadata
- **Error Categorization**: Categorize errors for better handling and reporting
- **Stack Traces**: Automatically capture stack traces for easier debugging
- **Error Wrapping**: Wrap errors to preserve context and add information
- **Error Context**: Add context to errors for better troubleshooting

## Installation

```bash
go get github.com/abitofhelp/servicelib/errors
```

## Quick Start

See the [Basic Error Creation example](../EXAMPLES/errors/basic_error_creation/README.md) for a complete, runnable example of how to use the errors component.

## Configuration

See the [Error Context example](../EXAMPLES/errors/error_context/README.md) for a complete, runnable example of how to configure the errors component.

## API Documentation

### Core Types

The errors component provides several core types for error handling.

#### Error

The main type that provides error handling functionality.

```
type Error struct {
    // Fields
}
```

#### ErrorCategory

Categorization for errors.

```
type ErrorCategory string
```

### Key Methods

The errors component provides several key methods for error handling.

#### New

Creates a new error with the specified message and category.

```
func New(message string, category ErrorCategory) *Error
```

#### Wrap

Wraps an existing error with additional context.

```
func Wrap(err error, message string) *Error
```

#### WithContext

Adds context to an error.

```
func (e *Error) WithContext(key string, value interface{}) *Error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Error Creation](../EXAMPLES/errors/basic_error_creation/README.md) - Creating structured errors
- [Error Categorization](../EXAMPLES/errors/error_categorization/README.md) - Categorizing errors
- [Error Context](../EXAMPLES/errors/error_context/README.md) - Adding context to errors
- [Error Wrapping](../EXAMPLES/errors/error_wrapping/README.md) - Wrapping errors
- [Stack Traces](../EXAMPLES/errors/stack_traces/README.md) - Working with stack traces

## Best Practices

1. **Use Structured Errors**: Always use structured errors for better context and debugging
2. **Categorize Errors**: Categorize errors for better handling and reporting
3. **Add Context**: Add context to errors for better troubleshooting
4. **Wrap Errors**: Wrap errors to preserve context and add information
5. **Check Stack Traces**: Use stack traces for easier debugging

## Troubleshooting

### Common Issues

#### Error Context Not Showing

If error context is not showing, ensure that you're using the WithContext method correctly.

#### Stack Traces Not Captured

If stack traces are not being captured, ensure that you're creating errors using the New or Wrap methods.

## Related Components

- [Logging](../logging/README.md) - Logging for error events
- [Middleware](../middleware/README.md) - HTTP middleware for error handling

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.