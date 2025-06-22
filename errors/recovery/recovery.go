// Copyright (c) 2025 A Bit of Help, Inc.

package recovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

var (
	// ErrCircuitBreakerOpen is returned when the circuit breaker is in the open state
	ErrCircuitBreakerOpen = errors.New(errors.InternalErrorCode, "circuit breaker is open")
	// ErrServiceUnavailable is returned when a service is not available
	ErrServiceUnavailable = errors.New(errors.NetworkErrorCode, "service unavailable")
	// ErrContextCanceled is returned when the context is canceled
	ErrContextCanceled = errors.New(errors.CanceledCode, "context canceled")
	// ErrContextDeadlineExceeded is returned when the context deadline is exceeded
	ErrContextDeadlineExceeded = errors.New(errors.TimeoutCode, "context deadline exceeded")
)

// RecoveryHandler handles panic recovery and error rate limiting
type RecoveryHandler struct {
	logger      *logging.ContextLogger
	rateLimiter *time.Ticker
}

// NewRecoveryHandler creates a new recovery handler with error rate limiting
func NewRecoveryHandler(logger *logging.ContextLogger, errorsPerSecond int) *RecoveryHandler {
	return &RecoveryHandler{
		logger:      logger,
		rateLimiter: time.NewTicker(time.Second / time.Duration(errorsPerSecond)),
	}
}

// WithRecovery wraps an operation with panic recovery and error rate limiting
func (h *RecoveryHandler) WithRecovery(ctx context.Context, operation string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			select {
			case <-h.rateLimiter.C:
				h.logger.Error(ctx, "panic recovered",
					zap.String("operation", operation),
					zap.Any("panic", r))
			default:
				// Rate limit exceeded, skip logging
			}
			err = errors.New(errors.InternalErrorCode, fmt.Sprintf("panic in %s: %v", operation, r))
		}
	}()

	err = fn()
	if err != nil {
		select {
		case <-h.rateLimiter.C:
			h.logger.Error(ctx, "Operation failed",
				zap.String("operation", operation),
				zap.Error(err))
		default:
			// Rate limit exceeded, skip logging
		}
	}
	return err
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	failureThreshold int
	cooldownPeriod   time.Duration
	failures         int
	openUntil        time.Time
	mu               sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold int, cooldownPeriod time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		cooldownPeriod:   cooldownPeriod,
	}
}

// Execute executes an operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, operation func() error) (err error) {
	return cb.ExecuteWithName(ctx, "circuit-breaker", operation)
}

// ExecuteWithName executes an operation with circuit breaker protection and a specific operation name
func (cb *CircuitBreaker) ExecuteWithName(ctx context.Context, operationName string, operation func() error) (err error) {
	// Check if context is canceled or deadline exceeded
	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.Canceled:
			return ErrContextCanceled
		case context.DeadlineExceeded:
			return ErrContextDeadlineExceeded
		default:
			return ctx.Err()
		}
	default:
		// Continue with execution
	}

	if cb.isOpen() {
		return ErrCircuitBreakerOpen
	}

	defer func() {
		if r := recover(); r != nil {
			cb.recordFailure()
			err = fmt.Errorf("panic in %s: %v", operationName, r)
		}
	}()

	err = operation()
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// isOpen checks if the circuit breaker is in the open state
func (cb *CircuitBreaker) isOpen() bool {
	cb.mu.RLock()

	// Check if the circuit breaker is not open
	if cb.openUntil.IsZero() {
		cb.mu.RUnlock()
		return false
	}

	// Check if the circuit breaker is still open
	if time.Now().Before(cb.openUntil) {
		cb.mu.RUnlock()
		return true
	}

	// Cooldown period has passed, reset the circuit breaker
	cb.mu.RUnlock()

	// Acquire a write lock to reset the circuit breaker
	cb.mu.Lock()
	cb.openUntil = time.Time{}
	cb.failures = 0
	cb.mu.Unlock()

	return false
}

// recordFailure records a failure and potentially opens the circuit
func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	if cb.failures >= cb.failureThreshold {
		cb.openUntil = time.Now().Add(cb.cooldownPeriod)
	}
}

// recordSuccess resets the failure counter
func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.openUntil = time.Time{}
}
