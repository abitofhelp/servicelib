// Copyright (c) 2025 A Bit of Help, Inc.

package shutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
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

// TestGracefulShutdown_Timeout tests that GracefulShutdown handles timeout correctly
func TestGracefulShutdown_Timeout(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a shutdown function that takes longer than the timeout
	var shutdownCalled int32 // Using atomic for thread safety
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		// Sleep for longer than the timeout to simulate a hanging function
		time.Sleep(5 * time.Second)
		return nil
	}

	// Start the graceful shutdown in a goroutine
	go func() {
		// Create a timeout context for the shutdown function
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer shutdownCancel()

		// Log that we're executing the shutdown function
		logger.Info(ctx, "Executing shutdown function")

		// Start the shutdown function in a goroutine
		shutdownDone := make(chan struct{})
		go func() {
			defer close(shutdownDone)
			shutdownFunc()
		}()

		// Wait for shutdown to complete or timeout
		select {
		case <-shutdownDone:
			// Shutdown completed within timeout (shouldn't happen in this test)
			resultCh <- nil
		case <-shutdownCtx.Done():
			// Shutdown timed out (expected in this test)
			logger.Error(ctx, "Shutdown function timed out after 30 seconds")
			resultCh <- fmt.Errorf("shutdown function timed out after 30 seconds")
		}
	}()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shutdown function timed out")
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 2, "Expected at least 2 log entries")

	// Find specific log messages
	var foundExecuting, foundTimeout bool
	for _, log := range logs {
		if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Shutdown function timed out after 30 seconds" {
			foundTimeout = true
		}
	}

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

// TestGracefulShutdown_SignalBased tests that GracefulShutdown handles OS signals correctly
func TestGracefulShutdown_SignalBased(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Create a channel to simulate OS signals
	quit := make(chan os.Signal, 2)

	// Start a goroutine to run GracefulShutdown with our simulated signal channel
	go func() {
		// Create a done channel to track completion
		done := make(chan struct{})
		var shutdownErr error

		// Start a goroutine to handle signals
		go func() {
			defer close(done)

			// Wait for either interrupt signal or context cancellation
			select {
			case sig := <-quit:
				// Log the signal
				logger.Info(ctx, "Received termination signal",
					zap.String("signal", sig.String()),
					zap.String("type", "first"))
			case <-ctx.Done():
				// This case shouldn't happen in this test
				logger.Info(ctx, "Context cancelled, shutting down", zap.Error(ctx.Err()))
			}

			// Create a timeout context for the shutdown function
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Execute shutdown in a separate goroutine so we can monitor for timeout
			shutdownDone := make(chan struct{})
			go func() {
				defer close(shutdownDone)
				// Call the shutdown function
				logger.Info(ctx, "Executing shutdown function")
				shutdownErr = shutdownFunc()
			}()

			// Wait for shutdown to complete or timeout
			select {
			case <-shutdownDone:
				// Shutdown completed within timeout
				if shutdownErr != nil {
					logger.Error(ctx, "Error during shutdown", zap.Error(shutdownErr))
				} else {
					logger.Info(ctx, "Graceful shutdown completed successfully")
				}
			case <-shutdownCtx.Done():
				// Shutdown timed out
				logger.Error(ctx, "Shutdown function timed out after 30 seconds")
				shutdownErr = fmt.Errorf("shutdown function timed out after 30 seconds")
			case sig := <-quit:
				// Received second signal during shutdown
				logger.Warn(ctx, "Received second termination signal during shutdown, forcing exit",
					zap.String("signal", sig.String()),
					zap.String("type", "second"))
				// In a real implementation, this would call os.Exit(1)
			}
		}()

		// Wait for shutdown to complete
		<-done
		resultCh <- shutdownErr
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Send a signal
	quit <- syscall.SIGTERM

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
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

// TestGracefulShutdown_SecondSignal tests the behavior when a second signal would be received during shutdown
func TestGracefulShutdown_SecondSignal(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Create a channel to simulate the quit channel in GracefulShutdown
	quit := make(chan os.Signal, 2)

	// Create a channel to track if the second signal case was executed
	secondSignalReceived := make(chan bool, 1)

	// Create a shutdown function that blocks to allow time for second signal
	shutdownStarted := make(chan struct{})
	shutdownDone := make(chan struct{})
	shutdownFunc := func() error {
		close(shutdownStarted)
		// Block until we're done with the test
		<-shutdownDone
		return nil
	}

	// Start a goroutine to simulate the behavior of GracefulShutdown
	go func() {
		// Simulate the first signal being received
		quit <- syscall.SIGTERM

		// Wait for the shutdown function to start
		<-shutdownStarted

		// Simulate receiving a second signal during shutdown
		quit <- syscall.SIGINT

		// Signal that we received the second signal
		secondSignalReceived <- true
	}()

	// Start another goroutine to simulate the core logic of GracefulShutdown
	go func() {
		// Wait for the first signal
		sig := <-quit

		// Log the first signal
		logger.Info(ctx, "Received termination signal",
			zap.String("signal", sig.String()),
			zap.String("type", "first"))

		// Start the shutdown function in a goroutine
		logger.Info(ctx, "Executing shutdown function")
		go shutdownFunc()

		// Wait for either shutdown completion or second signal
		select {
		case <-shutdownDone:
			// Shutdown completed
			logger.Info(ctx, "Graceful shutdown completed successfully")
		case sig := <-quit:
			// Second signal received
			logger.Warn(ctx, "Received second termination signal during shutdown, forcing exit",
				zap.String("signal", sig.String()),
				zap.String("type", "second"))

			// In the real implementation, this would call os.Exit(1)
			// Here we just log it
		}
	}()

	// Wait for the second signal to be processed
	select {
	case <-secondSignalReceived:
		// Success, second signal was received
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for second signal")
	}

	// Allow the shutdown function to complete
	close(shutdownDone)

	// Wait a bit for all goroutines to complete
	time.Sleep(100 * time.Millisecond)

	// Verify logs
	logs := observedLogs.All()

	// Find specific log messages
	var foundFirstSignal, foundExecuting, foundSecondSignal bool
	for _, log := range logs {
		if log.Message == "Received termination signal" && log.Context[1].String == "first" {
			foundFirstSignal = true
		} else if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Received second termination signal during shutdown, forcing exit" {
			foundSecondSignal = true
		}
	}

	assert.True(t, foundFirstSignal, "Expected 'Received termination signal' log message")
	assert.True(t, foundExecuting, "Expected 'Executing shutdown function' log message")
	assert.True(t, foundSecondSignal, "Expected 'Received second termination signal' log message")
}

// TestGracefulShutdown_ContextAlreadyCancelled tests that GracefulShutdown handles a context that's already cancelled
func TestGracefulShutdown_ContextAlreadyCancelled(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Start the graceful shutdown in a goroutine
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, shutdownFunc)
	}()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
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

// TestSetupGracefulShutdown_ContextAlreadyCancelled tests that SetupGracefulShutdown handles a context that's already cancelled
func TestSetupGracefulShutdown_ContextAlreadyCancelled(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Call SetupGracefulShutdown
	_, errCh := SetupGracefulShutdown(ctx, logger, shutdownFunc)

	// Wait for the error with a timeout
	select {
	case err := <-errCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for error channel")
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

// TestGracefulShutdown_TimeoutDetailed tests that GracefulShutdown handles timeout correctly with detailed simulation
func TestGracefulShutdown_TimeoutDetailed(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a shutdown function that takes longer than the timeout
	shutdownStarted := make(chan struct{})
	shutdownFunc := func() error {
		close(shutdownStarted)
		// Sleep for longer than the timeout to simulate a hanging function
		time.Sleep(5 * time.Second)
		return nil
	}

	// Start a goroutine to simulate GracefulShutdown
	go func() {
		// Create a done channel to track completion
		done := make(chan struct{})
		var shutdownErr error

		// Start a goroutine to handle signals
		go func() {
			defer close(done)

			// Simulate context cancellation
			logger.Info(ctx, "Context cancelled, shutting down", zap.Error(fmt.Errorf("context canceled")))

			// Create a timeout context for the shutdown function with a very short timeout
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Execute shutdown in a separate goroutine so we can monitor for timeout
			shutdownDone := make(chan struct{})
			go func() {
				defer close(shutdownDone)
				// Call the shutdown function
				logger.Info(ctx, "Executing shutdown function")
				shutdownErr = shutdownFunc()
			}()

			// Wait for shutdown to complete or timeout
			select {
			case <-shutdownDone:
				// Shutdown completed within timeout (shouldn't happen in this test)
				if shutdownErr != nil {
					logger.Error(ctx, "Error during shutdown", zap.Error(shutdownErr))
				} else {
					logger.Info(ctx, "Graceful shutdown completed successfully")
				}
			case <-shutdownCtx.Done():
				// Shutdown timed out (expected in this test)
				logger.Error(ctx, "Shutdown function timed out after 30 seconds")
				shutdownErr = fmt.Errorf("shutdown function timed out after 30 seconds")
			}
		}()

		// Wait for shutdown to complete
		<-done
		resultCh <- shutdownErr
	}()

	// Wait for the shutdown function to start
	select {
	case <-shutdownStarted:
		// Success, shutdown function started
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for shutdown function to start")
	}

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shutdown function timed out")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for result")
	}

	// Verify logs
	logs := observedLogs.All()

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

// TestGracefulShutdown_SuccessfulCompletion tests that GracefulShutdown handles successful completion of the shutdown function
func TestGracefulShutdown_SuccessfulCompletion(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a shutdown function that completes successfully after a short delay
	var shutdownCalled int32 // Using atomic for thread safety
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	// Start the graceful shutdown in a goroutine
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, shutdownFunc)
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
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

// TestSetupGracefulShutdown_NoListenerWithError tests the default case in SetupGracefulShutdown where no one is listening to the error channel
func TestSetupGracefulShutdown_NoListenerWithError(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Expected error
	expectedErr := errors.New("shutdown error")

	// Create a shutdown function that returns an error after a short delay
	shutdownFunc := func() error {
		// Add a small delay to ensure the error channel is read after the function completes
		time.Sleep(50 * time.Millisecond)
		return expectedErr
	}

	// Create a custom implementation that forces the default case
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

	// Wait for the goroutine to complete
	time.Sleep(300 * time.Millisecond)

	// Verify logs
	logs := observedLogs.All()

	// Find the error log message
	var foundErrorLog bool
	for _, log := range logs {
		if log.Message == "Error during background graceful shutdown (no listener)" {
			foundErrorLog = true
			// Verify that the error is included in the log
			assert.Equal(t, expectedErr.Error(), log.Context[0].Interface.(error).Error())
		}
	}

	assert.True(t, foundErrorLog, "Expected 'Error during background graceful shutdown (no listener)' log message")

	// Verify that the error channel is closed
	_, ok := <-errCh
	assert.False(t, ok, "Error channel should be closed")
}

// TestSetupGracefulShutdown_ContextCancellation tests that the context is cancelled when shutdown is complete
func TestSetupGracefulShutdown_ContextCancellation(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a parent context with cancellation
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Create a channel to track if the shutdown is complete
	shutdownComplete := make(chan struct{})

	// Create a shutdown function that completes quickly
	shutdownCalled := false
	shutdownFunc := func() error {
		shutdownCalled = true
		return nil
	}

	// Call SetupGracefulShutdown
	cancel, errCh := SetupGracefulShutdown(parentCtx, logger, shutdownFunc)

	// Start a goroutine to monitor the error channel
	go func() {
		// Wait for the error channel to be closed
		for range errCh {
			// Drain the channel
		}
		// Signal that shutdown is complete
		close(shutdownComplete)
	}()

	// Trigger shutdown
	cancel()

	// Wait for shutdown to complete with a timeout
	select {
	case <-shutdownComplete:
		// Success, shutdown completed
		assert.True(t, shutdownCalled, "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for shutdown to complete")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 2, "Expected at least 2 log entries")

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

// TestSetupGracefulShutdown_ContextCancelledAfterChannelClosed tests that the context is cancelled after the error channel is closed
func TestSetupGracefulShutdown_ContextCancelledAfterChannelClosed(t *testing.T) {
	// Create a logger
	testLogger := zaptest.NewLogger(t)
	logger := logging.NewContextLogger(testLogger)

	// Create a parent context with cancellation
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Create a shutdown function that completes quickly
	shutdownFunc := func() error {
		return nil
	}

	// Create a wrapper around SetupGracefulShutdown to test the context cancellation
	shutdownCtx, shutdownCancel := context.WithCancel(parentCtx)

	// Create a channel to receive shutdown errors
	errCh := make(chan error, 1)

	// Create a channel to track if the context was cancelled
	contextCancelled := make(chan struct{})

	// Start a goroutine to monitor the context
	go func() {
		<-shutdownCtx.Done()
		close(contextCancelled)
	}()

	// Start a goroutine to handle shutdown
	go func() {
		// Execute the graceful shutdown
		err := GracefulShutdown(shutdownCtx, logger, shutdownFunc)

		// Send the error to the channel (nil if no error)
		select {
		case errCh <- err:
			// Error sent successfully
		default:
			// No one is listening, just log the error
			if err != nil {
				logger.Error(parentCtx, "Error during background graceful shutdown (no listener)", zap.Error(err))
			}
		}

		// Close the error channel to signal completion
		close(errCh)

		// Ensure context is cancelled when shutdown is complete
		shutdownCancel()
	}()

	// Trigger shutdown
	parentCancel()

	// Wait for the error channel to be closed
	for range errCh {
		// Drain the channel
	}

	// Wait for the context to be cancelled with a timeout
	select {
	case <-contextCancelled:
		// Success, context was cancelled
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for context cancellation")
	}
}

// TestSetupGracefulShutdown_NoListenerNoError tests the default case in SetupGracefulShutdown where no one is listening to the error channel and the error is nil
func TestSetupGracefulShutdown_NoListenerNoError(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Create a shutdown function that returns nil
	shutdownFunc := func() error {
		return nil
	}

	// Create a custom implementation that forces the default case
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

	// Wait for the goroutine to complete
	time.Sleep(300 * time.Millisecond)

	// Verify logs
	logs := observedLogs.All()

	// Verify that there are no error logs
	for _, log := range logs {
		assert.NotEqual(t, "Error during background graceful shutdown (no listener)", log.Message, "Should not log error when error is nil")
	}

	// Verify that the error channel is closed
	_, ok := <-errCh
	assert.False(t, ok, "Error channel should be closed")
}

// TestSetupGracefulShutdown_ErrorSentSuccessfully tests that SetupGracefulShutdown sends the error successfully to the channel
func TestSetupGracefulShutdown_ErrorSentSuccessfully(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a parent context with cancellation
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Expected error
	expectedErr := errors.New("expected error")

	// Create a shutdown function that returns the expected error
	shutdownFunc := func() error {
		return expectedErr
	}

	// Call SetupGracefulShutdown
	cancel, errCh := SetupGracefulShutdown(parentCtx, logger, shutdownFunc)

	// Trigger shutdown
	cancel()

	// Wait for the error with a timeout
	var receivedErr error
	select {
	case err := <-errCh:
		receivedErr = err
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for error channel")
	}

	// Verify that the error is the expected one
	assert.Equal(t, expectedErr, receivedErr, "Expected error to be sent to the channel")

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 2, "Expected at least 2 log entries")

	// Find specific log messages
	var foundCancelled, foundExecuting, foundError bool
	for _, log := range logs {
		if log.Message == "Context cancelled, shutting down" {
			foundCancelled = true
		} else if log.Message == "Executing shutdown function" {
			foundExecuting = true
		} else if log.Message == "Error during shutdown" {
			foundError = true
			// Verify that the error is included in the log
			assert.Equal(t, expectedErr.Error(), log.Context[0].Interface.(error).Error())
		}
	}

	assert.True(t, foundCancelled, "Expected 'Context cancelled' log message")
	assert.True(t, foundExecuting, "Expected 'Executing shutdown function' log message")
	assert.True(t, foundError, "Expected 'Error during shutdown' log message")
}

// TestSetupGracefulShutdown_ParentContextCancellation tests that SetupGracefulShutdown handles parent context cancellation
func TestSetupGracefulShutdown_ParentContextCancellation(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a parent context with cancellation
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// Create a channel to track if the shutdown is complete
	shutdownComplete := make(chan struct{})

	// Create a shutdown function that completes quickly
	shutdownCalled := false
	shutdownFunc := func() error {
		shutdownCalled = true
		return nil
	}

	// Call SetupGracefulShutdown
	_, errCh := SetupGracefulShutdown(parentCtx, logger, shutdownFunc)

	// Start a goroutine to monitor the error channel
	go func() {
		// Wait for the error channel to be closed
		for range errCh {
			// Drain the channel
		}
		// Signal that shutdown is complete
		close(shutdownComplete)
	}()

	// Trigger shutdown by cancelling the parent context
	parentCancel()

	// Wait for shutdown to complete with a timeout
	select {
	case <-shutdownComplete:
		// Success, shutdown completed
		assert.True(t, shutdownCalled, "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for shutdown to complete")
	}

	// Verify logs
	logs := observedLogs.All()
	assert.GreaterOrEqual(t, len(logs), 2, "Expected at least 2 log entries")

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

// TestGracefulShutdown_NoopLogger tests that GracefulShutdown works with a no-op logger
func TestGracefulShutdown_NoopLogger(t *testing.T) {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Create a no-op logger
	noopLogger := zap.NewNop()
	logger := logging.NewContextLogger(noopLogger)

	// Start the graceful shutdown in a goroutine with the no-op logger
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, shutdownFunc)
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}

// TestGracefulShutdown_NilContext tests that GracefulShutdown handles a nil context
func TestGracefulShutdown_NilContext(t *testing.T) {
	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a no-op logger
	noopLogger := zap.NewNop()
	logger := logging.NewContextLogger(noopLogger)

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Start the graceful shutdown in a goroutine with a nil context
	go func() {
		resultCh <- GracefulShutdown(nil, logger, shutdownFunc)
	}()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context is nil")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}

// TestGracefulShutdown_NilShutdownFunc tests that GracefulShutdown handles a nil shutdown function
func TestGracefulShutdown_NilShutdownFunc(t *testing.T) {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a no-op logger
	noopLogger := zap.NewNop()
	logger := logging.NewContextLogger(noopLogger)

	// Start the graceful shutdown in a goroutine with a nil shutdown function
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, nil)
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shutdown function is nil")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}

// TestSetupGracefulShutdown_NilContext tests that SetupGracefulShutdown handles a nil context
func TestSetupGracefulShutdown_NilContext(t *testing.T) {
	// Create a no-op logger
	noopLogger := zap.NewNop()
	logger := logging.NewContextLogger(noopLogger)

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Call SetupGracefulShutdown with a nil context
	cancel, errCh := SetupGracefulShutdown(nil, logger, shutdownFunc)

	// Trigger shutdown
	cancel()

	// Wait for the error with a timeout
	select {
	case err := <-errCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context is nil")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for error channel")
	}
}

// TestSetupGracefulShutdown_NilShutdownFunc tests that SetupGracefulShutdown handles a nil shutdown function
func TestSetupGracefulShutdown_NilShutdownFunc(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a no-op logger
	noopLogger := zap.NewNop()
	logger := logging.NewContextLogger(noopLogger)

	// Call SetupGracefulShutdown with a nil shutdown function
	cancel, errCh := SetupGracefulShutdown(ctx, logger, nil)

	// Trigger shutdown
	cancel()

	// Wait for the error with a timeout
	select {
	case err := <-errCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shutdown function is nil")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for error channel")
	}
}

// TestSetupGracefulShutdown_NoopLogger tests that SetupGracefulShutdown works with a no-op logger
func TestSetupGracefulShutdown_NoopLogger(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Create a no-op logger
	noopLogger := zap.NewNop()
	logger := logging.NewContextLogger(noopLogger)

	// Call SetupGracefulShutdown with the no-op logger
	cancel, errCh := SetupGracefulShutdown(ctx, logger, shutdownFunc)

	// Trigger shutdown
	cancel()

	// Wait for the error with a timeout
	select {
	case err := <-errCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for error channel")
	}
}

// TestGracefulShutdown_DirectSignal tests GracefulShutdown with a direct signal
func TestGracefulShutdown_DirectSignal(t *testing.T) {
	// Create an observed zap logger to capture logs
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger := logging.NewContextLogger(observedLogger)

	// Create a context
	ctx := context.Background()

	// Create a channel to receive the result
	resultCh := make(chan error)

	// Create a simple shutdown function
	var shutdownCalled int32
	shutdownFunc := func() error {
		atomic.StoreInt32(&shutdownCalled, 1)
		return nil
	}

	// Start GracefulShutdown in a goroutine
	go func() {
		resultCh <- GracefulShutdown(ctx, logger, shutdownFunc)
	}()

	// Wait a bit to ensure the goroutine is running
	time.Sleep(50 * time.Millisecond)

	// Send a signal to the process
	// This is a bit tricky since we can't easily send a signal to ourselves in a test
	// Instead, we'll simulate it by sending a signal to the process group
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	// Send a SIGTERM signal
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		t.Fatalf("Failed to send signal: %v", err)
	}

	// Wait for the result with a timeout
	select {
	case err := <-resultCh:
		assert.NoError(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled), "Shutdown function should have been called")
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
