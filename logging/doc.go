// Copyright (c) 2025 A Bit of Help, Inc.

// Package logging provides centralized logging functionality for services.
//
// This package wraps the zap logging library and adds features like trace ID extraction
// from context and context-aware logging methods. It is designed to be used throughout
// the application to provide consistent, structured logging with proper context information.
//
// Key features:
//   - Context-aware logging with automatic trace ID extraction
//   - Support for different log levels (Debug, Info, Warn, Error, Fatal)
//   - Configurable for both development (console output) and production (JSON output) environments
//   - Integration with OpenTelemetry for distributed tracing
//   - Structured logging with additional fields
//
// The package provides two main components:
//   - ContextLogger: A wrapper around zap.Logger that automatically extracts trace information
//     from context and includes it in log entries
//   - Logger interface: An abstraction that can be implemented by different logging backends
//
// Example usage:
//
//	// Create a new logger
//	zapLogger, err := logging.NewLogger("info", true)
//	if err != nil {
//	    panic(err)
//	}
//	defer zapLogger.Sync()
//
//	// Create a context logger
//	logger := logging.NewContextLogger(zapLogger)
//
//	// Log with context
//	ctx := context.Background()
//	logger.Info(ctx, "Application started", zap.String("app_name", "example"))
//
//	// Log an error
//	if err := someOperation(); err != nil {
//	    logger.Error(ctx, "Operation failed", zap.Error(err))
//	}
//
// The package is designed to be used as a dependency by other packages in the application,
// providing a consistent logging interface throughout the codebase.
package logging
