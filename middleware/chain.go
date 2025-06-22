// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
)

// Middleware represents a middleware function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain applies multiple middleware to a http.Handler in the order they are provided
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order so that the first middleware in the list
	// is the outermost wrapper (first to process the request, last to process the response)
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		handler = middleware(handler)
	}
	return handler
}

// WithRequestID is a middleware that adds a request ID to the context
func WithRequestID(ctx context.Context) Middleware {
	return func(next http.Handler) http.Handler {
		return WithRequestContext(next)
	}
}

// Logging is a middleware that logs request information
func Logging(logger *logging.ContextLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return WithLogging(logger, next)
	}
}

// Recovery is a middleware that recovers from panics
func Recovery(logger *logging.ContextLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return WithRecovery(logger, next)
	}
}
