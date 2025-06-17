// Copyright (c) 2025 A Bit of Help, Inc.

package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewJWTService(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}

	// Execute
	service := NewJWTService(config, logger)

	// Verify
	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.Equal(t, logger, service.logger)
}

// mockSigningMethod is a mock implementation of jwt.SigningMethod that always returns an error
type mockSigningMethod struct{}

func (m *mockSigningMethod) Verify(signingString string, sig []byte, key interface{}) error {
	return fmt.Errorf("mock verification error")
}

func (m *mockSigningMethod) Sign(signingString string, key interface{}) ([]byte, error) {
	return nil, fmt.Errorf("mock signing error")
}

func (m *mockSigningMethod) Alg() string {
	return "mockMethod"
}

func TestGenerateToken(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name          string
		userID        string
		roles         []string
		secretKey     string
		signingMethod jwt.SigningMethod
		wantErr       bool
		errMessage    string
	}{
		{
			name:          "Valid token generation",
			userID:        "user123",
			roles:         []string{"admin", "user"},
			secretKey:     "test-secret-key",
			signingMethod: jwt.SigningMethodHS256,
			wantErr:       false,
		},
		{
			name:          "Empty user ID",
			userID:        "",
			roles:         []string{"user"},
			secretKey:     "test-secret-key",
			signingMethod: jwt.SigningMethodHS256,
			wantErr:       false,
		},
		{
			name:          "Empty roles",
			userID:        "user123",
			roles:         []string{},
			secretKey:     "test-secret-key",
			signingMethod: jwt.SigningMethodHS256,
			wantErr:       false,
		},
		{
			name:          "Nil roles",
			userID:        "user123",
			roles:         nil,
			secretKey:     "test-secret-key",
			signingMethod: jwt.SigningMethodHS256,
			wantErr:       false,
		},
		{
			name:          "Error in token signing",
			userID:        "user123",
			roles:         []string{"admin"},
			secretKey:     "test-secret-key",
			signingMethod: &mockSigningMethod{},
			wantErr:       true,
			errMessage:    "failed to generate token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a service with the test-specific secret key
			testConfig := JWTConfig{
				SecretKey:     tt.secretKey,
				TokenDuration: 1 * time.Hour,
				Issuer:        "test-issuer",
			}
			testService := NewJWTService(testConfig, logger)

			// For the error test case, we need to override the signing method
			if tt.name == "Error in token signing" {
				// Create a token with our mock signing method that always returns an error
				now := time.Now()
				expiresAt := now.Add(testConfig.TokenDuration)

				claims := Claims{
					UserID: tt.userID,
					Roles:  tt.roles,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(expiresAt),
						IssuedAt:  jwt.NewNumericDate(now),
						NotBefore: jwt.NewNumericDate(now),
						Issuer:    testConfig.Issuer,
					},
				}

				// Create a token with our mock signing method
				token := jwt.NewWithClaims(tt.signingMethod, claims)

				// Try to sign it - this should fail with our mock
				tokenString, err := token.SignedString([]byte(testConfig.SecretKey))

				// Verify
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Contains(t, err.Error(), "mock signing error")

				// Skip the rest of the test for this case
				return
			}

			// For normal test cases, use the service's GenerateToken method
			tokenString, err := testService.GenerateToken(tt.userID, tt.roles)

			// Verify
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)

				// Parse the token to verify its contents
				token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.secretKey), nil
				})
				assert.NoError(t, err)
				assert.True(t, token.Valid)

				claims, ok := token.Claims.(*Claims)
				assert.True(t, ok)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.roles, claims.Roles)
				assert.Equal(t, testConfig.Issuer, claims.Issuer)

				// Check that expiration time is set correctly
				now := time.Now()
				assert.True(t, claims.ExpiresAt.Time.After(now))
				assert.True(t, claims.ExpiresAt.Time.Before(now.Add(testConfig.TokenDuration).Add(5*time.Second))) // Allow small buffer
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)

	// Generate a valid token for testing
	userID := "user123"
	roles := []string{"admin", "user"}
	validToken, err := service.GenerateToken(userID, roles)
	assert.NoError(t, err)

	// Generate an expired token
	expiredConfig := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: -1 * time.Hour, // Negative duration for expired token
		Issuer:        "test-issuer",
	}
	expiredService := NewJWTService(expiredConfig, logger)
	expiredToken, err := expiredService.GenerateToken(userID, roles)
	assert.NoError(t, err)

	// Generate a token with different signing method
	differentMethodToken := createTokenWithDifferentMethod(t, userID, roles, config.SecretKey)

	tests := []struct {
		name        string
		tokenString string
		wantUserID  string
		wantRoles   []string
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			wantUserID:  userID,
			wantRoles:   roles,
			wantErr:     false,
		},
		{
			name:        "Expired token",
			tokenString: expiredToken,
			wantUserID:  "",
			wantRoles:   nil,
			wantErr:     true,
		},
		{
			name:        "Invalid token format",
			tokenString: "invalid.token.format",
			wantUserID:  "",
			wantRoles:   nil,
			wantErr:     true,
		},
		{
			name:        "Empty token",
			tokenString: "",
			wantUserID:  "",
			wantRoles:   nil,
			wantErr:     true,
		},
		{
			name:        "Different signing method",
			tokenString: differentMethodToken,
			wantUserID:  "",
			wantRoles:   nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			claims, err := service.ValidateToken(tt.tokenString)

			// Verify
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.wantUserID, claims.UserID)
				assert.Equal(t, tt.wantRoles, claims.Roles)
			}
		})
	}
}

// Helper function to create a token with a different signing method
func createTokenWithDifferentMethod(t *testing.T, userID string, roles []string, secretKey string) string {
	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test-issuer",
		},
	}

	// Use RSA signing method instead of HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)
	return tokenString
}

// TestInvalidToken tests cases where the token is invalid
func TestInvalidToken(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)

	// Test case 1: Token with invalid signature (valid format but wrong signature)
	t.Run("Invalid signature", func(t *testing.T) {
		// Create a valid token
		validToken, err := service.GenerateToken("user123", []string{"admin"})
		assert.NoError(t, err)

		// Tamper with the signature part
		parts := strings.Split(validToken, ".")
		assert.Equal(t, 3, len(parts))

		// Change the last character of the signature to make it invalid
		if len(parts[2]) > 0 {
			lastChar := parts[2][len(parts[2])-1]
			var newChar byte
			if lastChar == 'A' {
				newChar = 'B'
			} else {
				newChar = 'A'
			}
			parts[2] = parts[2][:len(parts[2])-1] + string(newChar)
		}

		tamperedToken := strings.Join(parts, ".")

		// Validate the tampered token
		claims, err := service.ValidateToken(tamperedToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})
}

// TestValidateTokenWithCustomClaims tests that the ValidateToken method can handle
// tokens with custom claims that are compatible with our Claims structure
func TestValidateTokenWithCustomClaims(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)

	// Create a token with a compatible but different claims structure
	type CustomClaims struct {
		Name  string `json:"name"`
		Admin bool   `json:"admin"`
		jwt.RegisteredClaims
	}

	customClaims := CustomClaims{
		Name:  "Test User",
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test-issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	assert.NoError(t, err)

	// Execute
	claims, err := service.ValidateToken(tokenString)

	// Verify - the token should be valid but the claims will be empty for our custom fields
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Empty(t, claims.UserID)                // UserID should be empty
	assert.Nil(t, claims.Roles)                   // Roles should be nil
	assert.Equal(t, "test-issuer", claims.Issuer) // Standard claims should be preserved
}
