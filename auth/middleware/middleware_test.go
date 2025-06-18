// Copyright (c) 2025 A Bit of Help, Inc.

package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestNewMiddleware tests the creation of a new middleware
func TestNewMiddleware(t *testing.T) {
	// Create a real JWT service for testing
	logger := zap.NewNop()
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	config := middleware.Config{
		SkipPaths:   []string{"/health", "/metrics"},
		RequireAuth: true,
	}

	// Test with logger
	mw := middleware.NewMiddleware(jwtService, config, logger)
	assert.NotNil(t, mw)

	// Test with nil logger (should use NopLogger)
	mw = middleware.NewMiddleware(jwtService, config, nil)
	assert.NotNil(t, mw)
}

// TestHandler_SkipPaths tests that paths in SkipPaths are not authenticated
func TestHandler_SkipPaths(t *testing.T) {
	logger := zap.NewNop()
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	config := middleware.Config{
		SkipPaths:   []string{"/health", "/metrics"},
		RequireAuth: true,
	}
	mw := middleware.NewMiddleware(jwtService, config, logger)

	// Create a test handler that will be wrapped by the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a request to a path that should skip authentication
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Call the middleware
	handler := mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

// TestHandler_NoAuthHeader tests behavior when no Authorization header is present
func TestHandler_NoAuthHeader(t *testing.T) {
	logger := zap.NewNop()
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	// Test with RequireAuth = true
	config := middleware.Config{
		SkipPaths:   []string{"/health"},
		RequireAuth: true,
	}
	mw := middleware.NewMiddleware(jwtService, config, logger)

	// Create a test handler that will be wrapped by the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a request without an Authorization header
	req := httptest.NewRequest("GET", "/api/resource", nil)
	w := httptest.NewRecorder()

	// Call the middleware
	handler := mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was rejected with 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test with RequireAuth = false
	config = middleware.Config{
		SkipPaths:   []string{"/health"},
		RequireAuth: false,
	}
	mw = middleware.NewMiddleware(jwtService, config, logger)

	// Create a new request and response recorder
	req = httptest.NewRequest("GET", "/api/resource", nil)
	w = httptest.NewRecorder()

	// Call the middleware
	handler = mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

// TestHandler_InvalidAuthHeader tests behavior with an invalid Authorization header
func TestHandler_InvalidAuthHeader(t *testing.T) {
	logger := zap.NewNop()
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	config := middleware.Config{
		RequireAuth: true,
	}
	mw := middleware.NewMiddleware(jwtService, config, logger)

	// Create a test handler that will be wrapped by the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a request with an invalid Authorization header
	req := httptest.NewRequest("GET", "/api/resource", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	// Call the middleware
	handler := mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was rejected with 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestHandler_ValidJWTToken tests behavior with a valid JWT token
func TestHandler_ValidJWTToken(t *testing.T) {
	logger := zap.NewNop()
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	config := middleware.Config{
		RequireAuth: true,
	}
	mw := middleware.NewMiddleware(jwtService, config, logger)

	// Generate a valid token
	ctx := context.Background()
	userID := "user123"
	roles := []string{"admin", "user"}
	token, err := jwtService.GenerateToken(ctx, userID, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Create a test handler that will be wrapped by the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that user info was added to the context
		userIDFromCtx, ok := middleware.GetUserID(r.Context())
		assert.True(t, ok)
		assert.Equal(t, userID, userIDFromCtx)

		rolesFromCtx, ok := middleware.GetUserRoles(r.Context())
		assert.True(t, ok)
		assert.Equal(t, roles, rolesFromCtx)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a request with a valid Authorization header
	req := httptest.NewRequest("GET", "/api/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Call the middleware
	handler := mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

// TestHandler_InvalidJWTToken tests behavior with an invalid JWT token
func TestHandler_InvalidJWTToken(t *testing.T) {
	logger := zap.NewNop()
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	config := middleware.Config{
		RequireAuth: true,
	}
	mw := middleware.NewMiddleware(jwtService, config, logger)

	// Create a test handler that will be wrapped by the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a request with an invalid token
	req := httptest.NewRequest("GET", "/api/resource", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	w := httptest.NewRecorder()

	// Call the middleware
	handler := mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was rejected with 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestHandler_ExpiredJWTToken tests behavior with an expired JWT token
func TestHandler_ExpiredJWTToken(t *testing.T) {
	logger := zap.NewNop()

	// Create a JWT service with negative token duration to generate expired tokens
	jwtService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: -1 * time.Hour, // Negative duration to create expired token
		Issuer:        "test-issuer",
	}, logger)

	// Generate an expired token
	ctx := context.Background()
	userID := "user123"
	roles := []string{"admin", "user"}
	expiredToken, err := jwtService.GenerateToken(ctx, userID, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, expiredToken)

	// Create a middleware with a different JWT service (with normal duration)
	validJWTService := jwt.NewService(jwt.Config{
		SecretKey:     "test-secret",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}, logger)

	config := middleware.Config{
		RequireAuth: true,
	}
	mw := middleware.NewMiddleware(validJWTService, config, logger)

	// Create a test handler that will be wrapped by the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a request with an expired token
	req := httptest.NewRequest("GET", "/api/resource", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	w := httptest.NewRecorder()

	// Call the middleware
	handler := mw.Handler(nextHandler)
	handler.ServeHTTP(w, req)

	// Check that the request was rejected with 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token expired")
}

// TestContextFunctions tests the context helper functions
func TestContextFunctions(t *testing.T) {
	// Test WithUserID and GetUserID
	ctx := context.Background()
	userID := "user123"
	ctx = middleware.WithUserID(ctx, userID)

	retrievedID, ok := middleware.GetUserID(ctx)
	assert.True(t, ok)
	assert.Equal(t, userID, retrievedID)

	// Test WithUserRoles and GetUserRoles
	roles := []string{"admin", "user"}
	ctx = middleware.WithUserRoles(ctx, roles)

	retrievedRoles, ok := middleware.GetUserRoles(ctx)
	assert.True(t, ok)
	assert.Equal(t, roles, retrievedRoles)

	// Test IsAuthenticated
	assert.True(t, middleware.IsAuthenticated(ctx))

	// Test with empty context
	emptyCtx := context.Background()
	_, ok = middleware.GetUserID(emptyCtx)
	assert.False(t, ok)

	_, ok = middleware.GetUserRoles(emptyCtx)
	assert.False(t, ok)

	assert.False(t, middleware.IsAuthenticated(emptyCtx))
}
