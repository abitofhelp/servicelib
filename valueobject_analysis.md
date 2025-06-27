# Value Objects vs. Primitive Types: Analysis and Best Practices

## Introduction

This document analyzes when to use value objects versus primitive types in Go applications, with specific reference to the patterns implemented in the `servicelib` codebase. It provides guidance on decision-making criteria and illustrates these with concrete examples from the codebase.

## What are Value Objects?

Value objects are immutable objects that represent concepts in the domain model. They are defined by their attributes rather than their identity, meaning two value objects with the same attributes are considered equal regardless of their memory location. In Domain-Driven Design (DDD), value objects are a fundamental building block for creating expressive, type-safe domain models.

## Key Characteristics of Value Objects in Our Codebase

1. **Immutability**: Value objects cannot be changed after creation
2. **Equality by Value**: Two value objects with the same attributes are considered equal
3. **Self-Validation**: Value objects validate their invariants upon creation
4. **Domain-Specific Behavior**: Value objects encapsulate domain logic related to their concept
5. **Type Safety**: Value objects provide compile-time type checking for domain concepts

## Decision Criteria: When to Use Value Objects vs. Primitives

### Use Value Objects When:

1. **The concept has validation rules**
   - Example: `Email` validates format, `Money` validates currency code
   - Benefit: Validation happens at creation time, preventing invalid states

2. **The concept has domain-specific behavior**
   - Example: `Money` provides arithmetic operations that respect currency
   - Benefit: Domain logic is encapsulated with the data it operates on

3. **The concept has multiple attributes that belong together**
   - Example: `Money` combines amount and currency
   - Benefit: Related attributes stay together, preventing inconsistent states

4. **The concept appears in multiple places in the domain**
   - Example: `Email` used in user profiles, notifications, etc.
   - Benefit: Consistent implementation across the codebase

5. **The concept needs special equality semantics**
   - Example: `Name` with case-insensitive equality
   - Benefit: Domain-appropriate equality checks

6. **The concept requires type safety beyond what primitives offer**
   - Example: Preventing confusion between different ID types
   - Benefit: Compile-time type checking prevents logical errors

### Use Primitives When:

1. **The concept is truly primitive with no validation or behavior**
   - Example: Simple counters, flags
   - Benefit: Simplicity and performance

2. **The concept is used only locally and doesn't cross boundaries**
   - Example: Loop counters, temporary variables
   - Benefit: Reduced overhead for implementation

3. **Performance is critical and the overhead of value objects is measurable**
   - Example: High-frequency operations where allocation matters
   - Benefit: Reduced memory allocation and CPU cycles

## Examples from the Codebase

### Example 1: `Name` Value Object

```
// Example code (not actual implementation)
// Name represents a person's name value object
type Name string

// NewName creates a new Name with validation
func NewName(name string) (Name, error) {
    // Trim whitespace
    trimmedName := strings.TrimSpace(name)

    // Validate name is not empty
    if trimmedName == "" {
        return "", errors.New("name cannot be empty")
    }

    // Validate name length
    if len(trimmedName) > 100 {
        return "", errors.New("name is too long")
    }

    return Name(trimmedName), nil
}
```

**Analysis**: The `Name` value object adds validation (non-empty, maximum length) and behavior (case-insensitive equality) to what would otherwise be a simple string. This ensures that names throughout the application are consistently validated and compared.

### Example 2: `IPAddress` Value Object

```
// Example code (not actual implementation)
// IPAddress represents an IP address value object that supports both IPv4 and IPv6 formats
type IPAddress string

// NewIPAddress creates a new IPAddress with validation for both IPv4 and IPv6 formats
func NewIPAddress(ip string) (IPAddress, error) {
    // Trim whitespace
    trimmedIP := strings.TrimSpace(ip)

    // Empty IP is allowed (optional field)
    if trimmedIP == "" {
        return "", nil
    }

    // Validate IP format (both IPv4 and IPv6)
    parsedIP := net.ParseIP(trimmedIP)
    if parsedIP == nil {
        return "", errors.New("invalid IP address format")
    }

    // Use the normalized string format from net.ParseIP
    return IPAddress(parsedIP.String()), nil
}
```

**Analysis**: The `IPAddress` value object adds validation and normalization to IP addresses. It also provides domain-specific methods like `IsIPv4()`, `IsIPv6()`, `IsLoopback()`, and `IsPrivate()` that encapsulate knowledge about IP addresses. This is significantly more valuable than using a string primitive.

### Example 3: `Money` Value Object

```
// Example code (not actual implementation)
// Money represents a monetary value object with amount and currency
type Money struct {
    base.BaseStructValueObject

    amount decimal.Decimal
    currency string
}

// Add adds another Money value and returns a new Money
// Both Money values must have the same currency
func (v Money) Add(other Money) (Money, error) {
    if v.currency != other.currency {
        return Money{}, errors.New("cannot add money with different currencies")
    }

    return Money{
        amount:   v.amount.Add(other.amount),
        currency: v.currency,
    }, nil
}
```

**Analysis**: The `Money` value object is a complex example that combines multiple attributes (amount and currency) and provides rich domain-specific behavior (arithmetic operations that respect currency, comparison operations, formatting). It ensures that money calculations follow business rules (e.g., can't add different currencies) and prevents common errors.

## Recommendations for Future Development

1. **Identify Domain Concepts**: Analyze your domain model to identify concepts that would benefit from value object implementation

2. **Start with High-Value Targets**: Begin with concepts that have validation rules or domain-specific behavior

3. **Consider API Boundaries**: Value objects are particularly valuable at API boundaries where validation is essential

4. **Balance Pragmatism**: Not everything needs to be a value object; use judgment based on the criteria above

5. **Consistency**: Once a concept is implemented as a value object, use it consistently throughout the codebase

6. **Testing**: Value objects are typically easy to test due to their immutability and self-contained behavior

## Conclusion

Value objects provide significant benefits in terms of code quality, domain expressiveness, and error prevention. By encapsulating validation and behavior with the data they represent, they create more robust and maintainable code. The decision to use value objects should be guided by the nature of the concept being modeled and the specific requirements of your application.

The `servicelib` codebase demonstrates effective use of value objects for concepts ranging from simple identifiers to complex monetary values. This pattern should be continued and expanded where appropriate to maintain the quality and expressiveness of the codebase.
