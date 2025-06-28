// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// AppRepositoryInitializer is a function type that initializes a repository
type AppRepositoryInitializer[T any] func(ctx context.Context, connectionString string, logger *zap.Logger) (T, error)

// GenericDomainServiceInitializer is a function type that initializes a domain service
type GenericDomainServiceInitializer[R any, D any] func(repository R) (D, error)

// GenericApplicationServiceInitializer is a function type that initializes an application service
type GenericApplicationServiceInitializer[R any, D any, A any] func(domainService D, repository R) (A, error)

// GenericAppContainer is a generic dependency injection container for any application
type GenericAppContainer[R any, D any, A any, C any] struct {
	*BaseContainer[C]
	repository         R
	domainService      D
	applicationService A
}

// NewGenericAppContainer creates a new generic application container
func NewGenericAppContainer[R any, D any, A any, C any](
	ctx context.Context,
	logger *zap.Logger,
	cfg C,
	connectionString string,
	initRepo AppRepositoryInitializer[R],
	initDomainService GenericDomainServiceInitializer[R, D],
	initAppService GenericApplicationServiceInitializer[R, D, A],
) (*GenericAppContainer[R, D, A, C], error) {
	// Create base container
	baseContainer, err := NewBaseContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

	// Create generic container
	container := &GenericAppContainer[R, D, A, C]{
		BaseContainer: baseContainer,
	}

	// Initialize repository
	container.repository, err = initRepo(ctx, connectionString, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Initialize domain service
	container.domainService, err = initDomainService(container.repository)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize domain service: %w", err)
	}

	// Initialize application service
	container.applicationService, err = initAppService(container.domainService, container.repository)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize application service: %w", err)
	}

	return container, nil
}

// GetRepository returns the repository
func (c *GenericAppContainer[R, D, A, C]) GetRepository() R {
	return c.repository
}

// GetDomainService returns the domain service
func (c *GenericAppContainer[R, D, A, C]) GetDomainService() D {
	return c.domainService
}

// GetApplicationService returns the application service
func (c *GenericAppContainer[R, D, A, C]) GetApplicationService() A {
	return c.applicationService
}

// GetRepositoryFactory returns the repository as an interface{}
func (c *GenericAppContainer[R, D, A, C]) GetRepositoryFactory() interface{} {
	return c.repository
}

// Close closes all resources
func (c *GenericAppContainer[R, D, A, C]) Close() error {
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
