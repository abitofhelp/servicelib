// Copyright (c) 2025 A Bit of Help, Inc.

// Package cache provides functionality for caching frequently accessed data.
//
// This package implements a generic in-memory cache with expiration, supporting
// various features for efficient data caching. It is designed to be flexible,
// configurable, and easy to integrate with existing code.
//
// Key features:
//   - Generic implementation using Go generics (Cache[T])
//   - Configurable time-to-live (TTL) for cache items
//   - Maximum cache size with eviction strategies
//   - Automatic cleanup of expired items
//   - OpenTelemetry integration for monitoring and tracing
//   - Thread-safe operations for concurrent access
//   - Helper functions for caching function results
//
// The package supports different eviction strategies when the cache reaches its
// maximum size:
//   - LRU (Least Recently Used)
//   - LFU (Least Frequently Used)
//   - FIFO (First In, First Out)
//   - Random
//
// Example usage:
//
//	// Create a cache configuration
//	config := cache.DefaultConfig().
//	    WithMaxSize(1000).
//	    WithTTL(5 * time.Minute)
//
//	// Create a cache for string values
//	stringCache := cache.NewCache[string](config, cache.DefaultOptions())
//	defer stringCache.Shutdown()
//
//	// Store a value in the cache
//	ctx := context.Background()
//	stringCache.Set(ctx, "key1", "value1")
//
//	// Retrieve a value from the cache
//	value, found := stringCache.Get(ctx, "key1")
//	if found {
//	    fmt.Println("Value:", value)
//	}
//
//	// Use the cache with a function
//	result, err := cache.WithCache(ctx, stringCache, "key2", func(ctx context.Context) (string, error) {
//	    // Expensive operation that returns a string
//	    return fetchDataFromDatabase(ctx)
//	})
//
// The cache package is designed to be used as a dependency by other packages in the
// application, providing a consistent caching interface throughout the codebase.
package cache