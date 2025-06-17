package telemetry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zaptest"
)

func TestInstrumentHandler(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Instrument the handler
	instrumentedHandler := InstrumentHandler(handler, "test-operation")

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	instrumentedHandler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestInstrumentClient(t *testing.T) {
	// Create a test client
	client := &http.Client{}

	// Instrument the client
	instrumentedClient := InstrumentClient(client)

	// Verify the client was instrumented
	assert.NotNil(t, instrumentedClient)
	assert.NotNil(t, instrumentedClient.Transport)

	// Test with nil client
	nilClient := InstrumentClient(nil)
	assert.NotNil(t, nilClient)
	assert.NotNil(t, nilClient.Transport)
}

func TestNewHTTPMiddleware(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create the middleware
	middleware := NewHTTPMiddleware(logger)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the context has a span
		span := trace.SpanFromContext(r.Context())
		assert.NotNil(t, span)

		// Verify that the trace ID is set in the response headers
		traceID := span.SpanContext().TraceID().String()
		assert.NotEmpty(t, traceID)
		assert.Equal(t, traceID, w.Header().Get("X-Trace-ID"))

		w.WriteHeader(http.StatusOK)
	})

	// Apply the middleware
	wrappedHandler := middleware(handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	wrappedHandler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("X-Trace-ID"))
}

func TestStartSpan(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Start a span
	spanCtx, span := StartSpan(ctx, "test-span")

	// Verify the span was created
	assert.NotNil(t, span)
	assert.NotEqual(t, ctx, spanCtx)

	// End the span
	span.End()
}

func TestAddSpanAttributes(t *testing.T) {
	// Create a context with a span
	ctx, span := StartSpan(context.Background(), "test-span")
	defer span.End()

	// Add attributes to the span
	AddSpanAttributes(ctx, attribute.String("key", "value"))

	// No assertion needed as we can't easily verify the attributes were added
}

func TestRecordErrorSpan(t *testing.T) {
	// Create a context with a span
	ctx, span := StartSpan(context.Background(), "test-span")
	defer span.End()

	// Record an error on the span
	err := errors.New("test error")
	RecordErrorSpan(ctx, err)

	// Test with nil error
	RecordErrorSpan(ctx, nil)

	// No assertion needed as we can't easily verify the error was recorded
}
