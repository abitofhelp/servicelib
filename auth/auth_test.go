// Copyright (c) 2025 A Bit of Help, Inc.

package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestNew tests the creation of a new Auth instance
func TestNew(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Test with valid configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	assert.NoError(t, err)
	assert.NotNil(t, authInstance)

	// Test with empty JWT secret key
	invalidConfig := auth.DefaultConfig()
	invalidConfig.JWT.SecretKey = "" // Empty secret key should cause validation error
	invalidConfig.JWT.TokenDuration = 1 * time.Hour
	invalidConfig.JWT.Issuer = "test-issuer"

	authInstance, err = auth.New(ctx, invalidConfig, logger)
	assert.Error(t, err)
	assert.Nil(t, authInstance)

	// Test with nil logger (should use NopLogger)
	authInstance, err = auth.New(ctx, config, nil)
	assert.NoError(t, err)
	assert.NotNil(t, authInstance)

	// Test with OIDC configuration
	oidcConfig := auth.DefaultConfig()
	oidcConfig.JWT.SecretKey = "test-secret-key"
	oidcConfig.JWT.TokenDuration = 1 * time.Hour
	oidcConfig.JWT.Issuer = "test-issuer"
	oidcConfig.OIDC.IssuerURL = "https://example.com"
	oidcConfig.OIDC.ClientID = "test-client-id"

	// Skip this test since it requires an external OIDC provider
	// authInstance, err = auth.New(ctx, oidcConfig, logger)
	// assert.NoError(t, err)
	// assert.NotNil(t, authInstance)
}

// TestDefaultConfig tests the DefaultConfig function
func TestDefaultConfig(t *testing.T) {
	config := auth.DefaultConfig()

	// Check JWT defaults
	assert.Equal(t, 24*time.Hour, config.JWT.TokenDuration)
	assert.Equal(t, "auth", config.JWT.Issuer)

	// Check OIDC defaults
	assert.Equal(t, 10*time.Second, config.OIDC.Timeout)
	assert.Equal(t, []string{"openid", "profile", "email"}, config.OIDC.Scopes)

	// Check Middleware defaults
	assert.True(t, config.Middleware.RequireAuth)

	// Check Service defaults
	assert.Equal(t, "admin", config.Service.AdminRoleName)
	assert.Equal(t, "authuser", config.Service.ReadOnlyRoleName)
	assert.Contains(t, config.Service.ReadOperationPrefixes, "read:")
	assert.Contains(t, config.Service.ReadOperationPrefixes, "get:")
	assert.Contains(t, config.Service.ReadOperationPrefixes, "list:")
	assert.Contains(t, config.Service.ReadOperationPrefixes, "find:")
	assert.Contains(t, config.Service.ReadOperationPrefixes, "query:")
	assert.Contains(t, config.Service.ReadOperationPrefixes, "count:")
}

// TestMiddleware tests the Middleware function
func TestMiddleware(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test that Middleware returns a function
	middlewareFunc := authInstance.Middleware()
	assert.NotNil(t, middlewareFunc)

	// Test that the middleware function can be used to wrap a handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middlewareFunc(handler)
	assert.NotNil(t, wrappedHandler)
}

// TestGenerateToken tests the GenerateToken function
func TestGenerateToken(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test generating a token
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := authInstance.GenerateToken(ctx, userID, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test with empty user ID
	token, err = authInstance.GenerateToken(ctx, "", roles)
	assert.Error(t, err)
	assert.Empty(t, token)
}

// TestValidateToken tests the ValidateToken function
func TestValidateToken(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Generate a valid token for testing
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := authInstance.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Test validating a valid token
	claims, err := authInstance.ValidateToken(ctx, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)

	// Test with empty token
	claims, err = authInstance.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Test with invalid token
	claims, err = authInstance.ValidateToken(ctx, "invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestIsAuthorized tests the IsAuthorized function
func TestIsAuthorized(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"
	config.Service.AdminRoleName = "admin"
	config.Service.ReadOnlyRoleName = "reader"
	config.Service.ReadOperationPrefixes = []string{"read:", "get:"}

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test with unauthenticated context
	authorized, err := authInstance.IsAuthorized(ctx, "read:resource")
	assert.Error(t, err)
	assert.False(t, authorized)

	// Test with admin role
	adminCtx := middleware.WithUserID(ctx, "admin123")
	adminCtx = middleware.WithUserRoles(adminCtx, []string{"admin"})

	authorized, err = authInstance.IsAuthorized(adminCtx, "write:resource")
	assert.NoError(t, err)
	assert.True(t, authorized)

	// Test with reader role and read operation
	readerCtx := middleware.WithUserID(ctx, "reader123")
	readerCtx = middleware.WithUserRoles(readerCtx, []string{"reader"})

	authorized, err = authInstance.IsAuthorized(readerCtx, "read:resource")
	assert.NoError(t, err)
	assert.True(t, authorized)

	// Test with reader role and write operation
	authorized, err = authInstance.IsAuthorized(readerCtx, "write:resource")
	assert.Error(t, err)
	assert.False(t, authorized)
}

// TestIsAdmin tests the IsAdmin function
func TestIsAdmin(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"
	config.Service.AdminRoleName = "admin"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test with unauthenticated context
	isAdmin, err := authInstance.IsAdmin(ctx)
	assert.Error(t, err)
	assert.False(t, isAdmin)

	// Test with admin role
	adminCtx := middleware.WithUserID(ctx, "admin123")
	adminCtx = middleware.WithUserRoles(adminCtx, []string{"admin"})

	isAdmin, err = authInstance.IsAdmin(adminCtx)
	assert.NoError(t, err)
	assert.True(t, isAdmin)

	// Test with non-admin role
	userCtx := middleware.WithUserID(ctx, "user123")
	userCtx = middleware.WithUserRoles(userCtx, []string{"user"})

	isAdmin, err = authInstance.IsAdmin(userCtx)
	assert.NoError(t, err)
	assert.False(t, isAdmin)
}

// TestHasRole tests the HasRole function
func TestHasRole(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test with unauthenticated context
	hasRole, err := authInstance.HasRole(ctx, "admin")
	assert.Error(t, err)
	assert.False(t, hasRole)

	// Test with matching role
	userCtx := middleware.WithUserID(ctx, "user123")
	userCtx = middleware.WithUserRoles(userCtx, []string{"admin", "user"})

	hasRole, err = authInstance.HasRole(userCtx, "admin")
	assert.NoError(t, err)
	assert.True(t, hasRole)

	// Test with non-matching role
	hasRole, err = authInstance.HasRole(userCtx, "editor")
	assert.NoError(t, err)
	assert.False(t, hasRole)
}

// TestGetUserID tests the GetUserID function
func TestGetUserID(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test with unauthenticated context
	userID, err := authInstance.GetUserID(ctx)
	assert.Error(t, err)
	assert.Empty(t, userID)

	// Test with authenticated context
	userCtx := middleware.WithUserID(ctx, "user123")

	userID, err = authInstance.GetUserID(userCtx)
	assert.NoError(t, err)
	assert.Equal(t, "user123", userID)
}

// TestGetUserRoles tests the GetUserRoles function
func TestGetUserRoles(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-issuer"

	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Test with unauthenticated context
	roles, err := authInstance.GetUserRoles(ctx)
	assert.Error(t, err)
	assert.Nil(t, roles)

	// Test with authenticated context
	userRoles := []string{"admin", "user"}
	userCtx := middleware.WithUserID(ctx, "user123")
	userCtx = middleware.WithUserRoles(userCtx, userRoles)

	roles, err = authInstance.GetUserRoles(userCtx)
	assert.NoError(t, err)
	assert.Equal(t, userRoles, roles)
}

// TestContextFunctions tests the context helper functions
func TestContextFunctions(t *testing.T) {
	// Test WithUserID and GetUserIDFromContext
	ctx := context.Background()
	userID := "user123"
	ctx = auth.WithUserID(ctx, userID)

	retrievedID, ok := auth.GetUserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, userID, retrievedID)

	// Test WithUserRoles and GetUserRolesFromContext
	roles := []string{"admin", "user"}
	ctx = auth.WithUserRoles(ctx, roles)

	retrievedRoles, ok := auth.GetUserRolesFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, roles, retrievedRoles)

	// Test IsAuthenticated
	assert.True(t, auth.IsAuthenticated(ctx))

	// Test with empty context
	emptyCtx := context.Background()
	_, ok = auth.GetUserIDFromContext(emptyCtx)
	assert.False(t, ok)

	_, ok = auth.GetUserRolesFromContext(emptyCtx)
	assert.False(t, ok)

	assert.False(t, auth.IsAuthenticated(emptyCtx))
}
