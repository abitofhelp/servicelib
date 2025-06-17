// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides generic configuration interfaces and adapters.
package config

import (
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
)

// ConfigAdapter adapts the existing config.Config to the new Config interface
type ConfigAdapter struct {
	config *config.Config
}

// NewConfigAdapter creates a new ConfigAdapter
func NewConfigAdapter(cfg *config.Config) *ConfigAdapter {
	return &ConfigAdapter{
		config: cfg,
	}
}

// GetApp returns the application configuration
func (a *ConfigAdapter) GetApp() AppConfig {
	return &AppConfigAdapter{
		config: a.config,
	}
}

// GetDatabase returns the database configuration
func (a *ConfigAdapter) GetDatabase() DatabaseConfig {
	return &DatabaseConfigAdapter{
		config: a.config,
	}
}

// AppConfigAdapter adapts the existing config.Config to the new AppConfig interface
type AppConfigAdapter struct {
	config *config.Config
}

// GetVersion returns the application version
func (a *AppConfigAdapter) GetVersion() string {
	return a.config.App.Version
}

// GetName returns the application name
func (a *AppConfigAdapter) GetName() string {
	// The existing config doesn't have a name field, so we return a default value
	return "family-service-graphql"
}

// GetEnvironment returns the application environment
func (a *AppConfigAdapter) GetEnvironment() string {
	// The existing config doesn't have an environment field, so we return a default value
	return "development"
}

// DatabaseConfigAdapter adapts the existing config.Config to the new DatabaseConfig interface
type DatabaseConfigAdapter struct {
	config *config.Config
}

// GetType returns the database type
func (a *DatabaseConfigAdapter) GetType() string {
	return a.config.Database.Type
}

// GetConnectionString returns the database connection string
func (a *DatabaseConfigAdapter) GetConnectionString() string {
	switch a.config.Database.Type {
	case "mongodb":
		return a.config.Database.MongoDB.URI
	case "postgres":
		return a.config.Database.Postgres.DSN
	case "sqlite":
		return a.config.Database.SQLite.URI
	default:
		return ""
	}
}

// GetDatabaseName returns the database name
func (a *DatabaseConfigAdapter) GetDatabaseName() string {
	// This is hardcoded in the current implementation
	return "family_service"
}

// GetCollectionName returns the collection/table name for a given entity type
func (a *DatabaseConfigAdapter) GetCollectionName(entityType string) string {
	// This is hardcoded in the current implementation
	switch entityType {
	case "family":
		return "families"
	default:
		return entityType + "s" // Simple pluralization
	}
}
