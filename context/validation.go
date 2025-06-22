// Copyright (c) 2025 A Bit of Help, Inc.

package context

import (
	"context"
	"time"

	"github.com/abitofhelp/servicelib/errors"
)

// ValidationOptions defines options for context validation
type ValidationOptions struct {
	RequireRequestID     bool
	RequireTraceID      bool
	RequireUserID       bool
	RequireTenantID     bool
	RequireOperation    bool
	RequireCorrelationID bool
	MaxTimeout          time.Duration
	MinTimeout          time.Duration
}

// DefaultValidationOptions returns default validation options
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		RequireRequestID:     true,
		RequireTraceID:      true,
		RequireCorrelationID: true,
		MaxTimeout:          24 * time.Hour,
		MinTimeout:          100 * time.Millisecond,
	}
}

// ValidateContext validates a context based on the provided options
func ValidateContext(ctx context.Context, opts ValidationOptions) error {
	if ctx == nil {
		return errors.InvalidInput("context cannot be nil")
	}

	// Check required fields
	if opts.RequireRequestID && GetRequestID(ctx) == "" {
		return errors.InvalidInput("missing request ID in context")
	}

	if opts.RequireTraceID && GetTraceID(ctx) == "" {
		return errors.InvalidInput("missing trace ID in context")
	}

	if opts.RequireUserID && GetUserID(ctx) == "" {
		return errors.InvalidInput("missing user ID in context")
	}

	if opts.RequireTenantID && GetTenantID(ctx) == "" {
		return errors.InvalidInput("missing tenant ID in context")
	}

	if opts.RequireOperation && GetOperation(ctx) == "" {
		return errors.InvalidInput("missing operation in context")
	}

	if opts.RequireCorrelationID && GetCorrelationID(ctx) == "" {
		return errors.InvalidInput("missing correlation ID in context")
	}

	// Validate deadline if present
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		if timeout < opts.MinTimeout {
			return errors.InvalidInput("context deadline is too short: %v < %v", timeout, opts.MinTimeout)
		}
		if timeout > opts.MaxTimeout {
			return errors.InvalidInput("context deadline is too long: %v > %v", timeout, opts.MaxTimeout)
		}
	}

	return nil
}

// ValidateContextInheritance validates that a child context properly inherits from its parent
func ValidateContextInheritance(parent, child context.Context) error {
	if parent == nil || child == nil {
		return errors.InvalidInput("both parent and child contexts must be non-nil")
	}

	// Check if child inherits deadline from parent
	parentDeadline, parentHasDeadline := parent.Deadline()
	childDeadline, childHasDeadline := child.Deadline()

	if parentHasDeadline {
		if !childHasDeadline {
			return errors.InvalidInput("child context must inherit deadline from parent")
		}
		if childDeadline.After(parentDeadline) {
			return errors.InvalidInput("child context deadline cannot be later than parent deadline")
		}
	}

	// Check if child inherits cancellation from parent
	select {
	case <-parent.Done():
		select {
		case <-child.Done():
			// Child is also done, which is correct
		default:
			return errors.InvalidInput("child context must be cancelled if parent is cancelled")
		}
	default:
		// Parent not cancelled, child can be in either state
	}

	// Check value inheritance
	parentReqID := GetRequestID(parent)
	childReqID := GetRequestID(child)
	if parentReqID != "" && childReqID != parentReqID {
		return errors.InvalidInput("child context must inherit request ID from parent")
	}

	parentTraceID := GetTraceID(parent)
	childTraceID := GetTraceID(child)
	if parentTraceID != "" && childTraceID != parentTraceID {
		return errors.InvalidInput("child context must inherit trace ID from parent")
	}

	return nil
}

// MustValidateContext validates a context and panics if validation fails
func MustValidateContext(ctx context.Context, opts ValidationOptions) {
	if err := ValidateContext(ctx, opts); err != nil {
		panic(err)
	}
}
