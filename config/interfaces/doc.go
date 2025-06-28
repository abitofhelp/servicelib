// Copyright (c) 2025 A Bit of Help, Inc.

// Package interfaces defines the core configuration interfaces used throughout the application.
//
// This package contains the fundamental interfaces that establish the contract for
// configuration access across the application. By defining these interfaces separately
// from their implementations, the package promotes loose coupling and enables
// dependency inversion.
//
// The interfaces defined here are:
//   - AppConfig: For accessing application-specific configuration
//   - DatabaseConfig: For accessing database-specific configuration
//   - Config: A composite interface that provides access to both AppConfig and DatabaseConfig
//
// These interfaces are designed to be implementation-agnostic, allowing different
// configuration sources (files, environment variables, remote services, etc.) to
// provide the same consistent access patterns.
//
// Example usage:
//
//	// Using the Config interface
//	func InitializeApp(cfg interfaces.Config) {
//	    appName := cfg.GetApp().GetName()
//	    appVersion := cfg.GetApp().GetVersion()
//	    dbConnStr := cfg.GetDatabase().GetConnectionString()
//	
//	    fmt.Printf("Initializing %s v%s\n", appName, appVersion)
//	    fmt.Printf("Connecting to database: %s\n", dbConnStr)
//	}
//
// The interfaces in this package are typically implemented by adapter types in the
// parent config package, but can be implemented by any type that satisfies the
// required methods.
package interfaces