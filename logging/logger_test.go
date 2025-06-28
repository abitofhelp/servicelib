// Copyright (c) 2025 A Bit of Help, Inc.

package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		level       string
		development bool
		expectError bool
	}{
		{
			name:        "Debug level in development mode",
			level:       "debug",
			development: true,
			expectError: false,
		},
		{
			name:        "Info level in production mode",
			level:       "info",
			development: false,
			expectError: false,
		},
		{
			name:        "Invalid level defaults to info",
			level:       "invalid",
			development: false,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute
			logger, err := NewLogger(tc.level, tc.development)

			// Verify
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}

func TestWithTraceID(t *testing.T) {
	// Setup
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	// Create a context with a trace
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	// Execute
	loggerWithTrace := WithTraceID(ctx, logger)

	// Log a message
	loggerWithTrace.Info("test message")

	// Verify
	logs := recorded.All()
	require.Len(t, logs, 1)

	// Check that trace_id and span_id fields are present
	assert.Equal(t, "test message", logs[0].Message)

	traceIDField, ok := logs[0].ContextMap()["trace_id"]
	assert.True(t, ok)
	assert.Equal(t, traceID.String(), traceIDField)

	spanIDField, ok := logs[0].ContextMap()["span_id"]
	assert.True(t, ok)
	assert.Equal(t, spanID.String(), spanIDField)
}

func TestWithTraceID_NoTrace(t *testing.T) {
	// Setup
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	// Create a context without a trace
	ctx := context.Background()

	// Execute
	loggerWithTrace := WithTraceID(ctx, logger)

	// Log a message
	loggerWithTrace.Info("test message")

	// Verify
	logs := recorded.All()
	require.Len(t, logs, 1)

	// Check that trace_id and span_id fields are not present
	assert.Equal(t, "test message", logs[0].Message)

	_, ok := logs[0].ContextMap()["trace_id"]
	assert.False(t, ok)

	_, ok = logs[0].ContextMap()["span_id"]
	assert.False(t, ok)
}

func TestNewContextLogger(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)

	// Execute
	contextLogger := NewContextLogger(logger)

	// Verify
	assert.NotNil(t, contextLogger)
	assert.Equal(t, logger, contextLogger.base)
}

func TestContextLogger_With(t *testing.T) {
	// Setup
	core, _ := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	contextLogger := NewContextLogger(logger)

	// Create a context with a trace
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	// Execute
	loggerWithContext := contextLogger.With(ctx)

	// Verify
	assert.NotNil(t, loggerWithContext)
	assert.NotEqual(t, logger, loggerWithContext)
}

func TestContextLogger_LogMethods(t *testing.T) {
	// Setup
	core, recorded := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	contextLogger := NewContextLogger(logger)

	// Create a context with a trace
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	// Test cases for different log levels
	testCases := []struct {
		name    string
		logFunc func(context.Context, string, ...zap.Field)
		level   zapcore.Level
		message string
	}{
		{
			name:    "Debug",
			logFunc: contextLogger.Debug,
			level:   zapcore.DebugLevel,
			message: "debug message",
		},
		{
			name:    "Info",
			logFunc: contextLogger.Info,
			level:   zapcore.InfoLevel,
			message: "info message",
		},
		{
			name:    "Warn",
			logFunc: contextLogger.Warn,
			level:   zapcore.WarnLevel,
			message: "warn message",
		},
		{
			name:    "Error",
			logFunc: contextLogger.Error,
			level:   zapcore.ErrorLevel,
			message: "error message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear previous logs
			recorded.TakeAll()

			// Execute
			tc.logFunc(ctx, tc.message, zap.String("key", "value"))

			// Verify
			logs := recorded.All()
			require.Len(t, logs, 1)

			assert.Equal(t, tc.level, logs[0].Level)
			assert.Equal(t, tc.message, logs[0].Message)

			// Check fields
			assert.Equal(t, "value", logs[0].ContextMap()["key"])
			assert.Equal(t, traceID.String(), logs[0].ContextMap()["trace_id"])
			assert.Equal(t, spanID.String(), logs[0].ContextMap()["span_id"])
		})
	}
}

func TestContextLogger_Sync(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)
	contextLogger := NewContextLogger(logger)

	// Execute
	err := contextLogger.Sync()

	// Verify
	assert.NoError(t, err)
}
