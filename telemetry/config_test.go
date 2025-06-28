// Copyright (c) 2025 A Bit of Help, Inc.

package telemetry

import (
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetTelemetryDefaults(t *testing.T) {
	// Test with default service name
	defaults := GetTelemetryDefaults()
	assert.Equal(t, "service", defaults["telemetry.service_name"])
	assert.Equal(t, true, defaults["telemetry.enabled"])
	assert.Equal(t, "development", defaults["telemetry.environment"])
	assert.Equal(t, "1.0.0", defaults["telemetry.version"])
	assert.Equal(t, 5, defaults["telemetry.shutdown_timeout"])
	assert.Equal(t, "localhost:4317", defaults["telemetry.otlp.endpoint"])
	assert.Equal(t, true, defaults["telemetry.otlp.insecure"])
	assert.Equal(t, 5, defaults["telemetry.otlp.timeout_seconds"])
	assert.Equal(t, true, defaults["telemetry.tracing.enabled"])
	assert.Equal(t, 1.0, defaults["telemetry.tracing.sampling_ratio"])
	assert.Equal(t, []string{"traceparent", "tracestate", "baggage"}, defaults["telemetry.tracing.propagation_keys"])
	assert.Equal(t, true, defaults["telemetry.metrics.enabled"])
	assert.Equal(t, 15, defaults["telemetry.metrics.reporting_frequency_seconds"])
	assert.Equal(t, true, defaults["telemetry.exporters.metrics.prometheus.enabled"])
	assert.Equal(t, "0.0.0.0:8089", defaults["telemetry.exporters.metrics.prometheus.listen"])
	assert.Equal(t, "/metrics", defaults["telemetry.exporters.metrics.prometheus.path"])
	assert.Equal(t, true, defaults["telemetry.http.tracing_enabled"])

	// Test with custom service name
	customDefaults := GetTelemetryDefaults("custom-service")
	assert.Equal(t, "custom-service", customDefaults["telemetry.service_name"])
}

func TestLoadConfig(t *testing.T) {
	// Create a new koanf instance
	k := koanf.New(".")

	// Load default values
	defaults := GetTelemetryDefaults("test-service")
	for key, value := range defaults {
		k.Set(key, value)
	}

	// Load configuration
	config := LoadConfig(k)

	// Verify configuration values
	assert.Equal(t, true, config.Enabled)
	assert.Equal(t, "test-service", config.ServiceName)
	assert.Equal(t, "development", config.Environment)
	assert.Equal(t, "1.0.0", config.Version)
	assert.Equal(t, 5, config.ShutdownTimeout)

	// Verify OTLP configuration
	assert.Equal(t, "localhost:4317", config.OTLP.Endpoint)
	assert.Equal(t, true, config.OTLP.Insecure)
	assert.Equal(t, 5, config.OTLP.Timeout)

	// Verify tracing configuration
	assert.Equal(t, true, config.Tracing.Enabled)
	assert.Equal(t, 1.0, config.Tracing.SamplingRatio)
	assert.Equal(t, []string{"traceparent", "tracestate", "baggage"}, config.Tracing.PropagationKeys)

	// Verify metrics configuration
	assert.Equal(t, true, config.Metrics.Enabled)
	assert.Equal(t, 15, config.Metrics.ReportingFreq)
	assert.Equal(t, true, config.Metrics.Prometheus.Enabled)
	assert.Equal(t, "0.0.0.0:8089", config.Metrics.Prometheus.Listen)
	assert.Equal(t, "/metrics", config.Metrics.Prometheus.Path)

	// Verify HTTP configuration
	assert.Equal(t, true, config.HTTP.TracingEnabled)

	// Test with custom values
	k.Set("telemetry.enabled", false)
	k.Set("telemetry.service_name", "custom-service")
	k.Set("telemetry.environment", "production")
	k.Set("telemetry.version", "2.0.0")
	k.Set("telemetry.shutdown_timeout", 10)
	k.Set("telemetry.otlp.endpoint", "otlp.example.com:4317")
	k.Set("telemetry.otlp.insecure", false)
	k.Set("telemetry.otlp.timeout_seconds", 15)
	k.Set("telemetry.tracing.enabled", false)
	k.Set("telemetry.tracing.sampling_ratio", 0.5)
	k.Set("telemetry.tracing.propagation_keys", []string{"custom-key"})
	k.Set("telemetry.metrics.enabled", false)
	k.Set("telemetry.metrics.reporting_frequency_seconds", 30)
	k.Set("telemetry.exporters.metrics.prometheus.enabled", false)
	k.Set("telemetry.exporters.metrics.prometheus.listen", "127.0.0.1:9090")
	k.Set("telemetry.exporters.metrics.prometheus.path", "/custom-metrics")
	k.Set("telemetry.http.tracing_enabled", false)

	// Load custom configuration
	customConfig := LoadConfig(k)

	// Verify custom configuration values
	assert.Equal(t, false, customConfig.Enabled)
	assert.Equal(t, "custom-service", customConfig.ServiceName)
	assert.Equal(t, "production", customConfig.Environment)
	assert.Equal(t, "2.0.0", customConfig.Version)
	assert.Equal(t, 10, customConfig.ShutdownTimeout)

	// Verify custom OTLP configuration
	assert.Equal(t, "otlp.example.com:4317", customConfig.OTLP.Endpoint)
	assert.Equal(t, false, customConfig.OTLP.Insecure)
	assert.Equal(t, 15, customConfig.OTLP.Timeout)

	// Verify custom tracing configuration
	assert.Equal(t, false, customConfig.Tracing.Enabled)
	assert.Equal(t, 0.5, customConfig.Tracing.SamplingRatio)
	assert.Equal(t, []string{"custom-key"}, customConfig.Tracing.PropagationKeys)

	// Verify custom metrics configuration
	assert.Equal(t, false, customConfig.Metrics.Enabled)
	assert.Equal(t, 30, customConfig.Metrics.ReportingFreq)
	assert.Equal(t, false, customConfig.Metrics.Prometheus.Enabled)
	assert.Equal(t, "127.0.0.1:9090", customConfig.Metrics.Prometheus.Listen)
	assert.Equal(t, "/custom-metrics", customConfig.Metrics.Prometheus.Path)

	// Verify custom HTTP configuration
	assert.Equal(t, false, customConfig.HTTP.TracingEnabled)
}

func TestIsMetricsEnabled(t *testing.T) {
	// Create a new koanf instance
	k := koanf.New(".")

	// Test with metrics enabled
	k.Set("telemetry.metrics.enabled", true)
	assert.True(t, IsMetricsEnabled(k))

	// Test with metrics disabled
	k.Set("telemetry.metrics.enabled", false)
	assert.False(t, IsMetricsEnabled(k))
}
