// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for collecting and exporting telemetry data,
// including distributed tracing and metrics.
//
// This file contains the configuration structures and loading functions for the telemetry
// package. It defines the configuration schema for all telemetry components, including
// tracing, metrics, and OTLP exporters. The configuration can be loaded from a koanf
// instance, and default values are provided for quick setup in development environments.
package telemetry

import (
	"github.com/knadh/koanf/v2"
)

// Config holds all telemetry configuration
type Config struct {
	// Enabled indicates whether telemetry is enabled
	Enabled bool `mapstructure:"enabled"`

	// ServiceName is the name of the service
	ServiceName string `mapstructure:"service_name"`

	// Environment is the environment the service is running in
	Environment string `mapstructure:"environment"`

	// Version is the version of the service
	Version string `mapstructure:"version"`

	// ShutdownTimeout is the timeout for shutting down telemetry in seconds
	ShutdownTimeout int `mapstructure:"shutdown_timeout"`

	// OTLP is the configuration for the OTLP exporter
	OTLP OTLPConfig `mapstructure:"otlp"`

	// Tracing is the configuration for tracing
	Tracing TracingConfig `mapstructure:"tracing"`

	// Metrics is the configuration for metrics
	Metrics MetricsConfig `mapstructure:"metrics"`

	// HTTP is the configuration for HTTP telemetry
	HTTP HTTPConfig `mapstructure:"http"`
}

// OTLPConfig holds configuration for OTLP exporter
type OTLPConfig struct {
	// Endpoint is the OTLP endpoint
	Endpoint string `mapstructure:"endpoint"`

	// Insecure indicates whether to use insecure connections
	Insecure bool `mapstructure:"insecure"`

	// Timeout is the timeout for OTLP operations in seconds
	Timeout int `mapstructure:"timeout_seconds"`
}

// TracingConfig holds configuration for tracing
type TracingConfig struct {
	// Enabled indicates whether tracing is enabled
	Enabled bool `mapstructure:"enabled"`

	// SamplingRatio is the ratio of traces to sample
	SamplingRatio float64 `mapstructure:"sampling_ratio"`

	// PropagationKeys are the keys to propagate in trace context
	PropagationKeys []string `mapstructure:"propagation_keys"`
}

// MetricsConfig holds configuration for metrics
type MetricsConfig struct {
	// Enabled indicates whether metrics are enabled
	Enabled bool `mapstructure:"enabled"`

	// ReportingFreq is the frequency of metrics reporting in seconds
	ReportingFreq int `mapstructure:"reporting_frequency_seconds"`

	// Prometheus is the configuration for Prometheus metrics
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

// PrometheusConfig holds configuration for Prometheus metrics
type PrometheusConfig struct {
	// Enabled indicates whether Prometheus metrics are enabled
	Enabled bool `mapstructure:"enabled"`

	// Listen is the address to listen on for Prometheus metrics
	Listen string `mapstructure:"listen"`

	// Path is the path to expose Prometheus metrics on
	Path string `mapstructure:"path"`
}

// HTTPConfig holds configuration for HTTP telemetry
type HTTPConfig struct {
	// TracingEnabled indicates whether HTTP tracing is enabled
	TracingEnabled bool `mapstructure:"tracing_enabled"`
}

// LoadConfig loads telemetry configuration from koanf
//
// Parameters:
//   - k: The koanf instance to load configuration from
//
// Returns:
//   - Config: The loaded telemetry configuration
func LoadConfig(k *koanf.Koanf) Config {
	return Config{
		Enabled:         k.Bool("telemetry.enabled"),
		ServiceName:     k.String("telemetry.service_name"),
		Environment:     k.String("telemetry.environment"),
		Version:         k.String("telemetry.version"),
		ShutdownTimeout: k.Int("telemetry.shutdown_timeout"),
		OTLP: OTLPConfig{
			Endpoint: k.String("telemetry.otlp.endpoint"),
			Insecure: k.Bool("telemetry.otlp.insecure"),
			Timeout:  k.Int("telemetry.otlp.timeout_seconds"),
		},
		Tracing: TracingConfig{
			Enabled:         k.Bool("telemetry.tracing.enabled"),
			SamplingRatio:   k.Float64("telemetry.tracing.sampling_ratio"),
			PropagationKeys: k.Strings("telemetry.tracing.propagation_keys"),
		},
		Metrics: MetricsConfig{
			Enabled:       k.Bool("telemetry.metrics.enabled"),
			ReportingFreq: k.Int("telemetry.metrics.reporting_frequency_seconds"),
			Prometheus: PrometheusConfig{
				Enabled: k.Bool("telemetry.exporters.metrics.prometheus.enabled"),
				Listen:  k.String("telemetry.exporters.metrics.prometheus.listen"),
				Path:    k.String("telemetry.exporters.metrics.prometheus.path"),
			},
		},
		HTTP: HTTPConfig{
			TracingEnabled: k.Bool("telemetry.http.tracing_enabled"),
		},
	}
}

// GetTelemetryDefaults returns default values for telemetry configuration
//
// Parameters:
//   - serviceName: The name of the service (optional, defaults to "service")
//
// Returns:
//   - map[string]interface{}: Default values for telemetry configuration
func GetTelemetryDefaults(serviceName ...string) map[string]interface{} {
	// Use the provided service name or default to "service"
	name := "service"
	if len(serviceName) > 0 && serviceName[0] != "" {
		name = serviceName[0]
	}

	return map[string]interface{}{
		"telemetry.enabled":          true,
		"telemetry.service_name":     name,
		"telemetry.environment":      "development",
		"telemetry.version":          "1.0.0",
		"telemetry.shutdown_timeout": 5,

		// OTLP defaults
		"telemetry.otlp.endpoint":        "localhost:4317",
		"telemetry.otlp.insecure":        true,
		"telemetry.otlp.timeout_seconds": 5,

		// Tracing defaults
		"telemetry.tracing.enabled":          true,
		"telemetry.tracing.sampling_ratio":   1.0, // Sample everything by default
		"telemetry.tracing.propagation_keys": []string{"traceparent", "tracestate", "baggage"},

		// Metrics defaults
		"telemetry.metrics.enabled":                      true,
		"telemetry.metrics.reporting_frequency_seconds":  15,
		"telemetry.exporters.metrics.prometheus.enabled": true,
		"telemetry.exporters.metrics.prometheus.listen":  "0.0.0.0:8089",
		"telemetry.exporters.metrics.prometheus.path":    "/metrics",

		// HTTP defaults
		"telemetry.http.tracing_enabled": true,
	}
}
