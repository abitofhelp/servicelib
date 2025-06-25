# Stringutil Module

## Overview

The Stringutil Module provides additional string manipulation utilities beyond what's available in the standard library. It offers a collection of helper functions for common string operations that are frequently needed in Go applications.


## Features

- **Case-Insensitive Operations**: Functions for case-insensitive string comparisons
- **Prefix Checking**: Utilities for checking string prefixes with various options
- **String Joining**: Advanced string joining with natural language formatting
- **Whitespace Handling**: Functions for checking and manipulating whitespace
- **String Truncation**: Safely truncate strings with ellipsis
- **Path Normalization**: Convert file paths to use consistent separators
- **Regular Expression Utilities**: String manipulation using regular expressions


## Installation

```bash
go get github.com/abitofhelp/servicelib/stringutil
```


## Quick Start

See the [Case-Insensitive Operations example](../EXAMPLES/stringutil/case_insensitive_operations_example/main.go) for a complete, runnable example of how to use the Stringutil module.


## Configuration

The stringutil package does not require any configuration.


## API Documentation


### Core Types

The stringutil package primarily provides utility functions rather than types.


### Key Methods

#### Case-Insensitive Operations

The module provides functions for case-insensitive string operations like prefix checking and substring searching.

```
func ContainsIgnoreCase(s, substr string) bool
func EqualIgnoreCase(s1, s2 string) bool
func HasPrefixIgnoreCase(s, prefix string) bool
func HasSuffixIgnoreCase(s, suffix string) bool
```

See the [Case-Insensitive Operations example](../EXAMPLES/stringutil/case_insensitive_operations_example/main.go) for a complete, runnable example of how to use case-insensitive string operations.

#### Prefix Checking

The `HasAnyPrefix` function checks if a string starts with any of the specified prefixes.

```
func HasAnyPrefix(s string, prefixes ...string) bool
func HasAnyPrefixIgnoreCase(s string, prefixes ...string) bool
```

See the [Prefix Checking example](../EXAMPLES/stringutil/prefix_checking_example/main.go) for a complete, runnable example of how to check string prefixes.

#### String Joining with Natural Language Formatting

The `JoinWithAnd` function joins a slice of strings with commas and "and", with optional Oxford comma support.

```
func JoinWithAnd(items []string, useOxfordComma bool) string
```

See the [String Joining example](../EXAMPLES/stringutil/string_joining_example/main.go) for a complete, runnable example of how to join strings with natural language formatting.

#### Whitespace Handling

The module provides functions for checking if a string is empty or contains only whitespace, and for removing whitespace.

```
func IsBlank(s string) bool
func RemoveWhitespace(s string) string
```

See the [Whitespace Handling example](../EXAMPLES/stringutil/whitespace_handling_example/main.go) for a complete, runnable example of how to handle whitespace in strings.

#### String Truncation

The `Truncate` function truncates a string to a specified maximum length, adding an ellipsis if needed.

```
func Truncate(s string, maxLength int) string
```

See the [String Truncation example](../EXAMPLES/stringutil/string_truncation_example/main.go) for a complete, runnable example of how to truncate strings.

#### Path Normalization

The `ForwardSlashPath` function converts backslashes in a path to forward slashes for consistent path handling across platforms.

```
func ForwardSlashPath(path string) string
```

See the [Path Normalization example](../EXAMPLES/stringutil/path_normalization_example/main.go) for a complete, runnable example of how to normalize file paths.


## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Case-Insensitive Operations](../EXAMPLES/stringutil/case_insensitive_operations_example/main.go) - Shows how to use case-insensitive string operations
- [Prefix Checking](../EXAMPLES/stringutil/prefix_checking_example/main.go) - Shows how to check string prefixes
- [String Joining](../EXAMPLES/stringutil/string_joining_example/main.go) - Shows how to join strings with natural language formatting
- [Whitespace Handling](../EXAMPLES/stringutil/whitespace_handling_example/main.go) - Shows how to handle whitespace in strings
- [String Truncation](../EXAMPLES/stringutil/string_truncation_example/main.go) - Shows how to truncate strings
- [Path Normalization](../EXAMPLES/stringutil/path_normalization_example/main.go) - Shows how to normalize file paths


## Best Practices

1. **Performance Considerations**: For performance-critical code, be aware that some functions like `ContainsIgnoreCase` create temporary strings.

2. **String Immutability**: Remember that strings in Go are immutable, so functions like `Truncate` and `RemoveWhitespace` return new strings.

3. **Unicode Support**: The functions in this package work with UTF-8 encoded strings, which is Go's default string encoding.

4. **Error Handling**: These utility functions don't return errors, so validate input when necessary before calling them.

5. **Path Handling**: For more comprehensive path manipulation, consider using the standard library's `path/filepath` package alongside `ForwardSlashPath`.

6. **Regular Expressions**: `RemoveWhitespace` uses regular expressions, which can be expensive for very large strings or frequent calls.

7. **String Builders**: For complex string manipulations, consider using `strings.Builder` for better performance.


## Troubleshooting

### Common Issues

#### Performance Issues

**Issue**: String operations are slow for large strings or frequent calls.

**Solution**: For performance-critical code, consider using `strings.Builder` for complex string manipulations, and be aware that some functions like `ContainsIgnoreCase` and `RemoveWhitespace` may create temporary strings or use regular expressions.

#### Unicode Handling

**Issue**: Unexpected behavior with non-ASCII characters.

**Solution**: The functions in this package work with UTF-8 encoded strings, which is Go's default string encoding. For more complex Unicode operations, consider using the `unicode` package from the standard library.


## Related Components

- [Validation](../validation/README.md) - The validation component can use string utilities for input validation.
- [Logging](../logging/README.md) - The logging component can use string utilities for formatting log messages.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
