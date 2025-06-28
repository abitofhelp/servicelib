// Copyright (c) 2025 A Bit of Help, Inc.

package jwt_test

import (
	"context"
	"errors"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
)

// MockValidator is a mock implementation of the TokenValidator interface for testing.
type MockValidator struct {
	// ShouldSucceed determines whether the validation should succeed or fail
	ShouldSucceed bool

	// ErrorToReturn is the error to return when validation fails
	ErrorToReturn error

	// ClaimsToReturn is the claims to return when validation succeeds
	ClaimsToReturn *jwt.Claims

	// Called tracks whether ValidateToken was called
	Called bool
}

// ValidateToken implements the TokenValidator interface for testing.
func (m *MockValidator) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	m.Called = true

	if tokenString == "" {
		return nil, autherrors.ErrMissingToken
	}

	if m.ShouldSucceed {
		return m.ClaimsToReturn, nil
	}

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	return nil, errors.New("mock validation failed")
}

// NewSuccessfulMockValidator creates a new MockValidator that succeeds with the given claims.
func NewSuccessfulMockValidator(claims *jwt.Claims) *MockValidator {
	return &MockValidator{
		ShouldSucceed:  true,
		ClaimsToReturn: claims,
	}
}

// NewFailingMockValidator creates a new MockValidator that fails with the given error.
func NewFailingMockValidator(err error) *MockValidator {
	return &MockValidator{
		ShouldSucceed: false,
		ErrorToReturn: err,
	}
}
