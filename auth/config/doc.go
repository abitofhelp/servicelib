// Copyright (c) 2025 A Bit of Help, Inc.

/*
Package config provides adapters for integrating the auth configuration with the config package.

The auth configuration adapter allows you to:
  - Adapt the auth.Config to the config package interfaces
  - Access JWT, OIDC, Middleware, and Service configurations through a unified interface
  - Convert auth configuration to generic configuration
  - Create specific configurations for JWT, OIDC, Middleware, and Service components

Basic usage:

	// Create an auth configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "example-secret-key"

	// Create an auth config adapter
	adapter := authconfig.NewAuthConfigAdapter(config)

	// Get the auth configuration
	authCfg := adapter.GetAuth()

	// Use the auth configuration to create JWT configuration
	jwtConfig := authconfig.CreateJWTConfig(authCfg)

	// Create JWT service
	jwtService := jwt.NewService(jwtConfig, logger)

For more examples, see the example_test.go file.
*/
package config
