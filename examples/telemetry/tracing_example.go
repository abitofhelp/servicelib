// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using tracing functionality in the telemetry package
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
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
	// For this example, we'll just use the global OpenTelemetry functions

	// Start a span
	ctx, span := telemetry.StartSpan(ctx, "main-operation")
	defer span.End()

	// Add attributes to the span
	telemetry.AddSpanAttributes(ctx,
		attribute.String("example.type", "tracing"),
		attribute.Int("example.version", 1),
	)

	// Log with the trace context
	logger.Info(ctx, "Starting tracing example")

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Use WithSpan to create a child span for a function
	err = telemetry.WithSpan(ctx, "child-operation", func(ctx context.Context) error {
		// Log with the child span context
		logger.Info(ctx, "Executing child operation")

		// Simulate some work
		time.Sleep(50 * time.Millisecond)

		// Add more attributes to the child span
		telemetry.AddSpanAttributes(ctx,
			attribute.String("child.status", "processing"),
		)

		// Simulate a successful operation
		return nil
	})

	if err != nil {
		logger.Error(ctx, "Error in child operation", zap.Error(err))
	}

	// Use WithSpanTimed to measure execution time
	duration, err := telemetry.WithSpanTimed(ctx, "timed-operation", func(ctx context.Context) error {
		// Log with the timed span context
		logger.Info(ctx, "Executing timed operation")

		// Simulate some work
		time.Sleep(75 * time.Millisecond)

		// Simulate an error
		return errors.New("simulated error in timed operation")
	})

	// Record the error on the span
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		logger.Error(ctx, "Error in timed operation",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
	}

	// Log completion
	logger.Info(ctx, "Tracing example completed")

	fmt.Println("Example completed. In a real application, spans would be exported to a tracing backend.")
}
