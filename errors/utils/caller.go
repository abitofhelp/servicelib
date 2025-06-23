// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

// Override runtime.Caller for testing
var runtimeCaller = runtime.Caller

// getCallerInfo returns the file name and line number of the caller.
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtimeCaller(skip + 1)
	if !ok {
		return "unknown", 0
	}
	return filepath.Base(file), line
}

// GetCallerPackage returns the package name of the caller.
func GetCallerPackage(skip int) string {
	_, file, _, ok := runtimeCaller(skip + 1)
	if !ok {
		return "unknown"
	}

	pkg := filepath.Dir(file)
	if strings.Contains(pkg, "github.com") {
		// Extract package name after github.com
		parts := strings.Split(pkg, "github.com/")
		if len(parts) > 1 {
			pkg = parts[1]
		}
	}

	// Remove leading slash if present
	pkg = strings.TrimPrefix(pkg, "/")

	// Replace path separators with dots
	return strings.ReplaceAll(pkg, "/", ".")
}
