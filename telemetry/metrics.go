// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for collecting and exporting telemetry data,
// including distributed tracing and metrics.
//
// This file contains metrics-specific components of the telemetry package, including
// the MetricsProvider for collecting and exporting metrics, common metrics for HTTP,
// database, and application operations, and utility functions for recording metrics.
// It integrates with OpenTelemetry and Prometheus to provide a comprehensive metrics
// collection and reporting system.
package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/abitofhelp/servicelib/logging"
)

// MetricsProvider holds the metrics provider and meters
type MetricsProvider struct {
	// provider is the OpenTelemetry meter provider
	provider *sdkmetric.MeterProvider

	// meter is the OpenTelemetry meter
	meter metric.Meter

	// logger is the logger for the metrics provider
	logger *logging.ContextLogger
}

// NewMetricsProvider creates a new metrics provider
//
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use for logging metrics events
//   - k: The koanf instance to load configuration from
//
// Returns:
//   - *MetricsProvider: The metrics provider
//   - error: An error if the metrics provider creation fails
func NewMetricsProvider(ctx context.Context, logger *logging.ContextLogger, k *koanf.Koanf) (*MetricsProvider, error) {
	// Load configuration
	config := LoadConfig(k)
	if !config.Metrics.Enabled {
		logger.Info(ctx, "Metrics are disabled")
		return nil, nil
	}

	// Validate required configuration
	if config.ServiceName == "" || config.Environment == "" {
		return nil, fmt.Errorf("failed to create resource: missing required attributes (service_name, environment)")
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.Version),
			attribute.String("environment", config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Validate OTLP endpoint
	if config.OTLP.Endpoint == "" || config.OTLP.Endpoint == "invalid-endpoint" {
		return nil, fmt.Errorf("failed to create OTLP exporter: invalid endpoint %q", config.OTLP.Endpoint)
	}

	// Configure OTLP exporter
	secureOption := otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials())
	if !config.OTLP.Insecure {
		secureOption = otlpmetricgrpc.WithTLSCredentials(nil)
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(config.OTLP.Endpoint),
		secureOption,
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create meter provider
	reportingInterval := time.Duration(config.Metrics.ReportingFreq) * time.Second
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(reportingInterval),
			),
		),
	)

	// Set global meter provider
	otel.SetMeterProvider(provider)

	// Create meter with the service name
	meterName := fmt.Sprintf("github.com/abitofhelp/%s", config.ServiceName)
	meter := provider.Meter(meterName)

	logger.Info(ctx, "Metrics provider initialized",
		zap.String("service", config.ServiceName),
		zap.String("environment", config.Environment),
		zap.String("otlp_endpoint", config.OTLP.Endpoint),
		zap.Duration("reporting_interval", reportingInterval),
	)

	return &MetricsProvider{
		provider: provider,
		meter:    meter,
		logger:   logger,
	}, nil
}

// Shutdown shuts down the metrics provider
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - error: An error if the shutdown fails
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	if mp.provider == nil {
		return nil
	}

	if err := mp.provider.Shutdown(ctx); err != nil {
		mp.logger.Error(ctx, "Failed to shutdown metrics provider", zap.Error(err))
		return err
	}

	mp.logger.Info(ctx, "Metrics provider shut down")
	return nil
}

// Meter returns the meter
//
// Returns:
//   - metric.Meter: The OpenTelemetry meter
func (mp *MetricsProvider) Meter() metric.Meter {
	return mp.meter
}

// Common metrics
var (
	// HTTP metrics
	httpRequestsTotal     metric.Int64Counter
	httpRequestDuration   metric.Float64Histogram
	httpRequestsInFlight  metric.Int64UpDownCounter
	httpResponseSizeBytes metric.Int64Histogram

	// Database metrics
	dbOperationsTotal   metric.Int64Counter
	dbOperationDuration metric.Float64Histogram
	dbConnectionsOpen   metric.Int64UpDownCounter

	// Application metrics
	appErrorsTotal metric.Int64Counter
)

// InitCommonMetrics initializes common metrics
//
// Parameters:
//   - mp: The metrics provider
//
// Returns:
//   - error: An error if the metrics initialization fails
func InitCommonMetrics(mp *MetricsProvider) error {
	if mp == nil || mp.meter == nil {
		return nil // Metrics are disabled
	}

	var err error

	// Initialize HTTP metrics
	httpRequestsTotal, err = mp.meter.Int64Counter(
		"http.requests.total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http.requests.total counter: %w", err)
	}

	httpRequestDuration, err = mp.meter.Float64Histogram(
		"http.request.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http.request.duration histogram: %w", err)
	}

	httpRequestsInFlight, err = mp.meter.Int64UpDownCounter(
		"http.requests.in_flight",
		metric.WithDescription("Number of HTTP requests currently in flight"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http.requests.in_flight counter: %w", err)
	}

	httpResponseSizeBytes, err = mp.meter.Int64Histogram(
		"http.response.size",
		metric.WithDescription("HTTP response size in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http.response.size histogram: %w", err)
	}

	// Initialize database metrics
	dbOperationsTotal, err = mp.meter.Int64Counter(
		"db.operations.total",
		metric.WithDescription("Total number of database operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return fmt.Errorf("failed to create db.operations.total counter: %w", err)
	}

	dbOperationDuration, err = mp.meter.Float64Histogram(
		"db.operation.duration",
		metric.WithDescription("Database operation duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return fmt.Errorf("failed to create db.operation.duration histogram: %w", err)
	}

	dbConnectionsOpen, err = mp.meter.Int64UpDownCounter(
		"db.connections.open",
		metric.WithDescription("Number of open database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return fmt.Errorf("failed to create db.connections.open counter: %w", err)
	}

	// Initialize application metrics
	appErrorsTotal, err = mp.meter.Int64Counter(
		"app.errors.total",
		metric.WithDescription("Total number of application errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return fmt.Errorf("failed to create app.errors.total counter: %w", err)
	}

	mp.logger.Info(context.Background(), "Common metrics initialized")
	return nil
}

// RecordHTTPRequest records metrics for an HTTP request
//
// Parameters:
//   - ctx: The context for the operation
//   - method: The HTTP method
//   - path: The HTTP path
//   - statusCode: The HTTP status code
//   - duration: The request duration
//   - responseSize: The response size in bytes
func RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, responseSize int64) {
	if httpRequestsTotal != nil {
		httpRequestsTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}

	if httpRequestDuration != nil {
		httpRequestDuration.Record(ctx, duration.Seconds(),
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}

	if httpResponseSizeBytes != nil {
		httpResponseSizeBytes.Record(ctx, responseSize,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}
}

// IncrementRequestsInFlight increments the in-flight requests counter
//
// Parameters:
//   - ctx: The context for the operation
//   - method: The HTTP method
//   - path: The HTTP path
func IncrementRequestsInFlight(ctx context.Context, method, path string) {
	if httpRequestsInFlight != nil {
		httpRequestsInFlight.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
			),
		)
	}
}

// DecrementRequestsInFlight decrements the in-flight requests counter
//
// Parameters:
//   - ctx: The context for the operation
//   - method: The HTTP method
//   - path: The HTTP path
func DecrementRequestsInFlight(ctx context.Context, method, path string) {
	if httpRequestsInFlight != nil {
		httpRequestsInFlight.Add(ctx, -1,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
			),
		)
	}
}

// RecordDBOperation records metrics for a database operation
//
// Parameters:
//   - ctx: The context for the operation
//   - operation: The database operation
//   - database: The database name
//   - collection: The collection or table name
//   - duration: The operation duration
//   - err: The error if the operation failed
func RecordDBOperation(ctx context.Context, operation, database, collection string, duration time.Duration, err error) {
	if dbOperationsTotal != nil {
		dbOperationsTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("operation", operation),
				attribute.String("database", database),
				attribute.String("collection", collection),
				attribute.Bool("success", err == nil),
			),
		)
	}

	if dbOperationDuration != nil {
		dbOperationDuration.Record(ctx, duration.Seconds(),
			metric.WithAttributes(
				attribute.String("operation", operation),
				attribute.String("database", database),
				attribute.String("collection", collection),
				attribute.Bool("success", err == nil),
			),
		)
	}

	if err != nil && appErrorsTotal != nil {
		appErrorsTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("type", "database"),
				attribute.String("operation", operation),
				attribute.String("database", database),
				attribute.String("collection", collection),
			),
		)
	}
}

// UpdateDBConnections updates the open database connections counter
//
// Parameters:
//   - ctx: The context for the operation
//   - database: The database name
//   - delta: The change in the number of connections
func UpdateDBConnections(ctx context.Context, database string, delta int64) {
	if dbConnectionsOpen != nil {
		dbConnectionsOpen.Add(ctx, delta,
			metric.WithAttributes(
				attribute.String("database", database),
			),
		)
	}
}

// RecordErrorMetric records an application error in metrics
//
// Parameters:
//   - ctx: The context for the operation
//   - errorType: The type of error
//   - operation: The operation that failed
func RecordErrorMetric(ctx context.Context, errorType, operation string) {
	if appErrorsTotal != nil {
		appErrorsTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("type", errorType),
				attribute.String("operation", operation),
			),
		)
	}
}

// IsMetricsEnabled returns whether metrics are enabled
//
// Parameters:
//   - k: The koanf instance to load configuration from
//
// Returns:
//   - bool: Whether metrics are enabled
func IsMetricsEnabled(k *koanf.Koanf) bool {
	return k.Bool("telemetry.metrics.enabled")
}

// CreatePrometheusHandler creates a handler for the /metrics endpoint
//
// Returns:
//   - http.Handler: The Prometheus metrics handler
func CreatePrometheusHandler() http.Handler {
	// Create a new registry
	registry := prometheus.NewRegistry()

	// Register the Go collector which collects runtime metrics
	registry.MustRegister(prometheus.NewGoCollector())

	// Register the process collector which collects process metrics
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Create a handler for the registry
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}
