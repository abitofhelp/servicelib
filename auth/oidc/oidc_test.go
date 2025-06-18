// Copyright (c) 2025 A Bit of Help, Inc.

package oidc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestNewService tests the creation of a new OIDC service
func TestNewService(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	logger := zap.NewNop()
	ctx := context.Background()

	// Test with valid configuration
	config := oidc.Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}

	service, err := oidc.NewService(ctx, config, logger)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test with empty issuer URL
	invalidConfig := oidc.Config{
		IssuerURL:    "",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	service, err = oidc.NewService(ctx, invalidConfig, logger)
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidConfig))

	// Test with empty client ID
	invalidConfig = oidc.Config{
		IssuerURL:    "https://accounts.google.com",
		ClientID:     "",
		ClientSecret: "test-client-secret",
	}

	service, err = oidc.NewService(ctx, invalidConfig, logger)
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidConfig))

	// Test with nil context
	service, err = oidc.NewService(nil, config, logger)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test with nil logger
	service, err = oidc.NewService(ctx, config, nil)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

// TestValidateToken tests token validation
func TestValidateToken(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	logger := zap.NewNop()
	ctx := context.Background()

	config := oidc.Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}

	service, err := oidc.NewService(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Test with empty token
	claims, err := service.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.True(t, errors.Is(err, autherrors.ErrMissingToken))

	// Test with invalid token
	claims, err = service.ValidateToken(ctx, "invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidToken))
}

// TestIsAdmin tests the IsAdmin function
func TestIsAdmin(t *testing.T) {
	// Create a real service for testing IsAdmin with "admin" role
	adminService, err := oidc.NewService(context.Background(), oidc.Config{
		IssuerURL:     "https://example.com", // Not used for IsAdmin
		ClientID:      "test-client",         // Not used for IsAdmin
		AdminRoleName: "admin",
	}, zap.NewNop())

	// Skip the test if we can't create the service
	if err != nil {
		t.Skip("Skipping test due to service creation error:", err)
	}

	// Test with admin role
	roles := []string{"user", "admin", "editor"}
	isAdmin := adminService.IsAdmin(roles)
	assert.True(t, isAdmin)

	// Test without admin role
	roles = []string{"user", "editor"}
	isAdmin = adminService.IsAdmin(roles)
	assert.False(t, isAdmin)

	// Test with empty roles
	roles = []string{}
	isAdmin = adminService.IsAdmin(roles)
	assert.False(t, isAdmin)

	// Test with nil roles
	isAdmin = adminService.IsAdmin(nil)
	assert.False(t, isAdmin)

	// Create another service with a different admin role name
	superuserService, err := oidc.NewService(context.Background(), oidc.Config{
		IssuerURL:     "https://example.com", // Not used for IsAdmin
		ClientID:      "test-client",         // Not used for IsAdmin
		AdminRoleName: "superuser",
	}, zap.NewNop())

	// Skip the test if we can't create the service
	if err != nil {
		t.Skip("Skipping test due to service creation error:", err)
	}

	// Test with admin role but different admin role name
	roles = []string{"user", "admin", "editor"}
	isAdmin = superuserService.IsAdmin(roles)
	assert.False(t, isAdmin)

	// Test with superuser role
	roles = []string{"user", "superuser", "editor"}
	isAdmin = superuserService.IsAdmin(roles)
	assert.True(t, isAdmin)
}

// TestGetAuthURL tests the GetAuthURL function
func TestGetAuthURL(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	logger := zap.NewNop()
	ctx := context.Background()

	config := oidc.Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}

	service, err := oidc.NewService(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Test GetAuthURL
	state := "random-state-string"
	url := service.GetAuthURL(state)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, config.ClientID)
	assert.Contains(t, url, config.RedirectURL)
	assert.Contains(t, url, state)
}

// TestExchange tests the Exchange function
func TestExchange(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	logger := zap.NewNop()
	ctx := context.Background()

	config := oidc.Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}

	service, err := oidc.NewService(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Test with invalid code
	token, err := service.Exchange(ctx, "invalid-code")
	assert.Error(t, err)
	assert.Nil(t, token)
}

// TestGetUserInfo tests the GetUserInfo function
func TestGetUserInfo(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	logger := zap.NewNop()
	ctx := context.Background()

	config := oidc.Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}

	service, err := oidc.NewService(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, service)

	// We can't test GetUserInfo without a valid token, which we can't get without a real OAuth flow
	// This test is included for completeness but is skipped
}
