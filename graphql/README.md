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
  AUTHUSER
  # Add more roles as needed
}

directive @isAuthorized(allowedRoles: [Role!]!) on FIELD_DEFINITION
```

2. Apply the directive to your queries and mutations:

```graphql
type Query {
  getItem(id: ID!): Item @isAuthorized(allowedRoles: [ADMIN, AUTHUSER])
}

type Mutation {
  createItem(input: ItemInput!): Item! @isAuthorized(allowedRoles: [ADMIN])
}
```

3. Register the directive in your GraphQL server:

```go
schema := generated.NewExecutableSchema(generated.Config{
    Resolvers: resolverInstance,
    Directives: generated.DirectiveRoot{
        IsAuthorized: func(ctx context.Context, obj interface{}, next graphql.Resolver, allowedRoles []string) (interface{}, error) {
            return graphql.IsAuthorizedDirective(ctx, obj, next, allowedRoles, logger)
        },
    },
})
```

### JWT Authentication

The RBAC implementation works with the JWT authentication middleware from the `servicelib/auth` package. The middleware extracts the user's roles from the JWT token and adds them to the request context, which is then used by the `@isAuthorized` directive to check if the user has the required roles.

To set up JWT authentication:

1. Configure the auth service in your dependency injection container:

```go
authConfig := auth.DefaultConfig()
authConfig.JWT.SecretKey = cfg.Auth.JWT.SecretKey
authConfig.JWT.Issuer = cfg.Auth.JWT.Issuer
authConfig.JWT.TokenDuration = cfg.Auth.JWT.TokenDuration
authConfig.Middleware.SkipPaths = []string{"/health", "/metrics", "/playground"}

authService, err := auth.New(ctx, authConfig, logger)
if err != nil {
    return nil, fmt.Errorf("failed to initialize auth service: %w", err)
}
```

2. Apply the auth middleware to your HTTP handler:

```go
handler = authService.Middleware()(handler)
```

### Helper Functions

The package provides several helper functions for working with RBAC:

- `IsAuthorizedDirective`: The implementation of the `@isAuthorized` directive
- `CheckAuthorization`: A helper function for checking authorization in resolvers
- `ConvertRolesToStrings`: A helper function for converting enum roles to strings

## Example Usage

```go
// In your resolver
func (r *Resolver) CreateItem(ctx context.Context, input model.ItemInput) (*model.Item, error) {
    // Check authorization manually if needed
    if err := graphql.CheckAuthorization(ctx, []string{"ADMIN"}, "CreateItem", r.logger); err != nil {
        return nil, err
    }

    // Proceed with the operation
    // ...
}
```

## Metrics

The package provides metrics for authorization checks:

- `authorization_check_duration_seconds`: Duration of authorization checks in seconds
- `authorization_failures_total`: Total number of authorization failures

These metrics are automatically registered with Prometheus and can be used to monitor the performance and security of your GraphQL API.