// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRepository is a mock implementation of a repository
type MockRepository struct {
	ID string
}

// GetID returns the repository ID
func (r *MockRepository) GetID() string {
	return r.ID
}

// MockDatabaseConfig is a mock implementation of the DatabaseConfig interface
type MockDatabaseConfig struct {
	mock.Mock
}

// GetType returns the database type
func (m *MockDatabaseConfig) GetType() string {
	args := m.Called()
	return args.String(0)
}

// GetConnectionString returns the database connection string
func (m *MockDatabaseConfig) GetConnectionString() string {
	args := m.Called()
	return args.String(0)
}

// GetDatabaseName returns the database name
func (m *MockDatabaseConfig) GetDatabaseName() string {
	args := m.Called()
	return args.String(0)
}

// GetCollectionName returns the collection name for the given entity type
func (m *MockDatabaseConfig) GetCollectionName(entityType string) string {
	args := m.Called(entityType)
	return args.String(0)
}

// MockAppConfig is a mock implementation of the AppConfig interface
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

// MockConfig is a mock implementation of the Config interface
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

// TestNewRepositoryContainer tests the NewRepositoryContainer function
func TestNewRepositoryContainer(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Create a mock app config
	mockAppConfig := new(MockAppConfig)
	mockAppConfig.On("GetVersion").Return("1.0.0")
	mockAppConfig.On("GetName").Return("test-app")
	mockAppConfig.On("GetEnvironment").Return("test")

	// Create a mock database config
	mockDBConfig := new(MockDatabaseConfig)
	mockDBConfig.On("GetConnectionString").Return("mongodb://localhost:27017")
	mockDBConfig.On("GetDatabaseName").Return("testdb")
	mockDBConfig.On("GetCollectionName", "user").Return("users")

	// Create a mock config
	mockConfig := new(MockConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)
	mockConfig.On("GetDatabase").Return(mockDBConfig)

	// Create a mock repository initializer function
	initRepo := func(
		ctx context.Context,
		connectionString string,
		databaseName string,
		collectionName string,
		logger *zap.Logger,
	) (*MockRepository, error) {
		if connectionString == "" {
			return nil, errors.New("connection string cannot be empty")
		}
		return &MockRepository{ID: "test-repo"}, nil
	}

	// Test case 1: Valid inputs
	t.Run("Valid inputs", func(t *testing.T) {
		// Create a container
		container, err := NewRepositoryContainer(
			context.Background(),
			logger,
			mockConfig,
			"user",
			initRepo,
		)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, container)
		assert.Equal(t, mockConfig, container.GetConfig())
		assert.NotNil(t, container.GetContext())
		assert.NotNil(t, container.GetLogger())
		assert.NotNil(t, container.GetRepository())
		assert.Equal(t, "test-repo", container.GetRepository().ID)
		assert.Equal(t, container.GetRepository(), container.GetRepositoryFactory())
	})

	// Test case 2: Nil context
	t.Run("Nil context", func(t *testing.T) {
		// Create a container with nil context
		container, err := NewRepositoryContainer(
			nil,
			logger,
			mockConfig,
			"user",
			initRepo,
		)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "context cannot be nil")
	})

	// Test case 3: Nil logger
	t.Run("Nil logger", func(t *testing.T) {
		// Create a container with nil logger
		container, err := NewRepositoryContainer(
			context.Background(),
			nil,
			mockConfig,
			"user",
			initRepo,
		)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "logger cannot be nil")
	})

	// Test case 4: Nil config
	t.Run("Nil config", func(t *testing.T) {
		// Create a container with nil config
		container, err := NewRepositoryContainer(
			context.Background(),
			logger,
			nil,
			"user",
			initRepo,
		)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})

	// Test case 5: Repository initialization error
	t.Run("Repository initialization error", func(t *testing.T) {
		// Create a mock app config
		mockAppConfigEmpty := new(MockAppConfig)
		mockAppConfigEmpty.On("GetVersion").Return("1.0.0")
		mockAppConfigEmpty.On("GetName").Return("test-app")
		mockAppConfigEmpty.On("GetEnvironment").Return("test")

		// Create a mock database config that returns an empty connection string
		mockDBConfigEmpty := new(MockDatabaseConfig)
		mockDBConfigEmpty.On("GetConnectionString").Return("")
		mockDBConfigEmpty.On("GetDatabaseName").Return("testdb")
		mockDBConfigEmpty.On("GetCollectionName", "user").Return("users")

		// Create a mock config
		mockConfigEmpty := new(MockConfig)
		mockConfigEmpty.On("GetApp").Return(mockAppConfigEmpty)
		mockConfigEmpty.On("GetDatabase").Return(mockDBConfigEmpty)

		// Create a container with a config that will cause the repository initialization to fail
		container, err := NewRepositoryContainer(
			context.Background(),
			logger,
			mockConfigEmpty,
			"user",
			initRepo,
		)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to initialize repository")
	})
}

// TestRepositoryContainerClose tests the Close method of RepositoryContainer
func TestRepositoryContainerClose(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Create a mock app config
	mockAppConfig := new(MockAppConfig)
	mockAppConfig.On("GetVersion").Return("1.0.0")
	mockAppConfig.On("GetName").Return("test-app")
	mockAppConfig.On("GetEnvironment").Return("test")

	// Create a mock database config
	mockDBConfig := new(MockDatabaseConfig)
	mockDBConfig.On("GetConnectionString").Return("mongodb://localhost:27017")
	mockDBConfig.On("GetDatabaseName").Return("testdb")
	mockDBConfig.On("GetCollectionName", "user").Return("users")

	// Create a mock config
	mockConfig := new(MockConfig)
	mockConfig.On("GetApp").Return(mockAppConfig)
	mockConfig.On("GetDatabase").Return(mockDBConfig)

	// Create a mock repository initializer function
	initRepo := func(
		ctx context.Context,
		connectionString string,
		databaseName string,
		collectionName string,
		logger *zap.Logger,
	) (*MockRepository, error) {
		return &MockRepository{ID: "test-repo"}, nil
	}

	// Create a container
	container, err := NewRepositoryContainer(
		context.Background(),
		logger,
		mockConfig,
		"user",
		initRepo,
	)
	assert.NoError(t, err)

	// Call the Close method
	err = container.Close()

	// Assertions
	assert.NoError(t, err)
}
