// Copyright (c) 2025 A Bit of Help, Inc.

// Package tracing provides functionality for distributed tracing in applications.
package tracing

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/abitofhelp/servicelib/logging"
)

// MetricsConfig holds configuration for metrics
type MetricsConfig struct {
	// Enabled indicates whether metrics are enabled
	Enabled bool

	// ServiceName is the name of the service
	ServiceName string

	// Environment is the environment the service is running in
	Environment string

	// Version is the version of the service
	Version string

	// PrometheusEnabled indicates whether Prometheus metrics are enabled
	PrometheusEnabled bool

	// PrometheusPath is the path to expose Prometheus metrics on
	PrometheusPath string
}

// MetricsProvider holds the metrics provider and meters
type MetricsProvider struct {
	// meter is the OpenTelemetry meter
	meter metric.Meter

	// logger is the logger for the metrics provider
	logger *logging.ContextLogger

	// config is the metrics configuration
	config MetricsConfig
}

// NewMetricsProvider creates a new metrics provider
//
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use for logging metrics events
//   - config: The metrics configuration
//   - meter: The OpenTelemetry meter
//
// Returns:
//   - *MetricsProvider: The metrics provider
func NewMetricsProvider(ctx context.Context, logger *logging.ContextLogger, config MetricsConfig, meter metric.Meter) *MetricsProvider {
	if !config.Enabled {
		logger.Info(ctx, "Metrics are disabled")
		return nil
	}

	logger.Info(ctx, "Metrics provider initialized",
		zap.String("service", config.ServiceName),
		zap.String("environment", config.Environment),
		zap.Bool("prometheus_enabled", config.PrometheusEnabled),
	)

	return &MetricsProvider{
		meter:  meter,
		logger: logger,
		config: config,
	}
}

// Meter returns the meter
//
// Returns:
//   - metric.Meter: The OpenTelemetry meter
func (mp *MetricsProvider) Meter() metric.Meter {
	return mp.meter
}

// CreateCounter creates a new counter metric
//
// Parameters:
//   - name: The name of the counter
//   - description: The description of the counter
//   - unit: The unit of the counter
//
// Returns:
//   - metric.Int64Counter: The counter
//   - error: An error if the counter creation fails
func (mp *MetricsProvider) CreateCounter(name, description, unit string) (metric.Int64Counter, error) {
	return mp.meter.Int64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// CreateHistogram creates a new histogram metric
//
// Parameters:
//   - name: The name of the histogram
//   - description: The description of the histogram
//   - unit: The unit of the histogram
//
// Returns:
//   - metric.Float64Histogram: The histogram
//   - error: An error if the histogram creation fails
func (mp *MetricsProvider) CreateHistogram(name, description, unit string) (metric.Float64Histogram, error) {
	return mp.meter.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// CreateUpDownCounter creates a new up/down counter metric
//
// Parameters:
//   - name: The name of the counter
//   - description: The description of the counter
//   - unit: The unit of the counter
//
// Returns:
//   - metric.Int64UpDownCounter: The counter
//   - error: An error if the counter creation fails
func (mp *MetricsProvider) CreateUpDownCounter(name, description, unit string) (metric.Int64UpDownCounter, error) {
	return mp.meter.Int64UpDownCounter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// RecordDuration records a duration metric
//
// Parameters:
//   - ctx: The context for the operation
//   - histogram: The histogram to record to
//   - duration: The duration to record
//   - attrs: The attributes to record with the duration
func RecordDuration(ctx context.Context, histogram metric.Float64Histogram, duration time.Duration, attrs ...attribute.KeyValue) {
	if histogram != nil {
		histogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
	}
}

// IncrementCounter increments a counter metric
//
// Parameters:
//   - ctx: The context for the operation
//   - counter: The counter to increment
//   - value: The value to increment by
//   - attrs: The attributes to record with the increment
func IncrementCounter(ctx context.Context, counter metric.Int64Counter, value int64, attrs ...attribute.KeyValue) {
	if counter != nil {
		counter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
}

// UpdateUpDownCounter updates an up/down counter metric
//
// Parameters:
//   - ctx: The context for the operation
//   - counter: The counter to update
//   - value: The value to update by
//   - attrs: The attributes to record with the update
func UpdateUpDownCounter(ctx context.Context, counter metric.Int64UpDownCounter, value int64, attrs ...attribute.KeyValue) {
	if counter != nil {
		counter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
}

// CreatePrometheusHandler creates a handler for the Prometheus metrics endpoint
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