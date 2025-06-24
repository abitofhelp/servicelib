# Cache Package

The `cache` package provides a generic in-memory cache implementation with expiration for frequently accessed data.

## Overview

This package implements a simple in-memory cache with the following features:

- Generic support for caching any type of data
- Configurable time-to-live (TTL) for cache items
- Automatic cleanup of expired items
- Maximum size limit with configurable eviction strategies
- Support for OpenTelemetry tracing
- Fluent interface for configuration
- Thread-safe implementation
- Middleware functions for easy integration

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/abitofhelp/servicelib/cache"
    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
    "time"
)

// Create a cache
cfg := cache.DefaultConfig().
    WithEnabled(true).
    WithTTL(5 * time.Minute).
    WithMaxSize(1000)

logger := logging.NewContextLogger(zap.NewNop())
options := cache.DefaultOptions().
    WithName("user-cache").
    WithLogger(logger)

// Create a cache for User objects
userCache := cache.NewCache[User](cfg, options)

// Set a value in the cache
ctx := context.Background()
userCache.Set(ctx, "user:123", user)

// Get a value from the cache
user, found := userCache.Get(ctx, "user:123")
if found {
    // Use the cached user
} else {
    // User not found in cache
}

// Set a value with a custom TTL
userCache.SetWithTTL(ctx, "user:456", user, 10 * time.Minute)

// Delete a value from the cache
userCache.Delete(ctx, "user:123")

// Clear the entire cache
userCache.Clear(ctx)

// Get the number of items in the cache
count := userCache.Size()

// Shutdown the cache (stops the cleanup goroutine)
userCache.Shutdown()
```

### Using Cache Middleware

The package provides middleware functions to simplify caching function results:

```go
// Cache the result of a function
user, err := cache.WithCache(ctx, userCache, "user:123", func(ctx context.Context) (User, error) {
    // This function will only be called if the key is not in the cache
    return userService.GetUser(ctx, "123")
})

// Cache the result with a custom TTL
user, err := cache.WithCacheTTL(ctx, userCache, "user:123", 10 * time.Minute, func(ctx context.Context) (User, error) {
    return userService.GetUser(ctx, "123")
})
```

## Configuration

The cache can be configured using the `Config` struct and the fluent interface:

```go
cfg := cache.DefaultConfig().
    WithEnabled(true).                  // Enable/disable the cache
    WithTTL(5 * time.Minute).           // Default time-to-live for cache items
    WithMaxSize(1000).                  // Maximum number of items in the cache
    WithPurgeInterval(1 * time.Minute)  // Interval at which expired items are purged
```

## Eviction Strategies

The cache supports different eviction strategies when the maximum size is reached:

- **LRU (Least Recently Used)**: Evicts the least recently used item
- **LFU (Least Frequently Used)**: Evicts the least frequently used item
- **FIFO (First In, First Out)**: Evicts the first item added to the cache
- **Random**: Evicts a random item

## OpenTelemetry Integration

The cache supports OpenTelemetry tracing:

```go
import (
    "go.opentelemetry.io/otel/trace"
)

// Create a tracer
tracer := otelTracer // Your OpenTelemetry tracer

// Configure the cache with the tracer
options := cache.DefaultOptions().
    WithName("user-cache").
    WithLogger(logger).
    WithOtelTracer(tracer)

userCache := cache.NewCache[User](cfg, options)
```

## Thread Safety

The cache is thread-safe and can be used concurrently from multiple goroutines.

## Performance Considerations

- The cache uses a map internally, so lookups are O(1)
- The cache uses a read-write mutex to ensure thread safety, so concurrent reads are allowed but writes are exclusive
- The cleanup goroutine runs at the configured interval to remove expired items
- When the maximum size is reached, an item will be evicted according to the configured eviction strategy