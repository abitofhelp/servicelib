// Copyright (c) 2025 A Bit of Help, Inc.

package config_test

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	authconfig "github.com/abitofhelp/servicelib/auth/config"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/service"
	"go.uber.org/zap"
)

func Example() {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create an auth configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "example-secret-key"
	config.OIDC.IssuerURL = "https://example.com/oidc"
	config.OIDC.ClientID = "example-client-id"
	config.OIDC.ClientSecret = "example-client-secret"

	// Create an auth config adapter
	adapter := authconfig.NewAuthConfigAdapter(config)

	// Get the auth configuration
	authCfg := adapter.GetAuth()

	// Use the auth configuration to create JWT, OIDC, middleware, and service configurations
	jwtConfig := authconfig.CreateJWTConfig(authCfg)
	jwtRemoteConfig := authconfig.CreateJWTRemoteConfig(authCfg)
	oidcConfig := authconfig.CreateOIDCConfig(authCfg)
	middlewareConfig := authconfig.CreateMiddlewareConfig(authCfg)
	serviceConfig := authconfig.CreateServiceConfig(authCfg)

	// Create JWT service
	jwtService := jwt.NewService(jwtConfig, logger)

	// Add remote validator if enabled
	if authCfg.GetJWT().GetRemote().GetEnabled() {
		jwtService.WithRemoteValidator(jwtRemoteConfig)
	}

	// Skip creating a real OIDC service since it requires an external provider
	// Instead, just check that the configurations were created correctly
	fmt.Println("JWT config created:", jwtConfig.SecretKey != "")
	fmt.Println("OIDC config created:", oidcConfig.IssuerURL != "")
	fmt.Println("Middleware config created:", middlewareConfig.RequireAuth)
	fmt.Println("Service config created:", serviceConfig.AdminRoleName != "")

	// Create JWT service
	fmt.Println("JWT service created:", jwtService != nil)

	// Create service
	authService := service.NewService(serviceConfig, logger)
	fmt.Println("Auth service created:", authService != nil)

	// Output:
	// JWT config created: true
	// OIDC config created: true
	// Middleware config created: true
	// Service config created: true
	// JWT service created: true
	// Auth service created: true
}

func ExampleAuthConfigAdapter_AsGenericConfig() {
	// Create an auth configuration
	config := auth.DefaultConfig()

	// Create an auth config adapter
	adapter := authconfig.NewAuthConfigAdapter(config)

	// Convert to generic config
	genericConfig := adapter.AsGenericConfig()

	// Use the generic config
	appConfig := genericConfig.GetApp()
	fmt.Println("App name:", appConfig.GetName())
	fmt.Println("App environment:", appConfig.GetEnvironment())

	// Output:
	// App name: auth
	// App environment: production
}

func ExampleCreateJWTConfig() {
	// Create an auth configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "example-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "example-issuer"

	// Create an auth config adapter
	adapter := authconfig.NewAuthConfigAdapter(config)

	// Get the auth configuration
	authCfg := adapter.GetAuth()

	// Create JWT configuration
	jwtConfig := authconfig.CreateJWTConfig(authCfg)

	// Use the JWT configuration
	fmt.Println("Secret key:", jwtConfig.SecretKey)
	fmt.Println("Token duration:", jwtConfig.TokenDuration)
	fmt.Println("Issuer:", jwtConfig.Issuer)

	// Output:
	// Secret key: example-secret-key
	// Token duration: 1h0m0s
	// Issuer: example-issuer
}

func ExampleCreateOIDCConfig() {
	// Create an auth configuration
	config := auth.DefaultConfig()
	config.OIDC.IssuerURL = "https://example.com/oidc"
	config.OIDC.ClientID = "example-client-id"
	config.OIDC.ClientSecret = "example-client-secret"
	config.OIDC.RedirectURL = "https://myapp.com/callback"
	config.OIDC.Scopes = []string{"openid", "profile", "email", "custom-scope"}
	config.OIDC.Timeout = 30 * time.Second

	// Create an auth config adapter
	adapter := authconfig.NewAuthConfigAdapter(config)

	// Get the auth configuration
	authCfg := adapter.GetAuth()

	// Create OIDC configuration
	oidcConfig := authconfig.CreateOIDCConfig(authCfg)

	// Use the OIDC configuration
	fmt.Println("Issuer URL:", oidcConfig.IssuerURL)
	fmt.Println("Client ID:", oidcConfig.ClientID)
	fmt.Println("Redirect URL:", oidcConfig.RedirectURL)
	fmt.Println("Timeout:", oidcConfig.Timeout)
	fmt.Println("Number of scopes:", len(oidcConfig.Scopes))

	// Output:
	// Issuer URL: https://example.com/oidc
	// Client ID: example-client-id
	// Redirect URL: https://myapp.com/callback
	// Timeout: 30s
	// Number of scopes: 4
}
