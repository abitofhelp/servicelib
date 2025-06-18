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
	"go.opentelemetry.io/otel"
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

	// Save the original metrics
	origHTTPRequestsTotal := httpRequestsTotal
	origHTTPRequestDuration := httpRequestDuration
	origHTTPResponseSizeBytes := httpResponseSizeBytes

	defer func() {
		// Restore the original metrics
		httpRequestsTotal = origHTTPRequestsTotal
		httpRequestDuration = origHTTPRequestDuration
		httpResponseSizeBytes = origHTTPResponseSizeBytes
	}()

	// Test with nil metrics
	httpRequestsTotal = nil
	httpRequestDuration = nil
	httpResponseSizeBytes = nil

	// This should not panic
	RecordHTTPRequest(ctx, "GET", "/test", 200, 100*time.Millisecond, 1024)

	// Test with non-nil metrics
	// Create dummy metrics (we can't verify they're called, but we can ensure no panics)
	httpRequestsTotal = origHTTPRequestsTotal
	httpRequestDuration = origHTTPRequestDuration
	httpResponseSizeBytes = origHTTPResponseSizeBytes

	// This should not panic
	RecordHTTPRequest(ctx, "GET", "/test", 200, 100*time.Millisecond, 1024)
}

func TestIncrementRequestsInFlight(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Save the original metric
	origHTTPRequestsInFlight := httpRequestsInFlight

	defer func() {
		// Restore the original metric
		httpRequestsInFlight = origHTTPRequestsInFlight
	}()

	// Test with nil metric
	httpRequestsInFlight = nil

	// This should not panic
	IncrementRequestsInFlight(ctx, "GET", "/test")

	// Test with non-nil metric
	// Create dummy metric (we can't verify it's called, but we can ensure no panics)
	httpRequestsInFlight = origHTTPRequestsInFlight

	// This should not panic
	IncrementRequestsInFlight(ctx, "GET", "/test")
}

func TestDecrementRequestsInFlight(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Save the original metric
	origHTTPRequestsInFlight := httpRequestsInFlight

	defer func() {
		// Restore the original metric
		httpRequestsInFlight = origHTTPRequestsInFlight
	}()

	// Test with nil metric
	httpRequestsInFlight = nil

	// This should not panic
	DecrementRequestsInFlight(ctx, "GET", "/test")

	// Test with non-nil metric
	// Create dummy metric (we can't verify it's called, but we can ensure no panics)
	httpRequestsInFlight = origHTTPRequestsInFlight

	// This should not panic
	DecrementRequestsInFlight(ctx, "GET", "/test")
}

func TestRecordDBOperation(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Save the original metrics
	origDBOperationsTotal := dbOperationsTotal
	origDBOperationDuration := dbOperationDuration
	origAppErrorsTotal := appErrorsTotal

	defer func() {
		// Restore the original metrics
		dbOperationsTotal = origDBOperationsTotal
		dbOperationDuration = origDBOperationDuration
		appErrorsTotal = origAppErrorsTotal
	}()

	// Test with nil metrics
	dbOperationsTotal = nil
	dbOperationDuration = nil
	appErrorsTotal = nil

	// This should not panic
	RecordDBOperation(ctx, "query", "postgres", "users", 100*time.Millisecond, nil)

	// Test with error and nil metrics
	err := errors.New("database error")
	RecordDBOperation(ctx, "query", "postgres", "users", 100*time.Millisecond, err)

	// Test with non-nil metrics
	// Create dummy metrics (we can't verify they're called, but we can ensure no panics)
	dbOperationsTotal = origDBOperationsTotal
	dbOperationDuration = origDBOperationDuration
	appErrorsTotal = origAppErrorsTotal

	// This should not panic
	RecordDBOperation(ctx, "query", "postgres", "users", 100*time.Millisecond, nil)

	// Test with error and non-nil metrics
	RecordDBOperation(ctx, "query", "postgres", "users", 100*time.Millisecond, err)
}

func TestUpdateDBConnections(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Save the original metric
	origDBConnectionsOpen := dbConnectionsOpen

	defer func() {
		// Restore the original metric
		dbConnectionsOpen = origDBConnectionsOpen
	}()

	// Test with nil metric
	dbConnectionsOpen = nil

	// This should not panic
	UpdateDBConnections(ctx, "postgres", 1)
	UpdateDBConnections(ctx, "postgres", -1)

	// Test with non-nil metric
	// Create dummy metric (we can't verify it's called, but we can ensure no panics)
	dbConnectionsOpen = origDBConnectionsOpen

	// This should not panic
	UpdateDBConnections(ctx, "postgres", 1)
	UpdateDBConnections(ctx, "postgres", -1)
}

func TestRecordErrorMetric(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Save the original metric
	origAppErrorsTotal := appErrorsTotal

	defer func() {
		// Restore the original metric
		appErrorsTotal = origAppErrorsTotal
	}()

	// Test with nil metric
	appErrorsTotal = nil

	// This should not panic
	RecordErrorMetric(ctx, "validation", "create_user")

	// Test with non-nil metric
	// Create dummy metric (we can't verify it's called, but we can ensure no panics)
	appErrorsTotal = origAppErrorsTotal

	// This should not panic
	RecordErrorMetric(ctx, "validation", "create_user")
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

	// Test with invalid resource creation
	k = koanf.New(".")
	k.Set("telemetry.metrics.enabled", true)
	// Missing required service name and other attributes

	// This should fail because we're missing required attributes for resource creation
	provider, err = NewMetricsProvider(context.Background(), logger, k)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "failed to create resource")
	}
	assert.Nil(t, provider)

	// Test with invalid OTLP endpoint
	k = koanf.New(".")
	k.Set("telemetry.metrics.enabled", true)
	k.Set("telemetry.service_name", "test-service")
	k.Set("telemetry.environment", "test")
	k.Set("telemetry.version", "1.0.0")
	k.Set("telemetry.otlp.endpoint", "invalid-endpoint") // Invalid endpoint
	k.Set("telemetry.otlp.insecure", true)
	k.Set("telemetry.metrics.reporting_frequency_seconds", 15)

	// This should fail because the OTLP endpoint is invalid
	provider, err = NewMetricsProvider(context.Background(), logger, k)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "failed to create OTLP exporter")
	}
	assert.Nil(t, provider)

	// Note: We can't easily test the successful case as it would try to connect to a real OTLP endpoint
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
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a metrics provider with a nil meter
	provider := &MetricsProvider{
		provider: nil,
		meter:    nil,
		logger:   logger,
	}

	// Test Meter() with nil meter
	meter := provider.Meter()
	assert.Nil(t, meter, "Meter should be nil when the provider's meter is nil")

	// Create a real meter
	realMeter := otel.Meter("test-meter")

	// Create a metrics provider with a non-nil meter
	provider = &MetricsProvider{
		provider: nil,
		meter:    realMeter,
		logger:   logger,
	}

	// Test Meter() with non-nil meter
	meter = provider.Meter()
	assert.Equal(t, realMeter, meter, "Meter should return the provider's meter")
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
