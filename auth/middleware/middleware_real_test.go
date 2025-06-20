// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestRealMiddlewareHandler tests the real middleware's Handler function
func TestRealMiddlewareHandler(t *testing.T) {
	// Create a real JWT service
	jwtConfig := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	jwtService := jwt.NewService(jwtConfig, zap.NewNop())

	// Generate a valid token
	ctx := context.Background()
	userID := "user123"
	roles := []string{"admin", "user"}
	scopes := []string{"read", "write"}
	resources := []string{"resource1", "resource2"}
	token, err := jwtService.GenerateToken(ctx, userID, roles, scopes, resources)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Create test cases
	tests := []struct {
		name           string
		skipPaths      []string
		requireAuth    bool
		authHeader     string
		path           string
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "Skip path",
			skipPaths:      []string{"/health"},
			requireAuth:    true,
			authHeader:     "",
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedUserID: "",
		},
		{
			name:           "No auth header, require auth",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "",
			path:           "/test",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name:           "No auth header, don't require auth",
			skipPaths:      []string{},
			requireAuth:    false,
			authHeader:     "",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedUserID: "",
		},
		{
			name:           "Invalid auth header format",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "InvalidFormat",
			path:           "/test",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name:           "Valid auth header, validation succeeds",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "Bearer " + token,
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedUserID: userID,
		},
		{
			name:           "Valid auth header, invalid token",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "Bearer invalid-token",
			path:           "/test",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a middleware
			middleware := NewMiddleware(jwtService, Config{
				SkipPaths:   tt.skipPaths,
				RequireAuth: tt.requireAuth,
			}, zap.NewNop())

			// Create test handler
			var userIDFromContext string
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, _ := GetUserID(r.Context())
				userIDFromContext = userID
				w.WriteHeader(http.StatusOK)
			})

			// Create test request
			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create test response recorder
			rr := httptest.NewRecorder()

			// Call the middleware
			handler := middleware.Handler(nextHandler)
			handler.ServeHTTP(rr, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedUserID, userIDFromContext)
		})
	}
}

// TestRealMiddlewareWithOIDC tests the middleware with OIDC support
func TestRealMiddlewareWithOIDC(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	// Create a real JWT service
	jwtConfig := jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	jwtService := jwt.NewService(jwtConfig, zap.NewNop())

	// Create a real OIDC service (this would fail in a real test without a real OIDC provider)
	oidcConfig := oidc.Config{
		IssuerURL:     "https://accounts.google.com",
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/callback",
		Scopes:        []string{"openid", "profile", "email"},
		AdminRoleName: "admin",
		Timeout:       10 * time.Second,
	}
	oidcService, err := oidc.NewService(context.Background(), oidcConfig, zap.NewNop())
	if err != nil {
		t.Skip("Skipping test due to OIDC service creation error:", err)
	}

	// Create a middleware with OIDC support
	middleware := NewMiddlewareWithOIDC(jwtService, oidcService, Config{
		SkipPaths:   []string{"/health"},
		RequireAuth: true,
	}, zap.NewNop())

	// Create test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	// Create test response recorder
	rr := httptest.NewRecorder()

	// Call the middleware
	handler := middleware.Handler(nextHandler)
	handler.ServeHTTP(rr, req)

	// Check the response (should fail with unauthorized since we're using an invalid token)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// TestTestMiddlewareWithOIDC tests the TestMiddleware with OIDC support
func TestTestMiddlewareWithOIDC(t *testing.T) {
	// Create a test middleware with OIDC support
	middleware := NewTestMiddlewareWithOIDC(
		func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
			return nil, errors.ErrInvalidToken
		},
		func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
			return &jwt.Claims{
				UserID:    "oidc-user",
				Roles:     []string{"admin"},
				Scopes:    []string{"read", "write"},
				Resources: []string{"resource1"},
			}, nil
		},
		Config{
			SkipPaths:   []string{"/health"},
			RequireAuth: true,
		},
		zap.NewNop(),
	)

	// Create test handler
	var userIDFromContext string
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetUserID(r.Context())
		userIDFromContext = userID
		w.WriteHeader(http.StatusOK)
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	// Create test response recorder
	rr := httptest.NewRecorder()

	// Call the middleware
	handler := middleware.Handler(nextHandler)
	handler.ServeHTTP(rr, req)

	// Check the response (should succeed with OIDC validation)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "oidc-user", userIDFromContext)
}

// TestTestMiddlewareHandleAuthError tests the TestMiddleware's handleAuthError method
func TestTestMiddlewareHandleAuthError(t *testing.T) {
	// Create test cases
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing token",
			err:            errors.ErrMissingToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization required\n",
		},
		{
			name:           "Expired token",
			err:            errors.ErrExpiredToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token expired\n",
		},
		{
			name:           "Invalid signature",
			err:            errors.ErrInvalidSignature,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token signature\n",
		},
		{
			name:           "Invalid claims",
			err:            errors.ErrInvalidClaims,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token claims\n",
		},
		{
			name:           "Other error",
			err:            errors.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test middleware
			middleware := NewTestMiddleware(
				func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
					return nil, tt.err
				},
				Config{},
				zap.NewNop(),
			)

			// Create test response recorder
			rr := httptest.NewRecorder()

			// Call handleAuthError
			middleware.handleAuthError(rr, tt.err)

			// Check the response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
