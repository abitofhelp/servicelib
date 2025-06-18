// Copyright (c) 2025 A Bit of Help, Inc.

// Package signal provides utilities for testing signal handling.
package signal

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSignalHandling tests basic signal handling functionality.
func TestSignalHandling(t *testing.T) {
	// Create a channel to receive signals
	sigCh := make(chan os.Signal, 1)

	// Send a signal
	sigCh <- os.Interrupt

	// Verify we can receive the signal
	select {
	case sig := <-sigCh:
		assert.Equal(t, os.Interrupt, sig)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for signal")
	}
}

// TestDefaultSignalHandler tests the DefaultSignalHandler.
func TestDefaultSignalHandler(t *testing.T) {
	// Create a DefaultSignalHandler
	handler := NewDefaultSignalHandler()

	// Verify the handler is not nil
	assert.NotNil(t, handler)

	// Verify no signals have been handled yet
	assert.False(t, handler.IsHandled(os.Interrupt))

	// Handle a signal
	err := handler.HandleSignal(os.Interrupt)
	assert.NoError(t, err)

	// Verify the signal has been handled
	assert.True(t, handler.IsHandled(os.Interrupt))

	// Verify other signals have not been handled
	assert.False(t, handler.IsHandled(os.Kill))
}

// TestSignalWaiter tests the SignalWaiter.
func TestSignalWaiter(t *testing.T) {
	// Create a SignalWaiter
	waiter := NewSignalWaiter()

	// Verify the waiter is not nil
	assert.NotNil(t, waiter)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Send a signal
	go func() {
		time.Sleep(10 * time.Millisecond)
		waiter.SendSignal(os.Interrupt)
	}()

	// Wait for the signal
	sig, err := waiter.WaitForSignal(ctx, 50*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, os.Interrupt, sig)
}

// TestSignalWaiterTimeout tests the SignalWaiter with a timeout.
func TestSignalWaiterTimeout(t *testing.T) {
	// Create a SignalWaiter
	waiter := NewSignalWaiter()

	// Verify the waiter is not nil
	assert.NotNil(t, waiter)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Wait for a signal with a very short timeout
	sig, err := waiter.WaitForSignal(ctx, 10*time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Nil(t, sig)
}

// TestSignalWaiterContextCancellation tests the SignalWaiter with a canceled context.
func TestSignalWaiterContextCancellation(t *testing.T) {
	// Create a SignalWaiter
	waiter := NewSignalWaiter()

	// Verify the waiter is not nil
	assert.NotNil(t, waiter)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context
	cancel()

	// Wait for a signal
	sig, err := waiter.WaitForSignal(ctx, 50*time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Nil(t, sig)
}

// TestContextCancellation tests that a context can be canceled.
func TestContextCancellation(t *testing.T) {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to track if the context is canceled
	done := make(chan struct{})

	// Start a goroutine that waits for the context to be canceled
	go func() {
		<-ctx.Done()
		close(done)
	}()

	// Cancel the context
	cancel()

	// Verify the context is canceled
	select {
	case <-done:
		// Context was canceled, test passes
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context was not canceled within timeout")
	}
}

// TestTestHelper tests the TestHelper.
func TestTestHelper(t *testing.T) {
	// Create a TestHelper
	helper := NewTestHelper()

	// Verify the helper is not nil
	assert.NotNil(t, helper)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Send a signal
	go func() {
		time.Sleep(10 * time.Millisecond)
		helper.SendSignal(os.Interrupt)
	}()

	// Wait for the signal
	sig, err := helper.WaitForSignal(ctx, 50*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, os.Interrupt, sig)
}

// TestTestHelperTimeout tests the TestHelper with a timeout.
func TestTestHelperTimeout(t *testing.T) {
	// Create a TestHelper
	helper := NewTestHelper()

	// Verify the helper is not nil
	assert.NotNil(t, helper)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Wait for a signal with a very short timeout
	sig, err := helper.WaitForSignal(ctx, 10*time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Nil(t, sig)
}

// TestTestHelperContextCancellation tests the TestHelper with a canceled context.
func TestTestHelperContextCancellation(t *testing.T) {
	// Create a TestHelper
	helper := NewTestHelper()

	// Verify the helper is not nil
	assert.NotNil(t, helper)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context
	cancel()

	// Wait for a signal
	sig, err := helper.WaitForSignal(ctx, 50*time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Nil(t, sig)
}
