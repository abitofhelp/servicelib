// Package signal provides utilities for handling OS signals and graceful shutdown.
package signal

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// TestableGracefulShutdown is a version of GracefulShutdown that can be used for testing
type TestableGracefulShutdown struct {
	timeout    time.Duration
	logger     *logging.ContextLogger
	callbacks  []ShutdownCallback
	callbackMu sync.Mutex
	signals    []os.Signal
	sigCh      chan os.Signal
	done       chan struct{}
}

// NewTestableGracefulShutdown creates a new TestableGracefulShutdown
func NewTestableGracefulShutdown(timeout time.Duration, logger *logging.ContextLogger) *TestableGracefulShutdown {
	return &TestableGracefulShutdown{
		timeout:   timeout,
		logger:    logger,
		callbacks: make([]ShutdownCallback, 0),
		signals:   []os.Signal{os.Interrupt},
		sigCh:     make(chan os.Signal, 1),
		done:      make(chan struct{}),
	}
}

// Close closes the TestableGracefulShutdown and cleans up resources
func (tgs *TestableGracefulShutdown) Close() {
	close(tgs.done)
}

// RegisterCallback registers a callback to be called during shutdown
func (tgs *TestableGracefulShutdown) RegisterCallback(callback ShutdownCallback) {
	tgs.callbackMu.Lock()
	defer tgs.callbackMu.Unlock()
	tgs.callbacks = append(tgs.callbacks, callback)
}

// HandleShutdownForTest is a version of HandleShutdown that uses a test-specific signal channel
func (tgs *TestableGracefulShutdown) HandleShutdownForTest() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Wait for a signal or done
		select {
		case sig := <-tgs.sigCh:
			tgs.logger.Info(ctx, "Received shutdown signal", zap.String("signal", sig.String()))
		case <-tgs.done:
			// Test is done, exit goroutine
			return
		}

		// Cancel the context to notify all services to shut down
		cancel()

		// Create a timeout context for shutdown
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), tgs.timeout)
		defer timeoutCancel()

		// Execute registered callbacks
		tgs.callbackMu.Lock()
		callbacks := make([]ShutdownCallback, len(tgs.callbacks))
		copy(callbacks, tgs.callbacks)
		tgs.callbackMu.Unlock()

		var wg sync.WaitGroup
		for i, callback := range callbacks {
			wg.Add(1)
			go func(i int, cb ShutdownCallback) {
				defer wg.Done()
				tgs.logger.Info(ctx, "Executing shutdown callback", zap.Int("callback_index", i))
				if err := cb(timeoutCtx); err != nil {
					tgs.logger.Error(ctx, "Shutdown callback failed",
						zap.Int("callback_index", i),
						zap.Error(err))
				} else {
					tgs.logger.Info(ctx, "Shutdown callback completed successfully",
						zap.Int("callback_index", i))
				}
			}(i, callback)
		}

		// Wait for callbacks to complete or timeout
		callbacksDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(callbacksDone)
		}()

		// Wait for callbacks to complete, timeout, second signal, or done
		select {
		case <-callbacksDone:
			tgs.logger.Info(ctx, "All shutdown callbacks completed successfully")
		case <-timeoutCtx.Done():
			tgs.logger.Warn(ctx, "Graceful shutdown timed out, forcing exit")
		case sig := <-tgs.sigCh:
			tgs.logger.Warn(ctx, "Received second signal, forcing exit",
				zap.String("signal", sig.String()))
		case <-tgs.done:
			// Test is done, exit goroutine
			return
		}
	}()

	return ctx, cancel
}

// SendSignal sends a signal to the TestableGracefulShutdown
func (tgs *TestableGracefulShutdown) SendSignal(sig os.Signal) {
	tgs.sigCh <- sig
}
