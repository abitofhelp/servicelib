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

## Best Practices

1. Use the auth configuration adapter to access auth configuration through a unified interface
2. Use the helper functions to create specific configurations for JWT, OIDC, Middleware, and Service components
3. Use the `AsGenericConfig` method to convert auth configuration to generic configuration when needed
4. Follow the examples in the `example_test.go` file for guidance on how to use the package