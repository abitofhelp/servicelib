// Copyright (c) 2025 A Bit of Help, Inc.

// Package oidc provides OpenID Connect integration for the auth module.
// It includes functionality for validating OIDC tokens and extracting user information.
package oidc

import (
	"context"
	"strconv"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/coreos/go-oidc/v3/oidc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Config holds the configuration for OIDC integration.
type Config struct {
	// IssuerURL is the URL of the OIDC provider
	IssuerURL string

	// ClientID is the client ID for the OIDC provider
	ClientID string

	// ClientSecret is the client secret for the OIDC provider
	ClientSecret string

	// RedirectURL is the redirect URL for the OIDC provider
	RedirectURL string

	// Scopes are the OAuth2 scopes to request
	Scopes []string

	// AdminRoleName is the name of the admin role
	AdminRoleName string

	// Timeout is the timeout for OIDC operations
	Timeout time.Duration
}

// Service handles OIDC operations.
type Service struct {
	// config contains the OIDC configuration parameters
	config Config

	// provider is the OIDC provider
	provider *oidc.Provider

	// verifier is the ID token verifier
	verifier *oidc.IDTokenVerifier

	// oauthConfig is the OAuth2 configuration
	oauthConfig *oauth2.Config

	// logger is used for logging OIDC operations and errors
	logger *zap.Logger

	// tracer is used for tracing OIDC operations
	tracer trace.Tracer
}

// NewService creates a new OIDC service with the provided configuration and logger.
func NewService(ctx context.Context, config Config, logger *zap.Logger) (*Service, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Validate context
	if ctx == nil {
		ctx = context.Background()
		logger.Warn("Nil context provided to NewOIDCService, using background context")
	}

	// Validate required configuration
	if config.IssuerURL == "" {
		err := errors.WithMessage(errors.ErrInvalidConfig, "issuer URL cannot be empty")
		logger.Error("Failed to create OIDC service: issuer URL is empty")
		return nil, err
	}

	if config.ClientID == "" {
		err := errors.WithMessage(errors.ErrInvalidConfig, "client ID cannot be empty")
		logger.Error("Failed to create OIDC service: client ID is empty")
		return nil, err
	}

	// Set default timeout if not provided
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
		logger.Debug("Using default OIDC timeout of 10 seconds")
	}

	// Create a context with timeout for OIDC provider discovery
	discoveryCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Discover OIDC provider
	provider, err := oidc.NewProvider(discoveryCtx, config.IssuerURL)
	if err != nil {
		err = errors.WithContext(errors.Wrap(err, "failed to discover OIDC provider"), "issuer_url", config.IssuerURL)
		err = errors.WithOp(err, "oidc.NewService")
		logger.Error("Failed to discover OIDC provider", zap.Error(err), zap.String("issuer_url", config.IssuerURL))
		return nil, err
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	// Create OAuth2 config
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       append([]string{oidc.ScopeOpenID, "profile", "email"}, config.Scopes...),
	}

	return &Service{
		config:      config,
		provider:    provider,
		verifier:    verifier,
		oauthConfig: oauthConfig,
		logger:      logger,
		tracer:      otel.Tracer("auth.oidc"),
	}, nil
}

// ValidateToken validates an OIDC token and returns the claims.
func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	ctx, span := s.tracer.Start(ctx, "oidc.Service.ValidateToken")
	defer span.End()

	span.SetAttributes(attribute.String("token.length", strconv.Itoa(len(tokenString))))

	if tokenString == "" {
		err := errors.WithOp(errors.ErrMissingToken, "oidc.Service.ValidateToken")
		s.logger.Debug("Token string is empty")
		return nil, err
	}

	// Parse and verify the ID token
	idToken, err := s.verifier.Verify(ctx, tokenString)
	if err != nil {
		// Check for specific error types
		var baseErr error
		var errMsg string

		if err.Error() == "oidc: token is expired" {
			baseErr = errors.ErrExpiredToken
			errMsg = "token has expired"
		} else {
			baseErr = errors.ErrInvalidToken
			errMsg = "invalid token"
		}

		err = errors.WithContext(baseErr, "error", err.Error())
		err = errors.WithOp(err, "oidc.Service.ValidateToken")
		err = errors.WithMessage(err, errMsg)
		s.logger.Debug("Failed to verify ID token", zap.Error(err))
		return nil, err
	}

	// Extract claims from the ID token
	var claims struct {
		Subject string   `json:"sub"`
		Email   string   `json:"email"`
		Name    string   `json:"name"`
		Roles   []string `json:"roles"`
	}
	if err := idToken.Claims(&claims); err != nil {
		err = errors.WithContext(errors.ErrInvalidClaims, "error", err.Error())
		err = errors.WithOp(err, "oidc.Service.ValidateToken")
		err = errors.WithMessage(err, "failed to extract claims from ID token")
		s.logger.Debug("Failed to extract claims from ID token", zap.Error(err))
		return nil, err
	}

	if claims.Subject == "" {
		err := errors.WithOp(errors.ErrInvalidClaims, "oidc.Service.ValidateToken")
		err = errors.WithMessage(err, "subject is missing from token claims")
		s.logger.Debug("Subject is missing from token claims")
		return nil, err
	}

	// Create JWT claims
	jwtClaims := &jwt.Claims{
		UserID: claims.Subject,
		Roles:  claims.Roles,
	}

	s.logger.Debug("Token validated successfully", zap.String("user_id", claims.Subject))
	return jwtClaims, nil
}

// IsAdmin checks if the user has the admin role.
func (s *Service) IsAdmin(roles []string) bool {
	for _, role := range roles {
		if role == s.config.AdminRoleName {
			return true
		}
	}
	return false
}

// GetAuthURL returns the URL for the OAuth2 authorization endpoint.
func (s *Service) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// Exchange exchanges an authorization code for a token.
func (s *Service) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	ctx, span := s.tracer.Start(ctx, "oidc.Service.Exchange")
	defer span.End()

	// Create a context with timeout for token exchange
	exchangeCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	token, err := s.oauthConfig.Exchange(exchangeCtx, code)
	if err != nil {
		err = errors.WithContext(errors.Wrap(err, "failed to exchange authorization code"), "code_length", len(code))
		err = errors.WithOp(err, "oidc.Service.Exchange")
		s.logger.Error("Failed to exchange authorization code", zap.Error(err))
		return nil, err
	}

	return token, nil
}

// GetUserInfo gets the user info from the OIDC provider.
func (s *Service) GetUserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	ctx, span := s.tracer.Start(ctx, "oidc.Service.GetUserInfo")
	defer span.End()

	// Create a context with timeout for user info request
	userInfoCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	userInfo, err := s.provider.UserInfo(userInfoCtx, oauth2.StaticTokenSource(token))
	if err != nil {
		err = errors.WithContext(errors.Wrap(err, "failed to get user info"), "token_valid", token.Valid())
		err = errors.WithOp(err, "oidc.Service.GetUserInfo")
		s.logger.Error("Failed to get user info", zap.Error(err))
		return nil, err
	}

	return userInfo, nil
}
