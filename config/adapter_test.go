// Copyright (c) 2025 A Bit of Help, Inc.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockConfig implements the necessary interfaces for testing
type MockConfig struct {
	AppVersion    string
	AppName       string
	AppEnv        string
	DbType        string
	MongoURI      string
	PostgresDSN   string
	SQLiteURI     string
}

// GetAppVersion returns the application version
func (c *MockConfig) GetAppVersion() string {
	return c.AppVersion
}

// GetAppName returns the application name
func (c *MockConfig) GetAppName() string {
	return c.AppName
}

// GetAppEnvironment returns the application environment
func (c *MockConfig) GetAppEnvironment() string {
	return c.AppEnv
}

// GetDatabaseType returns the database type
func (c *MockConfig) GetDatabaseType() string {
	return c.DbType
}

// GetDatabaseConnectionString returns the database connection string
func (c *MockConfig) GetDatabaseConnectionString(dbType string) string {
	switch dbType {
	case "mongodb":
		return c.MongoURI
	case "postgres":
		return c.PostgresDSN
	case "sqlite":
		return c.SQLiteURI
	default:
		return ""
	}
}

// createMockConfig creates a mock config for testing
func createMockConfig(appVersion, dbType, mongoURI, postgresDSN, sqliteURI string) *MockConfig {
	return &MockConfig{
		AppVersion:  appVersion,
		AppName:     "test-app",
		AppEnv:      "test",
		DbType:      dbType,
		MongoURI:    mongoURI,
		PostgresDSN: postgresDSN,
		SQLiteURI:   sqliteURI,
	}
}

func TestNewGenericConfigAdapter(t *testing.T) {
	// Create a mock config
	mockConfig := createMockConfig(
		"1.0.0",
		"sqlite",
		"mongodb://localhost:27017",
		"postgres://localhost:5432",
		"file:test.db",
	)

	// Create a config adapter
	adapter := NewGenericConfigAdapter(mockConfig)

	// Check that the adapter is not nil
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.config)
}

func TestGenericConfigAdapter_GetApp(t *testing.T) {
	// Create a mock config
	mockConfig := createMockConfig(
		"1.0.0",
		"",
		"",
		"",
		"",
	)

	// Create a config adapter
	adapter := NewGenericConfigAdapter(mockConfig)

	// Get the app config
	appConfig := adapter.GetApp()

	// Check that the app config is not nil
	assert.NotNil(t, appConfig)

	// Check that the app config has the expected values
	assert.Equal(t, "1.0.0", appConfig.GetVersion())
	assert.Equal(t, "test-app", appConfig.GetName())
	assert.Equal(t, "test", appConfig.GetEnvironment())
}

func TestGenericConfigAdapter_GetDatabase(t *testing.T) {
	// Test cases for different database types
	testCases := []struct {
		name            string
		dbType          string
		mongoURI        string
		postgresDSN     string
		sqliteURI       string
		expectedType    string
		expectedConnStr string
	}{
		{
			name:            "MongoDB",
			dbType:          "mongodb",
			mongoURI:        "mongodb://localhost:27017",
			expectedType:    "mongodb",
			expectedConnStr: "mongodb://localhost:27017",
		},
		{
			name:            "PostgreSQL",
			dbType:          "postgres",
			postgresDSN:     "postgres://localhost:5432",
			expectedType:    "postgres",
			expectedConnStr: "postgres://localhost:5432",
		},
		{
			name:            "SQLite",
			dbType:          "sqlite",
			sqliteURI:       "file:test.db",
			expectedType:    "sqlite",
			expectedConnStr: "file:test.db",
		},
		{
			name:            "Unknown",
			dbType:          "unknown",
			expectedType:    "unknown",
			expectedConnStr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock config
			mockConfig := createMockConfig(
				"1.0.0",
				tc.dbType,
				tc.mongoURI,
				tc.postgresDSN,
				tc.sqliteURI,
			)

			// Create a config adapter
			adapter := NewGenericConfigAdapter(mockConfig)

			// Get the database config
			dbConfig := adapter.GetDatabase()

			// Check that the database config is not nil
			assert.NotNil(t, dbConfig)

			// Check that the database config has the expected values
			assert.Equal(t, tc.expectedType, dbConfig.GetType())
			assert.Equal(t, tc.expectedConnStr, dbConfig.GetConnectionString())
			assert.Equal(t, "database", dbConfig.GetDatabaseName())
		})
	}
}

func TestGenericDatabaseConfigAdapter_GetCollectionName(t *testing.T) {
	// Test cases for different entity types
	testCases := []struct {
		name         string
		entityType   string
		expectedName string
	}{
		{
			name:         "Family",
			entityType:   "family",
			expectedName: "familys", // Simple pluralization
		},
		{
			name:         "User",
			entityType:   "user",
			expectedName: "users",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock config
			mockConfig := createMockConfig(
				"1.0.0",
				"mongodb",
				"mongodb://localhost:27017",
				"",
				"",
			)

			// Create a config adapter
			adapter := NewGenericConfigAdapter(mockConfig)

			// Get the database config
			dbConfig := adapter.GetDatabase()

			// Get the collection name
			collectionName := dbConfig.GetCollectionName(tc.entityType)

			// Check that the collection name is as expected
			assert.Equal(t, tc.expectedName, collectionName)
		})
	}
}
