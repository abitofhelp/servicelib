# Stringutil Module

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

See the [Case-Insensitive Operations example](../examples/stringutil/case_insensitive_operations_example.go) for a complete, runnable example of how to use the Stringutil module.

## API Documentation

### Case-Insensitive Operations

The module provides functions for case-insensitive string operations like prefix checking and substring searching.

#### Case-Insensitive String Operations

See the [Case-Insensitive Operations example](../examples/stringutil/case_insensitive_operations_example.go) for a complete, runnable example of how to use case-insensitive string operations.

### Prefix Checking

The `HasAnyPrefix` function checks if a string starts with any of the specified prefixes.

#### Checking String Prefixes

See the [Prefix Checking example](../examples/stringutil/prefix_checking_example.go) for a complete, runnable example of how to check string prefixes.

### String Joining with Natural Language Formatting

The `JoinWithAnd` function joins a slice of strings with commas and "and", with optional Oxford comma support.

#### Joining Strings with Natural Language Formatting

See the [String Joining example](../examples/stringutil/string_joining_example.go) for a complete, runnable example of how to join strings with natural language formatting.

### Whitespace Handling

The module provides functions for checking if a string is empty or contains only whitespace, and for removing whitespace.

#### Handling Whitespace in Strings

See the [Whitespace Handling example](../examples/stringutil/whitespace_handling_example.go) for a complete, runnable example of how to handle whitespace in strings.

### String Truncation

The `Truncate` function truncates a string to a specified maximum length, adding an ellipsis if needed.

#### Truncating Strings

See the [String Truncation example](../examples/stringutil/string_truncation_example.go) for a complete, runnable example of how to truncate strings.

### Path Normalization

The `ForwardSlashPath` function converts backslashes in a path to forward slashes for consistent path handling across platforms.

#### Normalizing File Paths

See the [Path Normalization example](../examples/stringutil/path_normalization_example.go) for a complete, runnable example of how to normalize file paths.

## Best Practices

1. **Performance Considerations**: For performance-critical code, be aware that some functions like `ContainsIgnoreCase` create temporary strings.

2. **String Immutability**: Remember that strings in Go are immutable, so functions like `Truncate` and `RemoveWhitespace` return new strings.

3. **Unicode Support**: The functions in this package work with UTF-8 encoded strings, which is Go's default string encoding.

4. **Error Handling**: These utility functions don't return errors, so validate input when necessary before calling them.

5. **Path Handling**: For more comprehensive path manipulation, consider using the standard library's `path/filepath` package alongside `ForwardSlashPath`.

6. **Regular Expressions**: `RemoveWhitespace` uses regular expressions, which can be expensive for very large strings or frequent calls.

7. **String Builders**: For complex string manipulations, consider using `strings.Builder` for better performance.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
