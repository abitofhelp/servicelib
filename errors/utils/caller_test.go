// Copyright (c) 2025 A Bit of Help, Inc.

package utils

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallerInfo(t *testing.T) {
	// Since getCallerInfo is not exported, we'll test it indirectly through a wrapper function
	file, line := getCallerInfoWrapper()

	// Get the expected file using runtime.Caller directly
	_, expectedFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "runtime.Caller should succeed")

	// The file should be the base name of the current file
	assert.Equal(t, filepath.Base(expectedFile), file, "File name should match")

	// The line number will be different because we're calling through a wrapper
	// but it should be a positive number
	assert.Greater(t, line, 0, "Line number should be positive")
}

func getCallerInfoWrapper() (string, int) {
	// Call getCallerInfo with skip=0 to get info about the caller of this function
	return getCallerInfo(0)
}

func TestGetCallerPackage(t *testing.T) {
	// Get the package name of the caller (this test function)
	pkg := GetCallerPackage(0)

	// The package should contain "errors.utils" since this test is in that package
	assert.Contains(t, pkg, "errors.utils", "Package name should contain errors.utils")

	// Test with a higher skip value (which should still give a valid package)
	higherSkipPkg := GetCallerPackage(1)
	assert.NotEmpty(t, higherSkipPkg, "Package name should not be empty with higher skip")

	// Test with a very high skip value (which should return "unknown")
	veryHighSkipPkg := GetCallerPackage(100)
	assert.Equal(t, "unknown", veryHighSkipPkg, "Very high skip should return 'unknown'")
}

func TestGetCallerPackageWithGithubPath(t *testing.T) {
	// This test simulates a GitHub path by creating a wrapper function
	// that manipulates the result of runtime.Caller

	// Define a test case with a GitHub path
	testCases := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "GitHub path",
			filePath: "/home/user/go/src/github.com/abitofhelp/servicelib/errors/utils/caller.go",
			expected: "abitofhelp.servicelib.errors.utils",
		},
		{
			name:     "Non-GitHub path",
			filePath: "/home/user/projects/myapp/utils/caller.go",
			expected: "home.user.projects.myapp.utils",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock function that returns a fixed file path
			originalRuntimeCaller := runtimeCaller
			defer func() { runtimeCaller = originalRuntimeCaller }()

			runtimeCaller = func(skip int) (uintptr, string, int, bool) {
				return 0, tc.filePath, 0, true
			}

			// Call GetCallerPackage with our mocked runtime.Caller
			result := GetCallerPackage(0)

			// Check the result
			assert.Equal(t, tc.expected, result, "Package name should be correctly extracted")
		})
	}
}

