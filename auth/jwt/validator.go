// Copyright (c) 2024 A Bit of Help, Inc.

package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	stderrors "errors" // Standard errors package with alias
	"strings"
	"time"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TokenValidator is an interface for validating JWT tokens.
// It allows for different validation strategies (local, remote, etc.)
type TokenValidator interface {
	// ValidateToken validates a JWT token and returns the claims if valid.
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
}

// parseToken parses and validates a JWT token.
func parseToken(tokenString string, secretKey string, signingMethod SigningMethod) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, errors.WithOp(errors.ErrMissingToken, "jwt.parseToken")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method based on the configured method
		var expectedMethod string
		switch signingMethod {
		case SigningMethodHS256, SigningMethodHS384, SigningMethodHS512:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				expectedMethod = string(signingMethod)
				err := errors.WithContext(errors.ErrInvalidSignature, "alg", token.Header["alg"])
				err = errors.WithOp(err, "jwt.parseToken")
				err = errors.WithMessage(err, fmt.Sprintf("unexpected signing method: expected %s", expectedMethod))
				return nil, err
			}
		case SigningMethodRS256, SigningMethodRS384, SigningMethodRS512:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				expectedMethod = string(signingMethod)
				err := errors.WithContext(errors.ErrInvalidSignature, "alg", token.Header["alg"])
				err = errors.WithOp(err, "jwt.parseToken")
				err = errors.WithMessage(err, fmt.Sprintf("unexpected signing method: expected %s", expectedMethod))
				return nil, err
			}
		case SigningMethodES256, SigningMethodES384, SigningMethodES512:
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				expectedMethod = string(signingMethod)
				err := errors.WithContext(errors.ErrInvalidSignature, "alg", token.Header["alg"])
				err = errors.WithOp(err, "jwt.parseToken")
				err = errors.WithMessage(err, fmt.Sprintf("unexpected signing method: expected %s", expectedMethod))
				return nil, err
			}
		default:
			// Default to HS256 for backward compatibility
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				expectedMethod = "HS256"
				err := errors.WithContext(errors.ErrInvalidSignature, "alg", token.Header["alg"])
				err = errors.WithOp(err, "jwt.parseToken")
				err = errors.WithMessage(err, fmt.Sprintf("unexpected signing method: expected %s", expectedMethod))
				return nil, err
			}
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		// Check for specific error types
		var baseErr error
		var errMsg string

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
		err = errors.WithOp(err, "jwt.parseToken")
		err = errors.WithMessage(err, errMsg)
		return nil, err
	}

	if !token.Valid {
		err := errors.WithOp(errors.ErrInvalidToken, "jwt.parseToken")
		err = errors.WithMessage(err, "token is not valid")
		return nil, err
	}

	return token, nil
}

// LocalValidator implements TokenValidator using local validation.
type LocalValidator struct {
	// config contains the JWT configuration parameters
	config Config

	// logger is used for logging token operations and errors
	logger *logging.ContextLogger

	// tracer is used for tracing token operations
	tracer trace.Tracer
}

// NewLocalValidator creates a new local validator with the provided configuration and logger.
func NewLocalValidator(config Config, logger *zap.Logger) *LocalValidator {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &LocalValidator{
		config: config,
		logger: logging.NewContextLogger(logger),
		tracer: otel.Tracer("auth.jwt.local"),
	}
}

// ValidateToken validates a JWT token locally and returns the claims if valid.
func (v *LocalValidator) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	ctx, span := v.tracer.Start(ctx, "jwt.LocalValidator.ValidateToken")
	defer span.End()

	span.SetAttributes(attribute.Int("token.length", len(tokenString)))

	if tokenString == "" {
		err := errors.WithOp(errors.ErrMissingToken, "jwt.LocalValidator.ValidateToken")
		v.logger.Debug(ctx, "Token string is empty")
		return nil, err
	}

	// Use the existing validation logic from the JWT service
	token, err := parseToken(tokenString, v.config.SecretKey, v.config.SigningMethod)
	if err != nil {
		v.logger.Debug(ctx, "Failed to parse token", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		err := errors.WithOp(errors.ErrInvalidClaims, "jwt.LocalValidator.ValidateToken")
		err = errors.WithMessage(err, "failed to extract claims from token")
		v.logger.Debug(ctx, "Failed to extract claims from token")
		return nil, err
	}

	if claims.UserID == "" {
		err := errors.WithOp(errors.ErrInvalidClaims, "jwt.LocalValidator.ValidateToken")
		err = errors.WithMessage(err, "user ID is missing from token claims")
		v.logger.Debug(ctx, "User ID is missing from token claims")
		return nil, err
	}

	v.logger.Debug(ctx, "Token validated successfully", zap.String("user_id", claims.UserID))
	return claims, nil
}

// RemoteValidator implements TokenValidator using remote validation.
type RemoteValidator struct {
	// config contains the remote validation configuration parameters
	config RemoteConfig

	// logger is used for logging token operations and errors
	logger *logging.ContextLogger

	// tracer is used for tracing token operations
	tracer trace.Tracer

	// httpClient is the HTTP client used for remote validation
	httpClient *RemoteClient
}

// RemoteConfig holds the configuration for remote JWT token validation.
type RemoteConfig struct {
	// ValidationURL is the URL of the remote validation endpoint
	ValidationURL string

	// ClientID is the client ID for the remote validation service
	ClientID string

	// ClientSecret is the client secret for the remote validation service
	ClientSecret string

	// Timeout is the timeout for remote validation operations
	Timeout time.Duration
}

// NewRemoteValidator creates a new remote validator with the provided configuration and logger.
func NewRemoteValidator(config RemoteConfig, logger *zap.Logger) *RemoteValidator {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Set default timeout if not provided
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	return &RemoteValidator{
		config:     config,
		logger:     logging.NewContextLogger(logger),
		tracer:     otel.Tracer("auth.jwt.remote"),
		httpClient: NewRemoteClient(config),
	}
}

// ValidateToken validates a JWT token remotely and returns the claims if valid.
func (v *RemoteValidator) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	ctx, span := v.tracer.Start(ctx, "jwt.RemoteValidator.ValidateToken")
	defer span.End()

	span.SetAttributes(attribute.Int("token.length", len(tokenString)))

	if tokenString == "" {
		err := errors.WithOp(errors.ErrMissingToken, "jwt.RemoteValidator.ValidateToken")
		v.logger.Debug(ctx, "Token string is empty")
		return nil, err
	}

	// Create a context with timeout for remote validation
	validationCtx, cancel := context.WithTimeout(ctx, v.config.Timeout)
	defer cancel()

	// Validate the token remotely
	claims, err := v.httpClient.ValidateToken(validationCtx, tokenString)
	if err != nil {
		v.logger.Debug(ctx, "Remote validation failed", zap.Error(err))
		return nil, err
	}

	v.logger.Debug(ctx, "Token validated successfully", zap.String("user_id", claims.UserID))
	return claims, nil
}

// RemoteClient handles HTTP communication with the remote validation service.
type RemoteClient struct {
	config RemoteConfig
	client *http.Client
}

// NewRemoteClient creates a new remote client with the provided configuration.
func NewRemoteClient(config RemoteConfig) *RemoteClient {
	return &RemoteClient{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// ValidationRequest represents a request to validate a token.
type ValidationRequest struct {
	Token string `json:"token"`
}

// ValidationResponse represents a response from the validation endpoint.
type ValidationResponse struct {
	Valid  bool    `json:"valid"`
	UserID string  `json:"user_id"`
	Roles  []string `json:"roles"`
	Scopes []string `json:"scopes"`
	Error  string  `json:"error,omitempty"`
}

// ValidateToken sends a validation request to the remote service.
func (c *RemoteClient) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	// Create the request body
	reqBody := ValidationRequest{
		Token: tokenString,
	}

	// Marshal the request to JSON
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.WithContext(errors.Wrap(err, "failed to marshal validation request"), "token_length", len(tokenString))
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.ValidationURL, strings.NewReader(string(reqJSON)))
	if err != nil {
		return nil, errors.WithContext(errors.Wrap(err, "failed to create validation request"), "url", c.config.ValidationURL)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication if provided
	if c.config.ClientID != "" && c.config.ClientSecret != "" {
		req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)
	}

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithContext(errors.Wrap(err, "failed to send validation request"), "url", c.config.ValidationURL)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithContext(errors.Wrap(err, "failed to read validation response"), "status", resp.Status)
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		err := errors.WithContext(errors.ErrInvalidToken, "status", resp.Status)
		err = errors.WithContext(err, "response", string(respBody))
		return nil, err
	}

	// Parse the response
	var validationResp ValidationResponse
	if err := json.Unmarshal(respBody, &validationResp); err != nil {
		return nil, errors.WithContext(errors.Wrap(err, "failed to parse validation response"), "response", string(respBody))
	}

	// Check if the token is valid
	if !validationResp.Valid {
		if validationResp.Error != "" {
			return nil, errors.WithMessage(errors.ErrInvalidToken, validationResp.Error)
		}
		return nil, errors.WithMessage(errors.ErrInvalidToken, "token is not valid")
	}

	// Create claims from the response
	claims := &Claims{
		UserID:    validationResp.UserID,
		Roles:     validationResp.Roles,
		Scopes:    validationResp.Scopes,
	}

	return claims, nil
}
