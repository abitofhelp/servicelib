# Cache

## Overview

The Cache component provides a generic in-memory cache implementation for storing and retrieving data with configurable Time-To-Live (TTL), size limits, and automatic cleanup.

## Features

- **Generic Implementation**: Works with any data type
- **TTL Support**: Automatically expires items after a configurable time
- **Size Limiting**: Limits the maximum number of items in the cache
- **Automatic Cleanup**: Periodically removes expired items
- **Thread-Safe**: Safe for concurrent use

## Installation

```bash
go get github.com/abitofhelp/servicelib/cache
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/cache/basic_usage/README.md) for a complete, runnable example of how to use the cache component.

## Configuration

See the [Custom Configuration example](../EXAMPLES/cache/custom_configuration/README.md) for a complete, runnable example of how to configure the cache component.

## API Documentation

### Core Types

The cache component provides several core types for implementing caching functionality.

#### Cache

The main type that implements the cache functionality.

```
type Cache[T any] struct {
    // Fields
}
```

#### Config

Configuration for the Cache.

```
type Config struct {
    Enabled        bool
    TTL            time.Duration
    MaxSize        int
    PurgeInterval  time.Duration
}
```

#### Item

Represents an item in the cache.

```
type Item struct {
    Value      interface{}
    Expiration int64
}
```

### Key Methods

The cache component provides several key methods for working with the cache.

#### Set

Sets a value in the cache with the default TTL.

```
func (c *Cache[T]) Set(ctx context.Context, key string, value T)
```

#### SetWithTTL

Sets a value in the cache with a custom TTL.

```
func (c *Cache[T]) SetWithTTL(ctx context.Context, key string, value T, ttl time.Duration)
```

#### Get

Gets a value from the cache.

```
func (c *Cache[T]) Get(ctx context.Context, key string) (T, bool)
```

#### Delete

Deletes a value from the cache.

```
func (c *Cache[T]) Delete(ctx context.Context, key string)
```

#### Clear

Clears all values from the cache.

```
func (c *Cache[T]) Clear(ctx context.Context)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/cache/basic_usage/README.md) - Shows basic usage of the cache
- [Custom Configuration](../EXAMPLES/cache/custom_configuration/README.md) - Shows how to configure the cache

## Best Practices

1. **Use Appropriate TTL**: Set appropriate TTL values based on your data's freshness requirements
2. **Limit Cache Size**: Set a reasonable maximum size to prevent memory issues
3. **Handle Cache Misses**: Always handle the case when an item is not found in the cache
4. **Use WithCache Helper**: Use the WithCache helper function for common cache patterns
5. **Shutdown Properly**: Call Shutdown() when you're done with the cache to stop background cleanup

## Troubleshooting

### Common Issues

#### Cache Items Expiring Too Quickly

If cache items are expiring too quickly, check that your TTL value is appropriate for your use case.

#### Cache Using Too Much Memory

If the cache is using too much memory, reduce the MaxSize configuration parameter.

## Related Components

- [Logging](../logging/README.md) - Logging for cache operations
- [Telemetry](../telemetry/README.md) - Telemetry for cache operations

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
