// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to get user information from the context
package main

import (
	"context"
	"fmt"
)

// This is a simplified example to demonstrate user information retrieval
// In a real application, you would use the auth module's methods

// UserInfo represents user information
type UserInfo struct {
	ID    string
	Roles []string
	Email string
}

// UserInfoService simulates a user information service
type UserInfoService struct {
	users map[string]UserInfo
}

// NewUserInfoService creates a new user information service
func NewUserInfoService() *UserInfoService {
	return &UserInfoService{
		users: map[string]UserInfo{
			"user123": {
				ID:    "user123",
				Roles: []string{"admin", "editor"},
				Email: "admin@example.com",
			},
			"user456": {
				ID:    "user456",
				Roles: []string{"reader"},
				Email: "reader@example.com",
			},
		},
	}
}

// GetUserID gets the user ID from the context
func (a *UserInfoService) GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetUserRoles gets the user roles from the context
func (a *UserInfoService) GetUserRoles(ctx context.Context) ([]string, error) {
	userID, err := a.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, ok := a.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return user.Roles, nil
}

// GetUserEmail gets the user email from the context
func (a *UserInfoService) GetUserEmail(ctx context.Context) (string, error) {
	userID, err := a.GetUserID(ctx)
	if err != nil {
		return "", err
	}

	user, ok := a.users[userID]
	if !ok {
		return "", fmt.Errorf("user not found")
	}

	return user.Email, nil
}

// AddUserIDToContext adds a user ID to the context
func AddUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func main() {
	// Create a user information service
	userInfo := NewUserInfoService()

	// Create a context with user ID
	ctx := context.Background()
	adminCtx := AddUserIDToContext(ctx, "user123")

	// Get user ID
	userID, err := userInfo.GetUserID(adminCtx)
	if err != nil {
		fmt.Printf("Error getting user ID: %v\n", err)
		return
	}
	fmt.Printf("User ID: %s\n", userID)

	// Get user roles
	roles, err := userInfo.GetUserRoles(adminCtx)
	if err != nil {
		fmt.Printf("Error getting user roles: %v\n", err)
		return
	}
	fmt.Printf("User roles: %v\n", roles)

	// Get user email
	email, err := userInfo.GetUserEmail(adminCtx)
	if err != nil {
		fmt.Printf("Error getting user email: %v\n", err)
		return
	}
	fmt.Printf("User email: %s\n", email)

	// Try with a non-existent user
	nonExistentCtx := AddUserIDToContext(ctx, "nonexistent")
	_, err = userInfo.GetUserRoles(nonExistentCtx)
	fmt.Printf("Error for non-existent user: %v\n", err)
}