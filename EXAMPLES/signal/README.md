# Signal Package Examples

This directory contains examples demonstrating how to use the `signal` package, which provides utilities for handling OS signals in Go applications. The package makes it easier to implement graceful shutdown and other signal-based behaviors in response to signals like SIGINT (Ctrl+C) and SIGTERM.

## Examples

### 1. Basic Signal Handling Example

[basic_signal_handling_example.go](basic_signal_handling_example.go)

Demonstrates the basic usage of signal handling for graceful shutdown.

Key concepts:
- Using WaitForShutdown to wait for termination signals
- Setting up an HTTP server
- Implementing graceful shutdown when signals are received
- Setting timeouts for the shutdown process
- Logging shutdown events

### 2. Manual Cancellation Example

[manual_cancellation_example.go](manual_cancellation_example.go)

Shows how to combine signal handling with manual cancellation.

Key concepts:
- Creating a cancellable context
- Handling both OS signals and programmatic cancellation
- Coordinating multiple cancellation sources
- Implementing custom cancellation logic
- Propagating cancellation through the application

### 3. Shutdown Callbacks Example

[shutdown_callbacks_example.go](shutdown_callbacks_example.go)

Demonstrates how to register and execute callbacks when shutdown signals are received.

Key concepts:
- Registering shutdown callbacks
- Executing callbacks in a specific order
- Handling errors in shutdown callbacks
- Managing resources during shutdown
- Implementing complex shutdown sequences

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/signal/basic_signal_handling_example.go
```

## Additional Resources

For more information about the signal package, see the [signal package documentation](../../signal/README.md).