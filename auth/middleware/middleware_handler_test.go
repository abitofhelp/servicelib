// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestHandleAuthError tests the handleAuthError method
func TestHandleAuthError(t *testing.T) {
	// Create test cases
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing token",
			err:            autherrors.ErrMissingToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization required\n",
		},
		{
			name:           "Expired token",
			err:            autherrors.ErrExpiredToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token expired\n",
		},
		{
			name:           "Invalid signature",
			err:            autherrors.ErrInvalidSignature,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token signature\n",
		},
		{
			name:           "Invalid claims",
			err:            autherrors.ErrInvalidClaims,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token claims\n",
		},
		{
			name:           "Other error",
			err:            autherrors.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a middleware
			middleware := &Middleware{
				logger: zap.NewNop(),
			}

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

// TestHandler tests the Handler function
func TestHandler(t *testing.T) {
	// Create test cases
	tests := []struct {
		name           string
		skipPaths      []string
		requireAuth    bool
		authHeader     string
		path           string
		validateResult *jwt.Claims
		validateError  error
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "Skip path",
			skipPaths:      []string{"/health"},
			requireAuth:    true,
			authHeader:     "",
			path:           "/health",
			validateResult: nil,
			validateError:  nil,
			expectedStatus: http.StatusOK,
			expectedUserID: "",
		},
		{
			name:           "No auth header, require auth",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "",
			path:           "/test",
			validateResult: nil,
			validateError:  nil,
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name:           "No auth header, don't require auth",
			skipPaths:      []string{},
			requireAuth:    false,
			authHeader:     "",
			path:           "/test",
			validateResult: nil,
			validateError:  nil,
			expectedStatus: http.StatusOK,
			expectedUserID: "",
		},
		{
			name:           "Invalid auth header format",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "InvalidFormat",
			path:           "/test",
			validateResult: nil,
			validateError:  nil,
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name:           "Valid auth header, validation succeeds",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "Bearer valid-token",
			path:           "/test",
			validateResult: &jwt.Claims{UserID: "user123", Roles: []string{"admin"}},
			validateError:  nil,
			expectedStatus: http.StatusOK,
			expectedUserID: "user123",
		},
		{
			name:           "Valid auth header, validation fails",
			skipPaths:      []string{},
			requireAuth:    true,
			authHeader:     "Bearer invalid-token",
			path:           "/test",
			validateResult: nil,
			validateError:  autherrors.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test middleware with a mock JWT validator
			middleware := NewTestMiddleware(
				func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
					return tt.validateResult, tt.validateError
				},
				Config{
					SkipPaths:   tt.skipPaths,
					RequireAuth: tt.requireAuth,
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
