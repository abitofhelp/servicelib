// Copyright (c) 2025 A Bit of Help, Inc.

// Package signal provides utilities for testing signal handling.
package signal

import (
	"context"
	"os"
	"time"
)

// TestHelper provides utilities for testing signal handling.
type TestHelper struct {
	sigCh chan os.Signal
}

// NewTestHelper creates a new TestHelper.
func NewTestHelper() *TestHelper {
	return &TestHelper{
		sigCh: make(chan os.Signal, 1),
	}
}

// SendSignal sends a signal to the test helper.
func (h *TestHelper) SendSignal(sig os.Signal) {
	h.sigCh <- sig
}

// WaitForSignal waits for a signal or until the context is canceled.
func (h *TestHelper) WaitForSignal(ctx context.Context, timeout time.Duration) (os.Signal, error) {
	select {
	case sig := <-h.sigCh:
		return sig, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, context.DeadlineExceeded
	}
}
