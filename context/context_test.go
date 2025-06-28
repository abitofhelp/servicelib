// Copyright (c) 2025 A Bit of Help, Inc.

package context

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/errors"
)

// TestWithTimeout tests the WithTimeout function
func TestWithTimeout(t *testing.T) {
	ctx := context.Background()
	timeout := 100 * time.Millisecond

	// Create a context with timeout
	timeoutCtx, cancel := WithTimeout(ctx, timeout, ContextOptions{})
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(timeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}

	// Wait for the context to timeout
	time.Sleep(timeout + 10*time.Millisecond)

	// Check that the context is done
	select {
	case <-timeoutCtx.Done():
		// Expected
	default:
		t.Error("Expected context to be done after timeout")
	}

	// Check the error
	if timeoutCtx.Err() != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", timeoutCtx.Err())
	}
}

// TestWithDefaultTimeout tests the WithDefaultTimeout function
func TestWithDefaultTimeout(t *testing.T) {
	ctx := context.Background()

	// Create a context with default timeout
	timeoutCtx, cancel := WithDefaultTimeout(ctx)
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(DefaultTimeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}
}

// TestWithShortTimeout tests the WithShortTimeout function
func TestWithShortTimeout(t *testing.T) {
	ctx := context.Background()

	// Create a context with short timeout
	timeoutCtx, cancel := WithShortTimeout(ctx)
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(ShortTimeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}
}

// TestWithLongTimeout tests the WithLongTimeout function
func TestWithLongTimeout(t *testing.T) {
	ctx := context.Background()

	// Create a context with long timeout
	timeoutCtx, cancel := WithLongTimeout(ctx)
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(LongTimeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}
}

// TestCheckContext tests the CheckContext function
func TestCheckContext(t *testing.T) {
	// Test with a background context
	t.Run("Background context", func(t *testing.T) {
		ctx := context.Background()
		err := CheckContext(ctx)
		if err != nil {
			t.Errorf("Expected no error for background context, got %v", err)
		}
	})

	// Test with a canceled context
	t.Run("Canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := CheckContext(ctx)
		if err == nil {
			t.Error("Expected error for canceled context, got nil")
		}

		if !errors.IsCancelled(err) {
			t.Errorf("Expected a cancelled error, got %v", err)
		}
	})

	// Test with a timed out context
	t.Run("Timed out context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait for the context to timeout
		time.Sleep(5 * time.Millisecond)

		err := CheckContext(ctx)
		if err == nil {
			t.Error("Expected error for timed out context, got nil")
		}

		if !errors.IsTimeout(err) {
			t.Errorf("Expected a timeout error, got %v", err)
		}
	})
}

// TestWithValues tests the WithValues function
func TestWithValues(t *testing.T) {
	ctx := context.Background()

	// Test with even number of arguments
	t.Run("Even number of arguments", func(t *testing.T) {
		key1 := "key1"
		value1 := "value1"
		key2 := "key2"
		value2 := "value2"

		newCtx := WithValues(ctx, key1, value1, key2, value2)

		// Check that the values are set
		if newCtx.Value(key1) != value1 {
			t.Errorf("Expected value %v for key %v, got %v", value1, key1, newCtx.Value(key1))
		}

		if newCtx.Value(key2) != value2 {
			t.Errorf("Expected value %v for key %v, got %v", value2, key2, newCtx.Value(key2))
		}
	})

	// Test with odd number of arguments
	t.Run("Odd number of arguments", func(t *testing.T) {
		key1 := "key1"
		value1 := "value1"
		key2 := "key2"

		newCtx := WithValues(ctx, key1, value1, key2) // Odd number of arguments

		// Check that the first key-value pair is set
		if newCtx.Value(key1) != value1 {
			t.Errorf("Expected value %v for key %v, got %v", value1, key1, newCtx.Value(key1))
		}

		// Check that the second key is not set
		if newCtx.Value(key2) != nil {
			t.Errorf("Expected nil value for key %v, got %v", key2, newCtx.Value(key2))
		}
	})
}

// TestBackground tests the Background function
func TestBackground(t *testing.T) {
	ctx := Background()

	// Check that the context is not nil
	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	// Check that the context is not done
	select {
	case <-ctx.Done():
		t.Error("Expected context to not be done")
	default:
		// Expected
	}

	// Check that the context has no deadline
	deadline, ok := ctx.Deadline()
	if ok {
		t.Errorf("Expected no deadline, got %v", deadline)
	}
}

// TestTODO tests the TODO function
func TestTODO(t *testing.T) {
	ctx := TODO()

	// Check that the context is not nil
	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	// Check that the context is not done
	select {
	case <-ctx.Done():
		t.Error("Expected context to not be done")
	default:
		// Expected
	}

	// Check that the context has no deadline
	deadline, ok := ctx.Deadline()
	if ok {
		t.Errorf("Expected no deadline, got %v", deadline)
	}
}

// TestRequestID tests the WithRequestID and GetRequestID functions
func TestRequestID(t *testing.T) {
	ctx := context.Background()

	// Check that the request ID is initially empty
	requestID := GetRequestID(ctx)
	if requestID != "" {
		t.Errorf("Expected empty request ID, got %v", requestID)
	}

	// Add a request ID
	ctx = WithRequestID(ctx)

	// Check that the request ID is set
	requestID = GetRequestID(ctx)
	if requestID == "" {
		t.Error("Expected non-empty request ID")
	}
}

// TestTraceID tests the WithTraceID and GetTraceID functions
func TestTraceID(t *testing.T) {
	ctx := context.Background()

	// Check that the trace ID is initially empty
	traceID := GetTraceID(ctx)
	if traceID != "" {
		t.Errorf("Expected empty trace ID, got %v", traceID)
	}

	// Add a trace ID
	ctx = WithTraceID(ctx)

	// Check that the trace ID is set
	traceID = GetTraceID(ctx)
	if traceID == "" {
		t.Error("Expected non-empty trace ID")
	}
}

// TestUserID tests the WithUserID and GetUserID functions
func TestUserID(t *testing.T) {
	ctx := context.Background()

	// Check that the user ID is initially empty
	userID := GetUserID(ctx)
	if userID != "" {
		t.Errorf("Expected empty user ID, got %v", userID)
	}

	// Add a user ID
	expectedUserID := "user123"
	ctx = WithUserID(ctx, expectedUserID)

	// Check that the user ID is set
	userID = GetUserID(ctx)
	if userID != expectedUserID {
		t.Errorf("Expected user ID %v, got %v", expectedUserID, userID)
	}
}

// TestTenantID tests the WithTenantID and GetTenantID functions
func TestTenantID(t *testing.T) {
	ctx := context.Background()

	// Check that the tenant ID is initially empty
	tenantID := GetTenantID(ctx)
	if tenantID != "" {
		t.Errorf("Expected empty tenant ID, got %v", tenantID)
	}

	// Add a tenant ID
	expectedTenantID := "tenant456"
	ctx = WithTenantID(ctx, expectedTenantID)

	// Check that the tenant ID is set
	tenantID = GetTenantID(ctx)
	if tenantID != expectedTenantID {
		t.Errorf("Expected tenant ID %v, got %v", expectedTenantID, tenantID)
	}
}

// TestNewContext tests the NewContext function
func TestNewContext(t *testing.T) {
	// Test with nil parent
	t.Run("Nil parent", func(t *testing.T) {
		opts := ContextOptions{
			Timeout:       100 * time.Millisecond,
			RequestID:     "req123",
			TraceID:       "trace456",
			UserID:        "user789",
			TenantID:      "tenant012",
			Operation:     "test-operation",
			CorrelationID: "corr345",
			ServiceName:   "test-service",
			Parent:        nil,
		}

		ctx, cancel := NewContext(opts)
		defer cancel()

		// Check that the context is not nil
		if ctx == nil {
			t.Error("Expected non-nil context")
		}

		// Check that the deadline is set
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Expected deadline to be set")
		}

		// Check that the deadline is approximately correct
		expectedDeadline := time.Now().Add(opts.Timeout)
		if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
			t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
		}

		// Check that the values are set
		if GetRequestID(ctx) != opts.RequestID {
			t.Errorf("Expected request ID %v, got %v", opts.RequestID, GetRequestID(ctx))
		}

		if GetTraceID(ctx) != opts.TraceID {
			t.Errorf("Expected trace ID %v, got %v", opts.TraceID, GetTraceID(ctx))
		}

		if GetUserID(ctx) != opts.UserID {
			t.Errorf("Expected user ID %v, got %v", opts.UserID, GetUserID(ctx))
		}

		if GetTenantID(ctx) != opts.TenantID {
			t.Errorf("Expected tenant ID %v, got %v", opts.TenantID, GetTenantID(ctx))
		}

		if GetOperation(ctx) != opts.Operation {
			t.Errorf("Expected operation %v, got %v", opts.Operation, GetOperation(ctx))
		}

		if GetCorrelationID(ctx) != opts.CorrelationID {
			t.Errorf("Expected correlation ID %v, got %v", opts.CorrelationID, GetCorrelationID(ctx))
		}

		if GetServiceName(ctx) != opts.ServiceName {
			t.Errorf("Expected service name %v, got %v", opts.ServiceName, GetServiceName(ctx))
		}
	})

	// Test with parent context
	t.Run("With parent", func(t *testing.T) {
		parent := context.Background()
		opts := ContextOptions{
			Timeout:       100 * time.Millisecond,
			RequestID:     "req123",
			TraceID:       "trace456",
			UserID:        "user789",
			TenantID:      "tenant012",
			Operation:     "test-operation",
			CorrelationID: "corr345",
			ServiceName:   "test-service",
			Parent:        parent,
		}

		ctx, cancel := NewContext(opts)
		defer cancel()

		// Check that the context is not nil
		if ctx == nil {
			t.Error("Expected non-nil context")
		}

		// Check that the deadline is set
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Expected deadline to be set")
		}

		// Check that the deadline is approximately correct
		expectedDeadline := time.Now().Add(opts.Timeout)
		if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
			t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
		}
	})

	// Test without timeout
	t.Run("Without timeout", func(t *testing.T) {
		opts := ContextOptions{
			RequestID:     "req123",
			TraceID:       "trace456",
			UserID:        "user789",
			TenantID:      "tenant012",
			Operation:     "test-operation",
			CorrelationID: "corr345",
			ServiceName:   "test-service",
			Parent:        context.Background(),
		}

		ctx, cancel := NewContext(opts)
		defer cancel()

		// Check that the context is not nil
		if ctx == nil {
			t.Error("Expected non-nil context")
		}

		// Check that the deadline is not set
		_, ok := ctx.Deadline()
		if ok {
			t.Error("Expected no deadline to be set")
		}
	})
}

// TestWithDatabaseTimeout tests the WithDatabaseTimeout function
func TestWithDatabaseTimeout(t *testing.T) {
	ctx := context.Background()

	// Create a context with database timeout
	timeoutCtx, cancel := WithDatabaseTimeout(ctx)
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(DatabaseTimeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}

	// Check that the operation is set
	operation := GetOperation(timeoutCtx)
	if operation != "database" {
		t.Errorf("Expected operation 'database', got %v", operation)
	}
}

// TestWithNetworkTimeout tests the WithNetworkTimeout function
func TestWithNetworkTimeout(t *testing.T) {
	ctx := context.Background()

	// Create a context with network timeout
	timeoutCtx, cancel := WithNetworkTimeout(ctx)
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(NetworkTimeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}

	// Check that the operation is set
	operation := GetOperation(timeoutCtx)
	if operation != "network" {
		t.Errorf("Expected operation 'network', got %v", operation)
	}
}

// TestWithExternalServiceTimeout tests the WithExternalServiceTimeout function
func TestWithExternalServiceTimeout(t *testing.T) {
	ctx := context.Background()
	serviceName := "test-external-service"

	// Create a context with external service timeout
	timeoutCtx, cancel := WithExternalServiceTimeout(ctx, serviceName)
	defer cancel()

	// Check that the deadline is set
	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Expected deadline to be set")
	}

	// Check that the deadline is approximately correct
	expectedDeadline := time.Now().Add(ExternalServiceTimeout)
	if deadline.Sub(expectedDeadline) > 10*time.Millisecond {
		t.Errorf("Deadline not set correctly, got %v, expected approximately %v", deadline, expectedDeadline)
	}

	// Check that the operation is set
	operation := GetOperation(timeoutCtx)
	if operation != "external_service" {
		t.Errorf("Expected operation 'external_service', got %v", operation)
	}

	// Check that the service name is set
	svcName := GetServiceName(timeoutCtx)
	if svcName != serviceName {
		t.Errorf("Expected service name %v, got %v", serviceName, svcName)
	}
}

// TestWithOperation tests the WithOperation and GetOperation functions
func TestWithOperation(t *testing.T) {
	ctx := context.Background()

	// Check that the operation is initially empty
	operation := GetOperation(ctx)
	if operation != "" {
		t.Errorf("Expected empty operation, got %v", operation)
	}

	// Add an operation
	expectedOperation := "test-operation"
	ctx = WithOperation(ctx, expectedOperation)

	// Check that the operation is set
	operation = GetOperation(ctx)
	if operation != expectedOperation {
		t.Errorf("Expected operation %v, got %v", expectedOperation, operation)
	}
}

// TestWithCorrelationID tests the WithCorrelationID and GetCorrelationID functions
func TestWithCorrelationID(t *testing.T) {
	ctx := context.Background()

	// Check that the correlation ID is initially empty
	correlationID := GetCorrelationID(ctx)
	if correlationID != "" {
		t.Errorf("Expected empty correlation ID, got %v", correlationID)
	}

	// Add a correlation ID
	expectedCorrelationID := "corr123"
	ctx = WithCorrelationID(ctx, expectedCorrelationID)

	// Check that the correlation ID is set
	correlationID = GetCorrelationID(ctx)
	if correlationID != expectedCorrelationID {
		t.Errorf("Expected correlation ID %v, got %v", expectedCorrelationID, correlationID)
	}
}

// TestWithServiceName tests the WithServiceName and GetServiceName functions
func TestWithServiceName(t *testing.T) {
	ctx := context.Background()

	// Check that the service name is initially empty
	serviceName := GetServiceName(ctx)
	if serviceName != "" {
		t.Errorf("Expected empty service name, got %v", serviceName)
	}

	// Add a service name
	expectedServiceName := "test-service"
	ctx = WithServiceName(ctx, expectedServiceName)

	// Check that the service name is set
	serviceName = GetServiceName(ctx)
	if serviceName != expectedServiceName {
		t.Errorf("Expected service name %v, got %v", expectedServiceName, serviceName)
	}
}

// TestMustCheck tests the MustCheck function
func TestMustCheck(t *testing.T) {
	// Test with a background context
	t.Run("Background context", func(t *testing.T) {
		ctx := context.Background()

		// This should not panic
		MustCheck(ctx)
	})

	// Test with a canceled context
	t.Run("Canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// This should panic
		defer func() {
			r := recover()
			if r == nil {
				t.Error("Expected panic for canceled context")
			}
		}()

		MustCheck(ctx)
	})
}

// TestContextInfo tests the ContextInfo function
func TestContextInfo(t *testing.T) {
	// Test with empty context
	t.Run("Empty context", func(t *testing.T) {
		ctx := context.Background()

		// Add request ID and trace ID
		ctx = WithRequestID(ctx)
		ctx = WithTraceID(ctx)

		info := ContextInfo(ctx)

		// Check that the info contains the request ID and trace ID
		if info == "" {
			t.Error("Expected non-empty context info")
		}

		if !contains(info, "RequestID:") || !contains(info, "TraceID:") {
			t.Errorf("Expected info to contain RequestID and TraceID, got %v", info)
		}
	})

	// Test with all values
	t.Run("All values", func(t *testing.T) {
		ctx := context.Background()

		// Add all values
		ctx = WithRequestID(ctx)
		ctx = WithTraceID(ctx)
		ctx = WithUserID(ctx, "user123")
		ctx = WithTenantID(ctx, "tenant456")
		ctx = WithOperation(ctx, "test-operation")
		ctx = WithCorrelationID(ctx, "corr789")
		ctx = WithServiceName(ctx, "test-service")

		// Add deadline
		ctx, cancel := context.WithTimeout(ctx, 1*time.Hour)
		defer cancel()

		info := ContextInfo(ctx)

		// Check that the info contains all values
		if !contains(info, "RequestID:") || !contains(info, "TraceID:") ||
			!contains(info, "UserID: user123") || !contains(info, "TenantID: tenant456") ||
			!contains(info, "Operation: test-operation") || !contains(info, "CorrelationID: corr789") ||
			!contains(info, "ServiceName: test-service") || !contains(info, "Deadline:") {
			t.Errorf("Expected info to contain all values, got %v", info)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && strings.Contains(s, substr)
}
