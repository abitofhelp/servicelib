//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of integrating environment variables with a configuration structure
package example_env

import (
	"fmt"
	"strings"

	"github.com/abitofhelp/servicelib/env"
)

// AppConfig holds the application configuration
type AppConfig struct {
	ServerPort  string
	DatabaseURL string
	LogLevel    string
	APIKey      string
	Environment string
	Features    []string
}

// LoadConfig loads the application configuration from environment variables
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

// Validate validates the configuration
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
