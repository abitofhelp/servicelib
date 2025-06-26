// Copyright (c) 2025 A Bit of Help, Inc.

// Package middleware provides a collection of HTTP middleware for Go web applications.
//
// Middleware are functions that wrap HTTP handlers to add functionality such as
// request logging, error handling, timeout management, and CORS support. This package
// offers ready-to-use middleware for common tasks and utilities for middleware chaining
// and context management.
//
// Key features:
//   - Request Context: Add request IDs and timing information to request contexts
//   - Logging: Log request details including method, path, status code, and duration
//   - Error Handling: Map errors to appropriate HTTP responses with status codes
//   - Panic Recovery: Catch and handle panics to prevent application crashes
//   - Timeout Management: Add request timeouts with proper cancellation handling
//   - CORS Support: Add Cross-Origin Resource Sharing headers to responses
//   - Context Cancellation: Detect and handle client disconnections
//   - Middleware Chaining: Apply multiple middleware in a specific order
//   - Thread Safety: Thread-safe response writing for concurrent operations
//
// The package provides both individual middleware functions and utilities for
// combining them. The ApplyMiddleware function applies a standard set of middleware
// in the recommended order, while the Chain function allows for custom combinations.
//
// Example usage:
//
//	// Create a logger
//	logger, _ := zap.NewProduction()
//
//	// Create a simple handler
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("Hello, World!"))
//	})
//
//	// Apply middleware
//	wrappedHandler := middleware.ApplyMiddleware(handler, logger)
//
//	// Start the server
//	http.ListenAndServe(":8080", wrappedHandler)
//
// For more advanced usage, middleware can be applied individually or chained
// in a custom order:
//
//	// Create a logger
//	logger, _ := zap.NewProduction()
//	contextLogger := logging.NewContextLogger(logger)
//
//	// Apply middleware in a custom order
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("Hello, World!"))
//	})
//
//	// Chain middleware in a specific order
//	wrappedHandler := middleware.Chain(handler,
//	    middleware.WithRequestID(context.Background()),
//	    middleware.Logging(contextLogger),
//	    middleware.Recovery(contextLogger),
//	)
//
//	// Start the server
//	http.ListenAndServe(":8080", wrappedHandler)
//
// The order of middleware application is important, as it determines the order
// of execution. Middleware is executed in reverse order of application, so the
// first middleware applied is the outermost wrapper (first to process the request,
// last to process the response).
package middleware