// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Mock types for testing
type ServiceMockRepository struct {
	ID string
}

func (r *ServiceMockRepository) GetID() string {
	return r.ID
}

type ServiceMockDomainService struct {
	ID string
}

func (s *ServiceMockDomainService) GetID() string {
	return s.ID
}

type ServiceMockApplicationService struct {
	ID string
}

func (s *ServiceMockApplicationService) GetID() string {
	return s.ID
}

type ServiceMockConfig struct {
	AppName string
}

func (c ServiceMockConfig) GetApp() config.AppConfig {
	return nil // Not used in tests
}

func (c ServiceMockConfig) GetDatabase() config.DatabaseConfig {
	return nil // Not used in tests
}

// Mock initializers
func serviceMockDomainServiceInitializer(repo *ServiceMockRepository) (*ServiceMockDomainService, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}
	if repo.ID == "error-domain" {
		return nil, errors.New("domain service error")
	}
	return &ServiceMockDomainService{ID: "domain-" + repo.ID}, nil
}

func serviceMockAppServiceInitializer(domain *ServiceMockDomainService, repo *ServiceMockRepository) (*ServiceMockApplicationService, error) {
	if domain == nil {
		return nil, errors.New("domain service cannot be nil")
	}
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}
	if domain.ID == "domain-error-app" {
		return nil, errors.New("application service error")
	}
	return &ServiceMockApplicationService{ID: "app-" + domain.ID}, nil
}

// TestNewServiceContainer tests the NewServiceContainer function
func TestNewServiceContainer(t *testing.T) {
	// Create a context and logger
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	cfg := ServiceMockConfig{AppName: "TestApp"}

	// Test cases
	testCases := []struct {
		name        string
		repository  *ServiceMockRepository
		expectError bool
		errorStage  string
	}{
		{
			name:        "Success",
			repository:  &ServiceMockRepository{ID: "success"},
			expectError: false,
		},
		{
			name:        "Domain service initialization error",
			repository:  &ServiceMockRepository{ID: "error-domain"},
			expectError: true,
			errorStage:  "domain",
		},
		{
			name:        "Application service initialization error",
			repository:  &ServiceMockRepository{ID: "error-app"},
			expectError: true,
			errorStage:  "app",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			container, err := NewServiceContainer(
				ctx,
				logger,
				cfg,
				tc.repository,
				serviceMockDomainServiceInitializer,
				serviceMockAppServiceInitializer,
			)

			// Assertions
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, container)

				// Check error message contains expected stage
				if tc.errorStage == "domain" {
					assert.Contains(t, err.Error(), "failed to initialize domain service")
				} else if tc.errorStage == "app" {
					assert.Contains(t, err.Error(), "failed to initialize application service")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, container)

				// Check that the container has the expected components
				assert.Equal(t, tc.repository, container.GetRepository())
				assert.Equal(t, "domain-"+tc.repository.ID, container.GetDomainService().ID)
				assert.Equal(t, "app-domain-"+tc.repository.ID, container.GetApplicationService().ID)

				// Check that GetRepositoryFactory returns the repository
				repoInterface := container.GetRepositoryFactory()
				assert.Equal(t, tc.repository, repoInterface)

				// Test Close method
				err = container.Close()
				assert.NoError(t, err)
			}
		})
	}
}

// TestNewServiceContainerWithNilInputs tests the NewServiceContainer function with nil inputs
func TestNewServiceContainerWithNilInputs(t *testing.T) {
	// Create a context and logger
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	cfg := ServiceMockConfig{AppName: "TestApp"}
	repo := &ServiceMockRepository{ID: "test"}

	// Test with nil context
	t.Run("Nil context", func(t *testing.T) {
		container, err := NewServiceContainer(
			nil,
			logger,
			cfg,
			repo,
			serviceMockDomainServiceInitializer,
			serviceMockAppServiceInitializer,
		)
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to create base container")
	})

	// Test with nil logger
	t.Run("Nil logger", func(t *testing.T) {
		container, err := NewServiceContainer(
			ctx,
			nil,
			cfg,
			repo,
			serviceMockDomainServiceInitializer,
			serviceMockAppServiceInitializer,
		)
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to create base container")
	})

	// Test with nil repository
	t.Run("Nil repository", func(t *testing.T) {
		container, err := NewServiceContainer(
			ctx,
			logger,
			cfg,
			nil,
			serviceMockDomainServiceInitializer,
			serviceMockAppServiceInitializer,
		)
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to initialize domain service")
	})
}
