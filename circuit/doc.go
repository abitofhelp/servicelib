// Copyright (c) 2025 A Bit of Help, Inc.

// Package circuit provides functionality for implementing the circuit breaker pattern.
//
// The circuit breaker pattern is a design pattern used to detect failures and prevent
// cascading failures in distributed systems. It works by "breaking the circuit" when
// a certain threshold of failures is reached, preventing further requests from being
// made to a failing component until it has had time to recover.
//
// This package implements a configurable circuit breaker with three states:
//   - Closed: The circuit is closed and requests are allowed through normally.
//   - Open: The circuit is open and requests are immediately rejected without being attempted.
//   - HalfOpen: After a configurable sleep window, the circuit transitions to half-open state,
//     allowing a limited number of requests through to test if the dependency is healthy.
//
// Key features:
//   - Configurable error threshold and volume threshold for tripping the circuit
//   - Automatic recovery with configurable sleep window
//   - Support for fallback functions when the circuit is open
//   - Integration with OpenTelemetry for tracing
//   - Comprehensive logging of circuit state changes
//
// Example usage:
//
//	// Create a circuit breaker with default configuration
//	cb := circuit.NewCircuitBreaker(circuit.DefaultConfig(), circuit.DefaultOptions())
//
//	// Execute a function with circuit breaking
//	result, err := circuit.Execute(ctx, cb, "database_query", func(ctx context.Context) (string, error) {
//	    return db.Query(ctx, "SELECT * FROM users")
//	})
//
//	// Execute a function with circuit breaking and fallback
//	result, err := circuit.ExecuteWithFallback(
//	    ctx,
//	    cb,
//	    "api_call",
//	    func(ctx context.Context) (string, error) {
//	        return api.Call(ctx, "get_data")
//	    },
//	    func(ctx context.Context, err error) (string, error) {
//	        return "fallback_data", nil
//	    },
//	)
//
// The circuit package is designed to be used as a dependency by other packages in the application,
// providing a consistent circuit breaking interface throughout the codebase.
package circuit
