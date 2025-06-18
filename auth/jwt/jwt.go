// Copyright (c) 2025 A Bit of Help, Inc.

// Package jwt provides JWT token handling for the auth module.
// It includes functionality for generating and validating JWT tokens.
package jwt

import (
	"context"
	stderrors "errors" // Standard errors package with alias
	"strconv"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config holds the configuration for JWT token handling.
type Config struct {
	// SecretKey is the key used to sign and verify JWT tokens
	SecretKey string

	// TokenDuration is the validity period for generated tokens
	TokenDuration time.Duration

	// Issuer identifies the entity that issued the token
	Issuer string
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
}

// NewService creates a new JWT service with the provided configuration and logger.
func NewService(config Config, logger *zap.Logger) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Create local validator
	localValidator := NewLocalValidator(config, logger)

	return &Service{
		config:         config,
		logger:         logger,
		tracer:         otel.Tracer("auth.jwt"),
		localValidator: localValidator,
	}
}

// WithRemoteValidator adds a remote validator to the JWT service.
func (s *Service) WithRemoteValidator(config RemoteConfig) *Service {
	s.remoteValidator = NewRemoteValidator(config, s.logger)
	return s
}

// Claims represents the JWT claims contained in a token.
type Claims struct {
	// UserID is the unique identifier of the user (stored in the 'sub' claim)
	UserID string `json:"sub"`

	// Roles contains the user's assigned roles for authorization
	Roles []string `json:"roles"`

	// RegisteredClaims contains the standard JWT claims like expiration time
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user with the specified roles.
func (s *Service) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "jwt.Service.GenerateToken")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.StringSlice("user.roles", roles),
	)

	if userID == "" {
		err := errors.WithContext(errors.ErrInvalidClaims, "user_id", userID)
		err = errors.WithOp(err, "jwt.Service.GenerateToken")
		err = errors.WithMessage(err, "user ID cannot be empty")
		s.logger.Error("Failed to generate token: user ID is empty")
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(s.config.TokenDuration)

	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		err = errors.WithContext(errors.Wrap(err, "failed to sign token"), "user_id", userID)
		err = errors.WithOp(err, "jwt.Service.GenerateToken")
		s.logger.Error("Failed to generate token", zap.Error(err), zap.String("user_id", userID))
		return "", err
	}

	s.logger.Debug("Token generated successfully", zap.String("user_id", userID))
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims if valid.
func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	ctx, span := s.tracer.Start(ctx, "jwt.Service.ValidateToken")
	defer span.End()

	span.SetAttributes(attribute.String("token.length", strconv.Itoa(len(tokenString))))

	if tokenString == "" {
		err := errors.WithOp(errors.ErrMissingToken, "jwt.Service.ValidateToken")
		s.logger.Debug("Token string is empty")
		return nil, err
	}

	// Try remote validation first if available
	if s.remoteValidator != nil {
		claims, err := s.remoteValidator.ValidateToken(ctx, tokenString)
		if err == nil {
			s.logger.Debug("Token validated successfully by remote validator", zap.String("user_id", claims.UserID))
			return claims, nil
		}

		// If remote validation fails with a "not implemented" error, log it but don't treat it as a fatal error
		if stderrors.Is(err, errors.ErrNotImplemented) {
			s.logger.Debug("Remote validation not implemented, falling back to local validation")
		} else {
			// For other errors, log the error but still try local validation
			s.logger.Debug("Remote validation failed, falling back to local validation", zap.Error(err))
		}
	}

	// Fall back to local validation
	claims, err := s.localValidator.ValidateToken(ctx, tokenString)
	if err != nil {
		s.logger.Debug("Local validation failed", zap.Error(err))
		return nil, err
	}

	s.logger.Debug("Token validated successfully by local validator", zap.String("user_id", claims.UserID))
	return claims, nil
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
