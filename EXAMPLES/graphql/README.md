# GraphQL Package Examples

This directory contains examples demonstrating how to use the `graphql` package, which provides utilities for integrating authentication and authorization with GraphQL in Go applications. The package offers functionality for securing GraphQL APIs, handling JWT tokens, and implementing authorization directives.

## Examples

### 1. Auth Configuration Example

[auth_configuration_example.go](auth_configuration_example.go)

Demonstrates how to configure the auth service for use with GraphQL.

Key concepts:
- Setting up JWT configuration (secret key, issuer, token duration)
- Configuring middleware to skip certain paths
- Initializing the auth service for GraphQL
- Understanding the auth service capabilities

### 2. Auth Middleware Example

[auth_middleware_example.go](auth_middleware_example.go)

Shows how to integrate authentication middleware with GraphQL.

Key concepts:
- Setting up auth middleware for GraphQL handlers
- Protecting GraphQL endpoints
- Handling authentication in GraphQL requests
- Passing authentication context to resolvers

### 3. Directive Registration Example

[directive_registration_example.go](directive_registration_example.go)

Demonstrates how to register and use GraphQL directives for authorization.

Key concepts:
- Creating custom GraphQL directives
- Registering directives with the GraphQL schema
- Using directives for field-level authorization
- Implementing directive logic

### 4. JWT Token Generation Example

[jwt_token_generation_example.go](jwt_token_generation_example.go)

Shows how to generate JWT tokens for GraphQL authentication.

Key concepts:
- Creating JWT tokens with appropriate claims
- Setting token expiration
- Handling token generation errors
- Using tokens for GraphQL authentication

### 5. Resolver Authorization Example

[resolver_authorization_example.go](resolver_authorization_example.go)

Demonstrates how to implement authorization checks in GraphQL resolvers.

Key concepts:
- Accessing authentication context in resolvers
- Performing permission checks
- Handling unauthorized access
- Implementing role-based access control

### 6. Authorization Check Example

[authorization/main.go](authorization/main.go)

Shows how to properly check authorization for GraphQL operations using the servicelib authorization middleware.

Key concepts:
- Using the CheckAuthorization function
- Handling authorization errors
- Logging authorization failures
- Implementing role and scope-based access control

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/graphql/auth_configuration_example.go
```

## Additional Resources

For more information about the graphql package, see the [graphql package documentation](../../graphql/README.md).
