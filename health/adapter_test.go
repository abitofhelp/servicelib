// Copyright (c) 2025 A Bit of Help, Inc.

package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoryFactoryProvider is a mock implementation of RepositoryFactoryProvider
type MockRepositoryFactoryProvider struct {
	mock.Mock
}

// GetRepositoryFactory is a mock implementation of GetRepositoryFactory
func (m *MockRepositoryFactoryProvider) GetRepositoryFactory() any {
	args := m.Called()
	return args.Get(0)
}

// TestNewGenericContainerAdapter tests the NewGenericContainerAdapter function
func TestNewGenericContainerAdapter(t *testing.T) {
	// Create a mock repository factory provider
	mockProvider := new(MockRepositoryFactoryProvider)
	mockRepo := "mock-repo"
	mockProvider.On("GetRepositoryFactory").Return(mockRepo)

	// Create a generic container adapter
	adapter := NewGenericContainerAdapter(mockProvider)

	// Verify that the adapter is not nil
	assert.NotNil(t, adapter)
	assert.Equal(t, mockProvider, adapter.Container)

	// Verify that the adapter's GetRepositoryFactory method returns the expected value
	repo := adapter.GetRepositoryFactory()
	assert.Equal(t, mockRepo, repo)

	// Verify that the mock's expectations were met
	mockProvider.AssertExpectations(t)
}
