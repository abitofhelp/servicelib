// Copyright (c) 2025 A Bit of Help, Inc.

// Package auth provides comprehensive authentication and authorization functionality for services.
//
// This package implements a complete authentication and authorization system that includes:
//   - JWT token generation and validation
//   - OpenID Connect (OIDC) integration
//   - HTTP middleware for securing API endpoints
//   - Role-based access control (RBAC)
//   - Context utilities for user information management
//
// The auth package is designed to be flexible and configurable, supporting both
// local JWT-based authentication and integration with external identity providers
// through OIDC. It provides a unified interface for authentication and authorization
// operations, regardless of the underlying implementation.
//
// Key components:
//   - Auth: The main service that provides authentication and authorization functionality
//   - Config: Configuration for JWT, OIDC, middleware, and authorization settings
//   - Middleware: HTTP middleware for authenticating requests
//   - JWT subpackage: Handles JWT token generation and validation
//   - OIDC subpackage: Integrates with OpenID Connect providers
//   - Service subpackage: Implements authorization logic
//
// Example usage:
//
//	// Create an auth service
//	config := auth.DefaultConfig()
//	config.JWT.SecretKey = "your-secret-key"
//	config.JWT.Issuer = "your-service"
//	
//	authService, err := auth.New(ctx, config, logger)
//	if err != nil {
//	    log.Fatalf("Failed to create auth service: %v", err)
//	}
//	
//	// Use the middleware to protect routes
//	router := http.NewServeMux()
//	router.Handle("/api/", authService.Middleware()(apiHandler))
//	
//	// Generate a token for a user
//	token, err := authService.GenerateToken(ctx, "user123", []string{"admin"}, []string{"read:users"}, []string{})
//	if err != nil {
//	    log.Fatalf("Failed to generate token: %v", err)
//	}
//	
//	// Check if a user is authorized for an operation
//	authorized, err := authService.IsAuthorized(ctx, "read:users")
//	if err != nil {
//	    log.Fatalf("Failed to check authorization: %v", err)
//	}
//	
//	// Get user information from context
//	userID, err := authService.GetUserID(ctx)
//	if err != nil {
//	    log.Fatalf("Failed to get user ID: %v", err)
//	}
//
// The auth package is designed to be used as a dependency by other packages in the application,
// providing a consistent authentication and authorization interface throughout the codebase.
package auth