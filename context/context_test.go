// Copyright (c) 2025 A Bit of Help, Inc.

package context

import (
	"context"
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
