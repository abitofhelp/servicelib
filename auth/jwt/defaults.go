// Copyright (c) 2025 A Bit of Help, Inc.

package jwt

import "time"

// DefaultConfig returns a default configuration for JWT token handling.
func DefaultConfig() Config {
	return Config{
		SecretKey:          "default-secret-key-please-change-in-production",
		TokenDuration:      24 * time.Hour,
		Issuer:             "servicelib",
		SigningMethod:      SigningMethodHS256,
		MinSecretKeyLength: 32,
	}
}

// DefaultRemoteConfig returns a default configuration for remote JWT token validation.
func DefaultRemoteConfig() RemoteConfig {
	return RemoteConfig{
		ValidationURL: "http://localhost:8080/validate",
		ClientID:      "",
		ClientSecret:  "",
		Timeout:       10 * time.Second,
	}
}