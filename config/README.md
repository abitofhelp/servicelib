# Configuration Package

## Overview

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

## Quick Start

See the [Basic Usage example](../EXAMPLES/config/basic_usage_example.go) for a complete, runnable example of how to use the config package.

## Configuration

See the [Configuration example](../examples/config/app_config_example.go) for a complete, runnable example of how to configure the config package.

## API Documentation

### Core Types

#### Config

The `Config` interface provides methods for accessing configuration values.

See the [Basic Usage example](../examples/config/basic_usage_example.go) for a complete, runnable example of how to use the Config interface.

#### Options

The `Options` struct provides configuration options for creating a new config instance.

See the [Custom Adapter example](../examples/config/custom_adapter_example.go) for a complete, runnable example of how to use the Options struct.

### Key Methods

#### New

The `New` function creates a new configuration instance from the specified files.

```go
cfg, err := config.New("config.yaml", "env.yaml")
```

See the [Basic Usage example](../examples/config/basic_usage_example.go) for a complete, runnable example.

#### NewWithOptions

The `NewWithOptions` function creates a new configuration instance with the specified options.

```go
cfg, err := config.NewWithOptions(config.Options{
    Files: []string{"config.yaml"},
    EnvPrefix: "APP_",
    EnvReplacer: "_",
})
```

See the [Custom Adapter example](../examples/config/custom_adapter_example.go) for a complete, runnable example.

#### Get Methods

The `Get*` methods retrieve configuration values of different types.

```go
// Get a string value
apiKey := cfg.GetString("api.key")

// Get an int value with default
port := cfg.GetInt("server.port", 8080)

// Get a boolean value
debug := cfg.GetBool("logging.debug", false)

// Get a duration value
timeout := cfg.GetDuration("server.timeout", "30s")
```

See the [Basic Usage example](../examples/config/basic_usage_example.go) for a complete, runnable example.

#### Unmarshal

The `Unmarshal` method binds configuration values to a struct.

```go
var appConfig AppConfig
if err := cfg.Unmarshal(&appConfig); err != nil {
    log.Fatalf("Failed to unmarshal configuration: %v", err)
}
```

See the [App Config example](../examples/config/app_config_example.go) for a complete, runnable example.

#### Watch

The `Watch` method watches for configuration changes and calls the specified function when changes are detected.

```go
cfg.Watch(func() {
    // Reload configuration when changes are detected
})
```

See the [Basic Usage example](../examples/config/basic_usage_example.go) for a complete, runnable example.

## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Usage](../examples/config/basic_usage_example.go) - Shows basic usage of the config package
- [App Config](../examples/config/app_config_example.go) - Shows how to bind configuration to a struct
- [Database Config](../examples/config/database_config_example.go) - Shows how to configure database connections
- [Custom Adapter](../examples/config/custom_adapter_example.go) - Shows how to create a custom configuration adapter

## Best Practices

1. **Use Environment Variables for Secrets**: Never store sensitive information like API keys or passwords in configuration files. Use environment variables instead.

2. **Default Values**: Always provide sensible default values for configuration that might be missing.

3. **Validation**: Validate configuration values to ensure they meet your requirements.

4. **Configuration Hierarchy**: Use a hierarchical approach to organize your configuration.

5. **Documentation**: Document all configuration options, including their purpose, type, and default values.

## Troubleshooting

### Common Issues

#### Missing Configuration Files

**Issue**: Configuration files are not found at the specified path.

**Solution**: Ensure that the configuration files exist at the specified path. Use absolute paths or paths relative to the working directory.

#### Type Conversion Errors

**Issue**: Type conversion fails when retrieving configuration values.

**Solution**: Ensure that the configuration values are of the expected type. Use the appropriate Get* method for the value type.

#### Environment Variables Not Applied

**Issue**: Environment variables are not being applied to the configuration.

**Solution**: Ensure that the environment variables are properly formatted with the correct prefix and separator. Check that the EnvPrefix and EnvReplacer options are set correctly.

## Related Components

- [Logging](../logging/README.md) - The logging component uses the config package for configuration.
- [Database](../db/README.md) - The database component uses the config package for database configuration.
- [Telemetry](../telemetry/README.md) - The telemetry component uses the config package for telemetry configuration.

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
