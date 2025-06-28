// Copyright (c) 2025 A Bit of Help, Inc.

package circuit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewCircuitBreaker(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		options  Options
		expected bool
	}{
		{
			name: "Enabled circuit breaker",
			config: DefaultConfig().
				WithEnabled(true),
			options:  DefaultOptions(),
			expected: true,
		},
		{
			name: "Disabled circuit breaker",
			config: DefaultConfig().
				WithEnabled(false),
			options:  DefaultOptions(),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cb := NewCircuitBreaker(tc.config, tc.options)
			if tc.expected {
				assert.NotNil(t, cb)
			} else {
				assert.Nil(t, cb)
			}
		})
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	// Create a circuit breaker
	cfg := DefaultConfig().
		WithEnabled(true).
		WithErrorThreshold(0.5).
		WithVolumeThreshold(3)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cb := NewCircuitBreaker(cfg, options)
	assert.NotNil(t, cb)

	// Execute a successful function
	successFunc := func(ctx context.Context) (bool, error) {
		return true, nil
	}

	result, err := Execute(context.Background(), cb, "TestSuccess", successFunc)
	assert.NoError(t, err)
	assert.True(t, result)

	// Circuit should still be closed
	assert.Equal(t, Closed, cb.GetState())
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	// Create a circuit breaker with a low volume threshold for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithErrorThreshold(0.5).
		WithVolumeThreshold(2)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cb := NewCircuitBreaker(cfg, options)
	assert.NotNil(t, cb)

	// Execute a failing function
	failFunc := func(ctx context.Context) (bool, error) {
		return false, errors.New("test error")
	}

	// First execution should fail but not trip the circuit
	result, err := Execute(context.Background(), cb, "TestFailure", failFunc)
	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, Closed, cb.GetState())

	// Second execution should fail and trip the circuit
	result, err = Execute(context.Background(), cb, "TestFailure", failFunc)
	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, Open, cb.GetState())

	// Third execution should fail immediately with circuit open error
	successFunc := func(ctx context.Context) (bool, error) {
		t.Fatal("This function should not be called when circuit is open")
		return true, nil
	}

	result, err = Execute(context.Background(), cb, "TestFailure", successFunc)
	assert.Error(t, err)
	assert.False(t, result)
}

func TestCircuitBreaker_ExecuteWithFallback(t *testing.T) {
	// Create a circuit breaker
	cfg := DefaultConfig().
		WithEnabled(true).
		WithErrorThreshold(0.5).
		WithVolumeThreshold(2)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cb := NewCircuitBreaker(cfg, options)
	assert.NotNil(t, cb)

	// Define a failing function
	failFunc := func(ctx context.Context) (string, error) {
		return "", errors.New("test error")
	}

	// Define a fallback function
	fallbackFunc := func(ctx context.Context, err error) (string, error) {
		return "fallback", nil
	}

	// Execute with fallback should use the fallback
	result, err := ExecuteWithFallback(context.Background(), cb, "TestFallback", failFunc, fallbackFunc)
	assert.NoError(t, err)
	assert.Equal(t, "fallback", result)

	// Second execution should fail and trip the circuit
	result, err = ExecuteWithFallback(context.Background(), cb, "TestFallback", failFunc, fallbackFunc)
	assert.NoError(t, err)
	assert.Equal(t, "fallback", result)
	assert.Equal(t, Open, cb.GetState())

	// Third execution should use fallback immediately without calling the function
	neverCalledFunc := func(ctx context.Context) (string, error) {
		t.Fatal("This function should not be called when circuit is open")
		return "", nil
	}

	result, err = ExecuteWithFallback(context.Background(), cb, "TestFallback", neverCalledFunc, fallbackFunc)
	assert.NoError(t, err)
	assert.Equal(t, "fallback", result)
}

func TestCircuitBreaker_HalfOpen(t *testing.T) {
	// Create a circuit breaker with a short sleep window for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithErrorThreshold(0.5).
		WithVolumeThreshold(2).
		WithSleepWindow(10 * time.Millisecond)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cb := NewCircuitBreaker(cfg, options)
	assert.NotNil(t, cb)

	// Trip the circuit
	failFunc := func(ctx context.Context) (bool, error) {
		return false, errors.New("test error")
	}

	// First execution
	Execute(context.Background(), cb, "TestHalfOpen", failFunc)
	// Second execution should trip the circuit
	Execute(context.Background(), cb, "TestHalfOpen", failFunc)
	assert.Equal(t, Open, cb.GetState())

	// Wait for the sleep window to elapse
	time.Sleep(20 * time.Millisecond)

	// Next execution should put the circuit in half-open state
	successFunc := func(ctx context.Context) (bool, error) {
		return true, nil
	}

	result, err := Execute(context.Background(), cb, "TestHalfOpen", successFunc)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, Closed, cb.GetState())

	// Circuit should be closed now
	result, err = Execute(context.Background(), cb, "TestHalfOpen", successFunc)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, Closed, cb.GetState())
}

func TestCircuitBreaker_Reset(t *testing.T) {
	// Create a circuit breaker
	cfg := DefaultConfig().
		WithEnabled(true).
		WithErrorThreshold(0.5).
		WithVolumeThreshold(2)

	logger := logging.NewContextLogger(zap.NewNop())
	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	cb := NewCircuitBreaker(cfg, options)
	assert.NotNil(t, cb)

	// Trip the circuit
	failFunc := func(ctx context.Context) (bool, error) {
		return false, errors.New("test error")
	}

	// First execution
	Execute(context.Background(), cb, "TestReset", failFunc)
	// Second execution should trip the circuit
	Execute(context.Background(), cb, "TestReset", failFunc)
	assert.Equal(t, Open, cb.GetState())

	// Reset the circuit
	cb.Reset()
	assert.Equal(t, Closed, cb.GetState())

	// Circuit should allow requests again
	successFunc := func(ctx context.Context) (bool, error) {
		return true, nil
	}

	result, err := Execute(context.Background(), cb, "TestReset", successFunc)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestCircuitBreaker_NilSafety(t *testing.T) {
	// Test with nil circuit breaker
	var cb *CircuitBreaker

	// Execute should work with nil circuit breaker
	successFunc := func(ctx context.Context) (bool, error) {
		return true, nil
	}

	result, err := Execute(context.Background(), cb, "TestNil", successFunc)
	assert.NoError(t, err)
	assert.True(t, result)

	// ExecuteWithFallback should work with nil circuit breaker
	failFunc := func(ctx context.Context) (bool, error) {
		return false, errors.New("test error")
	}

	fallbackFunc := func(ctx context.Context, err error) (bool, error) {
		return true, nil
	}

	result, err = ExecuteWithFallback(context.Background(), cb, "TestNil", failFunc, fallbackFunc)
	assert.NoError(t, err)
	assert.True(t, result)

	// GetState should return Closed for nil circuit breaker
	assert.Equal(t, Closed, cb.GetState())

	// Reset should not panic with nil circuit breaker
	assert.NotPanics(t, func() {
		cb.Reset()
	})
}