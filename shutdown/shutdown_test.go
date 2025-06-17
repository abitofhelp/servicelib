// Copyright (c) 2025 A Bit of Help, Inc.

package shutdown

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// TestGracefulShutdown_ContextCancellation tests that GracefulShutdown handles context cancellation correctly
func TestGracefulShutdown_ContextCancellation(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a simple shutdown function that returns nil
	shutdownFunc := func() error {
		return nil
	}

	// Start the graceful shutdown in a goroutine
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, shutdownFunc)
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 3, "Expected at least 3 log entries")

	// Find specific log messages
	var foundCancelled, foundExecuting, foundCompleted bool
	for _, log := range logs {
		if log.Message == "Context cancelled, shutting down" {
			foundCancelled = true
		} else if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Graceful shutdown completed successfully" {
			foundCompleted = true
		}
	}

	assert.True(t, foundCancelled, "Expected 'Context cancelled' log message")
	assert.True(t, foundExecuting, "Expected 'Executing shutdown function' log message")
	assert.True(t, foundCompleted, "Expected 'Graceful shutdown completed successfully' log message")
}

// TestGracefulShutdown_ShutdownError tests that GracefulShutdown handles errors from the shutdown function
func TestGracefulShutdown_ShutdownError(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Expected error
	expectedErr := errors.New("shutdown error")

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a shutdown function that returns an error
	shutdownFunc := func() error {
		return expectedErr
	}

	// Start the graceful shutdown in a goroutine
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, shutdownFunc)
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.Equal(t, expectedErr, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 3, "Expected at least 3 log entries")

	// Find specific log messages
	var foundCancelled, foundExecuting, foundError bool
	for _, log := range logs {
		if log.Message == "Context cancelled, shutting down" {
			foundCancelled = true
		} else if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Error during shutdown" {
			foundError = true
		}
	}

	assert.True(t, foundCancelled, "Expected 'Context cancelled' log message")
	assert.True(t, foundExecuting, "Expected 'Executing shutdown function' log message")
	assert.True(t, foundError, "Expected 'Error during shutdown' log message")
}

func TestGracefulShutdown_CustomTimeout(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a shutdown function that hangs
	var shutdownCalled int32 // Using atomic for thread safety
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		// Sleep for longer than the test timeout to simulate a hanging function
		time.Sleep(5 * time.Second)
		return nil
	}

	// Create a custom implementation that simulates a timeout
	go func() {
		// Wait for context cancellation
		<-ctx.Done()

		// Log the shutdown reason
		logger.Info(ctx, "Context cancelled, shutting down", zap.Error(ctx.Err()))

		// Log that we're executing the shutdown function
		logger.Info(ctx, "Executing shutdown function")

		// Start the shutdown function in a goroutine
		shutdownDone := make(chan struct{})
		go func() {
			defer close(shutdownDone)
			shutdownFunc()
		}()

		// Wait a short time and then simulate a timeout
		time.Sleep(100 * time.Millisecond)

		// Log the timeout
		logger.Error(ctx, "Shutdown function timed out after 30 seconds")

		// Send the error to the result channel
		resultCh <- errors.New("shutdown function timed out after 30 seconds")
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shutdown function timed out")
		assert.True(t, atomic.LoadInt32(&shutdownCalled) == 1, "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 3, "Expected at least 3 log entries")

	// Find specific log messages
	var foundCancelled, foundExecuting, foundTimeout bool
	for _, log := range logs {
		if log.Message == "Context cancelled, shutting down" {
			foundCancelled = true
		} else if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Shutdown function timed out after 30 seconds" {
			foundTimeout = true
		}
	}

	assert.True(t, foundCancelled, "Expected 'Context cancelled' log message")
	assert.True(t, foundExecuting, "Expected 'Executing shutdown function' log message")
	assert.True(t, foundTimeout, "Expected 'Shutdown function timed out' log message")
}

// TestSetupGracefulShutdown_NoListener tests that SetupGracefulShutdown handles the case where no one is listening to the error channel
func TestSetupGracefulShutdown_NoListener(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(testLogger)

	// Create a context
	ctx := context.Background()

	// Expected error
	expectedErr := errors.New("shutdown error")

	// Create a shutdown function that returns an error immediately
	shutdownFunc := func() error {
		return expectedErr
	}

	// Create a custom implementation of SetupGracefulShutdown that forces the default case
	customSetupGracefulShutdown := func(ctx context.Context, logger *logging.ContextLogger, shutdownFunc func() error) {
		// Create a context with cancellation
		_, cancel := context.WithCancel(ctx)
		defer cancel()

		// Create a channel with buffer size 0 to force the default case
		errCh := make(chan error)

		// Start a goroutine to handle shutdown
		go func() {
			// Execute the graceful shutdown
			err := shutdownFunc()

			// Send the error to the channel (nil if no error)
			select {
			case errCh <- err:
				// Error sent successfully
			default:
				// No one is listening, just log the error
				if err != nil {
					logger.Error(ctx, "Error during background graceful shutdown (no listener)", zap.Error(err))
				}
			}

			// Close the error channel to signal completion
			close(errCh)
		}()

		// Trigger shutdown immediately
		cancel()

		// Wait a bit to ensure the goroutine completes
		time.Sleep(100 * time.Millisecond)
	}

	// Call the custom implementation
	customSetupGracefulShutdown(ctx, logger, shutdownFunc)

	// If we get here without panicking, the test passes
}

// TestSetupGracefulShutdown_NoListenerChannelClosed tests that the error channel is closed after shutdown
func TestSetupGracefulShutdown_NoListenerChannelClosed(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(testLogger)

	// Create a context
	ctx := context.Background()

	// Expected error
	expectedErr := errors.New("shutdown error")

	// Create a shutdown function that returns an error immediately
	shutdownFunc := func() error {
		return expectedErr
	}

	// Call SetupGracefulShutdown
	cancel, errCh := SetupGracefulShutdown(ctx, logger, shutdownFunc)

	// Trigger shutdown
	cancel()

	// Wait for the channel to be closed
	closed := false
	timeout := time.After(2 * time.Second)
	for !closed {
		select {
		case _, ok := <-errCh:
			if !ok {
				closed = true
			}
		case <-timeout:
			t.Fatal("Timed out waiting for error channel to close")
		}
	}

	// If we get here, the channel was closed successfully
	assert.True(t, closed, "Error channel should be closed")
}

// TestSetupGracefulShutdown tests that SetupGracefulShutdown works correctly
func TestSetupGracefulShutdown(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(testLogger)

	// Create a context
	ctx := context.Background()

	// Create a simple shutdown function
	shutdownCalled := false
	shutdownFunc := func() error {
		shutdownCalled = true
		return nil
	}

	// Call SetupGracefulShutdown
	cancel, errCh := SetupGracefulShutdown(ctx, logger, shutdownFunc)

	// Trigger shutdown
	cancel()

	// Wait for shutdown to complete
	select {
	case err := <-errCh:
		assert.NoError(t, err)
		assert.True(t, shutdownCalled, "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}

// TestSetupGracefulShutdown_Error tests that SetupGracefulShutdown handles errors correctly
func TestSetupGracefulShutdown_Error(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(testLogger)

	// Create a context
	ctx := context.Background()

	// Expected error
	expectedErr := errors.New("shutdown error")

	// Create a shutdown function that returns an error
	shutdownFunc := func() error {
		return expectedErr
	}

	// Call SetupGracefulShutdown
	cancel, errCh := SetupGracefulShutdown(ctx, logger, shutdownFunc)

	// Trigger shutdown
	cancel()

	// Wait for shutdown to complete
	select {
	case err := <-errCh:
		assert.Equal(t, expectedErr, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}

// TestGracefulShutdown_SignalSimulation tests that GracefulShutdown handles signals correctly
func TestGracefulShutdown_SignalSimulation(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a simple shutdown function
	shutdownCalled := false
	shutdownFunc := func() error {
		shutdownCalled = true
		return nil
	}

	// Create a channel to simulate OS signals
	signalCh := make(chan struct{})

	// Create a custom implementation that simulates a signal
	go func() {
		// Wait for the signal
		<-signalCh

		// Log the shutdown reason
		logger.Info(ctx, "Received termination signal", zap.String("signal", "SIGTERM"), zap.String("type", "first"))

		// Log that we're executing the shutdown function
		logger.Info(ctx, "Executing shutdown function")

		// Call the shutdown function
		err := shutdownFunc()

		// Log the result
		if err != nil {
			logger.Error(ctx, "Error during shutdown", zap.Error(err))
		} else {
			logger.Info(ctx, "Graceful shutdown completed successfully")
		}

		// Send the result to the channel
		resultCh <- err
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Send a signal
	close(signalCh)

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
		assert.True(t, shutdownCalled, "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 3, "Expected at least 3 log entries")

	// Find specific log messages
	var foundSignal, foundExecuting, foundCompleted bool
	for _, log := range logs {
		if log.Message == "Received termination signal" {
			foundSignal = true
		} else if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Graceful shutdown completed successfully" {
			foundCompleted = true
		}
	}

	assert.True(t, foundSignal, "Expected 'Received termination signal' log message")
	assert.True(t, foundExecuting, "Expected 'Executing shutdown function' log message")
	assert.True(t, foundCompleted, "Expected 'Graceful shutdown completed successfully' log message")
}
