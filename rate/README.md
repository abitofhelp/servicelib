# Rate Limiter Package

## Overview

The `rate` package provides a generic implementation of a token bucket rate limiter to protect resources from being overwhelmed by too many requests. Rate limiting is a strategy to control the rate of requests sent or received by a system. It helps prevent resource exhaustion, maintain service quality, and ensure fair usage of shared resources.


## Features

- **Generic Support**: Works with any function that returns a value and an error
- **Token Bucket Algorithm**: Precise rate limiting with configurable parameters
- **Configurable Limits**: Set requests per second and burst size
- **Flexible Handling**: Support for both immediate rejection and waiting for available tokens
- **OpenTelemetry Integration**: Built-in support for distributed tracing
- **Fluent Configuration**: Easy-to-use fluent interface for configuration
- **Thread Safety**: Safe for concurrent use from multiple goroutines
- **Middleware Functions**: Simplify rate limiting of function calls


## Installation

```bash
go get github.com/abitofhelp/servicelib/rate
```


## Quick Start

See the [Basic Usage Example](../EXAMPLES/rate/basic_usage/main.go) for a complete, runnable example of how to use the rate package.


## Configuration

The rate limiter can be configured using the `Config` struct and the fluent interface.

See the [Custom Configuration Example](../EXAMPLES/rate/custom_configuration/main.go) for a complete, runnable example of how to configure the rate limiter.


## API Documentation


### Core Types

#### Config

The `Config` struct holds the configuration for the rate limiter.

```
type Config struct {
    Enabled           bool
    RequestsPerSecond int
    BurstSize         int
}
```

#### Options

The `Options` struct holds additional options for the rate limiter.

```
type Options struct {
    Name       string
    Logger     *logging.ContextLogger
    OtelTracer trace.Tracer
}
```

#### RateLimiter

The `RateLimiter` struct is the main type in the rate package. It provides methods for rate limiting.

```
type RateLimiter struct {
    // contains filtered or unexported fields
}
```


### Key Methods

#### NewRateLimiter

`NewRateLimiter` creates a new rate limiter with the specified configuration and options.

```
func NewRateLimiter(config Config, options Options) *RateLimiter
```

#### Allow

`Allow` checks if a request is allowed based on the rate limit.

```
func (rl *RateLimiter) Allow() bool
```

#### Execute

`Execute` executes a function with rate limiting.

```
func Execute[T any](ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (T, error)) (T, error)
```

#### ExecuteWithWait

`ExecuteWithWait` executes a function with rate limiting and waits for a token if necessary.

```
func ExecuteWithWait[T any](ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (T, error)) (T, error)
```

#### Reset

`Reset` resets the rate limiter to its initial state.

```
func (rl *RateLimiter) Reset()
```


## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Usage](../EXAMPLES/rate/basic_usage/main.go) - Shows basic usage of the rate package
- [Custom Configuration](../EXAMPLES/rate/custom_configuration/main.go) - Shows advanced configuration options


## Best Practices

1. **Choose Appropriate Limits**: Set rate limits based on the capacity of the resources you're protecting. Too low limits can unnecessarily restrict legitimate traffic, while too high limits may not provide adequate protection.

2. **Use Burst Allowance**: Configure a reasonable burst size to handle temporary spikes in traffic without rejecting legitimate requests.

3. **Consider Wait vs. Reject**: For user-facing operations, consider using `ExecuteWithWait` to improve user experience by waiting for a token rather than immediately rejecting the request.

4. **Monitor Rate Limiting**: Log rate limiting events and monitor them to adjust your rate limits as needed.

5. **Use Distributed Rate Limiting**: For distributed systems, consider using a distributed rate limiter that coordinates across multiple instances.


## Troubleshooting

### Common Issues

#### Rate Limit Too Restrictive

**Issue**: Legitimate requests are being rate limited too aggressively.

**Solution**: Increase the requests per second or burst size in the configuration. Monitor the rate of requests and adjust the limits accordingly.

#### Rate Limit Not Effective

**Issue**: Resources are still being overwhelmed despite rate limiting.

**Solution**: Decrease the requests per second or burst size in the configuration. Ensure that all access paths to the resource are protected by rate limiting.

#### High Latency with ExecuteWithWait

**Issue**: Using `ExecuteWithWait` is causing high latency for requests.

**Solution**: Consider using a queue or other asynchronous processing mechanism instead of waiting for tokens. Alternatively, adjust the rate limits to reduce waiting time.


## Related Components

- [Logging](../logging/README.md) - The logging component is used by the rate limiter for logging operations.
- [Telemetry](../telemetry/README.md) - The telemetry component provides tracing, which is used by the rate limiter for tracing operations.
- [Errors](../errors/README.md) - The errors component is used by the rate limiter for error handling.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
