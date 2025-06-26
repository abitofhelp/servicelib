# Shutdown

## Overview

The Shutdown component provides functionality for graceful application shutdown in Go applications. It handles OS signals (SIGINT, SIGTERM, SIGHUP) and context cancellation to trigger controlled shutdown processes, ensuring that resources are properly released and in-flight operations can complete before the application exits.

## Features

- **Signal Handling**: Automatically detect and respond to OS termination signals
- **Context Cancellation**: Support for context-based shutdown initiation
- **Timeout Management**: Apply configurable timeouts to prevent shutdown processes from hanging
- **Multiple Signal Handling**: Force immediate exit if a second signal is received during shutdown
- **Error Propagation**: Capture and return errors that occur during the shutdown process
- **Programmatic Shutdown**: Trigger shutdown programmatically in addition to signal-based shutdown
- **Logging Integration**: Comprehensive logging of the shutdown process

## Installation

```bash
go get github.com/abitofhelp/servicelib/shutdown
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/shutdown"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    contextLogger := logging.NewContextLogger(logger)
    
    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }
    
    // Define a shutdown function
    shutdownFunc := func() error {
        // Create a timeout for server shutdown
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        fmt.Println("Shutting down HTTP server...")
        return server.Shutdown(ctx)
    }
    
    // Set up graceful shutdown
    ctx := context.Background()
    cancel, errCh := shutdown.SetupGracefulShutdown(ctx, contextLogger, shutdownFunc)
    defer cancel() // Ensure cancellation in case of early return
    
    // Start the HTTP server
    go func() {
        fmt.Println("Starting HTTP server on :8080")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()
    
    // Wait for shutdown to complete
    if err := <-errCh; err != nil {
        log.Fatalf("Shutdown error: %v", err)
    }
    
    fmt.Println("Application has shut down gracefully")
}
```

## API Documentation

### Key Functions

#### GracefulShutdown

Waits for termination signals and calls the provided shutdown function.

```go
func GracefulShutdown(ctx context.Context, logger *logging.ContextLogger, shutdownFunc func() error) error
```

This function handles OS signals (SIGINT, SIGTERM, SIGHUP) and context cancellation to trigger graceful shutdown. It also handles multiple signals, forcing exit if a second signal is received during shutdown. A default timeout of 30 seconds is applied to the shutdown function to prevent hanging.

#### SetupGracefulShutdown

Sets up a goroutine that will handle graceful shutdown.

```go
func SetupGracefulShutdown(ctx context.Context, logger *logging.ContextLogger, shutdownFunc func() error) (context.CancelFunc, <-chan error)
```

This function creates a new context with cancellation and starts a background goroutine that calls GracefulShutdown. It returns a cancel function that can be called to trigger shutdown programmatically and a channel that will receive any error that occurs during shutdown.

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [HTTP Server](../EXAMPLES/shutdown/http_server/README.md) - Shows how to gracefully shut down an HTTP server
- [Database Connections](../EXAMPLES/shutdown/database_connections/README.md) - Shows how to gracefully close database connections
- [Worker Pool](../EXAMPLES/shutdown/worker_pool/README.md) - Shows how to gracefully shut down a worker pool

## Best Practices

1. **Define Clear Shutdown Order**: Close resources in the reverse order they were created
2. **Set Appropriate Timeouts**: Configure timeouts based on expected shutdown durations
3. **Handle In-Flight Operations**: Allow in-flight operations to complete before shutting down
4. **Log Shutdown Progress**: Log the start and completion of each shutdown step
5. **Propagate Errors**: Return errors from the shutdown function to help diagnose issues

## Troubleshooting

### Common Issues

#### Shutdown Hanging

If your application is hanging during shutdown:
- Check if any goroutines are deadlocked
- Ensure all resources have proper close methods
- Verify that timeouts are set appropriately
- Look for infinite loops or blocking operations in the shutdown function

#### Premature Exit

If your application exits before completing the shutdown process:
- Ensure you're waiting for the error channel from SetupGracefulShutdown
- Check if any panics are occurring during shutdown
- Verify that the shutdown function is properly implemented

## Related Components

- [Signal](../signal/README.md) - Lower-level signal handling utilities
- [Context](../context/README.md) - Context utilities for timeout and cancellation
- [Logging](../logging/README.md) - Logging integration for shutdown events
- [Errors](../errors/README.md) - Error handling for shutdown errors

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.