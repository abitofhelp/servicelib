# Shutdown Module
The Shutdown Module provides functionality for graceful application shutdown in Go applications. It helps ensure that applications terminate cleanly, allowing resources to be properly released and pending operations to complete.


## Overview

Brief description of the shutdown and its purpose in the ServiceLib library.

## Features

- **Signal Handling**: Captures OS termination signals (SIGINT, SIGTERM, SIGHUP)
- **Context Cancellation**: Supports shutdown via context cancellation
- **Timeout Management**: Applies timeouts to prevent hanging during shutdown
- **Multiple Signal Handling**: Forces exit if a second signal is received during shutdown
- **Error Propagation**: Returns errors from shutdown operations
- **Logging Integration**: Comprehensive logging of shutdown events


## Installation

```bash
go get github.com/abitofhelp/servicelib/shutdown
```


## Quick Start

See the [Basic Usage example](../EXAMPLES/shutdown/basic_usage_example.go) for a complete, runnable example of how to use the Shutdown module.


## Configuration

See the [Configuration example](../EXAMPLES/shutdown/configuration_example.go) for a complete, runnable example of how to configure the shutdown.

## API Documentation

### Graceful Shutdown

The `GracefulShutdown` function waits for termination signals and calls the provided shutdown function.

#### Basic Usage

See the [Basic Usage example](../EXAMPLES/shutdown/basic_usage_example.go) for a complete, runnable example of how to implement graceful shutdown.

### Programmatic Shutdown

The `SetupGracefulShutdown` function sets up a goroutine that will handle graceful shutdown, allowing for both signal-based and programmatic shutdown initiation.

#### Programmatic Shutdown Initiation

See the [Programmatic Shutdown example](../EXAMPLES/shutdown/programmatic_shutdown_example.go) for a complete, runnable example of how to implement programmatic shutdown.

### Multiple Resource Shutdown

When shutting down an application with multiple resources, it's important to shut them down in the correct order.

#### Shutting Down Multiple Resources

See the [Multiple Resource example](../EXAMPLES/shutdown/multiple_resource_example.go) for a complete, runnable example of how to shut down multiple resources in the correct order.

### Shutdown Function Signature

The shutdown function passed to the shutdown module should have the following signature:

```go
// Example of a shutdown function
package example

// ShutdownFunc is a function that performs shutdown operations
type ShutdownFunc func() error
```


### Core Types

Description of the main types provided by the shutdown.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/shutdown/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the shutdown.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/shutdown/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/shutdown/basic_usage_example.go) - Shows basic usage of the shutdown
- [Advanced Configuration](../EXAMPLES/shutdown/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/shutdown/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Resource Ordering**: Close resources in the reverse order they were created.

2. **Timeouts**: Set appropriate timeouts for shutdown operations to prevent hanging.

3. **Logging**: Log the beginning and completion of each shutdown step.

4. **Error Handling**: Properly handle and log errors during shutdown, but continue shutting down other resources.

5. **Signal Handling**: Be prepared to handle multiple termination signals.

6. **Context Usage**: Use contexts with timeouts for shutdown operations.

7. **Graceful Termination**: Allow in-flight operations to complete before shutting down.

8. **Health Checks**: Update health check status during shutdown to prevent new requests.

9. **Dependency Management**: Consider dependencies between resources when ordering shutdown.


## Troubleshooting

### Common Issues

#### Issue 1

Description of issue 1 and how to resolve it.

#### Issue 2

Description of issue 2 and how to resolve it.

## Related Components

- [Component 1](../shutdown1/README.md) - Description of how this shutdown relates to Component 1
- [Component 2](../shutdown2/README.md) - Description of how this shutdown relates to Component 2

## Contributing

Contributions to this shutdown are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
