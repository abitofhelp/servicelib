// Copyright (c) 2025 A Bit of Help, Inc.

package infra

import (
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

func TestNewInfrastructureError(t *testing.T) {
	// Test creating a new InfrastructureError
	err := NewInfrastructureError(core.DatabaseErrorCode, "Database connection failed", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.DatabaseErrorCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Database connection failed", err.BaseError.GetMessage())
	
	// Check that IsInfrastructureError returns true
	assert.True(t, err.IsInfrastructureError())
}

func TestNewDatabaseError(t *testing.T) {
	// Test creating a new DatabaseError
	err := NewDatabaseError("Failed to query database", "SELECT", "users", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.DatabaseErrorCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Failed to query database", err.BaseError.GetMessage())
	
	// Check that the operation and table are set correctly
	assert.Equal(t, "SELECT", err.Operation)
	assert.Equal(t, "users", err.Table)
	
	// Check that IsDatabaseError returns true
	assert.True(t, err.IsDatabaseError())
	
	// Check that IsInfrastructureError returns true (inheritance)
	assert.True(t, err.IsInfrastructureError())
}

func TestNewNetworkError(t *testing.T) {
	// Test creating a new NetworkError
	err := NewNetworkError("Failed to connect to server", "example.com", "8080", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.NetworkErrorCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Failed to connect to server", err.BaseError.GetMessage())
	
	// Check that the host and port are set correctly
	assert.Equal(t, "example.com", err.Host)
	assert.Equal(t, "8080", err.Port)
	
	// Check that IsNetworkError returns true
	assert.True(t, err.IsNetworkError())
	
	// Check that IsInfrastructureError returns true (inheritance)
	assert.True(t, err.IsInfrastructureError())
}

func TestNewExternalServiceError(t *testing.T) {
	// Test creating a new ExternalServiceError
	err := NewExternalServiceError("Failed to call external API", "PaymentService", "/api/payments", nil)
	
	// Check that the error is not nil
	assert.NotNil(t, err)
	
	// Check that the error code is set correctly
	assert.Equal(t, core.ExternalServiceErrorCode, err.BaseError.GetCode())
	
	// Check that the message is set correctly
	assert.Equal(t, "Failed to call external API", err.BaseError.GetMessage())
	
	// Check that the service name and endpoint are set correctly
	assert.Equal(t, "PaymentService", err.ServiceName)
	assert.Equal(t, "/api/payments", err.Endpoint)
	
	// Check that IsExternalServiceError returns true
	assert.True(t, err.IsExternalServiceError())
	
	// Check that IsInfrastructureError returns true (inheritance)
	assert.True(t, err.IsInfrastructureError())
}

func TestInfraErrorsAs(t *testing.T) {
	// Create errors of different types
	infraErr := NewInfrastructureError(core.DatabaseErrorCode, "Infrastructure error", nil)
	dbErr := NewDatabaseError("Failed to query database", "SELECT", "users", nil)
	networkErr := NewNetworkError("Failed to connect to server", "example.com", "8080", nil)
	externalErr := NewExternalServiceError("Failed to call external API", "PaymentService", "/api/payments", nil)
	
	// Test errors.As with InfrastructureError
	var ie *InfrastructureError
	assert.True(t, errors.As(infraErr, &ie))
	
	// Test errors.As with DatabaseError
	var dbe *DatabaseError
	assert.True(t, errors.As(dbErr, &dbe))
	
	// Test errors.As with NetworkError
	var ne *NetworkError
	assert.True(t, errors.As(networkErr, &ne))
	
	// Test errors.As with ExternalServiceError
	var ese *ExternalServiceError
	assert.True(t, errors.As(externalErr, &ese))
	
	// Test inheritance (DatabaseError is also an InfrastructureError)
	assert.True(t, errors.As(dbErr, &ie))
}

func TestInfraErrorWithCause(t *testing.T) {
	// Create a cause error
	cause := errors.New("original error")
	
	// Create an InfrastructureError with the cause
	err := NewInfrastructureError(core.DatabaseErrorCode, "Database connection failed", cause)
	
	// Check that the cause is set correctly
	assert.Equal(t, cause, err.BaseError.GetCause())
	
	// Check that the error message includes the cause
	assert.Contains(t, err.BaseError.Error(), "original error")
}