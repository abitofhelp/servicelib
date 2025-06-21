//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the config package
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/config"
)

// MyConfig is a simple configuration struct that implements the required interfaces
type MyConfig struct {
	Version     string
	Name        string
	Environment string
	DBType      string
	DBConnStr   map[string]string
}

// GetAppVersion implements the AppConfigProvider interface
func (c *MyConfig) GetAppVersion() string {
	return c.Version
}

// GetAppName implements the AppConfigProvider interface
func (c *MyConfig) GetAppName() string {
	return c.Name
}

// GetAppEnvironment implements the AppConfigProvider interface
func (c *MyConfig) GetAppEnvironment() string {
	return c.Environment
}

// GetDatabaseType implements the DatabaseConfigProvider interface
func (c *MyConfig) GetDatabaseType() string {
	return c.DBType
}

// GetDatabaseConnectionString implements the DatabaseConfigProvider interface
func (c *MyConfig) GetDatabaseConnectionString(dbType string) string {
	if connStr, ok := c.DBConnStr[dbType]; ok {
		return connStr
	}
	return ""
}

func main() {
	// Create a configuration
	myConfig := &MyConfig{
		Version:     "1.0.0",
		Name:        "MyApp",
		Environment: "development",
		DBType:      "postgres",
		DBConnStr: map[string]string{
			"postgres": "postgres://user:password@localhost:5432/mydb",
			"mongodb":  "mongodb://localhost:27017",
		},
	}

	// Create a config adapter
	adapter := config.NewGenericConfigAdapter(myConfig).
		WithAppName("MyApplication").
		WithAppEnvironment("production").
		WithDatabaseName("my_database")

	// Get the app configuration
	appConfig := adapter.GetApp()
	fmt.Println("Application Version:", appConfig.GetVersion())
	fmt.Println("Application Name:", appConfig.GetName())
	fmt.Println("Application Environment:", appConfig.GetEnvironment())

	// Get the database configuration
	dbConfig := adapter.GetDatabase()
	fmt.Println("Database Type:", dbConfig.GetType())
	fmt.Println("Database Connection String:", dbConfig.GetConnectionString())
	fmt.Println("Database Name:", dbConfig.GetDatabaseName())
	fmt.Println("Users Collection:", dbConfig.GetCollectionName("user"))

	// Expected output:
	// Application Version: 1.0.0
	// Application Name: MyApp
	// Application Environment: development
	// Database Type: postgres
	// Database Connection String: postgres://user:password@localhost:5432/mydb
	// Database Name: my_database
	// Users Collection: users
}
