// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasScope(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []string
		scope    string
		expected bool
	}{
		{
			name:     "user has scope",
			scopes:   []string{"read", "write"},
			scope:    "read",
			expected: true,
		},
		{
			name:     "user does not have scope",
			scopes:   []string{"read"},
			scope:    "write",
			expected: false,
		},
		{
			name:     "empty scopes",
			scopes:   []string{},
			scope:    "read",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = WithUserScopes(ctx, tt.scopes)

			if got := HasScope(ctx, tt.scope); got != tt.expected {
				t.Errorf("HasScope() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasResource(t *testing.T) {
	tests := []struct {
		name      string
		resources []string
		resource  string
		expected  bool
	}{
		{
			name:      "user has resource",
			resources: []string{"resource1", "resource2"},
			resource:  "resource1",
			expected:  true,
		},
		{
			name:      "user does not have resource",
			resources: []string{"resource1"},
			resource:  "resource2",
			expected:  false,
		},
		{
			name:      "empty resources",
			resources: []string{},
			resource:  "resource1",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = WithUserResources(ctx, tt.resources)

			if got := HasResource(ctx, tt.resource); got != tt.expected {
				t.Errorf("HasResource() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAuthorizedWithScopes(t *testing.T) {
	tests := []struct {
		name           string
		userRoles      []string
		userScopes     []string
		userResources  []string
		allowedRoles   []string
		requiredScopes []string
		resource       string
		expected       bool
	}{
		{
			name:           "user has role, scopes, and resource",
			userRoles:      []string{"admin", "user"},
			userScopes:     []string{"read", "write"},
			userResources:  []string{"resource1", "resource2"},
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{"read"},
			resource:       "resource1",
			expected:       true,
		},
		{
			name:           "user has role and resource, but not all scopes",
			userRoles:      []string{"admin"},
			userScopes:     []string{"read"},
			userResources:  []string{"resource1"},
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{"read", "write"},
			resource:       "resource1",
			expected:       false,
		},
		{
			name:           "user has role and scopes, but not resource",
			userRoles:      []string{"admin"},
			userScopes:     []string{"read", "write"},
			userResources:  []string{"resource2"},
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{"read"},
			resource:       "resource1",
			expected:       false,
		},
		{
			name:           "user does not have role",
			userRoles:      []string{"user"},
			userScopes:     []string{"read", "write"},
			userResources:  []string{"resource1"},
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{"read"},
			resource:       "resource1",
			expected:       false,
		},
		{
			name:           "no required scopes",
			userRoles:      []string{"admin"},
			userScopes:     []string{},
			userResources:  []string{"resource1"},
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{},
			resource:       "resource1",
			expected:       true,
		},
		{
			name:           "empty allowed roles",
			userRoles:      []string{"admin"},
			userScopes:     []string{"read"},
			userResources:  []string{"resource1"},
			allowedRoles:   []string{},
			requiredScopes: []string{"read"},
			resource:       "resource1",
			expected:       false,
		},
		{
			name:           "no user scopes",
			userRoles:      []string{"admin"},
			userScopes:     nil,
			userResources:  []string{"resource1"},
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{"read"},
			resource:       "resource1",
			expected:       false,
		},
		{
			name:           "no user resources",
			userRoles:      []string{"admin"},
			userScopes:     []string{"read"},
			userResources:  nil,
			allowedRoles:   []string{"admin"},
			requiredScopes: []string{"read"},
			resource:       "resource1",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = WithUserRoles(ctx, tt.userRoles)
			if tt.userScopes != nil {
				ctx = WithUserScopes(ctx, tt.userScopes)
			}
			if tt.userResources != nil {
				ctx = WithUserResources(ctx, tt.userResources)
			}

			if got := IsAuthorizedWithScopes(ctx, tt.allowedRoles, tt.requiredScopes, tt.resource); got != tt.expected {
				t.Errorf("IsAuthorizedWithScopes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected bool
	}{
		{
			name:     "user is authenticated",
			userID:   "user123",
			expected: true,
		},
		{
			name:     "user is not authenticated",
			userID:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.userID != "" {
				ctx = WithUserID(ctx, tt.userID)
			}

			if got := IsAuthenticated(ctx); got != tt.expected {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	ctx := context.Background()
	userID, ok := GetUserID(ctx)
	assert.Empty(t, userID)
	assert.False(t, ok)

	expectedUserID := "user123"
	ctx = WithUserID(ctx, expectedUserID)
	userID, ok = GetUserID(ctx)
	assert.Equal(t, expectedUserID, userID)
	assert.True(t, ok)
}

func TestGetUserRoles(t *testing.T) {
	ctx := context.Background()
	roles, ok := GetUserRoles(ctx)
	assert.Nil(t, roles)
	assert.False(t, ok)

	expectedRoles := []string{"admin", "user"}
	ctx = WithUserRoles(ctx, expectedRoles)
	roles, ok = GetUserRoles(ctx)
	assert.Equal(t, expectedRoles, roles)
	assert.True(t, ok)
}

func TestGetUserScopes(t *testing.T) {
	ctx := context.Background()
	scopes, ok := GetUserScopes(ctx)
	assert.Nil(t, scopes)
	assert.False(t, ok)

	expectedScopes := []string{"read", "write"}
	ctx = WithUserScopes(ctx, expectedScopes)
	scopes, ok = GetUserScopes(ctx)
	assert.Equal(t, expectedScopes, scopes)
	assert.True(t, ok)
}

func TestGetUserResources(t *testing.T) {
	ctx := context.Background()
	resources, ok := GetUserResources(ctx)
	assert.Nil(t, resources)
	assert.False(t, ok)

	expectedResources := []string{"resource1", "resource2"}
	ctx = WithUserResources(ctx, expectedResources)
	resources, ok = GetUserResources(ctx)
	assert.Equal(t, expectedResources, resources)
	assert.True(t, ok)
}
