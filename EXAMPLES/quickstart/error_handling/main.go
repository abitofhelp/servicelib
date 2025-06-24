// Copyright (c) 2025 A Bit of Help, Inc.

// This example demonstrates how to use the error handling system in servicelib.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/abitofhelp/servicelib/errors"
	errhttp "github.com/abitofhelp/servicelib/errors/http"
	errlog "github.com/abitofhelp/servicelib/errors/log"
	errmetrics "github.com/abitofhelp/servicelib/errors/metrics"
	errtrace "github.com/abitofhelp/servicelib/errors/trace"
	"github.com/abitofhelp/servicelib/logging"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// UserService is a simple service that demonstrates error handling
type UserService struct {
	logger *logging.ContextLogger
	meter  metric.Meter
}

// NewUserService creates a new UserService
func NewUserService(logger *logging.ContextLogger, meter metric.Meter) *UserService {
	return &UserService{
		logger: logger,
		meter:  meter,
	}
}

// User represents a user in the system
type User struct {
	ID    string
	Name  string
	Email string
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
	// Simulate a database query
	if id == "" {
		// Create a validation error
		err := errors.NewValidationError("User ID is required", "id", nil)

		// Log the error
		errlog.LogError(ctx, s.logger, err)

		// Record error metric
		errmetrics.RecordError(ctx, err)

		// Add error to trace
		errtrace.AddErrorToSpan(ctx, err)

		return nil, err
	}

	if id == "not-found" {
		// Create a not found error
		err := errors.NewNotFoundError("User", id, nil)

		// Log the error
		errlog.LogError(ctx, s.logger, err)

		// Record error metric
		errmetrics.RecordError(ctx, err)

		// Add error to trace
		errtrace.AddErrorToSpan(ctx, err)

		return nil, err
	}

	if id == "db-error" {
		// Create a database error
		dbErr := sql.ErrNoRows
		err := errors.NewDatabaseError("Failed to query user", "SELECT", "users", dbErr)

		// Log the error
		errlog.LogError(ctx, s.logger, err)

		// Record error metric
		errmetrics.RecordError(ctx, err)

		// Add error to trace
		errtrace.AddErrorToSpan(ctx, err)

		return nil, err
	}

	// Return a user
	return &User{
		ID:    id,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}, nil
}

// ValidateUser validates a user
func (s *UserService) ValidateUser(ctx context.Context, user *User) error {
	// Create a validation errors collection
	validationErrors := errors.NewValidationErrors("User validation failed")

	// Validate name
	if user.Name == "" {
		validationErrors.AddError(errors.NewValidationError("Name is required", "name", nil))
	}

	// Validate email
	if user.Email == "" {
		validationErrors.AddError(errors.NewValidationError("Email is required", "email", nil))
	} else if !isValidEmail(user.Email) {
		validationErrors.AddError(errors.NewValidationError("Email is invalid", "email", nil))
	}

	// If there are validation errors, return them
	if validationErrors.HasErrors() {
		// Log the error
		errlog.LogError(ctx, s.logger, validationErrors)

		// Record error metric
		errmetrics.RecordError(ctx, validationErrors)

		// Add error to trace
		errtrace.AddErrorToSpan(ctx, validationErrors)

		return validationErrors
	}

	return nil
}

// isValidEmail is a simple email validation function
func isValidEmail(email string) bool {
	// This is a very simple validation for demonstration purposes
	return len(email) > 0 && email[len(email)-4:] == ".com"
}

// UserHandler is a simple HTTP handler that demonstrates error handling
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetUser handles GET requests for users
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	id := r.URL.Query().Get("id")

	// Get user from service
	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		// Write error response
		errhttp.WriteError(w, err)
		return
	}

	// Write success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id":"%s","name":"%s","email":"%s"}`, user.ID, user.Name, user.Email)
}

// RunErrorHandlingExample demonstrates how to use the error handling system
func RunErrorHandlingExample() {
	// Create a logger
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a meter (this would normally come from your telemetry setup)
	var meter metric.Meter

	// Create the user service
	service := NewUserService(logger, meter)

	// Create the user handler
	handler := NewUserHandler(service)

	// Set up HTTP routes
	http.HandleFunc("/user", handler.GetUser)

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Try these URLs:")
	fmt.Println("- http://localhost:8080/user?id=123 (success)")
	fmt.Println("- http://localhost:8080/user?id= (validation error)")
	fmt.Println("- http://localhost:8080/user?id=not-found (not found error)")
	fmt.Println("- http://localhost:8080/user?id=db-error (database error)")
	http.ListenAndServe(":8080", nil)
}
