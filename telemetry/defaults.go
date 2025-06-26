// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for collecting and exporting telemetry data,
// including distributed tracing and metrics.
//
// This file contains default configuration functions for the telemetry package.
// These functions provide sensible default values for all telemetry components,
// making it easy to get started with telemetry without having to specify every
// configuration option. The defaults are suitable for development environments
// and can be overridden as needed for production use.
package telemetry

// DefaultConfig returns a default configuration for telemetry.
// This function provides a complete default configuration for all telemetry components,
// including tracing, metrics, and HTTP telemetry. The defaults are suitable for
// development environments and include:
//   - Service name: "service"
//   - Environment: "development"
//   - Version: "1.0.0"
//   - Enabled telemetry with a 5-second shutdown timeout
//   - Default OTLP, tracing, metrics, and HTTP configurations
//
// Returns:
//   - Config: A complete default configuration for telemetry
func DefaultConfig() Config {
	return Config{
		Enabled:         true,
		ServiceName:     "service",
		Environment:     "development",
		Version:         "1.0.0",
		ShutdownTimeout: 5,
		OTLP:            DefaultOTLPConfig(),
		Tracing:         DefaultTracingConfig(),
		Metrics:         DefaultMetricsConfig(),
		HTTP:            DefaultHTTPConfig(),
	}
}

// DefaultOTLPConfig returns a default configuration for the OTLP exporter.
// This function provides default values for the OpenTelemetry Protocol (OTLP) exporter,
// which is used to send telemetry data to an OTLP-compatible backend. The defaults include:
//   - Endpoint: "localhost:4317" (standard OTLP gRPC endpoint)
//   - Insecure: true (no TLS for development environments)
//   - Timeout: 5 seconds for OTLP operations
//
// These defaults are suitable for development with a local OpenTelemetry Collector.
// For production environments, you should configure a secure endpoint and adjust the timeout.
//
// Returns:
//   - OTLPConfig: A default configuration for the OTLP exporter
func DefaultOTLPConfig() OTLPConfig {
	return OTLPConfig{
		Endpoint: "localhost:4317",
		Insecure: true,
		Timeout:  5,
	}
}

// DefaultTracingConfig returns a default configuration for distributed tracing.
// This function provides default values for the tracing configuration, which controls
// how traces are sampled and propagated. The defaults include:
//   - Enabled: true (tracing is enabled)
//   - SamplingRatio: 1.0 (100% of traces are sampled)
//   - PropagationKeys: ["traceparent", "tracestate", "baggage"] (standard W3C trace context keys)
//
// These defaults are suitable for development environments where you want to see all traces.
// For production environments with high traffic, you may want to reduce the sampling ratio
// to avoid generating too much telemetry data.
//
// Returns:
//   - TracingConfig: A default configuration for distributed tracing
func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		Enabled:         true,
		SamplingRatio:   1.0,
		PropagationKeys: []string{"traceparent", "tracestate", "baggage"},
	}
}

// DefaultMetricsConfig returns a default configuration for metrics collection.
// This function provides default values for the metrics configuration, which controls
// how metrics are collected and reported. The defaults include:
//   - Enabled: true (metrics collection is enabled)
//   - ReportingFreq: 15 seconds (metrics are reported every 15 seconds)
//   - Prometheus: Default Prometheus configuration (see DefaultPrometheusConfig)
//
// These defaults are suitable for most environments. You may want to adjust the
// reporting frequency based on your monitoring needs and system load.
//
// Returns:
//   - MetricsConfig: A default configuration for metrics collection
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:       true,
		ReportingFreq: 15,
		Prometheus:    DefaultPrometheusConfig(),
	}
}

// DefaultPrometheusConfig returns a default configuration for Prometheus metrics.
// This function provides default values for the Prometheus metrics configuration,
// which controls how metrics are exposed for scraping by Prometheus. The defaults include:
//   - Enabled: true (Prometheus metrics endpoint is enabled)
//   - Listen: "0.0.0.0:8089" (listen on all interfaces, port 8089)
//   - Path: "/metrics" (standard Prometheus metrics endpoint path)
//
// These defaults allow Prometheus to scrape metrics from the application on the
// standard /metrics endpoint. You may need to adjust the listen address and port
// based on your network configuration and security requirements.
//
// Returns:
//   - PrometheusConfig: A default configuration for Prometheus metrics
func DefaultPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		Enabled: true,
		Listen:  "0.0.0.0:8089",
		Path:    "/metrics",
	}
}

// DefaultHTTPConfig returns a default configuration for HTTP telemetry.
// This function provides default values for the HTTP telemetry configuration,
// which controls how HTTP requests and responses are traced and monitored.
// The defaults include:
//   - TracingEnabled: true (HTTP request tracing is enabled)
//
// With these defaults, all HTTP requests will be traced, allowing you to see
// the full request flow in your distributed tracing system. This is useful for
// debugging and monitoring HTTP-based services.
//
// Returns:
//   - HTTPConfig: A default configuration for HTTP telemetry
func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		TracingEnabled: true,
	}
}
