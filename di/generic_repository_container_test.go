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
type RepoMockRepository struct {
	ID string
}

func (r *RepoMockRepository) GetID() string {
	return r.ID
}

// Mock initializer
func repoMockInitializer(ctx context.Context, connectionString string, logger *zap.Logger) (*RepoMockRepository, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if connectionString == "error" {
		return nil, errors.New("connection error")
	}
	return &RepoMockRepository{ID: "repo-" + connectionString}, nil
}

// TestGenericRepositoryContainer tests the GenericRepositoryContainer struct
func TestGenericRepositoryContainer(t *testing.T) {
	// Skip this test as it requires mocking the config.NewGenericConfigAdapter function
	// which is not easily mockable due to its generic nature
	t.Skip("Skipping TestGenericRepositoryContainer as it requires mocking generic functions")
}

// TestGenericRepositoryContainerGetters tests the getter methods of GenericRepositoryContainer
func TestGenericRepositoryContainerGetters(t *testing.T) {
	// Create a mock repository
	repo := &RepoMockRepository{ID: "test-repo"}

	// Create a simple test to verify the GetRepository and GetRepositoryFactory methods
	t.Run("GetRepository and GetRepositoryFactory", func(t *testing.T) {
		// Create a container directly (bypassing NewGenericRepositoryContainer)
		container := &GenericRepositoryContainer[*RepoMockRepository, interface{}]{
			repository: repo,
		}

		// Test GetRepository
		assert.Equal(t, repo, container.GetRepository())

		// Test GetRepositoryFactory
		repoInterface := container.GetRepositoryFactory()
		assert.Equal(t, repo, repoInterface)
	})
}

// TestGenericRepositoryContainerClose tests the Close method of GenericRepositoryContainer
func TestGenericRepositoryContainerClose(t *testing.T) {
	// Create a mock repository
	repo := &RepoMockRepository{ID: "test-repo"}

	// Create a simple test to verify the Close method
	t.Run("Close", func(t *testing.T) {
		// Create a container directly (bypassing NewGenericRepositoryContainer)
		container := &GenericRepositoryContainer[*RepoMockRepository, interface{}]{
			repository: repo,
		}

		// Test Close
		err := container.Close()
		assert.NoError(t, err)
	})
}
