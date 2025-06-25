// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Required sections in README.md files based on the template
var requiredSections = []string{
	"# ",                // Title
	"## Overview",
	"## Features",
	"## Installation",
	"## Quick Start",
	"## Configuration",
	"## API Documentation",
	"### Core Types",
	"### Key Methods",
	"## Examples",
	"## Best Practices",
	"## Troubleshooting",
	"## Related Components",
	"## Contributing",
	"## License",
}

// validateReadme checks if a README.md file follows the template structure
func validateReadme(filePath string) (bool, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, []string{fmt.Sprintf("Error opening file: %v", err)}
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, []string{fmt.Sprintf("Error reading file: %v", err)}
	}

	contentStr := string(content)
	issues := []string{}

	// Check for required sections
	for _, section := range requiredSections {
		if !strings.Contains(contentStr, section) {
			issues = append(issues, fmt.Sprintf("Missing section: %s", section))
		}
	}

	// Check for incorrect example paths
	// Look for "../examples/" (lowercase) instead of "../EXAMPLES/" (uppercase)
	examplesRegex := regexp.MustCompile(`\.\./examples/`)
	if examplesRegex.MatchString(contentStr) {
		issues = append(issues, "Incorrect example path format: using '../examples/' instead of '../EXAMPLES/'")

		// Find all instances of incorrect paths
		scanner := bufio.NewScanner(strings.NewReader(contentStr))
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if examplesRegex.MatchString(line) {
				issues = append(issues, fmt.Sprintf("  Line %d: %s", lineNum, strings.TrimSpace(line)))
			}
		}
	}

	return len(issues) == 0, issues
}

// findReadmeFiles finds all README.md files in component directories
func findReadmeFiles(rootDir string) ([]string, error) {
	var readmeFiles []string

	// Define directories to skip
	skipDirs := map[string]bool{
		"EXAMPLES": true,
		"vendor":   true,
		"DOCS":     true,
		"tools":    true,
	}

	// Walk through the directory tree
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Skip specified directories
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}

		// Check if this is a README.md file
		if !info.IsDir() && info.Name() == "README.md" {
			// Skip the root README.md
			if path == filepath.Join(rootDir, "README.md") {
				return nil
			}

			// Get the parent directory
			parentDir := filepath.Dir(path)

			// Check if the parent directory is a component directory (contains .go files)
			isComponent := false
			entries, err := os.ReadDir(parentDir)
			if err != nil {
				return nil
			}

			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
					isComponent = true
					break
				}
			}

			if isComponent {
				readmeFiles = append(readmeFiles, path)
			}
		}

		return nil
	})

	return readmeFiles, err
}

// validateTemplateFile checks if the template file itself has the correct example paths
func validateTemplateFile(templatePath string) (bool, []string) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return false, []string{fmt.Sprintf("Error reading template file: %v", err)}
	}

	contentStr := string(content)
	issues := []string{}

	// Check for incorrect example paths in the template
	examplesRegex := regexp.MustCompile(`\.\./examples/`)
	if examplesRegex.MatchString(contentStr) {
		issues = append(issues, "Template contains incorrect example path format: using '../examples/' instead of '../EXAMPLES/'")

		// Find all instances of incorrect paths
		scanner := bufio.NewScanner(strings.NewReader(contentStr))
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if examplesRegex.MatchString(line) {
				issues = append(issues, fmt.Sprintf("  Line %d: %s", lineNum, strings.TrimSpace(line)))
			}
		}
	}

	return len(issues) == 0, issues
}

func main() {
	// Get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// First, validate the template file
	templatePath := filepath.Join(rootDir, "COMPONENT_README_TEMPLATE.md")
	templateValid, templateIssues := validateTemplateFile(templatePath)
	if !templateValid {
		fmt.Println("Issues found in the template file:")
		for _, issue := range templateIssues {
			fmt.Printf("  %s\n", issue)
		}
		fmt.Println("Please fix the template file before validating README files.")
		os.Exit(1)
	}

	// Find all README.md files
	readmeFiles, err := findReadmeFiles(rootDir)
	if err != nil {
		fmt.Printf("Error finding README.md files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d README.md files to validate:\n", len(readmeFiles))
	for _, filePath := range readmeFiles {
		fmt.Printf("  %s\n", filePath)
	}
	fmt.Println()

	// Validate each README.md file
	allValid := true
	for _, filePath := range readmeFiles {
		valid, issues := validateReadme(filePath)
		if !valid {
			allValid = false
			fmt.Printf("Issues found in %s:\n", filePath)
			for _, issue := range issues {
				fmt.Printf("  %s\n", issue)
			}
			fmt.Println()
		}
	}

	if !allValid {
		fmt.Println("Validation failed. Please fix the issues above.")
		os.Exit(1)
	}

	fmt.Printf("All README.md files (%d) are valid and follow the template structure.\n", len(readmeFiles))
}
