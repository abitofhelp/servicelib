// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides a flexible and extensible configuration management system for applications.
//
// This package offers a set of interfaces and adapters that standardize how applications
// access configuration data, regardless of the underlying configuration source
// (files, environment variables, remote services, etc.).
//
// The package is designed around the adapter pattern, allowing different configuration
// implementations to be used interchangeably while providing a consistent interface
// for application code.
//
// Key components:
//   - Config: The main interface providing access to application and database configuration
//   - AppConfig: Interface for accessing application-specific configuration
//   - DatabaseConfig: Interface for accessing database-specific configuration
//   - GenericConfigAdapter: A flexible adapter that can wrap any configuration type
//
// The package supports:
//   - Type-safe configuration access
//   - Default values for missing configuration
//   - Fluent interface for adapter configuration
//   - Backward compatibility with older configuration patterns
//
// Example usage:
//
//	// Create a configuration struct
//	type MyConfig struct {
//	    AppVersion string
//	    DBType     string
//	    DBConnStr  string
//	}
//
//	// Implement the necessary provider interfaces
//	func (c *MyConfig) GetAppVersion() string {
//	    return c.AppVersion
//	}
//
//	func (c *MyConfig) GetDatabaseType() string {
//	    return c.DBType
//	}
//
//	func (c *MyConfig) GetDatabaseConnectionString(dbType string) string {
//	    return c.DBConnStr
//	}
//
//	// Create and use the adapter
//	cfg := &MyConfig{
//	    AppVersion: "1.0.0",
//	    DBType:     "postgres",
//	    DBConnStr:  "postgres://user:pass@localhost:5432/mydb",
//	}
//
//	adapter := config.NewGenericConfigAdapter(cfg).
//	    WithAppName("myapp").
//	    WithAppEnvironment("production")
//
//	// Use the adapter through the Config interface
//	var configInterface config.Config = adapter
//	appName := configInterface.GetApp().GetName()
//	dbConnStr := configInterface.GetDatabase().GetConnectionString()
package config
