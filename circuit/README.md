# Circuit Breaker Package
The `circuit` package provides a generic implementation of the circuit breaker pattern to protect against cascading failures when external dependencies are unavailable.


## Overview

The circuit breaker pattern is a design pattern used in modern software development to detect failures and encapsulate the logic of preventing a failure from constantly recurring during maintenance, temporary external system failure, or unexpected system difficulties.

This implementation provides:

- Generic support for any function that returns a value and an error
- Configurable error thresholds, volume thresholds, and sleep windows
- Support for OpenTelemetry tracing
- Fluent interface for configuration
- Thread-safe implementation
- Support for fallback functions

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/abitofhelp/servicelib/circuit"
    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
)

// Create a circuit breaker
cfg := circuit.DefaultConfig().
    WithEnabled(true).
    WithErrorThreshold(0.5).
    WithVolumeThreshold(10)

logger := logging.NewContextLogger(zap.NewNop())
options := circuit.DefaultOptions().
    WithName("my-service").
    WithLogger(logger)

cb := circuit.NewCircuitBreaker(cfg, options)

// Execute a function with circuit breaking
result, err := circuit.Execute(ctx, cb, "GetUserProfile", func(ctx context.Context) (UserProfile, error) {
    // Call external service
    return userService.GetProfile(ctx, userID)
})
```

### With Fallback

```go
// Execute a function with circuit breaking and fallback
result, err := circuit.ExecuteWithFallback(
    ctx, 
    cb, 
    "GetUserProfile", 
    func(ctx context.Context) (UserProfile, error) {
        // Call external service
        return userService.GetProfile(ctx, userID)
    },
    func(ctx context.Context, err error) (UserProfile, error) {
        // Fallback function
        return UserProfile{Name: "Default User"}, nil
    },
)
```


## Features

- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3

## Installation

```bash
go get github.com/abitofhelp/servicelib/circuit
```

## Quick Start

See the [Quick Start example](../EXAMPLES/circuit/quickstart_example.go) for a complete, runnable example of how to use the circuit.

## Configuration

The circuit breaker can be configured using the `Config` struct and the fluent interface:

```go
cfg := circuit.DefaultConfig().
    WithEnabled(true).                  // Enable/disable the circuit breaker
    WithTimeout(5 * time.Second).       // Maximum time allowed for a request
    WithMaxConcurrent(100).             // Maximum number of concurrent requests
    WithErrorThreshold(0.5).            // Percentage of errors that will trip the circuit (0.0-1.0)
    WithVolumeThreshold(10).            // Minimum number of requests before the error threshold is checked
    WithSleepWindow(1 * time.Second)    // Time to wait before allowing a single request through in half-open state
```

## States

The circuit breaker has three states:

1. **Closed**: The circuit is closed and requests are allowed through.
2. **Open**: The circuit is open and requests are not allowed through. All requests will immediately return an error.
3. **Half-Open**: After the sleep window has elapsed, the circuit enters the half-open state, allowing a single request through to test if the dependency is healthy. If the request succeeds, the circuit will close; if it fails, the circuit will open again.

## OpenTelemetry Integration

The circuit breaker supports OpenTelemetry tracing:

```go
import (
    "go.opentelemetry.io/otel/trace"
)

// Create a tracer
tracer := otelTracer // Your OpenTelemetry tracer

// Configure the circuit breaker with the tracer
options := circuit.DefaultOptions().
    WithName("my-service").
    WithLogger(logger).
    WithOtelTracer(tracer)

cb := circuit.NewCircuitBreaker(cfg, options)
```

## Thread Safety

The circuit breaker is thread-safe and can be used concurrently from multiple goroutines.

## API Documentation


### Core Types

Description of the main types provided by the circuit.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/circuit/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the circuit.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/circuit/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/circuit/basic_usage_example.go) - Shows basic usage of the circuit
- [Advanced Configuration](../EXAMPLES/circuit/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/circuit/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Best Practice 1**: Description of best practice 1
2. **Best Practice 2**: Description of best practice 2
3. **Best Practice 3**: Description of best practice 3

## Troubleshooting

### Common Issues

#### Issue 1

Description of issue 1 and how to resolve it.

#### Issue 2

Description of issue 2 and how to resolve it.

## Related Components

- [Component 1](../circuit1/README.md) - Description of how this circuit relates to Component 1
- [Component 2](../circuit2/README.md) - Description of how this circuit relates to Component 2

## Contributing

Contributions to this circuit are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
