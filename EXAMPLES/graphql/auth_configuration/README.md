# graphql auth_configuration Example

## Overview

This example demonstrates how to configure the authentication service for use with GraphQL in the ServiceLib library. It shows how to set up JWT authentication settings and configure middleware to bypass authentication for specific endpoints.

## Features

- **JWT Configuration**: Configure JWT token settings including secret key, issuer, and token duration
- **Middleware Skip Paths**: Set up paths that bypass authentication middleware
- **Auth Service Initialization**: Initialize the auth service with the configuration
- **GraphQL Integration**: Prepare auth service for use with GraphQL endpoints

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Auth Configuration Setup

The example shows how to create and configure the auth service with JWT settings:

```go
// Create a configuration for the auth service
authConfig := auth.DefaultConfig()
authConfig.JWT.SecretKey = "your-secret-key"
authConfig.JWT.Issuer = "your-service"
authConfig.JWT.TokenDuration = 24 * time.Hour
authConfig.Middleware.SkipPaths = []string{"/health", "/metrics", "/playground"}
```

### Auth Service Initialization

The example demonstrates how to initialize the auth service with the configuration:

```go
// Initialize the auth service
authService, err := auth.New(ctx, authConfig, logger)
if err != nil {
    fmt.Printf("Error initializing auth service: %v\n", err)
    return
}
```

### GraphQL Integration Considerations

The example specifically configures the auth service for GraphQL by:

1. Setting up skip paths for GraphQL-specific endpoints like `/playground`
2. Configuring JWT settings appropriate for GraphQL authentication
3. Preparing the auth service for integration with GraphQL resolvers

## Expected Output

When you run the example, you should see output similar to:

```
Example: Configuring auth service for GraphQL

Step 1: Configure the auth service

        authConfig := auth.DefaultConfig()
        authConfig.JWT.SecretKey = cfg.Auth.JWT.SecretKey
        authConfig.JWT.Issuer = cfg.Auth.JWT.Issuer
        authConfig.JWT.TokenDuration = cfg.Auth.JWT.TokenDuration
        authConfig.Middleware.SkipPaths = []string{"/health", "/metrics", "/playground"}


Step 2: Initialize the auth service

        authService, err := auth.New(ctx, authConfig, logger)
        if err != nil {
                return nil, fmt.Errorf("failed to initialize auth service: %w", err)
        }


Auth service configured successfully!
JWT Issuer: your-service
Token Duration: 24h0m0s
Skip Paths: [/health /metrics /playground]

The auth service provides:
- JWT token generation and validation
- Middleware for authenticating HTTP requests
- Functions for checking authorization
- Context utilities for working with user information
```

## Related Examples


- [auth_middleware](../auth_middleware/README.md) - Related example for auth_middleware
- [authorization](../authorization/README.md) - Related example for authorization
- [directive_registration](../directive_registration/README.md) - Related example for directive_registration

## Related Components

- [graphql Package](../../../graphql/README.md) - The graphql package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
