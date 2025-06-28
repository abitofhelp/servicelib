// Copyright (c) 2025 A Bit of Help, Inc.

// Package graphql provides utilities for working with GraphQL in Go applications.
//
// This package offers a comprehensive solution for implementing GraphQL APIs,
// including server configuration, middleware integration, error handling, and
// authorization directives. It builds on top of the gqlgen library to provide
// additional features and integrations with other components of the servicelib.
//
// Key features:
//   - GraphQL server configuration with sensible defaults
//   - Integration with logging, tracing, and metrics
//   - Authorization directives for securing GraphQL operations
//   - Standardized error handling and conversion
//   - Request timeout management
//   - Performance monitoring and metrics collection
//
// The package is designed to work seamlessly with other servicelib components,
// such as auth, logging, telemetry, and middleware.
//
// Example usage for creating a GraphQL server:
//
//	// Create a logger
//	logger := logging.NewContextLogger(zapLogger)
//
//	// Create server config with default settings
//	serverConfig := graphql.NewDefaultServerConfig()
//
//	// Create the GraphQL server
//	server := graphql.NewServer(
//	    generatedGraphQL.NewExecutableSchema(generatedGraphQL.Config{
//	        Resolvers: &resolvers.Resolver{},
//	        Directives: generatedGraphQL.DirectiveRoot{
//	            IsAuthorized: graphql.IsAuthorizedDirective,
//	        },
//	    }),
//	    logger,
//	    serverConfig,
//	)
//
//	// Use the server with an HTTP handler
//	http.Handle("/graphql", server)
//
// Example usage for authorization in resolvers:
//
//	func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
//	    // Check if the user is authorized to view user details
//	    err := graphql.CheckAuthorization(
//	        ctx,
//	        []string{"ADMIN", "USER_MANAGER"},
//	        []string{"users:read"},
//	        "user",
//	        "GetUser",
//	        r.logger,
//	    )
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Proceed with fetching the user
//	    user, err := r.userService.GetUser(ctx, id)
//	    if err != nil {
//	        // Handle and convert errors to GraphQL errors
//	        return nil, graphql.HandleError(ctx, err, "GetUser", r.logger)
//	    }
//
//	    return user, nil
//	}
//
// The package also provides utilities for working with GraphQL errors:
//
//	// Convert application errors to GraphQL errors
//	if err != nil {
//	    return nil, graphql.HandleError(ctx, err, "OperationName", logger)
//	}
//
// This error handling system automatically maps different error types to
// appropriate GraphQL errors with the right extensions and codes, while
// also ensuring that sensitive error details are not exposed to clients.
package graphql