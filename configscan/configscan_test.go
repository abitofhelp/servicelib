// Copyright (c) 2025 A Bit of Help, Inc.

package configscan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner_Scan(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "configscan_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	createTestFiles(t, tempDir)

	// Create a scanner with the temporary directory
	scanner := NewScanner(tempDir)

	// Scan packages for configuration requirements
	requirements, err := scanner.Scan()
	assert.NoError(t, err)

	// Verify the requirements
	assert.Len(t, requirements, 3)

	// Find the requirements for each config type
	var configAReq, configBReq, configCReq *ConfigRequirement
	for i, req := range requirements {
		if req.ConfigType == "ConfigA" {
			configAReq = &requirements[i]
		} else if req.ConfigType == "ConfigB" {
			configBReq = &requirements[i]
		} else if req.ConfigType == "ConfigC" {
			configCReq = &requirements[i]
		}
	}

	// Verify ConfigA requirement
	assert.NotNil(t, configAReq)
	assert.Equal(t, "packagea", configAReq.PackageName)
	assert.ElementsMatch(t, []string{"Field1", "Field2"}, configAReq.Fields)
	assert.True(t, configAReq.HasDefault)

	// Verify ConfigB requirement
	assert.NotNil(t, configBReq)
	assert.Equal(t, "packageb", configBReq.PackageName)
	assert.ElementsMatch(t, []string{"Field1", "Field2"}, configBReq.Fields)
	assert.False(t, configBReq.HasDefault)

	// Verify ConfigC requirement
	assert.NotNil(t, configCReq)
	assert.Equal(t, "packagec", configCReq.PackageName)
	assert.ElementsMatch(t, []string{"Field1", "Field2"}, configCReq.Fields)
	assert.True(t, configCReq.HasDefault)
}

func TestReport(t *testing.T) {
	// Create test requirements
	requirements := []ConfigRequirement{
		{
			PackageName: "packagea",
			ConfigType:  "ConfigA",
			Fields:      []string{"Field1", "Field2"},
			HasDefault:  true,
		},
		{
			PackageName: "packageb",
			ConfigType:  "ConfigB",
			Fields:      []string{"Field1", "Field2"},
			HasDefault:  false,
		},
		{
			PackageName: "packagec",
			ConfigType:  "ConfigC",
			Fields:      []string{"Field1", "Field2"},
			HasDefault:  true,
		},
	}

	// Generate a report
	missingDefaults := Report(requirements)

	// Verify the report
	assert.Len(t, missingDefaults, 1)
	assert.Equal(t, "packageb", missingDefaults[0].PackageName)
	assert.Equal(t, "ConfigB", missingDefaults[0].ConfigType)
}

func createTestFiles(t *testing.T, tempDir string) {
	// Create package directories
	packageADir := filepath.Join(tempDir, "packagea")
	packageBDir := filepath.Join(tempDir, "packageb")
	packageCDir := filepath.Join(tempDir, "packagec")
	assert.NoError(t, os.Mkdir(packageADir, 0755))
	assert.NoError(t, os.Mkdir(packageBDir, 0755))
	assert.NoError(t, os.Mkdir(packageCDir, 0755))

	// Create package A files
	packageAFile := filepath.Join(packageADir, "packagea.go")
	packageAContent := `
package packagea

// ConfigA is a configuration struct
type ConfigA struct {
	Field1 string
	Field2 int
}

// DefaultConfigA returns the default configuration
func DefaultConfigA() ConfigA {
	return ConfigA{
		Field1: "default",
		Field2: 42,
	}
}
`
	assert.NoError(t, os.WriteFile(packageAFile, []byte(packageAContent), 0644))

	// Create package B files
	packageBFile := filepath.Join(packageBDir, "packageb.go")
	packageBContent := `
package packageb

// ConfigB is a configuration struct
type ConfigB struct {
	Field1 string
	Field2 int
}
`
	assert.NoError(t, os.WriteFile(packageBFile, []byte(packageBContent), 0644))

	// Create package C files
	packageCFile := filepath.Join(packageCDir, "packagec.go")
	packageCContent := `
package packagec

// ConfigC is a configuration struct
type ConfigC struct {
	Field1 string
	Field2 int
}

// DefaultC returns the default configuration
func DefaultC() ConfigC {
	return ConfigC{
		Field1: "default",
		Field2: 42,
	}
}
`
	assert.NoError(t, os.WriteFile(packageCFile, []byte(packageCContent), 0644))
}