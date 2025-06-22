// Copyright (c) 2025 A Bit of Help, Inc.

// Example of initializing the telemetry package
package example_telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
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

	// Load configuration
	k := koanf.New(".")

	// In a real application, you would load configuration from a file
	// For this example, we'll use a hardcoded YAML configuration
	configYAML := []byte(`
telemetry:
  enabled: true
  service_name: "example-service"
  environment: "development"
  version: "1.0.0"
  shutdown_timeout: 5
  
  otlp:
    endpoint: "localhost:4317"
    insecure: true
    timeout_seconds: 5
  
  tracing:
    enabled: true
    sampling_ratio: 1.0
    propagation_keys:
      - "traceparent"
      - "tracestate"
      - "baggage"
  
  metrics:
    enabled: true
    reporting_frequency_seconds: 15
    prometheus:
      enabled: true
      listen: "0.0.0.0:8089"
      path: "/metrics"
  
  http:
    tracing_enabled: true
`)

	// Load the configuration
	if err := k.Load(rawbytes.Provider(configYAML), yaml.Parser()); err != nil {
		logger.Fatal(ctx, "Failed to load configuration", zap.Error(err))
	}

	// Create telemetry provider
	telemetryProvider, err := telemetry.NewTelemetryProvider(ctx, logger, k)
	if err != nil {
		logger.Fatal(ctx, "Failed to create telemetry provider", zap.Error(err))
	}

	// Set up graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Log that the telemetry provider is initialized
	logger.Info(ctx, "Telemetry provider initialized successfully")

	// In a real application, you would use the telemetry provider here
	// For this example, we'll just wait for a shutdown signal or timeout
	select {
	case <-shutdown:
		logger.Info(ctx, "Received shutdown signal")
	case <-time.After(10 * time.Second):
		logger.Info(ctx, "Example timeout reached")
	}

	// Shutdown the telemetry provider
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := telemetryProvider.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "Error shutting down telemetry provider", zap.Error(err))
	} else {
		logger.Info(ctx, "Telemetry provider shut down successfully")
	}

	fmt.Println("Example completed. In a real application, the telemetry provider would be used for tracing and metrics.")
}
