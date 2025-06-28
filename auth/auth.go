// Copyright (c) 2025 A Bit of Help, Inc.

// Package auth provides authentication and authorization functionality.
// It includes JWT token handling, OIDC integration, HTTP middleware, and role-based access control.
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/abitofhelp/servicelib/auth/service"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/abitofhelp/servicelib/validation"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// Config holds the configuration for the auth module.
// It contains settings for JWT token handling, OIDC integration, HTTP middleware,
// and authorization service behavior.
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
// It sets reasonable default values for all configuration options, including:
//   - JWT token duration (24 hours)
//   - JWT issuer ("auth")
//   - Remote validation timeout (10 seconds)
//   - OIDC timeout (10 seconds)
//   - OIDC scopes (openid, profile, email)
//   - Admin role name ("admin")
//   - Read-only role name ("authuser")
//   - Read operation prefixes for authorization
//
// Returns:
//   - Config: A configuration struct with default values.
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
// It checks that all required fields are set and have valid values, including:
//   - JWT secret key and issuer
//   - JWT token duration (must be positive)
//   - Remote validation URL and timeout if remote validation is enabled
//   - OIDC issuer URL and client ID if OIDC is configured
//   - Admin role name and read-only role name
//   - Read operation prefixes (must not be empty)
//
// Parameters:
//   - config: The configuration to validate.
//
// Returns:
//   - *validation.ValidationResult: A validation result containing any validation errors.
//     If the result's IsValid() method returns true, the configuration is valid.
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
// It encapsulates JWT token handling, OIDC integration, HTTP middleware,
// and role-based access control in a unified interface. The Auth struct
// is the main entry point for all authentication and authorization operations.
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
// It initializes the JWT service, OIDC service (if configured), authentication middleware,
// and authorization service based on the provided configuration. The function validates
// the configuration before creating any services and returns an error if the configuration
// is invalid.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - config: The configuration for the auth module.
//   - logger: The logger to use for logging operations and errors. If nil, a no-op logger will be used.
//
// Returns:
//   - *Auth: A new Auth instance if successful.
//   - error: An error if the configuration is invalid or if any service initialization fails.
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
		err := errors.NewConfigurationError("invalid configuration", "auth", validationResult.Error().Error(), nil)
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
		return nil, errors.WrapWithOperation(jwtErr, errors.InternalErrorCode, "JWT service creation failed", "auth.New")
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
			return nil, errors.WrapWithOperation(err, errors.InternalErrorCode, "OIDC service creation failed", "auth.New")
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
// This middleware intercepts HTTP requests, extracts and validates authentication tokens,
// and adds user information to the request context. It can be configured to skip
// authentication for specific paths and to require authentication for all other paths.
//
// Returns:
//   - func(http.Handler) http.Handler: A middleware function that can be used with
//     standard HTTP handlers or routers to add authentication to routes.
func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return a.authMiddleware.Handler
}

// GenerateToken generates a new JWT token for a user with the specified roles, scopes, and resources.
// The token includes claims for the user ID, roles, scopes, and resources, and is signed
// with the configured secret key. The token has an expiration time based on the configured
// token duration.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - userID: The unique identifier of the user for whom the token is being generated.
//   - roles: The roles assigned to the user (e.g., "admin", "user").
//   - scopes: The OAuth2 scopes granted to the user (e.g., "read:users", "write:users").
//   - resources: The resources the user has access to.
//
// Returns:
//   - string: The generated JWT token as a string.
//   - error: An error if token generation fails.
func (a *Auth) GenerateToken(ctx context.Context, userID string, roles []string, scopes []string, resources []string) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.GenerateToken")
	defer span.End()

	// Record metrics
	a.metrics.tokenGenerations++

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
// It verifies the token signature, checks the expiration time, and validates
// the issuer. If remote validation is configured, it also validates the token
// with the remote validation service.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for tracing and cancellation.
//   - tokenString: The JWT token string to validate.
//
// Returns:
//   - *jwt.Claims: The claims from the token if validation is successful.
//   - error: An error if token validation fails, such as if the token is expired,
//     has an invalid signature, or fails remote validation.
func (a *Auth) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.ValidateToken")
	defer span.End()

	// Record metrics
	a.metrics.tokenValidations++

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
// It uses the user's roles and the operation name to determine if the user
// is authorized. Users with the admin role are authorized for all operations.
// Users with the read-only role are authorized only for read operations.
//
// Parameters:
//   - ctx: The context for the operation, which must contain user information.
//     This is typically set by the authentication middleware.
//   - operation: The name of the operation to check authorization for.
//     This should be a string that identifies the operation, such as "read:users".
//
// Returns:
//   - bool: true if the user is authorized, false otherwise.
//   - error: An error if the authorization check fails, such as if the context
//     does not contain user information.
func (a *Auth) IsAuthorized(ctx context.Context, operation string) (bool, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.IsAuthorized")
	defer span.End()

	// Record metrics
	a.metrics.authorizationChecks++

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

// IsAdmin checks if the user has the admin role.
// This is a convenience method that checks if the user has the role configured
// as the admin role in the service configuration.
//
// Parameters:
//   - ctx: The context for the operation, which must contain user information.
//     This is typically set by the authentication middleware.
//
// Returns:
//   - bool: true if the user has the admin role, false otherwise.
//   - error: An error if the role check fails, such as if the context
//     does not contain user information.
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
// It retrieves the user's roles from the context and checks if the specified
// role is included in the user's roles. This method is case-sensitive.
//
// Parameters:
//   - ctx: The context for the operation, which must contain user information.
//     This is typically set by the authentication middleware.
//   - role: The role to check for, such as "admin" or "user".
//
// Returns:
//   - bool: true if the user has the specified role, false otherwise.
//   - error: An error if the role check fails, such as if the context
//     does not contain user information.
func (a *Auth) HasRole(ctx context.Context, role string) (bool, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.HasRole")
	defer span.End()

	// Add span attributes
	telemetry.AddSpanAttributes(ctx,
		attribute.String("role", role),
	)

	// Check role
	start := time.Now()
	hasRole, err := a.authService.HasRole(ctx, role)
	duration := time.Since(start)

	// Record result
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		telemetry.RecordErrorMetric(ctx, "auth", "role_check_error")
	}
	telemetry.AddSpanAttributes(ctx,
		attribute.Bool("has_role", hasRole),
		attribute.Float64("duration_ms", float64(duration.Milliseconds())),
	)

	return hasRole, err
}

// GetUserID retrieves the user ID from the context.
// This method extracts the user ID that was previously set in the context
// by the authentication middleware. It's a convenience method for accessing
// the authenticated user's identity.
//
// Parameters:
//   - ctx: The context for the operation, which must contain user information.
//     This is typically set by the authentication middleware.
//
// Returns:
//   - string: The user ID if found in the context.
//   - error: An error if the user ID cannot be retrieved, such as if the context
//     does not contain user information or if the user is not authenticated.
func (a *Auth) GetUserID(ctx context.Context) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.GetUserID")
	defer span.End()

	// Get user ID
	start := time.Now()
	userID, err := a.authService.GetUserID(ctx)
	duration := time.Since(start)

	// Record result
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		telemetry.RecordErrorMetric(ctx, "auth", "get_user_id_error")
	} else {
		telemetry.AddSpanAttributes(ctx,
			attribute.String("user.id", userID),
		)
	}
	telemetry.AddSpanAttributes(ctx,
		attribute.Float64("duration_ms", float64(duration.Milliseconds())),
	)

	return userID, err
}

// GetUserRoles retrieves the user roles from the context.
// This method extracts the user roles that were previously set in the context
// by the authentication middleware. It's a convenience method for accessing
// the authenticated user's roles for authorization purposes.
//
// Parameters:
//   - ctx: The context for the operation, which must contain user information.
//     This is typically set by the authentication middleware.
//
// Returns:
//   - []string: A slice of role strings if found in the context.
//   - error: An error if the user roles cannot be retrieved, such as if the context
//     does not contain user information or if the user is not authenticated.
func (a *Auth) GetUserRoles(ctx context.Context) ([]string, error) {
	ctx, span := telemetry.StartSpan(ctx, "auth.GetUserRoles")
	defer span.End()

	// Get user roles
	start := time.Now()
	roles, err := a.authService.GetUserRoles(ctx)
	duration := time.Since(start)

	// Record result
	if err != nil {
		telemetry.RecordErrorSpan(ctx, err)
		telemetry.RecordErrorMetric(ctx, "auth", "get_user_roles_error")
	} else {
		telemetry.AddSpanAttributes(ctx,
			attribute.StringSlice("user.roles", roles),
		)
	}
	telemetry.AddSpanAttributes(ctx,
		attribute.Float64("duration_ms", float64(duration.Milliseconds())),
	)

	return roles, err
}

// WithUserID returns a new context with the user ID.
// This function adds the user ID to the context, which can be later retrieved
// using GetUserIDFromContext. It's typically used in testing or when manually
// setting up an authenticated context without going through the normal
// authentication flow.
//
// Parameters:
//   - ctx: The parent context to which the user ID will be added.
//   - userID: The user ID to add to the context.
//
// Returns:
//   - context.Context: A new context containing the user ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return middleware.WithUserID(ctx, userID)
}

// WithUserRoles returns a new context with the user roles.
// This function adds the user roles to the context, which can be later retrieved
// using GetUserRolesFromContext. It's typically used in testing or when manually
// setting up an authenticated context without going through the normal
// authentication flow.
//
// Parameters:
//   - ctx: The parent context to which the user roles will be added.
//   - roles: A slice of role strings to add to the context.
//
// Returns:
//   - context.Context: A new context containing the user roles.
func WithUserRoles(ctx context.Context, roles []string) context.Context {
	return middleware.WithUserRoles(ctx, roles)
}

// GetUserIDFromContext retrieves the user ID from the context.
// This function extracts the user ID that was previously set in the context
// by either the authentication middleware or the WithUserID function.
//
// Parameters:
//   - ctx: The context from which to retrieve the user ID.
//
// Returns:
//   - string: The user ID if found in the context.
//   - bool: true if the user ID was found in the context, false otherwise.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	return middleware.GetUserID(ctx)
}

// GetUserRolesFromContext retrieves the user roles from the context.
// This function extracts the user roles that were previously set in the context
// by either the authentication middleware or the WithUserRoles function.
//
// Parameters:
//   - ctx: The context from which to retrieve the user roles.
//
// Returns:
//   - []string: A slice of role strings if found in the context.
//   - bool: true if the user roles were found in the context, false otherwise.
func GetUserRolesFromContext(ctx context.Context) ([]string, bool) {
	return middleware.GetUserRoles(ctx)
}

// IsAuthenticated checks if the user is authenticated.
// This function determines if the context contains authentication information,
// specifically a user ID. It's a convenience function for checking if a request
// has gone through the authentication process successfully.
//
// Parameters:
//   - ctx: The context to check for authentication information.
//
// Returns:
//   - bool: true if the user is authenticated (has a user ID in the context),
//     false otherwise.
func IsAuthenticated(ctx context.Context) bool {
	return middleware.IsAuthenticated(ctx)
}
