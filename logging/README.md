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

```go
// Example of the Logger interface
package example

import (
	"context"

	"go.uber.org/zap"
)

// Logger defines the interface for context-aware logging
type Logger interface {
	// Debug logs a debug-level message with context information
	Debug(ctx context.Context, msg string, fields ...zap.Field)

	// Info logs an info-level message with context information
	Info(ctx context.Context, msg string, fields ...zap.Field)

	// Warn logs a warning-level message with context information
	Warn(ctx context.Context, msg string, fields ...zap.Field)

	// Error logs an error-level message with context information
	Error(ctx context.Context, msg string, fields ...zap.Field)

	// Fatal logs a fatal-level message with context information
	Fatal(ctx context.Context, msg string, fields ...zap.Field)

	// Sync flushes any buffered log entries
	Sync() error
}
```

## Configuration

The logger can be configured with different log levels and output formats:

```go
// Example of logger configuration
package example

import (
	"github.com/abitofhelp/servicelib/logging"
)

func configureLoggers() {
	// Development mode (console output with colors)
	logger, err := logging.NewLogger("debug", true)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Production mode (JSON output)
	prodLogger, err := logging.NewLogger("info", false)
	if err != nil {
		panic(err)
	}
	defer prodLogger.Sync()
}
```

Available log levels:
- `debug`: Detailed information for debugging
- `info`: General operational information
- `warn`: Warning events that might cause issues
- `error`: Error events that might still allow the application to continue
- `fatal`: Severe error events that cause the application to terminate

## Logging Best Practices

1. **Use Structured Logging**: Always use structured logging with key-value pairs instead of string formatting.

   ```go
   // Example of structured logging
   package example

   import (
       "fmt"

       "go.uber.org/zap"
   )

   func logUserLogin(logger *zap.Logger, user struct{ Name string }, ip string) {
       // Good
       logger.Info("User logged in", zap.String("username", user.Name), zap.String("ip", ip))

       // Avoid
       logger.Info(fmt.Sprintf("User %s logged in from %s", user.Name, ip))
   }
   ```

2. **Include Context**: Always pass context to logging methods when available to include trace information.

3. **Log Levels**: Use appropriate log levels for different types of information.

4. **Error Logging**: When logging errors, always include the error object using `zap.Error(err)`.

5. **Performance**: In hot paths, check if the log level is enabled before constructing expensive log messages.

   ```go
   // Example of conditional logging
   package example

   import (
       "go.uber.org/zap"
       "go.uber.org/zap/zapcore"
   )

   func logExpensiveData(logger *zap.Logger) {
       if logger.Core().Enabled(zapcore.DebugLevel) {
           logger.Debug("Expensive debug info", zap.Any("data", generateExpensiveDebugData()))
       }
   }

   func generateExpensiveDebugData() interface{} {
       // This would be an expensive operation to generate debug data
       return map[string]interface{}{
           "complex": "data structure",
           "that": "would be expensive to compute",
       }
   }
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
