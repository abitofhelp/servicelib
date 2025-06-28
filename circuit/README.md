# Circuit Breaker

## Overview

The Circuit Breaker component provides a robust implementation of the Circuit Breaker pattern for handling failures in distributed systems. It helps prevent cascading failures by temporarily disabling operations that are likely to fail.

## Features

- **Failure Detection**: Automatically detects when operations are failing
- **Automatic Recovery**: Automatically recovers after a specified sleep window
- **Configurable Thresholds**: Customize error and volume thresholds
- **Concurrent Operation Limiting**: Limit the number of concurrent operations
- **Fallback Support**: Provide fallback mechanisms for when operations fail

## Installation

```bash
go get github.com/abitofhelp/servicelib/circuit
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/circuit/basic_usage/README.md) for a complete, runnable example of how to use the circuit breaker component.

## Configuration

See the [Custom Configuration example](../EXAMPLES/circuit/custom_configuration/README.md) for a complete, runnable example of how to configure the circuit breaker component.

## API Documentation

### Core Types

The circuit breaker component provides several core types for implementing the Circuit Breaker pattern.

#### CircuitBreaker

The main type that implements the Circuit Breaker pattern.

```
type CircuitBreaker struct {
    // Fields
}
```

#### Config

Configuration for the CircuitBreaker.

```
type Config struct {
    Enabled         bool
    Timeout         time.Duration
    MaxConcurrent   int
    ErrorThreshold  float64
    VolumeThreshold int
    SleepWindow     time.Duration
}
```

#### State

Represents the state of the circuit breaker.

```
type State int
```

### Key Methods

The circuit breaker component provides several key methods for implementing the Circuit Breaker pattern.

#### Execute

Executes an operation with circuit breaking functionality.

```
func Execute[T any](ctx context.Context, cb *CircuitBreaker, operation string, fn func(ctx context.Context) (T, error)) (T, error)
```

#### ExecuteWithFallback

Executes an operation with circuit breaking functionality and a fallback mechanism.

```
func ExecuteWithFallback[T any](ctx context.Context, cb *CircuitBreaker, operation string, fn func(ctx context.Context) (T, error), fallback func(ctx context.Context, err error) (T, error)) (T, error)
```

#### GetState

Gets the current state of the circuit breaker.

```
func (cb *CircuitBreaker) GetState() State
```

#### Reset

Resets the circuit breaker to its initial state.

```
func (cb *CircuitBreaker) Reset()
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/circuit/basic_usage/README.md) - Shows basic usage of the circuit breaker
- [Custom Configuration](../EXAMPLES/circuit/custom_configuration/README.md) - Shows how to configure the circuit breaker

## Best Practices

1. **Use Fallbacks**: Always provide fallback mechanisms for when operations fail
2. **Configure Thresholds**: Configure error and volume thresholds based on your application's needs
3. **Monitor States**: Monitor the state of your circuit breakers to detect issues
4. **Use Timeouts**: Always set appropriate timeouts for your operations
5. **Handle Errors**: Properly handle errors returned by the circuit breaker

## Troubleshooting

### Common Issues

#### Circuit Always Open

If the circuit is always open, check that your error threshold and volume threshold are not too low.

#### Circuit Never Opens

If the circuit never opens despite failures, check that your error threshold and volume threshold are not too high.

## Related Components

- [Errors](../errors/README.md) - Error handling for circuit breaker operations
- [Logging](../logging/README.md) - Logging for circuit breaker events
- [Telemetry](../telemetry/README.md) - Telemetry for circuit breaker operations

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
