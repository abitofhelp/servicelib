// Copyright (c) 2025 A Bit of Help, Inc.

// Package interfaces defines the interfaces for the logging package.
// It provides abstractions for logging functionality that can be implemented
// by different logging backends.
package interfaces

import (
	"context"

	"go.uber.org/zap"
)

// Logger defines the interface for context-aware logging
type Logger interface {
	// Debug logs a debug-level message with context information
	Debug(ctx context.Context, msg string, fields ...zap.Field)

	// Info logs an info-level message with context information
	Info(ctx context.Context, msg string, fields ...zap.Field)

	// Warn logs a warning-level message with context information
	Warn(ctx context.Context, msg string, fields ...zap.Field)

	// Error logs an error-level message with context information
	Error(ctx context.Context, msg string, fields ...zap.Field)

	// Fatal logs a fatal-level message with context information
	Fatal(ctx context.Context, msg string, fields ...zap.Field)

	// Sync flushes any buffered log entries
	Sync() error
}