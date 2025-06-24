# Package Dependencies Guide

## Overview

This guide provides recommendations for managing package dependencies in the ServiceLib library to reduce coupling and improve maintainability. It outlines best practices for package organization, interface design, and dependency management.

## Current Package Structure

ServiceLib is organized into a layered architecture:

1. **Core Layer**: Fundamental utilities like context, errors, and logging
2. **Infrastructure Layer**: Database, configuration, and telemetry
3. **Service Layer**: Authentication, health checks, and middleware
4. **Application Layer**: Validation and application-specific utilities

## Package Dependency Principles

To reduce coupling between packages, follow these principles:

1. **Depend on Abstractions, Not Implementations**: Packages should depend on interfaces rather than concrete implementations.
2. **Interface Segregation**: Define small, focused interfaces that clients can implement without unnecessary dependencies.
3. **Dependency Inversion**: High-level modules should not depend on low-level modules. Both should depend on abstractions.
4. **Package Cohesion**: Keep related functionality together in the same package.
5. **Acyclic Dependencies**: Avoid circular dependencies between packages.

## Interface Packages

To further reduce coupling, we recommend creating dedicated interface packages for key components:

1. **Core Interfaces**: Move interfaces from core packages to dedicated interface packages.
2. **Domain Interfaces**: Define interfaces for domain services and repositories.
3. **Infrastructure Interfaces**: Define interfaces for infrastructure services.

### Example Structure

```
servicelib/
├── auth/
│   ├── interfaces/     # Auth interfaces
│   ├── jwt/            # JWT implementation
│   └── oidc/           # OIDC implementation
├── db/
│   ├── interfaces/     # Database interfaces
│   ├── sql/            # SQL implementation
│   └── mongo/          # MongoDB implementation
├── logging/
│   ├── interfaces/     # Logging interfaces
│   └── zap/            # Zap implementation
```

## Implementation Guidelines

### 1. Create Interface Packages

For each major component, create a dedicated interface package:

```go
// auth/interfaces/authenticator.go
package interfaces

import "context"

// Authenticator defines the interface for authentication operations
type Authenticator interface {
    Authenticate(ctx context.Context, token string) (bool, error)
    GenerateToken(ctx context.Context, userID string, roles []string) (string, error)
}
```

### 2. Implement Interfaces in Concrete Packages

Implement the interfaces in concrete packages:

```go
// auth/jwt/jwt_authenticator.go
package jwt

import (
    "context"
    "github.com/abitofhelp/servicelib/auth/interfaces"
)

// Ensure JWTAuthenticator implements the Authenticator interface
var _ interfaces.Authenticator = (*JWTAuthenticator)(nil)

// JWTAuthenticator implements the Authenticator interface using JWT
type JWTAuthenticator struct {
    // implementation details
}

// Authenticate implements the Authenticator interface
func (a *JWTAuthenticator) Authenticate(ctx context.Context, token string) (bool, error) {
    // implementation
}

// GenerateToken implements the Authenticator interface
func (a *JWTAuthenticator) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
    // implementation
}
```

### 3. Depend on Interfaces, Not Implementations

Other packages should depend on the interfaces, not the concrete implementations:

```go
// middleware/auth_middleware.go
package middleware

import (
    "context"
    "net/http"
    
    "github.com/abitofhelp/servicelib/auth/interfaces"
)

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
    authenticator interfaces.Authenticator
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authenticator interfaces.Authenticator) *AuthMiddleware {
    return &AuthMiddleware{
        authenticator: authenticator,
    }
}

// Middleware returns an http.Handler middleware function
func (m *AuthMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Use the authenticator interface
            // implementation
        })
    }
}
```

## Package Dependency Documentation

For each package, document its dependencies in the README.md file:

```markdown
# Auth Package

## Dependencies

This package depends on the following interfaces:
- `logging.interfaces.Logger`: For logging
- `config.interfaces.Config`: For configuration

This package provides the following interfaces:
- `auth.interfaces.Authenticator`: For authentication operations
- `auth.interfaces.Authorizer`: For authorization operations
```

## Conclusion

By following these principles and guidelines, we can reduce coupling between packages, improve maintainability, and make the codebase more flexible and testable. The use of interface packages and dependency inversion will allow for easier changes and extensions in the future.