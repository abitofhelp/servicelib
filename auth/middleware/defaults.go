// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

// DefaultConfig returns a default configuration for the authentication middleware.
func DefaultConfig() Config {
	return Config{
		SkipPaths:   []string{"/health", "/metrics", "/swagger"},
		RequireAuth: true,
	}
}
