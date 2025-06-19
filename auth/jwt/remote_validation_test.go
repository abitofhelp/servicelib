// Copyright (c) 2025 A Bit of Help, Inc.

package jwt_test

import (
	"context"
	"testing"
	"time"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestRemoteValidationSuccess tests the case where remote validation succeeds
func TestRemoteValidationSuccess(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)
	ctx := context.Background()

	// Create expected claims
	expectedClaims := &jwt.Claims{
		UserID:    "remote-user-123",
		Roles:     []string{"admin", "remote-user"},
		Scopes:    []string{"read", "write"},
		Resources: []string{"resource1", "resource2"},
	}

	// Set up a mock validator that succeeds
	mockValidator := NewSuccessfulMockValidator(expectedClaims)

	// Set the mock validator as the remote validator
	service.WithRemoteValidator(jwt.RemoteConfig{})
	// Access the private field using reflection
	setRemoteValidator(t, service, mockValidator)

	// Test with any token (the mock will ignore the actual token)
	claims, err := service.ValidateToken(ctx, "any-token")

	// Verify the mock was called
	assert.True(t, mockValidator.Called)

	// Verify the result
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, expectedClaims.UserID, claims.UserID)
	assert.Equal(t, expectedClaims.Roles, claims.Roles)
	assert.Equal(t, expectedClaims.Scopes, claims.Scopes)
	assert.Equal(t, expectedClaims.Resources, claims.Resources)
}

// TestRemoteValidationFailure tests the case where remote validation fails with a non-NotImplemented error
func TestRemoteValidationFailure(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)
	ctx := context.Background()

	// Generate a valid token for local validation
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles, []string{}, []string{})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Set up a mock validator that fails with a custom error
	customError := autherrors.WithMessage(autherrors.ErrInvalidToken, "remote validation failed")
	mockValidator := NewFailingMockValidator(customError)

	// Set the mock validator as the remote validator
	service.WithRemoteValidator(jwt.RemoteConfig{})
	// Access the private field using reflection
	setRemoteValidator(t, service, mockValidator)

	// Test validation - should fall back to local validation and succeed
	claims, err := service.ValidateToken(ctx, token)

	// Verify the mock was called
	assert.True(t, mockValidator.Called)

	// Verify the result - should succeed with local validation
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
}

// TestRemoteValidationNotImplemented tests the case where remote validation returns ErrNotImplemented
func TestRemoteValidationNotImplemented(t *testing.T) {
	logger := zap.NewNop()
	config := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := jwt.NewService(config, logger)
	ctx := context.Background()

	// Generate a valid token for local validation
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := service.GenerateToken(ctx, userID, roles, []string{}, []string{})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Set up a mock validator that fails with ErrNotImplemented
	mockValidator := NewFailingMockValidator(autherrors.ErrNotImplemented)

	// Set the mock validator as the remote validator
	service.WithRemoteValidator(jwt.RemoteConfig{})
	// Access the private field using reflection
	setRemoteValidator(t, service, mockValidator)

	// Test validation - should fall back to local validation and succeed
	claims, err := service.ValidateToken(ctx, token)

	// Verify the mock was called
	assert.True(t, mockValidator.Called)

	// Verify the result - should succeed with local validation
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
}

// Helper function to set the remote validator for testing
func setRemoteValidator(t *testing.T, service *jwt.Service, validator jwt.TokenValidator) {
	// Use the exported method to set the remote validator for testing
	service.SetRemoteValidatorForTesting(validator)
}
