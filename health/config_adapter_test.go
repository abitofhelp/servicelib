// Copyright (c) 2025 A Bit of Help, Inc.

package health

import (
	"testing"

	"github.com/abitofhelp/servicelib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppConfig is a mock implementation of config.AppConfig
type MockAppConfig struct {
	mock.Mock
}

// GetVersion returns the application version
func (m *MockAppConfig) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

// GetName returns the application name
func (m *MockAppConfig) GetName() string {
	args := m.Called()
	return args.String(0)
}

// GetEnvironment returns the application environment
func (m *MockAppConfig) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

// MockConfig is a mock implementation of config.Config
type MockConfig struct {
	mock.Mock
}

// GetApp returns the application configuration
func (m *MockConfig) GetApp() config.AppConfig {
	args := m.Called()
	return args.Get(0).(config.AppConfig)
}

// GetDatabase returns the database configuration
func (m *MockConfig) GetDatabase() config.DatabaseConfig {
	args := m.Called()
	return args.Get(0).(config.DatabaseConfig)
}

// TestNewConfigAdapter tests the NewConfigAdapter function
func TestNewConfigAdapter(t *testing.T) {
	// Create a mock app config
	mockAppConfig := new(MockAppConfig)
	mockAppConfig.On("GetVersion").Return("1.0.0")
	mockAppConfig.On("GetName").Return("test-app")
	mockAppConfig.On("GetEnvironment").Return("test")

	// Create a mock config
	mockConfig := new(MockConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)

	// Create a config adapter
	adapter := NewConfigAdapter(mockConfig)

	// Verify that the adapter is not nil
	assert.NotNil(t, adapter)
	assert.Equal(t, mockConfig, adapter.config)

	// Verify that the adapter's GetVersion method returns the expected value
	version := adapter.GetVersion()
	assert.Equal(t, "1.0.0", version)

	// Verify that the adapter's GetName method returns the expected value
	name := adapter.GetName()
	assert.Equal(t, "test-app", name)

	// Verify that the adapter's GetEnvironment method returns the expected value
	env := adapter.GetEnvironment()
	assert.Equal(t, "test", env)

	// Verify that the adapter's GetTimeout method returns the expected value
	timeout := adapter.GetTimeout()
	assert.Equal(t, 5, timeout)

	// Verify that the mock's expectations were met
	mockAppConfig.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}
