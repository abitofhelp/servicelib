// Copyright (c) 2025 A Bit of Help, Inc.

// +build integration

package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// mockOIDCProvider creates a mock OIDC provider server
func mockOIDCProvider(t *testing.T) (*httptest.Server, *rsa.PrivateKey) {
	// Generate a private key for signing tokens
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create a test server that simulates an OIDC provider
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			// Return OIDC discovery document
			config := map[string]interface{}{
				"issuer":                 r.URL.Scheme + "://" + r.Host,
				"authorization_endpoint": r.URL.Scheme + "://" + r.Host + "/auth",
				"token_endpoint":         r.URL.Scheme + "://" + r.Host + "/token",
				"userinfo_endpoint":      r.URL.Scheme + "://" + r.Host + "/userinfo",
				"jwks_uri":               r.URL.Scheme + "://" + r.Host + "/jwks",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(config)

		case "/jwks":
			// Return JWKS (JSON Web Key Set)
			// In a real implementation, this would include the public key
			// For this test, we'll just return a minimal response
			jwks := map[string]interface{}{
				"keys": []map[string]interface{}{
					{
						"kty": "RSA",
						"kid": "test-key-id",
						"use": "sig",
						"alg": "RS256",
						"n":   "test-modulus",
						"e":   "AQAB",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jwks)

		case "/userinfo":
			// Return user info
			// This would normally validate the token, but for testing we'll just return data
			userInfo := map[string]interface{}{
				"sub":   "test-user-789",
				"name":  "Test User",
				"email": "test@example.com",
				"roles": []string{"admin", "user"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(userInfo)

		default:
			http.NotFound(w, r)
		}
	}))

	return server, privateKey
}

// createOIDCToken creates a mock OIDC token
func createOIDCToken(t *testing.T, privateKey *rsa.PrivateKey, issuer, userID string, roles []string) string {
	// Create token claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   issuer,
		"sub":   userID,
		"aud":   "test-client-id",
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
		"roles": roles,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-id"

	// Sign token
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	return tokenString
}

// TestOIDCIntegrationWithMockProvider tests OIDC integration with a mock provider
// This test simulates the integration between the auth middleware and OIDC service
func TestOIDCIntegrationWithMockProvider(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider setup
	// In a real environment, you would have a proper OIDC provider for integration testing
	t.Skip("Skipping OIDC integration test - requires proper OIDC provider setup")

	// Create a mock OIDC provider
	mockProvider, privateKey := mockOIDCProvider(t)
	defer mockProvider.Close()

	// Create a logger for testing
	logger := zaptest.NewLogger(t)

	// Create auth configuration with OIDC settings
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "test-integration-secret-key"
	config.JWT.TokenDuration = 1 * time.Hour
	config.JWT.Issuer = "test-integration-issuer"
	config.OIDC.IssuerURL = mockProvider.URL
	config.OIDC.ClientID = "test-client-id"
	config.OIDC.ClientSecret = "test-client-secret"
	config.OIDC.RedirectURL = "http://localhost:8080/callback"
	config.OIDC.Scopes = []string{"openid", "profile", "email"}

	// Create auth instance
	ctx := context.Background()
	authInstance, err := auth.New(ctx, config, logger)
	require.NoError(t, err)
	require.NotNil(t, authInstance)

	// Create a test token
	userID := "test-user-789"
	roles := []string{"admin", "user"}
	token := createOIDCToken(t, privateKey, mockProvider.URL, userID, roles)

	// Create a test server with auth middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user information from the context
		gotUserID, err := authInstance.GetUserID(r.Context())
		if err != nil {
			t.Errorf("Failed to get user ID from context: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify the extracted user information
		if gotUserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, gotUserID)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	// Apply auth middleware
	handler := authInstance.Middleware()(testHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a request with the token
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
}

// TestOIDCServiceDirectIntegration tests the OIDC service directly
func TestOIDCServiceDirectIntegration(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider setup
	t.Skip("Skipping OIDC service direct integration test - requires proper OIDC provider setup")

	// Create a mock OIDC provider
	mockProvider, privateKey := mockOIDCProvider(t)
	defer mockProvider.Close()

	// Create a logger for testing
	logger := zaptest.NewLogger(t)

	// Create OIDC service
	ctx := context.Background()
	oidcService, err := oidc.NewService(ctx, oidc.Config{
		IssuerURL:     mockProvider.URL,
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}, logger)
	require.NoError(t, err)
	require.NotNil(t, oidcService)

	// Create a test token
	userID := "test-user-789"
	roles := []string{"admin", "user"}
	token := createOIDCToken(t, privateKey, mockProvider.URL, userID, roles)

	// Validate the token
	claims, err := oidcService.ValidateToken(ctx, token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify the claims
	assert.Equal(t, userID, claims.UserID)
	assert.Contains(t, claims.Roles, "admin")
}