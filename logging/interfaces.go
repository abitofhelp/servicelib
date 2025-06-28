// Copyright (c) 2025 A Bit of Help, Inc.

// Package logging provides centralized logging functionality for services.
package logging

import (
	"github.com/abitofhelp/servicelib/logging/interfaces"
)

// Logger is an alias for interfaces.Logger for backward compatibility
type Logger = interfaces.Logger

// Ensure ContextLogger implements Logger interface
var _ Logger = (*ContextLogger)(nil)
