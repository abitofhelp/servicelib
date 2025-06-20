# Auth Module

The Auth Module provides comprehensive authentication and authorization functionality for Go applications. It includes JWT token handling, OIDC integration, HTTP middleware, and role-based access control.

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

See the [Quick Start example](../examples/auth/quickstart_example.go) for a complete, runnable example of how to use the Auth module.

## Configuration

The auth module can be configured using the `Config` struct. See the [Configuration example](../examples/auth/configuration_example.go) for a complete, runnable example of how to configure the Auth module.

## API Documentation

### Auth

The `Auth` struct is the main entry point for the auth module. It provides methods for authentication and authorization.

#### Creating an Auth Instance

See the [Auth Instance example](../examples/auth/auth_instance_example.go) for a complete, runnable example of how to create an Auth instance.

#### Middleware

See the [Middleware example](../examples/auth/middleware_example.go) for a complete, runnable example of how to use the Auth middleware.

#### Token Handling

See the [Token Handling example](../examples/auth/token_handling_example.go) for a complete, runnable example of how to generate and validate tokens.

#### Authorization

See the [Authorization example](../examples/auth/authorization_example.go) for a complete, runnable example of how to perform authorization checks.

#### User Information

See the [User Information example](../examples/auth/user_info_example.go) for a complete, runnable example of how to get user information from the context.

### Context Utilities

The auth module provides utilities for working with context. See the [Context Utilities example](../examples/auth/context_utilities_example.go) for a complete, runnable example of how to use these utilities.

## Error Handling

The auth module provides comprehensive error handling with context-aware errors. See the [Error Handling example](../examples/auth/error_handling_example.go) for a complete, runnable example of how to handle errors.

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
