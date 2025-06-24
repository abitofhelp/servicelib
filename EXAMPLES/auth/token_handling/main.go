// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating the concept of token handling in authentication
package main

import (
	"fmt"
	"time"
)

// This is a simplified example to demonstrate the concept of token handling
// In a real application, you would use the auth module's methods

// MockClaims represents the data stored in a JWT token
type MockClaims struct {
	Subject   string
	Roles     []string
	Issuer    string
	ExpiresAt time.Time
}

// MockAuth simulates an auth service
type MockAuth struct {
	secretKey string
	issuer    string
}

// NewMockAuth creates a new mock auth service
func NewMockAuth(secretKey, issuer string) *MockAuth {
	return &MockAuth{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

// GenerateToken creates a mock token
func (a *MockAuth) GenerateToken(userID string, roles []string) string {
	// In a real implementation, this would create a JWT token
	return fmt.Sprintf("mock.token.%s.%s", userID, a.secretKey)
}

// ValidateToken validates a mock token and returns claims
func (a *MockAuth) ValidateToken(token string) (*MockClaims, error) {
	// In a real implementation, this would validate the JWT token
	// and extract the claims

	// For this example, we'll just create mock claims
	return &MockClaims{
		Subject:   "user123",
		Roles:     []string{"admin", "editor"},
		Issuer:    a.issuer,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func main() {
	// Create a mock auth service
	auth := NewMockAuth("your-secret-key", "example-service")

	// Generate a token for a user with roles
	userID := "user123"
	roles := []string{"admin", "editor"}
	token := auth.GenerateToken(userID, roles)

	fmt.Printf("Generated token: %s\n", token)

	// Validate the token
	claims, err := auth.ValidateToken(token)
	if err != nil {
		fmt.Printf("Failed to validate token: %v\n", err)
		return
	}

	// Extract information from the claims
	fmt.Printf("Token validated successfully\n")
	fmt.Printf("User ID: %s\n", claims.Subject)
	fmt.Printf("Roles: %v\n", claims.Roles)
	fmt.Printf("Issuer: %s\n", claims.Issuer)
	fmt.Printf("Expiration: %v\n", claims.ExpiresAt)
}