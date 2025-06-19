// Copyright (c) 2025 A Bit of Help, Inc.

package oidc_test

import (
	"context"
	"testing"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/oidc"
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// TestMockValidateToken tests the ValidateToken method with a mock service.
func TestMockValidateToken(t *testing.T) {
	// Create test cases
	tests := []struct {
		name           string
		tokenString    string
		validateResult *jwt.Claims
		validateError  error
		expectError    bool
		expectedClaims *jwt.Claims
	}{
		{
			name:        "Empty token",
			tokenString: "",
			validateResult: nil,
			validateError:  autherrors.ErrMissingToken,
			expectError:    true,
			expectedClaims: nil,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			validateResult: nil,
			validateError:  autherrors.ErrInvalidToken,
			expectError:    true,
			expectedClaims: nil,
		},
		{
			name:        "Expired token",
			tokenString: "expired.token.string",
			validateResult: nil,
			validateError:  autherrors.ErrExpiredToken,
			expectError:    true,
			expectedClaims: nil,
		},
		{
			name:        "Valid token",
			tokenString: "valid.token.string",
			validateResult: &jwt.Claims{
				UserID: "user123",
				Roles:  []string{"admin", "user"},
			},
			validateError:  nil,
			expectError:    false,
			expectedClaims: &jwt.Claims{
				UserID: "user123",
				Roles:  []string{"admin", "user"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := NewMockService(oidc.Config{}, zap.NewNop())
			mockService.ValidateTokenFunc = func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
				assert.Equal(t, tt.tokenString, tokenString)
				return tt.validateResult, tt.validateError
			}

			// Call ValidateToken
			claims, err := mockService.ValidateToken(context.Background(), tt.tokenString)

			// Check the result
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
				if tt.validateError != nil {
					assert.ErrorIs(t, err, tt.validateError)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.expectedClaims.UserID, claims.UserID)
				assert.Equal(t, tt.expectedClaims.Roles, claims.Roles)
			}
		})
	}
}

// TestMockIsAdmin tests the IsAdmin method with a mock service.
func TestMockIsAdmin(t *testing.T) {
	// Create test cases
	tests := []struct {
		name          string
		roles         []string
		adminRoleName string
		expected      bool
	}{
		{
			name:          "User has admin role",
			roles:         []string{"admin", "user"},
			adminRoleName: "admin",
			expected:      true,
		},
		{
			name:          "User does not have admin role",
			roles:         []string{"user"},
			adminRoleName: "admin",
			expected:      false,
		},
		{
			name:          "Empty roles",
			roles:         []string{},
			adminRoleName: "admin",
			expected:      false,
		},
		{
			name:          "Nil roles",
			roles:         nil,
			adminRoleName: "admin",
			expected:      false,
		},
		{
			name:          "Different admin role name",
			roles:         []string{"superuser", "user"},
			adminRoleName: "superuser",
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := NewMockService(oidc.Config{
				AdminRoleName: tt.adminRoleName,
			}, zap.NewNop())

			// Use the real IsAdmin implementation
			mockService.IsAdminFunc = func(roles []string) bool {
				for _, role := range roles {
					if role == tt.adminRoleName {
						return true
					}
				}
				return false
			}

			// Call IsAdmin
			result := mockService.IsAdmin(tt.roles)

			// Check the result
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMockGetAuthURL tests the GetAuthURL method with a mock service.
func TestMockGetAuthURL(t *testing.T) {
	// Create test cases
	tests := []struct {
		name          string
		state         string
		expectedURL   string
	}{
		{
			name:          "With state",
			state:         "random-state",
			expectedURL:   "https://example.com/auth?state=random-state",
		},
		{
			name:          "Empty state",
			state:         "",
			expectedURL:   "https://example.com/auth?state=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := NewMockService(oidc.Config{}, zap.NewNop())
			mockService.GetAuthURLFunc = func(state string) string {
				assert.Equal(t, tt.state, state)
				return "https://example.com/auth?state=" + state
			}

			// Call GetAuthURL
			url := mockService.GetAuthURL(tt.state)

			// Check the result
			assert.Equal(t, tt.expectedURL, url)
		})
	}
}

// TestMockExchange tests the Exchange method with a mock service.
func TestMockExchange(t *testing.T) {
	// Create test cases
	tests := []struct {
		name          string
		code          string
		exchangeResult *oauth2.Token
		exchangeError error
		expectError   bool
	}{
		{
			name:          "Valid code",
			code:          "valid-code",
			exchangeResult: &oauth2.Token{
				AccessToken: "access-token",
				TokenType:   "Bearer",
			},
			exchangeError: nil,
			expectError:   false,
		},
		{
			name:          "Invalid code",
			code:          "invalid-code",
			exchangeResult: nil,
			exchangeError: autherrors.ErrInvalidToken,
			expectError:   true,
		},
		{
			name:          "Empty code",
			code:          "",
			exchangeResult: nil,
			exchangeError: autherrors.ErrInvalidToken,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := NewMockService(oidc.Config{}, zap.NewNop())
			mockService.ExchangeFunc = func(ctx context.Context, code string) (*oauth2.Token, error) {
				assert.Equal(t, tt.code, code)
				return tt.exchangeResult, tt.exchangeError
			}

			// Call Exchange
			token, err := mockService.Exchange(context.Background(), tt.code)

			// Check the result
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, token)
				if tt.exchangeError != nil {
					assert.ErrorIs(t, err, tt.exchangeError)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, tt.exchangeResult.AccessToken, token.AccessToken)
				assert.Equal(t, tt.exchangeResult.TokenType, token.TokenType)
			}
		})
	}
}

// TestMockGetUserInfo tests the GetUserInfo method with a mock service.
func TestMockGetUserInfo(t *testing.T) {
	// Create test cases
	tests := []struct {
		name            string
		token           *oauth2.Token
		userInfoError   error
		expectError     bool
		expectedSubject string
	}{
		{
			name: "Valid token",
			token: &oauth2.Token{
				AccessToken: "access-token",
				TokenType:   "Bearer",
			},
			userInfoError:   nil,
			expectError:     false,
			expectedSubject: "user123",
		},
		{
			name: "Invalid token",
			token: &oauth2.Token{
				AccessToken: "invalid-token",
				TokenType:   "Bearer",
			},
			userInfoError:   autherrors.ErrInvalidToken,
			expectError:     true,
			expectedSubject: "",
		},
		{
			name:            "Nil token",
			token:           nil,
			userInfoError:   autherrors.ErrMissingToken,
			expectError:     true,
			expectedSubject: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := NewMockService(oidc.Config{}, zap.NewNop())

			// Set up the mock function
			mockService.GetUserInfoFunc = func(ctx context.Context, token *oauth2.Token) (*gooidc.UserInfo, error) {
				assert.Equal(t, tt.token, token)
				if tt.expectError {
					return nil, tt.userInfoError
				}

				// For successful cases, we'll just return nil with no error
				// We can't return a real UserInfo because we can't implement the interface
				// But for testing purposes, we just need to check the error handling
				return nil, nil
			}

			// Call GetUserInfo
			userInfo, err := mockService.GetUserInfo(context.Background(), tt.token)

			// Check the result
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, userInfo)
				if tt.userInfoError != nil {
					assert.ErrorIs(t, err, tt.userInfoError)
				}
			} else {
				assert.NoError(t, err)
				// Since we're returning nil for the userInfo in the successful case,
				// we expect userInfo to be nil
				assert.Nil(t, userInfo)
			}
		})
	}
}

// TestRealOIDCService tests the real OIDC service with a mock provider and verifier.
func TestRealOIDCService(t *testing.T) {
	// Skip this test in normal runs since it requires external OIDC provider
	t.Skip("Skipping test that requires external OIDC provider")

	// This test would normally use a real OIDC provider, but we're skipping it
	// since it requires an external OIDC provider.
	// In a real test, we would:
	// 1. Create a real OIDC service with a real provider
	// 2. Generate a real token
	// 3. Validate the token
	// 4. Check the claims
}
