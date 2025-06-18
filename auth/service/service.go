// Copyright (c) 2025 A Bit of Help, Inc.

// Package service provides authorization services for the auth module.
// It includes functionality for checking if a user is authorized to perform specific operations.
package service

import (
	"context"
	"strings"

	"github.com/abitofhelp/servicelib/auth/errors"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config holds the configuration for the authorization service.
type Config struct {
	// AdminRoleName is the name of the admin role
	AdminRoleName string

	// ReadOnlyRoleName is the name of the read-only role
	ReadOnlyRoleName string

	// ReadOperationPrefixes are prefixes for read-only operations
	ReadOperationPrefixes []string
}

// DefaultConfig returns the default configuration for the authorization service.
func DefaultConfig() Config {
	return Config{
		AdminRoleName:    "admin",
		ReadOnlyRoleName: "authuser",
		ReadOperationPrefixes: []string{
			"read:",
			"list:",
			"get:",
			"find:",
			"query:",
			"count:",
		},
	}
}

// Service implements authorization functionality.
type Service struct {
	// config is the service configuration
	config Config

	// logger is used for logging service operations and errors
	logger *zap.Logger

	// tracer is used for tracing service operations
	tracer trace.Tracer
}

// NewService creates a new authorization service.
func NewService(config Config, logger *zap.Logger) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Service{
		config: config,
		logger: logger,
		tracer: otel.Tracer("auth.service"),
	}
}

// IsAuthorized checks if the user is authorized to perform the operation.
func (s *Service) IsAuthorized(ctx context.Context, operation string) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "service.IsAuthorized")
	defer span.End()

	span.SetAttributes(attribute.String("operation", operation))

	// Check if user is authenticated
	if !middleware.IsAuthenticated(ctx) {
		err := errors.WithOp(errors.ErrUnauthorized, "service.IsAuthorized")
		err = errors.WithMessage(err, "user is not authenticated")
		s.logger.Debug("User is not authenticated")
		return false, err
	}

	// Get user roles
	roles, ok := middleware.GetUserRoles(ctx)
	if !ok || len(roles) == 0 {
		err := errors.WithOp(errors.ErrUnauthorized, "service.IsAuthorized")
		err = errors.WithMessage(err, "user has no roles")
		s.logger.Debug("User has no roles")
		return false, err
	}

	// Check if user is admin
	isAdmin, err := s.IsAdmin(ctx)
	if err != nil {
		return false, err
	}

	if isAdmin {
		// Admins can do anything
		s.logger.Debug("User is admin, authorization granted", zap.String("operation", operation))
		return true, nil
	}

	// Check if user has read-only role
	hasReadOnlyRole := false
	for _, role := range roles {
		if role == s.config.ReadOnlyRoleName {
			hasReadOnlyRole = true
			break
		}
	}

	// Check if operation is read-only
	isReadOperation := false
	for _, prefix := range s.config.ReadOperationPrefixes {
		if strings.HasPrefix(operation, prefix) {
			isReadOperation = true
			break
		}
	}

	// Read-only users can only perform read operations
	if hasReadOnlyRole && isReadOperation {
		s.logger.Debug("User has read-only role and operation is read-only, authorization granted",
			zap.String("operation", operation))
		return true, nil
	}

	// If we get here, the user is not authorized
	err = errors.WithOp(errors.ErrForbidden, "service.IsAuthorized")
	err = errors.WithContext(err, "operation", operation)
	err = errors.WithMessage(err, "user is not authorized to perform this operation")
	s.logger.Debug("User is not authorized",
		zap.String("operation", operation),
		zap.Strings("roles", roles))
	return false, err
}

// IsAdmin checks if the user has admin role.
func (s *Service) IsAdmin(ctx context.Context) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "service.IsAdmin")
	defer span.End()

	// Get user roles
	roles, ok := middleware.GetUserRoles(ctx)
	if !ok {
		err := errors.WithOp(errors.ErrUnauthorized, "service.IsAdmin")
		err = errors.WithMessage(err, "user roles not found in context")
		s.logger.Debug("User roles not found in context")
		return false, err
	}

	// Check if user has admin role
	for _, role := range roles {
		if role == s.config.AdminRoleName {
			return true, nil
		}
	}

	return false, nil
}

// HasRole checks if the user has a specific role.
func (s *Service) HasRole(ctx context.Context, role string) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "service.HasRole")
	defer span.End()

	span.SetAttributes(attribute.String("role", role))

	// Get user roles
	roles, ok := middleware.GetUserRoles(ctx)
	if !ok {
		err := errors.WithOp(errors.ErrUnauthorized, "service.HasRole")
		err = errors.WithMessage(err, "user roles not found in context")
		s.logger.Debug("User roles not found in context")
		return false, err
	}

	// Check if user has the specified role
	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}

	return false, nil
}

// GetUserID retrieves the user ID from the context.
func (s *Service) GetUserID(ctx context.Context) (string, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		err := errors.WithOp(errors.ErrUnauthorized, "service.GetUserID")
		err = errors.WithMessage(err, "user ID not found in context")
		s.logger.Debug("User ID not found in context")
		return "", err
	}
	return userID, nil
}

// GetUserRoles retrieves the user roles from the context.
func (s *Service) GetUserRoles(ctx context.Context) ([]string, error) {
	roles, ok := middleware.GetUserRoles(ctx)
	if !ok {
		err := errors.WithOp(errors.ErrUnauthorized, "service.GetUserRoles")
		err = errors.WithMessage(err, "user roles not found in context")
		s.logger.Debug("User roles not found in context")
		return nil, err
	}
	return roles, nil
}
