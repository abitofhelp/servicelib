// Copyright (c) 2025 A Bit of Help, Inc.

// Package cache provides functionality for caching frequently accessed data.
//
// This package implements a generic in-memory cache with expiration, supporting
// various eviction strategies and thread-safe operations. It is designed to be
// flexible and easy to use, with configurable time-to-live (TTL) for cache items,
// maximum size limits, and automatic cleanup of expired items.
//
// Key features:
//   - Generic implementation that works with any data type
//   - Configurable item expiration (TTL)
//   - Maximum size limits with various eviction strategies
//   - Thread-safe operations for concurrent access
//   - Automatic cleanup of expired items
//   - Integration with OpenTelemetry for tracing
//   - Comprehensive logging of cache operations
//
// The package provides several main components:
//   - Cache: A generic in-memory cache with expiration
//   - Config: Configuration for cache behavior
//   - Options: Additional options for logging and tracing
//   - EvictionStrategy: Different strategies for evicting items when the cache is full
//
// Example usage:
//
//	// Create a cache for string values with default configuration
//	cache := cache.NewCache[string](cache.DefaultConfig(), cache.DefaultOptions())
//
//	// Store a value in the cache
//	cache.Set(ctx, "key1", "value1")
//
//	// Retrieve a value from the cache
//	value, found := cache.Get(ctx, "key1")
//	if found {
//	    fmt.Println("Value:", value)
//	}
//
//	// Store a value with a custom TTL
//	cache.SetWithTTL(ctx, "key2", "value2", 5*time.Minute)
//
//	// Delete a value from the cache
//	cache.Delete(ctx, "key1")
//
//	// Execute a function with caching
//	result, err := cache.WithCache(ctx, "key3", func(ctx context.Context) (string, error) {
//	    // Expensive operation to generate the value
//	    return "computed-value", nil
//	}, 10*time.Minute)
//
// The cache package is designed to be used as a dependency by other packages in the application,
// providing a consistent caching interface throughout the codebase.
package cache
