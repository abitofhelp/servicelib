# Signal Module

The Signal Module provides utilities for handling OS signals and implementing graceful shutdown in Go applications. It helps applications respond to termination signals and execute cleanup operations before exiting.

## Features

- **Signal Handling**: Captures OS termination signals (SIGINT, SIGTERM, SIGHUP, SIGQUIT)
- **Callback Registration**: Register multiple shutdown callbacks to be executed during shutdown
- **Concurrent Execution**: Execute shutdown callbacks concurrently
- **Timeout Management**: Apply timeouts to prevent hanging during shutdown
- **Context Cancellation**: Propagate shutdown events via context cancellation
- **Multiple Signal Handling**: Handle multiple signals with forced exit on second signal
- **Logging Integration**: Comprehensive logging of shutdown events

## Installation

```bash
go get github.com/abitofhelp/servicelib/signal
```

## Quick Start

See the [Basic Signal Handling example](../examples/signal/basic_signal_handling_example.go) for a complete, runnable example of how to use the Signal module.

## API Documentation

### Basic Signal Handling

The `WaitForShutdown` function blocks until a shutdown signal is received and returns a context that will be canceled when a shutdown signal is received.

#### Basic Usage

See the [Basic Signal Handling example](../examples/signal/basic_signal_handling_example.go) for a complete, runnable example of how to implement basic signal handling.

### Shutdown Callbacks

The `SetupSignalHandler` function sets up a signal handler for graceful shutdown and returns a context that will be canceled when a shutdown signal is received, along with a `GracefulShutdown` instance for registering callbacks.

#### Using Shutdown Callbacks

See the [Shutdown Callbacks example](../examples/signal/shutdown_callbacks_example.go) for a complete, runnable example of how to use shutdown callbacks.

### Manual Cancellation

The `HandleShutdown` method of the `GracefulShutdown` struct handles graceful shutdown and returns a context that will be canceled when a shutdown signal is received, along with a cancel function that can be called to cancel the context manually.

#### Manual Cancellation Example

See the [Manual Cancellation example](../examples/signal/manual_cancellation_example.go) for a complete, runnable example of how to implement manual cancellation.

## Best Practices

1. **Callback Ordering**: Register callbacks in the reverse order of resource creation to ensure proper dependency handling.

2. **Timeouts**: Set appropriate timeouts for shutdown operations to prevent hanging.

3. **Health Checks**: Update health check status early in the shutdown process to prevent new requests.

4. **Concurrent Execution**: Use concurrent execution for independent shutdown operations, but be careful with dependencies.

5. **Context Propagation**: Pass the shutdown context to all operations that need to be aware of shutdown.

6. **Error Handling**: Log errors during shutdown but continue with other shutdown operations.

7. **Signal Handling**: Be prepared to handle multiple termination signals, including forced termination.

8. **Resource Cleanup**: Ensure all resources are properly released during shutdown.

9. **Graceful Termination**: Allow in-flight operations to complete before shutting down services.

## Comparison with shutdown Package

The `signal` package is similar to the `shutdown` package but offers a different approach:

- **Callback-based**: The `signal` package uses a callback registration model
- **Concurrent Execution**: Executes shutdown callbacks concurrently
- **Object-Oriented**: Uses a `GracefulShutdown` struct with methods
- **More Flexible**: Allows for more customization of shutdown behavior

Choose the `signal` package when you need more control over the shutdown process or when you have multiple independent resources to shut down concurrently.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
