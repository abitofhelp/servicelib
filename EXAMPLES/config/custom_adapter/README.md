# config custom_adapter Example

## Overview

This example demonstrates how to create and use a custom configuration adapter with the ServiceLib config package. It shows how to implement a configuration provider that reads from environment variables and integrates with the existing configuration system.

## Features

- **Custom Environment Adapter**: Create a configuration adapter that reads from environment variables
- **Interface Implementation**: Implement the required interfaces for application and database configuration
- **Extended Functionality**: Add custom methods for retrieving typed configuration values (int, bool)
- **Default Values**: Provide fallback values when environment variables are not set

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Custom Environment Adapter

This example defines a custom configuration adapter that reads from environment variables with a specified prefix:

```
// EnvConfig is a configuration provider that reads from environment variables
type EnvConfig struct {
    prefix string
}

// NewEnvConfig creates a new EnvConfig with the given prefix
func NewEnvConfig(prefix string) *EnvConfig {
    return &EnvConfig{
        prefix: prefix,
    }
}

// Helper method to get environment variables with a prefix
func (c *EnvConfig) getEnv(key, defaultValue string) string {
    fullKey := c.prefix + "_" + key
    if value, exists := os.LookupEnv(fullKey); exists {
        return value
    }
    return defaultValue
}
```

### Interface Implementation

The custom adapter implements the required interfaces for application and database configuration:

```
// GetAppVersion implements the AppConfigProvider interface
func (c *EnvConfig) GetAppVersion() string {
    return c.getEnv("VERSION", "1.0.0")
}

// GetAppName implements the AppConfigProvider interface
func (c *EnvConfig) GetAppName() string {
    return c.getEnv("NAME", "DefaultApp")
}

// GetDatabaseType implements the DatabaseConfigProvider interface
func (c *EnvConfig) GetDatabaseType() string {
    return c.getEnv("DB_TYPE", "postgres")
}

// GetDatabaseConnectionString implements the DatabaseConfigProvider interface
func (c *EnvConfig) GetDatabaseConnectionString(dbType string) string {
    host := c.getEnv("DB_HOST", "localhost")
    port := c.getEnv("DB_PORT", "5432")
    user := c.getEnv("DB_USER", "postgres")
    pass := c.getEnv("DB_PASS", "postgres")
    name := c.getEnv("DB_NAME", "postgres")

    // Build connection string based on database type
    // ...
}
```

### Extended Functionality

The custom adapter adds methods for retrieving typed configuration values:

```
// Additional methods specific to EnvConfig
func (c *EnvConfig) GetInt(key string, defaultValue int) int {
    strValue := c.getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }

    intValue, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }

    return intValue
}

func (c *EnvConfig) GetBool(key string, defaultValue bool) bool {
    strValue := strings.ToLower(c.getEnv(key, ""))
    if strValue == "" {
        return defaultValue
    }

    return strValue == "true" || strValue == "yes" || strValue == "1"
}
```

## Expected Output

```
=== Application Configuration (from Environment) ===
Version: 3.0.0
Name: EnvApp
Environment: production

=== Database Configuration (from Environment) ===
Type: mysql
Connection String: mysql://admin:secret@db.example.com:3306/myapp
Database Name: myapp

=== Additional Settings (from Environment) ===
Debug Mode: true
Max Connections: 50
```

## Related Examples

- [app_config](../app_config/README.md) - Shows how to configure application-specific settings
- [basic_usage](../basic_usage/README.md) - Demonstrates basic configuration operations
- [database_config](../database_config/README.md) - Illustrates database-specific configuration options

## Related Components

- [config Package](../../../config/README.md) - The config package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
