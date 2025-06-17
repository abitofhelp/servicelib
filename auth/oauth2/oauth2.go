// Copyright (c) 2025 A Bit of Help, Inc.

// Package oauth2 provides OAuth2 authentication functionality for the family-service-graphql.
// This functionality is currently not enabled but will be in the future.
package oauth2

// Config represents the configuration for OAuth2 authentication.
type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow.
	RedirectURL string

	// Scopes specifies optional requested permissions.
	Scopes []string
}

// Provider represents an OAuth2 provider.
type Provider struct {
	// Config is the OAuth2 configuration.
	Config *Config
}

// NewProvider creates a new OAuth2 provider with the given configuration.
// Note: This functionality is not currently enabled.
func NewProvider(config *Config) *Provider {
	return &Provider{
		Config: config,
	}
}