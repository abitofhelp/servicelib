// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// ContextKey is a type for context keys to avoid collisions.
// Using a custom type for context keys helps prevent key collisions
// when multiple packages store values in the context.
type ContextKey string

const (
	// RequestIDKey is the key for the request ID in the context.
	// This key is used to store and retrieve the unique identifier for each request.
	RequestIDKey ContextKey = "request_id"

	// StartTimeKey is the key for the request start time in the context.
	// This key is used to store and retrieve the time when the request started processing.
	StartTimeKey ContextKey = "start_time"

	// UserIDKey is the key for the user ID in the context (for future auth).
	// This key is used to store and retrieve the ID of the authenticated user.
	UserIDKey ContextKey = "user_id"
)

// generateRequestID generates a simple unique ID for a request
func generateRequestID() string {
	// Initialize random number generator with current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random number
	randomNum := rand.Intn(1000000)

	// Format the request ID using timestamp and random number
	return fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), randomNum)
}

// RequestID returns the request ID from the context.
// This function safely extracts the request ID from the context, handling nil contexts
// and type assertions. If the context is nil or doesn't contain a request ID,
// an empty string is returned.
//
// Parameters:
//   - ctx: The context from which to extract the request ID.
//
// Returns:
//   - A string containing the request ID, or an empty string if not found.
func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// StartTime returns the request start time from the context.
// This function safely extracts the start time from the context, handling nil contexts
// and type assertions. If the context is nil or doesn't contain a start time,
// a zero time is returned.
//
// Parameters:
//   - ctx: The context from which to extract the start time.
//
// Returns:
//   - A time.Time representing when the request started, or a zero time if not found.
func StartTime(ctx context.Context) time.Time {
	if ctx == nil {
		return time.Time{}
	}
	if startTime, ok := ctx.Value(StartTimeKey).(time.Time); ok {
		return startTime
	}
	return time.Time{}
}

// RequestDuration returns the duration since the request started.
// This function calculates how long a request has been processing by comparing
// the start time from the context with the current time. If the start time
// is not found or is a zero time, a duration of 0 is returned.
//
// Parameters:
//   - ctx: The context containing the request start time.
//
// Returns:
//   - A time.Duration representing how long the request has been processing.
func RequestDuration(ctx context.Context) time.Duration {
	startTime := StartTime(ctx)
	if startTime.IsZero() {
		return 0
	}
	return time.Since(startTime)
}

// WithRequestContext adds request context information to the request.
// This middleware adds a unique request ID and start time to the request context.
// If the request already has an X-Request-ID header, that value is used as the
// request ID; otherwise, a new unique ID is generated. The request ID is also
// added to the response headers for correlation.
//
// This middleware is typically one of the first in the chain, as it provides
// context information that other middleware and handlers can use.
//
// Parameters:
//   - next: The next handler in the middleware chain.
//
// Returns:
//   - An http.Handler that wraps the next handler with request context functionality.
func WithRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a request ID if not already present
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Create a new context with request information
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		ctx = context.WithValue(ctx, StartTimeKey, time.Now())

		// Call the next handler with the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithTimeout adds a timeout to the request context.
// This middleware sets a maximum duration for the request to complete. If the request
// takes longer than the specified timeout, it is canceled and a 504 Gateway Timeout
// response is returned to the client.
//
// The middleware uses a goroutine to process the request and a select statement to
// wait for either the request to complete or the timeout to expire. It uses a
// synchronized response writer to avoid race conditions when writing the response.
//
// Parameters:
//   - timeout: The maximum duration allowed for the request to complete.
//
// Returns:
//   - A middleware function that adds timeout functionality to an HTTP handler.
func WithTimeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Create a channel to detect when the request is done
			done := make(chan struct{})

			// Create a mutex to protect access to the response writer
			var mu sync.Mutex

			// Create a flag to track whether a response has been written
			var responded bool

			// Create a wrapper for the response writer that's protected by the mutex
			syncWriter := &syncResponseWriter{
				ResponseWriter: w,
				mu:             &mu,
				responded:      &responded,
			}

			// Process the request in a goroutine
			go func() {
				next.ServeHTTP(syncWriter, r.WithContext(ctx))
				close(done)
			}()

			// Wait for either the request to complete or the timeout to expire
			select {
			case <-done:
				// Request completed normally
				return
			case <-ctx.Done():
				// Context timed out
				if ctx.Err() == context.DeadlineExceeded {
					// Use the synchronized writer to avoid race conditions
					syncWriter.WriteError("Request timeout", http.StatusGatewayTimeout)
				}
				return
			}
		})
	}
}

// WithRecovery adds panic recovery to the request.
// This middleware catches any panics that occur during request processing and
// converts them into 500 Internal Server Error responses. It also logs the panic
// details, including the stack trace, to help with debugging.
//
// Without this middleware, a panic in a handler would crash the entire server.
// With it, the panic is contained to just the current request, allowing the
// server to continue handling other requests.
//
// Parameters:
//   - logger: A context logger for logging panic details.
//   - next: The next handler in the middleware chain.
//
// Returns:
//   - An http.Handler that wraps the next handler with panic recovery functionality.
func WithRecovery(logger *logging.ContextLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				ctx := r.Context()
				requestID := RequestID(ctx)

				// Log the error with the request ID and stack trace
				logger.Error(ctx, "Panic recovered",
					zap.String("request_id", requestID),
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("url", r.URL.String()),
					zap.String("method", r.Method),
				)

				// Return a 500 Internal Server Error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, writeErr := w.Write([]byte(`{"error":"Internal Server Error","request_id":"` + requestID + `"}`))
				if writeErr != nil {
					logger.Error(ctx, "Failed to write error response", zap.Error(writeErr))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// WithLogging adds request logging.
// This middleware logs information about each request, including the HTTP method,
// path, status code, and duration. It uses a response writer wrapper to capture
// the status code of the response.
//
// The logging occurs after the request has been processed, so it includes the
// final status code and the total time taken to process the request. This is
// useful for monitoring and debugging.
//
// Parameters:
//   - logger: A context logger for logging request details.
//   - next: The next handler in the middleware chain.
//
// Returns:
//   - An http.Handler that wraps the next handler with logging functionality.
func WithLogging(logger *logging.ContextLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record start time
		start := time.Now()

		// Create a response wrapper to capture the status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the request
		ctx := r.Context()
		requestID := RequestID(ctx)

		// Use structured logger
		logger.Info(ctx, "Request completed",
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", duration))
	})
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the status code if not already set
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// syncResponseWriter is a thread-safe wrapper for http.ResponseWriter
type syncResponseWriter struct {
	http.ResponseWriter
	mu        *sync.Mutex
	responded *bool
}

// WriteHeader is a thread-safe implementation of http.ResponseWriter.WriteHeader
func (srw *syncResponseWriter) WriteHeader(code int) {
	srw.mu.Lock()
	defer srw.mu.Unlock()

	if !*srw.responded {
		*srw.responded = true
		srw.ResponseWriter.WriteHeader(code)
	}
}

// Write is a thread-safe implementation of http.ResponseWriter.Write
func (srw *syncResponseWriter) Write(b []byte) (int, error) {
	srw.mu.Lock()
	defer srw.mu.Unlock()

	if !*srw.responded {
		*srw.responded = true
		return srw.ResponseWriter.Write(b)
	}
	return len(b), nil // Pretend we wrote it
}

// WriteError is a convenience method for writing an error response
func (srw *syncResponseWriter) WriteError(msg string, statusCode int) {
	srw.mu.Lock()
	defer srw.mu.Unlock()

	if !*srw.responded {
		*srw.responded = true
		srw.ResponseWriter.Header().Set("Content-Type", "text/plain")
		srw.ResponseWriter.WriteHeader(statusCode)
		srw.ResponseWriter.Write([]byte(msg))
	}
}

// WithErrorHandling adds centralized error handling
func WithErrorHandling(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response wrapper to capture errors
		rw := &errorResponseWriter{ResponseWriter: w}

		// Process the request
		next.ServeHTTP(rw, r)

		// Handle any captured error
		if rw.err != nil {
			handleError(w, r, rw.err)
		}
	})
}

// errorResponseWriter is a wrapper for http.ResponseWriter that captures errors
type errorResponseWriter struct {
	http.ResponseWriter
	err error
}

// SetError sets the error to be handled
func (rw *errorResponseWriter) SetError(err error) {
	rw.err = err
}

// WithCORS adds CORS headers to allow cross-origin requests
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Origin header from the request
		origin := r.Header.Get("Origin")
		if origin == "" {
			// If no Origin header is present, use a wildcard
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// If Origin header is present, echo it back
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Set other CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID, X-Apollo-Operation-Name, Apollo-Require-Preflight, GraphQL-Query, GraphQL-Variables, GraphQL-Operation-Name, Origin, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// WithContextCancellation is a middleware that checks for context cancellation
// It detects when a client disconnects and logs the event
func WithContextCancellation(logger *logging.ContextLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a done channel to signal when the handler is complete
			done := make(chan struct{})

			// Create a copy of the request with a context that we can check
			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()
			r = r.WithContext(ctx)

			// Create a timeout to prevent goroutine leaks
			timeout := time.After(60 * time.Second)

			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
				// Request completed normally
				return
			case <-r.Context().Done():
				// Client disconnected or request was cancelled
				logger.Info(ctx, "Request cancelled by client",
					zap.String("url", r.URL.String()),
					zap.String("method", r.Method),
					zap.Error(r.Context().Err()),
				)
				// We don't need to send a response as the client has disconnected
				cancel()
				return
			case <-timeout:
				// Request took too long, log and cancel
				logger.Warn(ctx, "Request timed out after 60 seconds",
					zap.String("url", r.URL.String()),
					zap.String("method", r.Method),
				)
				// Cancel the context to signal the handler to stop
				cancel()
				// Return a 504 Gateway Timeout response
				w.WriteHeader(http.StatusGatewayTimeout)
				_, writeErr := w.Write([]byte(`{"error":"Request timed out"}`))
				if writeErr != nil {
					logger.Error(ctx, "Failed to write timeout response", zap.Error(writeErr))
				}
				return
			}
		})
	}
}

// ApplyMiddleware applies all middleware to a handler in the recommended order.
// This function provides a consistent way to apply middleware across different applications.
// It applies the following middleware in order (outermost first):
//  1. WithRequestContext - Adds request ID and timing information
//  2. WithCORS - Adds CORS headers for cross-origin requests
//  3. WithRecovery - Catches panics and converts them to 500 responses
//  4. WithContextCancellation - Handles client disconnections
//  5. WithLogging - Logs request details
//  6. WithErrorHandling - Maps errors to appropriate HTTP responses
//
// The order is important for proper middleware chaining. For example, the request context
// middleware is applied first so that all other middleware can access the request ID.
// The recovery middleware is applied early to catch panics in other middleware.
//
// Parameters:
//   - handler: The HTTP handler to wrap with middleware.
//   - logger: Logger for recording middleware events.
//
// Returns:
//   - A wrapped HTTP handler with all middleware applied in the recommended order.
func ApplyMiddleware(handler http.Handler, logger *zap.Logger) http.Handler {
	// Create a context logger from the zap logger
	contextLogger := logging.NewContextLogger(logger)

	// Apply middleware in order (outermost first)
	// The order is important for proper middleware chaining

	// First add request context information
	wrappedHandler := WithRequestContext(handler)

	// Add CORS headers
	wrappedHandler = WithCORS(wrappedHandler)

	// Add recovery to catch panics
	wrappedHandler = WithRecovery(contextLogger, wrappedHandler)

	// Add context cancellation handling
	wrappedHandler = WithContextCancellation(contextLogger)(wrappedHandler)

	// Add request logging
	wrappedHandler = WithLogging(contextLogger, wrappedHandler)

	// Add error handling
	wrappedHandler = WithErrorHandling(wrappedHandler)

	return wrappedHandler
}

// handleError maps errors to appropriate HTTP responses
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Map error types to status codes using the generic error interfaces
	statusCode := errors.GetHTTPStatus(err)

	// If statusCode is 0 (invalid), use a default status code
	if statusCode == 0 {
		// Try to determine a more specific status code based on error type
		// First check if the error implements IsNotFoundError() bool
		if e, ok := err.(interface{ IsNotFoundError() bool }); ok && e.IsNotFoundError() {
			statusCode = http.StatusNotFound
		} else if errors.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		} else if e, ok := err.(interface{ IsValidationError() bool }); ok && e.IsValidationError() {
			statusCode = http.StatusBadRequest
		} else if errors.IsValidationError(err) {
			statusCode = http.StatusBadRequest
		} else if e, ok := err.(interface{ IsDatabaseError() bool }); ok && e.IsDatabaseError() {
			statusCode = http.StatusInternalServerError
		} else if e, ok := err.(interface{ IsRepositoryError() bool }); ok && e.IsRepositoryError() {
			statusCode = http.StatusInternalServerError
		} else if errors.IsDatabaseError(err) {
			statusCode = http.StatusInternalServerError
		} else {
			// Default to 500 Internal Server Error
			statusCode = http.StatusInternalServerError
		}
	}

	errorMessage := "Internal server error"

	// Get appropriate error message based on error type
	// First check if the error implements the specific interface methods
	if e, ok := err.(interface{ IsValidationError() bool }); ok && e.IsValidationError() {
		errorMessage = err.Error()
	} else if e, ok := err.(interface{ IsNotFoundError() bool }); ok && e.IsNotFoundError() {
		errorMessage = err.Error()
	} else if e, ok := err.(interface{ IsApplicationError() bool }); ok && e.IsApplicationError() {
		errorMessage = err.Error()
	} else if errors.IsValidationError(err) || errors.IsNotFoundError(err) || errors.IsApplicationError(err) {
		errorMessage = err.Error()
	} else if e, ok := err.(interface{ IsDatabaseError() bool }); ok && e.IsDatabaseError() {
		errorMessage = "Database error"
	} else if e, ok := err.(interface{ IsRepositoryError() bool }); ok && e.IsRepositoryError() {
		errorMessage = "Database error"
	} else if errors.IsDatabaseError(err) {
		errorMessage = "Database error"
	}

	// Set status code
	w.WriteHeader(statusCode)

	// Write error response
	w.Write([]byte(`{"error":"` + errorMessage + `","request_id":"` + RequestID(r.Context()) + `"}`))
}
