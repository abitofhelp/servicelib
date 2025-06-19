# Logging Package

The `logging` package provides structured logging with Zap for high-performance logging in Go applications. It wraps the zap logging library and adds features like trace ID extraction from context and context-aware logging methods.

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

## Usage

### Basic Usage

```go
package main

import (
    "github.com/abitofhelp/servicelib/logging"
)

func main() {
    // Create a new logger
    logger, err := logging.NewLogger("info", true)
    if err != nil {
        panic("Failed to create logger: " + err.Error())
    }
    defer logger.Sync()

    // Log messages
    logger.Info("Starting application", zap.String("version", "1.0.0"))
    logger.Debug("Debug information")
    logger.Warn("Warning message", zap.Int("code", 123))
    logger.Error("Error occurred", zap.Error(errors.New("sample error")))
}
```

### Context-Aware Logging

```go
package main

import (
    "context"
    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
)

func main() {
    // Create a base logger
    baseLogger, _ := logging.NewLogger("info", true)
    
    // Create a context logger
    contextLogger := logging.NewContextLogger(baseLogger)
    
    // Create a context
    ctx := context.Background()
    
    // Log with context
    contextLogger.Info(ctx, "Processing request", zap.String("requestID", "123456"))
    
    // In a function with trace context
    processRequest(ctx, contextLogger)
}

func processRequest(ctx context.Context, logger *logging.ContextLogger) {
    // The trace ID from the context will be automatically included in the log
    logger.Info(ctx, "Processing data", zap.Int("items", 42))
    
    if err := doSomething(); err != nil {
        logger.Error(ctx, "Failed to process data", zap.Error(err))
    }
}
```

### Adding Trace IDs to Logs

```go
package main

import (
    "context"
    "github.com/abitofhelp/servicelib/logging"
    "go.opentelemetry.io/otel"
    "go.uber.org/zap"
)

func main() {
    // Create a base logger
    baseLogger, _ := logging.NewLogger("info", false)
    
    // Create a context with trace
    ctx, span := otel.Tracer("example").Start(context.Background(), "operation")
    defer span.End()
    
    // Add trace ID to logger
    loggerWithTrace := logging.WithTraceID(ctx, baseLogger)
    
    // Log with trace ID
    loggerWithTrace.Info("This log includes trace ID")
}
```

## Configuration

The logger can be configured with different log levels and output formats:

```go
// Development mode (console output with colors)
logger, err := logging.NewLogger("debug", true)

// Production mode (JSON output)
logger, err := logging.NewLogger("info", false)
```

Available log levels:
- `debug`: Detailed information for debugging
- `info`: General operational information
- `warn`: Warning events that might cause issues
- `error`: Error events that might still allow the application to continue
- `fatal`: Severe error events that cause the application to terminate

## Best Practices

1. **Use Structured Logging**: Always use structured logging with key-value pairs instead of string formatting.

   ```go
   // Good
   logger.Info("User logged in", zap.String("username", user.Name), zap.String("ip", ip))
   
   // Avoid
   logger.Info(fmt.Sprintf("User %s logged in from %s", user.Name, ip))
   ```

2. **Include Context**: Always pass context to logging methods when available to include trace information.

3. **Log Levels**: Use appropriate log levels for different types of information.

4. **Error Logging**: When logging errors, always include the error object using `zap.Error(err)`.

5. **Performance**: In hot paths, check if the log level is enabled before constructing expensive log messages.

   ```go
   if logger.Core().Enabled(zapcore.DebugLevel) {
       logger.Debug("Expensive debug info", zap.Any("data", generateExpensiveDebugData()))
   }
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.