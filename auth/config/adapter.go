// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides adapters for auth configuration.
package config

import (
	"time"

	"github.com/abitofhelp/servicelib/auth"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/abitofhelp/servicelib/auth/service"
	baseconfig "github.com/abitofhelp/servicelib/config"
)

// AuthConfigProvider defines the interface for accessing auth configuration
type AuthConfigProvider interface {
	// GetJWTSecretKey returns the JWT secret key
	GetJWTSecretKey() string

	// GetJWTTokenDuration returns the JWT token duration
	GetJWTTokenDuration() time.Duration

	// GetJWTIssuer returns the JWT issuer
	GetJWTIssuer() string

	// GetJWTRemoteEnabled returns whether JWT remote validation is enabled
	GetJWTRemoteEnabled() bool

	// GetJWTRemoteValidationURL returns the JWT remote validation URL
	GetJWTRemoteValidationURL() string

	// GetJWTRemoteClientID returns the JWT remote client ID
	GetJWTRemoteClientID() string

	// GetJWTRemoteClientSecret returns the JWT remote client secret
	GetJWTRemoteClientSecret() string

	// GetJWTRemoteTimeout returns the JWT remote timeout
	GetJWTRemoteTimeout() time.Duration

	// GetOIDCIssuerURL returns the OIDC issuer URL
	GetOIDCIssuerURL() string

	// GetOIDCClientID returns the OIDC client ID
	GetOIDCClientID() string

	// GetOIDCClientSecret returns the OIDC client secret
	GetOIDCClientSecret() string

	// GetOIDCRedirectURL returns the OIDC redirect URL
	GetOIDCRedirectURL() string

	// GetOIDCScopes returns the OIDC scopes
	GetOIDCScopes() []string

	// GetOIDCTimeout returns the OIDC timeout
	GetOIDCTimeout() time.Duration

	// GetMiddlewareSkipPaths returns the middleware skip paths
	GetMiddlewareSkipPaths() []string

	// GetMiddlewareRequireAuth returns whether middleware requires auth
	GetMiddlewareRequireAuth() bool

	// GetServiceAdminRoleName returns the service admin role name
	GetServiceAdminRoleName() string

	// GetServiceReadOnlyRoleName returns the service read-only role name
	GetServiceReadOnlyRoleName() string

	// GetServiceReadOperationPrefixes returns the service read operation prefixes
	GetServiceReadOperationPrefixes() []string
}

// AuthConfigAdapter adapts the auth.Config to the config package interfaces
type AuthConfigAdapter struct {
	config auth.Config
}

// NewAuthConfigAdapter creates a new AuthConfigAdapter
func NewAuthConfigAdapter(config auth.Config) *AuthConfigAdapter {
	return &AuthConfigAdapter{
		config: config,
	}
}

// GetAuth returns the auth configuration
func (a *AuthConfigAdapter) GetAuth() Auth {
	return &AuthAdapter{
		config: a.config,
	}
}

// AsGenericConfig returns the auth configuration as a generic config
func (a *AuthConfigAdapter) AsGenericConfig() baseconfig.Config {
	return baseconfig.NewGenericConfigAdapter(a).
		WithAppName("auth").
		WithAppEnvironment("production")
}

// Auth is the interface for auth configuration
type Auth interface {
	// GetJWT returns the JWT configuration
	GetJWT() JWT

	// GetOIDC returns the OIDC configuration
	GetOIDC() OIDC

	// GetMiddleware returns the middleware configuration
	GetMiddleware() Middleware

	// GetService returns the service configuration
	GetService() Service
}

// JWT is the interface for JWT configuration
type JWT interface {
	// GetSecretKey returns the JWT secret key
	GetSecretKey() string

	// GetTokenDuration returns the JWT token duration
	GetTokenDuration() time.Duration

	// GetIssuer returns the JWT issuer
	GetIssuer() string

	// GetRemote returns the JWT remote configuration
	GetRemote() JWTRemote
}

// JWTRemote is the interface for JWT remote configuration
type JWTRemote interface {
	// GetEnabled returns whether JWT remote validation is enabled
	GetEnabled() bool

	// GetValidationURL returns the JWT remote validation URL
	GetValidationURL() string

	// GetClientID returns the JWT remote client ID
	GetClientID() string

	// GetClientSecret returns the JWT remote client secret
	GetClientSecret() string

	// GetTimeout returns the JWT remote timeout
	GetTimeout() time.Duration
}

// OIDC is the interface for OIDC configuration
type OIDC interface {
	// GetIssuerURL returns the OIDC issuer URL
	GetIssuerURL() string

	// GetClientID returns the OIDC client ID
	GetClientID() string

	// GetClientSecret returns the OIDC client secret
	GetClientSecret() string

	// GetRedirectURL returns the OIDC redirect URL
	GetRedirectURL() string

	// GetScopes returns the OIDC scopes
	GetScopes() []string

	// GetTimeout returns the OIDC timeout
	GetTimeout() time.Duration
}

// Middleware is the interface for middleware configuration
type Middleware interface {
	// GetSkipPaths returns the middleware skip paths
	GetSkipPaths() []string

	// GetRequireAuth returns whether middleware requires auth
	GetRequireAuth() bool
}

// Service is the interface for service configuration
type Service interface {
	// GetAdminRoleName returns the service admin role name
	GetAdminRoleName() string

	// GetReadOnlyRoleName returns the service read-only role name
	GetReadOnlyRoleName() string

	// GetReadOperationPrefixes returns the service read operation prefixes
	GetReadOperationPrefixes() []string
}

// AuthAdapter adapts the auth.Config to the Auth interface
type AuthAdapter struct {
	config auth.Config
}

// GetJWT returns the JWT configuration
func (a *AuthAdapter) GetJWT() JWT {
	return &JWTAdapter{
		config: a.config.JWT,
	}
}

// GetOIDC returns the OIDC configuration
func (a *AuthAdapter) GetOIDC() OIDC {
	return &OIDCAdapter{
		config: a.config.OIDC,
	}
}

// GetMiddleware returns the middleware configuration
func (a *AuthAdapter) GetMiddleware() Middleware {
	return &MiddlewareAdapter{
		config: a.config.Middleware,
	}
}

// GetService returns the service configuration
func (a *AuthAdapter) GetService() Service {
	return &ServiceAdapter{
		config: a.config.Service,
	}
}

// JWTAdapter adapts the auth.Config.JWT to the JWT interface
type JWTAdapter struct {
	config struct {
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
}

// GetSecretKey returns the JWT secret key
func (a *JWTAdapter) GetSecretKey() string {
	return a.config.SecretKey
}

// GetTokenDuration returns the JWT token duration
func (a *JWTAdapter) GetTokenDuration() time.Duration {
	return a.config.TokenDuration
}

// GetIssuer returns the JWT issuer
func (a *JWTAdapter) GetIssuer() string {
	return a.config.Issuer
}

// GetRemote returns the JWT remote configuration
func (a *JWTAdapter) GetRemote() JWTRemote {
	return &JWTRemoteAdapter{
		config: a.config.Remote,
	}
}

// JWTRemoteAdapter adapts the auth.Config.JWT.Remote to the JWTRemote interface
type JWTRemoteAdapter struct {
	config struct {
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

// GetEnabled returns whether JWT remote validation is enabled
func (a *JWTRemoteAdapter) GetEnabled() bool {
	return a.config.Enabled
}

// GetValidationURL returns the JWT remote validation URL
func (a *JWTRemoteAdapter) GetValidationURL() string {
	return a.config.ValidationURL
}

// GetClientID returns the JWT remote client ID
func (a *JWTRemoteAdapter) GetClientID() string {
	return a.config.ClientID
}

// GetClientSecret returns the JWT remote client secret
func (a *JWTRemoteAdapter) GetClientSecret() string {
	return a.config.ClientSecret
}

// GetTimeout returns the JWT remote timeout
func (a *JWTRemoteAdapter) GetTimeout() time.Duration {
	return a.config.Timeout
}

// OIDCAdapter adapts the auth.Config.OIDC to the OIDC interface
type OIDCAdapter struct {
	config struct {
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
}

// GetIssuerURL returns the OIDC issuer URL
func (a *OIDCAdapter) GetIssuerURL() string {
	return a.config.IssuerURL
}

// GetClientID returns the OIDC client ID
func (a *OIDCAdapter) GetClientID() string {
	return a.config.ClientID
}

// GetClientSecret returns the OIDC client secret
func (a *OIDCAdapter) GetClientSecret() string {
	return a.config.ClientSecret
}

// GetRedirectURL returns the OIDC redirect URL
func (a *OIDCAdapter) GetRedirectURL() string {
	return a.config.RedirectURL
}

// GetScopes returns the OIDC scopes
func (a *OIDCAdapter) GetScopes() []string {
	return a.config.Scopes
}

// GetTimeout returns the OIDC timeout
func (a *OIDCAdapter) GetTimeout() time.Duration {
	return a.config.Timeout
}

// MiddlewareAdapter adapts the auth.Config.Middleware to the Middleware interface
type MiddlewareAdapter struct {
	config struct {
		// SkipPaths are paths that should skip authentication
		SkipPaths []string

		// RequireAuth determines if authentication is required for all requests
		RequireAuth bool
	}
}

// GetSkipPaths returns the middleware skip paths
func (a *MiddlewareAdapter) GetSkipPaths() []string {
	return a.config.SkipPaths
}

// GetRequireAuth returns whether middleware requires auth
func (a *MiddlewareAdapter) GetRequireAuth() bool {
	return a.config.RequireAuth
}

// ServiceAdapter adapts the auth.Config.Service to the Service interface
type ServiceAdapter struct {
	config struct {
		// AdminRoleName is the name of the admin role
		AdminRoleName string

		// ReadOnlyRoleName is the name of the read-only role
		ReadOnlyRoleName string

		// ReadOperationPrefixes are prefixes for read-only operations
		ReadOperationPrefixes []string
	}
}

// GetAdminRoleName returns the service admin role name
func (a *ServiceAdapter) GetAdminRoleName() string {
	return a.config.AdminRoleName
}

// GetReadOnlyRoleName returns the service read-only role name
func (a *ServiceAdapter) GetReadOnlyRoleName() string {
	return a.config.ReadOnlyRoleName
}

// GetReadOperationPrefixes returns the service read operation prefixes
func (a *ServiceAdapter) GetReadOperationPrefixes() []string {
	return a.config.ReadOperationPrefixes
}

// CreateJWTConfig creates a JWT configuration from the auth configuration
func CreateJWTConfig(authConfig Auth) jwt.Config {
	jwtConfig := jwt.Config{
		SecretKey:     authConfig.GetJWT().GetSecretKey(),
		TokenDuration: authConfig.GetJWT().GetTokenDuration(),
		Issuer:        authConfig.GetJWT().GetIssuer(),
	}
	return jwtConfig
}

// CreateJWTRemoteConfig creates a JWT remote configuration from the auth configuration
func CreateJWTRemoteConfig(authConfig Auth) jwt.RemoteConfig {
	remote := authConfig.GetJWT().GetRemote()
	remoteConfig := jwt.RemoteConfig{
		ValidationURL: remote.GetValidationURL(),
		ClientID:      remote.GetClientID(),
		ClientSecret:  remote.GetClientSecret(),
		Timeout:       remote.GetTimeout(),
	}
	return remoteConfig
}

// CreateOIDCConfig creates an OIDC configuration from the auth configuration
func CreateOIDCConfig(authConfig Auth) oidc.Config {
	oidcConfig := oidc.Config{
		IssuerURL:     authConfig.GetOIDC().GetIssuerURL(),
		ClientID:      authConfig.GetOIDC().GetClientID(),
		ClientSecret:  authConfig.GetOIDC().GetClientSecret(),
		RedirectURL:   authConfig.GetOIDC().GetRedirectURL(),
		Scopes:        authConfig.GetOIDC().GetScopes(),
		Timeout:       authConfig.GetOIDC().GetTimeout(),
		AdminRoleName: authConfig.GetService().GetAdminRoleName(),
	}
	return oidcConfig
}

// CreateMiddlewareConfig creates a middleware configuration from the auth configuration
func CreateMiddlewareConfig(authConfig Auth) middleware.Config {
	middlewareConfig := middleware.Config{
		SkipPaths:   authConfig.GetMiddleware().GetSkipPaths(),
		RequireAuth: authConfig.GetMiddleware().GetRequireAuth(),
	}
	return middlewareConfig
}

// CreateServiceConfig creates a service configuration from the auth configuration
func CreateServiceConfig(authConfig Auth) service.Config {
	serviceConfig := service.Config{
		AdminRoleName:         authConfig.GetService().GetAdminRoleName(),
		ReadOnlyRoleName:      authConfig.GetService().GetReadOnlyRoleName(),
		ReadOperationPrefixes: authConfig.GetService().GetReadOperationPrefixes(),
	}
	return serviceConfig
}