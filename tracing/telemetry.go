// Copyright (c) 2025 A Bit of Help, Inc.

// Package tracing provides functionality for distributed tracing in applications.
package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/abitofhelp/servicelib/logging"
)

// TelemetryProvider is a unified provider for tracing and metrics
type TelemetryProvider struct {
	// tracingProvider is the tracing provider
	tracingProvider *TracingProvider

	// metricsProvider is the metrics provider
	metricsProvider *MetricsProvider

	// logger is the logger for the telemetry provider
	logger *logging.ContextLogger
}

// TelemetryConfig holds configuration for telemetry
type TelemetryConfig struct {
	// Enabled indicates whether telemetry is enabled
	Enabled bool

	// ServiceName is the name of the service
	ServiceName string

	// Environment is the environment the service is running in
	Environment string

	// Version is the version of the service
	Version string

	// Tracing is the configuration for tracing
	Tracing struct {
		// Enabled indicates whether tracing is enabled
		Enabled bool

		// SamplingRatio is the ratio of traces to sample
		SamplingRatio float64

		// PropagationKeys are the keys to propagate in trace context
		PropagationKeys []string

		// OTLPEndpoint is the OTLP endpoint
		OTLPEndpoint string

		// OTLPInsecure indicates whether to use insecure connections
		OTLPInsecure bool

		// OTLPTimeout is the timeout for OTLP operations in seconds
		OTLPTimeout int

		// ShutdownTimeout is the timeout for shutting down tracing in seconds
		ShutdownTimeout int
	}

	// Metrics is the configuration for metrics
	Metrics struct {
		// Enabled indicates whether metrics are enabled
		Enabled bool

		// Prometheus is the configuration for Prometheus metrics
		Prometheus struct {
			// Enabled indicates whether Prometheus metrics are enabled
			Enabled bool

			// Path is the path to expose Prometheus metrics on
			Path string
		}
	}
}

// NewTelemetryProvider creates a new telemetry provider
//
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use for logging telemetry events
//   - config: The telemetry configuration
//
// Returns:
//   - *TelemetryProvider: The telemetry provider
//   - error: An error if the telemetry provider creation fails
func NewTelemetryProvider(ctx context.Context, logger *logging.ContextLogger, config TelemetryConfig) (*TelemetryProvider, error) {
	if !config.Enabled {
		logger.Info(ctx, "Telemetry is disabled")
		return nil, nil
	}

	// Create tracing provider
	tracingConfig := Config{
		Enabled:         config.Tracing.Enabled,
		ServiceName:     config.ServiceName,
		Environment:     config.Environment,
		Version:         config.Version,
		SamplingRatio:   config.Tracing.SamplingRatio,
		PropagationKeys: config.Tracing.PropagationKeys,
		OTLPEndpoint:    config.Tracing.OTLPEndpoint,
		OTLPInsecure:    config.Tracing.OTLPInsecure,
		OTLPTimeout:     config.Tracing.OTLPTimeout,
		ShutdownTimeout: config.Tracing.ShutdownTimeout,
	}

	tracingProvider, err := NewTracingProvider(ctx, logger, tracingConfig)
	if err != nil {
		return nil, err
	}

	// Create metrics provider
	metricsConfig := MetricsConfig{
		Enabled:          config.Metrics.Enabled,
		ServiceName:      config.ServiceName,
		Environment:      config.Environment,
		Version:          config.Version,
		PrometheusEnabled: config.Metrics.Prometheus.Enabled,
		PrometheusPath:   config.Metrics.Prometheus.Path,
	}

	// Create meter
	meter := otel.Meter("github.com/abitofhelp/servicelib/tracing")

	metricsProvider := NewMetricsProvider(ctx, logger, metricsConfig, meter)

	logger.Info(ctx, "Telemetry provider initialized",
		zap.String("service", config.ServiceName),
		zap.String("environment", config.Environment),
		zap.Bool("tracing_enabled", config.Tracing.Enabled),
		zap.Bool("metrics_enabled", config.Metrics.Enabled),
	)

	return &TelemetryProvider{
		tracingProvider: tracingProvider,
		metricsProvider: metricsProvider,
		logger:          logger,
	}, nil
}

// Shutdown shuts down the telemetry provider
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - error: An error if the shutdown fails
func (tp *TelemetryProvider) Shutdown(ctx context.Context) error {
	if tp.tracingProvider != nil {
		if err := tp.tracingProvider.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Tracer returns the tracer
//
// Returns:
//   - trace.Tracer: The OpenTelemetry tracer
func (tp *TelemetryProvider) Tracer() trace.Tracer {
	if tp.tracingProvider != nil {
		return tp.tracingProvider.Tracer()
	}
	return otel.Tracer("github.com/abitofhelp/servicelib/tracing")
}

// Meter returns the meter
//
// Returns:
//   - metric.Meter: The OpenTelemetry meter
func (tp *TelemetryProvider) Meter() metric.Meter {
	if tp.metricsProvider != nil {
		return tp.metricsProvider.Meter()
	}
	return otel.Meter("github.com/abitofhelp/servicelib/tracing")
}

// CreatePrometheusHandler creates a handler for the Prometheus metrics endpoint
//
// Returns:
//   - http.Handler: The Prometheus metrics handler
func (tp *TelemetryProvider) CreatePrometheusHandler() http.Handler {
	return CreatePrometheusHandler()
}

// InstrumentHandler wraps an http.Handler with OpenTelemetry instrumentation
//
// Parameters:
//   - handler: The HTTP handler to instrument
//   - operation: The name of the operation for tracing
//
// Returns:
//   - http.Handler: The instrumented HTTP handler
func (tp *TelemetryProvider) InstrumentHandler(handler http.Handler, operation string) http.Handler {
	return InstrumentHandler(handler, operation)
}

// NewHTTPMiddleware creates a new middleware for HTTP request tracing
//
// Returns:
//   - func(http.Handler) http.Handler: The middleware function
func (tp *TelemetryProvider) NewHTTPMiddleware() func(http.Handler) http.Handler {
	return NewHTTPMiddleware(tp.logger)
}