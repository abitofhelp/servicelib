//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the health package
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abitofhelp/servicelib/config"
	"github.com/abitofhelp/servicelib/health"
	"go.uber.org/zap"
)

// MyHealthProvider implements the health.HealthCheckProvider interface
type MyHealthProvider struct{}

// GetRepositoryFactory returns a mock repository factory
func (p *MyHealthProvider) GetRepositoryFactory() any {
	// In a real application, this would return an actual repository factory
	// For this example, we'll just return a non-nil value
	return &struct{}{}
}

// MyAppConfig implements the config.AppConfig interface
type MyAppConfig struct{}

// GetVersion returns the application version
func (a *MyAppConfig) GetVersion() string {
	return "1.0.0"
}

// GetName returns the application name
func (a *MyAppConfig) GetName() string {
	return "my-service"
}

// GetEnvironment returns the application environment
func (a *MyAppConfig) GetEnvironment() string {
	return "development"
}

// MyDatabaseConfig implements the config.DatabaseConfig interface
type MyDatabaseConfig struct{}

// GetType returns the database type
func (d *MyDatabaseConfig) GetType() string {
	return "postgres"
}

// GetConnectionString returns the database connection string
func (d *MyDatabaseConfig) GetConnectionString() string {
	return "postgres://user:password@localhost:5432/mydb?sslmode=disable"
}

// GetDatabaseName returns the database name
func (d *MyDatabaseConfig) GetDatabaseName() string {
	return "mydb"
}

// GetCollectionName returns the collection/table name for a given entity type
func (d *MyDatabaseConfig) GetCollectionName(entityType string) string {
	return entityType + "s" // Simple pluralization for this example
}

// MyConfig implements the config.Config interface
type MyConfig struct{}

// GetApp returns the application configuration
func (c *MyConfig) GetApp() config.AppConfig {
	return &MyAppConfig{}
}

// GetDatabase returns the database configuration
func (c *MyConfig) GetDatabase() config.DatabaseConfig {
	return &MyDatabaseConfig{}
}

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a health check provider
	provider := &MyHealthProvider{}

	// Create a configuration
	cfg := &MyConfig{}

	// Create a health check handler
	healthHandler := health.NewHandler(provider, logger, cfg)

	// Register the health check handler
	http.HandleFunc("/health", healthHandler)

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Health check endpoint: http://localhost:8080/health")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     logger.Fatal("Failed to start server", zap.Error(err))
	// }
}