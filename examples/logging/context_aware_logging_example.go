// Copyright (c) 2025 A Bit of Help, Inc.

// Example of context-aware logging
package main

import (
    "context"
    "errors"
    "fmt"
    "log"

    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
)

func main() {
    // Create a base logger
    baseLogger, err := logging.NewLogger("info", true)
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer baseLogger.Sync()
    
    // Create a context logger
    contextLogger := logging.NewContextLogger(baseLogger)
    
    // Create a context
    ctx := context.Background()
    
    // Add request ID to context (in a real application, this would come from the request)
    ctx = context.WithValue(ctx, "request_id", "req-123456")
    
    // Log with context
    contextLogger.Info(ctx, "Processing request", zap.String("endpoint", "/api/users"))
    
    // Simulate processing in another function
    processRequest(ctx, contextLogger)
    
    fmt.Println("Example completed. In a real application, logs would be visible in the console or log files.")
}

// processRequest simulates processing a request with context-aware logging
func processRequest(ctx context.Context, logger *logging.ContextLogger) {
    // Log the start of processing
    logger.Debug(ctx, "Starting request processing", zap.String("component", "processor"))
    
    // Simulate some work
    // ...
    
    // Log an informational message
    logger.Info(ctx, "Processing data", zap.Int("items", 42))
    
    // Simulate an error
    err := errors.New("database connection failed")
    
    // Log the error with context
    logger.Error(ctx, "Failed to process data", 
        zap.Error(err),
        zap.String("component", "database"),
    )
    
    // In a real application, you might have nested function calls
    // that all use the same context for logging
    validateData(ctx, logger)
}

// validateData simulates a nested function call with context-aware logging
func validateData(ctx context.Context, logger *logging.ContextLogger) {
    logger.Info(ctx, "Validating data", zap.String("validation_type", "schema"))
    
    // The trace ID from the context will be automatically included in all these logs
    // making it easy to correlate logs from the same request
}