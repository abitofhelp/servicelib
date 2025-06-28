// Copyright (c) 2025 A Bit of Help, Inc.

// Package rate provides functionality for rate limiting to protect resources.
//
// This package implements a token bucket rate limiter that helps protect services
// and resources from being overwhelmed by too many requests. Rate limiting is
// essential for maintaining system stability, preventing resource exhaustion,
// and ensuring fair usage of shared resources.
//
// The rate limiter uses a token bucket algorithm where:
//   - Tokens are added to the bucket at a fixed rate (RequestsPerSecond)
//   - Each request consumes one token
//   - If tokens are available, the request is allowed
//   - If no tokens are available, the request is rejected or delayed
//   - The bucket has a maximum capacity (BurstSize) to allow for bursts of traffic
//
// Key features:
//   - Configurable requests per second and burst size
//   - Support for blocking and non-blocking rate limiting
//   - Integration with OpenTelemetry for tracing
//   - Comprehensive logging of rate limiting decisions
//   - Thread-safe implementation for concurrent use
//
// Example usage:
//
//	// Create a rate limiter with default configuration
//	limiter := rate.NewRateLimiter(rate.DefaultConfig(), rate.DefaultOptions())
//
//	// Execute a function with rate limiting (non-blocking)
//	result, err := rate.Execute(ctx, limiter, "fetch_data", func(ctx context.Context) (string, error) {
//	    // Operation to be rate limited
//	    return fetchData(ctx)
//	})
//
//	// Execute a function with rate limiting (blocking until allowed)
//	result, err := rate.ExecuteWithWait(ctx, limiter, "process_item", func(ctx context.Context) (string, error) {
//	    // Operation to be rate limited
//	    return processItem(ctx)
//	})
//
// The package is designed to be flexible and can be used to rate limit various
// types of operations, including API calls, database queries, and resource-intensive
// computations.
package rate
