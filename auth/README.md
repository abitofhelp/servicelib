# Auth Module

The Auth Module provides comprehensive authentication and authorization functionality for Go applications. It includes JWT token handling, OIDC integration, HTTP middleware, and role-based access control.

## Features

- **JWT Token Handling**: Generate and validate JWT tokens
- **OIDC Integration**: Validate tokens from OpenID Connect providers
- **HTTP Middleware**: Authenticate HTTP requests
- **Role-Based Access Control**: Control access to resources based on user roles
- **Comprehensive Error Handling**: Context-aware error handling with detailed error information
- **Tracing Integration**: OpenTelemetry tracing for all operations
- **Logging Integration**: Structured logging with zap

## Installation

```bash
go get github.com/abitofhelp/servicelib/auth
```

## Quick Start

```go
package main

import (
    "context"
    "net/http"
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

    // Create a configuration
    config := auth.DefaultConfig()
    config.JWT.SecretKey = "your-secret-key"

    // Create an auth instance
    authInstance, err := auth.New(ctx, config, logger)
    if err != nil {
        logger.Fatal("Failed to create auth instance", zap.Error(err))
    }

    // Create an HTTP handler
    http.Handle("/", authInstance.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if the user is authorized to perform an operation
        authorized, err := authInstance.IsAuthorized(r.Context(), "read:resource")
        if err != nil {
            http.Error(w, "Authorization error", http.StatusInternalServerError)
            return
        }

        if !authorized {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        // Get the user ID
        userID, err := authInstance.GetUserID(r.Context())
        if err != nil {
            http.Error(w, "User ID not found", http.StatusInternalServerError)
            return
        }

        w.Write([]byte("Hello, " + userID))
    })))

    // Start the server
    http.ListenAndServe(":8080", nil)
}
```

## Configuration

The auth module can be configured using the `Config` struct:

```go
config := auth.Config{}

// JWT configuration
config.JWT.SecretKey = "your-secret-key"
config.JWT.TokenDuration = 24 * time.Hour
config.JWT.Issuer = "your-issuer"

// OIDC configuration
config.OIDC.IssuerURL = "https://your-oidc-provider.com"
config.OIDC.ClientID = "your-client-id"
config.OIDC.ClientSecret = "your-client-secret"
config.OIDC.RedirectURL = "https://your-app.com/callback"
config.OIDC.Scopes = []string{"openid", "profile", "email"}
config.OIDC.Timeout = 10 * time.Second

// Middleware configuration
config.Middleware.SkipPaths = []string{"/public", "/health"}
config.Middleware.RequireAuth = true

// Service configuration
config.Service.AdminRoleName = "admin"
config.Service.ReadOnlyRoleName = "reader"
config.Service.ReadOperationPrefixes = []string{"read:", "list:", "get:"}
```

## API Documentation

### Auth

The `Auth` struct is the main entry point for the auth module. It provides methods for authentication and authorization.

#### Creating an Auth Instance

```go
authInstance, err := auth.New(ctx, config, logger)
```

#### Middleware

```go
// Get the middleware function
middleware := authInstance.Middleware()

// Use the middleware with an HTTP handler
http.Handle("/", middleware(yourHandler))
```

#### Token Handling

```go
// Generate a token
token, err := authInstance.GenerateToken(ctx, "user123", []string{"admin"})

// Validate a token
claims, err := authInstance.ValidateToken(ctx, token)
```

#### Authorization

```go
// Check if the user is authorized to perform an operation
authorized, err := authInstance.IsAuthorized(ctx, "read:resource")

// Check if the user has admin role
isAdmin, err := authInstance.IsAdmin(ctx)

// Check if the user has a specific role
hasRole, err := authInstance.HasRole(ctx, "editor")
```

#### User Information

```go
// Get the user ID from the context
userID, err := authInstance.GetUserID(ctx)

// Get the user roles from the context
roles, err := authInstance.GetUserRoles(ctx)
```

### Context Utilities

The auth module provides utilities for working with context:

```go
// Add user ID to context
ctx = auth.WithUserID(ctx, "user123")

// Add user roles to context
ctx = auth.WithUserRoles(ctx, []string{"admin", "editor"})

// Get user ID from context
userID, ok := auth.GetUserIDFromContext(ctx)

// Get user roles from context
roles, ok := auth.GetUserRolesFromContext(ctx)

// Check if the user is authenticated
isAuthenticated := auth.IsAuthenticated(ctx)
```

## Error Handling

The auth module provides comprehensive error handling with context-aware errors:

```go
import "github.com/abitofhelp/servicelib/auth/errors"

// Check for specific error types
if errors.Is(err, errors.ErrInvalidToken) {
    // Handle invalid token error
}

// Get context information from an error
if value, ok := errors.GetContext(err, "key"); ok {
    // Use the context value
}

// Get the operation that caused the error
if op, ok := errors.GetOp(err); ok {
    // Use the operation
}

// Get the error message
if message, ok := errors.GetMessage(err); ok {
    // Use the message
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.