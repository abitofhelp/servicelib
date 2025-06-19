// Copyright (c) 2025 A Bit of Help, Inc.

// Package middleware provides HTTP middleware for authentication.
// It extracts and validates tokens from HTTP requests and adds user information to the request context.
package middleware

import (
	"context"
	stderrors "errors" // Standard errors package with alias
	"net/http"
	"strconv"
	"strings"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// contextKey is a private type for context keys
type contextKey int

const (
	// userIDKey is the context key for the user ID
	userIDKey contextKey = iota

	// userRolesKey is the context key for the user roles
	userRolesKey

	// userScopesKey is the context key for the user scopes
	userScopesKey

	// userResourcesKey is the context key for the user resources
	userResourcesKey
)

// Config holds the configuration for the authentication middleware.
type Config struct {
	// SkipPaths are paths that should skip authentication
	SkipPaths []string

	// RequireAuth determines if authentication is required for all requests
	RequireAuth bool
}

// Middleware is a middleware for handling authentication.
type Middleware struct {
	// jwtService is the JWT service for token validation
	jwtService *jwt.Service

	// oidcService is the OIDC service for token validation
	oidcService *oidc.Service

	// config is the middleware configuration
	config Config

	// logger is used for logging middleware operations and errors
	logger *zap.Logger

	// tracer is used for tracing middleware operations
	tracer trace.Tracer
}

// NewMiddleware creates a new authentication middleware with JWT support.
func NewMiddleware(jwtService *jwt.Service, config Config, logger *zap.Logger) *Middleware {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Middleware{
		jwtService: jwtService,
		config:     config,
		logger:     logger,
		tracer:     otel.Tracer("auth.middleware"),
	}
}

// NewMiddlewareWithOIDC creates a new authentication middleware with both JWT and OIDC support.
func NewMiddlewareWithOIDC(jwtService *jwt.Service, oidcService *oidc.Service, config Config, logger *zap.Logger) *Middleware {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Middleware{
		jwtService:  jwtService,
		oidcService: oidcService,
		config:      config,
		logger:      logger,
		tracer:      otel.Tracer("auth.middleware"),
	}
}

// Handler is the HTTP middleware function.
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := m.tracer.Start(r.Context(), "middleware.Handler")
		defer span.End()

		// Check if the path should skip authentication
		for _, path := range m.config.SkipPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				m.logger.Debug("Skipping authentication for path", zap.String("path", r.URL.Path))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if m.config.RequireAuth {
				err := errors.WithOp(errors.ErrMissingToken, "middleware.Handler")
				err = errors.WithMessage(err, "authorization header is required")
				m.logger.Debug("No Authorization header provided")
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// No token provided, continue as unauthenticated
			m.logger.Debug("No Authorization header provided, continuing as unauthenticated")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Extract token from header
		tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			err = errors.WithOp(err, "middleware.Handler")
			m.logger.Debug("Invalid Authorization header format", zap.Error(err))
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		span.SetAttributes(attribute.String("token.length", strconv.Itoa(len(tokenString))))

		var claims *jwt.Claims

		// Try OIDC validation first if available
		if m.oidcService != nil {
			claims, err = m.oidcService.ValidateToken(ctx, tokenString)
			if err != nil {
				m.logger.Debug("OIDC validation failed, trying JWT", zap.Error(err))
				// Fall back to JWT validation
				claims, err = m.jwtService.ValidateToken(ctx, tokenString)
				if err != nil {
					m.handleAuthError(w, err)
					return
				}
			}
		} else {
			// Use JWT validation
			claims, err = m.jwtService.ValidateToken(ctx, tokenString)
			if err != nil {
				m.handleAuthError(w, err)
				return
			}
		}

		// Add user info to context
		ctx = WithUserID(ctx, claims.UserID)
		ctx = WithUserRoles(ctx, claims.Roles)
		ctx = WithUserScopes(ctx, claims.Scopes)
		ctx = WithUserResources(ctx, claims.Resources)

		span.SetAttributes(
			attribute.String("user.id", claims.UserID),
			attribute.StringSlice("user.roles", claims.Roles),
			attribute.StringSlice("user.scopes", claims.Scopes),
			attribute.StringSlice("user.resources", claims.Resources),
		)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleAuthError handles authentication errors.
func (m *Middleware) handleAuthError(w http.ResponseWriter, err error) {
	status := http.StatusUnauthorized
	message := "Invalid token"

	if stderrors.Is(err, errors.ErrMissingToken) {
		message = "Authorization required"
	} else if stderrors.Is(err, errors.ErrExpiredToken) {
		message = "Token expired"
	} else if stderrors.Is(err, errors.ErrInvalidSignature) {
		message = "Invalid token signature"
	} else if stderrors.Is(err, errors.ErrInvalidClaims) {
		message = "Invalid token claims"
	}

	m.logger.Debug("Authentication error", zap.Error(err), zap.String("message", message))
	http.Error(w, message, status)
}

// WithUserID returns a new context with the user ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithUserRoles returns a new context with the user roles.
func WithUserRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, userRolesKey, roles)
}

// WithUserScopes returns a new context with the user scopes.
func WithUserScopes(ctx context.Context, scopes []string) context.Context {
	return context.WithValue(ctx, userScopesKey, scopes)
}

// WithUserResources returns a new context with the user resources.
func WithUserResources(ctx context.Context, resources []string) context.Context {
	return context.WithValue(ctx, userResourcesKey, resources)
}

// GetUserID retrieves the user ID from the context.
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// GetUserRoles retrieves the user roles from the context.
func GetUserRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(userRolesKey).([]string)
	return roles, ok
}

// GetUserScopes retrieves the user scopes from the context.
func GetUserScopes(ctx context.Context) ([]string, bool) {
	scopes, ok := ctx.Value(userScopesKey).([]string)
	return scopes, ok
}

// GetUserResources retrieves the user resources from the context.
func GetUserResources(ctx context.Context) ([]string, bool) {
	resources, ok := ctx.Value(userResourcesKey).([]string)
	return resources, ok
}

// IsAuthenticated checks if the user is authenticated.
func IsAuthenticated(ctx context.Context) bool {
	_, ok := GetUserID(ctx)
	return ok
}

// HasRole checks if the user has a specific role.
func HasRole(ctx context.Context, role string) bool {
	roles, ok := GetUserRoles(ctx)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

// IsAuthorized checks if the user is authorized to perform a specific action based on their roles.
// It takes a list of allowed roles and returns true if the user has at least one of them.
func IsAuthorized(ctx context.Context, allowedRoles []string) bool {
	// If no roles are specified, deny access
	if len(allowedRoles) == 0 {
		return false
	}

	// Get user roles from context
	userRoles, ok := GetUserRoles(ctx)
	if !ok {
		return false
	}

	// Check if the user has any of the allowed roles
	for _, allowedRole := range allowedRoles {
		for _, userRole := range userRoles {
			if userRole == allowedRole {
				return true
			}
		}
	}

	return false
}

// HasScope checks if the user has a specific scope.
func HasScope(ctx context.Context, scope string) bool {
	scopes, ok := GetUserScopes(ctx)
	if !ok {
		return false
	}

	for _, s := range scopes {
		if s == scope {
			return true
		}
	}

	return false
}

// HasResource checks if the user has access to a specific resource.
func HasResource(ctx context.Context, resource string) bool {
	resources, ok := GetUserResources(ctx)
	if !ok {
		return false
	}

	for _, r := range resources {
		if r == resource {
			return true
		}
	}

	return false
}

// IsAuthorizedWithScopes checks if the user is authorized to perform a specific action based on their roles and scopes.
// It takes a list of allowed roles, required scopes, and a resource, and returns true if the user has at least one of the allowed roles
// and all of the required scopes for the specified resource.
func IsAuthorizedWithScopes(ctx context.Context, allowedRoles []string, requiredScopes []string, resource string) bool {
	// First check if the user has any of the allowed roles
	if !IsAuthorized(ctx, allowedRoles) {
		return false
	}

	// If no scopes are required, the role check is sufficient
	if len(requiredScopes) == 0 {
		return true
	}

	// Get user scopes from context
	userScopes, ok := GetUserScopes(ctx)
	if !ok {
		return false
	}

	// Get user resources from context
	userResources, ok := GetUserResources(ctx)
	if !ok {
		return false
	}

	// Check if the user has access to the required resource
	hasResource := false
	for _, r := range userResources {
		if r == resource {
			hasResource = true
			break
		}
	}

	if !hasResource {
		return false
	}

	// Check if the user has all of the required scopes
	for _, requiredScope := range requiredScopes {
		hasScope := false
		for _, userScope := range userScopes {
			if userScope == requiredScope {
				hasScope = true
				break
			}
		}
		if !hasScope {
			return false
		}
	}

	return true
}
