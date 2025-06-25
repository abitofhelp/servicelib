// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

// Template content for missing sections
var sectionTemplates = map[string]string{
	"## Overview": `
Brief description of the component and its purpose in the ServiceLib library.
`,
	"## Features": `
- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3
`,
	"## Installation": `
` + "```bash" + `
go get github.com/abitofhelp/servicelib/component
` + "```" + `
`,
	"## Quick Start": `
See the [Quick Start example](../EXAMPLES/component/quickstart_example.go) for a complete, runnable example of how to use the component.
`,
	"## Configuration": `
See the [Configuration example](../EXAMPLES/component/configuration_example.go) for a complete, runnable example of how to configure the component.
`,
	"## API Documentation": `
`,
	"### Core Types": `
Description of the main types provided by the component.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/component/type1_example.go) for a complete, runnable example of how to use Type 1.
`,
	"### Key Methods": `
Description of the key methods provided by the component.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/component/method1_example.go) for a complete, runnable example of how to use Method 1.
`,
	"## Examples": `
For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/component/basic_usage_example.go) - Shows basic usage of the component
- [Advanced Configuration](../EXAMPLES/component/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/component/error_handling_example.go) - Shows how to handle errors
`,
	"## Best Practices": `
1. **Best Practice 1**: Description of best practice 1
2. **Best Practice 2**: Description of best practice 2
3. **Best Practice 3**: Description of best practice 3
`,
	"## Troubleshooting": `
### Common Issues

#### Issue 1

Description of issue 1 and how to resolve it.

#### Issue 2

Description of issue 2 and how to resolve it.
`,
	"## Related Components": `
- [Component 1](../component1/README.md) - Description of how this component relates to Component 1
- [Component 2](../component2/README.md) - Description of how this component relates to Component 2
`,
	"## Contributing": `
Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.
`,
	"## License": `
This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
`,
}

// fixReadme updates a README.md file to follow the template structure
func fixReadme(filePath string) error {
	// Read the existing content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	contentStr := string(content)

	// Fix example paths
	contentStr = strings.ReplaceAll(contentStr, "../examples/", "../EXAMPLES/")

	// Get the component name from the directory
	componentDir := filepath.Base(filepath.Dir(filePath))

	// Parse the existing content to identify existing sections
	existingSections := make(map[string]string)
	currentSection := ""
	var sectionContent strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(contentStr))
	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line starts a new section
		isNewSection := false
		for _, section := range requiredSections {
			if strings.HasPrefix(line, section) {
				isNewSection = true
				// Save the previous section content if any
				if currentSection != "" {
					existingSections[currentSection] = sectionContent.String()
					sectionContent.Reset()
				}
				currentSection = line
				break
			}
		}

		if !isNewSection && currentSection != "" {
			sectionContent.WriteString(line + "\n")
		}
	}

	// Save the last section
	if currentSection != "" {
		existingSections[currentSection] = sectionContent.String()
	}

	// Create the updated content with all required sections
	var updatedContent strings.Builder

	// Start with the title
	titleFound := false
	for key := range existingSections {
		if strings.HasPrefix(key, "# ") {
			updatedContent.WriteString(key + existingSections[key])
			titleFound = true
			break
		}
	}

	// If no title found, create one based on the component name
	if !titleFound {
		title := "# " + strings.Title(strings.ReplaceAll(componentDir, "_", " "))
		updatedContent.WriteString(title + "\n\n")
	}

	// Add all other required sections
	for _, section := range requiredSections[1:] { // Skip the title
		updatedContent.WriteString("\n" + section + "\n")

		// Use existing content if available, otherwise use template
		if content, ok := existingSections[section]; ok {
			updatedContent.WriteString(content)
		} else {
			// Replace "component" with the actual component name in templates
			template := sectionTemplates[section]
			template = strings.ReplaceAll(template, "component", componentDir)
			updatedContent.WriteString(template)
		}
	}

	// Write the updated content back to the file
	return os.WriteFile(filePath, []byte(updatedContent.String()), 0644)
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

func main() {
	// Get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Find all README.md files
	readmeFiles, err := findReadmeFiles(rootDir)
	if err != nil {
		fmt.Printf("Error finding README.md files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d README.md files to fix:\n", len(readmeFiles))

	// Fix each README.md file
	fixedCount := 0
	for _, filePath := range readmeFiles {
		fmt.Printf("Fixing %s...\n", filePath)
		err := fixReadme(filePath)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fixedCount++
		}
	}

	fmt.Printf("\nFixed %d out of %d README.md files.\n", fixedCount, len(readmeFiles))
	fmt.Println("Run 'make validate-readme' to verify the fixes.")
}
