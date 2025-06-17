// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"

	"github.com/abitofhelp/servicelib/config"
	"go.uber.org/zap"
)

// RepositoryInitializerFunc is a function type that initializes a repository
type RepositoryInitializerFunc[T any] func(
	ctx context.Context, 
	connectionString string, 
	databaseName string, 
	collectionName string, 
	logger *zap.Logger,
) (T, error)

// RepositoryContainer is a generic dependency injection container for any repository type
type RepositoryContainer[T any] struct {
	ctx        context.Context
	logger     *zap.Logger
	config     config.Config
	repository T
}

// NewRepositoryContainer creates a new repository container
func NewRepositoryContainer[T any](
	ctx context.Context,
	logger *zap.Logger,
	cfg config.Config,
	entityType string,
	initRepo RepositoryInitializerFunc[T],
) (*RepositoryContainer[T], error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Create container
	container := &RepositoryContainer[T]{
		ctx:    ctx,
		logger: logger,
		config: cfg,
	}

	// Get database configuration
	dbConfig := cfg.GetDatabase()
	connectionString := dbConfig.GetConnectionString()
	databaseName := dbConfig.GetDatabaseName()
	collectionName := dbConfig.GetCollectionName(entityType)

	// Initialize repository
	var err error
	container.repository, err = initRepo(ctx, connectionString, databaseName, collectionName, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return container, nil
}

// GetContext returns the context
func (c *RepositoryContainer[T]) GetContext() context.Context {
	return c.ctx
}

// GetLogger returns the logger
func (c *RepositoryContainer[T]) GetLogger() *zap.Logger {
	return c.logger
}

// GetConfig returns the configuration
func (c *RepositoryContainer[T]) GetConfig() config.Config {
	return c.config
}

// GetRepository returns the repository
func (c *RepositoryContainer[T]) GetRepository() T {
	return c.repository
}

// GetRepositoryFactory returns the repository as an interface{}
func (c *RepositoryContainer[T]) GetRepositoryFactory() interface{} {
	return c.repository
}

// Close closes all resources
func (c *RepositoryContainer[T]) Close() error {
	// No resources to close in this implementation
	return nil
}
