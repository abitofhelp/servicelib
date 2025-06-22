// Copyright (c) 2025 A Bit of Help, Inc.

package shutdown

import (
	"context"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/health"
)

// ShutdownHealthCheck monitors shutdown health
type ShutdownHealthCheck struct {
	ps           *PrioritizedShutdown
	healthStatus health.Status
	mu           sync.RWMutex
}

// NewShutdownHealthCheck creates a new shutdown health check
func NewShutdownHealthCheck(ps *PrioritizedShutdown) *ShutdownHealthCheck {
	return &ShutdownHealthCheck{
		ps:           ps,
		healthStatus: health.StatusUp,
	}
}

// Check implements health.Check interface
func (hc *ShutdownHealthCheck) Check(ctx context.Context) health.Result {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	status := hc.ps.GetStatus()
	details := map[string]interface{}{
		"state":         status.State,
		"total_ops":     status.TotalOps,
		"completed_ops": status.CompletedOps,
		"failed_ops":    status.FailedOps,
	}

	if status.StartTime != (time.Time{}) {
		details["start_time"] = status.StartTime
		if status.EndTime != (time.Time{}) {
			details["end_time"] = status.EndTime
			details["duration"] = status.EndTime.Sub(status.StartTime).String()
		}
	}

	if len(status.Errors) > 0 {
		var errMsgs []string
		for _, err := range status.Errors {
			errMsgs = append(errMsgs, err.Error())
		}
		details["errors"] = errMsgs
	}

	return health.Result{
		Name:    "shutdown",
		Status:  hc.healthStatus,
		Details: details,
	}
}

// Name implements health.Check interface
func (hc *ShutdownHealthCheck) Name() string {
	return "shutdown"
}

// UpdateStatus updates the health check status
func (hc *ShutdownHealthCheck) UpdateStatus(status health.Status) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.healthStatus = status
}

// MarkShuttingDown marks the service as shutting down in the health check
func (hc *ShutdownHealthCheck) MarkShuttingDown() {
	hc.UpdateStatus(health.StatusShuttingDown)
}
