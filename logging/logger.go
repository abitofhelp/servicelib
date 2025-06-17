// Copyright (c) 2025 A Bit of Help, Inc.

// Package logging provides centralized logging functionality for the family service.
// It wraps the zap logging library and adds features like trace ID extraction from context
// and context-aware logging methods. This package is part of the infrastructure layer
// and provides logging capabilities to all other layers of the application.
package logging

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger configured based on the provided level and environment.
// It sets up appropriate encoders and log levels for either development or production use.
// Parameters:
//   - level: The minimum log level as a string (e.g., "debug", "info", "warn", "error")
//   - development: Whether to use development mode with console output (true) or production mode with JSON output (false)
//
// Returns:
//   - *zap.Logger: A configured zap logger instance
//   - error: An error if logger creation fails
func NewLogger(level string, development bool) (*zap.Logger, error) {
	// Parse log level
	var zapLevel zapcore.Level
	err := zapLevel.UnmarshalText([]byte(level))
	if err != nil {
		// Default to info level if parsing fails
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	var core zapcore.Core
	if development {
		// In development mode, log to console with colored output
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLevel)
	} else {
		// In production mode, log as JSON
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapLevel)
	}

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// WithTraceID adds trace ID and span ID to the logger from the provided context.
// This enables correlation between logs and traces for distributed tracing.
// Parameters:
//   - ctx: The context containing trace information
//   - logger: The base logger to enhance with trace information
//
// Returns:
//   - *zap.Logger: A new logger with trace ID and span ID fields added if available
func WithTraceID(ctx context.Context, logger *zap.Logger) *zap.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return logger.With(
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
		)
	}
	return logger
}

// ContextLogger is a logger that includes context information in log entries.
// It wraps a zap.Logger and provides methods that accept a context parameter,
// automatically extracting and including trace information in log entries.
type ContextLogger struct {
	// base is the underlying zap logger instance
	base *zap.Logger
}

// NewContextLogger creates a new context-aware logger wrapping the provided base logger.
// If base is nil, a no-op logger will be used to prevent nil pointer panics.
// Parameters:
//   - base: The base zap logger to wrap
//
// Returns:
//   - *ContextLogger: A new context logger instance
func NewContextLogger(base *zap.Logger) *ContextLogger {
	// If base logger is nil, use a no-op logger to prevent nil pointer panics
	if base == nil {
		base = zap.NewNop()
	}
	return &ContextLogger{
		base: base,
	}
}

// With returns a logger with trace information from the given context.
// Parameters:
//   - ctx: The context containing trace information
//
// Returns:
//   - *zap.Logger: A logger with trace ID and span ID fields added if available
func (l *ContextLogger) With(ctx context.Context) *zap.Logger {
	return WithTraceID(ctx, l.base)
}

// Debug logs a debug-level message with context information.
// Parameters:
//   - ctx: The context containing trace information
//   - msg: The message to log
//   - fields: Additional fields to include in the log entry
func (l *ContextLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(ctx).Debug(msg, fields...)
}

// Info logs an info-level message with context information.
// Parameters:
//   - ctx: The context containing trace information
//   - msg: The message to log
//   - fields: Additional fields to include in the log entry
func (l *ContextLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(ctx).Info(msg, fields...)
}

// Warn logs a warning-level message with context information.
// Parameters:
//   - ctx: The context containing trace information
//   - msg: The message to log
//   - fields: Additional fields to include in the log entry
func (l *ContextLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(ctx).Warn(msg, fields...)
}

// Error logs an error-level message with context information.
// Parameters:
//   - ctx: The context containing trace information
//   - msg: The message to log
//   - fields: Additional fields to include in the log entry
func (l *ContextLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(ctx).Error(msg, fields...)
}

// Fatal logs a fatal-level message with context information.
// This will terminate the program after logging the message.
// Parameters:
//   - ctx: The context containing trace information
//   - msg: The message to log
//   - fields: Additional fields to include in the log entry
func (l *ContextLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(ctx).Fatal(msg, fields...)
}

// Sync flushes any buffered log entries to their destination.
// This should be called before program termination to ensure all logs are written.
// Returns:
//   - error: An error if flushing fails
func (l *ContextLogger) Sync() error {
	return l.base.Sync()
}
