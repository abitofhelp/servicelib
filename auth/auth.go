// Copyright (c) 2025 A Bit of Help, Inc.

// Package auth provides authentication and authorization functionality.
// It includes JWT token handling, OIDC integration, HTTP middleware, and role-based access control.
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/abitofhelp/servicelib/auth/service"
	"go.uber.org/zap"
)

// Config holds the configuration for the auth module.
type Config struct {
	// JWT configuration
	JWT struct {
		// SecretKey is the key used to sign and verify JWT tokens
		SecretKey string

		// TokenDuration is the validity period for generated tokens
		TokenDuration time.Duration

		// Issuer identifies the entity that issued the token
		Issuer string
	}

	// OIDC configuration
	OIDC struct {
		// IssuerURL is the URL of the OIDC provider
		IssuerURL string

		// ClientID is the client ID for the OIDC provider
		ClientID string

		// ClientSecret is the client secret for the OIDC provider
		ClientSecret string

		// RedirectURL is the redirect URL for the OIDC provider
		RedirectURL string

		// Scopes are the OAuth2 scopes to request
		Scopes []string

		// Timeout is the timeout for OIDC operations
		Timeout time.Duration
	}

	// Middleware configuration
	Middleware struct {
		// SkipPaths are paths that should skip authentication
		SkipPaths []string

		// RequireAuth determines if authentication is required for all requests
		RequireAuth bool
	}

	// Service configuration
	Service struct {
		// AdminRoleName is the name of the admin role
		AdminRoleName string

		// ReadOnlyRoleName is the name of the read-only role
		ReadOnlyRoleName string

		// ReadOperationPrefixes are prefixes for read-only operations
		ReadOperationPrefixes []string
	}
}

// DefaultConfig returns the default configuration for the auth module.
func DefaultConfig() Config {
	config := Config{}

	// JWT defaults
	config.JWT.TokenDuration = 24 * time.Hour
	config.JWT.Issuer = "auth"

	// OIDC defaults
	config.OIDC.Timeout = 10 * time.Second
	config.OIDC.Scopes = []string{"openid", "profile", "email"}

	// Middleware defaults
	config.Middleware.RequireAuth = true

	// Service defaults
	config.Service.AdminRoleName = "admin"
	config.Service.ReadOnlyRoleName = "authuser"
	config.Service.ReadOperationPrefixes = []string{
		"read:",
		"list:",
		"get:",
		"find:",
		"query:",
		"count:",
	}

	return config
}

// Auth provides authentication and authorization functionality.
type Auth struct {
	// jwtService is the JWT service for token handling
	jwtService *jwt.Service

	// oidcService is the OIDC service for token validation
	oidcService *oidc.Service

	// authMiddleware is the middleware for authentication
	authMiddleware *middleware.Middleware

	// authService is the service for authorization
	authService *service.Service

	// logger is used for logging operations and errors
	logger *zap.Logger
}

// New creates a new Auth instance with the provided configuration and logger.
func New(ctx context.Context, config Config, logger *zap.Logger) (*Auth, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Validate JWT configuration
	if config.JWT.SecretKey == "" {
		return nil, errors.WithMessage(errors.ErrInvalidConfig, "JWT secret key cannot be empty")
	}

	// Create JWT service
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     config.JWT.SecretKey,
		TokenDuration: config.JWT.TokenDuration,
		Issuer:        config.JWT.Issuer,
	}, logger)

	// Create OIDC service if configured
	var oidcService *oidc.Service
	var err error
	if config.OIDC.IssuerURL != "" && config.OIDC.ClientID != "" {
		oidcService, err = oidc.NewService(ctx, oidc.Config{
			IssuerURL:     config.OIDC.IssuerURL,
			ClientID:      config.OIDC.ClientID,
			ClientSecret:  config.OIDC.ClientSecret,
			RedirectURL:   config.OIDC.RedirectURL,
			Scopes:        config.OIDC.Scopes,
			AdminRoleName: config.Service.AdminRoleName,
			Timeout:       config.OIDC.Timeout,
		}, logger)
		if err != nil {
			return nil, errors.WithOp(err, "auth.New")
		}
	}

	// Create middleware
	var authMiddleware *middleware.Middleware
	if oidcService != nil {
		authMiddleware = middleware.NewMiddlewareWithOIDC(jwtService, oidcService, middleware.Config{
			SkipPaths:   config.Middleware.SkipPaths,
			RequireAuth: config.Middleware.RequireAuth,
		}, logger)
	} else {
		authMiddleware = middleware.NewMiddleware(jwtService, middleware.Config{
			SkipPaths:   config.Middleware.SkipPaths,
			RequireAuth: config.Middleware.RequireAuth,
		}, logger)
	}

	// Create authorization service
	authService := service.NewService(service.Config{
		AdminRoleName:         config.Service.AdminRoleName,
		ReadOnlyRoleName:      config.Service.ReadOnlyRoleName,
		ReadOperationPrefixes: config.Service.ReadOperationPrefixes,
	}, logger)

	return &Auth{
		jwtService:     jwtService,
		oidcService:    oidcService,
		authMiddleware: authMiddleware,
		authService:    authService,
		logger:         logger,
	}, nil
}

// Middleware returns the HTTP middleware for authentication.
func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return a.authMiddleware.Handler
}

// GenerateToken generates a new JWT token for a user with the specified roles.
func (a *Auth) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
	return a.jwtService.GenerateToken(ctx, userID, roles)
}

// ValidateToken validates a JWT token and returns the claims if valid.
func (a *Auth) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	return a.jwtService.ValidateToken(ctx, tokenString)
}

// IsAuthorized checks if the user is authorized to perform the operation.
func (a *Auth) IsAuthorized(ctx context.Context, operation string) (bool, error) {
	return a.authService.IsAuthorized(ctx, operation)
}

// IsAdmin checks if the user has admin role.
func (a *Auth) IsAdmin(ctx context.Context) (bool, error) {
	return a.authService.IsAdmin(ctx)
}

// HasRole checks if the user has a specific role.
func (a *Auth) HasRole(ctx context.Context, role string) (bool, error) {
	return a.authService.HasRole(ctx, role)
}

// GetUserID retrieves the user ID from the context.
func (a *Auth) GetUserID(ctx context.Context) (string, error) {
	return a.authService.GetUserID(ctx)
}

// GetUserRoles retrieves the user roles from the context.
func (a *Auth) GetUserRoles(ctx context.Context) ([]string, error) {
	return a.authService.GetUserRoles(ctx)
}

// WithUserID returns a new context with the user ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return middleware.WithUserID(ctx, userID)
}

// WithUserRoles returns a new context with the user roles.
func WithUserRoles(ctx context.Context, roles []string) context.Context {
	return middleware.WithUserRoles(ctx, roles)
}

// GetUserIDFromContext retrieves the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	return middleware.GetUserID(ctx)
}

// GetUserRolesFromContext retrieves the user roles from the context.
func GetUserRolesFromContext(ctx context.Context) ([]string, bool) {
	return middleware.GetUserRoles(ctx)
}

// IsAuthenticated checks if the user is authenticated.
func IsAuthenticated(ctx context.Context) bool {
	return middleware.IsAuthenticated(ctx)
}
