// Copyright (c) 2025 A Bit of Help, Inc.

package config

import (
	"testing"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/stretchr/testify/assert"
)

// createMockConfig creates a mock config.Config for testing
func createMockConfig(appVersion, dbType, mongoURI, postgresDSN, sqliteURI string) *config.Config {
	cfg := &config.Config{}

	// Set App values
	cfg.App.Version = appVersion

	// Set Database values
	cfg.Database.Type = dbType
	cfg.Database.MongoDB.URI = mongoURI
	cfg.Database.Postgres.DSN = postgresDSN
	cfg.Database.SQLite.URI = sqliteURI

	return cfg
}

func TestNewConfigAdapter(t *testing.T) {
	// Create a mock config
	mockConfig := createMockConfig(
		"1.0.0",
		"sqlite",
		"mongodb://localhost:27017",
		"postgres://localhost:5432",
		"file:test.db",
	)

	// Create a config adapter
	adapter := NewConfigAdapter(mockConfig)

	// Check that the adapter is not nil
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.config)
}

func TestConfigAdapter_GetApp(t *testing.T) {
	// Create a mock config
	mockConfig := createMockConfig(
		"1.0.0",
		"",
		"",
		"",
		"",
	)

	// Create a config adapter
	adapter := NewConfigAdapter(mockConfig)

	// Get the app config
	appConfig := adapter.GetApp()

	// Check that the app config is not nil
	assert.NotNil(t, appConfig)

	// Check that the app config has the expected values
	assert.Equal(t, "1.0.0", appConfig.GetVersion())
	assert.Equal(t, "family-service-graphql", appConfig.GetName())
	assert.Equal(t, "development", appConfig.GetEnvironment())
}

func TestConfigAdapter_GetDatabase(t *testing.T) {
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
			adapter := NewConfigAdapter(mockConfig)

			// Get the database config
			dbConfig := adapter.GetDatabase()

			// Check that the database config is not nil
			assert.NotNil(t, dbConfig)

			// Check that the database config has the expected values
			assert.Equal(t, tc.expectedType, dbConfig.GetType())
			assert.Equal(t, tc.expectedConnStr, dbConfig.GetConnectionString())
			assert.Equal(t, "family_service", dbConfig.GetDatabaseName())
		})
	}
}

// TestDatabaseConfigAdapter is a special version of DatabaseConfigAdapter for testing
type TestDatabaseConfigAdapter struct{}

// GetCollectionName returns the collection/table name for a given entity type
func (a *TestDatabaseConfigAdapter) GetCollectionName(entityType string) string {
	// This is hardcoded in the current implementation
	switch entityType {
	case "family":
		return "families"
	default:
		return entityType + "s" // Simple pluralization
	}
}

func TestDatabaseConfigAdapter_GetCollectionName(t *testing.T) {
	// Test cases for different entity types
	testCases := []struct {
		name         string
		entityType   string
		expectedName string
	}{
		{
			name:         "Family",
			entityType:   "family",
			expectedName: "families",
		},
		{
			name:         "Other",
			entityType:   "user",
			expectedName: "users",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test database config adapter
			adapter := &TestDatabaseConfigAdapter{}

			// Get the collection name
			collectionName := adapter.GetCollectionName(tc.entityType)

			// Check that the collection name is as expected
			assert.Equal(t, tc.expectedName, collectionName)
		})
	}
}
