//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of implementing and using the AppConfig interface
package example_config

import (
	"fmt"

	"github.com/abitofhelp/servicelib/config"
)

// AppSettings is a custom configuration struct that implements the AppConfigProvider interface
type AppSettings struct {
	AppVersion     string
	AppName        string
	AppEnvironment string
	LogLevel       string
	Features       map[string]bool
}

// GetAppVersion implements the AppConfigProvider interface
func (s *AppSettings) GetAppVersion() string {
	return s.AppVersion
}

// GetAppName implements the AppConfigProvider interface
func (s *AppSettings) GetAppName() string {
	return s.AppName
}

// GetAppEnvironment implements the AppConfigProvider interface
func (s *AppSettings) GetAppEnvironment() string {
	return s.AppEnvironment
}

// Additional methods specific to AppSettings
func (s *AppSettings) GetLogLevel() string {
	return s.LogLevel
}

func (s *AppSettings) IsFeatureEnabled(featureName string) bool {
	if enabled, ok := s.Features[featureName]; ok {
		return enabled
	}
	return false
}

func main() {
	// Create app settings
	settings := &AppSettings{
		AppVersion:     "2.1.0",
		AppName:        "FeatureApp",
		AppEnvironment: "staging",
		LogLevel:       "debug",
		Features: map[string]bool{
			"darkMode":      true,
			"betaFeatures":  true,
			"notifications": false,
		},
	}

	// Create a config adapter
	adapter := config.NewGenericConfigAdapter(settings)

	// Get the app configuration through the adapter
	appConfig := adapter.GetApp()

	// Use the standard AppConfig interface methods
	fmt.Println("=== Application Configuration ===")
	fmt.Println("Version:", appConfig.GetVersion())
	fmt.Println("Name:", appConfig.GetName())
	fmt.Println("Environment:", appConfig.GetEnvironment())

	// Use the original settings object for additional functionality
	fmt.Println("\n=== Additional Settings ===")
	fmt.Println("Log Level:", settings.GetLogLevel())
	fmt.Println("Dark Mode Enabled:", settings.IsFeatureEnabled("darkMode"))
	fmt.Println("Beta Features Enabled:", settings.IsFeatureEnabled("betaFeatures"))
	fmt.Println("Notifications Enabled:", settings.IsFeatureEnabled("notifications"))

	// Expected output:
	// === Application Configuration ===
	// Version: 2.1.0
	// Name: FeatureApp
	// Environment: staging
	//
	// === Additional Settings ===
	// Log Level: debug
	// Dark Mode Enabled: true
	// Beta Features Enabled: true
	// Notifications Enabled: false
}
