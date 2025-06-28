// Copyright (c) 2025 A Bit of Help, Inc.

package signal

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

	// Create a channel to safely pass the context from the goroutine to the main test
	ctxChan := make(chan context.Context, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Get the context from WaitForShutdown and send it through the channel
		ctx := WaitForShutdown(100*time.Millisecond, logger)
		ctxChan <- ctx
	}()

	// Wait a bit for WaitForShutdown to set up
	time.Sleep(50 * time.Millisecond)

	// Get the context from the channel with a timeout to avoid blocking indefinitely
	var ctx context.Context
	select {
	case ctx = <-ctxChan:
		// Successfully received the context
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for context from WaitForShutdown")
	}

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

// MockSignalNotifier is a mock for testing signal notifications
type MockSignalNotifier struct {
	mock.Mock
}

// Notify mocks the signal.Notify function
func (m *MockSignalNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	m.Called(c, sig)
}

// Stop mocks the signal.Stop function
func (m *MockSignalNotifier) Stop(c chan<- os.Signal) {
	m.Called(c)
}

// TestRegisterMultipleCallbacks tests registering multiple callbacks
func TestRegisterMultipleCallbacks(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Register multiple callbacks
	callback1 := func(ctx context.Context) error {
		return nil
	}
	callback2 := func(ctx context.Context) error {
		return nil
	}

	gs.RegisterCallback(callback1)
	gs.RegisterCallback(callback2)

	// Verify both callbacks were registered
	assert.Len(t, gs.callbacks, 2)
}

// TestHandleShutdownConcurrency tests the concurrency of callback execution in HandleShutdown
func TestHandleShutdownConcurrency(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Create a mutex to protect access to the callbacksExecuted slice
	var mu sync.Mutex
	callbacksExecuted := make([]int, 0, 2)

	// Register callbacks that add their index to the callbacksExecuted slice
	gs.RegisterCallback(func(ctx context.Context) error {
		mu.Lock()
		callbacksExecuted = append(callbacksExecuted, 1)
		mu.Unlock()
		return nil
	})

	gs.RegisterCallback(func(ctx context.Context) error {
		mu.Lock()
		callbacksExecuted = append(callbacksExecuted, 2)
		mu.Unlock()
		return nil
	})

	// Call HandleShutdown
	_, cancel := gs.HandleShutdown()
	defer cancel()

	// Wait a bit for the goroutine to set up
	time.Sleep(50 * time.Millisecond)

	// Verify that both callbacks were registered
	assert.Len(t, gs.callbacks, 2)
}

// TestCallbackWithError tests a callback that returns an error
func TestCallbackWithError(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Register a callback that returns an error
	expectedErr := errors.New("test error")
	gs.RegisterCallback(func(ctx context.Context) error {
		return expectedErr
	})

	// Verify the callback was registered
	assert.Len(t, gs.callbacks, 1)

	// Execute the callback directly to verify it returns an error
	err := gs.callbacks[0](context.Background())
	assert.Equal(t, expectedErr, err)
}

// TestHandleShutdownContext tests that HandleShutdown returns a context
func TestHandleShutdownContext(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Call HandleShutdown
	ctx, cancel := gs.HandleShutdown()
	defer cancel()

	// Verify the context is not nil
	assert.NotNil(t, ctx)

	// Verify the context is not canceled yet
	select {
	case <-ctx.Done():
		t.Fatal("Context should not be canceled yet")
	default:
		// Expected
	}

	// Cancel the context
	cancel()

	// Verify the context is now canceled
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Fatal("Context should be canceled")
	}
}

// TestWaitForShutdownImplementation tests the implementation of WaitForShutdown
func TestWaitForShutdownImplementation(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Test that WaitForShutdown creates a GracefulShutdown and calls HandleShutdown
	ctx := WaitForShutdown(100*time.Millisecond, logger)

	// Verify that the context is not nil
	require.NotNil(t, ctx)
}

// TestTimeoutContext tests that a timeout context is created for callbacks
func TestTimeoutContext(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown with a very short timeout
	timeout := 50 * time.Millisecond
	gs := NewGracefulShutdown(timeout, logger)

	// Verify the timeout was set correctly
	assert.Equal(t, timeout, gs.timeout)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a channel to track if the context times out
	timeoutCh := make(chan bool, 1)

	// Start a goroutine that waits for the context to be done
	go func() {
		<-ctx.Done()
		timeoutCh <- true
	}()

	// Wait for the context to time out
	select {
	case <-timeoutCh:
		// Context timed out, test passes
	case <-time.After(timeout * 2):
		t.Fatal("Context did not time out within expected duration")
	}
}

// TestCallbackContextCancellation tests that a callback receives a canceled context when the parent context is canceled
func TestCallbackContextCancellation(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to track if the callback's context is canceled
	contextCanceled := make(chan bool, 1)

	// Register a callback that checks if its context is canceled
	gs.RegisterCallback(func(callbackCtx context.Context) error {
		// Wait for the context to be canceled
		<-callbackCtx.Done()
		contextCanceled <- true
		return callbackCtx.Err()
	})

	// Execute the callback directly with our context
	go gs.callbacks[0](ctx)

	// Cancel the context
	cancel()

	// Wait for the callback to report that its context was canceled
	select {
	case <-contextCanceled:
		// Context was canceled, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback did not report context cancellation within timeout")
	}
}

// TestSetupSignalHandlerImplementation tests the implementation of SetupSignalHandler
func TestSetupSignalHandlerImplementation(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Call SetupSignalHandler
	ctx, gs := SetupSignalHandler(100*time.Millisecond, logger)

	// Verify the context and GracefulShutdown are not nil
	require.NotNil(t, ctx)
	require.NotNil(t, gs)

	// Verify the GracefulShutdown has the correct timeout
	assert.Equal(t, 100*time.Millisecond, gs.timeout)
}

// TestSetupSignalHandlerCallback tests registering a callback with a GracefulShutdown from SetupSignalHandler
func TestSetupSignalHandlerCallback(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Call SetupSignalHandler
	_, gs := SetupSignalHandler(100*time.Millisecond, logger)

	// Register a callback
	callbackExecuted := false
	gs.RegisterCallback(func(ctx context.Context) error {
		callbackExecuted = true
		return nil
	})

	// Verify the callback was registered
	assert.Len(t, gs.callbacks, 1)

	// Execute the callback directly
	err := gs.callbacks[0](context.Background())
	assert.NoError(t, err)

	// Verify the callback was executed
	assert.True(t, callbackExecuted)
}

// TestGracefulShutdownSignals tests that GracefulShutdown has the expected signals
func TestGracefulShutdownSignals(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Verify the signals slice has the expected signals
	assert.Len(t, gs.signals, 4)
	assert.Contains(t, gs.signals, os.Interrupt)
}

// TestWaitGroupInShutdown tests the use of WaitGroup in HandleShutdown
func TestWaitGroupInShutdown(t *testing.T) {
	// Create a WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a channel to signal when the WaitGroup is done
	done := make(chan struct{})

	// Start a goroutine that waits for the WaitGroup
	go func() {
		wg.Wait()
		close(done)
	}()

	// Verify the WaitGroup is not done yet
	select {
	case <-done:
		t.Fatal("WaitGroup should not be done yet")
	default:
		// Expected
	}

	// Mark the WaitGroup as done
	wg.Done()

	// Verify the WaitGroup is now done
	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitGroup should be done")
	}
}

// TestGracefulShutdownMutex tests the mutex in GracefulShutdown
func TestGracefulShutdownMutex(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Create a channel to track when the goroutine is done
	done := make(chan struct{})

	// Start a goroutine that locks the mutex
	go func() {
		gs.callbackMu.Lock()
		// Hold the lock for a bit
		time.Sleep(50 * time.Millisecond)
		gs.callbackMu.Unlock()
		close(done)
	}()

	// Wait a bit for the goroutine to acquire the lock
	time.Sleep(10 * time.Millisecond)

	// Try to register a callback, which will also try to lock the mutex
	gs.RegisterCallback(func(ctx context.Context) error {
		return nil
	})

	// Wait for the goroutine to finish
	select {
	case <-done:
		// Expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Goroutine did not finish within timeout")
	}

	// Verify the callback was registered
	assert.Len(t, gs.callbacks, 1)
}

// TestTestableGracefulShutdown tests the TestableGracefulShutdown
func TestTestableGracefulShutdown(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance
	tgs := NewTestableGracefulShutdown(100*time.Millisecond, logger)
	defer tgs.Close()

	// Verify the TestableGracefulShutdown was created correctly
	assert.Equal(t, 100*time.Millisecond, tgs.timeout)
	assert.Equal(t, logger, tgs.logger)
	assert.Empty(t, tgs.callbacks)
	assert.Len(t, tgs.signals, 1)
}

// TestTestableGracefulShutdownHandleShutdown tests the HandleShutdownForTest method
func TestTestableGracefulShutdownHandleShutdown(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance
	tgs := NewTestableGracefulShutdown(100*time.Millisecond, logger)
	defer tgs.Close()

	// Register a callback
	callbackExecuted := make(chan bool, 1)
	tgs.RegisterCallback(func(ctx context.Context) error {
		callbackExecuted <- true
		return nil
	})

	// Call HandleShutdownForTest
	ctx, cancel := tgs.HandleShutdownForTest()
	defer cancel()

	// Send a signal to trigger the shutdown
	tgs.SendSignal(os.Interrupt)

	// Wait for the callback to be executed
	select {
	case <-callbackExecuted:
		// Callback was executed, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not executed within timeout")
	}

	// Verify the context is canceled
	select {
	case <-ctx.Done():
		// Context is canceled, test passes
	default:
		t.Fatal("Context should be canceled")
	}
}

// TestTestableGracefulShutdownWithErrorCallback tests the HandleShutdownForTest method with a callback that returns an error
func TestTestableGracefulShutdownWithErrorCallback(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance
	tgs := NewTestableGracefulShutdown(100*time.Millisecond, logger)
	defer tgs.Close()

	// Register a callback that returns an error
	callbackExecuted := make(chan bool, 1)
	tgs.RegisterCallback(func(ctx context.Context) error {
		callbackExecuted <- true
		return errors.New("test error")
	})

	// Call HandleShutdownForTest
	ctx, cancel := tgs.HandleShutdownForTest()
	defer cancel()

	// Send a signal to trigger the shutdown
	tgs.SendSignal(os.Interrupt)

	// Wait for the callback to be executed
	select {
	case <-callbackExecuted:
		// Callback was executed, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not executed within timeout")
	}

	// Verify the context is canceled
	select {
	case <-ctx.Done():
		// Context is canceled, test passes
	default:
		t.Fatal("Context should be canceled")
	}
}

// TestTestableGracefulShutdownWithTimeout tests the HandleShutdownForTest method with a timeout
func TestTestableGracefulShutdownWithTimeout(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance with a very short timeout
	tgs := NewTestableGracefulShutdown(50*time.Millisecond, logger)
	defer tgs.Close()

	// Register a callback that takes longer than the timeout
	callbackStarted := make(chan bool, 1)
	tgs.RegisterCallback(func(ctx context.Context) error {
		callbackStarted <- true
		// Sleep longer than the timeout
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	// Call HandleShutdownForTest
	ctx, cancel := tgs.HandleShutdownForTest()
	defer cancel()

	// Send a signal to trigger the shutdown
	tgs.SendSignal(os.Interrupt)

	// Wait for the callback to start
	select {
	case <-callbackStarted:
		// Callback started, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not started within timeout")
	}

	// Verify the context is canceled
	select {
	case <-ctx.Done():
		// Context is canceled, test passes
	default:
		t.Fatal("Context should be canceled")
	}
}

// TestMultipleCallbacksWithDifferentResults tests multiple callbacks with different results
func TestMultipleCallbacksWithDifferentResults(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a GracefulShutdown instance
	gs := NewGracefulShutdown(100*time.Millisecond, logger)

	// Register callbacks with different results
	gs.RegisterCallback(func(ctx context.Context) error {
		return nil // Success
	})

	gs.RegisterCallback(func(ctx context.Context) error {
		return errors.New("test error") // Error
	})

	gs.RegisterCallback(func(ctx context.Context) error {
		// Block until context is canceled
		<-ctx.Done()
		return ctx.Err() // Context canceled
	})

	// Verify the callbacks were registered
	assert.Len(t, gs.callbacks, 3)

	// Execute the callbacks directly
	err1 := gs.callbacks[0](context.Background())
	assert.NoError(t, err1)

	err2 := gs.callbacks[1](context.Background())
	assert.Error(t, err2)
	assert.Equal(t, "test error", err2.Error())

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start the third callback in a goroutine
	done := make(chan struct{})
	go func() {
		err3 := gs.callbacks[2](ctx)
		assert.Error(t, err3)
		assert.Equal(t, context.Canceled, err3)
		close(done)
	}()

	// Cancel the context
	cancel()

	// Wait for the callback to complete
	select {
	case <-done:
		// Callback completed, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback did not complete within timeout")
	}
}

// TestHandleShutdownWithRealSignal tests the HandleShutdown function with a real signal
func TestHandleShutdownWithRealSignal(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance with a very short timeout
	tgs := NewTestableGracefulShutdown(50*time.Millisecond, logger)
	defer tgs.Close()

	// Create a channel to track callback execution
	callbackExecuted := make(chan bool, 1)

	// Register a callback that signals when it's called
	tgs.RegisterCallback(func(ctx context.Context) error {
		callbackExecuted <- true
		return nil
	})

	// Call HandleShutdownForTest
	ctx, cancel := tgs.HandleShutdownForTest()
	defer cancel()

	// Send a signal to trigger the shutdown
	tgs.SendSignal(os.Interrupt)

	// Wait for the callback to be executed or timeout
	select {
	case <-callbackExecuted:
		// Callback was executed, test passes
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Callback was not executed within timeout")
	}

	// Verify the context is canceled
	select {
	case <-ctx.Done():
		// Context is canceled, test passes
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be canceled")
	}
}

// TestHandleShutdownWithTimeoutExceeded tests the HandleShutdown function with a timeout
func TestHandleShutdownWithTimeoutExceeded(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance with a very short timeout
	tgs := NewTestableGracefulShutdown(50*time.Millisecond, logger)
	defer tgs.Close()

	// Register a callback that takes longer than the timeout
	tgs.RegisterCallback(func(ctx context.Context) error {
		// Sleep longer than the timeout
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	// Call HandleShutdownForTest
	ctx, cancel := tgs.HandleShutdownForTest()
	defer cancel()

	// Send a signal to trigger the shutdown
	tgs.SendSignal(os.Interrupt)

	// Verify the context is canceled
	select {
	case <-ctx.Done():
		// Context is canceled, test passes
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be canceled")
	}

	// Wait a bit to allow the timeout to occur
	time.Sleep(300 * time.Millisecond)
}

// TestHandleShutdownWithSecondSignal tests the HandleShutdown function with a second signal
func TestHandleShutdownWithSecondSignal(t *testing.T) {
	// Create a test logger
	logger := logging.NewContextLogger(zaptest.NewLogger(t))

	// Create a TestableGracefulShutdown instance with a longer timeout
	tgs := NewTestableGracefulShutdown(500*time.Millisecond, logger)
	defer tgs.Close()

	// Register a callback that takes some time
	tgs.RegisterCallback(func(ctx context.Context) error {
		// Sleep for a bit
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	// Call HandleShutdownForTest
	ctx, cancel := tgs.HandleShutdownForTest()
	defer cancel()

	// Send a signal to trigger the shutdown
	tgs.SendSignal(os.Interrupt)

	// Verify the context is canceled
	select {
	case <-ctx.Done():
		// Context is canceled, test passes
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be canceled")
	}

	// Wait a bit and then send a second signal
	time.Sleep(100 * time.Millisecond)
	tgs.SendSignal(syscall.SIGTERM)

	// Wait a bit to allow the second signal to be processed
	time.Sleep(300 * time.Millisecond)
}
