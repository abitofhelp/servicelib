// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating context utilities for authentication
package main

import (
	"context"
	"fmt"
)

// This is a simplified example to demonstrate context utilities
// In a real application, you would use the auth module's methods

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

// Context keys
const (
	UserIDKey    ContextKey = "userID"
	UserRolesKey ContextKey = "userRoles"
)

// AddUserIDToCtx adds a user ID to the context
func AddUserIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// AddUserRolesToCtx adds user roles to the context
func AddUserRolesToCtx(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, UserRolesKey, roles)
}

// GetUserIDFromCtx gets the user ID from the context
func GetUserIDFromCtx(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserRolesFromCtx gets the user roles from the context
func GetUserRolesFromCtx(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(UserRolesKey).([]string)
	return roles, ok
}

// IsUserAuthenticated checks if the user is authenticated
func IsUserAuthenticated(ctx context.Context) bool {
	_, ok := GetUserIDFromCtx(ctx)
	return ok
}

func main() {
	// Create a base context
	ctx := context.Background()

	// Check if user is authenticated (should be false)
	fmt.Printf("Is authenticated (before): %v\n", IsUserAuthenticated(ctx))

	// Add user ID to context
	userID := "user123"
	ctx = AddUserIDToCtx(ctx, userID)

	// Check if user is authenticated (should be true)
	fmt.Printf("Is authenticated (after adding user ID): %v\n", IsUserAuthenticated(ctx))

	// Get user ID from context
	retrievedUserID, ok := GetUserIDFromCtx(ctx)
	if ok {
		fmt.Printf("User ID from context: %s\n", retrievedUserID)
	} else {
		fmt.Println("User ID not found in context")
	}

	// Add user roles to context
	roles := []string{"admin", "editor"}
	ctx = AddUserRolesToCtx(ctx, roles)

	// Get user roles from context
	retrievedRoles, ok := GetUserRolesFromCtx(ctx)
	if ok {
		fmt.Printf("User roles from context: %v\n", retrievedRoles)
	} else {
		fmt.Println("User roles not found in context")
	}

	// Create a new context with a different user ID
	newCtx := AddUserIDToCtx(context.Background(), "user456")

	// Check if the new context has the same user ID (should be different)
	newUserID, _ := GetUserIDFromCtx(newCtx)
	fmt.Printf("New user ID: %s\n", newUserID)

	// Check if the new context has roles (should not have roles)
	_, hasRoles := GetUserRolesFromCtx(newCtx)
	fmt.Printf("New context has roles: %v\n", hasRoles)
}