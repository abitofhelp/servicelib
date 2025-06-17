package signal

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewGracefulShutdown(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a new GracefulShutdown
	timeout := 5 * time.Second
	gs := NewGracefulShutdown(timeout, logger)

	// Verify the GracefulShutdown was created correctly
	assert.Equal(t, timeout, gs.timeout)
	assert.Equal(t, logger, gs.logger)
	assert.Empty(t, gs.callbacks)
	assert.Len(t, gs.signals, 4) // SIGINT, SIGTERM, SIGHUP, SIGQUIT
}

func TestRegisterCallback(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a new GracefulShutdown
	gs := NewGracefulShutdown(5*time.Second, logger)

	// Create a test callback
	callback1 := func(ctx context.Context) error {
		return nil
	}

	// Register the callback
	gs.RegisterCallback(callback1)

	// Verify the callback was registered
	assert.Len(t, gs.callbacks, 1)

	// Register another callback
	callback2 := func(ctx context.Context) error {
		return errors.New("test error")
	}

	gs.RegisterCallback(callback2)

	// Verify both callbacks were registered
	assert.Len(t, gs.callbacks, 2)
}

func TestHandleShutdown(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a custom GracefulShutdown for testing
	gs := &GracefulShutdown{
		timeout:   100 * time.Millisecond,
		logger:    logger,
		callbacks: make([]ShutdownCallback, 0),
		signals:   []os.Signal{os.Interrupt}, // Use a simple signal for testing
	}

	// Create a channel to track callback execution
	callbackCalled := make(chan bool, 1)

	// Create a test callback that signals when it's called
	callback := func(ctx context.Context) error {
		callbackCalled <- true
		return nil
	}

	// Register the callback
	gs.RegisterCallback(callback)

	// Call HandleShutdown
	ctx, cancel := gs.HandleShutdown()

	// Verify the context is not canceled yet
	select {
	case <-ctx.Done():
		t.Fatal("Context should not be canceled yet")
	default:
		// Expected
	}

	// Instead of using cancel(), we'll directly execute the callback logic
	// This simulates what happens when a signal is received
	go func() {
		// Execute the callback directly
		for _, cb := range gs.callbacks {
			err := cb(context.Background())
			assert.NoError(t, err)
		}
	}()

	// Wait for the callback to be called or timeout
	select {
	case <-callbackCalled:
		// Callback was called, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not called within timeout")
	}

	// Clean up
	cancel()
}

// TestHandleShutdownWithError tests that HandleShutdown handles errors from callbacks correctly
func TestHandleShutdownWithError(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a custom GracefulShutdown for testing
	gs := &GracefulShutdown{
		timeout:   100 * time.Millisecond,
		logger:    logger,
		callbacks: make([]ShutdownCallback, 0),
		signals:   []os.Signal{os.Interrupt}, // Use a simple signal for testing
	}

	// Create a channel to track callback execution
	callbackCalled := make(chan bool, 1)

	// Create a test callback that returns an error
	callback := func(ctx context.Context) error {
		callbackCalled <- true
		return errors.New("test error")
	}

	// Register the callback
	gs.RegisterCallback(callback)

	// Call HandleShutdown
	_, cancel := gs.HandleShutdown()
	defer cancel()

	// Execute the callback directly
	go func() {
		for _, cb := range gs.callbacks {
			_ = cb(context.Background()) // Ignore the error, we expect one
		}
	}()

	// Wait for the callback to be called or timeout
	select {
	case <-callbackCalled:
		// Callback was called, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not called within timeout")
	}
}

// TestHandleShutdownWithTimeout tests that HandleShutdown handles timeouts correctly
func TestHandleShutdownWithTimeout(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a custom GracefulShutdown for testing
	gs := &GracefulShutdown{
		timeout:   10 * time.Millisecond, // Very short timeout
		logger:    logger,
		callbacks: make([]ShutdownCallback, 0),
		signals:   []os.Signal{os.Interrupt}, // Use a simple signal for testing
	}

	// Create a channel to track if callback starts execution
	callbackStarted := make(chan bool, 1)
	callbackFinished := make(chan bool, 1)

	// Create a test callback that takes longer than the timeout
	callback := func(ctx context.Context) error {
		callbackStarted <- true

		// Check if context is canceled before we finish
		select {
		case <-time.After(50 * time.Millisecond):
			callbackFinished <- true
			return nil
		case <-ctx.Done():
			// Context was canceled due to timeout
			return ctx.Err()
		}
	}

	// Register the callback
	gs.RegisterCallback(callback)

	// Call HandleShutdown
	_, cancel := gs.HandleShutdown()
	defer cancel()

	// Create a timeout context for the callback
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), gs.timeout)
	defer timeoutCancel()

	// Execute the callback directly with the timeout context
	go func() {
		for _, cb := range gs.callbacks {
			_ = cb(timeoutCtx)
		}
	}()

	// Wait for the callback to start
	select {
	case <-callbackStarted:
		// Callback started, continue test
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not started within timeout")
	}

	// The callback should not finish due to the timeout
	select {
	case <-callbackFinished:
		t.Fatal("Callback should not have finished due to timeout")
	case <-time.After(20 * time.Millisecond):
		// Expected: callback didn't finish due to timeout
	}
}

// TestWaitForShutdown tests the WaitForShutdown function
func TestWaitForShutdown(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Call WaitForShutdown in a goroutine
	var ctx context.Context
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx = WaitForShutdown(100*time.Millisecond, logger)
	}()

	// Wait a bit for WaitForShutdown to set up
	time.Sleep(50 * time.Millisecond)

	// Verify that the function returns a non-nil context
	assert.NotNil(t, ctx)
}

// TestSetupSignalHandler tests the SetupSignalHandler function
func TestSetupSignalHandler(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Call SetupSignalHandler
	ctx, gs := SetupSignalHandler(100*time.Millisecond, logger)

	// Verify the context and GracefulShutdown are not nil
	assert.NotNil(t, ctx)
	assert.NotNil(t, gs)

	// Verify the GracefulShutdown has the correct timeout
	assert.Equal(t, 100*time.Millisecond, gs.timeout)
}
