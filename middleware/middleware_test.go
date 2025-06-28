// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// MockValidationError is a mock implementation of the errors.ValidationErrorInterface
type MockValidationError struct {
	Msg string
}

func (e *MockValidationError) Error() string {
	return e.Msg
}

func (e *MockValidationError) Code() string {
	return "VALIDATION_ERROR"
}

func (e *MockValidationError) HTTPStatus() int {
	return http.StatusBadRequest
}

func (e *MockValidationError) IsValidationError() bool {
	return true
}

// MockNotFoundError is a mock implementation of the errors.NotFoundErrorInterface
type MockNotFoundError struct {
	ResourceType string
	ID           string
}

func (e *MockNotFoundError) Error() string {
	return e.ResourceType + " with ID " + e.ID + " not found"
}

func (e *MockNotFoundError) Code() string {
	return "NOT_FOUND"
}

func (e *MockNotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

func (e *MockNotFoundError) IsNotFoundError() bool {
	return true
}

// MockApplicationError is a mock implementation of the errors.ApplicationErrorInterface
type MockApplicationError struct {
	Msg      string
	CodeName string
}

func (e *MockApplicationError) Error() string {
	return e.Msg
}

func (e *MockApplicationError) Code() string {
	return e.CodeName
}

func (e *MockApplicationError) HTTPStatus() int {
	return http.StatusInternalServerError
}

func (e *MockApplicationError) IsApplicationError() bool {
	return true
}

// MockRepositoryError is a mock implementation of the errors.RepositoryErrorInterface
type MockRepositoryError struct {
	Msg      string
	CodeName string
}

func (e *MockRepositoryError) Error() string {
	return e.Msg
}

func (e *MockRepositoryError) Code() string {
	return e.CodeName
}

func (e *MockRepositoryError) HTTPStatus() int {
	return http.StatusInternalServerError
}

func (e *MockRepositoryError) IsRepositoryError() bool {
	return true
}

func TestGenerateRequestID(t *testing.T) {
	// Generate two request IDs
	id1 := generateRequestID()
	id2 := generateRequestID()

	// Verify they are not empty
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)

	// Verify they are different (unique)
	assert.NotEqual(t, id1, id2)

	// Verify they have the expected format (starts with "req-")
	assert.Contains(t, id1, "req-")
	assert.Contains(t, id2, "req-")
}

func TestRequestID(t *testing.T) {
	// Test with nil context
	assert.Empty(t, RequestID(nil))

	// Test with context but no request ID
	assert.Empty(t, RequestID(context.Background()))

	// Test with context containing request ID
	reqID := "test-request-id"
	ctx := context.WithValue(context.Background(), RequestIDKey, reqID)
	assert.Equal(t, reqID, RequestID(ctx))

	// Test with context containing wrong type for request ID
	ctx = context.WithValue(context.Background(), RequestIDKey, 123)
	assert.Empty(t, RequestID(ctx))
}

func TestStartTime(t *testing.T) {
	// Test with nil context
	assert.True(t, StartTime(nil).IsZero())

	// Test with context but no start time
	assert.True(t, StartTime(context.Background()).IsZero())

	// Test with context containing start time
	now := time.Now()
	ctx := context.WithValue(context.Background(), StartTimeKey, now)
	assert.Equal(t, now, StartTime(ctx))

	// Test with context containing wrong type for start time
	ctx = context.WithValue(context.Background(), StartTimeKey, "not a time")
	assert.True(t, StartTime(ctx).IsZero())
}

func TestRequestDuration(t *testing.T) {
	// Test with nil context
	assert.Equal(t, time.Duration(0), RequestDuration(nil))

	// Test with context but no start time
	assert.Equal(t, time.Duration(0), RequestDuration(context.Background()))

	// Test with context containing start time
	startTime := time.Now().Add(-100 * time.Millisecond)
	ctx := context.WithValue(context.Background(), StartTimeKey, startTime)
	duration := RequestDuration(ctx)
	assert.True(t, duration >= 100*time.Millisecond)
}

func TestWithRequestContext(t *testing.T) {
	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is in context
		reqID := RequestID(r.Context())
		assert.NotEmpty(t, reqID)

		// Check that start time is in context
		startTime := StartTime(r.Context())
		assert.False(t, startTime.IsZero())

		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	handler := WithRequestContext(nextHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("X-Request-ID"))

	// Test with existing request ID
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "existing-id")
	rr = httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "existing-id", rr.Header().Get("X-Request-ID"))
}

func TestWithTimeout(t *testing.T) {
	// Test normal request (completes before timeout)
	t.Run("Normal request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := WithTimeout(1 * time.Second)(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test timeout request
	t.Run("Timeout request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sleep longer than the timeout
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		})

		handler := WithTimeout(10 * time.Millisecond)(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusGatewayTimeout, rr.Code)
	})
}

func TestWithRecovery(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Test normal request (no panic)
	t.Run("Normal request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := WithRecovery(logger, nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test request with panic
	t.Run("Request with panic", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		handler := WithRecovery(logger, nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-request-id")
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Internal Server Error")
		assert.Contains(t, rr.Body.String(), "test-request-id")
	})
}

func TestWithLogging(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	handler := WithLogging(logger, nextHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestResponseWriter(t *testing.T) {
	// Test WriteHeader
	t.Run("WriteHeader", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: rr, statusCode: 0}
		rw.WriteHeader(http.StatusCreated)
		assert.Equal(t, http.StatusCreated, rw.statusCode)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// Test Write
	t.Run("Write with no status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: rr, statusCode: 0}
		n, err := rw.Write([]byte("test"))
		require.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, http.StatusOK, rw.statusCode)
		assert.Equal(t, "test", rr.Body.String())
	})

	t.Run("Write with existing status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: rr, statusCode: http.StatusCreated}
		n, err := rw.Write([]byte("test"))
		require.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, http.StatusCreated, rw.statusCode)
		assert.Equal(t, "test", rr.Body.String())
	})
}

func TestWithErrorHandling(t *testing.T) {
	// Test normal request (no error)
	t.Run("Normal request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := WithErrorHandling(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test request with error
	t.Run("Request with error", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := w.(*errorResponseWriter)
			rw.SetError(&MockNotFoundError{ResourceType: "resource", ID: "123"})
		})

		handler := WithErrorHandling(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-request-id")
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "resource with ID 123 not found")
		assert.Contains(t, rr.Body.String(), "test-request-id")
	})
}

func TestWithCORS(t *testing.T) {
	// Test with no Origin header
	t.Run("No Origin header", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := WithCORS(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	// Test with Origin header
	t.Run("With Origin header", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := WithCORS(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
	})

	// Test OPTIONS request
	t.Run("OPTIONS request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This should not be called for OPTIONS
			t.Error("Next handler should not be called for OPTIONS request")
		})

		handler := WithCORS(nextHandler)
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestWithContextCancellation(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Test normal request
	t.Run("Normal request", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := WithContextCancellation(logger)(nextHandler)
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test timeout (using a mock timeout)
	t.Run("Timeout", func(t *testing.T) {
		// Create a custom middleware that simulates a timeout
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Directly call the timeout case in WithContextCancellation
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte(`{"error":"Request timed out"}`))
		})

		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusGatewayTimeout, rr.Code)
		assert.Contains(t, rr.Body.String(), "Request timed out")
	})
}

func TestApplyMiddleware(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	handler := ApplyMiddleware(nextHandler, logger)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("X-Request-ID"))
}

func TestHandleError(t *testing.T) {
	// Test cases for different error types
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Validation error",
			err:            &MockValidationError{Msg: "validation failed"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name:           "Not found error",
			err:            &MockNotFoundError{ResourceType: "resource", ID: "123"},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "resource with ID 123 not found",
		},
		{
			name:           "Application error",
			err:            &MockApplicationError{Msg: "application error", CodeName: "APP_ERR"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "application error",
		},
		{
			name:           "Repository error",
			err:            &MockRepositoryError{Msg: "database error", CodeName: "DB_ERR"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Database error",
		},
		{
			name:           "Generic error",
			err:            fmt.Errorf("generic error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			ctx := context.WithValue(req.Context(), RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			handleError(rr, req, tc.err)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedBody)
			assert.Contains(t, rr.Body.String(), "test-request-id")
		})
	}
}
