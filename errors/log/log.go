// Copyright (c) 2025 A Bit of Help, Inc.

// Package log provides logging integration for the error handling system.
package log

import (
	"context"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogError logs an error with context information.
// It extracts error details like code, operation, and details if available.
func LogError(ctx context.Context, logger *logging.ContextLogger, err error) {
	if err == nil {
		return
	}

	// Create base fields
	fields := []zap.Field{
		zap.Error(err),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		fields = append(fields, zap.String("error_code", string(e.GetCode())))
	}

	// Add operation if available
	if e, ok := err.(interface{ GetOperation() string }); ok && e.GetOperation() != "" {
		fields = append(fields, zap.String("operation", e.GetOperation()))
	}

	// Add source location if available
	if e, ok := err.(interface {
		GetSource() string
		GetLine() int
	}); ok && e.GetSource() != "" {
		fields = append(fields, zap.String("source", e.GetSource()))
		if e.GetLine() > 0 {
			fields = append(fields, zap.Int("line", e.GetLine()))
		}
	}

	// Add details if available
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		details := e.GetDetails()
		if details != nil {
			for k, v := range details {
				fields = append(fields, zap.Any("detail_"+k, v))
			}
		}
	}

	// Add HTTP status if available
	if e, ok := err.(interface{ GetHTTPStatus() int }); ok && e.GetHTTPStatus() != 0 {
		fields = append(fields, zap.Int("http_status", e.GetHTTPStatus()))
	}

	// Log the error
	logger.Error(ctx, "Error occurred", fields...)
}

// LogErrorWithLevel logs an error with a specific log level.
func LogErrorWithLevel(ctx context.Context, logger *logging.ContextLogger, err error, level zapcore.Level) {
	if err == nil {
		return
	}

	// Create base fields
	fields := []zap.Field{
		zap.Error(err),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		fields = append(fields, zap.String("error_code", string(e.GetCode())))
	}

	// Add operation if available
	if e, ok := err.(interface{ GetOperation() string }); ok && e.GetOperation() != "" {
		fields = append(fields, zap.String("operation", e.GetOperation()))
	}

	// Add source location if available
	if e, ok := err.(interface {
		GetSource() string
		GetLine() int
	}); ok && e.GetSource() != "" {
		fields = append(fields, zap.String("source", e.GetSource()))
		if e.GetLine() > 0 {
			fields = append(fields, zap.Int("line", e.GetLine()))
		}
	}

	// Add details if available
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		details := e.GetDetails()
		if details != nil {
			for k, v := range details {
				fields = append(fields, zap.Any("detail_"+k, v))
			}
		}
	}

	// Add HTTP status if available
	if e, ok := err.(interface{ GetHTTPStatus() int }); ok && e.GetHTTPStatus() != 0 {
		fields = append(fields, zap.Int("http_status", e.GetHTTPStatus()))
	}

	// Log the error with the specified level
	switch level {
	case zapcore.DebugLevel:
		logger.Debug(ctx, "Error occurred", fields...)
	case zapcore.InfoLevel:
		logger.Info(ctx, "Error occurred", fields...)
	case zapcore.WarnLevel:
		logger.Warn(ctx, "Error occurred", fields...)
	case zapcore.ErrorLevel:
		logger.Error(ctx, "Error occurred", fields...)
	case zapcore.FatalLevel:
		logger.Fatal(ctx, "Error occurred", fields...)
	default:
		logger.Error(ctx, "Error occurred", fields...)
	}
}

// LogErrorWithMessage logs an error with a custom message.
func LogErrorWithMessage(ctx context.Context, logger *logging.ContextLogger, err error, message string) {
	if err == nil {
		return
	}

	// Create base fields
	fields := []zap.Field{
		zap.Error(err),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		fields = append(fields, zap.String("error_code", string(e.GetCode())))
	}

	// Add operation if available
	if e, ok := err.(interface{ GetOperation() string }); ok && e.GetOperation() != "" {
		fields = append(fields, zap.String("operation", e.GetOperation()))
	}

	// Add source location if available
	if e, ok := err.(interface {
		GetSource() string
		GetLine() int
	}); ok && e.GetSource() != "" {
		fields = append(fields, zap.String("source", e.GetSource()))
		if e.GetLine() > 0 {
			fields = append(fields, zap.Int("line", e.GetLine()))
		}
	}

	// Add details if available
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		details := e.GetDetails()
		if details != nil {
			for k, v := range details {
				fields = append(fields, zap.Any("detail_"+k, v))
			}
		}
	}

	// Add HTTP status if available
	if e, ok := err.(interface{ GetHTTPStatus() int }); ok && e.GetHTTPStatus() != 0 {
		fields = append(fields, zap.Int("http_status", e.GetHTTPStatus()))
	}

	// Log the error with the custom message
	logger.Error(ctx, message, fields...)
}
