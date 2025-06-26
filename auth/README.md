# Authentication & Authorization

## Overview

The Authentication & Authorization component provides a comprehensive solution for securing your services with JWT-based authentication and role-based authorization. It integrates with OpenID Connect (OIDC) providers and offers middleware for HTTP request authentication.

## Features

- **JWT Token Management**: Generate and validate JWT tokens
- **Role-Based Authorization**: Control access based on user roles
- **OpenID Connect Integration**: Authenticate users via OIDC providers
- **HTTP Middleware**: Protect HTTP endpoints with authentication middleware
- **Context Utilities**: Access user information from context

## Installation

```bash
go get github.com/abitofhelp/servicelib/auth
```

## Quick Start

See the [Quick Start example](../EXAMPLES/auth/quickstart/README.md) for a complete, runnable example of how to use the auth component.

## Configuration

See the [Configuration example](../EXAMPLES/auth/configuration/README.md) for a complete, runnable example of how to configure the auth component.

## API Documentation

### Core Types

The auth component provides several core types for authentication and authorization.

#### Auth

The main type that provides authentication and authorization functionality.

```
type Auth struct {
    // Fields
}
```

#### Config

Configuration for the Auth component.

```
type Config struct {
    // Fields
}
```

### Key Methods

The auth component provides several key methods for authentication and authorization.

#### GenerateToken

Generates a JWT token for a user with specified roles, scopes, and resources.

```
func (a *Auth) GenerateToken(ctx context.Context, userID string, roles []string, scopes []string, resources []string) (string, error)
```

#### ValidateToken

Validates a JWT token and returns the claims.

```
func (a *Auth) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error)
```

#### IsAuthorized

Checks if the user in the context is authorized for a specific operation.

```
func (a *Auth) IsAuthorized(ctx context.Context, operation string) (bool, error)
```

#### HasRole

Checks if the user in the context has a specific role.

```
func (a *Auth) HasRole(ctx context.Context, role string) (bool, error)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Auth Instance](../EXAMPLES/auth/auth_instance/README.md) - Creating and configuring auth instances
- [Authorization](../EXAMPLES/auth/authorization/README.md) - Implementing authorization checks
- [Configuration](../EXAMPLES/auth/configuration/README.md) - Configuring auth components
- [Context Utilities](../EXAMPLES/auth/context_utilities/README.md) - Working with auth context
- [Middleware](../EXAMPLES/auth/middleware/README.md) - Using auth middleware
- [Token Handling](../EXAMPLES/auth/token_handling/README.md) - Working with JWT tokens

## Best Practices

1. **Use Middleware**: Always use the auth middleware to protect your HTTP endpoints
2. **Check Authorization**: Always check authorization before performing sensitive operations
3. **Validate Tokens**: Always validate tokens before trusting their contents
4. **Use Context**: Use context to pass user information between functions
5. **Handle Errors**: Properly handle authentication and authorization errors

## Troubleshooting

### Common Issues

#### Token Validation Failures

If token validation fails, check that the token is not expired and that the signing key is correct.

#### Authorization Failures

If authorization checks fail, ensure that the user has the required roles and that the context contains the user information.

## Related Components

- [Errors](../errors/README.md) - Error handling for authentication and authorization
- [Logging](../logging/README.md) - Logging for authentication and authorization events
- [Middleware](../middleware/README.md) - HTTP middleware for authentication

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
