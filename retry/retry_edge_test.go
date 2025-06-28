// Copyright (c) 2025 A Bit of Help, Inc.

package retry

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJitterFactorExtremes(t *testing.T) {
	// Test with jitter factor = 0 (no jitter)
	t.Run("jitter factor = 0", func(t *testing.T) {
		options := Options{
			Tracer: &SimpleMockTracer{},
		}

		config := DefaultConfig().WithJitterFactor(0)
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

		err := DoWithOptions(context.Background(), fn, config, isRetryable, options)
		assert.NoError(t, err)
		assert.Equal(t, maxAttempts+1, attempts)
	})

	// Test with jitter factor = 1 (maximum jitter)
	t.Run("jitter factor = 1", func(t *testing.T) {
		options := Options{
			Tracer: &SimpleMockTracer{},
		}

		config := DefaultConfig().WithJitterFactor(1)
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

		err := DoWithOptions(context.Background(), fn, config, isRetryable, options)
		assert.NoError(t, err)
		assert.Equal(t, maxAttempts+1, attempts)
	})
}

func TestBackoffFactorExtremes(t *testing.T) {
	// Test with very small backoff factor
	t.Run("small backoff factor", func(t *testing.T) {
		options := Options{
			Tracer: &SimpleMockTracer{},
		}

		config := DefaultConfig().
			WithBackoffFactor(0.1).
			WithInitialBackoff(10 * time.Millisecond).
			WithMaxBackoff(100 * time.Millisecond)

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

		err := DoWithOptions(context.Background(), fn, config, isRetryable, options)
		assert.NoError(t, err)
		assert.Equal(t, maxAttempts+1, attempts)
	})

	// Test with very large backoff factor
	t.Run("large backoff factor", func(t *testing.T) {
		options := Options{
			Tracer: &SimpleMockTracer{},
		}

		config := DefaultConfig().
			WithBackoffFactor(10.0).
			WithInitialBackoff(1 * time.Millisecond).
			WithMaxBackoff(100 * time.Millisecond)

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

		err := DoWithOptions(context.Background(), fn, config, isRetryable, options)
		assert.NoError(t, err)
		assert.Equal(t, maxAttempts+1, attempts)
	})
}

func TestConcurrentUsage(t *testing.T) {
	// Test concurrent usage
	t.Run("concurrent usage", func(t *testing.T) {
		options := Options{
			Tracer: &SimpleMockTracer{},
		}

		config := DefaultConfig().
			WithMaxRetries(2).
			WithInitialBackoff(1 * time.Millisecond).
			WithMaxBackoff(10 * time.Millisecond)

		const numGoroutines = 10
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()

				attempts := 0
				maxAttempts := 1

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

				err := DoWithOptions(context.Background(), fn, config, isRetryable, options)
				assert.NoError(t, err)
				assert.Equal(t, maxAttempts+1, attempts)
			}(i)
		}

		wg.Wait()
	})
}

func TestDoWithOptionsUsingMocks(t *testing.T) {
	options := Options{
		Tracer: &SimpleMockTracer{},
	}

	config := DefaultConfig()

	// Function that succeeds on first attempt
	fn := func(ctx context.Context) error {
		return nil
	}

	err := DoWithOptions(context.Background(), fn, config, nil, options)
	assert.NoError(t, err)
}
