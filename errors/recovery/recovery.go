// Copyright (c) 2025 A Bit of Help, Inc.

// Package recovery provides error recovery and rate limiting mechanisms.
package recovery

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"golang.org/x/time/rate"
)

// RecoveryHandler handles panic recovery and error tracking
type RecoveryHandler struct {
	logger    *logging.ContextLogger
	limiter   *rate.Limiter
	errorsMu  sync.RWMutex
	errorRate map[string]*rate.Limiter
}

// NewRecoveryHandler creates a new RecoveryHandler
func NewRecoveryHandler(logger *logging.ContextLogger, maxBurst int) *RecoveryHandler {
	return &RecoveryHandler{
		logger:    logger,
		limiter:   rate.NewLimiter(rate.Limit(maxBurst), maxBurst),
		errorRate: make(map[string]*rate.Limiter),
	}
}

// WithRecovery wraps an operation with panic recovery and error tracking
func (h *RecoveryHandler) WithRecovery(ctx context.Context, operation string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic in %s: %v\n%s", operation, r, stack)
			h.logger.Error(ctx, "Recovered from panic",
				"operation", operation,
				"error", err,
				"stack", string(stack),
			)
		}
	}()

	err = fn()
	if err != nil {
		if !h.shouldLogError(operation) {
			return err
		}
		h.logger.Error(ctx, "Operation failed",
			"operation", operation,
			"error", err,
		)
	}
	return err
}

// shouldLogError implements rate limiting for error logging
func (h *RecoveryHandler) shouldLogError(operation string) bool {
	h.errorsMu.RLock()
	limiter, exists := h.errorRate[operation]
	h.errorsMu.RUnlock()

	if !exists {
		h.errorsMu.Lock()
		limiter = rate.NewLimiter(rate.Every(time.Minute), 10)
		h.errorRate[operation] = limiter
		h.errorsMu.Unlock()
	}

	return limiter.Allow()
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	failures         int
	lastFailure      time.Time
	mu              sync.RWMutex
	state           string
}

// NewCircuitBreaker creates a new CircuitBreaker
func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		state:           "closed",
	}
}

// Execute executes an operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(operation func() error) error {
	if !cb.isAllowed() {
		return errors.ServiceUnavailable("circuit breaker is open")
	}

	err := operation()
	cb.recordResult(err)
	return err
}

// isAllowed checks if the operation should be allowed
func (cb *CircuitBreaker) isAllowed() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == "closed" {
		return true
	}

	if time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.mu.RUnlock()
		cb.mu.Lock()
		cb.state = "half-open"
		cb.mu.Unlock()
		cb.mu.RLock()
		return true
	}

	return false
}

// recordResult records the result of an operation
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		if cb.failures >= cb.failureThreshold {
			cb.state = "open"
		}
	} else {
		cb.failures = 0
		cb.state = "closed"
	}
}
