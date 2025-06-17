// Copyright (c) 2025 A Bit of Help, Inc.

package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	// Setup
	config := &Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8089/callback",
		Scopes:       []string{"profile", "email"},
	}

	// Execute
	provider := NewProvider(config)

	// Verify
	assert.NotNil(t, provider)
	assert.Equal(t, config, provider.Config)
}
