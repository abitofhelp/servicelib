# Cache Package

## Overview

The `cache` package provides a generic in-memory cache implementation with expiration for frequently accessed data. This package implements a simple in-memory cache that can be used to cache any type of data with configurable time-to-live (TTL) and automatic cleanup of expired items.


## Features

- **Generic Support**: Cache any type of data using Go generics
- **Configurable TTL**: Set default and per-item time-to-live
- **Automatic Cleanup**: Expired items are automatically removed
- **Size Limits**: Configure maximum cache size with eviction strategies
- **Thread Safety**: Safe for concurrent use from multiple goroutines
- **OpenTelemetry Integration**: Built-in support for distributed tracing
- **Fluent Configuration**: Easy-to-use fluent interface for configuration
- **Middleware Functions**: Simplify caching of function results


## Installation

```bash
go get github.com/abitofhelp/servicelib/cache
```


## Quick Start

See the [Basic Usage Example](../EXAMPLES/cache/basic_usage/main.go) for a complete, runnable example of how to use the cache package.


## Configuration

The cache can be configured using the `Config` struct and the fluent interface.

See the [Custom Configuration Example](../EXAMPLES/cache/custom_configuration/main.go) for a complete, runnable example of how to configure the cache.


## API Documentation


### Core Types

#### Cache

The `Cache` struct is the main type in the cache package. It provides methods for storing and retrieving values from the cache.

```
// Cache is a generic in-memory cache with expiration
type Cache[T any] struct {
    // contains filtered or unexported fields
}
```

See the [Basic Usage Example](../EXAMPLES/cache/basic_usage/main.go) for a complete, runnable example of how to use the Cache type.

#### Config

The `Config` struct holds the configuration for the cache.

```
// Config holds the configuration for the cache
type Config struct {
    // contains filtered or unexported fields
}
```

See the [Custom Configuration Example](../EXAMPLES/cache/custom_configuration/main.go) for a complete, runnable example of how to use the Config type.

#### Options

The `Options` struct holds additional options for the cache.

```
// Options holds additional options for the cache
type Options struct {
    // contains filtered or unexported fields
}
```


### Key Methods

#### NewCache

`NewCache` creates a new cache with the specified configuration and options.

```
// NewCache creates a new cache with the specified configuration and options
func NewCache[T any](cfg Config, options Options) *Cache[T]
```

#### Get

`Get` retrieves a value from the cache.

```
// Get retrieves a value from the cache
func (c *Cache[T]) Get(ctx context.Context, key string) (T, bool)
```

#### Set

`Set` stores a value in the cache with the default TTL.

```
// Set stores a value in the cache with the default TTL
func (c *Cache[T]) Set(ctx context.Context, key string, value T)
```

#### SetWithTTL

`SetWithTTL` stores a value in the cache with a custom TTL.

```
// SetWithTTL stores a value in the cache with a custom TTL
func (c *Cache[T]) SetWithTTL(ctx context.Context, key string, value T, ttl time.Duration)
```

#### Delete

`Delete` removes a value from the cache.

```
// Delete removes a value from the cache
func (c *Cache[T]) Delete(ctx context.Context, key string)
```

#### Clear

`Clear` removes all values from the cache.

```
// Clear removes all values from the cache
func (c *Cache[T]) Clear(ctx context.Context)
```

#### Size

`Size` returns the number of items in the cache.

```
// Size returns the number of items in the cache
func (c *Cache[T]) Size() int
```

#### Shutdown

`Shutdown` stops the cleanup goroutine.

```
// Shutdown stops the cleanup goroutine
func (c *Cache[T]) Shutdown()
```

#### WithCache

`WithCache` is a middleware function that caches the result of a function.

```
// WithCache caches the result of a function
func WithCache[T any, E error](ctx context.Context, cache *Cache[T], key string, fn func(ctx context.Context) (T, E)) (T, E)
```

#### WithCacheTTL

`WithCacheTTL` is a middleware function that caches the result of a function with a custom TTL.

```
// WithCacheTTL caches the result of a function with a custom TTL
func WithCacheTTL[T any, E error](ctx context.Context, cache *Cache[T], key string, ttl time.Duration, fn func(ctx context.Context) (T, E)) (T, E)
```


## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Usage](../EXAMPLES/cache/basic_usage/main.go) - Shows basic usage of the cache package
- [Custom Configuration](../EXAMPLES/cache/custom_configuration/main.go) - Shows advanced configuration options


## Best Practices

1. **Choose Appropriate TTL**: Set TTL values based on how frequently your data changes. Shorter TTLs for frequently changing data, longer TTLs for more static data.

2. **Use Meaningful Keys**: Choose cache keys that are meaningful and unique to avoid collisions.

3. **Handle Cache Misses Gracefully**: Always handle the case where a value is not found in the cache.

4. **Monitor Cache Performance**: Keep track of cache hit rates and adjust TTL values accordingly.

5. **Set Size Limits**: Configure a maximum size for your cache to prevent memory issues.

6. **Use Middleware Functions**: Use the provided middleware functions to simplify caching of function results.


## Troubleshooting

### Common Issues

#### Cache Misses

**Issue**: Values are not being found in the cache when expected.

**Solution**: Check that the keys being used are consistent and that the TTL is not too short. Also, ensure that the cache is enabled.

#### Memory Usage

**Issue**: The cache is using too much memory.

**Solution**: Configure a maximum size for the cache and choose an appropriate eviction strategy. Also, consider using shorter TTL values for large items.

## Eviction Strategies

The cache supports different eviction strategies when the maximum size is reached:

- **LRU (Least Recently Used)**: Evicts the least recently used item
- **LFU (Least Frequently Used)**: Evicts the least frequently used item
- **FIFO (First In, First Out)**: Evicts the first item added to the cache
- **Random**: Evicts a random item


## Related Components

- [Logging](../logging/README.md) - The logging component is used by the cache for logging operations.
- [Telemetry](../telemetry/README.md) - The telemetry component provides tracing, which is used by the cache for tracing operations.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
