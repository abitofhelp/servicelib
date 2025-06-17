// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a generic dependency injection container that can be used across different applications.
package di

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/core/application/ports"
	application "github.com/abitofhelp/family-service/core/application/services"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"go.uber.org/zap"
)

// RepositoryInitializer is a function type that initializes a repository
type RepositoryInitializer[T domainports.FamilyRepository] func(ctx context.Context, connectionString string, logger *zap.Logger) (T, error)

// GenericContainer is a generic dependency injection container for any database type
type GenericContainer[T domainports.FamilyRepository] struct {
	*Container
	familyRepo          T
	familyDomainService *domainservices.FamilyDomainService
	familyAppService    ports.FamilyApplicationService
}

// NewGenericContainer creates a new generic dependency injection container
func NewGenericContainer[T domainports.FamilyRepository](
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	initRepo RepositoryInitializer[T],
	connectionString string,
) (*GenericContainer[T], error) {
	// Create base container
	baseContainer, err := NewContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

	// Create generic container
	container := &GenericContainer[T]{
		Container: baseContainer,
	}

	// Initialize repository
	container.familyRepo, err = initRepo(ctx, connectionString, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Initialize domain service
	container.familyDomainService = domainservices.NewFamilyDomainService(container.familyRepo)

	// Initialize application service
	container.familyAppService = application.NewFamilyApplicationService(
		container.familyDomainService,
		container.familyRepo,
	)

	return container, nil
}

// GetFamilyRepository returns the family repository
func (c *GenericContainer[T]) GetFamilyRepository() domainports.FamilyRepository {
	return c.familyRepo
}

// GetFamilyDomainService returns the family domain service
func (c *GenericContainer[T]) GetFamilyDomainService() *domainservices.FamilyDomainService {
	return c.familyDomainService
}

// GetFamilyApplicationService returns the family application service
func (c *GenericContainer[T]) GetFamilyApplicationService() ports.FamilyApplicationService {
	return c.familyAppService
}

// GetRepositoryFactory returns the family repository
func (c *GenericContainer[T]) GetRepositoryFactory() interface{} {
	return c.familyRepo
}

// For backward compatibility with the test file
func (c *GenericContainer[T]) GetFamilyService() interface{} {
	return c.familyAppService
}

// For backward compatibility with the test file
func (c *GenericContainer[T]) GetAuthorizationService() interface{} {
	return nil
}

// Close closes all resources
func (c *GenericContainer[T]) Close() error {
	var errs []error

	// Add resource cleanup here as needed
	// For example, close database connections if they implement a Close method

	// Close base container
	if err := c.Container.Close(); err != nil {
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
