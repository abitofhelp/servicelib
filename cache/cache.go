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

// Cache is a simple in-memory cache with expiration
type Cache[T any] struct {
	name             string
	items            map[string]Item[T]
	mu               sync.RWMutex
	defaultTTL       time.Duration
	maxSize          int
	cleanupInterval  time.Duration
	logger           *logging.ContextLogger
	tracer           telemetry.Tracer
	stopCleanup      chan bool
	evictionStrategy EvictionStrategy
}

// EvictionStrategy defines the strategy for evicting items when the cache is full
type EvictionStrategy int

const (
	// LRU evicts the least recently used item
	LRU EvictionStrategy = iota
	// LFU evicts the least frequently used item
	LFU
	// FIFO evicts the first item added to the cache
	FIFO
	// Random evicts a random item
	Random
)

// NewCache creates a new cache with the given configuration
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

// Set adds an item to the cache with the default expiration time
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

// SetWithTTL adds an item to the cache with a custom expiration time
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

// Get retrieves an item from the cache
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

// Delete removes an item from the cache
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

// Clear removes all items from the cache
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

// Size returns the number of items in the cache
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

// Shutdown stops the cleanup timer
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

// WithCache is a middleware that adds caching to a function
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

// WithCacheTTL is a middleware that adds caching with a custom TTL to a function
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
