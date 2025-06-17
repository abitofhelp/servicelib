// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Default timeout for database operations
const DefaultTimeout = 30 * time.Second

// Container is a generic dependency injection container for backward compatibility
// It uses the BaseContainer with an interface{} config type
type Container struct {
	*BaseContainer[interface{}]
}

// NewContainer creates a new generic dependency injection container
// This function is kept for backward compatibility
func NewContainer(ctx context.Context, logger *zap.Logger, cfg interface{}) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	baseContainer, err := NewBaseContainer(ctx, logger, cfg)
	if err != nil {
		return nil, err
	}

	return &Container{
		BaseContainer: baseContainer,
	}, nil
}

// GetRepositoryFactory returns the repository factory
// This is a placeholder method that should be overridden by derived containers
func (c *Container) GetRepositoryFactory() interface{} {
	return nil
}
