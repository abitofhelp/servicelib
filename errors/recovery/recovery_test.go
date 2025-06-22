// Copyright (c) 2025 A Bit of Help, Inc.

package recovery

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewRecoveryHandler(t *testing.T) {
	// Arrange
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	errorsPerSecond := 10

	// Act
	handler := NewRecoveryHandler(logger, errorsPerSecond)

	// Assert
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
	assert.NotNil(t, handler.rateLimiter)
}

func TestRecoveryHandler_WithRecovery_Success(t *testing.T) {
	// Arrange
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	handler := NewRecoveryHandler(logger, 10)
	ctx := context.Background()
	operationCalled := false

	// Act
	err := handler.WithRecovery(ctx, "test-operation", func() error {
		operationCalled = true
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.True(t, operationCalled)
}

func TestRecoveryHandler_WithRecovery_Error(t *testing.T) {
	// Arrange
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	handler := NewRecoveryHandler(logger, 10)
	ctx := context.Background()
	expectedErr := errors.New("operation failed")

	// Act
	err := handler.WithRecovery(ctx, "test-operation", func() error {
		return expectedErr
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestRecoveryHandler_WithRecovery_Panic(t *testing.T) {
	// Arrange
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	handler := NewRecoveryHandler(logger, 10)
	ctx := context.Background()
	panicMsg := "panic occurred"

	// Act
	err := handler.WithRecovery(ctx, "test-operation", func() error {
		panic(panicMsg)
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), panicMsg)
	assert.Contains(t, err.Error(), "test-operation")
}

func TestNewCircuitBreaker(t *testing.T) {
	// Arrange
	failureThreshold := 3
	cooldownPeriod := 5 * time.Second

	// Act
	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod)

	// Assert
	assert.NotNil(t, cb)
	assert.Equal(t, failureThreshold, cb.failureThreshold)
	assert.Equal(t, cooldownPeriod, cb.cooldownPeriod)
	assert.Equal(t, 0, cb.failures)
	assert.True(t, cb.openUntil.IsZero())
}

func TestCircuitBreaker_Execute(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	ctx := context.Background()
	operationCalled := false

	// Act
	err := cb.Execute(ctx, func() error {
		operationCalled = true
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.True(t, operationCalled)
}

func TestCircuitBreaker_ExecuteWithName_Success(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	ctx := context.Background()
	operationCalled := false

	// Act
	err := cb.ExecuteWithName(ctx, "test-operation", func() error {
		operationCalled = true
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.True(t, operationCalled)
}

func TestCircuitBreaker_ExecuteWithName_Error(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	ctx := context.Background()
	expectedErr := errors.New("operation failed")

	// Act
	err := cb.ExecuteWithName(ctx, "test-operation", func() error {
		return expectedErr
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 1, cb.failures)
}

func TestCircuitBreaker_ExecuteWithName_Panic(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	ctx := context.Background()
	panicMsg := "panic occurred"

	// Act
	err := cb.ExecuteWithName(ctx, "test-operation", func() error {
		panic(panicMsg)
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), panicMsg)
	assert.Contains(t, err.Error(), "test-operation")
	assert.Equal(t, 1, cb.failures)
}

func TestCircuitBreaker_ExecuteWithName_ContextCanceled(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context

	// Act
	err := cb.ExecuteWithName(ctx, "test-operation", func() error {
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrContextCanceled, err)
}

func TestCircuitBreaker_ExecuteWithName_ContextDeadlineExceeded(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond) // Ensure deadline is exceeded

	// Act
	err := cb.ExecuteWithName(ctx, "test-operation", func() error {
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrContextDeadlineExceeded, err)
}

func TestCircuitBreaker_ExecuteWithName_CircuitOpen(t *testing.T) {
	// Arrange
	failureThreshold := 3
	cooldownPeriod := 100 * time.Millisecond
	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod)
	ctx := context.Background()

	// Trigger failures to open the circuit
	for i := 0; i < failureThreshold; i++ {
		_ = cb.ExecuteWithName(ctx, "test-operation", func() error {
			return errors.New("operation failed")
		})
	}

	// Act
	err := cb.ExecuteWithName(ctx, "test-operation", func() error {
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrCircuitBreakerOpen, err)

	// Wait for cooldown period to pass
	time.Sleep(cooldownPeriod + 10*time.Millisecond)

	// Try again after cooldown
	err = cb.ExecuteWithName(ctx, "test-operation", func() error {
		return nil
	})

	// Assert circuit is closed
	assert.NoError(t, err)
}

func TestCircuitBreaker_recordSuccess(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 5*time.Second)
	cb.failures = 2
	cb.openUntil = time.Now().Add(5 * time.Second)

	// Act
	cb.recordSuccess()

	// Assert
	assert.Equal(t, 0, cb.failures)
	assert.True(t, cb.openUntil.IsZero())
}

func TestCircuitBreaker_recordFailure(t *testing.T) {
	// Arrange
	failureThreshold := 3
	cooldownPeriod := 5 * time.Second
	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod)

	// Act - Record failures below threshold
	for i := 0; i < failureThreshold-1; i++ {
		cb.recordFailure()
	}

	// Assert - Circuit should still be closed
	assert.Equal(t, failureThreshold-1, cb.failures)
	assert.True(t, cb.openUntil.IsZero())

	// Act - Record one more failure to reach threshold
	cb.recordFailure()

	// Assert - Circuit should be open
	assert.Equal(t, failureThreshold, cb.failures)
	assert.False(t, cb.openUntil.IsZero())
	assert.True(t, cb.openUntil.After(time.Now()))
}

func TestCircuitBreaker_isOpen(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	// Test when circuit is closed
	assert.False(t, cb.isOpen())

	// Open the circuit
	cb.failures = 3
	cb.openUntil = time.Now().Add(100 * time.Millisecond)

	// Test when circuit is open
	assert.True(t, cb.isOpen())

	// Wait for cooldown period to pass
	time.Sleep(150 * time.Millisecond)

	// Test after cooldown period
	assert.False(t, cb.isOpen())
	assert.Equal(t, 0, cb.failures)
	assert.True(t, cb.openUntil.IsZero())
}
