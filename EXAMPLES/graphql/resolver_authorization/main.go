// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to check authorization in a GraphQL resolver
package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// This is a simplified example to demonstrate authorization in a GraphQL resolver
// In a real application, you would use generated code from gqlgen

// Mock Item type for demonstration
type Item struct {
	ID    string
	Name  string
	Owner string
}

// Mock ItemInput type for demonstration
type ItemInput struct {
	Name  string
	Owner string
}

// Resolver is a simplified resolver for demonstration
type Resolver struct {
	logger *zap.Logger
}

// CreateItem is a resolver method that demonstrates manual authorization checking
func (r *Resolver) CreateItem(ctx context.Context, input ItemInput) (*Item, error) {
	// This example shows how to manually check authorization in a resolver
	fmt.Println("Example: Checking authorization in a GraphQL resolver")

	// Step 1: Check authorization manually
	fmt.Println("\nStep 1: Check authorization manually")
	fmt.Println(`
	// Check authorization manually if needed
	if err := graphql.CheckAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"CREATE"}, "ITEM", "CreateItem", r.logger); err != nil {
		return nil, err
	}
	`)

	// Step 2: Proceed with the operation if authorized
	fmt.Println("\nStep 2: Proceed with the operation if authorized")
	fmt.Println(`
	// Proceed with the operation
	item := &model.Item{
		ID:    uuid.New().String(),
		Name:  input.Name,
		Owner: input.Owner,
	}

	// Save the item to the database
	if err := r.db.SaveItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	return item, nil
	`)

	// Explain the authorization check
	fmt.Println("\nThe CheckAuthorization function:")
	fmt.Println("1. Extracts user roles, scopes, and resources from the context")
	fmt.Println("2. Checks if the user has any of the allowed roles")
	fmt.Println("3. Checks if the user has all the required scopes")
	fmt.Println("4. Checks if the user has access to the specified resource")
	fmt.Println("5. Returns an error if any of the checks fail")

	// Return a mock item for demonstration
	return &Item{
		ID:    "123",
		Name:  input.Name,
		Owner: input.Owner,
	}, nil
}

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a resolver
	resolver := &Resolver{
		logger: logger,
	}

	// Create a context
	ctx := context.Background()

	// Create an item input
	input := ItemInput{
		Name:  "Test Item",
		Owner: "user123",
	}

	// Call the resolver method
	item, err := resolver.CreateItem(ctx, input)
	if err != nil {
		fmt.Printf("Error creating item: %v\n", err)
		return
	}

	fmt.Printf("\nCreated item: ID=%s, Name=%s, Owner=%s\n", item.ID, item.Name, item.Owner)
	fmt.Println("\nNote: In a real application, the authorization check would use the actual user information from the context.")
}