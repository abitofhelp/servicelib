// Copyright (c) 2025 A Bit of Help, Inc.

// Example of integrating with Prometheus in the telemetry package
package main

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

	// Log the start of the example
	logger.Info(ctx, "Starting Prometheus integration example")

	// Create a Prometheus handler
	prometheusHandler := telemetry.CreatePrometheusHandler()

	// Register the Prometheus handler
	http.Handle("/metrics", prometheusHandler)

	// Create a simple HTTP handler that generates some metrics
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		// Record the start time
		startTime := time.Now()

		// Get the context from the request
		ctx := r.Context()

		// Log the request
		logger.Info(ctx, "Handling request", zap.String("path", r.URL.Path))

		// Increment requests in flight
		telemetry.IncrementRequestsInFlight(ctx, r.Method, r.URL.Path)
		defer telemetry.DecrementRequestsInFlight(ctx, r.Method, r.URL.Path)

		// Simulate a database query
		telemetry.RecordDBOperation(ctx, "query", "postgres", "users", 15*time.Millisecond, nil)

		// Simulate some work
		time.Sleep(50 * time.Millisecond)

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody := []byte(`{"users":[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]}`)
		w.Write(responseBody)

		// Record HTTP request metrics
		duration := time.Since(startTime)
		telemetry.RecordHTTPRequest(ctx, r.Method, r.URL.Path, http.StatusOK, duration, int64(len(responseBody)))
	})

	// Start the server in a goroutine
	go func() {
		logger.Info(ctx, "Starting server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "Server error", zap.Error(err))
		}
	}()

	// Simulate some traffic to generate metrics
	go func() {
		// Wait for the server to start
		time.Sleep(100 * time.Millisecond)

		// Make some requests to generate metrics
		for i := 0; i < 5; i++ {
			resp, err := http.Get("http://localhost:8080/api/users")
			if err != nil {
				logger.Error(ctx, "Request failed", zap.Error(err))
				continue
			}
			resp.Body.Close()
			logger.Info(ctx, "Request completed", zap.Int("status", resp.StatusCode))
			time.Sleep(100 * time.Millisecond)
		}

		// Make a request to the metrics endpoint
		resp, err := http.Get("http://localhost:8080/metrics")
		if err != nil {
			logger.Error(ctx, "Metrics request failed", zap.Error(err))
			return
		}
		defer resp.Body.Close()

		logger.Info(ctx, "Metrics endpoint accessed", zap.Int("status", resp.StatusCode))
	}()

	// In a real application, you would wait for a shutdown signal
	// For this example, we'll just wait a moment
	fmt.Println("Server running on :8080")
	fmt.Println("Metrics endpoint: http://localhost:8080/metrics")
	fmt.Println("API endpoint: http://localhost:8080/api/users")
	fmt.Println("Generating some traffic...")

	// Wait for the example to complete
	time.Sleep(2 * time.Second)

	fmt.Println("Example completed. In a real application, metrics would be scraped by Prometheus.")
}