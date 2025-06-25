# Logging Module
The Logging Module provides structured logging with Zap for high-performance logging in Go applications. It wraps the zap logging library and adds features like trace ID extraction from context and context-aware logging methods.


## Overview

The Logging module provides a comprehensive structured logging system for Go applications. It serves as a central component in the ServiceLib library for recording application events, errors, and diagnostic information. The module wraps the high-performance Zap logging library, adding features like context-aware logging, trace ID extraction, and standardized log formatting.

This module is designed to integrate seamlessly with other ServiceLib components, providing consistent logging patterns across your application. It supports both development and production environments with appropriate log formatting and levels for each.

## Features

- **Log Levels**: Debug, Info, Warn, Error, Fatal
- **Structured Logging**: Key-value pairs for better searchability
- **Output Formats**: JSON, console
- **Context-Aware Logging**: Log with context information including trace IDs
- **Performance**: High-performance logging with minimal allocations
- **Trace Integration**: Automatic extraction of trace IDs from OpenTelemetry context


## Installation

```bash
go get github.com/abitofhelp/servicelib/logging
```


## Quick Start

See the [Basic Usage example](../EXAMPLES/logging/basic_usage/main.go) for a complete, runnable example of how to use the Logging module.


## Configuration

The logger can be configured with different log levels and output formats.

See the [Basic Usage example](../EXAMPLES/logging/basic_usage/main.go) for a complete, runnable example of how to configure the logger.

Available log levels:
- `debug`: Detailed information for debugging
- `info`: General operational information
- `warn`: Warning events that might cause issues
- `error`: Error events that might still allow the application to continue
- `fatal`: Severe error events that cause the application to terminate


## API Documentation

### Basic Logging

The `NewLogger` function creates a new zap logger configured for either development or production use.

#### Basic Usage

See the [Basic Usage example](../EXAMPLES/logging/basic_usage/main.go) for a complete, runnable example of how to use the basic logging functionality.

### Context-Aware Logging

The `ContextLogger` struct provides context-aware logging methods that automatically extract trace information from the context.

#### Context-Aware Logging Example

See the [Context-Aware Logging example](../EXAMPLES/logging/context_aware_logging/main.go) for a complete, runnable example of how to use context-aware logging.

### Trace Integration

The `WithTraceID` function adds trace ID and span ID to the logger from the provided context.

#### Adding Trace IDs to Logs

See the [Trace ID example](../EXAMPLES/logging/trace_id/main.go) for a complete, runnable example of how to add trace IDs to logs.

### Logger Interface

The `Logger` interface defines the methods for context-aware logging.

See the [Context-Aware Logging example](../EXAMPLES/logging/context_aware_logging/main.go) for a complete, runnable example of how to use the Logger interface.


### Core Types

The logging module provides several core types for structured logging:

#### Logger

The `Logger` interface defines the methods for context-aware logging. It provides methods for logging at different levels (Debug, Info, Warn, Error, Fatal) with context information.

The interface includes methods like Debug(), Info(), Warn(), Error(), and Fatal(), all of which take a context and message, along with optional fields for structured logging.

See the [Context-Aware Logging example](../EXAMPLES/logging/context_aware_logging_example.go) for a complete, runnable example of how to use the Logger interface.

#### ContextLogger

The `ContextLogger` struct implements the Logger interface and provides context-aware logging methods that automatically extract trace information from the context.

The ContextLogger wraps a zap.Logger and adds functionality to extract trace IDs from the context and include them in log entries.

See the [Context-Aware Logging example](../EXAMPLES/logging/context_aware_logging_example.go) for a complete, runnable example of how to use the ContextLogger.

### Key Methods

The logging module provides several key methods for structured logging:

#### NewLogger

The `NewLogger` function creates a new zap logger configured for either development or production use.

```
func NewLogger(config Config) (*zap.Logger, error)
```

This function takes a Config struct that specifies the log level, output format, and other options. It returns a configured zap.Logger that can be used directly or wrapped in a ContextLogger.

See the [Basic Usage example](../EXAMPLES/logging/basic_usage_example.go) for a complete, runnable example of how to use the NewLogger function.

#### NewContextLogger

The `NewContextLogger` function creates a new ContextLogger that wraps a zap.Logger.

```
func NewContextLogger(logger *zap.Logger) *ContextLogger
```

This function takes a zap.Logger and returns a ContextLogger that provides context-aware logging methods.

See the [Context-Aware Logging example](../EXAMPLES/logging/context_aware_logging_example.go) for a complete, runnable example of how to use the NewContextLogger function.

#### WithTraceID

The `WithTraceID` function adds trace ID and span ID to the logger from the provided context.

```
func WithTraceID(ctx context.Context, logger *zap.Logger) *zap.Logger
```

This function takes a context and a zap.Logger, extracts trace information from the context, and returns a new logger with the trace information added as fields.

See the [Trace ID example](../EXAMPLES/logging/trace_id_example.go) for a complete, runnable example of how to use the WithTraceID function.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/logging/basic_usage_example.go) - Shows basic usage of the logging
- [Advanced Configuration](../EXAMPLES/logging/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/logging/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Use Structured Logging**: Always use structured logging with key-value pairs instead of string formatting.

2. **Include Context**: Always pass context to logging methods when available to include trace information.

3. **Log Levels**: Use appropriate log levels for different types of information.

4. **Error Logging**: When logging errors, always include the error object using `zap.Error(err)`.

5. **Performance**: In hot paths, check if the log level is enabled before constructing expensive log messages.

See the [Basic Usage example](../EXAMPLES/logging/basic_usage_example.go) for examples of these best practices.


## Troubleshooting

### Common Issues

#### Logger Not Initialized

**Issue**: Attempting to use a logger that hasn't been properly initialized.

**Solution**: Always check for errors when creating a logger and ensure that the logger is initialized before use.

#### Log Messages Not Appearing

**Issue**: Log messages are not appearing in the expected output.

**Solution**: Check that the log level is set appropriately. For example, debug messages won't appear if the log level is set to info or higher.

#### Performance Issues

**Issue**: Logging is causing performance issues in the application.

**Solution**: Use conditional logging for expensive operations, ensure that debug logging is disabled in production, and consider using sampling for high-volume logs.


## Related Components

- [Telemetry](../telemetry/README.md) - The telemetry component uses the logging component for logging telemetry events.
- [Middleware](../middleware/README.md) - The middleware component uses the logging component for request logging.
- [Health](../health/README.md) - The health component uses the logging component for logging health check results.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
