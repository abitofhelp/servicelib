//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of implementing and using the DatabaseConfig interface
package main

import (
	"fmt"
	"strings"

	"github.com/abitofhelp/servicelib/config"
)

// DatabaseSettings is a custom configuration struct that implements the DatabaseConfigProvider interface
type DatabaseSettings struct {
	Type         string
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	Options      map[string]string
	TablePrefix  string
}

// GetDatabaseType implements the DatabaseConfigProvider interface
func (s *DatabaseSettings) GetDatabaseType() string {
	return s.Type
}

// GetDatabaseConnectionString implements the DatabaseConfigProvider interface
func (s *DatabaseSettings) GetDatabaseConnectionString(dbType string) string {
	if dbType != s.Type {
		return ""
	}

	switch dbType {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			s.Username, s.Password, s.Host, s.Port, s.DatabaseName)
	case "mysql":
		return fmt.Sprintf("mysql://%s:%s@%s:%d/%s",
			s.Username, s.Password, s.Host, s.Port, s.DatabaseName)
	case "mongodb":
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			s.Username, s.Password, s.Host, s.Port, s.DatabaseName)
	default:
		return ""
	}
}

// Additional methods specific to DatabaseSettings
func (s *DatabaseSettings) GetOption(key string, defaultValue string) string {
	if value, ok := s.Options[key]; ok {
		return value
	}
	return defaultValue
}

func (s *DatabaseSettings) GetTableName(entityName string) string {
	return s.TablePrefix + strings.ToLower(entityName)
}

func main() {
	// Create database settings
	dbSettings := &DatabaseSettings{
		Type:         "postgres",
		Host:         "localhost",
		Port:         5432,
		Username:     "admin",
		Password:     "password123",
		DatabaseName: "myapp_db",
		Options: map[string]string{
			"sslmode":   "disable",
			"pool_size": "10",
			"timeout":   "30s",
		},
		TablePrefix: "app_",
	}

	// Create a config adapter
	adapter := config.NewGenericConfigAdapter(dbSettings).
		WithDatabaseName(dbSettings.DatabaseName)

	// Get the database configuration through the adapter
	dbConfig := adapter.GetDatabase()

	// Use the standard DatabaseConfig interface methods
	fmt.Println("=== Database Configuration ===")
	fmt.Println("Type:", dbConfig.GetType())
	fmt.Println("Connection String:", dbConfig.GetConnectionString())
	fmt.Println("Database Name:", dbConfig.GetDatabaseName())
	fmt.Println("Users Collection:", dbConfig.GetCollectionName("user"))

	// Use the original settings object for additional functionality
	fmt.Println("\n=== Additional Database Settings ===")
	fmt.Println("SSL Mode:", dbSettings.GetOption("sslmode", "require"))
	fmt.Println("Pool Size:", dbSettings.GetOption("pool_size", "5"))
	fmt.Println("Timeout:", dbSettings.GetOption("timeout", "10s"))
	fmt.Println("Users Table:", dbSettings.GetTableName("User"))
	fmt.Println("Products Table:", dbSettings.GetTableName("Product"))

	// Expected output:
	// === Database Configuration ===
	// Type: postgres
	// Connection String: postgres://admin:password123@localhost:5432/myapp_db
	// Database Name: myapp_db
	// Users Collection: users
	//
	// === Additional Database Settings ===
	// SSL Mode: disable
	// Pool Size: 10
	// Timeout: 30s
	// Users Table: app_user
	// Products Table: app_product
}