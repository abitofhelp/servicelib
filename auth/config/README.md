# Auth Configuration Adapter
This package provides adapters for integrating the auth configuration with the config package in the servicelib library.


## Overview

The auth configuration adapter allows you to:

1. Adapt the auth.Config to the config package interfaces
2. Access JWT, OIDC, Middleware, and Service configurations through a unified interface
3. Convert auth configuration to generic configuration
4. Create specific configurations for JWT, OIDC, Middleware, and Service components

## Usage

### Basic Usage

```go
// Create an auth configuration
config := auth.DefaultConfig()
config.JWT.SecretKey = "example-secret-key"
config.OIDC.IssuerURL = "https://example.com/oidc"
config.OIDC.ClientID = "example-client-id"
config.OIDC.ClientSecret = "example-client-secret"

// Create an auth config adapter
adapter := authconfig.NewAuthConfigAdapter(config)

// Get the auth configuration
authCfg := adapter.GetAuth()

// Use the auth configuration to create JWT, OIDC, middleware, and service configurations
jwtConfig := authconfig.CreateJWTConfig(authCfg)
jwtRemoteConfig := authconfig.CreateJWTRemoteConfig(authCfg)
oidcConfig := authconfig.CreateOIDCConfig(authCfg)
middlewareConfig := authconfig.CreateMiddlewareConfig(authCfg)
serviceConfig := authconfig.CreateServiceConfig(authCfg)

// Create services using these configurations
jwtService := jwt.NewService(jwtConfig, logger)
oidcService, _ := oidc.NewService(context.Background(), oidcConfig, logger)
middleware := middleware.NewMiddlewareWithOIDC(jwtService, oidcService, middlewareConfig, logger)
authService := service.NewService(serviceConfig, logger)
```

### Converting to Generic Config

```go
// Create an auth configuration
config := auth.DefaultConfig()

// Create an auth config adapter
adapter := authconfig.NewAuthConfigAdapter(config)

// Convert to generic config
genericConfig := adapter.AsGenericConfig()

// Use the generic config
appConfig := genericConfig.GetApp()
fmt.Println("App name:", appConfig.GetName())
fmt.Println("App environment:", appConfig.GetEnvironment())
```

## Interfaces

The package provides the following interfaces:

- `Auth`: Interface for auth configuration
- `JWT`: Interface for JWT configuration
- `JWTRemote`: Interface for JWT remote configuration
- `OIDC`: Interface for OIDC configuration
- `Middleware`: Interface for middleware configuration
- `Service`: Interface for service configuration

## Helper Functions

The package provides the following helper functions:

- `CreateJWTConfig`: Creates a JWT configuration from the auth configuration
- `CreateJWTRemoteConfig`: Creates a JWT remote configuration from the auth configuration
- `CreateOIDCConfig`: Creates an OIDC configuration from the auth configuration
- `CreateMiddlewareConfig`: Creates a middleware configuration from the auth configuration
- `CreateServiceConfig`: Creates a service configuration from the auth configuration


## Features

- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3

## Installation

```bash
go get github.com/abitofhelp/servicelib/config
```

## Quick Start

See the [Quick Start example](../EXAMPLES/config/quickstart_example.go) for a complete, runnable example of how to use the config.

## Configuration

See the [Configuration example](../EXAMPLES/config/configuration_example.go) for a complete, runnable example of how to configure the config.

## API Documentation


### Core Types

Description of the main types provided by the config.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/config/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the config.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/config/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/config/basic_usage_example.go) - Shows basic usage of the config
- [Advanced Configuration](../EXAMPLES/config/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/config/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. Use the auth configuration adapter to access auth configuration through a unified interface
2. Use the helper functions to create specific configurations for JWT, OIDC, Middleware, and Service components
3. Use the `AsGenericConfig` method to convert auth configuration to generic configuration when needed
4. Follow the examples in the `example_test.go` file for guidance on how to use the package

## Troubleshooting

### Common Issues

#### Issue 1

Description of issue 1 and how to resolve it.

#### Issue 2

Description of issue 2 and how to resolve it.

## Related Components

- [Component 1](../config1/README.md) - Description of how this config relates to Component 1
- [Component 2](../config2/README.md) - Description of how this config relates to Component 2

## Contributing

Contributions to this config are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
