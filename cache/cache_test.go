// Copyright (c) 2025 A Bit of Help, Inc.

package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		options  Options
		expected bool
	}{
		{
			name: "Enabled cache",
			config: DefaultConfig().
				WithEnabled(true),
			options:  DefaultOptions(),
			expected: true,
		},
		{
			name: "Disabled cache",
			config: DefaultConfig().
				WithEnabled(false),
			options:  DefaultOptions(),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCache[string](tc.config, tc.options)
			if tc.expected {
				assert.NotNil(t, cache)
			} else {
				assert.Nil(t, cache)
			}
		})
	}
}

func TestCache_SetAndGet(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true).
		WithTTL(1 * time.Hour)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Set a value
	cache.Set(ctx, "key1", "value1")

	// Get the value
	value, found := cache.Get(ctx, "key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	// Get a non-existent value
	value, found = cache.Get(ctx, "key2")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_SetWithTTL(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true).
		WithTTL(1 * time.Hour)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Set a value with a short TTL
	cache.SetWithTTL(ctx, "key1", "value1", 10*time.Millisecond)

	// Get the value immediately
	value, found := cache.Get(ctx, "key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	// Wait for the TTL to expire
	time.Sleep(20 * time.Millisecond)

	// Get the value after expiration
	value, found = cache.Get(ctx, "key1")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_Delete(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Set a value
	cache.Set(ctx, "key1", "value1")

	// Delete the value
	cache.Delete(ctx, "key1")

	// Get the deleted value
	value, found := cache.Get(ctx, "key1")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_Clear(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Set multiple values
	cache.Set(ctx, "key1", "value1")
	cache.Set(ctx, "key2", "value2")
	cache.Set(ctx, "key3", "value3")

	// Verify the cache size
	assert.Equal(t, 3, cache.Size())

	// Clear the cache
	cache.Clear(ctx)

	// Verify the cache is empty
	assert.Equal(t, 0, cache.Size())

	// Get a value after clearing
	value, found := cache.Get(ctx, "key1")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_Size(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Initially empty
	assert.Equal(t, 0, cache.Size())

	// Add items
	cache.Set(ctx, "key1", "value1")
	assert.Equal(t, 1, cache.Size())

	cache.Set(ctx, "key2", "value2")
	assert.Equal(t, 2, cache.Size())

	// Delete an item
	cache.Delete(ctx, "key1")
	assert.Equal(t, 1, cache.Size())

	// Clear the cache
	cache.Clear(ctx)
	assert.Equal(t, 0, cache.Size())
}

func TestCache_MaxSize(t *testing.T) {
	// Create a cache with a small max size
	cfg := DefaultConfig().
		WithEnabled(true).
		WithMaxSize(2)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Add items up to max size
	cache.Set(ctx, "key1", "value1")
	cache.Set(ctx, "key2", "value2")
	assert.Equal(t, 2, cache.Size())

	// Add one more item, should evict an existing item
	cache.Set(ctx, "key3", "value3")
	assert.Equal(t, 2, cache.Size())

	// At least one of the original keys should be gone
	count := 0
	for _, key := range []string{"key1", "key2"} {
		if _, found := cache.Get(ctx, key); found {
			count++
		}
	}
	assert.Less(t, count, 2)

	// The new key should be present
	value, found := cache.Get(ctx, "key3")
	assert.True(t, found)
	assert.Equal(t, "value3", value)
}

func TestCache_Cleanup(t *testing.T) {
	// Create a cache with a short cleanup interval
	cfg := DefaultConfig().
		WithEnabled(true).
		WithTTL(10 * time.Millisecond).
		WithPurgeInterval(20 * time.Millisecond)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Set a value
	cache.Set(ctx, "key1", "value1")
	assert.Equal(t, 1, cache.Size())

	// Wait for the cleanup to run
	time.Sleep(50 * time.Millisecond)

	// The item should be gone
	assert.Equal(t, 0, cache.Size())
}

func TestCache_Shutdown(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	// Shutdown should not panic
	assert.NotPanics(t, func() {
		cache.Shutdown()
	})
}

func TestWithCache(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Define a function that returns a value
	callCount := 0
	fn := func(ctx context.Context) (string, error) {
		callCount++
		return "computed value", nil
	}

	// First call should compute the value
	result, err := WithCache(ctx, cache, "key1", fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 1, callCount)

	// Second call should use the cached value
	result, err = WithCache(ctx, cache, "key1", fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 1, callCount) // Function should not be called again
}

func TestWithCache_Error(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Define a function that returns an error
	fn := func(ctx context.Context) (string, error) {
		return "", errors.New("test error")
	}

	// Call should return the error and not cache anything
	result, err := WithCache(ctx, cache, "key1", fn)
	assert.Error(t, err)
	assert.Equal(t, "", result)

	// Verify nothing was cached
	value, found := cache.Get(ctx, "key1")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestWithCacheTTL(t *testing.T) {
	// Create a cache
	cfg := DefaultConfig().
		WithEnabled(true)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cache := NewCache[string](cfg, options)
	assert.NotNil(t, cache)

	ctx := context.Background()

	// Define a function that returns a value
	callCount := 0
	fn := func(ctx context.Context) (string, error) {
		callCount++
		return "computed value", nil
	}

	// First call should compute the value with a short TTL
	result, err := WithCacheTTL(ctx, cache, "key1", 10*time.Millisecond, fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 1, callCount)

	// Second call should use the cached value
	result, err = WithCacheTTL(ctx, cache, "key1", 10*time.Millisecond, fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 1, callCount) // Function should not be called again

	// Wait for the TTL to expire
	time.Sleep(20 * time.Millisecond)

	// Third call should compute the value again
	result, err = WithCacheTTL(ctx, cache, "key1", 10*time.Millisecond, fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 2, callCount) // Function should be called again
}

func TestCache_NilSafety(t *testing.T) {
	// Test with nil cache
	var cache *Cache[string]

	ctx := context.Background()

	// Set should not panic
	assert.NotPanics(t, func() {
		cache.Set(ctx, "key1", "value1")
	})

	// Get should return not found
	value, found := cache.Get(ctx, "key1")
	assert.False(t, found)
	assert.Equal(t, "", value)

	// Delete should not panic
	assert.NotPanics(t, func() {
		cache.Delete(ctx, "key1")
	})

	// Clear should not panic
	assert.NotPanics(t, func() {
		cache.Clear(ctx)
	})

	// Size should return 0
	assert.Equal(t, 0, cache.Size())

	// Shutdown should not panic
	assert.NotPanics(t, func() {
		cache.Shutdown()
	})

	// WithCache should call the function
	callCount := 0
	fn := func(ctx context.Context) (string, error) {
		callCount++
		return "computed value", nil
	}

	result, err := WithCache(ctx, cache, "key1", fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 1, callCount)

	// WithCacheTTL should call the function
	result, err = WithCacheTTL(ctx, cache, "key1", 10*time.Millisecond, fn)
	assert.NoError(t, err)
	assert.Equal(t, "computed value", result)
	assert.Equal(t, 2, callCount)
}
