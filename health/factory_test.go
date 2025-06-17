// Copyright (c) 2025 A Bit of Help, Inc.

package health

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// TestNewHealthHandler tests the NewHealthHandler function
func TestNewHealthHandler(t *testing.T) {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a mock repository factory
	mockFactory := new(MockRepositoryFactory)
	mockRepo := "mock-repo"
	mockFactory.On("GetRepository").Return(mockRepo)

	// Create a mock app config
	mockAppConfig := new(MockAppConfig)
	mockAppConfig.On("GetVersion").Return("1.0.0")
	mockAppConfig.On("GetName").Return("test-app")
	mockAppConfig.On("GetEnvironment").Return("test")

	// Create a mock config
	mockConfig := new(MockConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)

	// Call the function
	handler := NewHealthHandler(mockFactory, mockConfig, logger)

	// Verify that the handler is not nil
	assert.NotNil(t, handler)
	assert.IsType(t, (http.HandlerFunc)(nil), handler)

	// Verify that the mock's expectations were met
	mockFactory.AssertExpectations(t)
	mockAppConfig.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}