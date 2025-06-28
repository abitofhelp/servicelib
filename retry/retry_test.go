// Copyright (c) 2025 A Bit of Help, Inc.

package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	serviceErrors "github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 100*time.Millisecond, config.InitialBackoff)
	assert.Equal(t, 2*time.Second, config.MaxBackoff)
	assert.Equal(t, 2.0, config.BackoffFactor)
	assert.Equal(t, 0.2, config.JitterFactor)
	assert.Nil(t, config.RetryableErrors)
}

func TestConfigWithMethods(t *testing.T) {
	config := DefaultConfig()

	// Test WithMaxRetries
	newConfig := config.WithMaxRetries(5)
	assert.Equal(t, 5, newConfig.MaxRetries)

	// Test WithInitialBackoff
	newConfig = config.WithInitialBackoff(200 * time.Millisecond)
	assert.Equal(t, 200*time.Millisecond, newConfig.InitialBackoff)

	// Test WithMaxBackoff
	newConfig = config.WithMaxBackoff(5 * time.Second)
	assert.Equal(t, 5*time.Second, newConfig.MaxBackoff)

	// Test WithBackoffFactor
	newConfig = config.WithBackoffFactor(3.0)
	assert.Equal(t, 3.0, newConfig.BackoffFactor)

	// Test WithJitterFactor
	newConfig = config.WithJitterFactor(0.5)
	assert.Equal(t, 0.5, newConfig.JitterFactor)

	// Test WithRetryableErrors
	retryableErrors := []error{errors.New("error1"), errors.New("error2")}
	newConfig = config.WithRetryableErrors(retryableErrors)
	assert.Equal(t, retryableErrors, newConfig.RetryableErrors)
}

func TestDoSuccessFirstAttempt(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	// Function that succeeds on first attempt
	fn := func(ctx context.Context) error {
		return nil
	}

	err := Do(ctx, fn, config, nil)
	assert.NoError(t, err)
}

func TestDoSuccessAfterRetries(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	attempts := 0
	maxAttempts := 2

	// Function that succeeds after maxAttempts attempts
	fn := func(ctx context.Context) error {
		attempts++
		if attempts <= maxAttempts {
			return errors.New("temporary error")
		}
		return nil
	}

	// Always retry
	isRetryable := func(err error) bool {
		return true
	}

	err := Do(ctx, fn, config, isRetryable)
	assert.NoError(t, err)
	assert.Equal(t, maxAttempts+1, attempts)
}

func TestDoMaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	// Function that always fails
	fn := func(ctx context.Context) error {
		return errors.New("persistent error")
	}

	// Always retry
	isRetryable := func(err error) bool {
		return true
	}

	err := Do(ctx, fn, config, isRetryable)
	assert.Error(t, err)
	assert.True(t, serviceErrors.IsRetryError(err))
}

func TestDoNonRetryableError(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	// Function that always fails
	fn := func(ctx context.Context) error {
		return errors.New("non-retryable error")
	}

	// Never retry
	isRetryable := func(err error) bool {
		return false
	}

	err := Do(ctx, fn, config, isRetryable)
	assert.Error(t, err)
	assert.Equal(t, "non-retryable error", err.Error())
}

func TestDoContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := DefaultConfig()

	// Cancel the context immediately
	cancel()

	// Function that should never be called
	fn := func(ctx context.Context) error {
		t.Fatal("Function should not be called after context cancellation")
		return nil
	}

	err := Do(ctx, fn, config, nil)
	assert.Error(t, err)
	assert.True(t, serviceErrors.IsContextError(err))
}

func TestDoContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	config := DefaultConfig().
		WithInitialBackoff(50 * time.Millisecond).
		WithMaxRetries(5)

	// Function that always fails
	fn := func(ctx context.Context) error {
		return errors.New("error")
	}

	// Always retry
	isRetryable := func(err error) bool {
		return true
	}

	// Sleep to ensure context times out
	time.Sleep(20 * time.Millisecond)

	err := Do(ctx, fn, config, isRetryable)
	assert.Error(t, err)
	assert.True(t, serviceErrors.IsContextError(err))
}

func TestDoBackoffAndJitter(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig().
		WithMaxRetries(3).
		WithInitialBackoff(10 * time.Millisecond).
		WithMaxBackoff(100 * time.Millisecond).
		WithBackoffFactor(2.0).
		WithJitterFactor(0.1)

	attempts := 0
	startTimes := make([]time.Time, 0, config.MaxRetries+1)

	// Function that records attempt times and always fails
	fn := func(ctx context.Context) error {
		startTimes = append(startTimes, time.Now())
		attempts++
		return errors.New("error")
	}

	// Always retry
	isRetryable := func(err error) bool {
		return true
	}

	err := Do(ctx, fn, config, isRetryable)
	assert.Error(t, err)
	assert.Equal(t, config.MaxRetries+1, attempts)

	// Verify that backoff is increasing
	for i := 1; i < len(startTimes); i++ {
		elapsed := startTimes[i].Sub(startTimes[i-1])
		assert.True(t, elapsed > 0, "Backoff should be greater than 0")

		// For the last attempt, check that backoff doesn't exceed max
		if i == len(startTimes)-1 {
			maxWithJitter := time.Duration(float64(config.MaxBackoff) * (1 + config.JitterFactor))
			assert.True(t, elapsed <= maxWithJitter,
				"Last backoff should not exceed max backoff plus jitter")
		}
	}
}

func TestIsNetworkError(t *testing.T) {
	// Test that the deprecated IsNetworkError function correctly delegates to errors.IsNetworkError
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "nil error",
			err:  nil,
		},
		{
			name: "connection refused error",
			err:  errors.New("connection refused"),
		},
		{
			name: "connection reset error",
			err:  errors.New("connection reset by peer"),
		},
		{
			name: "non-network error",
			err:  errors.New("some other error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify that the deprecated function returns the same result as the errors package function
			expected := serviceErrors.IsNetworkError(tc.err)
			result := IsNetworkError(tc.err)
			assert.Equal(t, expected, result, "IsNetworkError should delegate to errors.IsNetworkError")
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "timeout error",
			err:      errors.New("operation timed out"),
			expected: true,
		},
		{
			name:     "deadline exceeded error",
			err:      errors.New("context deadline exceeded"),
			expected: true,
		},
		{
			name:     "non-timeout error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsTimeoutError(tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsTransientError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "network error",
			err:      errors.New("connection refused"),
			expected: true,
		},
		{
			name:     "timeout error",
			err:      errors.New("operation timed out"),
			expected: true,
		},
		{
			name:     "rate limit error",
			err:      errors.New("rate limit exceeded"),
			expected: true,
		},
		{
			name:     "non-transient error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsTransientError(tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDeprecatedErrorFunctions(t *testing.T) {
	// Test that the deprecated error detection functions correctly delegate to the errors package functions
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "nil error",
			err:  nil,
		},
		{
			name: "network error",
			err:  errors.New("connection refused"),
		},
		{
			name: "timeout error",
			err:  errors.New("operation timed out"),
		},
		{
			name: "context deadline exceeded",
			err:  context.DeadlineExceeded,
		},
		{
			name: "transient error",
			err:  errors.New("rate limit exceeded"),
		},
		{
			name: "non-special error",
			err:  errors.New("some other error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test IsNetworkError
			expected := serviceErrors.IsNetworkError(tc.err)
			result := IsNetworkError(tc.err)
			assert.Equal(t, expected, result, "IsNetworkError should delegate to serviceErrors.IsNetworkError")

			// Test IsTimeoutError
			expected = serviceErrors.IsTimeout(tc.err)
			result = IsTimeoutError(tc.err)
			assert.Equal(t, expected, result, "IsTimeoutError should delegate to serviceErrors.IsTimeout")

			// Test IsTransientError
			expected = serviceErrors.IsTransientError(tc.err)
			result = IsTransientError(tc.err)
			assert.Equal(t, expected, result, "IsTransientError should delegate to serviceErrors.IsTransientError")
		})
	}
}

func TestDoWithOptionsCustomLoggerAndTracer(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	// Create a custom logger
	logger := logging.NewContextLogger(nil) // Using nil will create a no-op logger

	// Create a custom tracer (using no-op tracer for testing)
	tracer := telemetry.NewNoopTracer()

	// Create options with custom logger and tracer
	options := Options{
		Logger: logger,
		Tracer: tracer,
	}

	// Function that succeeds on first attempt
	fn := func(ctx context.Context) error {
		return nil
	}

	// Execute with custom options
	err := DoWithOptions(ctx, fn, config, nil, options)
	assert.NoError(t, err)

	// Test with a function that fails
	attempts := 0
	maxAttempts := 2

	// Function that succeeds after maxAttempts attempts
	fnWithRetries := func(ctx context.Context) error {
		attempts++
		if attempts <= maxAttempts {
			return errors.New("temporary error")
		}
		return nil
	}

	// Always retry
	isRetryable := func(err error) bool {
		return true
	}

	// Reset attempts counter
	attempts = 0

	// Execute with custom options
	err = DoWithOptions(ctx, fnWithRetries, config, isRetryable, options)
	assert.NoError(t, err)
	assert.Equal(t, maxAttempts+1, attempts)
}
