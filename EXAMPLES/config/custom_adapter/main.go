// Copyright (c) 2025 A Bit of Help, Inc.

// Example of creating a custom adapter for the config package
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/abitofhelp/servicelib/config"
)

// EnvConfig is a configuration provider that reads from environment variables
type EnvConfig struct {
	prefix string
}

// NewEnvConfig creates a new EnvConfig with the given prefix
func NewEnvConfig(prefix string) *EnvConfig {
	return &EnvConfig{
		prefix: prefix,
	}
}

// GetAppVersion implements the AppConfigProvider interface
func (c *EnvConfig) GetAppVersion() string {
	return c.getEnv("VERSION", "1.0.0")
}

// GetAppName implements the AppConfigProvider interface
func (c *EnvConfig) GetAppName() string {
	return c.getEnv("NAME", "DefaultApp")
}

// GetAppEnvironment implements the AppConfigProvider interface
func (c *EnvConfig) GetAppEnvironment() string {
	return c.getEnv("ENV", "development")
}

// GetDatabaseType implements the DatabaseConfigProvider interface
func (c *EnvConfig) GetDatabaseType() string {
	return c.getEnv("DB_TYPE", "postgres")
}

// GetDatabaseConnectionString implements the DatabaseConfigProvider interface
func (c *EnvConfig) GetDatabaseConnectionString(dbType string) string {
	host := c.getEnv("DB_HOST", "localhost")
	port := c.getEnv("DB_PORT", "5432")
	user := c.getEnv("DB_USER", "postgres")
	pass := c.getEnv("DB_PASS", "postgres")
	name := c.getEnv("DB_NAME", "postgres")

	switch dbType {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name)
	case "mysql":
		return fmt.Sprintf("mysql://%s:%s@%s:%s/%s", user, pass, host, port, name)
	default:
		return ""
	}
}

// Helper method to get environment variables with a prefix
func (c *EnvConfig) getEnv(key, defaultValue string) string {
	fullKey := c.prefix + "_" + key
	if value, exists := os.LookupEnv(fullKey); exists {
		return value
	}
	return defaultValue
}

// Additional methods specific to EnvConfig
func (c *EnvConfig) GetInt(key string, defaultValue int) int {
	strValue := c.getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func (c *EnvConfig) GetBool(key string, defaultValue bool) bool {
	strValue := strings.ToLower(c.getEnv(key, ""))
	if strValue == "" {
		return defaultValue
	}

	return strValue == "true" || strValue == "yes" || strValue == "1"
}

func main() {
	// Set some environment variables for testing
	// In a real application, these would be set outside the program
	os.Setenv("APP_VERSION", "3.0.0")
	os.Setenv("APP_NAME", "EnvApp")
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_DB_TYPE", "mysql")
	os.Setenv("APP_DB_HOST", "db.example.com")
	os.Setenv("APP_DB_PORT", "3306")
	os.Setenv("APP_DB_USER", "admin")
	os.Setenv("APP_DB_PASS", "secret")
	os.Setenv("APP_DB_NAME", "myapp")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_MAX_CONNECTIONS", "50")

	// Create an environment-based configuration
	envConfig := NewEnvConfig("APP")

	// Create a config adapter
	adapter := config.NewGenericConfigAdapter(envConfig).
		WithDatabaseName(envConfig.getEnv("DB_NAME", "default"))

	// Get the app configuration
	appConfig := adapter.GetApp()
	fmt.Println("=== Application Configuration (from Environment) ===")
	fmt.Println("Version:", appConfig.GetVersion())
	fmt.Println("Name:", appConfig.GetName())
	fmt.Println("Environment:", appConfig.GetEnvironment())

	// Get the database configuration
	dbConfig := adapter.GetDatabase()
	fmt.Println("\n=== Database Configuration (from Environment) ===")
	fmt.Println("Type:", dbConfig.GetType())
	fmt.Println("Connection String:", dbConfig.GetConnectionString())
	fmt.Println("Database Name:", dbConfig.GetDatabaseName())

	// Use additional methods from the custom config
	fmt.Println("\n=== Additional Settings (from Environment) ===")
	fmt.Println("Debug Mode:", envConfig.GetBool("DEBUG", false))
	fmt.Println("Max Connections:", envConfig.GetInt("MAX_CONNECTIONS", 10))

	// Clean up environment variables
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_DB_TYPE")
	os.Unsetenv("APP_DB_HOST")
	os.Unsetenv("APP_DB_PORT")
	os.Unsetenv("APP_DB_USER")
	os.Unsetenv("APP_DB_PASS")
	os.Unsetenv("APP_DB_NAME")
	os.Unsetenv("APP_DEBUG")
	os.Unsetenv("APP_MAX_CONNECTIONS")

	// Expected output:
	// === Application Configuration (from Environment) ===
	// Version: 3.0.0
	// Name: EnvApp
	// Environment: production
	//
	// === Database Configuration (from Environment) ===
	// Type: mysql
	// Connection String: mysql://admin:secret@db.example.com:3306/myapp
	// Database Name: myapp
	//
	// === Additional Settings (from Environment) ===
	// Debug Mode: true
	// Max Connections: 50
}