// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides generic configuration interfaces and adapters.
package config

// AppConfigProvider defines the interface for accessing application configuration
type AppConfigProvider interface {
	// GetAppVersion returns the application version
	GetAppVersion() string

	// GetAppName returns the application name (optional)
	GetAppName() string

	// GetAppEnvironment returns the application environment (optional)
	GetAppEnvironment() string
}

// DatabaseConfigProvider defines the interface for accessing database configuration
type DatabaseConfigProvider interface {
	// GetDatabaseType returns the database type
	GetDatabaseType() string

	// GetDatabaseConnectionString returns the database connection string
	GetDatabaseConnectionString(dbType string) string
}

// GenericConfigAdapter is a generic adapter for any config type that provides the necessary methods
type GenericConfigAdapter[T any] struct {
	config         T
	appName        string
	appEnvironment string
	dbName         string
}

// NewGenericConfigAdapter creates a new GenericConfigAdapter with default values
func NewGenericConfigAdapter[T any](cfg T) *GenericConfigAdapter[T] {
	return &GenericConfigAdapter[T]{
		config:         cfg,
		appName:        "application",
		appEnvironment: "development",
		dbName:         "database",
	}
}

// WithAppName sets the application name
func (a *GenericConfigAdapter[T]) WithAppName(name string) *GenericConfigAdapter[T] {
	a.appName = name
	return a
}

// WithAppEnvironment sets the application environment
func (a *GenericConfigAdapter[T]) WithAppEnvironment(env string) *GenericConfigAdapter[T] {
	a.appEnvironment = env
	return a
}

// WithDatabaseName sets the database name
func (a *GenericConfigAdapter[T]) WithDatabaseName(name string) *GenericConfigAdapter[T] {
	a.dbName = name
	return a
}

// GetApp returns the application configuration
func (a *GenericConfigAdapter[T]) GetApp() AppConfig {
	return &GenericAppConfigAdapter[T]{
		config:         a.config,
		appName:        a.appName,
		appEnvironment: a.appEnvironment,
	}
}

// GetDatabase returns the database configuration
func (a *GenericConfigAdapter[T]) GetDatabase() DatabaseConfig {
	return &GenericDatabaseConfigAdapter[T]{
		config: a.config,
		dbName: a.dbName,
	}
}

// GenericAppConfigAdapter is a generic adapter for application configuration
type GenericAppConfigAdapter[T any] struct {
	config         T
	appName        string
	appEnvironment string
}

// GetVersion returns the application version
func (a *GenericAppConfigAdapter[T]) GetVersion() string {
	if provider, ok := any(a.config).(AppConfigProvider); ok {
		return provider.GetAppVersion()
	}
	return "1.0.0" // Default version
}

// GetName returns the application name
func (a *GenericAppConfigAdapter[T]) GetName() string {
	if provider, ok := any(a.config).(AppConfigProvider); ok {
		return provider.GetAppName()
	}
	return a.appName
}

// GetEnvironment returns the application environment
func (a *GenericAppConfigAdapter[T]) GetEnvironment() string {
	if provider, ok := any(a.config).(AppConfigProvider); ok {
		return provider.GetAppEnvironment()
	}
	return a.appEnvironment
}

// GenericDatabaseConfigAdapter is a generic adapter for database configuration
type GenericDatabaseConfigAdapter[T any] struct {
	config T
	dbName string
}

// GetType returns the database type
func (a *GenericDatabaseConfigAdapter[T]) GetType() string {
	if provider, ok := any(a.config).(DatabaseConfigProvider); ok {
		return provider.GetDatabaseType()
	}
	return "unknown"
}

// GetConnectionString returns the database connection string
func (a *GenericDatabaseConfigAdapter[T]) GetConnectionString() string {
	if provider, ok := any(a.config).(DatabaseConfigProvider); ok {
		return provider.GetDatabaseConnectionString(a.GetType())
	}
	return ""
}

// GetDatabaseName returns the database name
func (a *GenericDatabaseConfigAdapter[T]) GetDatabaseName() string {
	return a.dbName
}

// GetCollectionName returns the collection/table name for a given entity type
func (a *GenericDatabaseConfigAdapter[T]) GetCollectionName(entityType string) string {
	// Simple pluralization
	return entityType + "s"
}
