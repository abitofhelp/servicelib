// Copyright (c) 2025 A Bit of Help, Inc.

// +build integration

package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

// TestInstrumentedHTTPClient tests the instrumented HTTP client
func TestInstrumentedHTTPClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has trace headers
		traceID := r.Header.Get("Traceparent")
		assert.NotEmpty(t, traceID, "Request should have trace headers")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	// Create an instrumented client
	client := telemetry.InstrumentClient(&http.Client{})
	require.NotNil(t, client, "Instrumented client should not be nil")

	// Create a context with a span
	ctx, span := telemetry.StartSpan(context.Background(), "test-client-request")
	defer span.End()

	// Make a request with the instrumented client
	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	require.NoError(t, err, "Creating request should not fail")

	resp, err := client.Do(req)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

// TestInstrumentedHTTPHandler tests the instrumented HTTP handler
func TestInstrumentedHTTPHandler(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request has a span in its context
		span := telemetry.GetSpanFromContext(r.Context())
		assert.NotNil(t, span, "Request context should have a span")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	// Instrument the handler
	instrumentedHandler := telemetry.InstrumentHandler(testHandler, "test-handler")
	require.NotNil(t, instrumentedHandler, "Instrumented handler should not be nil")

	// Create a test server with the instrumented handler
	server := httptest.NewServer(instrumentedHandler)
	defer server.Close()

	// Make a request to the server
	resp, err := http.Get(server.URL)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

// TestHTTPMiddlewareWithAttributes tests the HTTP middleware with custom attributes
func TestHTTPMiddlewareWithAttributes(t *testing.T) {
	// Create a test handler that checks for attributes in the span
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the span from the context
		span := telemetry.GetSpanFromContext(r.Context())
		require.NotNil(t, span, "Span should be present in context")

		// Add custom attributes to the span
		telemetry.AddSpanAttributes(r.Context(),
			attribute.String("custom.attribute", "test-value"),
			attribute.Int("custom.count", 42),
		)

		// Write a response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	})

	// Create a custom middleware that adds request-specific attributes
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start a new span for this request
			ctx, span := telemetry.StartSpan(r.Context(), "custom-middleware")
			defer span.End()

			// Add request-specific attributes
			telemetry.AddSpanAttributes(ctx,
				attribute.String("request.path", r.URL.Path),
				attribute.String("request.method", r.Method),
				attribute.String("request.user_agent", r.UserAgent()),
			)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	// Apply the middleware
	handler := middleware(testHandler)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a request to the server
	req, err := http.NewRequest("GET", server.URL+"/test-path", nil)
	require.NoError(t, err, "Creating request should not fail")
	req.Header.Set("User-Agent", "test-agent")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type should be application/json")
}

// TestInstrumentedClientWithOptions tests the instrumented client with custom options
func TestInstrumentedClientWithOptions(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create an instrumented client with custom options
	client := telemetry.InstrumentClient(&http.Client{}, 
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return "custom-" + r.Method + "-" + r.URL.Path
		}),
	)
	require.NotNil(t, client, "Instrumented client should not be nil")

	// Make a request with the instrumented client
	resp, err := client.Get(server.URL + "/custom-path")
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

// TestInstrumentedHandlerWithOptions tests the instrumented handler with custom options
func TestInstrumentedHandlerWithOptions(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Instrument the handler with custom options
	instrumentedHandler := telemetry.InstrumentHandler(testHandler, "test-handler",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return "custom-" + operation + "-" + r.Method
		}),
	)
	require.NotNil(t, instrumentedHandler, "Instrumented handler should not be nil")

	// Create a test server with the instrumented handler
	server := httptest.NewServer(instrumentedHandler)
	defer server.Close()

	// Make a request to the server
	resp, err := http.Get(server.URL)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}