# Environment Variables Package

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

See the [Basic Usage Example](../examples/env/basic_usage_example.go) for a complete, runnable example of how to use the env package.

## Configuration

See the [Configuration Integration Example](../examples/env/config_integration_example.go) for a complete, runnable example of how to integrate environment variables with a configuration structure.

## Best Practices

1. **Sensitive Information**: Use environment variables for sensitive information like API keys, passwords, and tokens.

2. **Default Values**: Provide sensible default values for non-critical environment variables.

3. **Documentation**: Document all environment variables used by your application, including their purpose and default values.

4. **Validation**: Validate environment variables after retrieving them to ensure they meet your requirements.

5. **Configuration Structure**: Use a structured approach to loading environment variables into a configuration object.

6. **Error Handling**: Implement proper error handling for missing critical environment variables.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
