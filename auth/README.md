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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
