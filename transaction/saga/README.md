# Transaction Saga

## Overview

The Transaction Saga component provides functionality for implementing the Saga pattern for distributed transactions in the ServiceLib library.

## Features

- **Saga Pattern**: Implements the Saga pattern for managing distributed transactions
- **Compensation Actions**: Supports defining compensation actions for each step in a saga
- **Transaction Coordination**: Coordinates multiple steps in a distributed transaction
- **Error Handling**: Provides robust error handling for transaction failures
- **Rollback Support**: Automatically rolls back completed steps when a step fails

## Installation

```bash
go get github.com/abitofhelp/servicelib/transaction/saga
```

## Quick Start

See the [Quick Start example](../../EXAMPLES/transaction/saga/README.md) for a complete, runnable example of how to use the transaction saga component.

## API Documentation

### Core Types

Description of the main types provided by the component.

#### Saga

The main type that represents a saga transaction.

```
type Saga struct {
    // Fields
}
```

#### Step

Represents a step in a saga transaction.

```
type Step struct {
    // Fields
}
```

### Key Methods

Description of the key methods provided by the component.

#### NewSaga

Creates a new saga transaction.

```
func NewSaga(ctx context.Context, options ...Option) (*Saga, error)
```

#### AddStep

Adds a step to a saga transaction.

```
func (s *Saga) AddStep(action Action, compensation Compensation) error
```

#### Execute

Executes a saga transaction.

```
func (s *Saga) Execute(ctx context.Context) error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Usage](../../EXAMPLES/transaction/saga/basic_usage/README.md) - Shows basic usage of the saga component
- [Compensation](../../EXAMPLES/transaction/saga/compensation/README.md) - Shows how to define compensation actions

## Best Practices

1. **Define Compensation Actions**: Always define compensation actions for each step in a saga
2. **Use Context**: Pass context through the saga to enable proper cancellation and timeouts
3. **Handle Errors**: Properly handle errors from saga execution
4. **Idempotent Operations**: Ensure that saga steps are idempotent
5. **Transaction Boundaries**: Clearly define transaction boundaries

## Troubleshooting

### Common Issues

#### Compensation Failures

If compensation actions fail, the saga may be left in an inconsistent state. Ensure that compensation actions are robust and handle errors properly.

#### Timeout Issues

If saga execution times out, ensure that the context timeout is set appropriately and that individual steps have appropriate timeouts.

## Related Components

- [Transaction](../README.md) - The main transaction component
- [Errors](../../errors/README.md) - Error handling for transactions

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.