// Copyright (c) 2025 A Bit of Help, Inc.

package config_test

import (
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	authconfig "github.com/abitofhelp/servicelib/auth/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthConfigAdapter(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 2 * time.Hour
	config.JWT.Issuer = "test-issuer"
	config.JWT.Remote.Enabled = true
	config.JWT.Remote.ValidationURL = "https://test.com/validate"
	config.JWT.Remote.ClientID = "test-client-id"
	config.JWT.Remote.ClientSecret = "test-client-secret"
	config.JWT.Remote.Timeout = 15 * time.Second

	config.OIDC.IssuerURL = "https://test.com/oidc"
	config.OIDC.ClientID = "test-oidc-client-id"
	config.OIDC.ClientSecret = "test-oidc-client-secret"
	config.OIDC.RedirectURL = "https://myapp.com/callback"
	config.OIDC.Scopes = []string{"openid", "profile", "email", "custom-scope"}
	config.OIDC.Timeout = 20 * time.Second

	config.Middleware.SkipPaths = []string{"/health", "/metrics"}
	config.Middleware.RequireAuth = true

	config.Service.AdminRoleName = "test-admin"
	config.Service.ReadOnlyRoleName = "test-reader"
	config.Service.ReadOperationPrefixes = []string{"read:", "get:", "list:"}

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	require.NotNil(t, adapter, "Adapter should not be nil")

	// Test GetAuth
	authCfg := adapter.GetAuth()
	require.NotNil(t, authCfg, "Auth config should not be nil")

	// Test AsGenericConfig
	genericConfig := adapter.AsGenericConfig()
	require.NotNil(t, genericConfig, "Generic config should not be nil")
	appConfig := genericConfig.GetApp()
	assert.Equal(t, "auth", appConfig.GetName())
	assert.Equal(t, "production", appConfig.GetEnvironment())
}

func TestAuthAdapter(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 2 * time.Hour
	config.JWT.Issuer = "test-issuer"
	config.JWT.Remote.Enabled = true
	config.JWT.Remote.ValidationURL = "https://test.com/validate"
	config.JWT.Remote.ClientID = "test-client-id"
	config.JWT.Remote.ClientSecret = "test-client-secret"
	config.JWT.Remote.Timeout = 15 * time.Second

	config.OIDC.IssuerURL = "https://test.com/oidc"
	config.OIDC.ClientID = "test-oidc-client-id"
	config.OIDC.ClientSecret = "test-oidc-client-secret"
	config.OIDC.RedirectURL = "https://myapp.com/callback"
	config.OIDC.Scopes = []string{"openid", "profile", "email", "custom-scope"}
	config.OIDC.Timeout = 20 * time.Second

	config.Middleware.SkipPaths = []string{"/health", "/metrics"}
	config.Middleware.RequireAuth = true

	config.Service.AdminRoleName = "test-admin"
	config.Service.ReadOnlyRoleName = "test-reader"
	config.Service.ReadOperationPrefixes = []string{"read:", "get:", "list:"}

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// Test JWT interface
	jwt := authCfg.GetJWT()
	assert.Equal(t, "test-secret-key", jwt.GetSecretKey())
	assert.Equal(t, 2*time.Hour, jwt.GetTokenDuration())
	assert.Equal(t, "test-issuer", jwt.GetIssuer())

	// Test JWT Remote interface
	jwtRemote := jwt.GetRemote()
	assert.True(t, jwtRemote.GetEnabled())
	assert.Equal(t, "https://test.com/validate", jwtRemote.GetValidationURL())
	assert.Equal(t, "test-client-id", jwtRemote.GetClientID())
	assert.Equal(t, "test-client-secret", jwtRemote.GetClientSecret())
	assert.Equal(t, 15*time.Second, jwtRemote.GetTimeout())

	// Test OIDC interface
	oidc := authCfg.GetOIDC()
	assert.Equal(t, "https://test.com/oidc", oidc.GetIssuerURL())
	assert.Equal(t, "test-oidc-client-id", oidc.GetClientID())
	assert.Equal(t, "test-oidc-client-secret", oidc.GetClientSecret())
	assert.Equal(t, "https://myapp.com/callback", oidc.GetRedirectURL())
	assert.Equal(t, []string{"openid", "profile", "email", "custom-scope"}, oidc.GetScopes())
	assert.Equal(t, 20*time.Second, oidc.GetTimeout())

	// Test Middleware interface
	middleware := authCfg.GetMiddleware()
	assert.Equal(t, []string{"/health", "/metrics"}, middleware.GetSkipPaths())
	assert.True(t, middleware.GetRequireAuth())

	// Test Service interface
	service := authCfg.GetService()
	assert.Equal(t, "test-admin", service.GetAdminRoleName())
	assert.Equal(t, "test-reader", service.GetReadOnlyRoleName())
	assert.Equal(t, []string{"read:", "get:", "list:"}, service.GetReadOperationPrefixes())
}

func TestCreateJWTConfig(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 2 * time.Hour
	config.JWT.Issuer = "test-issuer"

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// Create JWT config
	jwtConfig := authconfig.CreateJWTConfig(authCfg)
	assert.Equal(t, "test-secret-key", jwtConfig.SecretKey)
	assert.Equal(t, 2*time.Hour, jwtConfig.TokenDuration)
	assert.Equal(t, "test-issuer", jwtConfig.Issuer)
}

func TestCreateJWTRemoteConfig(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.JWT.Remote.Enabled = true
	config.JWT.Remote.ValidationURL = "https://test.com/validate"
	config.JWT.Remote.ClientID = "test-client-id"
	config.JWT.Remote.ClientSecret = "test-client-secret"
	config.JWT.Remote.Timeout = 15 * time.Second

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// Create JWT remote config
	remoteConfig := authconfig.CreateJWTRemoteConfig(authCfg)
	assert.Equal(t, "https://test.com/validate", remoteConfig.ValidationURL)
	assert.Equal(t, "test-client-id", remoteConfig.ClientID)
	assert.Equal(t, "test-client-secret", remoteConfig.ClientSecret)
	assert.Equal(t, 15*time.Second, remoteConfig.Timeout)
}

func TestCreateOIDCConfig(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.OIDC.IssuerURL = "https://test.com/oidc"
	config.OIDC.ClientID = "test-oidc-client-id"
	config.OIDC.ClientSecret = "test-oidc-client-secret"
	config.OIDC.RedirectURL = "https://myapp.com/callback"
	config.OIDC.Scopes = []string{"openid", "profile", "email", "custom-scope"}
	config.OIDC.Timeout = 20 * time.Second
	config.Service.AdminRoleName = "test-admin"

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// Create OIDC config
	oidcConfig := authconfig.CreateOIDCConfig(authCfg)
	assert.Equal(t, "https://test.com/oidc", oidcConfig.IssuerURL)
	assert.Equal(t, "test-oidc-client-id", oidcConfig.ClientID)
	assert.Equal(t, "test-oidc-client-secret", oidcConfig.ClientSecret)
	assert.Equal(t, "https://myapp.com/callback", oidcConfig.RedirectURL)
	assert.Equal(t, []string{"openid", "profile", "email", "custom-scope"}, oidcConfig.Scopes)
	assert.Equal(t, 20*time.Second, oidcConfig.Timeout)
	assert.Equal(t, "test-admin", oidcConfig.AdminRoleName)
}

func TestCreateMiddlewareConfig(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.Middleware.SkipPaths = []string{"/health", "/metrics"}
	config.Middleware.RequireAuth = true

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// Create middleware config
	middlewareConfig := authconfig.CreateMiddlewareConfig(authCfg)
	assert.Equal(t, []string{"/health", "/metrics"}, middlewareConfig.SkipPaths)
	assert.True(t, middlewareConfig.RequireAuth)
}

func TestCreateServiceConfig(t *testing.T) {
	// Create a test configuration
	config := auth.DefaultConfig()
	config.Service.AdminRoleName = "test-admin"
	config.Service.ReadOnlyRoleName = "test-reader"
	config.Service.ReadOperationPrefixes = []string{"read:", "get:", "list:"}

	// Create the adapter
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// Create service config
	serviceConfig := authconfig.CreateServiceConfig(authCfg)
	assert.Equal(t, "test-admin", serviceConfig.AdminRoleName)
	assert.Equal(t, "test-reader", serviceConfig.ReadOnlyRoleName)
	assert.Equal(t, []string{"read:", "get:", "list:"}, serviceConfig.ReadOperationPrefixes)
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	// Test with empty config
	config := auth.Config{}
	adapter := authconfig.NewAuthConfigAdapter(config)
	authCfg := adapter.GetAuth()

	// JWT with empty values
	jwt := authCfg.GetJWT()
	assert.Empty(t, jwt.GetSecretKey())
	assert.Zero(t, jwt.GetTokenDuration())
	assert.Empty(t, jwt.GetIssuer())

	// JWT Remote with empty values
	jwtRemote := jwt.GetRemote()
	assert.False(t, jwtRemote.GetEnabled())
	assert.Empty(t, jwtRemote.GetValidationURL())
	assert.Empty(t, jwtRemote.GetClientID())
	assert.Empty(t, jwtRemote.GetClientSecret())
	assert.Zero(t, jwtRemote.GetTimeout())

	// OIDC with empty values
	oidc := authCfg.GetOIDC()
	assert.Empty(t, oidc.GetIssuerURL())
	assert.Empty(t, oidc.GetClientID())
	assert.Empty(t, oidc.GetClientSecret())
	assert.Empty(t, oidc.GetRedirectURL())
	assert.Empty(t, oidc.GetScopes())
	assert.Zero(t, oidc.GetTimeout())

	// Middleware with empty values
	middleware := authCfg.GetMiddleware()
	assert.Empty(t, middleware.GetSkipPaths())
	assert.False(t, middleware.GetRequireAuth())

	// Service with empty values
	service := authCfg.GetService()
	assert.Empty(t, service.GetAdminRoleName())
	assert.Empty(t, service.GetReadOnlyRoleName())
	assert.Empty(t, service.GetReadOperationPrefixes())

	// Test with nil values in slices
	configWithNil := auth.DefaultConfig()
	configWithNil.OIDC.Scopes = nil
	configWithNil.Middleware.SkipPaths = nil
	configWithNil.Service.ReadOperationPrefixes = nil

	adapterWithNil := authconfig.NewAuthConfigAdapter(configWithNil)
	authCfgWithNil := adapterWithNil.GetAuth()

	assert.Nil(t, authCfgWithNil.GetOIDC().GetScopes())
	assert.Nil(t, authCfgWithNil.GetMiddleware().GetSkipPaths())
	assert.Nil(t, authCfgWithNil.GetService().GetReadOperationPrefixes())
}
