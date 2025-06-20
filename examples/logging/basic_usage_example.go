// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the logging package
package main

import (
    "errors"
    "log"

    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
)

func main() {
    // Create a new logger
    logger, err := logging.NewLogger("info", true)
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Sync()

    // Log messages with different levels
    logger.Info("Starting application", zap.String("version", "1.0.0"))
    logger.Debug("Debug information") // Won't be displayed with info level
    logger.Warn("Warning message", zap.Int("code", 123))
    logger.Error("Error occurred", zap.Error(errors.New("sample error")))

    // Log with structured data
    logger.Info("User action",
        zap.String("user_id", "user123"),
        zap.String("action", "login"),
        zap.String("ip_address", "192.168.1.1"),
    )

    // Conditional logging
    if logger.Core().Enabled(zap.DebugLevel) {
        // This code will only execute if debug level is enabled
        logger.Debug("Expensive debug info", zap.Any("data", map[string]interface{}{
            "complex": "data structure",
            "that": "would be expensive to compute",
        }))
    }

    // In a real application, you might log when the application is shutting down
    logger.Info("Application shutting down")

    // Note: logger.Fatal would terminate the program, so we don't call it in this example
    // logger.Fatal("Fatal error", zap.Error(errors.New("fatal error")))
}