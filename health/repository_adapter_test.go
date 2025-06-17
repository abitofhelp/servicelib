// Copyright (c) 2025 A Bit of Help, Inc.

package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoryFactory is a mock implementation of repository.RepositoryFactory
type MockRepositoryFactory struct {
	mock.Mock
}

// GetRepository returns a repository
func (m *MockRepositoryFactory) GetRepository() any {
	args := m.Called()
	return args.Get(0)
}

// TestNewRepositoryAdapter tests the NewRepositoryAdapter function
func TestNewRepositoryAdapter(t *testing.T) {
	// Create a mock repository factory
	mockFactory := new(MockRepositoryFactory)
	mockRepo := "mock-repo"
	mockFactory.On("GetRepository").Return(mockRepo)

	// Create a repository adapter
	adapter := NewRepositoryAdapter(mockFactory)

	// Verify that the adapter is not nil
	assert.NotNil(t, adapter)
	assert.Equal(t, mockFactory, adapter.factory)

	// Verify that the adapter's GetRepositoryFactory method returns the expected value
	repo := adapter.GetRepositoryFactory()
	assert.Equal(t, mockRepo, repo)

	// Verify that the mock's expectations were met
	mockFactory.AssertExpectations(t)
}
