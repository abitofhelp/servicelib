// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
)

// Middleware represents a middleware function that wraps an http.Handler.
// A middleware takes an http.Handler and returns a new http.Handler that adds
// some functionality before and/or after calling the original handler.
//
// This type allows middleware to be composed and chained together in a flexible way.
type Middleware func(http.Handler) http.Handler

// Chain applies multiple middleware to a http.Handler in the order they are provided.
// The first middleware in the list becomes the outermost wrapper, meaning it is the
// first to process the incoming request and the last to process the outgoing response.
//
// The middleware are applied in reverse order (from last to first) to achieve this,
// because each middleware wraps the handler that comes after it in the chain.
//
// Parameters:
//   - handler: The base HTTP handler to wrap with middleware.
//   - middlewares: A variadic list of middleware functions to apply.
//
// Returns:
//   - An http.Handler with all the middleware applied in the specified order.
//
// Example:
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("Hello, World!"))
//	})
//	
//	// Apply middleware in order: logging, then recovery, then request ID
//	wrappedHandler := middleware.Chain(handler,
//	    middleware.WithRequestID(context.Background()),
//	    middleware.Recovery(logger),
//	    middleware.Logging(logger),
//	)
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order so that the first middleware in the list
	// is the outermost wrapper (first to process the request, last to process the response)
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		handler = middleware(handler)
	}
	return handler
}

// WithRequestID is a middleware that adds a request ID to the context.
// This is a convenience wrapper around WithRequestContext that returns a Middleware
// function, making it easier to use with the Chain function.
//
// Parameters:
//   - ctx: The parent context to use (typically context.Background()).
//
// Returns:
//   - A Middleware function that adds request ID functionality.
func WithRequestID(ctx context.Context) Middleware {
	return func(next http.Handler) http.Handler {
		return WithRequestContext(next)
	}
}

// Logging is a middleware that logs request information.
// This is a convenience wrapper around WithLogging that returns a Middleware
// function, making it easier to use with the Chain function.
//
// Parameters:
//   - logger: A context logger for logging request details.
//
// Returns:
//   - A Middleware function that adds logging functionality.
func Logging(logger *logging.ContextLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return WithLogging(logger, next)
	}
}

// Recovery is a middleware that recovers from panics.
// This is a convenience wrapper around WithRecovery that returns a Middleware
// function, making it easier to use with the Chain function.
//
// Parameters:
//   - logger: A context logger for logging panic details.
//
// Returns:
//   - A Middleware function that adds panic recovery functionality.
func Recovery(logger *logging.ContextLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return WithRecovery(logger, next)
	}
}
