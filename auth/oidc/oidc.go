// Copyright (c) 2025 A Bit of Help, Inc.

// Package oidc provides OpenID Connect (OIDC) authorization functionality for the family-service-graphql.
// This functionality is currently not enabled but will be in the future.
package oidc

// Config represents the configuration for OIDC authorization.
type Config struct {
	// Issuer is the URL of the OpenID Connect provider.
	Issuer string

	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// RedirectURL is the URL to redirect users going through
	// the OIDC flow.
	RedirectURL string

	// Scopes specifies optional requested permissions.
	Scopes []string
}

// Provider represents an OIDC provider.
type Provider struct {
	// Config is the OIDC configuration.
	Config *Config
}

// NewProvider creates a new OIDC provider with the given configuration.
// Note: This functionality is not currently enabled.
func NewProvider(config *Config) *Provider {
	return &Provider{
		Config: config,
	}
}
