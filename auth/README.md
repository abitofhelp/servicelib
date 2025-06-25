# Auth Module
The Auth Module provides comprehensive authentication and authorization functionality for Go applications. It includes JWT token handling, OIDC integration, HTTP middleware, and role-based access control.


## Overview

The Auth module provides a comprehensive authentication and authorization system for Go applications. It serves as a central component in the ServiceLib library for securing APIs and services. The module implements industry-standard security protocols including JWT (JSON Web Tokens) and OIDC (OpenID Connect), offering a flexible and robust solution for identity management and access control.

This module is designed to be easily integrated with HTTP services, providing middleware for request authentication, role-based access control for authorization decisions, and utilities for working with user identity information throughout your application.

## Features

- **JWT Token Handling**: Generate and validate JWT tokens
- **OIDC Integration**: Validate tokens from OpenID Connect providers
- **HTTP Middleware**: Authenticate HTTP requests
- **Role-Based Access Control**: Control access to resources based on user roles
- **Comprehensive Error Handling**: Context-aware error handling with detailed error information
- **Tracing Integration**: OpenTelemetry tracing for all operations
- **Logging Integration**: Structured logging with zap


## Installation

```bash
go get github.com/abitofhelp/servicelib/auth
```


## Quick Start

See the [Quick Start example](../EXAMPLES/auth/quickstart_example.go) for a complete, runnable example of how to use the Auth module.


## Configuration

The auth module can be configured using the `Config` struct. See the [Configuration example](../EXAMPLES/auth/configuration_example.go) for a complete, runnable example of how to configure the Auth module.


## API Documentation

### Auth

The `Auth` struct is the main entry point for the auth module. It provides methods for authentication and authorization.

#### Creating an Auth Instance

See the [Auth Instance example](../EXAMPLES/auth/auth_instance_example.go) for a complete, runnable example of how to create an Auth instance.

#### Middleware

See the [Middleware example](../EXAMPLES/auth/middleware_example.go) for a complete, runnable example of how to use the Auth middleware.

#### Token Handling

See the [Token Handling example](../EXAMPLES/auth/token_handling_example.go) for a complete, runnable example of how to generate and validate tokens.

#### Authorization

See the [Authorization example](../EXAMPLES/auth/authorization_example.go) for a complete, runnable example of how to perform authorization checks.

#### User Information

See the [User Information example](../EXAMPLES/auth/user_info_example.go) for a complete, runnable example of how to get user information from the context.

### Context Utilities

The auth module provides utilities for working with context. See the [Context Utilities example](../EXAMPLES/auth/context_utilities_example.go) for a complete, runnable example of how to use these utilities.

## Error Handling

The auth module provides comprehensive error handling with context-aware errors. See the [Error Handling example](../EXAMPLES/auth/error_handling_example.go) for a complete, runnable example of how to handle errors.


### Core Types

The auth module provides several core types for authentication and authorization:

#### Auth

The `Auth` struct is the main entry point for the auth module. It encapsulates all the functionality for authentication and authorization, including JWT token handling, OIDC integration, and role-based access control.

#### Config

The `Config` struct holds the configuration for the auth module. It includes settings for JWT, OIDC, middleware, and service components:

- **JWT**: Configuration for JWT token generation and validation
- **OIDC**: Configuration for OpenID Connect integration
- **Middleware**: Configuration for HTTP middleware
- **Service**: Configuration for authorization services

#### Claims

The `Claims` struct represents the claims in a JWT token, including standard claims like subject, issuer, and expiration time, as well as custom claims like roles and scopes.

See the [Auth Instance example](../EXAMPLES/auth/auth_instance_example.go) for a complete, runnable example of how to create and use these types.

### Key Methods

The auth module provides several key methods for authentication and authorization:

#### Middleware

The `Middleware()` method returns an HTTP middleware function that authenticates incoming requests. The middleware extracts and validates JWT tokens from request headers, populating the request context with user information if authentication is successful.

See the [Middleware example](../EXAMPLES/auth/middleware_example.go) for a complete, runnable example of how to use the middleware.

#### GenerateToken

The `GenerateToken()` method generates a JWT token for the specified user with the given roles and scopes. The token includes standard claims like issuer, subject, and expiration time, as well as custom claims for roles and scopes.

See the [Token Handling example](../EXAMPLES/auth/token_handling_example.go) for a complete, runnable example of how to generate tokens.

#### IsAuthorized

The `IsAuthorized()` method checks if the user in the current context is authorized to perform the specified operation. The authorization decision is based on the user's roles and the operation being performed.

See the [Authorization example](../EXAMPLES/auth/authorization_example.go) for a complete, runnable example of how to perform authorization checks.

#### User Information Methods

The auth module provides several methods for retrieving user information from the context:

- `GetUserID()`: Returns the user ID from the context
- `GetUserRoles()`: Returns the user roles from the context
- `IsAdmin()`: Checks if the user has the admin role
- `HasRole()`: Checks if the user has a specific role

See the [User Information example](../EXAMPLES/auth/user_info_example.go) for a complete, runnable example of how to use these methods.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/auth/basic_usage_example.go) - Shows basic usage of the auth
- [Advanced Configuration](../EXAMPLES/auth/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/auth/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Secure Secret Keys**: Always store JWT secret keys securely, preferably in environment variables or a secure key management system.

2. **Token Expiration**: Set appropriate expiration times for tokens based on your security requirements.

3. **HTTPS**: Always use HTTPS in production to protect tokens in transit.

4. **Validate All Tokens**: Always validate tokens before trusting their contents.

5. **Role-Based Access Control**: Use role-based access control to limit access to sensitive operations.

6. **Error Handling**: Implement proper error handling for authentication and authorization failures.

7. **Logging**: Log authentication and authorization events for audit purposes, but be careful not to log sensitive information.


## Troubleshooting

### Common Issues

#### Token Validation Failures

**Issue**: JWT token validation fails with "invalid signature" error.

**Solution**: Ensure that the same secret key is used for both token generation and validation. Check that the token hasn't been tampered with.

#### OIDC Provider Connection Issues

**Issue**: Unable to connect to the OIDC provider.

**Solution**: Check network connectivity, firewall settings, and that the OIDC provider is running and accessible.

#### Authorization Failures

**Issue**: User has the correct token but is not authorized to perform an operation.

**Solution**: Check that the user has the required roles or permissions for the operation. Verify that the role-based access control configuration is correct.


## Related Components

- [Middleware](../middleware/README.md) - The middleware component includes HTTP middleware for authentication and authorization.
- [Context](../context/README.md) - The context component provides utilities for working with context, which is used extensively in the auth module.
- [Logging](../logging/README.md) - The logging component provides structured logging, which is used by the auth module for logging authentication and authorization events.
- [Telemetry](../telemetry/README.md) - The telemetry component provides tracing, which is used by the auth module for tracing authentication and authorization operations.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
