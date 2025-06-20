// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating authorization concepts
package main

import (
	"context"
	"fmt"
)

// This is a simplified example to demonstrate authorization concepts
// In a real application, you would use the auth module's methods

// AuthService simulates an auth service with authorization capabilities
type AuthService struct {
	userRoles map[string][]string
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		userRoles: map[string][]string{
			"user123": {"admin", "editor"},
			"user456": {"reader"},
		},
	}
}

// IsAuthorized checks if the user is authorized to perform an operation
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

// IsAdmin checks if the user has admin role
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

// HasRole checks if the user has a specific role
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

// Helper functions

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// isReadOperation checks if an operation is read-only
func isReadOperation(operation string) bool {
	readPrefixes := []string{"read:", "get:", "list:"}
	for _, prefix := range readPrefixes {
		if len(operation) >= len(prefix) && operation[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// isEditOperation checks if an operation is an edit operation
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
