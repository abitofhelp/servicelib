//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Package main demonstrates how to integrate environment variables with a configuration structure.
//
// This example shows a more advanced usage of the env package by:
// - Creating a configuration structure to hold application settings
// - Loading environment variables into the configuration structure
// - Validating the configuration values
// - Handling complex types like slices from environment variables
//
// This pattern is useful for applications that need to manage multiple
// configuration values in a structured way, while still allowing for
// easy configuration through environment variables.
package main

import (
	"fmt"
	"strings"

	"github.com/abitofhelp/servicelib/env"
)

// AppConfig holds the application configuration settings loaded from environment variables.
//
// This struct centralizes all configuration settings in one place, making it easier
// to manage and access application configuration. Each field corresponds to a
// specific environment variable with a default value.
//
// Using a struct for configuration provides several benefits:
// - Type safety for configuration values
// - Centralized documentation of available settings
// - Easy validation of the entire configuration
// - Simple passing of configuration to different parts of the application
type AppConfig struct {
	ServerPort  string   // Port the server will listen on (from SERVER_PORT)
	DatabaseURL string   // Connection string for the database (from DATABASE_URL)
	LogLevel    string   // Logging level: debug, info, warn, error (from LOG_LEVEL)
	APIKey      string   // API key for external service authentication (from API_KEY)
	Environment string   // Application environment: development, testing, staging, production (from APP_ENV)
	Features    []string // Enabled features as a slice of strings (from comma-separated FEATURES)
}

// LoadConfig loads the application configuration from environment variables.
//
// This function retrieves values from environment variables using the env package
// and populates an AppConfig struct with those values. If an environment variable
// is not set, it uses the specified default value.
//
// For the FEATURES environment variable, it expects a comma-separated list of
// feature names, which it splits into a slice of strings.
//
// Returns:
//   - An AppConfig struct populated with values from environment variables
//
// Example usage:
//
//	config := LoadConfig()
//	fmt.Printf("Server will run on port: %s\n", config.ServerPort)
func LoadConfig() AppConfig {
	// Load features from a comma-separated environment variable
	featuresStr := env.GetEnv("FEATURES", "basic,standard")
	features := strings.Split(featuresStr, ",")

	return AppConfig{
		ServerPort:  env.GetEnv("SERVER_PORT", "8080"),
		DatabaseURL: env.GetEnv("DATABASE_URL", "postgres://localhost:5432/mydb"),
		LogLevel:    env.GetEnv("LOG_LEVEL", "info"),
		APIKey:      env.GetEnv("API_KEY", ""),
		Environment: env.GetEnv("APP_ENV", "development"),
		Features:    features,
	}
}

// Validate checks the configuration values for correctness and returns any validation errors.
//
// This method performs several validation checks on the AppConfig:
// - Ensures required fields like APIKey are not empty
// - Validates that LogLevel is one of the allowed values
// - Confirms that Environment is a valid environment name
//
// Having a dedicated validation method allows for complex validation logic
// that can't be expressed through simple type constraints.
//
// Returns:
//   - A slice of strings containing error messages, or an empty slice if validation passes
//
// Example usage:
//
//	config := LoadConfig()
//	errors := config.Validate()
//	if len(errors) > 0 {
//	    fmt.Println("Configuration errors:", errors)
//	    os.Exit(1)
//	}
func (c *AppConfig) Validate() []string {
	var errors []string

	// Check required values
	if c.APIKey == "" {
		errors = append(errors, "API_KEY environment variable is required")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(c.LogLevel)] {
		errors = append(errors, fmt.Sprintf("Invalid log level: %s. Must be one of: debug, info, warn, error", c.LogLevel))
	}

	// Validate environment
	validEnvs := map[string]bool{
		"development": true,
		"testing":     true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvs[strings.ToLower(c.Environment)] {
		errors = append(errors, fmt.Sprintf("Invalid environment: %s. Must be one of: development, testing, staging, production", c.Environment))
	}

	return errors
}

func main() {
	fmt.Println("Environment Variables Configuration Integration Example")
	fmt.Println("=====================================================")

	// Load configuration from environment variables
	config := LoadConfig()

	// Display the configuration
	fmt.Printf("Server Port: %s\n", config.ServerPort)
	fmt.Printf("Database URL: %s\n", config.DatabaseURL)
	fmt.Printf("Log Level: %s\n", config.LogLevel)
	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("Features: %v\n", config.Features)

	if config.APIKey == "" {
		fmt.Println("Warning: API_KEY environment variable is not set")
	} else {
		fmt.Println("API Key: [REDACTED]")
	}

	// Validate the configuration
	errors := config.Validate()
	if len(errors) > 0 {
		fmt.Println("\nConfiguration Errors:")
		for _, err := range errors {
			fmt.Printf("- %s\n", err)
		}
	} else {
		fmt.Println("\nConfiguration is valid.")
	}

	fmt.Println("\nTry running this example with different environment variables set:")
	fmt.Println("export SERVER_PORT=9090")
	fmt.Println("export DATABASE_URL=\"postgres://user:password@localhost:5432/mydb\"")
	fmt.Println("export LOG_LEVEL=\"debug\"")
	fmt.Println("export API_KEY=\"your-api-key\"")
	fmt.Println("export APP_ENV=\"production\"")
	fmt.Println("export FEATURES=\"basic,premium,advanced\"")
 fmt.Println("go run EXAMPLES/env/config_integration_example.go")
}
