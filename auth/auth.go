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
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/abitofhelp/servicelib/validation"
	"go.opentelemetry.io/otel/attribute"
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

		// Remote validation configuration
		Remote struct {
			// Enabled determines if remote validation should be used
			Enabled bool

			// ValidationURL is the URL of the remote validation endpoint
			ValidationURL string

			// ClientID is the client ID for the remote validation service
			ClientID string

			// ClientSecret is the client secret for the remote validation service
			ClientSecret string

			// Timeout is the timeout for remote validation operations
			Timeout time.Duration
		}
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

	// JWT Remote validation defaults
	config.JWT.Remote.Enabled = false
	config.JWT.Remote.Timeout = 10 * time.Second

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

// ValidateConfig validates the configuration for the auth module.
func ValidateConfig(config Config) *validation.ValidationResult {
	result := validation.NewValidationResult()

	// Validate JWT configuration
	validation.Required(config.JWT.SecretKey, "JWT.SecretKey", result)
	if config.JWT.TokenDuration <= 0 {
		result.AddError("must be positive", "JWT.TokenDuration")
	}
	validation.Required(config.JWT.Issuer, "JWT.Issuer", result)

	// Validate JWT Remote configuration if enabled
	if config.JWT.Remote.Enabled {
		validation.Required(config.JWT.Remote.ValidationURL, "JWT.Remote.ValidationURL", result)
		if config.JWT.Remote.Timeout <= 0 {
			result.AddError("must be positive", "JWT.Remote.Timeout")
		}
	}

	// Validate OIDC configuration if provided
	if config.OIDC.IssuerURL != "" || config.OIDC.ClientID != "" {
		validation.Required(config.OIDC.IssuerURL, "OIDC.IssuerURL", result)
		validation.Required(config.OIDC.ClientID, "OIDC.ClientID", result)
		if config.OIDC.Timeout <= 0 {
			result.AddError("must be positive", "OIDC.Timeout")
		}
	}

	// Validate Service configuration
	validation.Required(config.Service.AdminRoleName, "Service.AdminRoleName", result)
	validation.Required(config.Service.ReadOnlyRoleName, "Service.ReadOnlyRoleName", result)
	if len(config.Service.ReadOperationPrefixes) == 0 {
		result.AddError("must not be empty", "Service.ReadOperationPrefixes")
	}

	return result
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
	logger logging.Logger

	// metrics for tracking authentication and authorization operations
	metrics struct {
		tokenGenerations    int64
		tokenValidations    int64
		authorizationChecks int64
	}
}

// New creates a new Auth instance with the provided configuration and logger.
func New(ctx context.Context, config Config, logger *zap.Logger) (*Auth, error) {
	// Create a span for the initialization process
	ctx, span := telemetry.StartSpan(ctx, "auth.New")
	defer span.End()

	if logger == nil {
		logger = zap.NewNop()
	}

	// Create a context logger
	ctxLogger := logging.NewContextLogger(logger)

	// Add attributes to the span
	telemetry.AddSpanAttributes(ctx,
		attribute.String("jwt.issuer", config.JWT.Issuer),
		attribute.Bool("jwt.remote.enabled", config.JWT.Remote.Enabled),
		attribute.Bool("middleware.require_auth", config.Middleware.RequireAuth),
	)

	// Validate configuration
	validationResult := ValidateConfig(config)
	if !validationResult.IsValid() {
		err := errors.WithMessage(errors.ErrInvalidConfig, validationResult.Error().Error())
		telemetry.RecordErrorSpan(ctx, err)
		return nil, err
	}

	// Create JWT service
	jwtService, jwtErr := jwt.NewService(jwt.Config{
		SecretKey:     config.JWT.SecretKey,
		TokenDuration: config.JWT.TokenDuration,
		Issuer:        config.JWT.Issuer,
	}, logger)
	if jwtErr != nil {
		return nil, errors.WithOp(jwtErr, "auth.New")
	}

	// Add remote validator if enabled
	if config.JWT.Remote.Enabled {
		jwtService.WithRemoteValidator(jwt.RemoteConfig{
			ValidationURL: config.JWT.Remote.ValidationURL,
			ClientID:      config.JWT.Remote.ClientID,
			ClientSecret:  config.JWT.Remote.ClientSecret,
			Timeout:       config.JWT.Remote.Timeout,
		})

		ctxLogger.Info(ctx, "JWT remote validation enabled",
			zap.String("validation_url", config.JWT.Remote.ValidationURL),
			zap.Duration("timeout", config.JWT.Remote.Timeout))
	}

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
		logger:         ctxLogger,
	}, nil
}

// Middleware returns the HTTP middleware for authentication.
func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return a.authMiddleware.Handler
}

// GenerateToken generates a new JWT token for a user with the specified roles, scopes, and resources.
func (a *Auth) GenerateToken(ctx context.Context, userID string, roles []string, scopes []string, resources []string) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.GenerateToken")
	defer span.End()

	// Record metrics
	a.metrics.tokenGenerations++
	telemetry.RecordErrorMetric(ctx, "auth", "token_generation")

	// Add span attributes
	telemetry.AddSpanAttributes(ctx,
		attribute.String("user.id", userID),
		attribute.StringSlice("user.roles", roles),
		attribute.StringSlice("user.scopes", scopes),
		attribute.StringSlice("user.resources", resources),
	)

	// Generate token
	start := time.Now()
	token, err := a.jwtService.GenerateToken(ctx, userID, roles, scopes, resources)
	duration := time.Since(start)

	// Record result
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		telemetry.RecordErrorMetric(ctx, "auth", "token_generation_error")
	}

	// Add duration attribute
	telemetry.AddSpanAttributes(ctx,
		attribute.Float64("duration_ms", float64(duration.Milliseconds())),
	)

	return token, err
}

// ValidateToken validates a JWT token and returns the claims if valid.
func (a *Auth) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.ValidateToken")
	defer span.End()

	// Record metrics
	a.metrics.tokenValidations++
	telemetry.RecordErrorMetric(ctx, "auth", "token_validation")

	// Validate token
	start := time.Now()
	claims, err := a.jwtService.ValidateToken(ctx, tokenString)
	duration := time.Since(start)

	// Record result
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		telemetry.RecordErrorMetric(ctx, "auth", "token_validation_error")
	} else if claims != nil {
		telemetry.AddSpanAttributes(ctx,
			attribute.String("user.id", claims.UserID),
			attribute.StringSlice("user.roles", claims.Roles),
		)
	}

	// Add duration attribute
	telemetry.AddSpanAttributes(ctx,
		attribute.Float64("duration_ms", float64(duration.Milliseconds())),
	)

	return claims, err
}

// IsAuthorized checks if the user is authorized to perform the operation.
func (a *Auth) IsAuthorized(ctx context.Context, operation string) (bool, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.IsAuthorized")
	defer span.End()

	// Record metrics
	a.metrics.authorizationChecks++
	telemetry.RecordErrorMetric(ctx, "auth", "authorization_check")

	// Add span attributes
	telemetry.AddSpanAttributes(ctx,
		attribute.String("operation", operation),
	)

	// Check authorization
	start := time.Now()
	authorized, err := a.authService.IsAuthorized(ctx, operation)
	duration := time.Since(start)

	// Record result
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		telemetry.RecordErrorMetric(ctx, "auth", "authorization_check_error")
	}
	telemetry.AddSpanAttributes(ctx,
		attribute.Bool("authorized", authorized),
		attribute.Float64("duration_ms", float64(duration.Milliseconds())),
	)

	return authorized, err
}

// IsAdmin checks if the user has admin role.
func (a *Auth) IsAdmin(ctx context.Context) (bool, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.IsAdmin")
	defer span.End()

	isAdmin, err := a.authService.IsAdmin(ctx)
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
	}
	telemetry.AddSpanAttributes(ctx,
		attribute.Bool("is_admin", isAdmin),
	)
	return isAdmin, err
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
