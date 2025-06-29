# Auth Errors

## Overview

The Auth Errors package provides specialized error handling for the authentication and authorization components. It extends the core errors package with auth-specific error types and utilities.

## Features

- **Auth-Specific Errors**: Specialized error types for authentication and authorization failures
- **Context Information**: Add and retrieve context information from errors
- **Operation Tracking**: Track the operation that caused an error
- **Error Wrapping**: Wrap errors with additional information
- **Error Inspection**: Extract information from errors

## Installation

```bash
go get github.com/abitofhelp/servicelib/auth/errors
```

## Quick Start

See the [Error Handling example](../../EXAMPLES/auth/error_handling/README.md) for a complete, runnable example of how to use the auth errors package.

## API Documentation

### Core Types

#### AuthError

The main error type for authentication and authorization errors.

```go
type AuthError struct {
    // Fields
}
```

### Key Methods

#### NewAuthError

Creates a new auth error.

```go
func NewAuthError(err error, op string, message string, context map[string]interface{}) *AuthError
```

#### WithContext

Adds context information to an error.

```go
func WithContext(err error, key string, value interface{}) error
```

#### WithOp

Adds operation information to an error.

```go
func WithOp(err error, op string) error
```

#### WithMessage

Adds a message to an error.

```go
func WithMessage(err error, message string) error
```

#### Wrap

Wraps an error with a message.

```go
func Wrap(err error, message string) error
```

#### GetContext

Gets context information from an error.

```go
func GetContext(err error, key string) (interface{}, bool)
```

#### GetOp

Gets operation information from an error.

```go
func GetOp(err error) (string, bool)
```

#### GetMessage

Gets the message from an error.

```go
func GetMessage(err error) (string, bool)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Error Handling](../../EXAMPLES/auth/error_handling/README.md) - Shows how to handle auth errors

## Best Practices

1. **Use Specific Error Types**: Use the most specific error type for the situation
2. **Add Context Information**: Always add relevant context information to errors
3. **Track Operations**: Always track the operation that caused an error
4. **Check Error Types**: Check error types before handling them
5. **Log Errors**: Always log errors with their context information

## Troubleshooting

### Common Issues

#### Error Type Checking

If error type checking fails, ensure you're using the Is method or errors.Is function.

## Related Components

- [Auth](../README.md) - The parent auth package
- [Errors](../../errors/README.md) - The core errors package
- [Logging](../../logging/README.md) - Logging for errors

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.