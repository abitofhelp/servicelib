// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to register the @isAuthorized directive in a GraphQL server
package example_graphql

import (
	"fmt"
)

// This is a simplified example to demonstrate the concept of registering
// the @isAuthorized directive in a GraphQL server.
// In a real application, you would use generated code from gqlgen.

func main() {

	// In a real application with gqlgen, you would register the directive like this:
	fmt.Println("Example: Registering @isAuthorized directive in GraphQL server")

	// This is a simplified representation of what the generated code would look like
	fmt.Println(`
schema := generated.NewExecutableSchema(generated.Config{
    Resolvers: resolverInstance,
    Directives: generated.DirectiveRoot{
        IsAuthorized: func(ctx context.Context, obj interface{}, next graphql.Resolver, allowedRoles []string, requiredScopes []string, resource string) (interface{}, error) {
            return graphql.IsAuthorizedDirective(ctx, obj, next, allowedRoles, requiredScopes, resource, logger)
        },
    },
})`)

	// Demonstrate how to use the IsAuthorizedDirective function
	fmt.Println("\nThe IsAuthorizedDirective function checks if the user has the required roles, scopes, and access to the resource.")
	fmt.Println("It extracts user information from the context and performs the authorization check.")

	// Example of what happens inside IsAuthorizedDirective (shown as code)
	fmt.Println(`
	// Example of checking authorization manually
	ctx := context.Background()
	allowedRoles := []string{"ADMIN", "EDITOR"}
	requiredScopes := []string{"READ"}
	resource := "ITEM"

	err := graphql.CheckAuthorization(ctx, allowedRoles, requiredScopes, resource, "ExampleOperation", logger)
	if err != nil {
		// Handle authorization failure
		fmt.Printf("Authorization failed: %v\n", err)
		return nil, err
	}

	// Proceed with the operation if authorized
	return next(ctx)
	`)
}
