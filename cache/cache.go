// Copyright (c) 2025 A Bit of Help, Inc.

// Package cache provides functionality for caching frequently accessed data.
//
// This package implements a simple in-memory cache with expiration.
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Item represents a cached item with its value and expiration time.
// It is a generic struct that can hold any type of value along with its expiration timestamp.
type Item[T any] struct {
	// Value is the cached data of type T
	Value T
	// Expiration is the Unix timestamp in nanoseconds when this item expires
	Expiration int64
}

// Config contains cache configuration parameters.
// It defines the behavior of the cache, including whether it's enabled,
// how long items are kept, the maximum size, and how often cleanup occurs.
type Config struct {
	// Enabled determines if the cache is enabled.
	// If set to false, cache operations become no-ops.
	Enabled bool

	// TTL is the default time-to-live for cache items.
	// Items older than this duration will be automatically removed during cleanup.
	TTL time.Duration

	// MaxSize is the maximum number of items in the cache.
	// When this limit is reached, new items will cause old items to be evicted.
	MaxSize int

	// PurgeInterval is the interval at which expired items are purged.
	// A background goroutine runs at this interval to remove expired items.
	PurgeInterval time.Duration
}

// DefaultConfig returns a default cache configuration with reasonable values.
// The default configuration includes:
//   - Enabled: true (cache is enabled)
//   - TTL: 5 minutes (items expire after 5 minutes)
//   - MaxSize: 1000 (maximum of 1000 items in the cache)
//   - PurgeInterval: 1 minute (expired items are purged every minute)
//
// Returns:
//   - A Config instance with default values.
func DefaultConfig() Config {
	return Config{
		Enabled:       true,
		TTL:           5 * time.Minute,
		MaxSize:       1000,
		PurgeInterval: 1 * time.Minute,
	}
}

// WithEnabled sets whether the cache is enabled.
// If enabled is set to false, cache operations become no-ops.
//
// Parameters:
//   - enabled: A boolean indicating whether the cache should be enabled.
//
// Returns:
//   - A new Config instance with the updated Enabled value.
func (c Config) WithEnabled(enabled bool) Config {
	c.Enabled = enabled
	return c
}

// WithTTL sets the default time-to-live for cache items.
// This determines how long items remain in the cache before they expire.
// If a non-positive value is provided, it will be set to 1 millisecond.
//
// Parameters:
//   - ttl: The time-to-live duration for cache items.
//
// Returns:
//   - A new Config instance with the updated TTL value.
func (c Config) WithTTL(ttl time.Duration) Config {
	if ttl <= 0 {
		ttl = 1 * time.Millisecond
	}
	c.TTL = ttl
	return c
}

// WithMaxSize sets the maximum number of items in the cache.
// When this limit is reached, new items will cause old items to be evicted.
// If a non-positive value is provided, it will be set to 1.
//
// Parameters:
//   - maxSize: The maximum number of items allowed in the cache.
//
// Returns:
//   - A new Config instance with the updated MaxSize value.
func (c Config) WithMaxSize(maxSize int) Config {
	if maxSize <= 0 {
		maxSize = 1
	}
	c.MaxSize = maxSize
	return c
}

// WithPurgeInterval sets the interval at which expired items are purged.
// A background goroutine runs at this interval to remove expired items.
// If a non-positive value is provided, it will be set to 1 millisecond.
//
// Parameters:
//   - purgeInterval: The interval between purge operations.
//
// Returns:
//   - A new Config instance with the updated PurgeInterval value.
func (c Config) WithPurgeInterval(purgeInterval time.Duration) Config {
	if purgeInterval <= 0 {
		purgeInterval = 1 * time.Millisecond
	}
	c.PurgeInterval = purgeInterval
	return c
}

// Options contains additional options for the cache.
// These options are not directly related to the cache behavior itself,
// but provide additional functionality like logging, tracing, and identification.
type Options struct {
	// Logger is used for logging cache operations.
	// If nil, a no-op logger will be used.
	Logger *logging.ContextLogger

	// Tracer is used for tracing cache operations.
	// It provides integration with OpenTelemetry for distributed tracing.
	Tracer telemetry.Tracer

	// Name is the name of the cache.
	// This is useful for identifying the cache in logs and traces,
	// especially when multiple caches are used in the same application.
	Name string
}

// DefaultOptions returns default options for cache operations.
// The default options include:
//   - No logger (a no-op logger will be used)
//   - A no-op tracer (no OpenTelemetry integration)
//   - Name: "default"
//
// Returns:
//   - An Options instance with default values.
func DefaultOptions() Options {
	return Options{
		Logger: nil,
		Tracer: telemetry.NewNoopTracer(),
		Name:   "default",
	}
}

// WithLogger sets the logger for the cache.
// The logger is used to log cache operations, such as initialization,
// evictions, and shutdown events.
//
// Parameters:
//   - logger: A ContextLogger instance for logging cache operations.
//
// Returns:
//   - A new Options instance with the updated Logger value.
func (o Options) WithLogger(logger *logging.ContextLogger) Options {
	o.Logger = logger
	return o
}

// WithOtelTracer returns Options with an OpenTelemetry tracer.
// This allows users to opt-in to OpenTelemetry tracing if they need it.
// The tracer is used to create spans for cache operations, which can be
// viewed in a distributed tracing system.
//
// Parameters:
//   - tracer: An OpenTelemetry trace.Tracer instance.
//
// Returns:
//   - A new Options instance with the provided OpenTelemetry tracer.
func (o Options) WithOtelTracer(tracer trace.Tracer) Options {
	o.Tracer = telemetry.NewOtelTracer(tracer)
	return o
}

// WithName sets the name of the cache.
// The name is used to identify the cache in logs and traces,
// which is especially useful when multiple caches are used in the same application.
//
// Parameters:
//   - name: A string identifier for the cache.
//
// Returns:
//   - A new Options instance with the updated Name value.
func (o Options) WithName(name string) Options {
	o.Name = name
	return o
}

// Cache is a generic in-memory cache with expiration.
// It provides thread-safe operations for storing and retrieving values of any type,
// with automatic expiration and cleanup of expired items. The cache supports
// configurable size limits with eviction strategies, and integrates with
// OpenTelemetry for tracing and monitoring.
//
// Cache is implemented as a generic type, allowing it to store values of any type
// while maintaining type safety.
type Cache[T any] struct {
	// name is the identifier for this cache instance
	name string

	// items is the map that stores the cached data
	items map[string]Item[T]

	// mu protects the items map for concurrent access
	mu sync.RWMutex

	// defaultTTL is the default time-to-live for cache items
	defaultTTL time.Duration

	// maxSize is the maximum number of items allowed in the cache
	maxSize int

	// cleanupInterval is how often the cleanup routine runs
	cleanupInterval time.Duration

	// logger is used for logging cache operations
	logger *logging.ContextLogger

	// tracer is used for tracing cache operations
	tracer telemetry.Tracer

	// stopCleanup is a channel used to signal the cleanup goroutine to stop
	stopCleanup chan bool

	// evictionStrategy determines how items are evicted when the cache is full
	evictionStrategy EvictionStrategy
}

// EvictionStrategy defines the strategy for evicting items when the cache is full.
// It determines which items are removed when the cache reaches its maximum size
// and a new item needs to be added.
type EvictionStrategy int

const (
	// LRU evicts the least recently used item.
	// This strategy removes items that haven't been accessed for the longest time,
	// which is often the most efficient approach for many caching scenarios.
	LRU EvictionStrategy = iota

	// LFU evicts the least frequently used item.
	// This strategy removes items that have been accessed the fewest times,
	// which can be beneficial when access frequency is more important than recency.
	LFU

	// FIFO evicts the first item added to the cache.
	// This strategy implements a simple first-in, first-out queue,
	// removing the oldest items regardless of how often they've been accessed.
	FIFO

	// Random evicts a random item.
	// This strategy provides a simple and computationally inexpensive approach
	// that can work well when access patterns are unpredictable.
	Random
)

// NewCache creates a new cache with the given configuration and options.
// It initializes the cache with the specified parameters and starts a background
// goroutine to periodically clean up expired items. If the cache is disabled
// (config.Enabled is false), it returns nil and no cache operations will be performed.
//
// The function uses the provided logger and tracer from options, or creates
// no-op versions if they are nil. It also logs the initialization process,
// including the cache name, TTL, maximum size, and purge interval.
//
// Type Parameters:
//   - T: The type of values to be stored in the cache.
//
// Parameters:
//   - config: The configuration parameters for the cache.
//   - options: Additional options for the cache, such as logging and tracing.
//
// Returns:
//   - *Cache[T]: A new cache instance configured according to the provided parameters,
//     or nil if the cache is disabled.
func NewCache[T any](config Config, options Options) *Cache[T] {
	if !config.Enabled {
		if options.Logger != nil {
			options.Logger.Info(context.Background(), "Cache is disabled", zap.String("name", options.Name))
		}
		return nil
	}

	// Use the provided logger or create a no-op logger
	logger := options.Logger
	if logger == nil {
		logger = logging.NewContextLogger(zap.NewNop())
	}

	// Use the provided tracer or create a no-op tracer
	tracer := options.Tracer
	if tracer == nil {
		tracer = telemetry.NewNoopTracer()
	}

	logger.Info(context.Background(), "Initializing cache",
		zap.String("name", options.Name),
		zap.Duration("ttl", config.TTL),
		zap.Int("max_size", config.MaxSize),
		zap.Duration("purge_interval", config.PurgeInterval))

	cache := &Cache[T]{
		name:             options.Name,
		items:            make(map[string]Item[T]),
		defaultTTL:       config.TTL,
		maxSize:          config.MaxSize,
		cleanupInterval:  config.PurgeInterval,
		logger:           logger,
		tracer:           tracer,
		stopCleanup:      make(chan bool),
		evictionStrategy: LRU,
	}

	// Start the cleanup goroutine
	go cache.startCleanupTimer()

	logger.Info(context.Background(), "Cache initialized successfully", zap.String("name", options.Name))
	return cache
}

// Set adds an item to the cache with the default expiration time.
// The item will be stored in the cache until it expires or is explicitly deleted.
// If the cache is full and the key doesn't already exist, an existing item will
// be evicted according to the cache's eviction strategy.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method is a no-op.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - key: The key under which to store the value.
//   - value: The value to store in the cache.
func (c *Cache[T]) Set(ctx context.Context, key string, value T) {
	if c == nil {
		return
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = c.tracer.Start(ctx, "cache.Set")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", c.name),
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "set"),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict an item
	if len(c.items) >= c.maxSize && !c.itemExists(key) {
		c.evictItem()
	}

	c.items[key] = Item[T]{
		Value:      value,
		Expiration: time.Now().Add(c.defaultTTL).UnixNano(),
	}

	span.SetAttributes(attribute.Int("cache.size", len(c.items)))
}

// SetWithTTL adds an item to the cache with a custom expiration time.
// This method works like Set but allows specifying a custom time-to-live (TTL)
// for the item, overriding the default TTL configured for the cache.
//
// The item will be stored in the cache until it expires according to the provided TTL
// or is explicitly deleted. If the cache is full and the key doesn't already exist,
// an existing item will be evicted according to the cache's eviction strategy.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method is a no-op.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - key: The key under which to store the value.
//   - value: The value to store in the cache.
//   - ttl: The time-to-live duration for this specific item.
func (c *Cache[T]) SetWithTTL(ctx context.Context, key string, value T, ttl time.Duration) {
	if c == nil {
		return
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = c.tracer.Start(ctx, "cache.SetWithTTL")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", c.name),
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "set_with_ttl"),
		attribute.Int64("cache.ttl_ms", ttl.Milliseconds()),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict an item
	if len(c.items) >= c.maxSize && !c.itemExists(key) {
		c.evictItem()
	}

	c.items[key] = Item[T]{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}

	span.SetAttributes(attribute.Int("cache.size", len(c.items)))
}

// Get retrieves an item from the cache.
// It returns the value and a boolean indicating whether the value was found.
// If the key doesn't exist or the item has expired, the zero value of type T
// and false are returned.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method returns
// the zero value of type T and false.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - key: The key to look up in the cache.
//
// Returns:
//   - T: The value associated with the key, or the zero value of type T if not found.
//   - bool: true if the key was found and the item hasn't expired, false otherwise.
func (c *Cache[T]) Get(ctx context.Context, key string) (T, bool) {
	if c == nil {
		var zero T
		return zero, false
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = c.tracer.Start(ctx, "cache.Get")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", c.name),
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "get"),
	)

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		span.SetAttributes(attribute.Bool("cache.hit", false))
		var zero T
		return zero, false
	}

	// Check if the item has expired
	if time.Now().UnixNano() > item.Expiration {
		span.SetAttributes(attribute.Bool("cache.hit", false))
		span.SetAttributes(attribute.Bool("cache.expired", true))
		var zero T
		return zero, false
	}

	span.SetAttributes(attribute.Bool("cache.hit", true))
	return item.Value, true
}

// Delete removes an item from the cache.
// This method removes the item with the specified key from the cache, if it exists.
// If the key doesn't exist in the cache, this method is a no-op.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method is a no-op.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - key: The key of the item to remove from the cache.
func (c *Cache[T]) Delete(ctx context.Context, key string) {
	if c == nil {
		return
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = c.tracer.Start(ctx, "cache.Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", c.name),
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "delete"),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[key]
	delete(c.items, key)

	span.SetAttributes(
		attribute.Bool("cache.found", found),
		attribute.Int("cache.size", len(c.items)),
	)
}

// Clear removes all items from the cache.
// This method removes all items from the cache, effectively resetting it to an empty state.
// It creates a new empty map to replace the existing items map, allowing the garbage
// collector to reclaim the memory used by the old items.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method is a no-op.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
func (c *Cache[T]) Clear(ctx context.Context) {
	if c == nil {
		return
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = c.tracer.Start(ctx, "cache.Clear")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", c.name),
		attribute.String("cache.operation", "clear"),
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	oldSize := len(c.items)
	c.items = make(map[string]Item[T])

	span.SetAttributes(
		attribute.Int("cache.old_size", oldSize),
		attribute.Int("cache.size", 0),
	)
}

// Size returns the number of items in the cache.
// This method counts all items currently in the cache, including those that may have
// expired but haven't been removed by the cleanup process yet. For an accurate count
// of non-expired items, you would need to implement a custom counting method.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method returns 0.
//
// Returns:
//   - int: The number of items currently in the cache.
func (c *Cache[T]) Size() int {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// startCleanupTimer starts the cleanup timer
func (c *Cache[T]) startCleanupTimer() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanup removes expired items from the cache
func (c *Cache[T]) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range c.items {
		if now > v.Expiration {
			delete(c.items, k)
		}
	}
}

// Shutdown stops the cleanup timer and gracefully shuts down the cache.
// This method should be called when the cache is no longer needed to ensure
// that the background cleanup goroutine is properly terminated. It sends a signal
// to the cleanup goroutine to stop and logs the shutdown process.
//
// This method is safe to call multiple times and from multiple goroutines.
// If the cache is nil (which happens when the cache is disabled), this method is a no-op.
//
// It's recommended to call this method when the application is shutting down to prevent
// goroutine leaks and ensure proper resource cleanup.
func (c *Cache[T]) Shutdown() {
	if c == nil {
		return
	}

	c.logger.Info(context.Background(), "Shutting down cache", zap.String("name", c.name))
	c.stopCleanup <- true
	c.logger.Info(context.Background(), "Cache shut down successfully", zap.String("name", c.name))
}

// itemExists checks if an item exists in the cache
func (c *Cache[T]) itemExists(key string) bool {
	_, found := c.items[key]
	return found
}

// evictItem evicts an item from the cache based on the eviction strategy
func (c *Cache[T]) evictItem() {
	// For now, just evict a random item
	// In a real implementation, this would use the eviction strategy
	for k := range c.items {
		delete(c.items, k)
		return
	}
}

// WithCache executes a function with caching.
// This utility function implements a caching middleware pattern. It first tries to retrieve
// the value from the cache using the provided key. If the value is found in the cache, it is
// returned immediately without executing the function. If the value is not found, the function
// is executed, and its result is stored in the cache before being returned.
//
// This pattern is useful for expensive operations that are called frequently with the same
// parameters, such as database queries or API calls. The cached results use the default TTL
// configured for the cache. For custom TTL, use WithCacheTTL instead.
//
// If the cache is nil (which happens when the cache is disabled), the function is executed
// directly without any caching.
//
// Type Parameters:
//   - T: The type of value to be cached and returned by the function.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - cache: The cache instance to use for storing and retrieving values.
//   - key: The key under which to store and retrieve the value in the cache.
//   - fn: The function to execute if the value is not found in the cache.
//
// Returns:
//   - T: The value retrieved from the cache or returned by the function.
//   - error: An error if the function execution fails, or nil if successful.
func WithCache[T any](ctx context.Context, cache *Cache[T], key string, fn func(ctx context.Context) (T, error)) (T, error) {
	if cache == nil {
		return fn(ctx)
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = cache.tracer.Start(ctx, "cache.WithCache")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", cache.name),
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "with_cache"),
	)

	// Try to get from cache
	if value, found := cache.Get(ctx, key); found {
		span.SetAttributes(attribute.Bool("cache.hit", true))
		return value, nil
	}

	span.SetAttributes(attribute.Bool("cache.hit", false))

	// If not found, call the function
	value, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		return value, err
	}

	// Store the result in cache
	cache.Set(ctx, key, value)
	return value, nil
}

// WithCacheTTL executes a function with caching and a custom time-to-live.
// This utility function is similar to WithCache but allows specifying a custom TTL
// for the cached result. It first tries to retrieve the value from the cache using
// the provided key. If the value is found in the cache, it is returned immediately
// without executing the function. If the value is not found, the function is executed,
// and its result is stored in the cache with the specified TTL before being returned.
//
// This pattern is useful for expensive operations that are called frequently with the same
// parameters, such as database queries or API calls. The custom TTL allows for fine-grained
// control over how long results should be cached, which can be useful for data with
// different freshness requirements.
//
// If the cache is nil (which happens when the cache is disabled), the function is executed
// directly without any caching.
//
// Type Parameters:
//   - T: The type of value to be cached and returned by the function.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - cache: The cache instance to use for storing and retrieving values.
//   - key: The key under which to store and retrieve the value in the cache.
//   - ttl: The time-to-live duration for the cached value.
//   - fn: The function to execute if the value is not found in the cache.
//
// Returns:
//   - T: The value retrieved from the cache or returned by the function.
//   - error: An error if the function execution fails, or nil if successful.
func WithCacheTTL[T any](ctx context.Context, cache *Cache[T], key string, ttl time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	if cache == nil {
		return fn(ctx)
	}

	// Create a span for the cache operation
	var span telemetry.Span
	ctx, span = cache.tracer.Start(ctx, "cache.WithCacheTTL")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.name", cache.name),
		attribute.String("cache.key", key),
		attribute.String("cache.operation", "with_cache_ttl"),
		attribute.Int64("cache.ttl_ms", ttl.Milliseconds()),
	)

	// Try to get from cache
	if value, found := cache.Get(ctx, key); found {
		span.SetAttributes(attribute.Bool("cache.hit", true))
		return value, nil
	}

	span.SetAttributes(attribute.Bool("cache.hit", false))

	// If not found, call the function
	value, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		return value, err
	}

	// Store the result in cache with custom TTL
	cache.SetWithTTL(ctx, key, value, ttl)
	return value, nil
}
