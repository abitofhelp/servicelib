# GraphQL Package

This package provides utilities for working with GraphQL in Go services, including server configuration, error handling, and role-based access control (RBAC).

## Features

- GraphQL server configuration with sensible defaults
- Error handling with proper context and logging
- Role-based access control (RBAC) with the `@isAuthorized` directive
- Metrics and tracing for GraphQL operations
- Timeout and cancellation handling

## Role-Based Access Control (RBAC)

The package provides a generic implementation of the `@isAuthorized` directive for GraphQL, which can be used to restrict access to GraphQL operations based on user roles.

### Setup

1. Define a Role enum in your GraphQL schema:

```graphql
enum Role {
  ADMIN
  EDITOR
  VIEWER
  # Add more roles as needed
}

enum Scope {
  READ
  WRITE
  DELETE
  CREATE
}

enum Resource {
  FAMILY
  PARENT
  CHILD
  ITEM
}

directive @isAuthorized(
  allowedRoles: [Role!]!, 
  requiredScopes: [Scope!] = [READ], 
  resource: Resource = FAMILY
) on FIELD_DEFINITION
```

2. Apply the directive to your queries and mutations:

```graphql
type Query {
  getItem(id: ID!): Item @isAuthorized(
    allowedRoles: [ADMIN, EDITOR, VIEWER], 
    requiredScopes: [READ], 
    resource: ITEM
  )
}

type Mutation {
  createItem(input: ItemInput!): Item! @isAuthorized(
    allowedRoles: [ADMIN, EDITOR], 
    requiredScopes: [CREATE], 
    resource: ITEM
  )
}
```

3. Register the directive in your GraphQL server:

See the [Directive Registration example](../examples/graphql/directive_registration_example.go) for a complete, runnable example of how to register the @isAuthorized directive in your GraphQL server.

### JWT Authentication

The RBAC implementation works with the JWT authentication middleware from the `servicelib/auth` package. The middleware extracts the user's roles, scopes, and resources from the JWT token and adds them to the request context, which is then used by the `@isAuthorized` directive to check if the user has the required roles, scopes, and access to the specified resource.

To set up JWT authentication:

1. Configure the auth service in your dependency injection container:

See the [Auth Configuration example](../examples/graphql/auth_configuration_example.go) for a complete, runnable example of how to configure the auth service for GraphQL.

2. Apply the auth middleware to your HTTP handler:

See the [Auth Middleware example](../examples/graphql/auth_middleware_example.go) for a complete, runnable example of how to apply the auth middleware to a GraphQL handler.

### Helper Functions

The package provides several helper functions for working with RBAC:

- `IsAuthorizedDirective`: The implementation of the `@isAuthorized` directive that checks roles, scopes, and resources
- `CheckAuthorization`: A helper function for checking authorization in resolvers with roles, scopes, and resources
- `ConvertRolesToStrings`: A generic helper function for converting enum types to strings
- `IsAuthorizedWithScopes`: A function that checks if a user has the required roles, scopes, and access to a resource
- `HasScope`: A function that checks if a user has a specific scope
- `HasResource`: A function that checks if a user has access to a specific resource

## Example Usage

See the [Resolver Authorization example](../examples/graphql/resolver_authorization_example.go) for a complete, runnable example of how to check authorization in a GraphQL resolver.

## Generating JWT Tokens

To test the RBAC implementation, you can use the `genjwt` tool to generate JWT tokens with different roles, scopes, and resources:

See the [JWT Token Generation example](../examples/graphql/jwt_token_generation_example.go) for a complete, runnable example of how to generate JWT tokens for testing GraphQL RBAC.

## Metrics

The package provides metrics for authorization checks:

- `authorization_check_duration_seconds`: Duration of authorization checks in seconds
- `authorization_failures_total`: Total number of authorization failures

These metrics are automatically registered with Prometheus and can be used to monitor the performance and security of your GraphQL API.
