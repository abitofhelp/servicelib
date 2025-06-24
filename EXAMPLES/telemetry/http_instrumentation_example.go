// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using HTTP instrumentation in the telemetry package
package example_telemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	// Create a simple HTTP handler
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current span from the request context
		ctx := r.Context()

		// Log with the trace context
		logger.Info(ctx, "Handling request", zap.String("path", r.URL.Path))

		// Simulate some work
		time.Sleep(50 * time.Millisecond)

		// Write response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
	})

	// Instrument the handler with tracing
	instrumentedHandler := telemetry.InstrumentHandler(helloHandler, "hello-handler")

	// Create an HTTP middleware for tracing
	middleware := telemetry.NewHTTPMiddleware(logger)

	// Apply the middleware to the instrumented handler
	handlerWithMiddleware := middleware(instrumentedHandler)

	// Register the handler
	http.Handle("/hello", handlerWithMiddleware)

	// Instrument an HTTP client
	client := &http.Client{}
	instrumentedClient := telemetry.InstrumentClient(client)

	// Example of using the instrumented client
	go func() {
		// Wait a moment for the server to start
		time.Sleep(100 * time.Millisecond)

		// Create a request with the trace context
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/hello", nil)
		if err != nil {
			logger.Error(ctx, "Failed to create request", zap.Error(err))
			return
		}

		// Make the request with the instrumented client
		resp, err := instrumentedClient.Do(req)
		if err != nil {
			logger.Error(ctx, "Failed to make request", zap.Error(err))
			return
		}
		defer resp.Body.Close()

		// Log the response
		logger.Info(ctx, "Received response",
			zap.Int("status", resp.StatusCode),
			zap.String("trace_id", resp.Header.Get("X-Trace-ID")),
		)
	}()

	// Start the server
	fmt.Println("Starting server on :8080")
	fmt.Println("Press Ctrl+C to stop")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     logger.Fatal(ctx, "Failed to start server", zap.Error(err))
	// }

	// For this example, we'll just wait a moment
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Example completed. In a real application, HTTP requests would be traced.")
}
