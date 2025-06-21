// Copyright (c) 2025 A Bit of Help, Inc.

// Package oidc provides OpenID Connect integration for the auth module.
// It includes functionality for validating OIDC tokens and extracting user information.
package oidc

import (
	"context"
	stderrors "errors"
	"strconv"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"github.com/abitofhelp/servicelib/logging"
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

	// RetryConfig is the configuration for retry operations
	RetryConfig RetryConfig
}

// RetryConfig holds configuration for retry operations.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// InitialBackoff is the initial backoff duration
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration
	MaxBackoff time.Duration

	// BackoffFactor is the factor by which the backoff increases
	BackoffFactor float64
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
		BackoffFactor:  2.0,
	}
}

// retryWithBackoff retries a function with exponential backoff.
// It returns the result of the function if successful, or the last error if all retries fail.
func retryWithBackoff[T any](
	ctx context.Context,
	logger *logging.ContextLogger,
	retryConfig RetryConfig,
	operation string,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var result T
	var lastErr error

	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		// If this is a retry, wait with exponential backoff
		if attempt > 0 {
			backoff := time.Duration(float64(retryConfig.InitialBackoff) * float64(retryConfig.BackoffFactor) * float64(attempt-1))
			if backoff > retryConfig.MaxBackoff {
				backoff = retryConfig.MaxBackoff
			}

			logger.Debug(ctx, "Retrying operation",
				zap.String("operation", operation),
				zap.Int("attempt", attempt),
				zap.Duration("backoff", backoff),
				zap.Error(lastErr))

			select {
			case <-ctx.Done():
				var zero T
				return zero, errors.WithContext(ctx.Err(), "operation", operation)
			case <-time.After(backoff):
				// Continue with retry
			}
		}

		// Execute the function
		var err error
		result, err = fn(ctx)
		if err == nil {
			// Operation succeeded
			if attempt > 0 {
				logger.Debug(ctx, "Operation succeeded after retry",
					zap.String("operation", operation),
					zap.Int("attempt", attempt))
			}
			return result, nil
		}

		lastErr = err
	}

	// If we get here, we've exhausted our retries
	logger.Warn(ctx, "Operation failed after retries",
		zap.String("operation", operation),
		zap.Int("max_retries", retryConfig.MaxRetries),
		zap.Error(lastErr))

	var zero T
	return zero, lastErr
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
	logger *logging.ContextLogger

	// tracer is used for tracing OIDC operations
	tracer trace.Tracer
}

// NewService creates a new OIDC service with the provided configuration and logger.
func NewService(ctx context.Context, config Config, logger *zap.Logger) (*Service, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Create a context logger
	ctxLogger := logging.NewContextLogger(logger)

	// Validate context
	if ctx == nil {
		ctx = context.Background()
		ctxLogger.Warn(ctx, "Nil context provided to NewOIDCService, using background context")
	}

	// Validate required configuration
	if config.IssuerURL == "" {
		err := errors.WithMessage(errors.ErrInvalidConfig, "issuer URL cannot be empty")
		ctxLogger.Error(ctx, "Failed to create OIDC service: issuer URL is empty")
		return nil, err
	}

	if config.ClientID == "" {
		err := errors.WithMessage(errors.ErrInvalidConfig, "client ID cannot be empty")
		ctxLogger.Error(ctx, "Failed to create OIDC service: client ID is empty")
		return nil, err
	}

	// Set default timeout if not provided
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
		ctxLogger.Debug(ctx, "Using default OIDC timeout of 10 seconds")
	}

	// Set default retry config if not provided
	if config.RetryConfig.MaxRetries == 0 {
		config.RetryConfig = DefaultRetryConfig()
		ctxLogger.Debug(ctx, "Using default OIDC retry configuration",
			zap.Int("max_retries", config.RetryConfig.MaxRetries),
			zap.Duration("initial_backoff", config.RetryConfig.InitialBackoff),
			zap.Duration("max_backoff", config.RetryConfig.MaxBackoff),
			zap.Float64("backoff_factor", config.RetryConfig.BackoffFactor))
	}

	// Discover OIDC provider with retries
	provider, err := retryWithBackoff(ctx, ctxLogger, config.RetryConfig, "oidc_provider_discovery", func(retryCtx context.Context) (*oidc.Provider, error) {
		// Create a context with timeout for OIDC provider discovery
		discoveryCtx, cancel := context.WithTimeout(retryCtx, config.Timeout)
		defer cancel()

		// Discover OIDC provider
		provider, err := oidc.NewProvider(discoveryCtx, config.IssuerURL)
		if err != nil {
			err = errors.WithContext(errors.Wrap(err, "failed to discover OIDC provider"), "issuer_url", config.IssuerURL)
			err = errors.WithOp(err, "oidc.NewService")
			return nil, err
		}
		return provider, nil
	})

	if err != nil {
		ctxLogger.Error(ctx, "Failed to discover OIDC provider after retries", zap.Error(err), zap.String("issuer_url", config.IssuerURL))
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
		logger:      logging.NewContextLogger(logger),
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
		s.logger.Debug(ctx, "Token string is empty")
		return nil, err
	}

	// Parse and verify the ID token with retries
	idToken, err := retryWithBackoff(ctx, s.logger, s.config.RetryConfig, "verify_token", func(retryCtx context.Context) (*oidc.IDToken, error) {
		// Verify the token
		idToken, err := s.verifier.Verify(retryCtx, tokenString)
		if err != nil {
			// Check for specific error types
			var baseErr error
			var errMsg string

			if err.Error() == "oidc: token is expired" {
				baseErr = errors.ErrExpiredToken
				errMsg = "token has expired"
				// Don't retry expired tokens
				return nil, errors.WithMessage(errors.WithOp(errors.WithContext(baseErr, "error", err.Error()), "oidc.Service.ValidateToken"), errMsg)
			} else {
				baseErr = errors.ErrInvalidToken
				errMsg = "invalid token"
			}

			err = errors.WithContext(baseErr, "error", err.Error())
			err = errors.WithOp(err, "oidc.Service.ValidateToken")
			err = errors.WithMessage(err, errMsg)
			return nil, err
		}
		return idToken, nil
	})

	if err != nil {
		// If the token is expired, we don't want to log it as an error
		if stderrors.Is(err, errors.ErrExpiredToken) {
			s.logger.Debug(ctx, "Token is expired", zap.Error(err))
		} else {
			s.logger.Debug(ctx, "Failed to verify ID token after retries", zap.Error(err))
		}
		return nil, err
	}

	// Extract claims from the ID token
	var claims struct {
		Subject string   `json:"sub"`
		Email   string   `json:"email"`
		Name    string   `json:"name"`
		Roles   []string `json:"roles"`
	}

	// Extract claims with a single attempt (no need for retries here as it's a local operation)
	if err := idToken.Claims(&claims); err != nil {
		err = errors.WithContext(errors.ErrInvalidClaims, "error", err.Error())
		err = errors.WithOp(err, "oidc.Service.ValidateToken")
		err = errors.WithMessage(err, "failed to extract claims from ID token")
		s.logger.Debug(ctx, "Failed to extract claims from ID token", zap.Error(err))
		return nil, err
	}

	// Validate claims
	if claims.Subject == "" {
		err := errors.WithOp(errors.ErrInvalidClaims, "oidc.Service.ValidateToken")
		err = errors.WithMessage(err, "subject is missing from token claims")
		s.logger.Debug(ctx, "Subject is missing from token claims")
		return nil, err
	}

	// Create JWT claims
	jwtClaims := &jwt.Claims{
		UserID: claims.Subject,
		Roles:  claims.Roles,
	}

	s.logger.Debug(ctx, "Token validated successfully", zap.String("user_id", claims.Subject))
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

	// Exchange authorization code for token with retries
	token, err := retryWithBackoff(ctx, s.logger, s.config.RetryConfig, "token_exchange", func(retryCtx context.Context) (*oauth2.Token, error) {
		// Create a context with timeout for token exchange
		exchangeCtx, cancel := context.WithTimeout(retryCtx, s.config.Timeout)
		defer cancel()

		token, err := s.oauthConfig.Exchange(exchangeCtx, code)
		if err != nil {
			err = errors.WithContext(errors.Wrap(err, "failed to exchange authorization code"), "code_length", len(code))
			err = errors.WithOp(err, "oidc.Service.Exchange")
			return nil, err
		}
		return token, nil
	})

	if err != nil {
		s.logger.Error(ctx, "Failed to exchange authorization code after retries", zap.Error(err))
		return nil, err
	}

	return token, nil
}

// GetUserInfo gets the user info from the OIDC provider.
func (s *Service) GetUserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	ctx, span := s.tracer.Start(ctx, "oidc.Service.GetUserInfo")
	defer span.End()

	// Get user info with retries
	userInfo, err := retryWithBackoff(ctx, s.logger, s.config.RetryConfig, "get_user_info", func(retryCtx context.Context) (*oidc.UserInfo, error) {
		// Create a context with timeout for user info request
		userInfoCtx, cancel := context.WithTimeout(retryCtx, s.config.Timeout)
		defer cancel()

		userInfo, err := s.provider.UserInfo(userInfoCtx, oauth2.StaticTokenSource(token))
		if err != nil {
			err = errors.WithContext(errors.Wrap(err, "failed to get user info"), "token_valid", token.Valid())
			err = errors.WithOp(err, "oidc.Service.GetUserInfo")
			return nil, err
		}
		return userInfo, nil
	})

	if err != nil {
		s.logger.Error(ctx, "Failed to get user info after retries", zap.Error(err))
		return nil, err
	}

	return userInfo, nil
}
