# Shutdown Package

The `shutdown` package provides functionality for graceful application shutdown in Go applications. It helps ensure that applications terminate cleanly, allowing resources to be properly released and pending operations to complete.

## Features

- **Signal Handling**: Captures OS termination signals (SIGINT, SIGTERM, SIGHUP)
- **Context Cancellation**: Supports shutdown via context cancellation
- **Timeout Management**: Applies timeouts to prevent hanging during shutdown
- **Multiple Signal Handling**: Forces exit if a second signal is received during shutdown
- **Error Propagation**: Returns errors from shutdown operations
- **Logging Integration**: Comprehensive logging of shutdown events

## Installation

```bash
go get github.com/abitofhelp/servicelib/shutdown
```

## Usage

### Basic Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/shutdown"
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
    
    // Define shutdown function
    shutdownFunc := func() error {
        // Create a context with timeout for shutdown
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        logger.Info(ctx, "Shutting down HTTP server")
        return server.Shutdown(ctx)
    }
    
    // Wait for shutdown signal
    ctx := context.Background()
    err := shutdown.GracefulShutdown(ctx, logger, shutdownFunc)
    if err != nil {
        logger.Error(ctx, "Error during shutdown", zap.Error(err))
    }
    
    logger.Info(ctx, "Server stopped")
}
```

### Programmatic Shutdown Initiation

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/shutdown"
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
    
    // Register handlers
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    // Add a shutdown endpoint
    var shutdownTrigger context.CancelFunc
    http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Shutting down server..."))
        
        // Trigger shutdown after response is sent
        go func() {
            time.Sleep(100 * time.Millisecond)
            if shutdownTrigger != nil {
                shutdownTrigger()
            }
        }()
    })
    
    // Start the server in a goroutine
    go func() {
        logger.Info(context.Background(), "Starting server on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error(context.Background(), "Server error", zap.Error(err))
        }
    }()
    
    // Define shutdown function
    shutdownFunc := func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        logger.Info(ctx, "Shutting down HTTP server")
        return server.Shutdown(ctx)
    }
    
    // Setup graceful shutdown with programmatic trigger
    ctx := context.Background()
    shutdownTrigger, errCh := shutdown.SetupGracefulShutdown(ctx, logger, shutdownFunc)
    
    // Wait for shutdown to complete
    err := <-errCh
    if err != nil {
        logger.Error(ctx, "Error during shutdown", zap.Error(err))
    }
    
    logger.Info(ctx, "Server stopped")
}
```

### Multiple Resource Shutdown

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "time"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/shutdown"
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
        // Use the database connection
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
    
    // Define shutdown function for multiple resources
    shutdownFunc := func() error {
        // Create a context with timeout for shutdown
        ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
        defer cancel()
        
        // Shutdown HTTP server first
        logger.Info(ctx, "Shutting down HTTP server")
        if err := server.Shutdown(ctx); err != nil {
            logger.Error(ctx, "Error shutting down server", zap.Error(err))
            return err
        }
        
        // Then close database connection
        logger.Info(ctx, "Closing database connection")
        if err := db.Close(); err != nil {
            logger.Error(ctx, "Error closing database", zap.Error(err))
            return err
        }
        
        return nil
    }
    
    // Wait for shutdown signal
    ctx := context.Background()
    err = shutdown.GracefulShutdown(ctx, logger, shutdownFunc)
    if err != nil {
        logger.Error(ctx, "Error during shutdown", zap.Error(err))
    }
    
    logger.Info(ctx, "Application stopped")
}
```

## Best Practices

1. **Resource Ordering**: Close resources in the reverse order they were created.

2. **Timeouts**: Set appropriate timeouts for shutdown operations to prevent hanging.

3. **Logging**: Log the beginning and completion of each shutdown step.

4. **Error Handling**: Properly handle and log errors during shutdown, but continue shutting down other resources.

5. **Signal Handling**: Be prepared to handle multiple termination signals.

6. **Context Usage**: Use contexts with timeouts for shutdown operations.

7. **Graceful Termination**: Allow in-flight operations to complete before shutting down.

8. **Health Checks**: Update health check status during shutdown to prevent new requests.

9. **Dependency Management**: Consider dependencies between resources when ordering shutdown.

## License

This project is licensed under the MIT License - see the LICENSE file for details.