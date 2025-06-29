# env config_integration Example

## Overview

This example demonstrates how to integrate environment variables with a configuration structure using the ServiceLib env package. It shows how to create a structured configuration object that loads its values from environment variables, with validation and type conversion.

## Features

- **Configuration Structure**: Create a typed configuration structure to hold application settings
- **Environment Variable Loading**: Load values from environment variables into the configuration structure
- **Validation**: Validate configuration values for correctness
- **Complex Types**: Handle complex types like slices from environment variables

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### AppConfig Structure

The AppConfig structure centralizes all configuration settings in one place:

```go
type AppConfig struct {
    ServerPort  string   // Port the server will listen on (from SERVER_PORT)
    DatabaseURL string   // Connection string for the database (from DATABASE_URL)
    LogLevel    string   // Logging level: debug, info, warn, error (from LOG_LEVEL)
    APIKey      string   // API key for external service authentication (from API_KEY)
    Environment string   // Application environment: development, testing, staging, production (from APP_ENV)
    Features    []string // Enabled features as a slice of strings (from comma-separated FEATURES)
}
```

### Loading Configuration

The LoadConfig function retrieves values from environment variables and populates the AppConfig struct:

```go
func LoadConfig() AppConfig {
    // Load features from a comma-separated environment variable
    featuresStr := env.GetEnv("FEATURES", "basic,standard")
    features := strings.Split(featuresStr, ",")

    return AppConfig{
        ServerPort:  env.GetEnv("SERVER_PORT", "8080"),
        DatabaseURL: env.GetEnv("DATABASE_URL", "postgres://localhost:5432/mydb"),
        LogLevel:    env.GetEnv("LOG_LEVEL", "info"),
        APIKey:      env.GetEnv("API_KEY", ""),
        Environment: env.GetEnv("APP_ENV", "development"),
        Features:    features,
    }
}
```

### Configuration Validation

The Validate method checks the configuration values for correctness:

```go
func (c *AppConfig) Validate() []string {
    var errors []string

    // Check required values
    if c.APIKey == "" {
        errors = append(errors, "API_KEY environment variable is required")
    }

    // Validate log level
    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    if !validLogLevels[strings.ToLower(c.LogLevel)] {
        errors = append(errors, fmt.Sprintf("Invalid log level: %s. Must be one of: debug, info, warn, error", c.LogLevel))
    }

    // Validate environment
    validEnvs := map[string]bool{
        "development": true,
        "testing":     true,
        "staging":     true,
        "production":  true,
    }
    if !validEnvs[strings.ToLower(c.Environment)] {
        errors = append(errors, fmt.Sprintf("Invalid environment: %s. Must be one of: development, testing, staging, production", c.Environment))
    }

    return errors
}
```

## Expected Output

```
Environment Variables Configuration Integration Example
=====================================================
Server Port: 8080
Database URL: postgres://localhost:5432/mydb
Log Level: info
Environment: development
Features: [basic standard]
Warning: API_KEY environment variable is not set

Configuration Errors:
- API_KEY environment variable is required

Try running this example with different environment variables set:
export SERVER_PORT=9090
export DATABASE_URL="postgres://user:password@localhost:5432/mydb"
export LOG_LEVEL="debug"
export API_KEY="your-api-key"
export APP_ENV="production"
export FEATURES="basic,premium,advanced"
go run EXAMPLES/env/config_integration_example.go
```

## Related Examples

- [basic_usage](../basic_usage/README.md) - Shows basic usage of environment variables

## Related Components

- [env Package](../../../env/README.md) - The env package documentation.
- [Config Package](../../../config/README.md) - The config package for more advanced configuration.
- [Validation Package](../../../validation/README.md) - The validation package for input validation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
