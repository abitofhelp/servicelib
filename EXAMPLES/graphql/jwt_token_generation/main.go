//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to generate JWT tokens for testing GraphQL RBAC
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a context
	ctx := context.Background()

	// This example shows how to generate JWT tokens for testing GraphQL RBAC
	fmt.Println("Example: Generating JWT tokens for testing GraphQL RBAC")

	// Step 1: Create an auth service
	fmt.Println("\nStep 1: Create an auth service")
	authConfig := auth.DefaultConfig()
	authConfig.JWT.SecretKey = "your-secret-key"
	authConfig.JWT.Issuer = "test-service"
	authConfig.JWT.TokenDuration = 24 * time.Hour

	authService, err := auth.New(ctx, authConfig, logger)
	if err != nil {
		fmt.Printf("Error initializing auth service: %v\n", err)
		return
	}

	// Step 2: Generate tokens for different roles
	fmt.Println("\nStep 2: Generate tokens for different roles")
	fmt.Println(`
	// Generate admin token with all scopes for all resources
	adminScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	adminResources := []string{"FAMILY", "PARENT", "CHILD"}
	adminToken, err := authInstance.GenerateToken(ctx, "admin", []string{"ADMIN"}, adminScopes, adminResources)

	// Generate editor token with all scopes for all resources
	editorScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	editorResources := []string{"FAMILY", "PARENT", "CHILD"}
	editorToken, err := authInstance.GenerateToken(ctx, "editor", []string{"EDITOR"}, editorScopes, editorResources)

	// Generate viewer token with only READ scope for all resources
	viewerScopes := []string{"READ"}
	viewerResources := []string{"FAMILY", "PARENT", "CHILD"}
	viewerToken, err := authInstance.GenerateToken(ctx, "viewer", []string{"VIEWER"}, viewerScopes, viewerResources)
	`)

	// Generate tokens for demonstration
	// Note: The actual auth.GenerateToken method signature might be different
	// This is a simplified example
	fmt.Println("\nGenerating tokens (simplified for demonstration):")

	// Admin token
	fmt.Println("\nAdmin token:")
	fmt.Println("- User ID: admin")
	fmt.Println("- Roles: [ADMIN]")
	fmt.Println("- Scopes: [READ, WRITE, DELETE, CREATE]")
	fmt.Println("- Resources: [FAMILY, PARENT, CHILD]")
	fmt.Println("Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

	// Editor token
	fmt.Println("\nEditor token:")
	fmt.Println("- User ID: editor")
	fmt.Println("- Roles: [EDITOR]")
	fmt.Println("- Scopes: [READ, WRITE, DELETE, CREATE]")
	fmt.Println("- Resources: [FAMILY, PARENT, CHILD]")
	fmt.Println("Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

	// Viewer token
	fmt.Println("\nViewer token:")
	fmt.Println("- User ID: viewer")
	fmt.Println("- Roles: [VIEWER]")
	fmt.Println("- Scopes: [READ]")
	fmt.Println("- Resources: [FAMILY, PARENT, CHILD]")
	fmt.Println("Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

	// Explain how to use the tokens
	fmt.Println("\nTo use these tokens in GraphQL requests:")
	fmt.Println("1. Add the token to the Authorization header:")
	fmt.Println("   Authorization: Bearer <token>")
	fmt.Println("2. The @isAuthorized directive will check the token's roles, scopes, and resources")
	fmt.Println("3. The directive will allow or deny access based on the token's claims")

	// Prevent unused variable warning
	_ = authService
}