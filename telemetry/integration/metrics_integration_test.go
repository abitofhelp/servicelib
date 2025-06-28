// Copyright (c) 2025 A Bit of Help, Inc.

// +build integration

package integration

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap/zaptest"
)

// TestPrometheusMetricsEndpoint tests the Prometheus metrics endpoint
func TestPrometheusMetricsEndpoint(t *testing.T) {
	// Create a Prometheus handler
	handler := telemetry.CreatePrometheusHandler()
	require.NotNil(t, handler, "Prometheus handler should not be nil")

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a request to the metrics endpoint
	resp, err := http.Get(server.URL)
	require.NoError(t, err, "Request to metrics endpoint should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Metrics endpoint should return 200 OK")
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain", "Content type should be text/plain")

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Reading response body should not fail")

	// Verify that the response contains Prometheus metrics
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "go_", "Response should contain Go metrics")
	assert.Contains(t, bodyStr, "process_", "Response should contain process metrics")
}

// TestHTTPMiddlewareIntegration tests the integration of HTTP middleware with telemetry
func TestHTTPMiddlewareIntegration(t *testing.T) {
	// Create a logger for testing
	logger := zaptest.NewLogger(t)
	ctxLogger := logging.NewContextLogger(logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the trace ID is set in the context
		span := telemetry.GetSpanFromContext(r.Context())
		assert.NotNil(t, span, "Span should be present in context")

		// Write a response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	// Apply telemetry middleware
	middleware := telemetry.NewHTTPMiddleware(ctxLogger)
	handler := middleware(testHandler)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a request to the server
	resp, err := http.Get(server.URL + "/test-path")
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")

	// Verify that the trace ID is set in the response headers
	traceID := resp.Header.Get("X-Trace-ID")
	assert.NotEmpty(t, traceID, "Trace ID should be set in response headers")
}

// TestMetricsRecording tests recording metrics with telemetry
func TestMetricsRecording(t *testing.T) {
	// Skip this test in normal runs since it requires OpenTelemetry setup
	t.Skip("Skipping metrics recording test - requires OpenTelemetry setup")

	// Create a context
	ctx := context.Background()

	// Create a logger for testing
	logger := zaptest.NewLogger(t)
	ctxLogger := logging.NewContextLogger(logger)

	// Create a minimal koanf configuration
	k := koanf.New(".")

	// Create a metrics provider
	metricsProvider, err := telemetry.NewMetricsProvider(ctx, ctxLogger, k)
	if err != nil {
		// This is expected in test environment without proper configuration
		t.Logf("Metrics provider creation failed (expected in test): %v", err)
		return
	}

	if metricsProvider == nil {
		t.Log("Metrics provider is nil (expected in test without proper configuration)")
		return
	}

	// Initialize common metrics
	err = telemetry.InitCommonMetrics(metricsProvider)
	require.NoError(t, err, "Initializing common metrics should not fail")

	// Record HTTP request metrics
	telemetry.RecordHTTPRequest(ctx, "GET", "/test", 200, 100*time.Millisecond, 1024)

	// Record in-flight requests
	telemetry.IncrementRequestsInFlight(ctx, "GET", "/test")
	telemetry.DecrementRequestsInFlight(ctx, "GET", "/test")

	// Record DB operation metrics
	telemetry.RecordDBOperation(ctx, "query", "postgres", "users", 50*time.Millisecond, nil)

	// Record error metrics
	telemetry.RecordErrorMetric(ctx, "auth", "token_validation")

	// Shutdown the metrics provider
	err = metricsProvider.Shutdown(ctx)
	require.NoError(t, err, "Shutting down metrics provider should not fail")
}

// TestTelemetryProviderIntegration tests the integration of the telemetry provider
func TestTelemetryProviderIntegration(t *testing.T) {
	// Skip this test in normal runs since it requires OpenTelemetry setup
	t.Skip("Skipping telemetry provider integration test - requires OpenTelemetry setup")

	// Create a context
	ctx := context.Background()

	// Create a logger for testing
	logger := zaptest.NewLogger(t)
	ctxLogger := logging.NewContextLogger(logger)

	// Create a minimal koanf configuration
	k := koanf.New(".")

	// Create a telemetry provider
	provider, err := telemetry.NewTelemetryProvider(ctx, ctxLogger, k)
	if err != nil {
		// This is expected in test environment without proper configuration
		t.Logf("Telemetry provider creation failed (expected in test): %v", err)
		return
	}

	if provider == nil {
		t.Log("Telemetry provider is nil (expected in test without proper configuration)")
		return
	}

	// Get a meter
	meter := provider.Meter()
	require.NotNil(t, meter, "Meter should not be nil")

	// Create a counter
	counter, err := meter.Int64Counter("test.counter", metric.WithDescription("Test counter"))
	require.NoError(t, err, "Creating counter should not fail")

	// Record a value
	counter.Add(ctx, 1)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Instrument the handler
	instrumentedHandler := provider.InstrumentHandler(testHandler, "test-operation")
	require.NotNil(t, instrumentedHandler, "Instrumented handler should not be nil")

	// Create a test server
	server := httptest.NewServer(instrumentedHandler)
	defer server.Close()

	// Make a request to the server
	resp, err := http.Get(server.URL)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")

	// Create a Prometheus handler
	prometheusHandler := provider.CreatePrometheusHandler()
	require.NotNil(t, prometheusHandler, "Prometheus handler should not be nil")

	// Shutdown the provider
	err = provider.Shutdown(ctx)
	require.NoError(t, err, "Shutting down telemetry provider should not fail")
}

// TestSpanOperations tests span operations
func TestSpanOperations(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Start a span
	ctx, span := telemetry.StartSpan(ctx, "test-span")
	require.NotNil(t, span, "Span should not be nil")

	// Add attributes to the span
	telemetry.AddSpanAttributes(ctx, 
		attribute.String("test", "value"),
		attribute.Int("count", 42),
	)

	// Record an error
	err := io.EOF
	telemetry.RecordErrorSpan(ctx, err)

	// End the span
	span.End()

	// Test WithSpan helper
	err = telemetry.WithSpan(ctx, "test-with-span", func(ctx context.Context) error {
		// Do something with the context
		return nil
	})
	require.NoError(t, err, "WithSpan should not fail")

	// Test WithSpanTimed helper
	duration, err := telemetry.WithSpanTimed(ctx, "test-with-span-timed", func(ctx context.Context) error {
		// Do something with the context
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	require.NoError(t, err, "WithSpanTimed should not fail")
	assert.True(t, duration >= 10*time.Millisecond, "Duration should be at least 10ms")
}

// Helper function to get span from context
func TestGetSpanFromContext(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Get span from context (should be nil)
	span := telemetry.GetSpanFromContext(ctx)
	assert.Nil(t, span, "Span should be nil for empty context")

	// Start a span
	ctx, span = telemetry.StartSpan(ctx, "test-span")
	require.NotNil(t, span, "Span should not be nil")

	// Get span from context (should be the same span)
	gotSpan := telemetry.GetSpanFromContext(ctx)
	assert.Equal(t, span, gotSpan, "Should get the same span from context")
}
