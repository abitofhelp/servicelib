# Auth Package Examples

This directory contains examples demonstrating how to use the `auth` package, which provides authentication and authorization utilities for Go applications. The package offers functionality for JWT token handling, middleware integration, user information retrieval, and authorization checks.

## Examples

### 1. Quickstart Example

[quickstart_example.go](quickstart_example.go)

Demonstrates the basic setup and usage of the auth package.

Key concepts:
- Creating an auth configuration
- Initializing an auth instance
- Using auth middleware for HTTP handlers
- Checking if a user is authorized to perform an operation
- Getting a user ID from the context

### 2. Auth Instance Example

[auth_instance_example.go](auth_instance_example.go)

Shows how to create and configure an auth instance.

Key concepts:
- Creating an auth instance with custom configuration
- Handling initialization errors
- Auth instance lifecycle management

### 3. Authorization Example

[authorization_example.go](authorization_example.go)

Demonstrates how to perform authorization checks.

Key concepts:
- Checking permissions and roles
- Role-based access control
- Permission-based authorization
- Handling authorization errors

### 4. Configuration Example

[configuration_example.go](configuration_example.go)

Shows how to configure the auth package.

Key concepts:
- Setting up JWT configuration
- Configuring token validation parameters
- Setting up OIDC integration
- Using default and custom configurations

### 5. Context Utilities Example

[context_utilities_example.go](context_utilities_example.go)

Demonstrates working with auth-related information in context.

Key concepts:
- Adding auth information to context
- Retrieving auth information from context
- Context propagation with auth data

### 6. Error Handling Example

[error_handling_example.go](error_handling_example.go)

Shows how to handle various auth-related errors.

Key concepts:
- Handling token validation errors
- Dealing with authorization failures
- Graceful error recovery
- Error types and classification

### 7. Middleware Example

[middleware_example.go](middleware_example.go)

Demonstrates how to use auth middleware with HTTP servers.

Key concepts:
- Integrating auth middleware with HTTP handlers
- Protecting routes with authentication
- Customizing middleware behavior

### 8. Token Handling Example

[token_handling_example.go](token_handling_example.go)

Shows how to work with JWT tokens.

Key concepts:
- Generating JWT tokens
- Validating tokens
- Extracting claims from tokens
- Token refresh strategies

### 9. User Info Example

[user_info_example.go](user_info_example.go)

Demonstrates how to work with user information.

Key concepts:
- Retrieving user details
- Working with user profiles
- Handling user-related operations

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/auth/quickstart_example.go
```

## Additional Resources

For more information about the auth package, see the [auth package documentation](../../auth/README.md).