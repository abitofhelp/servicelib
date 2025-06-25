# Date Package Examples

This directory contains examples demonstrating how to use the `date` package, which provides utilities for working with dates and times in Go applications. The package offers a consistent way to parse, format, and manipulate dates using a standardized format.

## Examples

### 1. Basic Usage Example

[basic_usage/main.go](basic_usage/main.go)

Demonstrates basic date parsing and formatting.

Key concepts:
- Parsing date strings into time.Time objects
- Extracting components from parsed dates (year, month, day, etc.)
- Formatting time.Time objects into standardized strings
- Parsing formatted dates back into time.Time objects

### 2. Optional Date Example

[optional_date/main.go](optional_date/main.go)

Shows how to work with optional (nullable) dates.

Key concepts:
- Parsing nil and non-nil date strings
- Handling nil date pointers
- Formatting nil and non-nil dates
- Using optional dates in structs

### 3. Error Handling Example

[error_handling/main.go](error_handling/main.go)

Demonstrates error handling for invalid dates.

Key concepts:
- Handling invalid date formats
- Detecting validation errors
- Graceful error handling in functions
- Processing multiple date parsing attempts

### 4. Time Zone Example

[time_zone/main.go](time_zone/main.go)

Shows how to work with dates in different time zones.

Key concepts:
- Parsing dates with different time zones
- Comparing dates from different time zones
- Converting between time zones
- Formatting dates with time zone information
- Working with local time zone

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run basic_usage/main.go
```

## Key Features of the Date Package

1. **Standardized Date Format**: Uses RFC3339 format consistently throughout the application.

2. **Date Parsing**: Parse date strings into time.Time objects with proper error handling.

3. **Date Formatting**: Format time.Time objects into standardized strings.

4. **Optional Date Handling**: Special handling for nil date pointers, allowing for nullable dates in your application.

5. **Error Handling**: Proper validation and error handling for invalid date formats.

## Best Practices

1. **Consistent Date Format**: Always use the package's standard date format for consistency across your application.

2. **Error Handling**: Always check for errors when parsing dates.

3. **Nil Handling**: Use the optional date functions when dealing with dates that might be nil.

4. **Time Zones**: Be aware of time zone implications when parsing and formatting dates.

5. **Date Comparisons**: When comparing dates, ensure they are in the same time zone to avoid unexpected results.

## Additional Resources

For more information about the date package, see the [date package documentation](../../date/README.md).
