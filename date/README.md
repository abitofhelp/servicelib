# Date

## Overview

The Date component provides utilities for working with dates and times in a standardized format. It simplifies common date operations like parsing, formatting, and handling optional (nullable) dates.

## Features

- **Standardized Date Format**: Uses RFC3339 as the standard date format throughout an application
- **Date Parsing**: Parse date strings into time.Time objects with proper error handling
- **Date Formatting**: Format time.Time objects as strings in the standard format
- **Optional Date Support**: Handle nullable dates with pointer-based functions
- **Validation Error Integration**: Return structured validation errors for invalid date formats

## Installation

```bash
go get github.com/abitofhelp/servicelib/date
```

## Quick Start

See the [Basic Usage example](../EXAMPLES/date/basic_usage/README.md) for a complete, runnable example of how to use the date component.

## API Documentation

### Core Constants

#### StandardDateFormat

The standard date format used throughout the application.

```go
const StandardDateFormat = time.RFC3339
```

### Key Methods

#### ParseDate

Parses a date string in the standard format.

```go
func ParseDate(dateStr string) (time.Time, error)
```

#### ParseOptionalDate

Parses an optional date string in the standard format.

```go
func ParseOptionalDate(dateStr *string) (*time.Time, error)
```

#### FormatDate

Formats a time.Time as a string in the standard format.

```go
func FormatDate(date time.Time) string
```

#### FormatOptionalDate

Formats an optional time.Time as a string in the standard format.

```go
func FormatOptionalDate(date *time.Time) *string
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/date/basic_usage/README.md) - Shows basic usage of the date component
- [Error Handling](../EXAMPLES/date/error_handling/README.md) - Shows how to handle date parsing errors
- [Time Zones](../EXAMPLES/date/time_zone/README.md) - Shows how to work with time zones
- [Optional Dates](../EXAMPLES/date/optional_date/README.md) - Shows how to work with optional (nullable) dates

## Best Practices

1. **Use Standard Format**: Always use the StandardDateFormat constant for consistency
2. **Handle Parsing Errors**: Always check for errors when parsing date strings
3. **Use Optional Functions**: Use ParseOptionalDate and FormatOptionalDate for nullable dates
4. **Time Zone Awareness**: Be aware of time zones when parsing and formatting dates
5. **Validate User Input**: Always validate user-provided date strings before parsing

## Troubleshooting

### Common Issues

#### Invalid Date Format

If you're seeing "invalid date format" errors, ensure the date string follows the RFC3339 format (e.g., "2023-01-02T15:04:05Z"). Common issues include:
- Missing time component
- Missing timezone indicator
- Incorrect separators between date parts

#### Time Zone Confusion

If dates are being parsed correctly but show unexpected times, check the time zone:
- RFC3339 dates with "Z" suffix are in UTC
- Dates with explicit offsets (e.g., "+01:00") are in that specific time zone
- Consider converting to a specific time zone if needed for display or comparison

## Related Components

- [Validation](../validation/README.md) - Validation utilities that work with date validation
- [Errors](../errors/README.md) - Error handling for date parsing errors

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
