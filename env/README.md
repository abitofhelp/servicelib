# Environment

## Overview

The Environment component provides utilities for working with environment variables in Go applications. It simplifies the process of retrieving environment variables with fallback values, making configuration via environment variables more robust and convenient.

## Features

- **Simple API**: Easy-to-use functions for retrieving environment variables
- **Fallback Values**: Automatic fallback to default values when environment variables are not set
- **Zero Dependencies**: No external dependencies beyond the Go standard library
- **Type Safety**: Clear function signatures for type-safe environment variable handling
- **Testability**: Easy to mock for testing purposes

## Installation

```bash
go get github.com/abitofhelp/servicelib/env
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/env/basic_usage/README.md) for a complete, runnable example of how to use the environment component.

## API Documentation

### Key Methods

#### GetEnv

Retrieves the value of the environment variable named by the key. If the variable is not present, the fallback value is returned.

```go
func GetEnv(key, fallback string) string
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/env/basic_usage/README.md) - Shows basic usage of environment variables with fallbacks
- [Config Integration](../EXAMPLES/env/config_integration/README.md) - Shows how to integrate environment variables with a configuration structure

## Best Practices

1. **Use Descriptive Keys**: Choose clear, descriptive names for environment variables
2. **Provide Sensible Defaults**: Always provide reasonable fallback values for non-critical settings
3. **Validate Values**: Validate environment variable values before using them in your application
4. **Document Required Variables**: Clearly document which environment variables are required and which are optional
5. **Use Prefixes**: Use prefixes for environment variables to avoid naming conflicts (e.g., APP_PORT instead of PORT)

## Troubleshooting

### Common Issues

#### Environment Variables Not Being Read

If your environment variables are not being read:
- Verify that the variables are actually set in the environment
- Check for typos in the variable names
- Ensure the variables are exported properly in your shell
- Remember that environment variables are case-sensitive

#### Unexpected Values

If you're getting unexpected values:
- Check if the variable is being overridden elsewhere
- Verify that the variable is set at the correct scope (system vs. user vs. process)
- Ensure there are no leading or trailing whitespaces in the variable value

## Related Components

- [Config](../config/README.md) - Configuration management that can use environment variables
- [Logging](../logging/README.md) - Logging component that can be configured via environment variables

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
