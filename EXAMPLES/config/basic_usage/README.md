# config basic_usage Example

## Overview

This example demonstrates the basic usage of the ServiceLib config package, showing how to create a custom configuration structure, use the GenericConfigAdapter, and access application and database configuration values through a unified interface.

## Features

- **Custom Configuration Structure**: Create a custom configuration struct that implements required interfaces
- **Configuration Adapter**: Use the GenericConfigAdapter to access configuration values
- **Application and Database Settings**: Access application and database configuration through a unified interface

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Custom Configuration Structure

This example defines a custom configuration struct that implements the required interfaces for application and database configuration:

```
// MyConfig is a simple configuration struct that implements the required interfaces
type MyConfig struct {
    Version     string
    Name        string
    Environment string
    DBType      string
    DBConnStr   map[string]string
}

// GetAppVersion implements the AppConfigProvider interface
func (c *MyConfig) GetAppVersion() string {
    return c.Version
}

// Additional interface implementation methods...
```

### Configuration Adapter

The example uses the GenericConfigAdapter to provide a standardized way to access configuration values:

```
// Create a config adapter
adapter := config.NewGenericConfigAdapter(myConfig).
    WithAppName("MyApplication").
    WithAppEnvironment("production").
    WithDatabaseName("my_database")
```

### Accessing Configuration Values

The example demonstrates how to retrieve application and database configuration values through the adapter:

```
// Get the app configuration
appConfig := adapter.GetApp()
fmt.Println("Application Version:", appConfig.GetVersion())
fmt.Println("Application Name:", appConfig.GetName())
fmt.Println("Application Environment:", appConfig.GetEnvironment())

// Get the database configuration
dbConfig := adapter.GetDatabase()
fmt.Println("Database Type:", dbConfig.GetType())
fmt.Println("Database Connection String:", dbConfig.GetConnectionString())
```

## Expected Output

```
Application Version: 1.0.0
Application Name: MyApp
Application Environment: development
Database Type: postgres
Database Connection String: postgres://user:password@localhost:5432/mydb
Database Name: my_database
Users Collection: users
```

## Related Examples

- [app_config](../app_config/README.md) - Shows how to configure application-specific settings
- [custom_adapter](../custom_adapter/README.md) - Demonstrates creating custom configuration adapters
- [database_config](../database_config/README.md) - Illustrates database-specific configuration options

## Related Components

- [config Package](../../../config/README.md) - The config package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
