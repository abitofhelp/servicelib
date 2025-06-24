// Copyright (c) 2025 A Bit of Help, Inc.

// Example of adding trace IDs to logs
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abitofhelp/servicelib/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main() {
	// Create a base logger
	baseLogger, err := logging.NewLogger("info", false)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer baseLogger.Sync()

	// In a real application, you would initialize the OpenTelemetry SDK
	// For this example, we'll just use the global tracer provider

	// Create a tracer
	tracer := otel.Tracer("example-tracer")

	// Create a context with a trace
	ctx, span := tracer.Start(context.Background(), "example-operation")
	defer span.End()

	// Get the trace ID for display
	spanCtx := trace.SpanContextFromContext(ctx)
	traceID := spanCtx.TraceID().String()

	fmt.Printf("Generated Trace ID: %s\n", traceID)

	// Add trace ID to logger
	loggerWithTrace := logging.WithTraceID(ctx, baseLogger)

	// Log with trace ID
	loggerWithTrace.Info("This log includes trace ID")

	// Simulate a function call that uses the context
	processWithTrace(ctx, baseLogger)

	fmt.Println("Example completed. In a real application, logs would include the trace ID.")
}

// processWithTrace simulates a function that processes data and logs with trace ID
func processWithTrace(ctx context.Context, baseLogger *zap.Logger) {
	// Get a logger with trace ID for this function
	logger := logging.WithTraceID(ctx, baseLogger)

	// Log with the trace ID
	logger.Info("Processing data with trace",
		zap.String("component", "processor"),
		zap.Int("items", 42),
	)

	// Create a child span for a sub-operation
	_, childSpan := otel.Tracer("example-tracer").Start(ctx, "sub-operation")
	defer childSpan.End()

	// Log within the child span - will have the same trace ID but different span ID
	logger.Info("Sub-operation completed")

	// In a real application, you might pass the context to other functions or services
	// and they would all log with the same trace ID, making it easy to follow the
	// request flow across the system
}