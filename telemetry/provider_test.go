package telemetry

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewTelemetryProvider(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a koanf instance with telemetry disabled
	k := koanf.New(".")
	k.Set("telemetry.enabled", false)

	// Create a telemetry provider with telemetry disabled
	provider, err := NewTelemetryProvider(context.Background(), logger, k)
	assert.NoError(t, err)
	assert.Nil(t, provider)

	// Create a koanf instance with telemetry enabled but metrics disabled
	k = koanf.New(".")
	k.Set("telemetry.enabled", true)
	k.Set("telemetry.service_name", "test-service")
	k.Set("telemetry.environment", "test")
	k.Set("telemetry.version", "1.0.0")
	k.Set("telemetry.metrics.enabled", false)
	k.Set("telemetry.tracing.enabled", false)

	// Skip creating a telemetry provider with telemetry enabled as it would try to connect to a real OTLP endpoint
}

func TestTelemetryProviderShutdown(t *testing.T) {
	// Create a test provider
	provider := &TelemetryProvider{
		metricsProvider: nil,
		tracingProvider: nil,
		logger:          logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Test shutdown with nil metrics provider
	err := provider.Shutdown(context.Background())
	assert.NoError(t, err)

	// Create a metrics provider that returns an error on shutdown
	metricsProvider := &MetricsProvider{
		provider: nil,
		meter:    nil,
		logger:   logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Test shutdown with metrics provider
	provider.metricsProvider = metricsProvider
	err = provider.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestTelemetryProviderMeter(t *testing.T) {
	// Create a test provider with nil metrics provider
	provider := &TelemetryProvider{
		metricsProvider: nil,
		tracingProvider: nil,
		logger:          logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Test meter with nil metrics provider
	meter := provider.Meter()
	assert.NotNil(t, meter)

	// We can't easily test with a non-nil metrics provider because
	// we can't create a real meter in a test environment
}

func TestTelemetryProviderTracer(t *testing.T) {
	// Create a test provider with nil tracing provider
	provider := &TelemetryProvider{
		metricsProvider: nil,
		tracingProvider: nil,
		logger:          logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Test tracer with nil tracing provider
	tracer := provider.Tracer()
	assert.NotNil(t, tracer)

	// We can't easily test with a non-nil tracing provider because
	// we can't create a real tracer in a test environment
}

func TestTelemetryProviderCreatePrometheusHandler(t *testing.T) {
	// Create a test provider
	provider := &TelemetryProvider{
		metricsProvider: nil,
		tracingProvider: nil,
		logger:          logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Test create prometheus handler
	handler := provider.CreatePrometheusHandler()
	assert.NotNil(t, handler)
}

func TestTelemetryProviderInstrumentHandler(t *testing.T) {
	// Create a test provider
	provider := &TelemetryProvider{
		metricsProvider: nil,
		tracingProvider: nil,
		logger:          logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test instrument handler
	instrumentedHandler := provider.InstrumentHandler(handler, "test-operation")
	assert.NotNil(t, instrumentedHandler)
}

func TestTelemetryProviderNewHTTPMiddleware(t *testing.T) {
	// Create a test provider
	provider := &TelemetryProvider{
		metricsProvider: nil,
		tracingProvider: nil,
		logger:          logging.NewContextLogger(zaptest.NewLogger(t)),
	}

	// Test new HTTP middleware
	middleware := provider.NewHTTPMiddleware()
	assert.NotNil(t, middleware)
}

func TestWithSpan(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Test with span
	err := WithSpan(ctx, "test-span", func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)

	// Test with span that returns an error
	testErr := errors.New("test error")
	err = WithSpan(ctx, "test-span", func(ctx context.Context) error {
		return testErr
	})
	assert.Equal(t, testErr, err)
}

func TestWithSpanTimed(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Test with span timed
	duration, err := WithSpanTimed(ctx, "test-span", func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, duration >= 10*time.Millisecond)

	// Test with span timed that returns an error
	testErr := errors.New("test error")
	duration, err = WithSpanTimed(ctx, "test-span", func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return testErr
	})
	assert.Equal(t, testErr, err)
	assert.True(t, duration >= 10*time.Millisecond)
}
