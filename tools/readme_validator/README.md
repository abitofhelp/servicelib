# README Validator

## Overview

The README Validator is a tool that validates README.md files in component directories against the template structure defined in COMPONENT_README_TEMPLATE.md. It ensures that all component README.md files follow a consistent structure and use the correct paths for example files.

## Features

- **Template Structure Validation**: Checks if README.md files contain all required sections from the template
- **Example Path Validation**: Ensures that example paths use the correct format (../EXAMPLES/... instead of ../examples/...)
- **Component Directory Detection**: Automatically identifies component directories (directories containing .go files)
- **Selective Validation**: Only validates README.md files in component directories, skipping special directories like EXAMPLES, vendor, DOCS, and tools

## Usage

### Command Line

```bash
go run tools/readme_validator/main.go
```

### Makefile

```bash
make validate-readme
```

## Integration with CI

The README validator is integrated with the CI process through the Makefile. The `ci` target includes the `validate-readme` target, which ensures that all README.md files are validated as part of the continuous integration process.

```bash
make ci
```

## How It Works

1. The validator first checks if the template file (COMPONENT_README_TEMPLATE.md) itself uses the correct example paths.
2. It then finds all README.md files in component directories (directories containing .go files).
3. For each README.md file, it checks if it contains all required sections from the template.
4. It also checks if the README.md file uses the correct example paths (../EXAMPLES/... instead of ../examples/...).
5. If any issues are found, it reports them and exits with a non-zero status code.

## Required Sections

The following sections are required in all component README.md files:

- Title (# Component Name)
- Overview (## Overview)
- Features (## Features)
- Installation (## Installation)
- Quick Start (## Quick Start)
- Configuration (## Configuration)
- API Documentation (## API Documentation)
- Core Types (### Core Types)
- Key Methods (### Key Methods)
- Examples (## Examples)
- Best Practices (## Best Practices)
- Troubleshooting (## Troubleshooting)
- Related Components (## Related Components)
- Contributing (## Contributing)
- License (## License)

## Example Path Format

All example paths in README.md files should use the format `../EXAMPLES/...` (uppercase) instead of `../examples/...` (lowercase). This ensures consistency with the actual directory structure of the project.