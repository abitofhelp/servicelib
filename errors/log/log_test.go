// Copyright (c) 2025 A Bit of Help, Inc.

package log

import (
	"context"
	"testing"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// TestLogError tests the LogError function
func TestLogError(t *testing.T) {
	// Create a test logger
	zapLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(zapLogger)

	// Create a context
	ctx := context.Background()

	// Test with nil error
	LogError(ctx, logger, nil)
	// No assertions needed as this just verifies it doesn't panic

	// Test with standard error
	stdErr := errors.New(core.ValidationErrorCode, "Standard error")
	LogError(ctx, logger, stdErr)
	// This test just verifies the function doesn't panic

	// Test with error with operation
	opErr := errors.WrapWithOperation(stdErr, core.DatabaseErrorCode, "Operation error", "GetUserByID")
	LogError(ctx, logger, opErr)
	// This test just verifies the function doesn't panic

	// Test with error with details
	details := map[string]interface{}{
		"user_id": "123",
		"action":  "create",
	}
	detailsErr := errors.WrapWithDetails(stdErr, core.NotFoundCode, "Details error", details)
	LogError(ctx, logger, detailsErr)
	// This test just verifies the function doesn't panic
}

// TestLogErrorWithLevel tests the LogErrorWithLevel function
func TestLogErrorWithLevel(t *testing.T) {
	// Create a test logger
	zapLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(zapLogger)

	// Create a context
	ctx := context.Background()

	// Test with nil error
	LogErrorWithLevel(ctx, logger, nil, zapcore.ErrorLevel)
	// No assertions needed as this just verifies it doesn't panic

	// Test with debug level
	err := errors.New(core.ValidationErrorCode, "Debug error")
	LogErrorWithLevel(ctx, logger, err, zapcore.DebugLevel)
	// This test just verifies the function doesn't panic

	// Test with info level
	err = errors.New(core.ValidationErrorCode, "Info error")
	LogErrorWithLevel(ctx, logger, err, zapcore.InfoLevel)
	// This test just verifies the function doesn't panic

	// Test with warn level
	err = errors.New(core.ValidationErrorCode, "Warn error")
	LogErrorWithLevel(ctx, logger, err, zapcore.WarnLevel)
	// This test just verifies the function doesn't panic

	// Test with error level
	err = errors.New(core.ValidationErrorCode, "Error error")
	LogErrorWithLevel(ctx, logger, err, zapcore.ErrorLevel)
	// This test just verifies the function doesn't panic

	// Skip fatal level test as it would exit the program
	// Instead, test with a custom level that doesn't exist
	err = errors.New(core.ValidationErrorCode, "Invalid level error")
	LogErrorWithLevel(ctx, logger, err, zapcore.Level(100))
	// This test just verifies the function doesn't panic
}

// TestLogErrorWithMessage tests the LogErrorWithMessage function
func TestLogErrorWithMessage(t *testing.T) {
	// Create a test logger
	zapLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(zapLogger)

	// Create a context
	ctx := context.Background()

	// Test with nil error
	LogErrorWithMessage(ctx, logger, nil, "Custom message")
	// No assertions needed as this just verifies it doesn't panic

	// Test with standard error
	err := errors.New(core.ValidationErrorCode, "Standard error")
	LogErrorWithMessage(ctx, logger, err, "Custom message")
	// This test just verifies the function doesn't panic
}
