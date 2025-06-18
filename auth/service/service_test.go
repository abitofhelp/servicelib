// Copyright (c) 2025 A Bit of Help, Inc.

package service_test

import (
	"context"
	"testing"

	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/abitofhelp/servicelib/auth/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestNewService tests the creation of a new service
func TestNewService(t *testing.T) {
	// Test with logger
	logger := zap.NewExample()
	config := service.Config{
		AdminRoleName:         "admin",
		ReadOnlyRoleName:      "reader",
		ReadOperationPrefixes: []string{"read:", "get:"},
	}
	svc := service.NewService(config, logger)
	assert.NotNil(t, svc)

	// Test with nil logger (should use NopLogger)
	svc = service.NewService(config, nil)
	assert.NotNil(t, svc)

	// Test with default config
	defaultConfig := service.DefaultConfig()
	svc = service.NewService(defaultConfig, logger)
	assert.NotNil(t, svc)
}

// TestIsAuthorized tests the IsAuthorized function
func TestIsAuthorized(t *testing.T) {
	logger := zap.NewNop()
	config := service.Config{
		AdminRoleName:         "admin",
		ReadOnlyRoleName:      "reader",
		ReadOperationPrefixes: []string{"read:", "get:"},
	}
	svc := service.NewService(config, logger)

	// Test cases
	tests := []struct {
		name       string
		setupCtx   func() context.Context
		operation  string
		authorized bool
		hasError   bool
	}{
		{
			name: "Unauthenticated user",
			setupCtx: func() context.Context {
				return context.Background()
			},
			operation:  "read:resource",
			authorized: false,
			hasError:   true,
		},
		{
			name: "User with no roles",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				return ctx
			},
			operation:  "read:resource",
			authorized: false,
			hasError:   true,
		},
		{
			name: "Admin user can do anything",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "admin123")
				ctx = middleware.WithUserRoles(ctx, []string{"admin"})
				return ctx
			},
			operation:  "write:resource",
			authorized: true,
			hasError:   false,
		},
		{
			name: "Reader user can read",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "reader123")
				ctx = middleware.WithUserRoles(ctx, []string{"reader"})
				return ctx
			},
			operation:  "read:resource",
			authorized: true,
			hasError:   false,
		},
		{
			name: "Reader user cannot write",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "reader123")
				ctx = middleware.WithUserRoles(ctx, []string{"reader"})
				return ctx
			},
			operation:  "write:resource",
			authorized: false,
			hasError:   true,
		},
		{
			name: "Regular user cannot do anything",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				ctx = middleware.WithUserRoles(ctx, []string{"user"})
				return ctx
			},
			operation:  "read:resource",
			authorized: false,
			hasError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setupCtx()
			authorized, err := svc.IsAuthorized(ctx, tc.operation)
			assert.Equal(t, tc.authorized, authorized)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestIsAdmin tests the IsAdmin function
func TestIsAdmin(t *testing.T) {
	logger := zap.NewNop()
	config := service.Config{
		AdminRoleName: "admin",
	}
	svc := service.NewService(config, logger)

	// Test cases
	tests := []struct {
		name     string
		setupCtx func() context.Context
		isAdmin  bool
		hasError bool
	}{
		{
			name: "Unauthenticated user",
			setupCtx: func() context.Context {
				return context.Background()
			},
			isAdmin:  false,
			hasError: true,
		},
		{
			name: "User with admin role",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "admin123")
				ctx = middleware.WithUserRoles(ctx, []string{"admin"})
				return ctx
			},
			isAdmin:  true,
			hasError: false,
		},
		{
			name: "User without admin role",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				ctx = middleware.WithUserRoles(ctx, []string{"user"})
				return ctx
			},
			isAdmin:  false,
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setupCtx()
			isAdmin, err := svc.IsAdmin(ctx)
			assert.Equal(t, tc.isAdmin, isAdmin)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestHasRole tests the HasRole function
func TestHasRole(t *testing.T) {
	logger := zap.NewNop()
	config := service.Config{}
	svc := service.NewService(config, logger)

	// Test cases
	tests := []struct {
		name     string
		setupCtx func() context.Context
		role     string
		hasRole  bool
		hasError bool
	}{
		{
			name: "Unauthenticated user",
			setupCtx: func() context.Context {
				return context.Background()
			},
			role:     "user",
			hasRole:  false,
			hasError: true,
		},
		{
			name: "User with the role",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				ctx = middleware.WithUserRoles(ctx, []string{"user", "editor"})
				return ctx
			},
			role:     "editor",
			hasRole:  true,
			hasError: false,
		},
		{
			name: "User without the role",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				ctx = middleware.WithUserRoles(ctx, []string{"user"})
				return ctx
			},
			role:     "admin",
			hasRole:  false,
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setupCtx()
			hasRole, err := svc.HasRole(ctx, tc.role)
			assert.Equal(t, tc.hasRole, hasRole)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetUserID tests the GetUserID function
func TestGetUserID(t *testing.T) {
	logger := zap.NewNop()
	config := service.Config{}
	svc := service.NewService(config, logger)

	// Test cases
	tests := []struct {
		name     string
		setupCtx func() context.Context
		userID   string
		hasError bool
	}{
		{
			name: "Unauthenticated user",
			setupCtx: func() context.Context {
				return context.Background()
			},
			userID:   "",
			hasError: true,
		},
		{
			name: "Authenticated user",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				return ctx
			},
			userID:   "user123",
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setupCtx()
			userID, err := svc.GetUserID(ctx)
			assert.Equal(t, tc.userID, userID)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetUserRoles tests the GetUserRoles function
func TestGetUserRoles(t *testing.T) {
	logger := zap.NewNop()
	config := service.Config{}
	svc := service.NewService(config, logger)

	// Test cases
	tests := []struct {
		name     string
		setupCtx func() context.Context
		roles    []string
		hasError bool
	}{
		{
			name: "Unauthenticated user",
			setupCtx: func() context.Context {
				return context.Background()
			},
			roles:    nil,
			hasError: true,
		},
		{
			name: "User with roles",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				ctx = middleware.WithUserRoles(ctx, []string{"user", "editor"})
				return ctx
			},
			roles:    []string{"user", "editor"},
			hasError: false,
		},
		{
			name: "User with no roles",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = middleware.WithUserID(ctx, "user123")
				ctx = middleware.WithUserRoles(ctx, []string{})
				return ctx
			},
			roles:    []string{},
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setupCtx()
			roles, err := svc.GetUserRoles(ctx)
			assert.Equal(t, tc.roles, roles)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
