# Logging

## Overview

The Logging component provides a structured logging solution for your services. It offers a consistent API for logging messages with different severity levels, structured context, and integration with popular logging backends.

## Features

- **Structured Logging**: Log messages with structured context
- **Multiple Severity Levels**: Debug, Info, Warn, Error, Fatal
- **Context Integration**: Propagate context values to log entries
- **Trace ID Support**: Automatically include trace IDs in log entries
- **Multiple Backends**: Support for multiple logging backends

## Installation

```bash
go get github.com/abitofhelp/servicelib/logging
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/logging/basic_usage/README.md) for a complete, runnable example of how to use the logging component.

## Configuration

See the [Context Aware Logging example](../EXAMPLES/logging/context_aware_logging/README.md) for a complete, runnable example of how to configure the logging component.

## API Documentation

### Core Types

The logging component provides several core types for logging.

#### Logger

The main interface for logging.

```
type Logger interface {
    // Methods
}
```

#### Entry

Represents a log entry.

```
type Entry struct {
    // Fields
}
```

### Key Methods

The logging component provides several key methods for logging.

#### Debug

Logs a message at the Debug level.

```
func (l *Logger) Debug(ctx context.Context, msg string, fields ...Field)
```

#### Info

Logs a message at the Info level.

```
func (l *Logger) Info(ctx context.Context, msg string, fields ...Field)
```

#### Warn

Logs a message at the Warn level.

```
func (l *Logger) Warn(ctx context.Context, msg string, fields ...Field)
```

#### Error

Logs a message at the Error level.

```
func (l *Logger) Error(ctx context.Context, msg string, fields ...Field)
```

#### Fatal

Logs a message at the Fatal level and then calls os.Exit(1).

```
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...Field)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/logging/basic_usage/README.md) - Basic logging operations
- [Context Aware Logging](../EXAMPLES/logging/context_aware_logging/README.md) - Context-aware logging
- [Trace ID](../EXAMPLES/logging/trace_id/README.md) - Working with trace IDs

## Best Practices

1. **Use Structured Logging**: Always use structured logging for better searchability
2. **Include Context**: Always include context in log calls
3. **Use Appropriate Levels**: Use appropriate log levels for different types of messages
4. **Include Trace IDs**: Include trace IDs in log entries for distributed tracing
5. **Handle Errors**: Log errors with appropriate context

## Troubleshooting

### Common Issues

#### Missing Context Values

If context values are not appearing in log entries, ensure that you're using the correct context and that the values are set correctly.

#### Performance Issues

If you're experiencing performance issues with logging, consider reducing the log level or batching log entries.

## Related Components

- [Context](../context/README.md) - Context utilities for logging
- [Telemetry](../telemetry/README.md) - Telemetry integration for logging

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.