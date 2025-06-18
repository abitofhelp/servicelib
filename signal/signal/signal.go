// Copyright (c) 2025 A Bit of Help, Inc.

// Package signal provides utilities for testing signal handling.
package signal

import (
	"context"
	"os"
	"sync"
	"time"
)

// SignalHandler is an interface for handling signals.
type SignalHandler interface {
	// HandleSignal handles a signal.
	HandleSignal(sig os.Signal) error
}

// DefaultSignalHandler is a default implementation of SignalHandler.
type DefaultSignalHandler struct {
	mu      sync.Mutex
	handled map[os.Signal]bool
}

// NewDefaultSignalHandler creates a new DefaultSignalHandler.
func NewDefaultSignalHandler() *DefaultSignalHandler {
	return &DefaultSignalHandler{
		handled: make(map[os.Signal]bool),
	}
}

// HandleSignal handles a signal.
func (h *DefaultSignalHandler) HandleSignal(sig os.Signal) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handled[sig] = true
	return nil
}

// IsHandled returns true if the signal has been handled.
func (h *DefaultSignalHandler) IsHandled(sig os.Signal) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.handled[sig]
}

// SignalWaiter waits for signals.
type SignalWaiter struct {
	sigCh chan os.Signal
}

// NewSignalWaiter creates a new SignalWaiter.
func NewSignalWaiter() *SignalWaiter {
	return &SignalWaiter{
		sigCh: make(chan os.Signal, 1),
	}
}

// WaitForSignal waits for a signal or until the context is canceled.
func (w *SignalWaiter) WaitForSignal(ctx context.Context, timeout time.Duration) (os.Signal, error) {
	select {
	case sig := <-w.sigCh:
		return sig, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, context.DeadlineExceeded
	}
}

// SendSignal sends a signal to the SignalWaiter.
func (w *SignalWaiter) SendSignal(sig os.Signal) {
	w.sigCh <- sig
}
