package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockHealthCheckProvider is a mock implementation of HealthCheckProvider
type MockHealthCheckProvider struct {
	mock.Mock
}

// GetRepositoryFactory is a mock implementation of GetRepositoryFactory
func (m *MockHealthCheckProvider) GetRepositoryFactory() any {
	args := m.Called()
	return args.Get(0)
}

// MockVersionProvider is a mock implementation of VersionProvider
type MockVersionProvider struct {
	mock.Mock
}

// GetVersion is a mock implementation of GetVersion
func (m *MockVersionProvider) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

// TestNewGenericHandler_HealthyService tests the NewGenericHandler function with a healthy service
func TestNewGenericHandler_HealthyService(t *testing.T) {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a mock health check provider
	mockProvider := new(MockHealthCheckProvider)
	mockRepo := "mock-repo"
	mockProvider.On("GetRepositoryFactory").Return(mockRepo)

	// Create a mock version provider
	mockVersionProvider := new(MockVersionProvider)
	mockVersion := "1.0.0"
	mockVersionProvider.On("GetVersion").Return(mockVersion)

	// Create a handler
	handler := NewGenericHandler(mockProvider, mockVersionProvider, logger, 5)

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
	var response GenericHealthStatus
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
	mockVersionProvider.AssertExpectations(t)
}

// TestNewGenericHandler_DegradedService tests the NewGenericHandler function with a degraded service
func TestNewGenericHandler_DegradedService(t *testing.T) {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a mock health check provider
	mockProvider := new(MockHealthCheckProvider)
	mockProvider.On("GetRepositoryFactory").Return(nil)

	// Create a mock version provider
	mockVersionProvider := new(MockVersionProvider)
	mockVersion := "1.0.0"
	mockVersionProvider.On("GetVersion").Return(mockVersion)

	// Create a handler
	handler := NewGenericHandler(mockProvider, mockVersionProvider, logger, 5)

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
	var response GenericHealthStatus
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
	mockVersionProvider.AssertExpectations(t)
}