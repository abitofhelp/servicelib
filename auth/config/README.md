# Auth Configuration Adapter

## Overview

The Auth Configuration Adapter provides adapters for integrating the auth configuration with the config package. It allows for a unified interface to access JWT, OIDC, Middleware, and Service configurations.

## Features

- **Configuration Adaptation**: Adapt auth.Config to the config package interfaces
- **Unified Interface**: Access JWT, OIDC, Middleware, and Service configurations through a unified interface
- **Configuration Conversion**: Convert auth configuration to generic configuration
- **Component Configuration**: Create specific configurations for JWT, OIDC, Middleware, and Service components

## Installation

```bash
go get github.com/abitofhelp/servicelib/auth/config
```

## Quick Start

See the [Auth Configuration example](../../EXAMPLES/auth/configuration/README.md) for a complete, runnable example of how to use the auth configuration adapter.

## API Documentation

### Core Types

#### AuthConfigAdapter

The main type that provides adaptation between auth configuration and the config package.

```go
type AuthConfigAdapter struct {
    // Fields
}
```

### Key Methods

#### NewAuthConfigAdapter

Creates a new auth configuration adapter.

```go
func NewAuthConfigAdapter(config auth.Config) *AuthConfigAdapter
```

#### GetAuth

Gets the auth configuration.

```go
func (a *AuthConfigAdapter) GetAuth() auth.Config
```

#### CreateJWTConfig

Creates JWT configuration from auth configuration.

```go
func CreateJWTConfig(config auth.Config) jwt.Config
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Configuration](../../EXAMPLES/auth/configuration/README.md) - Shows how to configure auth components

## Best Practices

1. **Use Default Configuration**: Start with the default configuration and customize as needed
2. **Validate Configuration**: Always validate configuration before using it
3. **Secure Secret Keys**: Store secret keys securely and never hardcode them in your application
4. **Use Environment Variables**: Use environment variables for configuration values that change between environments

## Troubleshooting

### Common Issues

#### Configuration Validation Failures

If configuration validation fails, check that all required fields are set and that the values are valid.

## Related Components

- [Auth](../README.md) - The parent auth package
- [JWT](../jwt/README.md) - JWT token handling
- [OIDC](../oidc/README.md) - OpenID Connect integration
- [Config](../../config/README.md) - Configuration management

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.