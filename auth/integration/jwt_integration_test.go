// Copyright (c) 2025 A Bit of Help, Inc.

//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestJWTAuthenticationFlow tests the complete JWT authentication flow:
// 1. Generate a token
// 2. Use the token in a request
// 3. Validate the token in middleware
// 4. Extract user information from the token
func TestJWTAuthenticationFlow(t *testing.T) {
	// Create a logger for testing
	logger := zaptest.NewLogger(t)

	// Create auth configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-integration-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-integration-issuer"

	// Create auth instance
	ctx := context.Background()
	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Step 1: Generate a token
	userID := "test-user-123"
	roles := []string{"admin", "user"}
	token, err := authInstance.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Step 2: Create a test server with auth middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 4: Extract user information from the context
		gotUserID, err := authInstance.GetUserID(r.Context())
		if err != nil {
			t.Errorf("Failed to get user ID from context: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		gotRoles, err := authInstance.GetUserRoles(r.Context())
		if err != nil {
			t.Errorf("Failed to get user roles from context: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify the extracted user information
		if gotUserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, gotUserID)
		}

		// Verify roles
		require.Equal(t, len(roles), len(gotRoles), "Role count mismatch")
		for i, role := range roles {
			if i < len(gotRoles) && role != gotRoles[i] {
				t.Errorf("Expected role %s at position %d, got %s", role, i, gotRoles[i])
			}
		}

		// Check if user is admin
		isAdmin, err := authInstance.IsAdmin(r.Context())
		require.NoError(t, err)
		require.True(t, isAdmin, "User should be admin")

		// Check if user is authorized for an operation
		authorized, err := authInstance.IsAuthorized(r.Context(), "read:resource")
		require.NoError(t, err)
		require.True(t, authorized, "User should be authorized for read operation")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	// Apply auth middleware
	handler := authInstance.Middleware()(testHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Step 3: Make a request with the token
	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test with invalid token
	req, err = http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should be unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Test with missing token
	req, err = http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should be unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestJWTTokenValidation tests the token validation functionality
func TestJWTTokenValidation(t *testing.T) {
	// Create a logger for testing
	logger := zaptest.NewLogger(t)

	// Create JWT service directly
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-integration-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-integration-issuer",
	}, logger)

	// Generate a token
	ctx := context.Background()
	userID := "test-user-456"
	roles := []string{"admin", "user"}

	// Create a token directly using the JWT service
	tokenString, err := jwtService.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Validate the token
	validatedClaims, err := jwtService.ValidateToken(ctx, tokenString)
	require.NoError(t, err)
	require.NotNil(t, validatedClaims)

	// Verify the claims
	assert.Equal(t, userID, validatedClaims.UserID)
	assert.Equal(t, roles, validatedClaims.Roles)
	assert.Equal(t, "test-integration-issuer", validatedClaims.Issuer)

	// Test with expired token
	// Create a custom JWT service with a negative token duration to generate an expired token
	expiredJwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-integration-secret-key",
		TokenDuration: -1 * time.Hour, // Negative duration means token is already expired
		Issuer:        "test-integration-issuer",
	}, logger)

	expiredTokenString, err := expiredJwtService.GenerateToken(ctx, userID, roles)
	require.NoError(t, err)

	// Validate the expired token
	_, err = jwtService.ValidateToken(ctx, expiredTokenString)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "token has expired")
}
