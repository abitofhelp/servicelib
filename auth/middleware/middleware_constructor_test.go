// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"testing"

	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestNewMiddleware tests the NewMiddleware function
func TestNewMiddleware(t *testing.T) {
	// Create a mock JWT service
	jwtService := &jwt.Service{}

	// Create a config
	config := Config{
		SkipPaths:   []string{"/health", "/metrics"},
		RequireAuth: true,
	}

	// Test with logger
	logger := zap.NewNop()
	middleware := NewMiddleware(jwtService, config, logger)

	// Verify the middleware was created correctly
	assert.NotNil(t, middleware)
	assert.Equal(t, config, middleware.config)
	assert.Equal(t, jwtService, middleware.jwtService)
	assert.Nil(t, middleware.oidcService)

	// Test with nil logger (should use NopLogger)
	middleware = NewMiddleware(jwtService, config, nil)
	assert.NotNil(t, middleware)
}

// TestNewMiddlewareWithOIDC tests the NewMiddlewareWithOIDC function
func TestNewMiddlewareWithOIDC(t *testing.T) {
	// Create mock services
	jwtService := &jwt.Service{}
	oidcService := &oidc.Service{}

	// Create a config
	config := Config{
		SkipPaths:   []string{"/health", "/metrics"},
		RequireAuth: true,
	}

	// Test with logger
	logger := zap.NewNop()
	middleware := NewMiddlewareWithOIDC(jwtService, oidcService, config, logger)

	// Verify the middleware was created correctly
	assert.NotNil(t, middleware)
	assert.Equal(t, config, middleware.config)
	assert.Equal(t, jwtService, middleware.jwtService)
	assert.Equal(t, oidcService, middleware.oidcService)

	// Test with nil logger (should use NopLogger)
	middleware = NewMiddlewareWithOIDC(jwtService, oidcService, config, nil)
	assert.NotNil(t, middleware)
}
