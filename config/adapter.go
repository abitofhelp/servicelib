// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides generic configuration interfaces and adapters.
package config

// AppConfigProvider defines the interface for accessing application configuration.
// This interface should be implemented by configuration types that provide
// application-specific configuration values.
type AppConfigProvider interface {
	// GetAppVersion returns the application version.
	// This is typically a semantic version string (e.g., "1.2.3").
	GetAppVersion() string

	// GetAppName returns the application name.
	// This is optional and defaults to "application" if not provided.
	GetAppName() string

	// GetAppEnvironment returns the application environment.
	// Common values include "development", "staging", "production".
	// This is optional and defaults to "development" if not provided.
	GetAppEnvironment() string
}

// DatabaseConfigProvider defines the interface for accessing database configuration.
// This interface should be implemented by configuration types that provide
// database-specific configuration values.
type DatabaseConfigProvider interface {
	// GetDatabaseType returns the database type.
	// Common values include "mongodb", "postgres", "mysql", "sqlite", etc.
	GetDatabaseType() string

	// GetDatabaseConnectionString returns the database connection string.
	// The connection string format depends on the database type provided.
	// The dbType parameter allows for different connection strings based on the database type.
	GetDatabaseConnectionString(dbType string) string
}

// GenericConfigAdapter is a generic adapter for any config type that provides the necessary methods.
// It implements the Config interface and adapts any type T to provide standardized
// configuration access through the AppConfig and DatabaseConfig interfaces.
type GenericConfigAdapter[T any] struct {
	config         T
	appName        string
	appEnvironment string
	dbName         string
}

// NewGenericConfigAdapter creates a new GenericConfigAdapter with default values.
// It initializes the adapter with the provided configuration object and sets
// default values for application name, environment, and database name.
//
// The default values are:
//   - appName: "application"
//   - appEnvironment: "development"
//   - dbName: "database"
//
// These defaults can be overridden using the With* methods.
func NewGenericConfigAdapter[T any](cfg T) *GenericConfigAdapter[T] {
	return &GenericConfigAdapter[T]{
		config:         cfg,
		appName:        "application",
		appEnvironment: "development",
		dbName:         "database",
	}
}

// WithAppName sets the application name.
// This method overrides the default application name ("application") or
// the name provided by the underlying configuration object.
// It returns the adapter itself to allow for method chaining.
func (a *GenericConfigAdapter[T]) WithAppName(name string) *GenericConfigAdapter[T] {
	a.appName = name
	return a
}

// WithAppEnvironment sets the application environment.
// This method overrides the default environment ("development") or
// the environment provided by the underlying configuration object.
// Common values include "development", "staging", and "production".
// It returns the adapter itself to allow for method chaining.
func (a *GenericConfigAdapter[T]) WithAppEnvironment(env string) *GenericConfigAdapter[T] {
	a.appEnvironment = env
	return a
}

// WithDatabaseName sets the database name.
// This method overrides the default database name ("database").
// The database name is used when generating collection/table names
// and in other database-related operations.
// It returns the adapter itself to allow for method chaining.
func (a *GenericConfigAdapter[T]) WithDatabaseName(name string) *GenericConfigAdapter[T] {
	a.dbName = name
	return a
}

// GetApp returns the application configuration.
// This method implements the Config interface and returns an AppConfig
// that provides access to application-specific configuration values.
// The returned AppConfig is a GenericAppConfigAdapter that wraps the
// underlying configuration object.
func (a *GenericConfigAdapter[T]) GetApp() AppConfig {
	return &GenericAppConfigAdapter[T]{
		config:         a.config,
		appName:        a.appName,
		appEnvironment: a.appEnvironment,
	}
}

// GetDatabase returns the database configuration.
// This method implements the Config interface and returns a DatabaseConfig
// that provides access to database-specific configuration values.
// The returned DatabaseConfig is a GenericDatabaseConfigAdapter that wraps
// the underlying configuration object.
func (a *GenericConfigAdapter[T]) GetDatabase() DatabaseConfig {
	return &GenericDatabaseConfigAdapter[T]{
		config: a.config,
		dbName: a.dbName,
	}
}

// GenericAppConfigAdapter is a generic adapter for application configuration.
// It implements the AppConfig interface and adapts any type T to provide
// standardized access to application configuration values.
type GenericAppConfigAdapter[T any] struct {
	config         T
	appName        string
	appEnvironment string
}

// GetVersion returns the application version.
// This method implements the AppConfig interface.
// If the underlying configuration object implements AppConfigProvider,
// its GetAppVersion method is called. Otherwise, a default version "1.0.0" is returned.
func (a *GenericAppConfigAdapter[T]) GetVersion() string {
	if provider, ok := any(a.config).(AppConfigProvider); ok {
		return provider.GetAppVersion()
	}
	return "1.0.0" // Default version
}

// GetName returns the application name.
// This method implements the AppConfig interface.
// If the underlying configuration object implements AppConfigProvider,
// its GetAppName method is called. Otherwise, the name provided during
// adapter creation is returned (defaults to "application").
func (a *GenericAppConfigAdapter[T]) GetName() string {
	if provider, ok := any(a.config).(AppConfigProvider); ok {
		return provider.GetAppName()
	}
	return a.appName
}

// GetEnvironment returns the application environment.
// This method implements the AppConfig interface.
// If the underlying configuration object implements AppConfigProvider,
// its GetAppEnvironment method is called. Otherwise, the environment provided
// during adapter creation is returned (defaults to "development").
func (a *GenericAppConfigAdapter[T]) GetEnvironment() string {
	if provider, ok := any(a.config).(AppConfigProvider); ok {
		return provider.GetAppEnvironment()
	}
	return a.appEnvironment
}

// GenericDatabaseConfigAdapter is a generic adapter for database configuration.
// It implements the DatabaseConfig interface and adapts any type T to provide
// standardized access to database configuration values.
type GenericDatabaseConfigAdapter[T any] struct {
	config T
	dbName string
}

// GetType returns the database type.
// This method implements the DatabaseConfig interface.
// If the underlying configuration object implements DatabaseConfigProvider,
// its GetDatabaseType method is called. Otherwise, "unknown" is returned.
// Common return values include "mongodb", "postgres", "mysql", "sqlite", etc.
func (a *GenericDatabaseConfigAdapter[T]) GetType() string {
	if provider, ok := any(a.config).(DatabaseConfigProvider); ok {
		return provider.GetDatabaseType()
	}
	return "unknown"
}

// GetConnectionString returns the database connection string.
// This method implements the DatabaseConfig interface.
// If the underlying configuration object implements DatabaseConfigProvider,
// its GetDatabaseConnectionString method is called with the database type.
// Otherwise, an empty string is returned.
func (a *GenericDatabaseConfigAdapter[T]) GetConnectionString() string {
	if provider, ok := any(a.config).(DatabaseConfigProvider); ok {
		return provider.GetDatabaseConnectionString(a.GetType())
	}
	return ""
}

// GetDatabaseName returns the database name.
// This method implements the DatabaseConfig interface.
// It returns the database name provided during adapter creation
// (defaults to "database").
func (a *GenericDatabaseConfigAdapter[T]) GetDatabaseName() string {
	return a.dbName
}

// GetCollectionName returns the collection/table name for a given entity type.
// This method implements the DatabaseConfig interface.
// It performs simple pluralization of the entity type to derive the collection name.
// Special cases:
//   - "family" becomes "families"
//   - All other types have "s" appended (e.g., "user" becomes "users")
//
// For more complex pluralization rules, this method should be overridden.
func (a *GenericDatabaseConfigAdapter[T]) GetCollectionName(entityType string) string {
	// Special case for "family" -> "families"
	if entityType == "family" {
		return "families"
	}
	// Simple pluralization for other cases
	return entityType + "s"
}
