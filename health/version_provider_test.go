// Copyright (c) 2025 A Bit of Help, Inc.

package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfigVersionAdapter tests the NewConfigVersionAdapter function
func TestNewConfigVersionAdapter(t *testing.T) {
	// Create a config version adapter
	version := "1.0.0"
	adapter := NewConfigVersionAdapter(version)

	// Verify that the adapter is not nil
	assert.NotNil(t, adapter)
	assert.Equal(t, version, adapter.version)

	// Verify that the adapter's GetVersion method returns the expected value
	returnedVersion := adapter.GetVersion()
	assert.Equal(t, version, returnedVersion)
}
