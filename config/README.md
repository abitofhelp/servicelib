# Configuration

## Overview

The Configuration component provides a flexible and extensible system for managing application configuration from multiple sources and formats. It supports environment variables, files, and remote configuration services.

## Features

- **Multiple Sources**: Load configuration from files, environment variables, and remote services
- **Multiple Formats**: Support for JSON, YAML, TOML, and other formats
- **Dynamic Reloading**: Automatically reload configuration when changes are detected
- **Validation**: Validate configuration against schemas
- **Defaults**: Provide default values for configuration options

## Installation

```bash
go get github.com/abitofhelp/servicelib/config
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/config/basic_usage/README.md) for a complete, runnable example of how to use the configuration component.

## Configuration

See the [Custom Adapter example](../EXAMPLES/config/custom_adapter/README.md) for a complete, runnable example of how to configure the configuration component.

## API Documentation

### Core Types

The configuration component provides several core types for managing configuration.

#### Config

The main type that provides configuration functionality.

```
type Config struct {
    // Fields
}
```

#### Adapter

Interface for configuration adapters.

```
type Adapter interface {
    // Methods
}
```

### Key Methods

The configuration component provides several key methods for managing configuration.

#### Load

Loads configuration from a source.

```
func (c *Config) Load(ctx context.Context, source string) error
```

#### Get

Gets a configuration value.

```
func (c *Config) Get(ctx context.Context, key string) (interface{}, error)
```

#### Set

Sets a configuration value.

```
func (c *Config) Set(ctx context.Context, key string, value interface{}) error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [App Config](../EXAMPLES/config/app_config/README.md) - Application configuration
- [Basic Usage](../EXAMPLES/config/basic_usage/README.md) - Basic configuration operations
- [Custom Adapter](../EXAMPLES/config/custom_adapter/README.md) - Creating custom config adapters
- [Database Config](../EXAMPLES/config/database_config/README.md) - Database configuration

## Best Practices

1. **Use Environment Variables**: Use environment variables for sensitive configuration
2. **Validate Configuration**: Always validate configuration against schemas
3. **Provide Defaults**: Always provide default values for configuration options
4. **Handle Errors**: Properly handle configuration errors
5. **Use Typed Getters**: Use typed getters (GetString, GetInt, etc.) for type safety

## Troubleshooting

### Common Issues

#### Configuration Not Loading

If configuration is not loading, check that the source exists and is accessible.

#### Type Conversion Errors

If you're getting type conversion errors, use typed getters (GetString, GetInt, etc.) instead of Get.

## Related Components

- [Logging](../logging/README.md) - Logging for configuration events
- [Errors](../errors/README.md) - Error handling for configuration

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.