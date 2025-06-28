// Copyright (c) 2025 A Bit of Help, Inc.

package retry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestWithOtelTracer(t *testing.T) {
	// Create a mock tracer provider
	mockTracer := trace.NewNoopTracerProvider().Tracer("test")

	// Create default options
	options := DefaultOptions()

	// Apply WithOtelTracer
	newOptions := options.WithOtelTracer(mockTracer)

	// Verify that the tracer was set
	assert.NotNil(t, newOptions.Tracer)
	assert.NotEqual(t, options.Tracer, newOptions.Tracer)
}

func TestWithConfigEdgeCases(t *testing.T) {
	// Test WithMaxRetries with negative value
	t.Run("WithMaxRetries negative", func(t *testing.T) {
		config := DefaultConfig()
		newConfig := config.WithMaxRetries(-1)
		assert.Equal(t, 0, newConfig.MaxRetries)
	})

	// Test WithInitialBackoff with zero value
	t.Run("WithInitialBackoff zero", func(t *testing.T) {
		config := DefaultConfig()
		newConfig := config.WithInitialBackoff(0)
		assert.Equal(t, int64(1*1000*1000), newConfig.InitialBackoff.Nanoseconds())
	})

	// Test WithMaxBackoff with zero value
	t.Run("WithMaxBackoff zero", func(t *testing.T) {
		config := DefaultConfig()
		newConfig := config.WithMaxBackoff(0)
		assert.Equal(t, config.InitialBackoff, newConfig.MaxBackoff)
	})

	// Test WithBackoffFactor with zero value
	t.Run("WithBackoffFactor zero", func(t *testing.T) {
		config := DefaultConfig()
		newConfig := config.WithBackoffFactor(0)
		assert.Equal(t, 1.0, newConfig.BackoffFactor)
	})

	// Test WithJitterFactor with negative value
	t.Run("WithJitterFactor negative", func(t *testing.T) {
		config := DefaultConfig()
		newConfig := config.WithJitterFactor(-0.5)
		assert.Equal(t, 0.0, newConfig.JitterFactor)
	})

	// Test WithJitterFactor with value > 1
	t.Run("WithJitterFactor > 1", func(t *testing.T) {
		config := DefaultConfig()
		newConfig := config.WithJitterFactor(1.5)
		assert.Equal(t, 1.0, newConfig.JitterFactor)
	})
}
