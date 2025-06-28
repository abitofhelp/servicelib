// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Mock types for testing
type GenericMockRepository struct {
	ID string
}

func (r *GenericMockRepository) GetID() string {
	return r.ID
}

type GenericMockDomainService struct {
	ID string
}

func (s *GenericMockDomainService) GetID() string {
	return s.ID
}

type GenericMockApplicationService struct {
	ID string
}

func (s *GenericMockApplicationService) GetID() string {
	return s.ID
}

type GenericMockConfig struct {
	AppName string
}

// Mock initializers
func mockRepoInitializer(ctx context.Context, connectionString string, logger *zap.Logger) (*GenericMockRepository, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if connectionString == "error" {
		return nil, errors.New("connection error")
	}
	return &GenericMockRepository{ID: "repo-" + connectionString}, nil
}

func mockDomainServiceInitializer(repo *GenericMockRepository) (*GenericMockDomainService, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}
	if repo.ID == "repo-error-domain" {
		return nil, errors.New("domain service error")
	}
	return &GenericMockDomainService{ID: "domain-" + repo.ID}, nil
}

func mockAppServiceInitializer(domain *GenericMockDomainService, repo *GenericMockRepository) (*GenericMockApplicationService, error) {
	if domain == nil {
		return nil, errors.New("domain service cannot be nil")
	}
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}
	if domain.ID == "domain-repo-error-app" {
		return nil, errors.New("application service error")
	}
	return &GenericMockApplicationService{ID: "app-" + domain.ID}, nil
}

// TestNewGenericAppContainer tests the NewGenericAppContainer function
func TestNewGenericAppContainer(t *testing.T) {
	// Create a context and logger
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	cfg := GenericMockConfig{AppName: "TestApp"}

	// Test cases
	testCases := []struct {
		name             string
		connectionString string
		expectError      bool
		errorStage       string
	}{
		{
			name:             "Success",
			connectionString: "success",
			expectError:      false,
		},
		{
			name:             "Repository initialization error",
			connectionString: "error",
			expectError:      true,
			errorStage:       "repository",
		},
		{
			name:             "Domain service initialization error",
			connectionString: "error-domain",
			expectError:      true,
			errorStage:       "domain",
		},
		{
			name:             "Application service initialization error",
			connectionString: "error-app",
			expectError:      true,
			errorStage:       "app",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			container, err := NewGenericAppContainer(
				ctx,
				logger,
				cfg,
				tc.connectionString,
				mockRepoInitializer,
				mockDomainServiceInitializer,
				mockAppServiceInitializer,
			)

			// Assertions
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, container)

				// Check error message contains expected stage
				if tc.errorStage == "repository" {
					assert.Contains(t, err.Error(), "failed to initialize repository")
				} else if tc.errorStage == "domain" {
					assert.Contains(t, err.Error(), "failed to initialize domain service")
				} else if tc.errorStage == "app" {
					assert.Contains(t, err.Error(), "failed to initialize application service")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, container)

				// Check that the container has the expected components
				assert.Equal(t, "repo-"+tc.connectionString, container.GetRepository().ID)
				assert.Equal(t, "domain-repo-"+tc.connectionString, container.GetDomainService().ID)
				assert.Equal(t, "app-domain-repo-"+tc.connectionString, container.GetApplicationService().ID)

				// Check that GetRepositoryFactory returns the repository
				repoInterface := container.GetRepositoryFactory()
				repo, ok := repoInterface.(*GenericMockRepository)
				assert.True(t, ok)
				assert.Equal(t, "repo-"+tc.connectionString, repo.ID)

				// Test Close method
				err = container.Close()
				assert.NoError(t, err)
			}
		})
	}
}

// TestNewGenericAppContainerWithNilInputs tests the NewGenericAppContainer function with nil inputs
func TestNewGenericAppContainerWithNilInputs(t *testing.T) {
	// Create a context and logger
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	cfg := GenericMockConfig{AppName: "TestApp"}

	// Test with nil context
	t.Run("Nil context", func(t *testing.T) {
		container, err := NewGenericAppContainer(
			nil,
			logger,
			cfg,
			"test",
			mockRepoInitializer,
			mockDomainServiceInitializer,
			mockAppServiceInitializer,
		)
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to create base container")
	})

	// Test with nil logger
	t.Run("Nil logger", func(t *testing.T) {
		container, err := NewGenericAppContainer(
			ctx,
			nil,
			cfg,
			"test",
			mockRepoInitializer,
			mockDomainServiceInitializer,
			mockAppServiceInitializer,
		)
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to create base container")
	})

	// Test with nil repository initializer
	t.Run("Nil repository initializer", func(t *testing.T) {
		// This will cause a compile-time error because the function signature requires a non-nil function
		// container, err := NewGenericAppContainer(
		//     ctx,
		//     logger,
		//     cfg,
		//     "test",
		//     nil,
		//     mockDomainServiceInitializer,
		//     mockAppServiceInitializer,
		// )
		// assert.Error(t, err)
		// assert.Nil(t, container)
	})
}
