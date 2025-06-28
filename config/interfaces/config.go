// Copyright (c) 2025 A Bit of Help, Inc.

// Package interfaces provides generic configuration interfaces that can be used across different applications.
package interfaces

// AppConfig is a generic interface for application configuration
type AppConfig interface {
	// GetVersion returns the application version
	GetVersion() string

	// GetName returns the application name
	GetName() string

	// GetEnvironment returns the application environment (e.g., "development", "production")
	GetEnvironment() string
}

// DatabaseConfig is a generic interface for database configuration
type DatabaseConfig interface {
	// GetType returns the database type (e.g., "mongodb", "postgres", "sqlite")
	GetType() string

	// GetConnectionString returns the database connection string
	GetConnectionString() string

	// GetDatabaseName returns the database name
	GetDatabaseName() string

	// GetCollectionName returns the collection/table name for a given entity type
	GetCollectionName(entityType string) string
}

// Config is a generic interface for application configuration
type Config interface {
	// GetApp returns the application configuration
	GetApp() AppConfig

	// GetDatabase returns the database configuration
	GetDatabase() DatabaseConfig
}
