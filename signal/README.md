# Signal Package

The `signal` package provides utilities for handling OS signals and implementing graceful shutdown in Go applications. It helps applications respond to termination signals and execute cleanup operations before exiting.

## Features

- **Signal Handling**: Captures OS termination signals (SIGINT, SIGTERM, SIGHUP, SIGQUIT)
- **Callback Registration**: Register multiple shutdown callbacks to be executed during shutdown
- **Concurrent Execution**: Execute shutdown callbacks concurrently
- **Timeout Management**: Apply timeouts to prevent hanging during shutdown
- **Context Cancellation**: Propagate shutdown events via context cancellation
- **Multiple Signal Handling**: Handle multiple signals with forced exit on second signal
- **Logging Integration**: Comprehensive logging of shutdown events

## Installation

```bash
go get github.com/abitofhelp/servicelib/signal
```

## Usage

### Basic Signal Handling

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/signal"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    baseLogger, _ := zap.NewProduction()
    logger := logging.NewContextLogger(baseLogger)
    
    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }
    
    // Register a simple handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    // Start the server in a goroutine
    go func() {
        logger.Info(context.Background(), "Starting server on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error(context.Background(), "Server error", zap.Error(err))
        }
    }()
    
    // Wait for shutdown signal
    shutdownTimeout := 30 * time.Second
    ctx := signal.WaitForShutdown(shutdownTimeout, logger)
    
    // Perform shutdown when context is canceled
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    logger.Info(ctx, "Shutting down HTTP server")
    if err := server.Shutdown(shutdownCtx); err != nil {
        logger.Error(ctx, "Error shutting down server", zap.Error(err))
    }
    
    logger.Info(ctx, "Server stopped")
}
```

### Using Shutdown Callbacks

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/signal"
    _ "github.com/jackc/pgx/v5/stdlib"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    baseLogger, _ := zap.NewProduction()
    logger := logging.NewContextLogger(baseLogger)
    
    // Create a database connection
    db, err := sql.Open("pgx", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
    if err != nil {
        logger.Fatal(context.Background(), "Failed to connect to database", zap.Error(err))
    }
    
    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }
    
    // Register a simple handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var count int
        err := db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM users").Scan(&count)
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
        
        fmt.Fprintf(w, "Hello, World! There are %d users.", count)
    })
    
    // Start the server in a goroutine
    go func() {
        logger.Info(context.Background(), "Starting server on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error(context.Background(), "Server error", zap.Error(err))
        }
    }()
    
    // Setup graceful shutdown with callbacks
    shutdownTimeout := 30 * time.Second
    ctx, gs := signal.SetupSignalHandler(shutdownTimeout, logger)
    
    // Register server shutdown callback
    gs.RegisterCallback(func(ctx context.Context) error {
        logger.Info(ctx, "Shutting down HTTP server")
        return server.Shutdown(ctx)
    })
    
    // Register database shutdown callback
    gs.RegisterCallback(func(ctx context.Context) error {
        logger.Info(ctx, "Closing database connection")
        return db.Close()
    })
    
    // Wait for context to be canceled (when signal is received)
    <-ctx.Done()
    
    // Wait a moment for callbacks to complete
    time.Sleep(100 * time.Millisecond)
    
    logger.Info(ctx, "Application stopped")
}
```

### Custom Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/signal"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    baseLogger, _ := zap.NewProduction()
    logger := logging.NewContextLogger(baseLogger)
    
    // Create a custom graceful shutdown handler
    gs := signal.NewGracefulShutdown(30*time.Second, logger)
    
    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }
    
    // Register handlers
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    // Add a health check endpoint that reports shutdown status
    var isShuttingDown bool
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if isShuttingDown {
            w.WriteHeader(http.StatusServiceUnavailable)
            w.Write([]byte("Service is shutting down"))
            return
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Service is healthy"))
    })
    
    // Start the server in a goroutine
    go func() {
        logger.Info(context.Background(), "Starting server on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error(context.Background(), "Server error", zap.Error(err))
        }
    }()
    
    // Register shutdown callbacks in order
    
    // First, update health status
    gs.RegisterCallback(func(ctx context.Context) error {
        logger.Info(ctx, "Updating health status")
        isShuttingDown = true
        // Wait a moment for load balancers to detect the status change
        time.Sleep(2 * time.Second)
        return nil
    })
    
    // Then, shut down the server
    gs.RegisterCallback(func(ctx context.Context) error {
        logger.Info(ctx, "Shutting down HTTP server")
        return server.Shutdown(ctx)
    })
    
    // Handle shutdown and get a cancellable context
    ctx, _ := gs.HandleShutdown()
    
    // Wait for context to be canceled (when signal is received)
    <-ctx.Done()
    
    logger.Info(ctx, "Main function exiting")
}
```

## Best Practices

1. **Callback Ordering**: Register callbacks in the reverse order of resource creation to ensure proper dependency handling.

2. **Timeouts**: Set appropriate timeouts for shutdown operations to prevent hanging.

3. **Health Checks**: Update health check status early in the shutdown process to prevent new requests.

4. **Concurrent Execution**: Use concurrent execution for independent shutdown operations, but be careful with dependencies.

5. **Context Propagation**: Pass the shutdown context to all operations that need to be aware of shutdown.

6. **Error Handling**: Log errors during shutdown but continue with other shutdown operations.

7. **Signal Handling**: Be prepared to handle multiple termination signals, including forced termination.

8. **Resource Cleanup**: Ensure all resources are properly released during shutdown.

9. **Graceful Termination**: Allow in-flight operations to complete before shutting down services.

## Comparison with shutdown Package

The `signal` package is similar to the `shutdown` package but offers a different approach:

- **Callback-based**: The `signal` package uses a callback registration model
- **Concurrent Execution**: Executes shutdown callbacks concurrently
- **Object-Oriented**: Uses a `GracefulShutdown` struct with methods
- **More Flexible**: Allows for more customization of shutdown behavior

Choose the `signal` package when you need more control over the shutdown process or when you have multiple independent resources to shut down concurrently.

## License

This project is licensed under the MIT License - see the LICENSE file for details.