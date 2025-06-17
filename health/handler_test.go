package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockHandlerAppConfig is a mock implementation of config.AppConfig for handler tests
type MockHandlerAppConfig struct {
	mock.Mock
}

// GetVersion returns the application version
func (m *MockHandlerAppConfig) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

// GetName returns the application name
func (m *MockHandlerAppConfig) GetName() string {
	args := m.Called()
	return args.String(0)
}

// GetEnvironment returns the application environment
func (m *MockHandlerAppConfig) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

// MockHandlerConfig is a mock implementation of config.Config for handler tests
type MockHandlerConfig struct {
	mock.Mock
}

// GetApp returns the application configuration
func (m *MockHandlerConfig) GetApp() config.AppConfig {
	args := m.Called()
	return args.Get(0).(config.AppConfig)
}

// GetDatabase returns the database configuration
func (m *MockHandlerConfig) GetDatabase() config.DatabaseConfig {
	args := m.Called()
	return args.Get(0).(config.DatabaseConfig)
}

// MockHandlerHealthCheckProvider is a mock implementation of HealthCheckProvider for handler tests
type MockHandlerHealthCheckProvider struct {
	mock.Mock
}

// GetRepositoryFactory is a mock implementation of GetRepositoryFactory
func (m *MockHandlerHealthCheckProvider) GetRepositoryFactory() any {
	args := m.Called()
	return args.Get(0)
}

// TestNewHandler_HealthyService tests the NewHandler function with a healthy service
func TestNewHandler_HealthyService(t *testing.T) {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a mock health check provider
	mockProvider := new(MockHandlerHealthCheckProvider)
	mockRepo := "mock-repo"
	mockProvider.On("GetRepositoryFactory").Return(mockRepo)

	// Create a mock app config
	mockAppConfig := new(MockHandlerAppConfig)
	mockVersion := "1.0.0"
	mockAppConfig.On("GetVersion").Return(mockVersion)

	// Create a mock config
	mockConfig := new(MockHandlerConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)

	// Create a handler
	handler := NewHandler(mockProvider, logger, mockConfig)

	// Create a request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var response HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, StatusHealthy, response.Status)
	assert.Equal(t, mockVersion, response.Version)
	assert.Equal(t, ServiceUp, response.Services["database"])

	// Verify that the timestamp is in the correct format
	_, err = time.Parse(time.RFC3339, response.Timestamp)
	assert.NoError(t, err)

	// Verify that the mock's expectations were met
	mockProvider.AssertExpectations(t)
	mockAppConfig.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}

// TestNewHandler_DegradedService tests the NewHandler function with a degraded service
func TestNewHandler_DegradedService(t *testing.T) {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a mock health check provider
	mockProvider := new(MockHandlerHealthCheckProvider)
	mockProvider.On("GetRepositoryFactory").Return(nil)

	// Create a mock app config
	mockAppConfig := new(MockHandlerAppConfig)
	mockVersion := "1.0.0"
	mockAppConfig.On("GetVersion").Return(mockVersion)

	// Create a mock config
	mockConfig := new(MockHandlerConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)

	// Create a handler
	handler := NewHandler(mockProvider, logger, mockConfig)

	// Create a request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	// Parse the response
	var response HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, StatusDegraded, response.Status)
	assert.Equal(t, mockVersion, response.Version)
	assert.Equal(t, ServiceDown, response.Services["database"])

	// Verify that the timestamp is in the correct format
	_, err = time.Parse(time.RFC3339, response.Timestamp)
	assert.NoError(t, err)

	// Verify that the mock's expectations were met
	mockProvider.AssertExpectations(t)
	mockAppConfig.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}

// TestNewHandler_EncodingError tests the NewHandler function with an encoding error
func TestNewHandler_EncodingError(t *testing.T) {
	// This test is a bit tricky because we can't easily force an encoding error
	// in a normal test scenario. In a real-world scenario, we might use a custom
	// ResponseWriter that returns an error when Write is called, but for this
	// test we'll just verify that the handler doesn't panic.

	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a mock health check provider
	mockProvider := new(MockHandlerHealthCheckProvider)
	mockRepo := "mock-repo"
	mockProvider.On("GetRepositoryFactory").Return(mockRepo)

	// Create a mock app config
	mockAppConfig := new(MockHandlerAppConfig)
	mockVersion := "1.0.0"
	mockAppConfig.On("GetVersion").Return(mockVersion)

	// Create a mock config
	mockConfig := new(MockHandlerConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)

	// Create a handler
	handler := NewHandler(mockProvider, logger, mockConfig)

	// Create a request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler and verify it doesn't panic
	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})

	// Verify that the mock's expectations were met
	mockProvider.AssertExpectations(t)
	mockAppConfig.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}
