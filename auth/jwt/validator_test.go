// Copyright (c) 2025 A Bit of Help, Inc.

package jwt_test

import (
	"context"
	"testing"
	"time"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewLocalValidator(t *testing.T) {
	// Test with logger
	logger := zap.NewExample()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	validator := jwt.NewLocalValidator(config, logger)
	assert.NotNil(t, validator)

	// Test with nil logger (should use NopLogger)
	validator = jwt.NewLocalValidator(config, nil)
	assert.NotNil(t, validator)
}

func TestLocalValidator_ValidateToken(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	validator := jwt.NewLocalValidator(config, logger)
	ctx := context.Background()

	// Create a JWT service to generate tokens for testing
	service := jwt.NewService(config, logger)

	// Test with valid token
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := validator.ValidateToken(ctx, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, config.Issuer, claims.Issuer)

	// Test with empty token
	claims, err = validator.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrMissingToken)

	// Test with invalid token
	claims, err = validator.ValidateToken(ctx, "invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrInvalidToken)

	// Test with expired token
	expiredConfig := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: -1 * time.Hour, // Negative duration to create expired token
		Issuer:        "test-issuer",
	}
	expiredService := jwt.NewService(expiredConfig, logger)
	expiredToken, err := expiredService.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)

	claims, err = validator.ValidateToken(ctx, expiredToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrExpiredToken)

	// Test with token signed with different key
	differentConfig := jwt.Config{
		SecretKey:     "different-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	differentService := jwt.NewService(differentConfig, logger)
	differentToken, err := differentService.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)

	claims, err = validator.ValidateToken(ctx, differentToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrInvalidSignature)

	// Test with token that has empty user ID
	// Create a custom token with empty user ID
	customClaims := &gojwt.RegisteredClaims{
		Issuer:    config.Issuer,
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}
	token = createCustomToken(t, customClaims, config.SecretKey)

	claims, err = validator.ValidateToken(ctx, token)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrInvalidClaims)

	// Test with malformed token
	malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ"
	claims, err = validator.ValidateToken(ctx, malformedToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrInvalidToken)

	// Test with token using a different signing method
	// Create a token with a different signing method (RS256 instead of HS256)
	differentAlgToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ"
	claims, err = validator.ValidateToken(ctx, differentAlgToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrInvalidToken)
}

func TestNewRemoteValidator(t *testing.T) {
	// Test with logger
	logger := zap.NewExample()
	config := jwt.RemoteConfig{
		ValidationURL: "https://test.com/validate",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		Timeout:       15 * time.Second,
	}
	validator := jwt.NewRemoteValidator(config, logger)
	assert.NotNil(t, validator)

	// Test with nil logger (should use NopLogger)
	validator = jwt.NewRemoteValidator(config, nil)
	assert.NotNil(t, validator)

	// Test with zero timeout (should use default)
	configWithZeroTimeout := jwt.RemoteConfig{
		ValidationURL: "https://test.com/validate",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		Timeout:       0,
	}
	validator = jwt.NewRemoteValidator(configWithZeroTimeout, logger)
	assert.NotNil(t, validator)
}

func TestRemoteValidator_ValidateToken(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.RemoteConfig{
		ValidationURL: "https://test.com/validate",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		Timeout:       15 * time.Second,
	}
	validator := jwt.NewRemoteValidator(config, logger)
	ctx := context.Background()

	// Test with any token (should return not implemented error)
	claims, err := validator.ValidateToken(ctx, "any-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrNotImplemented)

	// Test with empty token
	claims, err = validator.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.ErrorIs(t, err, autherrors.ErrMissingToken)
}

// Helper function to create a custom token with specific claims
func createCustomToken(t *testing.T, claims gojwt.Claims, secretKey string) string {
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)
	return tokenString
}
