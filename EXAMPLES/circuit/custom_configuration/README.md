# circuit custom_configuration Example

## Overview

This example demonstrates how to create and configure a circuit breaker with custom parameters using the ServiceLib circuit package. It shows how to fine-tune the circuit breaker's behavior to meet specific application requirements.

## Features

- **Custom Configuration**: Configure timeout, error threshold, volume threshold, and sleep window
- **Logger Integration**: Integrate with the logging system for better visibility
- **Fallback Functionality**: Implement fallback mechanisms for when the circuit is open
- **Parameter Tuning**: Demonstrate individual parameter configuration for specific needs

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Custom Configuration

The example demonstrates how to create a circuit breaker with custom configuration parameters:

```
// Create a circuit breaker with custom configuration
cfg := circuit.DefaultConfig().
    WithEnabled(true).                // Explicitly enable the circuit breaker
    WithTimeout(2 * time.Second).     // Set timeout to 2 seconds
    WithMaxConcurrent(50).            // Set max concurrent requests to 50
    WithErrorThreshold(0.3).          // Trip the circuit at 30% error rate
    WithVolumeThreshold(5).           // Require at least 5 requests before tripping
    WithSleepWindow(5 * time.Second)  // Wait 5 seconds before trying again

options := circuit.DefaultOptions().
    WithName("custom-example").       // Set a custom name
    WithLogger(contextLogger)         // Use a custom logger

cb := circuit.NewCircuitBreaker(cfg, options)
```

### Fallback Functionality

The example shows how to implement fallback functionality for when the circuit is open:

```
// Define a fallback function
fallbackFn := func(ctx context.Context, err error) (string, error) {
    fmt.Printf("Fallback called with error: %v\n", err)
    return "fallback result", nil
}

// Execute with fallback
result, err = circuit.ExecuteWithFallback(ctx, cb, "failing-operation", failFn, fallbackFn)
if err != nil {
    fmt.Printf("Operation with fallback failed: %v\n", err)
} else {
    fmt.Printf("Operation with fallback succeeded with result: %s\n", result)
}
```

### Individual Parameter Configuration

The example demonstrates how to configure individual parameters:

```
// Start with default config
cfg = circuit.DefaultConfig()

// Configure timeout only
cfg = cfg.WithTimeout(3 * time.Second)
fmt.Printf("Timeout: %v\n", cfg.Timeout)

// Configure error threshold only
cfg = cfg.WithErrorThreshold(0.2)
fmt.Printf("Error threshold: %.1f\n", cfg.ErrorThreshold)

// Configure volume threshold only
cfg = cfg.WithVolumeThreshold(10)
fmt.Printf("Volume threshold: %d\n", cfg.VolumeThreshold)

// Configure sleep window only
cfg = cfg.WithSleepWindow(10 * time.Second)
fmt.Printf("Sleep window: %v\n", cfg.SleepWindow)
```

## Expected Output

```
Executing function with custom circuit breaking...
Executing operation with custom configuration...
Operation succeeded with result: success

Executing failing function with circuit breaking...
With error threshold of 0.3 and volume threshold of 5,
the circuit should trip after 2 failures (30% of 5 requests).
Executing operation that will fail...
Attempt 1: Operation failed: simulated failure
Executing operation that will fail...
Attempt 2: Operation failed: simulated failure
Executing operation that will fail...
Attempt 3: Operation failed: circuit breaker is open
Attempt 4: Operation failed: circuit breaker is open
Attempt 5: Operation failed: circuit breaker is open

Circuit state: open

Demonstrating fallback functionality...
Fallback called with error: circuit breaker is open
Operation with fallback succeeded with result: fallback result

Demonstrating individual parameter configuration:
Timeout: 3s
Error threshold: 0.2
Volume threshold: 10
Sleep window: 10s
```

## Related Examples

- [basic_usage](../basic_usage/README.md) - Demonstrates the basic usage of circuit breakers, including creating a circuit breaker with default configuration, executing functions, and handling circuit state

## Related Components

- [circuit Package](../../../circuit/README.md) - The circuit package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
