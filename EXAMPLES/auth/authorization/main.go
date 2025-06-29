//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Package main demonstrates authorization concepts
//
// This example shows how to implement basic role-based authorization:
// - Checking if a user has permission to perform specific operations
// - Verifying user roles and privileges
// - Using context to pass user identity information
// - Implementing different authorization strategies based on operation types
//
// The example creates a simple AuthService that manages user roles and provides
// methods to check if users are authorized to perform different operations based
// on their assigned roles.
package main

import (
	"context"
	"fmt"
)

// AuthService simulates an auth service with authorization capabilities.
//
// This struct represents a simplified authorization service that manages user roles
// and provides methods to check if users are authorized to perform different operations.
// It stores user roles in a map where the key is the user ID and the value is a slice
// of role names assigned to that user.
//
// In a real application, you would use the auth module's methods instead of
// implementing this functionality yourself. This example is for demonstration purposes only.
type AuthService struct {
	userRoles map[string][]string // Maps user IDs to their assigned roles
}

// NewAuthService creates a new auth service with predefined user roles.
//
// This function initializes an AuthService with a set of predefined user roles:
// - user123: has both "admin" and "editor" roles
// - user456: has only the "reader" role
//
// These predefined roles allow the example to demonstrate different authorization
// scenarios without requiring user input or external configuration.
//
// Returns:
//   - A pointer to a new AuthService instance with predefined user roles
func NewAuthService() *AuthService {
	return &AuthService{
		userRoles: map[string][]string{
			"user123": {"admin", "editor"},
			"user456": {"reader"},
		},
	}
}

// IsAuthorized checks if the user is authorized to perform an operation.
//
// This method implements a role-based authorization strategy with the following rules:
// 1. Users with the "admin" role can perform any operation
// 2. Users with the "reader" or "editor" role can perform read operations
// 3. Users with the "editor" role can perform edit operations
//
// The method extracts the user ID from the context, retrieves the user's roles,
// and then applies the authorization rules based on the requested operation.
//
// Parameters:
//   - ctx: The context containing the user ID
//   - operation: The operation to check authorization for (e.g., "read:document", "edit:document")
//
// Returns:
//   - bool: true if the user is authorized, false otherwise
//   - error: An error if the user ID is not found in the context or if another error occurs
//
// Example usage:
//
//	authorized, err := authService.IsAuthorized(ctx, "read:document")
//	if err != nil {
//	    // Handle error
//	}
//	if authorized {
//	    // Proceed with the operation
//	}
func (a *AuthService) IsAuthorized(ctx context.Context, operation string) (bool, error) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return false, fmt.Errorf("user ID not found in context")
	}

	// Get user roles
	roles, ok := a.userRoles[userID]
	if !ok {
		return false, nil
	}

	// Check if user has admin role (admins can do anything)
	for _, role := range roles {
		if role == "admin" {
			return true, nil
		}
	}

	// Check if operation is read-only and user has reader role
	if isReadOperation(operation) {
		for _, role := range roles {
			if role == "reader" || role == "editor" {
				return true, nil
			}
		}
	}

	// Check if user has editor role for edit operations
	if isEditOperation(operation) {
		for _, role := range roles {
			if role == "editor" {
				return true, nil
			}
		}
	}

	return false, nil
}

// IsAdmin checks if the user has the admin role.
//
// This method provides a convenient way to check if a user has administrative
// privileges without having to specify an operation. It extracts the user ID
// from the context, retrieves the user's roles, and checks if the "admin" role
// is present.
//
// Parameters:
//   - ctx: The context containing the user ID
//
// Returns:
//   - bool: true if the user has the admin role, false otherwise
//   - error: An error if the user ID is not found in the context or if another error occurs
//
// Example usage:
//
//	isAdmin, err := authService.IsAdmin(ctx)
//	if err != nil {
//	    // Handle error
//	}
//	if isAdmin {
//	    // Show admin features
//	}
func (a *AuthService) IsAdmin(ctx context.Context) (bool, error) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return false, fmt.Errorf("user ID not found in context")
	}

	// Get user roles
	roles, ok := a.userRoles[userID]
	if !ok {
		return false, nil
	}

	// Check if user has admin role
	for _, role := range roles {
		if role == "admin" {
			return true, nil
		}
	}

	return false, nil
}

// HasRole checks if the user has a specific role.
//
// This method allows checking for any arbitrary role, not just predefined ones
// like "admin". It extracts the user ID from the context, retrieves the user's
// roles, and checks if the specified role is present in the user's role list.
//
// Parameters:
//   - ctx: The context containing the user ID
//   - role: The role name to check for (e.g., "admin", "editor", "reader")
//
// Returns:
//   - bool: true if the user has the specified role, false otherwise
//   - error: An error if the user ID is not found in the context or if another error occurs
//
// Example usage:
//
//	hasEditorRole, err := authService.HasRole(ctx, "editor")
//	if err != nil {
//	    // Handle error
//	}
//	if hasEditorRole {
//	    // Show editor features
//	}
func (a *AuthService) HasRole(ctx context.Context, role string) (bool, error) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return false, fmt.Errorf("user ID not found in context")
	}

	// Get user roles
	roles, ok := a.userRoles[userID]
	if !ok {
		return false, nil
	}

	// Check if user has the specified role
	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}

	return false, nil
}

// Helper functions for context management

// WithUserID adds a user ID to the context.
//
// This function creates a new context derived from the parent context with
// the user ID stored as a value. This is a common pattern for passing user
// identity information through the call stack without having to explicitly
// pass it as a parameter to every function.
//
// Parameters:
//   - ctx: The parent context
//   - userID: The user ID to store in the context
//
// Returns:
//   - A new context containing the user ID
//
// Example usage:
//
//	// Create a context with user ID
//	ctx := WithUserID(context.Background(), "user123")
//	
//	// Pass the context to functions that need the user ID
//	authorized, err := authService.IsAuthorized(ctx, "read:document")
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

// GetUserIDFromContext gets the user ID from the context.
//
// This function extracts the user ID that was previously stored in the context
// using WithUserID. It returns the user ID and a boolean indicating whether
// the user ID was found in the context.
//
// Parameters:
//   - ctx: The context to extract the user ID from
//
// Returns:
//   - string: The user ID if found
//   - bool: true if the user ID was found, false otherwise
//
// Example usage:
//
//	userID, ok := GetUserIDFromContext(ctx)
//	if !ok {
//	    return fmt.Errorf("user ID not found in context")
//	}
//	fmt.Printf("User ID: %s\n", userID)
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// isReadOperation checks if an operation is read-only.
//
// This function determines if an operation is read-only by checking if it starts
// with one of the predefined read operation prefixes: "read:", "get:", or "list:".
// Read operations typically don't modify data and have less strict authorization
// requirements.
//
// Parameters:
//   - operation: The operation string to check
//
// Returns:
//   - true if the operation is a read operation, false otherwise
func isReadOperation(operation string) bool {
	readPrefixes := []string{"read:", "get:", "list:"}
	for _, prefix := range readPrefixes {
		if len(operation) >= len(prefix) && operation[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// isEditOperation checks if an operation is an edit operation.
//
// This function determines if an operation is an edit operation by checking if it
// starts with one of the predefined edit operation prefixes: "edit:", "update:",
// "create:", or "delete:". Edit operations modify data and typically require
// higher levels of authorization.
//
// Parameters:
//   - operation: The operation string to check
//
// Returns:
//   - true if the operation is an edit operation, false otherwise
func isEditOperation(operation string) bool {
	editPrefixes := []string{"edit:", "update:", "create:", "delete:"}
	for _, prefix := range editPrefixes {
		if len(operation) >= len(prefix) && operation[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func main() {
	// Create an auth service
	auth := NewAuthService()

	// Create a context with user ID
	ctx := context.Background()

	// Test with admin user
	adminCtx := WithUserID(ctx, "user123")

	// Check if admin user is authorized to perform operations
	authorized, _ := auth.IsAuthorized(adminCtx, "read:document")
	fmt.Printf("Admin authorized for read:document: %v\n", authorized)

	authorized, _ = auth.IsAuthorized(adminCtx, "edit:document")
	fmt.Printf("Admin authorized for edit:document: %v\n", authorized)

	isAdmin, _ := auth.IsAdmin(adminCtx)
	fmt.Printf("Is admin: %v\n", isAdmin)

	hasRole, _ := auth.HasRole(adminCtx, "editor")
	fmt.Printf("Admin has editor role: %v\n", hasRole)

	// Test with reader user
	readerCtx := WithUserID(ctx, "user456")

	// Check if reader user is authorized to perform operations
	authorized, _ = auth.IsAuthorized(readerCtx, "read:document")
	fmt.Printf("Reader authorized for read:document: %v\n", authorized)

	authorized, _ = auth.IsAuthorized(readerCtx, "edit:document")
	fmt.Printf("Reader authorized for edit:document: %v\n", authorized)

	isAdmin, _ = auth.IsAdmin(readerCtx)
	fmt.Printf("Is admin: %v\n", isAdmin)

	hasRole, _ = auth.HasRole(readerCtx, "reader")
	fmt.Printf("Reader has reader role: %v\n", hasRole)
}
