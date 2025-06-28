package middleware

import (
	"context"
	"testing"
)

func TestHasRole(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		role     string
		expected bool
	}{
		{
			name:     "user has role",
			roles:    []string{"admin", "user"},
			role:     "admin",
			expected: true,
		},
		{
			name:     "user does not have role",
			roles:    []string{"user"},
			role:     "admin",
			expected: false,
		},
		{
			name:     "empty roles",
			roles:    []string{},
			role:     "admin",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = WithUserRoles(ctx, tt.roles)

			if got := HasRole(ctx, tt.role); got != tt.expected {
				t.Errorf("HasRole() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAuthorized(t *testing.T) {
	tests := []struct {
		name         string
		userRoles    []string
		allowedRoles []string
		expected     bool
	}{
		{
			name:         "user has one of the allowed roles",
			userRoles:    []string{"admin", "user"},
			allowedRoles: []string{"admin", "superuser"},
			expected:     true,
		},
		{
			name:         "user has none of the allowed roles",
			userRoles:    []string{"user"},
			allowedRoles: []string{"admin", "superuser"},
			expected:     false,
		},
		{
			name:         "empty user roles",
			userRoles:    []string{},
			allowedRoles: []string{"admin", "superuser"},
			expected:     false,
		},
		{
			name:         "empty allowed roles",
			userRoles:    []string{"admin", "user"},
			allowedRoles: []string{},
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = WithUserRoles(ctx, tt.userRoles)

			if got := IsAuthorized(ctx, tt.allowedRoles); got != tt.expected {
				t.Errorf("IsAuthorized() = %v, want %v", got, tt.expected)
			}
		})
	}
}
