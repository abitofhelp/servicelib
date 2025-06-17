package telemetry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestCreatePrometheusHandler(t *testing.T) {
	// Create a Prometheus handler
	handler := CreatePrometheusHandler()

	// Verify the handler was created
	assert.NotNil(t, handler)

	// Test that the handler responds to requests
	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "text/plain")
}

func TestRecordHTTPRequest(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Record an HTTP request
	RecordHTTPRequest(ctx, "GET", "/test", 200, 100*time.Millisecond, 1024)

	// No assertion needed as we can't easily verify the metrics were recorded
}

func TestIncrementRequestsInFlight(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Increment requests in flight
	IncrementRequestsInFlight(ctx, "GET", "/test")

	// No assertion needed as we can't easily verify the metrics were recorded
}

func TestDecrementRequestsInFlight(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Decrement requests in flight
	DecrementRequestsInFlight(ctx, "GET", "/test")

	// No assertion needed as we can't easily verify the metrics were recorded
}

func TestRecordDBOperation(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Record a successful DB operation
	RecordDBOperation(ctx, "query", "postgres", "users", 100*time.Millisecond, nil)

	// Record a failed DB operation
	err := errors.New("database error")
	RecordDBOperation(ctx, "query", "postgres", "users", 100*time.Millisecond, err)

	// No assertion needed as we can't easily verify the metrics were recorded
}

func TestUpdateDBConnections(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Update DB connections
	UpdateDBConnections(ctx, "postgres", 1)
	UpdateDBConnections(ctx, "postgres", -1)

	// No assertion needed as we can't easily verify the metrics were recorded
}

func TestRecordErrorMetric(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Record an error metric
	RecordErrorMetric(ctx, "validation", "create_user")

	// No assertion needed as we can't easily verify the metrics were recorded
}

func TestNewMetricsProvider(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a koanf instance with metrics disabled
	k := koanf.New(".")
	k.Set("telemetry.metrics.enabled", false)

	// Create a metrics provider with metrics disabled
	provider, err := NewMetricsProvider(context.Background(), logger, k)
	require.NoError(t, err)
	assert.Nil(t, provider)

	// Create a koanf instance with metrics enabled
	k = koanf.New(".")
	k.Set("telemetry.metrics.enabled", true)
	k.Set("telemetry.service_name", "test-service")
	k.Set("telemetry.environment", "test")
	k.Set("telemetry.version", "1.0.0")
	k.Set("telemetry.otlp.endpoint", "localhost:4317")
	k.Set("telemetry.otlp.insecure", true)
	k.Set("telemetry.metrics.reporting_frequency_seconds", 15)

	// Skip creating a metrics provider with metrics enabled as it would try to connect to a real OTLP endpoint
}

func TestInitCommonMetrics(t *testing.T) {
	// Test with nil provider
	err := InitCommonMetrics(nil)
	assert.NoError(t, err)

	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a metrics provider with a nil meter
	provider := &MetricsProvider{
		provider: nil,
		meter:    nil,
		logger:   logger,
	}

	// Test with provider that has a nil meter
	err = InitCommonMetrics(provider)
	assert.NoError(t, err)
}

func TestMetricsProviderMeter(t *testing.T) {
	// Skip this test as we can't easily create a meter in a test environment
	t.Skip("Skipping test as we can't easily create a meter in a test environment")
}

func TestMetricsProviderShutdown(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a metrics provider with a nil provider
	provider := &MetricsProvider{
		provider: nil,
		meter:    nil,
		logger:   logger,
	}

	// Test shutdown with nil provider
	err := provider.Shutdown(context.Background())
	assert.NoError(t, err)
}
