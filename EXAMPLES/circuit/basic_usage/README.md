# circuit basic_usage Example

## Overview

This example demonstrates the basic usage of the ServiceLib circuit package, implementing the Circuit Breaker pattern. The Circuit Breaker pattern prevents cascading failures in distributed systems by temporarily disabling operations that are likely to fail, allowing the system to recover.

## Features

- **Circuit Breaker Creation**: Create and configure a circuit breaker with customizable parameters
- **Function Execution**: Execute functions with circuit breaking protection
- **Failure Handling**: Demonstrate how the circuit trips after multiple failures
- **Circuit State Management**: Show how to check and reset the circuit state

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Creating a Circuit Breaker

The example starts by creating a circuit breaker with default configuration:

```
// Create a circuit breaker with default configuration
cfg := circuit.DefaultConfig()
options := circuit.DefaultOptions().WithName("example")
cb := circuit.NewCircuitBreaker(cfg, options)
```

### Executing Functions with Circuit Breaking

Functions are executed with circuit breaking protection using the `Execute` method:

```
// Define a function to execute with circuit breaking
fn := func(ctx context.Context) (string, error) {
    // Your operation here
    fmt.Println("Executing operation...")
    return "success", nil
}

// Execute the function with circuit breaking
result, err := circuit.Execute(ctx, cb, "example-operation", fn)
if err != nil {
    fmt.Printf("Operation failed: %v\n", err)
} else {
    fmt.Printf("Operation succeeded with result: %s\n", result)
}
```

### Managing Circuit State

The example demonstrates how to check and reset the circuit state:

```
// Check the circuit state
fmt.Printf("Circuit state: %s\n", cb.GetState())

// Reset the circuit
fmt.Println("Resetting the circuit...")
cb.Reset()
fmt.Printf("Circuit state after reset: %s\n", cb.GetState())
```

## Expected Output

```
Executing function with circuit breaking...
Executing operation...
Operation succeeded with result: success

Executing failing function with circuit breaking...
Executing operation that will fail...
Operation failed as expected: simulated failure

Executing failing function multiple times to trip the circuit...
Executing operation that will fail...
Attempt 1: Operation failed: simulated failure
Executing operation that will fail...
Attempt 2: Operation failed: simulated failure
Executing operation that will fail...
Attempt 3: Operation failed: simulated failure
Executing operation that will fail...
Attempt 4: Operation failed: simulated failure
Executing operation that will fail...
Attempt 5: Operation failed: simulated failure
Attempt 6: Operation failed: circuit breaker is open
Attempt 7: Operation failed: circuit breaker is open
Attempt 8: Operation failed: circuit breaker is open
Attempt 9: Operation failed: circuit breaker is open
Attempt 10: Operation failed: circuit breaker is open

Circuit state: open

Trying to execute a successful function after the circuit is open...
Operation failed because circuit is open: circuit breaker is open

Resetting the circuit...
Circuit state after reset: closed

Executing a successful function after reset...
Executing operation...
Operation succeeded with result: success
```

## Related Examples

- [custom_configuration](../custom_configuration/README.md) - Demonstrates how to create a circuit breaker with custom configuration options, including timeout, error threshold, volume threshold, and fallback functionality

## Related Components

- [circuit Package](../../../circuit/README.md) - The circuit package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
