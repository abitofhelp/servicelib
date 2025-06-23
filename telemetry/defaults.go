// Copyright (c) 2025 A Bit of Help, Inc.

package telemetry

// DefaultConfig returns a default configuration for telemetry.
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

// DefaultOTLPConfig returns a default configuration for OTLP exporter.
func DefaultOTLPConfig() OTLPConfig {
	return OTLPConfig{
		Endpoint: "localhost:4317",
		Insecure: true,
		Timeout:  5,
	}
}

// DefaultTracingConfig returns a default configuration for tracing.
func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		Enabled:         true,
		SamplingRatio:   1.0,
		PropagationKeys: []string{"traceparent", "tracestate", "baggage"},
	}
}

// DefaultMetricsConfig returns a default configuration for metrics.
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:       true,
		ReportingFreq: 15,
		Prometheus:    DefaultPrometheusConfig(),
	}
}

// DefaultPrometheusConfig returns a default configuration for Prometheus metrics.
func DefaultPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		Enabled: true,
		Listen:  "0.0.0.0:8089",
		Path:    "/metrics",
	}
}

// DefaultHTTPConfig returns a default configuration for HTTP telemetry.
func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		TracingEnabled: true,
	}
}
