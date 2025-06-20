// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestNewContainer tests the NewContainer function
func TestNewContainer(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Test case 1: Valid inputs
	t.Run("Valid inputs", func(t *testing.T) {
		// Create a config
		cfg := map[string]interface{}{
			"key": "value",
		}

		// Create a container
		container, err := NewContainer(context.Background(), logger, cfg)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, container)
		assert.Equal(t, cfg, container.GetConfig())
	})

	// Test case 2: Nil context
	t.Run("Nil context", func(t *testing.T) {
		// Create a config
		cfg := map[string]interface{}{
			"key": "value",
		}

		// Create a container with nil context
		container, err := NewContainer(nil, logger, cfg)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "context cannot be nil")
	})

	// Test case 3: Nil logger
	t.Run("Nil logger", func(t *testing.T) {
		// Create a config
		cfg := map[string]interface{}{
			"key": "value",
		}

		// Create a container with nil logger
		container, err := NewContainer(context.Background(), nil, cfg)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "logger cannot be nil")
	})

	// Test case 4: Nil config
	t.Run("Nil config", func(t *testing.T) {
		// Create a container with nil config
		container, err := NewContainer(context.Background(), logger, nil)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})
}

// TestGetRepositoryFactory tests the GetRepositoryFactory method
func TestGetRepositoryFactory(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Create a config
	cfg := map[string]interface{}{
		"key": "value",
	}

	// Create a container
	container, err := NewContainer(context.Background(), logger, cfg)
	assert.NoError(t, err)

	// Call the method
	factory := container.GetRepositoryFactory()

	// Assertions
	assert.Nil(t, factory)
}
