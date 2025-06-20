# Logging Module

The Logging Module provides structured logging with Zap for high-performance logging in Go applications. It wraps the zap logging library and adds features like trace ID extraction from context and context-aware logging methods.

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

See the [Basic Usage example](../examples/logging/basic_usage_example.go) for a complete, runnable example of how to use the Logging module.

## API Documentation

### Basic Logging

The `NewLogger` function creates a new zap logger configured for either development or production use.

#### Basic Usage

See the [Basic Usage example](../examples/logging/basic_usage_example.go) for a complete, runnable example of how to use the basic logging functionality.

### Context-Aware Logging

The `ContextLogger` struct provides context-aware logging methods that automatically extract trace information from the context.

#### Context-Aware Logging Example

See the [Context-Aware Logging example](../examples/logging/context_aware_logging_example.go) for a complete, runnable example of how to use context-aware logging.

### Trace Integration

The `WithTraceID` function adds trace ID and span ID to the logger from the provided context.

#### Adding Trace IDs to Logs

See the [Trace ID example](../examples/logging/trace_id_example.go) for a complete, runnable example of how to add trace IDs to logs.

### Logger Interface

The `Logger` interface defines the methods for context-aware logging.

See the [Context-Aware Logging example](../examples/logging/context_aware_logging_example.go) for a complete, runnable example of how to use the Logger interface.

## Configuration

The logger can be configured with different log levels and output formats.

See the [Basic Usage example](../examples/logging/basic_usage_example.go) for a complete, runnable example of how to configure the logger.

Available log levels:
- `debug`: Detailed information for debugging
- `info`: General operational information
- `warn`: Warning events that might cause issues
- `error`: Error events that might still allow the application to continue
- `fatal`: Severe error events that cause the application to terminate

## Best Practices

1. **Use Structured Logging**: Always use structured logging with key-value pairs instead of string formatting.

2. **Include Context**: Always pass context to logging methods when available to include trace information.

3. **Log Levels**: Use appropriate log levels for different types of information.

4. **Error Logging**: When logging errors, always include the error object using `zap.Error(err)`.

5. **Performance**: In hot paths, check if the log level is enabled before constructing expensive log messages.

See the [Basic Usage example](../examples/logging/basic_usage_example.go) for examples of these best practices.

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
