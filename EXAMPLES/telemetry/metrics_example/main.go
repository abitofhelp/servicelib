//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using metrics functionality in the telemetry package
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	baseLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	logger := logging.NewContextLogger(baseLogger)
	defer baseLogger.Sync()

	// Create a context
	ctx := context.Background()

	// In a real application, you would initialize the telemetry provider
	// For this example, we'll just use the functions directly

	// Log the start of the example
	logger.Info(ctx, "Starting metrics example")

	// Record HTTP request metrics
	telemetry.RecordHTTPRequest(ctx, "GET", "/api/users", 200, 50*time.Millisecond, 1024)
	logger.Info(ctx, "Recorded HTTP request metrics",
		zap.String("method", "GET"),
		zap.String("path", "/api/users"),
		zap.Int("status", 200),
		zap.Duration("duration", 50*time.Millisecond),
		zap.Int64("size", 1024),
	)

	// Record HTTP request metrics for an error case
	telemetry.RecordHTTPRequest(ctx, "POST", "/api/users", 400, 30*time.Millisecond, 512)
	logger.Info(ctx, "Recorded HTTP request metrics for error case",
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
		zap.Int("status", 400),
		zap.Duration("duration", 30*time.Millisecond),
		zap.Int64("size", 512),
	)

	// Increment and decrement requests in flight
	telemetry.IncrementRequestsInFlight(ctx, "GET", "/api/products")
	logger.Info(ctx, "Incremented requests in flight",
		zap.String("method", "GET"),
		zap.String("path", "/api/products"),
	)

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	telemetry.DecrementRequestsInFlight(ctx, "GET", "/api/products")
	logger.Info(ctx, "Decremented requests in flight",
		zap.String("method", "GET"),
		zap.String("path", "/api/products"),
	)

	// Record database operation metrics for a successful query
	telemetry.RecordDBOperation(ctx, "query", "postgres", "users", 10*time.Millisecond, nil)
	logger.Info(ctx, "Recorded database operation metrics",
		zap.String("operation", "query"),
		zap.String("database", "postgres"),
		zap.String("collection", "users"),
		zap.Duration("duration", 10*time.Millisecond),
		zap.Bool("success", true),
	)

	// Record database operation metrics for a failed query
	dbErr := errors.New("database connection error")
	telemetry.RecordDBOperation(ctx, "update", "postgres", "orders", 5*time.Millisecond, dbErr)
	logger.Info(ctx, "Recorded database operation metrics for error case",
		zap.String("operation", "update"),
		zap.String("database", "postgres"),
		zap.String("collection", "orders"),
		zap.Duration("duration", 5*time.Millisecond),
		zap.Bool("success", false),
		zap.Error(dbErr),
	)

	// Update database connection count
	telemetry.UpdateDBConnections(ctx, "postgres", 1) // +1 connection
	logger.Info(ctx, "Incremented database connections",
		zap.String("database", "postgres"),
		zap.Int64("delta", 1),
	)

	// Simulate some work with the connection
	time.Sleep(50 * time.Millisecond)

	telemetry.UpdateDBConnections(ctx, "postgres", -1) // -1 connection
	logger.Info(ctx, "Decremented database connections",
		zap.String("database", "postgres"),
		zap.Int64("delta", -1),
	)

	// Record error metrics
	telemetry.RecordErrorMetric(ctx, "validation", "create-user")
	logger.Info(ctx, "Recorded error metric",
		zap.String("error_type", "validation"),
		zap.String("operation", "create-user"),
	)

	// Log completion
	logger.Info(ctx, "Metrics example completed")

	fmt.Println("Example completed. In a real application, metrics would be exported to a monitoring system.")
}
