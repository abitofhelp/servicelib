// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestNewBaseContainer tests the NewBaseContainer function
func TestNewBaseContainer(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Test case 1: Valid inputs
	t.Run("Valid inputs", func(t *testing.T) {
		// Create a config
		cfg := "test-config"

		// Create a container
		container, err := NewBaseContainer(context.Background(), logger, cfg)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, container)
		assert.Equal(t, cfg, container.GetConfig())
		assert.NotNil(t, container.GetContext())
		assert.NotNil(t, container.GetLogger())
		assert.NotNil(t, container.GetContextLogger())
		assert.NotNil(t, container.GetValidator())
	})

	// Test case 2: Nil context
	t.Run("Nil context", func(t *testing.T) {
		// Create a config
		cfg := "test-config"

		// Create a container with nil context
		container, err := NewBaseContainer(nil, logger, cfg)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "context cannot be nil")
	})

	// Test case 3: Nil logger
	t.Run("Nil logger", func(t *testing.T) {
		// Create a config
		cfg := "test-config"

		// Create a container with nil logger
		container, err := NewBaseContainer(context.Background(), nil, cfg)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "logger cannot be nil")
	})
}

// TestBaseContainerGetters tests the getter methods of BaseContainer
func TestBaseContainerGetters(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Create a context
	ctx := context.Background()

	// Create a config
	cfg := "test-config"

	// Create a container
	container, err := NewBaseContainer(ctx, logger, cfg)
	assert.NoError(t, err)

	// Test GetContext
	t.Run("GetContext", func(t *testing.T) {
		assert.Equal(t, ctx, container.GetContext())
	})

	// Test GetLogger
	t.Run("GetLogger", func(t *testing.T) {
		assert.Equal(t, logger, container.GetLogger())
	})

	// Test GetContextLogger
	t.Run("GetContextLogger", func(t *testing.T) {
		assert.NotNil(t, container.GetContextLogger())
	})

	// Test GetValidator
	t.Run("GetValidator", func(t *testing.T) {
		assert.NotNil(t, container.GetValidator())
	})

	// Test GetConfig
	t.Run("GetConfig", func(t *testing.T) {
		assert.Equal(t, cfg, container.GetConfig())
	})
}

// TestBaseContainerClose tests the Close method of BaseContainer
func TestBaseContainerClose(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Create a config
	cfg := "test-config"

	// Create a container
	container, err := NewBaseContainer(context.Background(), logger, cfg)
	assert.NoError(t, err)

	// Call the Close method
	err = container.Close()

	// Assertions
	assert.NoError(t, err)
}
