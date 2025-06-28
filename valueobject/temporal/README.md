# Value Object Temporal

## Overview

The Value Object Temporal component provides immutable value objects for representing and manipulating temporal concepts such as dates, times, durations, and intervals in the ServiceLib library.

## Features

- **Immutable Values**: All temporal value objects are immutable
- **Type Safety**: Strong typing for temporal concepts
- **Validation**: Built-in validation for temporal values
- **Comparison Operations**: Methods for comparing temporal values
- **Formatting**: Methods for formatting temporal values

## Installation

```bash
go get github.com/abitofhelp/servicelib/valueobject/temporal
```

## Quick Start

See the [Quick Start example](../../EXAMPLES/valueobject/temporal/README.md) for a complete, runnable example of how to use the temporal value objects.

## API Documentation

### Core Types

Description of the main types provided by the component.

#### DateVO

Represents a date value object.

```
type DateVO struct {
    // Fields
}
```

#### TimeVO

Represents a time value object.

```
type TimeVO struct {
    // Fields
}
```

#### DateTimeVO

Represents a date-time value object.

```
type DateTimeVO struct {
    // Fields
}
```

#### DurationVO

Represents a duration value object.

```
type DurationVO struct {
    // Fields
}
```

#### IntervalVO

Represents a time interval value object.

```
type IntervalVO struct {
    // Fields
}
```

### Key Methods

Description of the key methods provided by the component.

#### NewDateVO

Creates a new date value object.

```
func NewDateVO(year int, month time.Month, day int) (*DateVO, error)
```

#### NewTimeVO

Creates a new time value object.

```
func NewTimeVO(hour, minute, second, nanosecond int) (*TimeVO, error)
```

#### NewDateTimeVO

Creates a new date-time value object.

```
func NewDateTimeVO(date *DateVO, time *TimeVO) (*DateTimeVO, error)
```

#### NewDurationVO

Creates a new duration value object.

```
func NewDurationVO(duration time.Duration) (*DurationVO, error)
```

#### NewIntervalVO

Creates a new interval value object.

```
func NewIntervalVO(start, end *DateTimeVO) (*IntervalVO, error)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../../EXAMPLES/valueobject/temporal/basic_usage/README.md) - Shows basic usage of temporal value objects
- [Comparison](../../EXAMPLES/valueobject/temporal/comparison/README.md) - Shows how to compare temporal value objects
- [Formatting](../../EXAMPLES/valueobject/temporal/formatting/README.md) - Shows how to format temporal value objects

## Best Practices

1. **Immutability**: Never modify temporal value objects, always create new ones
2. **Validation**: Always validate input when creating temporal value objects
3. **Comparison**: Use the provided comparison methods instead of comparing fields directly
4. **Time Zones**: Be aware of time zone implications when working with date-time value objects
5. **Formatting**: Use the provided formatting methods for consistent output

## Troubleshooting

### Common Issues

#### Invalid Date/Time Values

If you encounter validation errors when creating temporal value objects, ensure that the input values are valid (e.g., month between 1-12, day appropriate for the month, etc.).

#### Time Zone Issues

If you're experiencing unexpected behavior with date-time value objects, check the time zone settings. Consider using UTC for consistency.

## Related Components

- [Value Object](../README.md) - The main value object component
- [Value Object Base](../base/README.md) - Base functionality for value objects

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.