// Copyright (c) 2025 A Bit of Help, Inc.

package shutdown

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// ShutdownPriority defines priority levels for shutdown operations
type ShutdownPriority int

const (
	// PriorityCritical is for operations that must be shut down first
	PriorityCritical ShutdownPriority = 100
	// PriorityHigh is for important but non-critical operations
	PriorityHigh ShutdownPriority = 75
	// PriorityNormal is for standard operations
	PriorityNormal ShutdownPriority = 50
	// PriorityLow is for operations that can be shut down last
	PriorityLow ShutdownPriority = 25
)

// ShutdownOperation represents a prioritized shutdown operation
type ShutdownOperation struct {
	Name     string
	Priority ShutdownPriority
	Func     func(ctx context.Context) error
	Status   string
	Error    error
	Duration time.Duration
}

// PrioritizedShutdown manages prioritized shutdown operations
type PrioritizedShutdown struct {
	operations []ShutdownOperation
	mu         sync.RWMutex
	logger     *logging.ContextLogger
	status     *ShutdownStatus
}

// ShutdownStatus tracks the overall shutdown progress
type ShutdownStatus struct {
	StartTime    time.Time
	EndTime      time.Time
	State        string
	TotalOps     int
	CompletedOps int
	FailedOps    int
	Errors       []error
	mu           sync.RWMutex
}

// NewPrioritizedShutdown creates a new PrioritizedShutdown
func NewPrioritizedShutdown(logger *logging.ContextLogger) *PrioritizedShutdown {
	return &PrioritizedShutdown{
		logger: logger,
		status: &ShutdownStatus{
			State: "ready",
		},
	}
}

// RegisterOperation registers a shutdown operation with a priority
func (ps *PrioritizedShutdown) RegisterOperation(name string, priority ShutdownPriority, fn func(ctx context.Context) error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.operations = append(ps.operations, ShutdownOperation{
		Name:     name,
		Priority: priority,
		Func:     fn,
		Status:   "pending",
	})

	ps.status.mu.Lock()
	ps.status.TotalOps++
	ps.status.mu.Unlock()
}

// ExecuteShutdown executes all registered shutdown operations in priority order
func (ps *PrioritizedShutdown) ExecuteShutdown(ctx context.Context) error {
	ps.mu.Lock()
	ps.status.StartTime = time.Now()
	ps.status.State = "in_progress"
	
	// Sort operations by priority (highest first)
	sort.Slice(ps.operations, func(i, j int) bool {
		return ps.operations[i].Priority > ps.operations[j].Priority
	})
	ops := ps.operations
	ps.mu.Unlock()

	var lastErr error
	for i := range ops {
		op := &ops[i]
		ps.logger.Info(ctx, "Starting shutdown operation",
			zap.String("operation", op.Name),
			zap.Int("priority", int(op.Priority)))

		start := time.Now()
		op.Status = "running"
		
		if err := op.Func(ctx); err != nil {
			op.Status = "failed"
			op.Error = err
			lastErr = err
			
			ps.status.mu.Lock()
			ps.status.FailedOps++
			ps.status.Errors = append(ps.status.Errors, err)
			ps.status.mu.Unlock()

			ps.logger.Error(ctx, "Shutdown operation failed",
				zap.String("operation", op.Name),
				zap.Error(err))
		} else {
			op.Status = "completed"
			ps.status.mu.Lock()
			ps.status.CompletedOps++
			ps.status.mu.Unlock()
		}

		op.Duration = time.Since(start)
		ps.logger.Info(ctx, "Completed shutdown operation",
			zap.String("operation", op.Name),
			zap.String("status", op.Status),
			zap.Duration("duration", op.Duration))
	}

	ps.status.mu.Lock()
	ps.status.EndTime = time.Now()
	ps.status.State = "completed"
	ps.status.mu.Unlock()

	return lastErr
}

// GetStatus returns the current shutdown status
func (ps *PrioritizedShutdown) GetStatus() ShutdownStatus {
	ps.status.mu.RLock()
	defer ps.status.mu.RUnlock()
	return *ps.status
}

// GetOperationStatus returns the status of a specific operation
func (ps *PrioritizedShutdown) GetOperationStatus(name string) (ShutdownOperation, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, op := range ps.operations {
		if op.Name == name {
			return op, nil
		}
	}

	return ShutdownOperation{}, errors.NotFound("operation %s not found", name)
}
