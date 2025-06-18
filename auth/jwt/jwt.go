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
}

// NewService creates a new JWT service with the provided configuration and logger.
func NewService(config Config, logger *zap.Logger) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Service{
		config: config,
		logger: logger,
		tracer: otel.Tracer("auth.jwt"),
	}
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

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.WithContext(errors.ErrInvalidSignature, "alg", token.Header["alg"])
			err = errors.WithOp(err, "jwt.Service.ValidateToken")
			err = errors.WithMessage(err, "unexpected signing method")
			return nil, err
		}

		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		// Check for specific error types
		var errMsg string
		var baseErr error

		switch {
		case stderrors.Is(err, jwt.ErrTokenExpired):
			baseErr = errors.ErrExpiredToken
			errMsg = "token has expired"
		case stderrors.Is(err, jwt.ErrTokenSignatureInvalid):
			baseErr = errors.ErrInvalidSignature
			errMsg = "invalid token signature"
		case stderrors.Is(err, jwt.ErrTokenMalformed):
			baseErr = errors.ErrInvalidToken
			errMsg = "malformed token"
		default:
			baseErr = errors.ErrInvalidToken
			errMsg = "invalid token"
		}

		err = errors.WithContext(baseErr, "error", err.Error())
		err = errors.WithOp(err, "jwt.Service.ValidateToken")
		err = errors.WithMessage(err, errMsg)
		s.logger.Debug("Failed to parse token", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		err := errors.WithOp(errors.ErrInvalidToken, "jwt.Service.ValidateToken")
		err = errors.WithMessage(err, "token is not valid")
		s.logger.Debug("Invalid token")
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		err := errors.WithOp(errors.ErrInvalidClaims, "jwt.Service.ValidateToken")
		err = errors.WithMessage(err, "failed to extract claims from token")
		s.logger.Debug("Failed to extract claims from token")
		return nil, err
	}

	if claims.UserID == "" {
		err := errors.WithOp(errors.ErrInvalidClaims, "jwt.Service.ValidateToken")
		err = errors.WithMessage(err, "user ID is missing from token claims")
		s.logger.Debug("User ID is missing from token claims")
		return nil, err
	}

	s.logger.Debug("Token validated successfully", zap.String("user_id", claims.UserID))
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
