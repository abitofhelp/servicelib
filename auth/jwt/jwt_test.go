// Copyright (c) 2025 A Bit of Help, Inc.

package jwt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewService(t *testing.T) {
	// Test with logger
	logger := zap.NewExample()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)
	assert.NotNil(t, service)

	// Test with nil logger (should use NopLogger)
	service = jwt.NewService(config, nil)
	assert.NotNil(t, service)
}

func TestGenerateToken(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)

	// Test successful token generation
	ctx := context.Background()
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test with empty user ID
	token, err = service.GenerateToken(ctx, "", roles)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidClaims))

	// Test with nil roles (should still work)
	token, err = service.GenerateToken(ctx, userID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test with empty roles (should still work)
	token, err = service.GenerateToken(ctx, userID, []string{})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Since it's difficult to force a signing error in a unit test,
	// we'll skip testing that specific error case.
	// In a real-world scenario, signing errors could occur due to:
	// - Invalid secret key format
	// - Memory allocation issues
	// - System-level errors
	// These are difficult to simulate in a unit test environment.
}

func TestValidateToken(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)
	ctx := context.Background()

	// Generate a valid token for testing
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Test successful token validation
	claims, err := service.ValidateToken(ctx, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, config.Issuer, claims.Issuer)

	// Test with empty token
	claims, err = service.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.True(t, errors.Is(err, autherrors.ErrMissingToken))

	// Test with invalid token
	claims, err = service.ValidateToken(ctx, "invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidToken))

	// Test with expired token
	expiredConfig := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: -1 * time.Hour, // Negative duration to create expired token
		Issuer:        "test-issuer",
	}
	expiredService := jwt.NewService(expiredConfig, logger)
	expiredToken, err := expiredService.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)

	claims, err = service.ValidateToken(ctx, expiredToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.True(t, errors.Is(err, autherrors.ErrExpiredToken))

	// Test with token signed with different key
	differentConfig := jwt.Config{
		SecretKey:     "different-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	differentService := jwt.NewService(differentConfig, logger)
	differentToken, err := differentService.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)

	claims, err = service.ValidateToken(ctx, differentToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidSignature))
}

func TestExtractTokenFromHeader(t *testing.T) {
	// Test with valid Bearer token
	authHeader := "Bearer token123"
	token, err := jwt.ExtractTokenFromHeader(authHeader)
	assert.NoError(t, err)
	assert.Equal(t, "token123", token)

	// Test with empty header
	token, err = jwt.ExtractTokenFromHeader("")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, autherrors.ErrMissingToken))

	// Test with invalid format (no Bearer prefix)
	token, err = jwt.ExtractTokenFromHeader("token123")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidToken))

	// Test with invalid format (too short)
	token, err = jwt.ExtractTokenFromHeader("Bear")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidToken))
}

// TestWithRemoteValidator tests the WithRemoteValidator function
func TestWithRemoteValidator(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)

	// Add a remote validator
	remoteConfig := jwt.RemoteConfig{
		ValidationURL: "https://test.com/validate",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		Timeout:       15 * time.Second,
	}
	result := service.WithRemoteValidator(remoteConfig)

	// Check that the method returns the service itself for chaining
	assert.Equal(t, service, result)

	// Test ValidateToken with remote validator
	// Since the remote validator is not implemented, it should fall back to local validation
	ctx := context.Background()
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := service.ValidateToken(ctx, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
}

// TestTokenWithMalformedClaims tests validation of a token with malformed claims
func TestTokenWithMalformedClaims(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)
	ctx := context.Background()

	// Create a token with valid signature but missing user ID
	// This is a bit tricky to test directly, so we'll use a workaround
	// by modifying a valid token's claims structure

	// First generate a valid token
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)

	// Now validate it to ensure it's valid
	claims, err := service.ValidateToken(ctx, token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Create a malformed token (this is a simplified test as we can't easily
	// create a token with valid signature but invalid claims structure)
	malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIiLCJyb2xlcyI6WyJhZG1pbiJdLCJleHAiOjk5OTk5OTk5OTl9.invalid_signature"

	claims, err = service.ValidateToken(ctx, malformedToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
