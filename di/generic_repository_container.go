// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"

	pkgconfig "github.com/abitofhelp/servicelib/config"
	"go.uber.org/zap"
)

// GenericRepositoryInitializer is a function type that initializes a repository
type GenericRepositoryInitializer[T any] func(ctx context.Context, connectionString string, logger *zap.Logger) (T, error)

// GenericRepositoryContainer is a generic dependency injection container for any repository type
type GenericRepositoryContainer[T any, C any] struct {
	*BaseContainer[C]
	config     pkgconfig.Config
	repository T
}

// NewGenericRepositoryContainer creates a new generic repository container
func NewGenericRepositoryContainer[T any, C any](
	ctx context.Context,
	logger *zap.Logger,
	cfg C,
	entityType string,
	initRepo GenericRepositoryInitializer[T],
) (*GenericRepositoryContainer[T, C], error) {
	// Create base container
	baseContainer, err := NewBaseContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

	// Create a config adapter
	configAdapter := pkgconfig.NewGenericConfigAdapter(cfg)

	// Create generic container
	container := &GenericRepositoryContainer[T, C]{
		BaseContainer: baseContainer,
		config:        configAdapter,
	}

	// Get database configuration
	dbConfig := configAdapter.GetDatabase()
	connectionString := dbConfig.GetConnectionString()

	// Initialize repository
	container.repository, err = initRepo(ctx, connectionString, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return container, nil
}

// GetRepository returns the repository
func (c *GenericRepositoryContainer[T, C]) GetRepository() T {
	return c.repository
}

// GetRepositoryFactory returns the repository as an interface{}
func (c *GenericRepositoryContainer[T, C]) GetRepositoryFactory() interface{} {
	return c.repository
}

// Close closes all resources
func (c *GenericRepositoryContainer[T, C]) Close() error {
	var errs []error

	// Add resource cleanup here as needed
	// For example, close database connections if they implement a Close method

	// Close base container
	if err := c.BaseContainer.Close(); err != nil {
		errs = append(errs, err)
	}

	// Return a combined error if any occurred
	if len(errs) > 0 {
		errMsg := "failed to close one or more resources:"
		for _, err := range errs {
			errMsg += " " + err.Error()
		}
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}
