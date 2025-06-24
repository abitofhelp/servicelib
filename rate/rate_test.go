// Copyright (c) 2025 A Bit of Help, Inc.

package rate

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestNewRateLimiter(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Test cases
	tests := []struct {
		name     string
		config   Config
		options  Options
		expected bool
	}{
		{
			name: "Enabled rate limiter",
			config: DefaultConfig().
				WithEnabled(true),
			options: DefaultOptions().
				WithName("test").
				WithLogger(logger),
			expected: true,
		},
		{
			name: "Disabled rate limiter",
			config: DefaultConfig().
				WithEnabled(false),
			options: DefaultOptions().
				WithName("test").
				WithLogger(logger),
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRateLimiter(tc.config, tc.options)

			if tc.expected {
				assert.NotNil(t, rl)
			} else {
				assert.Nil(t, rl)
			}
		})
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a rate limiter with a small burst size for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithRequestsPerSecond(10).
		WithBurstSize(5) // Small burst size for testing

	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	rl := NewRateLimiter(cfg, options)
	require.NotNil(t, rl)

	// First 5 requests should be allowed (burst size)
	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow(), "Request %d should be allowed", i+1)
	}

	// Next request should be denied (burst size exceeded)
	assert.False(t, rl.Allow(), "Request 6 should be denied")

	// Wait for token refill (at least 100ms for 1 token at 10 RPS)
	time.Sleep(200 * time.Millisecond)

	// Should be allowed again after refill
	assert.True(t, rl.Allow(), "Request after wait should be allowed")
}

func TestExecute_Success(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a rate limiter
	cfg := DefaultConfig().
		WithEnabled(true).
		WithRequestsPerSecond(10).
		WithBurstSize(5)

	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	rl := NewRateLimiter(cfg, options)
	require.NotNil(t, rl)

	// Test successful execution
	ctx := context.Background()
	callCount := 0

	// Function that should succeed
	fn := func(ctx context.Context) (string, error) {
		callCount++
		return "success", nil
	}

	// Execute the function
	result, err := Execute(ctx, rl, "test-operation", fn)
	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 1, callCount)
}

func TestExecute_RateLimitExceeded(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a rate limiter with a small burst size for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithRequestsPerSecond(10).
		WithBurstSize(1) // Small burst size for testing

	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	rl := NewRateLimiter(cfg, options)
	require.NotNil(t, rl)

	// Test rate limit exceeded
	ctx := context.Background()

	// First request should succeed
	result, err := Execute(ctx, rl, "test-operation", func(ctx context.Context) (string, error) {
		return "success", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "success", result)

	// Second request should fail with rate limit exceeded
	result, err = Execute(ctx, rl, "test-operation", func(ctx context.Context) (string, error) {
		t.Fatal("This function should not be called when rate limit is exceeded")
		return "", nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Equal(t, "", result)
	assert.True(t, errors.Is(err, errors.New(errors.ResourceExhaustedCode, "")))
}

func TestExecuteWithWait(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a rate limiter with a small burst size for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithRequestsPerSecond(10).
		WithBurstSize(1) // Small burst size for testing

	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	rl := NewRateLimiter(cfg, options)
	require.NotNil(t, rl)

	// Test execute with wait
	ctx := context.Background()
	callCount := 0

	// First request should succeed immediately
	start := time.Now()
	result, err := ExecuteWithWait(ctx, rl, "test-operation", func(ctx context.Context) (string, error) {
		callCount++
		return "success", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 1, callCount)
	assert.Less(t, time.Since(start), 50*time.Millisecond, "First request should not wait")

	// Second request should wait for a token
	start = time.Now()
	result, err = ExecuteWithWait(ctx, rl, "test-operation", func(ctx context.Context) (string, error) {
		callCount++
		return "success", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 2, callCount)
	assert.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond, "Second request should wait for a token")
}

func TestExecuteWithWait_ContextCancellation(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a rate limiter with a small burst size for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithRequestsPerSecond(1).
		WithBurstSize(1) // Small burst size for testing

	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	rl := NewRateLimiter(cfg, options)
	require.NotNil(t, rl)

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// First request should succeed
	result, err := Execute(ctx, rl, "test-operation", func(ctx context.Context) (string, error) {
		return "success", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "success", result)

	// Cancel the context
	cancel()

	// Second request should fail with context canceled
	result, err = ExecuteWithWait(ctx, rl, "test-operation", func(ctx context.Context) (string, error) {
		t.Fatal("This function should not be called when context is canceled")
		return "", nil
	})
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, "", result)
}

func TestRateLimiter_Reset(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a rate limiter with a small burst size for testing
	cfg := DefaultConfig().
		WithEnabled(true).
		WithRequestsPerSecond(10).
		WithBurstSize(2) // Small burst size for testing

	options := DefaultOptions().
		WithName("test").
		WithLogger(logger)

	rl := NewRateLimiter(cfg, options)
	require.NotNil(t, rl)

	// Use up all tokens
	assert.True(t, rl.Allow())
	assert.True(t, rl.Allow())
	assert.False(t, rl.Allow())

	// Reset the rate limiter
	rl.Reset()

	// Should have full burst size again
	assert.True(t, rl.Allow())
	assert.True(t, rl.Allow())
	assert.False(t, rl.Allow())
}

func TestRateLimiter_NilSafety(t *testing.T) {
	// Test that nil rate limiter operations don't panic
	var rl *RateLimiter

	assert.NotPanics(t, func() {
		_ = rl.Allow()
	})

	assert.NotPanics(t, func() {
		_, _ = Execute(context.Background(), rl, "test-operation", func(ctx context.Context) (string, error) {
			return "success", nil
		})
	})

	assert.NotPanics(t, func() {
		_, _ = ExecuteWithWait(context.Background(), rl, "test-operation", func(ctx context.Context) (string, error) {
			return "success", nil
		})
	})

	assert.NotPanics(t, func() {
		rl.Reset()
	})
}

func TestConfig_WithMethods(t *testing.T) {
	// Test the WithEnabled method
	cfg := DefaultConfig()
	assert.True(t, cfg.Enabled)
	cfg = cfg.WithEnabled(false)
	assert.False(t, cfg.Enabled)

	// Test the WithRequestsPerSecond method
	cfg = DefaultConfig()
	assert.Equal(t, 100, cfg.RequestsPerSecond)
	cfg = cfg.WithRequestsPerSecond(200)
	assert.Equal(t, 200, cfg.RequestsPerSecond)
	cfg = cfg.WithRequestsPerSecond(0) // Should set to 1
	assert.Equal(t, 1, cfg.RequestsPerSecond)
	cfg = cfg.WithRequestsPerSecond(-1) // Should set to 1
	assert.Equal(t, 1, cfg.RequestsPerSecond)

	// Test the WithBurstSize method
	cfg = DefaultConfig()
	assert.Equal(t, 50, cfg.BurstSize)
	cfg = cfg.WithBurstSize(100)
	assert.Equal(t, 100, cfg.BurstSize)
	cfg = cfg.WithBurstSize(0) // Should set to 1
	assert.Equal(t, 1, cfg.BurstSize)
	cfg = cfg.WithBurstSize(-1) // Should set to 1
	assert.Equal(t, 1, cfg.BurstSize)
}

func TestOptions_WithMethods(t *testing.T) {
	// Test the WithLogger method
	options := DefaultOptions()
	assert.Nil(t, options.Logger)
	logger := logging.NewContextLogger(zap.NewNop())
	options = options.WithLogger(logger)
	assert.Equal(t, logger, options.Logger)

	// Test the WithName method
	options = DefaultOptions()
	assert.Equal(t, "default", options.Name)
	options = options.WithName("test")
	assert.Equal(t, "test", options.Name)
}