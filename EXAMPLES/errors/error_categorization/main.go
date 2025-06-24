// Copyright (c) 2025 A Bit of Help, Inc.

// This example demonstrates how to categorize and handle different types of errors.
package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/abitofhelp/servicelib/errors"
)

// UserRepository simulates a user repository with potential errors
type UserRepository struct{}

// GetUserByID simulates getting a user by ID from a database
func (r *UserRepository) GetUserByID(id string) (string, error) {
	// Simulate different error scenarios based on the ID
	switch id {
	case "not-found":
		return "", errors.NewNotFoundError("User", id, nil)
	case "db-error":
		return "", errors.NewDatabaseError("database connection failed", "SELECT", "users", sql.ErrConnDone)
	case "validation-error":
		return "", errors.NewValidationError("invalid user ID format", "id", nil)
	case "business-rule":
		return "", errors.NewBusinessRuleError("user is inactive", "ActiveUserRule", nil)
	case "auth-error":
		return "", errors.NewAuthorizationError("insufficient permissions", "john.doe", "users", "read", nil)
	case "network-error":
		return "", errors.NewNetworkError("connection timeout", "db.example.com", "5432", nil)
	case "config-error":
		return "", errors.NewConfigurationError("invalid database URL", "DB_URL", "invalid://url", nil)
	case "external-service":
		return "", errors.NewExternalServiceError("payment service unavailable", "PaymentAPI", "/process", nil)
	default:
		return "User data for " + id, nil
	}
}

// UserCategoryService uses the repository and handles errors
type UserCategoryService struct {
	repo *UserRepository
}

// NewUserCategoryService creates a new UserCategoryService
func NewUserCategoryService() *UserCategoryService {
	return &UserCategoryService{
		repo: &UserRepository{},
	}
}

// GetUser gets a user by ID and handles different error categories
func (s *UserCategoryService) GetUser(id string) (string, error) {
	data, err := s.repo.GetUserByID(id)
	if err != nil {
		// Just pass the error up, we'll handle categorization at the API level
		return "", err
	}
	return data, nil
}

// UserAPI simulates an API that uses the service and categorizes errors
type UserAPI struct {
	service *UserCategoryService
}

// NewUserAPI creates a new UserAPI
func NewUserAPI() *UserAPI {
	return &UserAPI{
		service: NewUserCategoryService(),
	}
}

// HandleGetUser handles a user request and categorizes errors for appropriate responses
func (a *UserAPI) HandleGetUser(id string) (int, string) {
	data, err := a.service.GetUser(id)
	if err != nil {
		// Categorize the error and return appropriate HTTP status and message
		return a.categorizeAndHandleError(err)
	}
	return http.StatusOK, data
}

// categorizeAndHandleError categorizes errors and returns appropriate HTTP status and message
func (a *UserAPI) categorizeAndHandleError(err error) (int, string) {
	// Get HTTP status from error
	status := errors.GetHTTPStatus(err)
	if status != 0 {
		return status, err.Error()
	}

	// Categorize the error using type checking functions
	switch {
	case errors.IsNotFoundError(err):
		return http.StatusNotFound, fmt.Sprintf("Resource not found: %v", err)

	case errors.IsValidationError(err):
		return http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err)

	case errors.IsBusinessRuleError(err):
		return http.StatusUnprocessableEntity, fmt.Sprintf("Business rule violation: %v", err)

	case errors.IsAuthenticationError(err):
		return http.StatusUnauthorized, fmt.Sprintf("Authentication error: %v", err)

	case errors.IsAuthorizationError(err):
		return http.StatusForbidden, fmt.Sprintf("Authorization error: %v", err)

	case errors.IsDatabaseError(err):
		// Log the database error but don't expose details to the client
		fmt.Printf("Database error: %v\n", err)
		return http.StatusInternalServerError, "Internal server error"

	case errors.IsNetworkError(err):
		// Log the network error but don't expose details to the client
		fmt.Printf("Network error: %v\n", err)
		return http.StatusServiceUnavailable, "Service temporarily unavailable"

	case errors.IsConfigurationError(err):
		// Log the configuration error but don't expose details to the client
		fmt.Printf("Configuration error: %v\n", err)
		return http.StatusInternalServerError, "Internal server error"

	case errors.IsExternalServiceError(err):
		// Log the external service error but don't expose details to the client
		fmt.Printf("External service error: %v\n", err)
		return http.StatusBadGateway, "External service unavailable"

	case errors.IsTimeout(err):
		return http.StatusGatewayTimeout, "Request timed out"

	case errors.IsCancelled(err):
		return http.StatusRequestTimeout, "Request was cancelled"

	default:
		// For any other error, return a generic internal server error
		fmt.Printf("Uncategorized error: %v\n", err)
		return http.StatusInternalServerError, "Internal server error"
	}
}

func main() {
	fmt.Println("Error Categorization Example")
	fmt.Println("===========================")

	api := NewUserAPI()

	// Test different error scenarios
	testIDs := []string{
		"valid-id",
		"not-found",
		"db-error",
		"validation-error",
		"business-rule",
		"auth-error",
		"network-error",
		"config-error",
		"external-service",
	}

	for _, id := range testIDs {
		fmt.Printf("\nTesting with ID: %s\n", id)
		status, message := api.HandleGetUser(id)
		fmt.Printf("HTTP Status: %d\n", status)
		fmt.Printf("Response: %s\n", message)

		// For demonstration, also show how to check error types directly
		_, err := api.service.GetUser(id)
		if err != nil {
			fmt.Println("Error categorization:")

			if errors.IsNotFoundError(err) {
				fmt.Println("- Is a NotFoundError")
			}

			if errors.IsValidationError(err) {
				fmt.Println("- Is a ValidationError")
			}

			if errors.IsBusinessRuleError(err) {
				fmt.Println("- Is a BusinessRuleError")
			}

			if errors.IsAuthorizationError(err) {
				fmt.Println("- Is an AuthorizationError")
			}

			if errors.IsDatabaseError(err) {
				fmt.Println("- Is a DatabaseError")
			}

			if errors.IsNetworkError(err) {
				fmt.Println("- Is a NetworkError")
			}

			if errors.IsConfigurationError(err) {
				fmt.Println("- Is a ConfigurationError")
			}

			if errors.IsExternalServiceError(err) {
				fmt.Println("- Is an ExternalServiceError")
			}

			if errors.IsDomainError(err) {
				fmt.Println("- Is a DomainError")
			}

			if errors.IsInfrastructureError(err) {
				fmt.Println("- Is an InfrastructureError")
			}

			if errors.IsApplicationError(err) {
				fmt.Println("- Is an ApplicationError")
			}
		}
	}
}

// To run this example:
// go run examples/errors/error_categorization_example.go
