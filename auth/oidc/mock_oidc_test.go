// Copyright (c) 2025 A Bit of Help, Inc.

package oidc_test

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/auth/oidc"
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// MockIDToken is a mock implementation of the OIDC IDToken.
type MockIDToken struct {
	// ClaimsData is the data to return when Claims is called.
	ClaimsData map[string]interface{}
}

// Claims implements the OIDC IDToken interface for testing.
func (m *MockIDToken) Claims(v interface{}) error {
	if m.ClaimsData == nil {
		return errors.New("mock claims not implemented")
	}

	jsonData, err := json.Marshal(m.ClaimsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}

// MockUserInfo is a mock implementation of the OIDC UserInfo.
type MockUserInfo struct {
	// SubjectValue is the subject of the user info.
	SubjectValue string

	// EmailValue is the email of the user info.
	EmailValue string

	// NameValue is the name of the user info.
	NameValue string

	// ClaimsData is the data to return when Claims is called.
	ClaimsData map[string]interface{}
}

// Subject implements the OIDC UserInfo interface for testing.
func (m *MockUserInfo) Subject() string {
	return m.SubjectValue
}

// Email implements the OIDC UserInfo interface for testing.
func (m *MockUserInfo) Email() string {
	return m.EmailValue
}

// Claims implements the OIDC UserInfo interface for testing.
func (m *MockUserInfo) Claims(v interface{}) error {
	if m.ClaimsData == nil {
		// Default implementation
		data := map[string]interface{}{
			"sub":   m.SubjectValue,
			"email": m.EmailValue,
			"name":  m.NameValue,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonData, v)
	}

	jsonData, err := json.Marshal(m.ClaimsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}

// MockService is a mock implementation of the OIDC service for testing.
type MockService struct {
	// Config is the OIDC configuration.
	Config oidc.Config

	// Logger is the logger.
	Logger *zap.Logger

	// Tracer is the tracer.
	Tracer trace.Tracer

	// ValidateTokenFunc is the function to call when ValidateToken is called.
	ValidateTokenFunc func(ctx context.Context, tokenString string) (*jwt.Claims, error)

	// IsAdminFunc is the function to call when IsAdmin is called.
	IsAdminFunc func(roles []string) bool

	// GetAuthURLFunc is the function to call when GetAuthURL is called.
	GetAuthURLFunc func(state string) string

	// ExchangeFunc is the function to call when Exchange is called.
	ExchangeFunc func(ctx context.Context, code string) (*oauth2.Token, error)

	// GetUserInfoFunc is the function to call when GetUserInfo is called.
	GetUserInfoFunc func(ctx context.Context, token *oauth2.Token) (*gooidc.UserInfo, error)
}

// NewMockService creates a new mock OIDC service.
func NewMockService(config oidc.Config, logger *zap.Logger) *MockService {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &MockService{
		Config: config,
		Logger: logger,
		Tracer: otel.Tracer("auth.oidc.mock"),
	}
}

// ValidateToken implements the OIDC service interface for testing.
func (m *MockService) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(ctx, tokenString)
	}
	return nil, errors.New("mock validate token not implemented")
}

// IsAdmin implements the OIDC service interface for testing.
func (m *MockService) IsAdmin(roles []string) bool {
	if m.IsAdminFunc != nil {
		return m.IsAdminFunc(roles)
	}
	return false
}

// GetAuthURL implements the OIDC service interface for testing.
func (m *MockService) GetAuthURL(state string) string {
	if m.GetAuthURLFunc != nil {
		return m.GetAuthURLFunc(state)
	}
	return "https://mock.example.com/auth?state=" + state
}

// Exchange implements the OIDC service interface for testing.
func (m *MockService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if m.ExchangeFunc != nil {
		return m.ExchangeFunc(ctx, code)
	}
	return nil, errors.New("mock exchange not implemented")
}

// GetUserInfo implements the OIDC service interface for testing.
func (m *MockService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*gooidc.UserInfo, error) {
	if m.GetUserInfoFunc != nil {
		return m.GetUserInfoFunc(ctx, token)
	}
	return nil, errors.New("mock get user info not implemented")
}
