// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"

	"github.com/abitofhelp/servicelib/config"
	"go.uber.org/zap"
)

// ServiceContainer is a generic dependency injection container for services
type ServiceContainer[R Repository, D DomainService, A ApplicationService, C config.Config] struct {
	*BaseContainer[C]
	repository         R
	domainService      D
	applicationService A
}

// DomainServiceInitializerFunc is a function type that initializes a domain service
type DomainServiceInitializerFunc[R Repository, D DomainService] func(repository R) (D, error)

// ApplicationServiceInitializerFunc is a function type that initializes an application service
type ApplicationServiceInitializerFunc[R Repository, D DomainService, A ApplicationService] func(
	domainService D,
	repository R,
) (A, error)

// NewServiceContainer creates a new service container
func NewServiceContainer[R Repository, D DomainService, A ApplicationService, C config.Config](
	ctx context.Context,
	logger *zap.Logger,
	cfg C,
	repository R,
	initDomainService DomainServiceInitializerFunc[R, D],
	initAppService ApplicationServiceInitializerFunc[R, D, A],
) (*ServiceContainer[R, D, A, C], error) {
	// Create base container
	baseContainer, err := NewBaseContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

	// Create service container
	container := &ServiceContainer[R, D, A, C]{
		BaseContainer: baseContainer,
		repository:    repository,
	}

	// Initialize domain service
	container.domainService, err = initDomainService(repository)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize domain service: %w", err)
	}

	// Initialize application service
	container.applicationService, err = initAppService(container.domainService, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize application service: %w", err)
	}

	return container, nil
}

// GetRepository returns the repository
func (c *ServiceContainer[R, D, A, C]) GetRepository() R {
	return c.repository
}

// GetDomainService returns the domain service
func (c *ServiceContainer[R, D, A, C]) GetDomainService() D {
	return c.domainService
}

// GetApplicationService returns the application service
func (c *ServiceContainer[R, D, A, C]) GetApplicationService() A {
	return c.applicationService
}

// GetRepositoryFactory returns the repository as an interface{}
func (c *ServiceContainer[R, D, A, C]) GetRepositoryFactory() interface{} {
	return c.repository
}

// Close closes all resources
func (c *ServiceContainer[R, D, A, C]) Close() error {
	// Close base container
	return c.BaseContainer.Close()
}
