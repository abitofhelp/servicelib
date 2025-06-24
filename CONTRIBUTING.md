# Contributing to ServiceLib

Thank you for your interest in contributing to ServiceLib! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and considerate of others.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue on GitHub with the following information:

- A clear, descriptive title
- A detailed description of the issue
- Steps to reproduce the bug
- Expected behavior
- Actual behavior
- Any relevant logs or screenshots
- Your environment (Go version, OS, etc.)

### Suggesting Enhancements

If you have an idea for an enhancement, please create an issue on GitHub with the following information:

- A clear, descriptive title
- A detailed description of the enhancement
- Any relevant examples or use cases
- Any potential implementation details

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature-name`)
3. Make your changes
4. Run tests and linting
5. Commit your changes with a descriptive commit message
6. Push to your branch (`git push origin feature/your-feature-name`)
7. Create a Pull Request

## Development Workflow

### Prerequisites

- Go 1.24 or higher
- Make (optional, for using the Makefile)

### Setting Up the Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/abitofhelp/servicelib.git
   cd servicelib
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific tests
go test ./path/to/package -run TestName
```

For information on running integration tests, see the [Integration Tests Guide](DOCS/Integration_Tests.md).

### Linting

```bash
# Run linter
make lint
```

### Building

```bash
# Build the library
make build
```

## Coding Standards

- Follow Go best practices and style guidelines
- Write tests for new functionality
- Document public APIs
- Keep backward compatibility in mind
- Use meaningful variable and function names
- Keep functions small and focused on a single responsibility
- Use comments to explain complex logic

## Documentation

- Update documentation for any changes to the API
- Add examples for new features
- Keep README.md files up to date
- Document any breaking changes

## Versioning

ServiceLib follows semantic versioning (SemVer):

- Major version (X.y.z): Incompatible API changes
- Minor version (x.Y.z): Backwards-compatible functionality
- Patch version (x.y.Z): Backwards-compatible bug fixes

## License

By contributing to ServiceLib, you agree that your contributions will be licensed under the project's MIT License.
