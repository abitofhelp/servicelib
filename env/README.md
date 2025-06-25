# Environment Variables Package

## Overview

The `env` package provides utilities for working with environment variables in Go applications. It simplifies the process of retrieving environment variables with fallback values.


## Features

- **Environment Variable Retrieval**: Easy retrieval of environment variables
- **Default Values**: Support for fallback values when environment variables are not set
- **Simple API**: Clean and straightforward API for environment variable handling


## Installation

```bash
go get github.com/abitofhelp/servicelib/env
```


## Quick Start

See the [Basic Usage Example](../EXAMPLES/env/basic_usage/main.go) for a complete, runnable example of how to use the env package.


## Configuration

See the [Configuration Integration Example](../EXAMPLES/env/config_integration/main.go) for a complete, runnable example of how to integrate environment variables with a configuration structure.


## API Documentation


### Core Types

The env package does not define any custom types. It provides utility functions for working with environment variables.


### Key Methods

#### GetEnv

`GetEnv` retrieves the value of the environment variable named by the key. If the variable is not present, the fallback value is returned.

```go
package env

func GetEnv(key, fallback string) string
```

Parameters:
- `key`: The name of the environment variable
- `fallback`: The value to return if the environment variable is not set

Returns:
- `string`: The value of the environment variable or the fallback value

See the [Basic Usage Example](../EXAMPLES/env/basic_usage/main.go) for a complete, runnable example of how to use the GetEnv method.


## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Usage](../EXAMPLES/env/basic_usage/main.go) - Shows basic usage of the env package
- [Configuration Integration](../EXAMPLES/env/config_integration/main.go) - Shows how to integrate environment variables with a configuration structure


## Best Practices

1. **Sensitive Information**: Use environment variables for sensitive information like API keys, passwords, and tokens.

2. **Default Values**: Provide sensible default values for non-critical environment variables.

3. **Documentation**: Document all environment variables used by your application, including their purpose and default values.

4. **Validation**: Validate environment variables after retrieving them to ensure they meet your requirements.

5. **Configuration Structure**: Use a structured approach to loading environment variables into a configuration object.

6. **Error Handling**: Implement proper error handling for missing critical environment variables.


## Troubleshooting

### Common Issues

#### Missing Environment Variables

**Issue**: Critical environment variables are missing, causing the application to fail.

**Solution**: Implement proper validation of environment variables at application startup and provide clear error messages when required variables are missing.

#### Incorrect Environment Variable Names

**Issue**: Environment variables are not being found due to incorrect names.

**Solution**: Double-check the names of environment variables and ensure they match between the application code and the environment where the application is running.


## Related Components

- [Configuration](../config/README.md) - The configuration component can use environment variables as a source for configuration values.
- [Logging](../logging/README.md) - The logging component can use environment variables to configure log levels and output destinations.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
