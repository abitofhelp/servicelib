# Signal

## Overview

The Signal component provides utilities for handling OS signals and implementing graceful shutdown in Go applications. It offers a robust mechanism for detecting termination signals, executing registered shutdown callbacks, and ensuring that applications shut down cleanly even in the face of errors or timeouts.

## Features

- **Signal Handling**: Detect and respond to OS termination signals (SIGINT, SIGTERM, SIGHUP, SIGQUIT)
- **Callback Registration**: Register functions to be called during shutdown
- **Concurrent Execution**: Execute shutdown callbacks concurrently for faster shutdown
- **Timeout Management**: Apply configurable timeouts to prevent shutdown processes from hanging
- **Context Cancellation**: Provide contexts that are canceled when shutdown signals are received
- **Multiple Signal Handling**: Force immediate exit if a second signal is received during shutdown
- **Logging Integration**: Comprehensive logging of the shutdown process
- **Thread Safety**: Thread-safe callback registration and execution

## Installation

```bash
go get github.com/abitofhelp/servicelib/signal
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/signal"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    contextLogger := logging.NewContextLogger(logger)
    
    // Set up signal handling with a 30-second timeout
    ctx, gs := signal.SetupSignalHandler(30*time.Second, contextLogger)
    
    // Register shutdown callbacks
    gs.RegisterCallback(func(ctx context.Context) error {
        fmt.Println("Closing database connections...")
        // Close database connections here
        return nil
    })
    
    gs.RegisterCallback(func(ctx context.Context) error {
        fmt.Println("Flushing in-memory cache...")
        // Flush cache here
        return nil
    })
    
    // Run your application
    fmt.Println("Application started. Press Ctrl+C to shut down.")
    
    // Wait for shutdown signal
    <-ctx.Done()
    fmt.Println("Shutdown signal received, application is shutting down...")
    
    // The registered callbacks will be executed automatically
    // The main goroutine can exit now, the signal handler will
    // ensure all callbacks are executed before final exit
}
```

## API Documentation

### Core Types

#### ShutdownCallback

A function that is called during shutdown.

```go
type ShutdownCallback func(ctx context.Context) error
```

#### GracefulShutdown

Represents a graceful shutdown handler.

```go
type GracefulShutdown struct {
    // Internal fields
}
```

### Key Methods

#### NewGracefulShutdown

Creates a new graceful shutdown handler.

```go
func NewGracefulShutdown(timeout time.Duration, logger *logging.ContextLogger) *GracefulShutdown
```

#### RegisterCallback

Registers a callback to be called during shutdown.

```go
func (gs *GracefulShutdown) RegisterCallback(callback ShutdownCallback)
```

#### HandleShutdown

Handles graceful shutdown and returns a context that will be canceled when a shutdown signal is received.

```go
func (gs *GracefulShutdown) HandleShutdown() (context.Context, context.CancelFunc)
```

#### WaitForShutdown

Blocks until a shutdown signal is received and returns a context that will be canceled when that happens.

```go
func WaitForShutdown(timeout time.Duration, logger *logging.ContextLogger) context.Context
```

#### SetupSignalHandler

Sets up a signal handler for graceful shutdown and returns a context and the GracefulShutdown instance.

```go
func SetupSignalHandler(timeout time.Duration, logger *logging.ContextLogger) (context.Context, *GracefulShutdown)
```

## Examples

Currently, there are no dedicated examples for the signal package in the EXAMPLES directory. The following code snippets demonstrate common usage patterns:

### Basic Signal Handling

```go
// Create a logger
logger, _ := zap.NewProduction()
contextLogger := logging.NewContextLogger(logger)

// Set up signal handling with a 30-second timeout
ctx, gs := signal.SetupSignalHandler(30*time.Second, contextLogger)

// Wait for shutdown signal
<-ctx.Done()
fmt.Println("Shutdown signal received")
```

### Registering Shutdown Callbacks

```go
// Create a graceful shutdown handler
gs := signal.NewGracefulShutdown(30*time.Second, contextLogger)

// Register callbacks
gs.RegisterCallback(func(ctx context.Context) error {
    fmt.Println("Closing database connections...")
    return db.Close()
})

gs.RegisterCallback(func(ctx context.Context) error {
    fmt.Println("Stopping HTTP server...")
    return server.Shutdown(ctx)
})

// Handle shutdown
ctx, cancel := gs.HandleShutdown()
defer cancel()

// Wait for shutdown signal
<-ctx.Done()
```

### Simple Waiting for Shutdown

```go
// Wait for shutdown with a 10-second timeout
ctx := signal.WaitForShutdown(10*time.Second, contextLogger)

// Run your application until shutdown
select {
case <-ctx.Done():
    fmt.Println("Shutting down...")
    // Perform cleanup
}
```

## Best Practices

1. **Register Callbacks in Reverse Order**: Register shutdown callbacks in the reverse order of resource creation
2. **Set Appropriate Timeouts**: Configure timeouts based on expected shutdown durations
3. **Handle Callback Errors**: Ensure callbacks handle errors properly and don't panic
4. **Use Context Cancellation**: Use the provided context to propagate cancellation to ongoing operations
5. **Log Shutdown Progress**: Log the start and completion of each shutdown step

## Troubleshooting

### Common Issues

#### Shutdown Hanging

If your application is hanging during shutdown:
- Check if any callbacks are deadlocked
- Ensure all resources have proper close methods
- Verify that timeouts are set appropriately
- Look for infinite loops or blocking operations in callbacks

#### Premature Exit

If your application exits before completing the shutdown process:
- Ensure you're waiting for the context to be canceled
- Check if any panics are occurring during shutdown
- Verify that callbacks are properly registered

## Related Components

- [Shutdown](../shutdown/README.md) - Higher-level shutdown utilities built on top of this package
- [Context](../context/README.md) - Context utilities for timeout and cancellation
- [Logging](../logging/README.md) - Logging integration for shutdown events

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
