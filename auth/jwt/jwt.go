// Copyright (c) 2024 A Bit of Help, Inc.

// Package jwt provides JWT token handling for the auth module.
// It includes functionality for generating and validating JWT tokens.
package jwt

import (
	"context"
	stderrors "errors" // Standard errors package with alias
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// SigningMethod represents the algorithm used to sign JWT tokens
type SigningMethod string

const (
	// SigningMethodHS256 represents HMAC using SHA-256
	SigningMethodHS256 SigningMethod = "HS256"
	// SigningMethodHS384 represents HMAC using SHA-384
	SigningMethodHS384 SigningMethod = "HS384"
	// SigningMethodHS512 represents HMAC using SHA-512
	SigningMethodHS512 SigningMethod = "HS512"
	// SigningMethodRS256 represents RSASSA-PKCS1-v1_5 using SHA-256
	SigningMethodRS256 SigningMethod = "RS256"
	// SigningMethodRS384 represents RSASSA-PKCS1-v1_5 using SHA-384
	SigningMethodRS384 SigningMethod = "RS384"
	// SigningMethodRS512 represents RSASSA-PKCS1-v1_5 using SHA-512
	SigningMethodRS512 SigningMethod = "RS512"
	// SigningMethodES256 represents ECDSA using P-256 and SHA-256
	SigningMethodES256 SigningMethod = "ES256"
	// SigningMethodES384 represents ECDSA using P-384 and SHA-384
	SigningMethodES384 SigningMethod = "ES384"
	// SigningMethodES512 represents ECDSA using P-521 and SHA-512
	SigningMethodES512 SigningMethod = "ES512"
)

// Config holds the configuration for JWT token handling.
type Config struct {
	// SecretKey is the key used to sign and verify JWT tokens
	SecretKey string

	// TokenDuration is the validity period for generated tokens
	TokenDuration time.Duration

	// Issuer identifies the entity that issued the token
	Issuer string

	// SigningMethod is the algorithm used to sign JWT tokens
	// Default is HS256 if not specified
	SigningMethod SigningMethod

	// MinSecretKeyLength is the minimum length required for the secret key
	// Default is 32 if not specified
	MinSecretKeyLength int
}

// Service handles JWT token operations including generation and validation.
type Service struct {
	// config contains the JWT configuration parameters
	config Config

	// logger is used for logging token operations and errors
	logger *zap.Logger

	// tracer is used for tracing token operations
	tracer trace.Tracer

	// localValidator is used for local token validation
	localValidator TokenValidator

	// remoteValidator is used for remote token validation (optional)
	remoteValidator TokenValidator

	// revokedTokens is a map of revoked token IDs to their expiration time
	revokedTokens map[string]time.Time

	// mutex protects the revokedTokens map
	mutex sync.RWMutex
}

// NewService creates a new JWT service with the provided configuration and logger.
func NewService(config Config, logger *zap.Logger) (*Service, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Create a background context for logging
	ctx := context.Background()

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Set default values if not provided
	if config.SigningMethod == "" {
		config.SigningMethod = SigningMethodHS256
	}

	if config.MinSecretKeyLength == 0 {
		config.MinSecretKeyLength = 32
	}

	// Validate secret key length
	if len(config.SecretKey) < config.MinSecretKeyLength {
		err := errors.WithContext(errors.ErrInvalidConfig, "secret_key_length", len(config.SecretKey))
		err = errors.WithOp(err, "jwt.NewService")
		err = errors.WithMessage(err, fmt.Sprintf("secret key must be at least %d characters", config.MinSecretKeyLength))
		contextLogger.Error(ctx, "Invalid JWT configuration", zap.Error(err))
		return nil, err
	}

	// Create local validator
	localValidator := NewLocalValidator(config, logger)

	return &Service{
		config:         config,
		logger:         logger,
		tracer:         otel.Tracer("auth.jwt"),
		localValidator: localValidator,
		revokedTokens:  make(map[string]time.Time),
	}, nil
}

// WithRemoteValidator adds a remote validator to the JWT service.
func (s *Service) WithRemoteValidator(config RemoteConfig) (*Service, error) {
	// Create a background context for logging
	ctx := context.Background()

	// Get a context logger for this operation
	logger := s.getContextLogger(ctx)

	// Validate the ValidationURL
	if config.ValidationURL == "" {
		err := errors.WithContext(errors.ErrInvalidConfig, "validation_url", config.ValidationURL)
		err = errors.WithOp(err, "jwt.Service.WithRemoteValidator")
		err = errors.WithMessage(err, "validation URL cannot be empty")
		logger.Error(ctx, "Invalid remote validator configuration", zap.Error(err))
		return nil, err
	}

	// Parse and validate the URL
	_, err := url.Parse(config.ValidationURL)
	if err != nil {
		err = errors.WithContext(errors.Wrap(err, "invalid validation URL"), "validation_url", config.ValidationURL)
		err = errors.WithOp(err, "jwt.Service.WithRemoteValidator")
		logger.Error(ctx, "Invalid remote validator configuration", zap.Error(err))
		return nil, err
	}

	s.remoteValidator = NewRemoteValidator(config, s.logger)
	return s, nil
}

// SetRemoteValidatorForTesting sets the remote validator for testing purposes.
// This method should only be used in tests.
func (s *Service) SetRemoteValidatorForTesting(validator TokenValidator) {
	s.remoteValidator = validator
}

// getContextLogger returns a ContextLogger for the given context.
// This is a helper method to ensure that wherever a context is passed, the code uses a ContextLogger.
func (s *Service) getContextLogger(ctx context.Context) *logging.ContextLogger {
	return logging.NewContextLogger(s.logger)
}

// Claims represents the JWT claims contained in a token.
type Claims struct {
	// UserID is the unique identifier of the user (stored in the 'sub' claim)
	UserID string `json:"sub"`

	// Roles contains the user's assigned roles for authorization
	Roles []string `json:"roles"`

	// Scopes contains the user's assigned permission scopes
	Scopes []string `json:"scopes"`

	// Resources contains the resources the user has access to
	Resources []string `json:"resources"`

	// RegisteredClaims contains the standard JWT claims like expiration time
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user with the specified roles, scopes, and resources.
func (s *Service) GenerateToken(ctx context.Context, userID string, roles []string, scopes []string, resources []string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "jwt.Service.GenerateToken")
	defer span.End()

	// Get a context logger for this operation
	logger := s.getContextLogger(ctx)

	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.StringSlice("user.roles", roles),
		attribute.StringSlice("user.scopes", scopes),
		attribute.StringSlice("user.resources", resources),
	)

	if userID == "" {
		err := errors.WithContext(errors.ErrInvalidClaims, "user_id", userID)
		err = errors.WithOp(err, "jwt.Service.GenerateToken")
		err = errors.WithMessage(err, "user ID cannot be empty")
		logger.Error(ctx, "Failed to generate token: user ID is empty")
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(s.config.TokenDuration)

	// Generate a unique token ID
	tokenID := fmt.Sprintf("%s-%d", userID, now.UnixNano())

	claims := Claims{
		UserID:    userID,
		Roles:     roles,
		Scopes:    scopes,
		Resources: resources,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
			ID:        tokenID, // Add a unique token ID for revocation
		},
	}

	// Select the signing method based on the configuration
	var signingMethod jwt.SigningMethod
	switch s.config.SigningMethod {
	case SigningMethodHS256:
		signingMethod = jwt.SigningMethodHS256
	case SigningMethodHS384:
		signingMethod = jwt.SigningMethodHS384
	case SigningMethodHS512:
		signingMethod = jwt.SigningMethodHS512
	case SigningMethodRS256:
		signingMethod = jwt.SigningMethodRS256
	case SigningMethodRS384:
		signingMethod = jwt.SigningMethodRS384
	case SigningMethodRS512:
		signingMethod = jwt.SigningMethodRS512
	case SigningMethodES256:
		signingMethod = jwt.SigningMethodES256
	case SigningMethodES384:
		signingMethod = jwt.SigningMethodES384
	case SigningMethodES512:
		signingMethod = jwt.SigningMethodES512
	default:
		signingMethod = jwt.SigningMethodHS256
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		err = errors.WithContext(errors.Wrap(err, "failed to sign token"), "user_id", userID)
		err = errors.WithOp(err, "jwt.Service.GenerateToken")
		logger.Error(ctx, "Failed to generate token", zap.Error(err), zap.String("user_id", userID))
		return "", err
	}

	logger.Debug(ctx, "Token generated successfully", zap.String("user_id", userID), zap.String("token_id", tokenID))
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims if valid.
func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	ctx, span := s.tracer.Start(ctx, "jwt.Service.ValidateToken")
	defer span.End()

	// Get a context logger for this operation
	logger := s.getContextLogger(ctx)

	span.SetAttributes(attribute.String("token.length", strconv.Itoa(len(tokenString))))

	if tokenString == "" {
		err := errors.WithOp(errors.ErrMissingToken, "jwt.Service.ValidateToken")
		logger.Debug(ctx, "Token string is empty")
		return nil, err
	}

	// Try remote validation first if available
	if s.remoteValidator != nil {
		claims, err := s.remoteValidator.ValidateToken(ctx, tokenString)
		if err == nil {
			// Check if the token has been revoked
			if s.isTokenRevoked(claims.ID) {
				err := errors.WithContext(errors.ErrInvalidToken, "token_id", claims.ID)
				err = errors.WithOp(err, "jwt.Service.ValidateToken")
				err = errors.WithMessage(err, "token has been revoked")
				logger.Debug(ctx, "Token has been revoked", zap.String("token_id", claims.ID))
				return nil, err
			}

			logger.Debug(ctx, "Token validated successfully by remote validator", zap.String("user_id", claims.UserID))
			return claims, nil
		}

		// If remote validation fails with a "not implemented" error, log it but don't treat it as a fatal error
		if stderrors.Is(err, errors.ErrNotImplemented) {
			logger.Debug(ctx, "Remote validation not implemented, falling back to local validation")
		} else {
			// For other errors, log the error but still try local validation
			logger.Debug(ctx, "Remote validation failed, falling back to local validation", zap.Error(err))
		}
	}

	// Fall back to local validation
	claims, err := s.localValidator.ValidateToken(ctx, tokenString)
	if err != nil {
		logger.Debug(ctx, "Local validation failed", zap.Error(err))
		return nil, err
	}

	// Check if the token has been revoked
	if s.isTokenRevoked(claims.ID) {
		err := errors.WithContext(errors.ErrInvalidToken, "token_id", claims.ID)
		err = errors.WithOp(err, "jwt.Service.ValidateToken")
		err = errors.WithMessage(err, "token has been revoked")
		logger.Debug(ctx, "Token has been revoked", zap.String("token_id", claims.ID))
		return nil, err
	}

	logger.Debug(ctx, "Token validated successfully by local validator", zap.String("user_id", claims.UserID))
	return claims, nil
}

// RevokeToken revokes a token by its ID.
func (s *Service) RevokeToken(ctx context.Context, tokenID string, expiresAt time.Time) error {
	ctx, span := s.tracer.Start(ctx, "jwt.Service.RevokeToken")
	defer span.End()

	// Get a context logger for this operation
	logger := s.getContextLogger(ctx)

	span.SetAttributes(attribute.String("token.id", tokenID))

	if tokenID == "" {
		err := errors.WithContext(errors.ErrInvalidToken, "token_id", tokenID)
		err = errors.WithOp(err, "jwt.Service.RevokeToken")
		err = errors.WithMessage(err, "token ID cannot be empty")
		logger.Error(ctx, "Failed to revoke token: token ID is empty")
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Add the token to the revoked tokens map
	s.revokedTokens[tokenID] = expiresAt
	logger.Debug(ctx, "Token revoked", zap.String("token_id", tokenID))

	// Clean up expired revoked tokens
	s.cleanupRevokedTokens()

	return nil
}

// RevokeTokenByString revokes a token by its string representation.
func (s *Service) RevokeTokenByString(ctx context.Context, tokenString string) error {
	ctx, span := s.tracer.Start(ctx, "jwt.Service.RevokeTokenByString")
	defer span.End()

	// Get a context logger for this operation
	logger := s.getContextLogger(ctx)

	span.SetAttributes(attribute.String("token.length", strconv.Itoa(len(tokenString))))

	// Parse the token to get its ID and expiration time
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		logger.Error(ctx, "Failed to revoke token: token validation failed", zap.Error(err))
		return err
	}

	// Revoke the token
	return s.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time)
}

// isTokenRevoked checks if a token has been revoked.
func (s *Service) isTokenRevoked(tokenID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, revoked := s.revokedTokens[tokenID]
	return revoked
}

// cleanupRevokedTokens removes expired revoked tokens from the map.
func (s *Service) cleanupRevokedTokens() {
	// Create a background context for logging
	ctx := context.Background()

	// Get a context logger for this operation
	logger := s.getContextLogger(ctx)

	now := time.Now()
	for tokenID, expiresAt := range s.revokedTokens {
		if now.After(expiresAt) {
			delete(s.revokedTokens, tokenID)
			logger.Debug(ctx, "Removed expired revoked token", zap.String("token_id", tokenID))
		}
	}
}

// ExtractTokenFromHeader extracts a JWT token from an Authorization header.
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.ErrMissingToken
	}

	// Check if the header has the Bearer prefix
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		err := errors.WithMessage(errors.ErrInvalidToken, "invalid authorization header format")
		return "", err
	}

	return authHeader[len(prefix):], nil
}
