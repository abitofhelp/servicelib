// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	stderrors "errors" // Standard errors package with alias
	"net/http"
	"strings"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/jwt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TestMiddleware is a middleware for testing.
// It allows us to inject mock services without using reflection.
type TestMiddleware struct {
	// jwtValidator is the function to call when ValidateToken is called.
	jwtValidator func(ctx context.Context, tokenString string) (*jwt.Claims, error)

	// oidcValidator is the function to call when ValidateToken is called.
	oidcValidator func(ctx context.Context, tokenString string) (*jwt.Claims, error)

	// config is the middleware configuration
	config Config

	// logger is used for logging middleware operations and errors
	logger *zap.Logger

	// tracer is used for tracing middleware operations
	tracer trace.Tracer
}

// NewTestMiddleware creates a new test middleware.
func NewTestMiddleware(
	jwtValidator func(ctx context.Context, tokenString string) (*jwt.Claims, error),
	config Config,
	logger *zap.Logger,
) *TestMiddleware {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &TestMiddleware{
		jwtValidator: jwtValidator,
		config:       config,
		logger:       logger,
		tracer:       otel.Tracer("auth.middleware.test"),
	}
}

// NewTestMiddlewareWithOIDC creates a new test middleware with OIDC support.
func NewTestMiddlewareWithOIDC(
	jwtValidator func(ctx context.Context, tokenString string) (*jwt.Claims, error),
	oidcValidator func(ctx context.Context, tokenString string) (*jwt.Claims, error),
	config Config,
	logger *zap.Logger,
) *TestMiddleware {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &TestMiddleware{
		jwtValidator:  jwtValidator,
		oidcValidator: oidcValidator,
		config:        config,
		logger:        logger,
		tracer:        otel.Tracer("auth.middleware.test"),
	}
}

// Handler is the HTTP middleware function.
// It mimics the behavior of the real middleware but uses the injected validators.
func (m *TestMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := m.tracer.Start(r.Context(), "middleware.Handler")
		defer span.End()

		// Add attributes to the span
		span.SetAttributes(attribute.String("path", r.URL.Path))

		// Check if the path should skip authentication
		for _, path := range m.config.SkipPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				m.logger.Debug("Skipping authentication for path", zap.String("path", r.URL.Path))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if m.config.RequireAuth {
				err := errors.WithOp(errors.ErrMissingToken, "middleware.Handler")
				err = errors.WithMessage(err, "authorization header is required")
				m.logger.Debug("No Authorization header provided")
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// No token provided, continue as unauthenticated
			m.logger.Debug("No Authorization header provided, continuing as unauthenticated")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Extract token from header
		tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			err = errors.WithOp(err, "middleware.Handler")
			m.logger.Debug("Invalid Authorization header format", zap.Error(err))
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		var claims *jwt.Claims

		// Try OIDC validation first if available
		if m.oidcValidator != nil {
			claims, err = m.oidcValidator(ctx, tokenString)
			if err != nil {
				m.logger.Debug("OIDC validation failed, trying JWT", zap.Error(err))
				// Fall back to JWT validation
				claims, err = m.jwtValidator(ctx, tokenString)
				if err != nil {
					m.handleAuthError(w, err)
					return
				}
			}
		} else {
			// Use JWT validation
			claims, err = m.jwtValidator(ctx, tokenString)
			if err != nil {
				m.handleAuthError(w, err)
				return
			}
		}

		// Add user info to context
		ctx = WithUserID(ctx, claims.UserID)
		ctx = WithUserRoles(ctx, claims.Roles)
		ctx = WithUserScopes(ctx, claims.Scopes)
		ctx = WithUserResources(ctx, claims.Resources)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleAuthError handles authentication errors.
func (m *TestMiddleware) handleAuthError(w http.ResponseWriter, err error) {
	status := http.StatusUnauthorized
	message := "Invalid token"

	if stderrors.Is(err, errors.ErrMissingToken) {
		message = "Authorization required"
	} else if stderrors.Is(err, errors.ErrExpiredToken) {
		message = "Token expired"
	} else if stderrors.Is(err, errors.ErrInvalidSignature) {
		message = "Invalid token signature"
	} else if stderrors.Is(err, errors.ErrInvalidClaims) {
		message = "Invalid token claims"
	}

	m.logger.Debug("Authentication error", zap.Error(err), zap.String("message", message))
	http.Error(w, message, status)
}
