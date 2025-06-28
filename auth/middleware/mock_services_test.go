// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"errors"

	"github.com/abitofhelp/servicelib/auth/jwt"
)

// MockJWTService is a mock implementation of the JWT service for testing.
type MockJWTService struct {
	// ValidateTokenFunc is the function to call when ValidateToken is called.
	ValidateTokenFunc func(ctx context.Context, tokenString string) (*jwt.Claims, error)
}

// ValidateToken implements the JWT service interface for testing.
func (m *MockJWTService) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(ctx, tokenString)
	}
	return nil, errors.New("mock JWT validation not implemented")
}

// NewMockJWTService creates a new mock JWT service.
func NewMockJWTService() *MockJWTService {
	return &MockJWTService{}
}

// MockOIDCService is a mock implementation of the OIDC service for testing.
type MockOIDCService struct {
	// ValidateTokenFunc is the function to call when ValidateToken is called.
	ValidateTokenFunc func(ctx context.Context, tokenString string) (*jwt.Claims, error)
}

// ValidateToken implements the OIDC service interface for testing.
func (m *MockOIDCService) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(ctx, tokenString)
	}
	return nil, errors.New("mock OIDC validation not implemented")
}

// NewMockOIDCService creates a new mock OIDC service.
func NewMockOIDCService() *MockOIDCService {
	return &MockOIDCService{}
}

// NewSuccessfulJWTService creates a mock JWT service that always succeeds with the given claims.
func NewSuccessfulJWTService(claims *jwt.Claims) *MockJWTService {
	return &MockJWTService{
		ValidateTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
			return claims, nil
		},
	}
}

// NewFailingJWTService creates a mock JWT service that always fails with the given error.
func NewFailingJWTService(err error) *MockJWTService {
	return &MockJWTService{
		ValidateTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
			return nil, err
		},
	}
}

// NewSuccessfulOIDCService creates a mock OIDC service that always succeeds with the given claims.
func NewSuccessfulOIDCService(claims *jwt.Claims) *MockOIDCService {
	return &MockOIDCService{
		ValidateTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
			return claims, nil
		},
	}
}

// NewFailingOIDCService creates a mock OIDC service that always fails with the given error.
func NewFailingOIDCService(err error) *MockOIDCService {
	return &MockOIDCService{
		ValidateTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Claims, error) {
			return nil, err
		},
	}
}

// CreateTestClaims creates test claims for testing.
func CreateTestClaims(userID string, roles []string, scopes []string, resources []string) *jwt.Claims {
	return &jwt.Claims{
		UserID:    userID,
		Roles:     roles,
		Scopes:    scopes,
		Resources: resources,
	}
}
