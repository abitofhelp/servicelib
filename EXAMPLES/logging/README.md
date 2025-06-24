# Logging Package Examples

This directory contains examples demonstrating how to use the `logging` package, which provides structured logging utilities for Go applications. The package is built on top of the popular `zap` logging library and offers enhanced functionality for context-aware logging and trace ID propagation.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates the basic setup and usage of the logging package.

Key concepts:
- Creating a new logger with a specified log level
- Logging messages at different levels (info, debug, warn, error)
- Including structured data in log messages
- Conditional logging based on the enabled log level
- Proper cleanup with defer logger.Sync()

### 2. Context-Aware Logging Example

[context_aware_logging_example.go](context_aware_logging_example.go)

Shows how to use context-aware logging features.

Key concepts:
- Adding logger to context
- Retrieving logger from context
- Enriching logs with context information
- Propagating logger through function calls
- Maintaining consistent logging context

### 3. Trace ID Example

[trace_id_example.go](trace_id_example.go)

Demonstrates how to use trace IDs for request tracking.

Key concepts:
- Generating trace IDs for requests
- Adding trace IDs to the logging context
- Propagating trace IDs through service calls
- Correlating logs across different components
- Tracking request flow through distributed systems

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/logging/basic_usage_example.go
```

## Additional Resources

For more information about the logging package, see the [logging package documentation](../../logging/README.md).