# README Fixer

## Overview

The README Fixer is a tool that automatically updates README.md files in component directories to follow the template structure defined in COMPONENT_README_TEMPLATE.md. It preserves existing content while adding missing sections and fixing incorrect example paths.

## Features

- **Automatic Section Addition**: Adds missing sections required by the template
- **Content Preservation**: Preserves existing content in sections that are already present
- **Example Path Correction**: Fixes paths to example files (changes ../examples/... to ../EXAMPLES/...)
- **Component Name Detection**: Automatically detects the component name from the directory
- **Selective Processing**: Only processes README.md files in component directories

## Usage

### Command Line

```bash
go run tools/readme_fixer/main.go
```

## How It Works

1. The fixer finds all README.md files in component directories (directories containing .go files).
2. For each README.md file, it:
   - Reads the existing content
   - Fixes any incorrect example paths (../examples/... to ../EXAMPLES/...)
   - Identifies existing sections and their content
   - Creates a new file with all required sections, using existing content where available
   - For missing sections, it uses templates with the component name substituted
   - Writes the updated content back to the file

## Template Sections

The fixer ensures that each README.md file contains the following sections:

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

## Verification

After running the fixer, you can verify that all README.md files follow the template structure by running:

```bash
make validate-readme
```

## Integration with CI

The README validator is integrated with the CI process through the Makefile. The `ci` target includes the `validate-readme` target, which ensures that all README.md files are validated as part of the continuous integration process.

```bash
make ci
```

This ensures that all README.md files in the repository follow the standardized structure.