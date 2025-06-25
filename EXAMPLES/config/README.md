# Config Package Examples

This directory contains examples demonstrating how to use the `config` package, which provides generic configuration interfaces and adapters for Go applications. The package offers a flexible way to manage application and database configuration.

## Examples

### 1. Basic Usage Example

[basic_usage/main.go](basic_usage/main.go)

Demonstrates basic usage of the GenericConfigAdapter with a custom configuration struct.

Key concepts:
- Creating a custom configuration struct
- Implementing the AppConfigProvider and DatabaseConfigProvider interfaces
- Creating a GenericConfigAdapter
- Accessing application and database configuration

### 2. App Config Example

[app_config/main.go](app_config/main.go)

Shows how to implement and use the AppConfig interface.

Key concepts:
- Creating a custom AppSettings struct
- Implementing the AppConfigProvider interface
- Adding application-specific configuration methods
- Accessing configuration through both the adapter and the original settings object

### 3. Database Config Example

[database_config/main.go](database_config/main.go)

Demonstrates how to implement and use the DatabaseConfig interface.

Key concepts:
- Creating a custom DatabaseSettings struct
- Implementing the DatabaseConfigProvider interface
- Generating connection strings for different database types
- Adding database-specific configuration methods
- Accessing database configuration through both the adapter and the original settings object

### 4. Custom Adapter Example

[custom_adapter/main.go](custom_adapter/main.go)

Shows how to create a custom adapter that reads configuration from environment variables.

Key concepts:
- Creating a custom EnvConfig adapter
- Reading configuration from environment variables
- Implementing type conversion methods (GetInt, GetBool)
- Using the adapter with the GenericConfigAdapter
- Accessing configuration through both adapters

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run basic_usage/main.go
```

## Additional Resources

For more information about the config package, see the [config package documentation](../../config/README.md).
