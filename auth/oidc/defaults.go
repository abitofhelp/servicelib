// Copyright (c) 2025 A Bit of Help, Inc.

package oidc

import "time"

// DefaultConfig returns a default configuration for OIDC integration.
func DefaultConfig() Config {
	return Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "",
		ClientSecret:  "",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
		RetryConfig:   DefaultRetryConfig(),
	}
}