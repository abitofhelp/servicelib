# Environment Variables Package Examples

This directory contains examples demonstrating how to use the `env` package, which provides utilities for working with environment variables in Go applications. It simplifies the process of retrieving environment variables with fallback values.

## Examples

### 1. Basic Usage Example

[basic_usage_example.go](basic_usage_example.go)

Demonstrates how to retrieve environment variables with fallback values.

Key concepts:
- Using `GetEnv` to retrieve environment variables with default values
- Handling missing environment variables
- Setting up fallback values for configuration

### 2. Configuration Integration Example

[config_integration_example.go](config_integration_example.go)

Shows how to integrate environment variables with a configuration structure.

Key concepts:
- Creating a configuration struct
- Loading environment variables into the configuration
- Providing default values for non-critical settings
- Validating configuration after loading

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/env/basic_usage_example.go
```

You can also set environment variables before running the examples:

```bash
# Unix/Linux/macOS
export PORT=9090
export DATABASE_URL="postgres://user:password@localhost:5432/mydb"
go run examples/env/basic_usage_example.go

# Windows (Command Prompt)
set PORT=9090
set DATABASE_URL=postgres://user:password@localhost:5432/mydb
go run examples/env/basic_usage_example.go
```

## Best Practices

1. **Sensitive Information**: Use environment variables for sensitive information like API keys, passwords, and tokens.

2. **Default Values**: Provide sensible default values for non-critical environment variables.

3. **Documentation**: Document all environment variables used by your application, including their purpose and default values.

4. **Validation**: Validate environment variables after retrieving them to ensure they meet your requirements.

5. **Configuration Structure**: Use a structured approach to loading environment variables into a configuration object.

## Additional Resources

For more information about the env package, see the [env package documentation](../../env/README.md).