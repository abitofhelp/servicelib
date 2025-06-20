# GraphQL Package

## Overview

The `graphql` package provides utilities for working with GraphQL in Go services, including server configuration, error handling, and role-based access control (RBAC).

## Features

- **GraphQL Server Configuration**: Sensible defaults for GraphQL server setup
- **Error Handling**: Proper context and logging for GraphQL errors
- **Role-Based Access Control**: RBAC with the `@isAuthorized` directive
- **Metrics and Tracing**: Performance monitoring for GraphQL operations
- **Timeout and Cancellation**: Proper handling of timeouts and cancellations

## Installation

```bash
go get github.com/abitofhelp/servicelib/graphql
```

## Quick Start

See the [Directive Registration example](../examples/graphql/directive_registration_example.go) for a complete, runnable example of how to set up a GraphQL server with the `@isAuthorized` directive.

## Configuration

See the [Auth Configuration example](../examples/graphql/auth_configuration_example.go) for a complete, runnable example of how to configure the auth service for GraphQL.

## API Documentation

### Core Types

#### IsAuthorizedDirective

The `IsAuthorizedDirective` function implements the `@isAuthorized` directive for GraphQL, which can be used to restrict access to GraphQL operations based on user roles.

See the [Directive Registration example](../examples/graphql/directive_registration_example.go) for a complete, runnable example of how to use the IsAuthorizedDirective.

#### CheckAuthorization

The `CheckAuthorization` function is a helper for checking authorization in resolvers with roles, scopes, and resources.

See the [Resolver Authorization example](../examples/graphql/resolver_authorization_example.go) for a complete, runnable example of how to use the CheckAuthorization function.

### Key Methods

#### IsAuthorizedWithScopes

The `IsAuthorizedWithScopes` function checks if a user has the required roles, scopes, and access to a resource.

See the [Resolver Authorization example](../examples/graphql/resolver_authorization_example.go) for a complete, runnable example of how to use the IsAuthorizedWithScopes function.

#### HasScope

The `HasScope` function checks if a user has a specific scope.

#### HasResource

The `HasResource` function checks if a user has access to a specific resource.

## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Directive Registration](../examples/graphql/directive_registration_example.go) - Shows how to register the @isAuthorized directive
- [Auth Configuration](../examples/graphql/auth_configuration_example.go) - Shows how to configure the auth service
- [Auth Middleware](../examples/graphql/auth_middleware_example.go) - Shows how to apply the auth middleware
- [Resolver Authorization](../examples/graphql/resolver_authorization_example.go) - Shows how to check authorization in resolvers
- [JWT Token Generation](../examples/graphql/jwt_token_generation_example.go) - Shows how to generate JWT tokens for testing

## Best Practices

1. **Schema Design**: Define clear Role, Scope, and Resource enums in your GraphQL schema.

2. **Directive Usage**: Apply the `@isAuthorized` directive to all operations that require authorization.

3. **JWT Configuration**: Configure the JWT authentication middleware to extract roles, scopes, and resources from tokens.

4. **Error Handling**: Implement proper error handling for authorization failures.

5. **Testing**: Generate test tokens with different roles, scopes, and resources to verify your authorization logic.

## Troubleshooting

### Common Issues

#### Authorization Failures

**Issue**: Users with valid tokens are not authorized to perform operations.

**Solution**: Check that the token contains the correct roles, scopes, and resources. Verify that the `@isAuthorized` directive is configured correctly.

#### Directive Registration Errors

**Issue**: Errors when registering the `@isAuthorized` directive.

**Solution**: Ensure that your GraphQL schema includes the Role, Scope, and Resource enums and the directive definition.

#### Performance Issues

**Issue**: Authorization checks are causing performance issues.

**Solution**: Use the provided metrics to identify bottlenecks. Consider caching authorization results for frequently accessed operations.

## Related Components

- [Auth](../auth/README.md) - The auth component provides JWT authentication middleware used by the graphql package.
- [Middleware](../middleware/README.md) - The middleware component includes HTTP middleware for authentication and authorization.
- [Telemetry](../telemetry/README.md) - The telemetry component provides metrics and tracing used by the graphql package.

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
