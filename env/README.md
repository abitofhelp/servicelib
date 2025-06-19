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

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/env"
)

func main() {
    // Get an environment variable with a fallback value
    port := env.GetEnv("PORT", "8080")
    fmt.Printf("Server will run on port: %s\n", port)
    
    // Get a database URL with a fallback
    dbURL := env.GetEnv("DATABASE_URL", "postgres://localhost:5432/mydb")
    fmt.Printf("Database URL: %s\n", dbURL)
    
    // Get API keys (sensitive information)
    apiKey := env.GetEnv("API_KEY", "")
    if apiKey == "" {
        fmt.Println("Warning: API_KEY environment variable is not set")
    }
}
```

### Integration with Configuration

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/env"
)

// AppConfig holds the application configuration
type AppConfig struct {
    ServerPort  string
    DatabaseURL string
    LogLevel    string
    APIKey      string
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig() AppConfig {
    return AppConfig{
        ServerPort:  env.GetEnv("SERVER_PORT", "8080"),
        DatabaseURL: env.GetEnv("DATABASE_URL", "postgres://localhost:5432/mydb"),
        LogLevel:    env.GetEnv("LOG_LEVEL", "info"),
        APIKey:      env.GetEnv("API_KEY", ""),
    }
}

func main() {
    config := LoadConfig()
    
    fmt.Printf("Server Port: %s\n", config.ServerPort)
    fmt.Printf("Database URL: %s\n", config.DatabaseURL)
    fmt.Printf("Log Level: %s\n", config.LogLevel)
    
    if config.APIKey == "" {
        fmt.Println("Warning: API_KEY environment variable is not set")
    }
}
```

## Best Practices

1. **Sensitive Information**: Use environment variables for sensitive information like API keys, passwords, and tokens.

2. **Default Values**: Provide sensible default values for non-critical environment variables.

3. **Documentation**: Document all environment variables used by your application, including their purpose and default values.

4. **Validation**: Validate environment variables after retrieving them to ensure they meet your requirements.

5. **Configuration Structure**: Use a structured approach to loading environment variables into a configuration object.

6. **Error Handling**: Implement proper error handling for missing critical environment variables.

## License

This project is licensed under the MIT License - see the LICENSE file for details.