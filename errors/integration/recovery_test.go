// Copyright (c) 2025 A Bit of Help, Inc.

package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/errors/recovery"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRecoveryIntegration(t *testing.T) {
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	ctx := context.Background()
	handler := recovery.NewRecoveryHandler(logger, 10)
	cb := recovery.NewCircuitBreaker(2, time.Second)

	// Test successful operation with recovery and circuit breaker
	err := handler.WithRecovery(ctx, "test-operation", func() error {
		return cb.ExecuteWithName(ctx, "test-operation", func() error {
			return nil
		})
	})
	require.NoError(t, err)

	// Test operation with error
	expectedErr := errors.New("operation failed")
	err = handler.WithRecovery(ctx, "test-operation", func() error {
		return cb.ExecuteWithName(ctx, "test-operation", func() error {
			return expectedErr
		})
	})
	assert.ErrorIs(t, err, expectedErr)

	// Test circuit breaker opening
	for i := 0; i < 3; i++ {
		err = handler.WithRecovery(ctx, "test-operation", func() error {
			return cb.ExecuteWithName(ctx, "test-operation", func() error {
				return errors.New("failure")
			})
		})
		assert.Error(t, err)
	}

	// Verify circuit is open
	err = handler.WithRecovery(ctx, "test-operation", func() error {
		return cb.ExecuteWithName(ctx, "test-operation", func() error {
			return nil // This should not be executed
		})
	})
	assert.ErrorIs(t, err, recovery.ErrCircuitBreakerOpen)

	// Wait for circuit breaker to reset
	time.Sleep(time.Second)

	// Test successful operation after reset
	err = handler.WithRecovery(ctx, "test-operation", func() error {
		return cb.ExecuteWithName(ctx, "test-operation", func() error {
			return nil
		})
	})
	assert.NoError(t, err)
}

func TestPanicRecoveryWithCircuitBreaker(t *testing.T) {
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	ctx := context.Background()
	handler := recovery.NewRecoveryHandler(logger, 10)
	cb := recovery.NewCircuitBreaker(2, time.Second)

	// Test panic recovery
	err := handler.WithRecovery(ctx, "panic-operation", func() error {
		return cb.ExecuteWithName(ctx, "panic-operation", func() error {
			panic("something went wrong")
		})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic in panic-operation")

	// Second panic should count as a failure
	err = handler.WithRecovery(ctx, "panic-operation", func() error {
		return cb.ExecuteWithName(ctx, "panic-operation", func() error {
			panic("another panic")
		})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic in panic-operation")

	// Verify circuit is open after panics
	err = handler.WithRecovery(ctx, "panic-operation", func() error {
		return cb.ExecuteWithName(ctx, "panic-operation", func() error {
			return nil // This should not be executed
		})
	})
	assert.ErrorIs(t, err, recovery.ErrCircuitBreakerOpen)

	// Wait for circuit breaker to reset
	time.Sleep(time.Second)

	// Test successful operation after reset
	err = handler.WithRecovery(ctx, "panic-operation", func() error {
		return cb.ExecuteWithName(ctx, "panic-operation", func() error {
			return nil
		})
	})
	assert.NoError(t, err)
}

func TestContextCancellationWithCircuitBreaker(t *testing.T) {
	logger := logging.NewContextLogger(zaptest.NewLogger(t))
	handler := recovery.NewRecoveryHandler(logger, 10)
	cb := recovery.NewCircuitBreaker(2, time.Second)

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	err := handler.WithRecovery(ctx, "cancel-operation", func() error {
		return cb.ExecuteWithName(ctx, "cancel-operation", func() error {
			return nil // This should not be executed
		})
	})
	assert.ErrorIs(t, err, recovery.ErrContextCanceled)

	// Test context deadline exceeded
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // Ensure deadline is exceeded

	err = handler.WithRecovery(ctx, "timeout-operation", func() error {
		return cb.ExecuteWithName(ctx, "timeout-operation", func() error {
			return nil // This should not be executed
		})
	})
	assert.ErrorIs(t, err, recovery.ErrContextDeadlineExceeded)
}
