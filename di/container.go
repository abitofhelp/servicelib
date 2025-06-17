// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Default timeout for database operations
const DefaultTimeout = 30 * time.Second

// Container is a generic dependency injection container
type Container struct {
	ctx           context.Context
	logger        *zap.Logger
	contextLogger *logging.ContextLogger
	validator     *validator.Validate
	config        *config.Config
}

// NewContainer creates a new generic dependency injection container
func NewContainer(ctx context.Context, logger *zap.Logger, cfg *config.Config) (*Container, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	container := &Container{
		ctx:    ctx,
		logger: logger,
		config: cfg,
	}

	// Initialize context logger
	container.contextLogger = logging.NewContextLogger(logger)

	// Initialize validator
	container.validator = validator.New()

	return container, nil
}

// GetContext returns the context
func (c *Container) GetContext() context.Context {
	return c.ctx
}

// GetLogger returns the logger
func (c *Container) GetLogger() *zap.Logger {
	return c.logger
}

// GetContextLogger returns the context logger
func (c *Container) GetContextLogger() *logging.ContextLogger {
	return c.contextLogger
}

// GetValidator returns the validator
func (c *Container) GetValidator() *validator.Validate {
	return c.validator
}

// GetConfig returns the configuration
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// Close closes all resources
func (c *Container) Close() error {
	// Base container has no resources to close
	return nil
}
