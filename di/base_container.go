// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// BaseContainer is a generic dependency injection container that can be embedded in other containers
type BaseContainer[C any] struct {
	ctx           context.Context
	logger        *zap.Logger
	contextLogger *logging.ContextLogger
	validator     *validator.Validate
	config        C
}

// NewBaseContainer creates a new base dependency injection container
func NewBaseContainer[C any](ctx context.Context, logger *zap.Logger, cfg C) (*BaseContainer[C], error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	container := &BaseContainer[C]{
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
func (c *BaseContainer[C]) GetContext() context.Context {
	return c.ctx
}

// GetLogger returns the logger
func (c *BaseContainer[C]) GetLogger() *zap.Logger {
	return c.logger
}

// GetContextLogger returns the context logger
func (c *BaseContainer[C]) GetContextLogger() *logging.ContextLogger {
	return c.contextLogger
}

// GetValidator returns the validator
func (c *BaseContainer[C]) GetValidator() *validator.Validate {
	return c.validator
}

// GetConfig returns the configuration
func (c *BaseContainer[C]) GetConfig() C {
	return c.config
}

// Close closes all resources
func (c *BaseContainer[C]) Close() error {
	// Base container has no resources to close
	return nil
}
