# Configuration Package

The `config` package provides a flexible configuration system that supports multiple sources and formats for Go applications. It allows you to manage application configuration in a hierarchical manner with support for various data sources.

## Features

- **Multiple Sources**:
  - YAML and JSON files
  - Environment variables
  - Command-line flags
  - In-memory values

- **Hierarchical Configuration**: Access nested values using dot notation
- **Default Values**: Specify fallback values when configuration is missing
- **Type Conversion**: Automatic conversion to the appropriate type
- **Configuration Reloading**: Watch for changes and reload configuration
- **Validation**: Validate configuration against schemas
- **Adapters**: Easily create custom adapters for different configuration sources

## Installation

```bash
go get github.com/abitofhelp/servicelib/config
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/abitofhelp/servicelib/config"
)

func main() {
    // Create a new configuration from files
    cfg, err := config.New("config.yaml", "env.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Get a string value
    apiKey := cfg.GetString("api.key")
    fmt.Println("API Key:", apiKey)

    // Get an int value with default
    port := cfg.GetInt("server.port", 8080)
    fmt.Println("Server Port:", port)

    // Get a nested value
    dbURL := cfg.GetString("database.url")
    fmt.Println("Database URL:", dbURL)

    // Get a boolean value
    debug := cfg.GetBool("logging.debug", false)
    fmt.Println("Debug Mode:", debug)

    // Get a duration value
    timeout := cfg.GetDuration("server.timeout", "30s")
    fmt.Println("Server Timeout:", timeout)
}
```

### Binding to Structs

```go
package main

import (
    "fmt"
    "log"
    "github.com/abitofhelp/servicelib/config"
)

// Example configuration struct
type AppConfig struct {
    Server struct {
        Port    int    `yaml:"port"`
        Host    string `yaml:"host"`
        Timeout int    `yaml:"timeout"`
    } `yaml:"server"`
    Database struct {
        URL      string `yaml:"url"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Pool     int    `yaml:"pool"`
    } `yaml:"database"`
    API struct {
        Key      string `yaml:"key"`
        Endpoint string `yaml:"endpoint"`
        Version  string `yaml:"version"`
    } `yaml:"api"`
    Logging struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
        Path   string `yaml:"path"`
    } `yaml:"logging"`
}

func main() {
    // Create a new configuration
    cfg, err := config.New("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Bind configuration to a struct
    var appConfig AppConfig
    if err := cfg.Unmarshal(&appConfig); err != nil {
        log.Fatalf("Failed to unmarshal configuration: %v", err)
    }

    fmt.Printf("Server Configuration: %+v\n", appConfig.Server)
    fmt.Printf("Database Configuration: %+v\n", appConfig.Database)
}
```

### Watching for Changes

```go
package main

import (
    "log"
    "github.com/abitofhelp/servicelib/config"
)

func main() {
    // Create a new configuration
    cfg, err := config.New("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Watch for configuration changes
    cfg.Watch(func() {
        // Reload configuration when changes are detected
        var appConfig AppConfig
        if err := cfg.Unmarshal(&appConfig); err != nil {
            log.Printf("Failed to reload configuration: %v", err)
            return
        }
        log.Println("Configuration reloaded successfully")
        
        // Apply the new configuration
        applyConfiguration(appConfig)
    })
    
    // Keep the application running
    select {}
}

func applyConfiguration(config AppConfig) {
    // Apply the configuration to your application
    // For example, update database connections, change log levels, etc.
}
```

### Environment Variables

The config package can also load configuration from environment variables:

```go
// Create a configuration that includes environment variables
cfg, err := config.NewWithOptions(config.Options{
    Files: []string{"config.yaml"},
    EnvPrefix: "APP_",
    EnvReplacer: "_",
})

// This will look for environment variables like:
// APP_SERVER_PORT, APP_DATABASE_URL, etc.
```

## Configuration File Example

Here's an example of a YAML configuration file:

```yaml
server:
  port: 8080
  host: "localhost"
  timeout: 30

database:
  url: "postgres://user:password@localhost:5432/mydb"
  username: "user"
  password: "password"
  pool: 10

api:
  key: "your-api-key"
  endpoint: "https://api.example.com"
  version: "v1"

logging:
  level: "info"
  format: "json"
  path: "/var/log/app.log"
```

## Best Practices

1. **Use Environment Variables for Secrets**: Never store sensitive information like API keys or passwords in configuration files. Use environment variables instead.

2. **Default Values**: Always provide sensible default values for configuration that might be missing.

3. **Validation**: Validate configuration values to ensure they meet your requirements.

4. **Configuration Hierarchy**: Use a hierarchical approach to organize your configuration.

5. **Documentation**: Document all configuration options, including their purpose, type, and default values.

## License

This project is licensed under the MIT License - see the LICENSE file for details.