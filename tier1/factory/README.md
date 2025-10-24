# Factory Pattern

## Overview

The Factory Pattern is a creational design pattern that provides an interface for creating objects without specifying their exact classes. It encapsulates object creation logic, making code more flexible and maintainable.

## Problem

When object creation involves complex logic or when the exact type of object needed depends on runtime conditions, direct instantiation becomes problematic:

- **Tight coupling**: Clients know about concrete types and their dependencies
- **Complex initialization**: Creation logic scattered throughout the codebase
- **Difficult to extend**: Adding new types requires modifying client code
- **Centralization**: No single place to manage object creation logic
- **Dependency management**: Clients responsible for wiring up dependencies

### Real-World Context

Consider a payment processing system supporting multiple payment methods (credit card, PayPal, Bitcoin). Clients shouldn't need to know about each payment processor's specific implementation details. A factory encapsulates the logic of choosing and creating the right payment processor based on the payment method.

## Why Use This Pattern?

- **Decoupling**: Clients depend on abstractions, not concrete implementations
- **Centralized logic**: Single place to manage object creation
- **Easy extension**: Add new types without modifying client code
- **Dependency injection**: Factories can manage dependencies elegantly
- **Go conventions**: Constructor functions are the Go idiomatic factory pattern

## When to Use

- Creating objects with complex initialization
- Object type depends on runtime conditions
- Multiple related types that share an interface
- Encapsulating creation logic that might change
- Dependency management and configuration
- Plugin systems or extensible architectures

## When NOT to Use

- Simple object creation (direct constructor is clearer)
- Single concrete type needed
- No variation in object creation needed
- Performance-critical code where factory overhead matters
- Trivial objects with obvious initialization

## Implementation Guidelines

1. **Common interface**: All created objects implement a shared interface
2. **Concrete types**: Multiple implementations of the interface
3. **Factory function**: Takes parameters and returns interface type
4. **Creation logic**: Handles dependency injection and initialization
5. **Error handling**: Return errors for invalid parameters or failed creation

## Go Idioms

In Go, the factory pattern is implemented through **constructor functions**:

- Functions starting with `New` or `NewType` are factories
- Return concrete type or interface, not `*SomeType`
- Can take configuration parameters
- Often combined with the Options pattern for flexible configuration

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│       Factory Pattern Architecture           │
└──────────────────────────────────────────────┘

Client Code Request
        │
        ▼
┌─────────────────────┐
│  Factory Function   │
│  (Create based on   │
│   runtime params)   │
└────────────────────┬┘
        │
        ├─ Check Type/Config
        │
    ┌───┴─────────────────────┐
    │                         │
    ▼                         ▼
┌──────────────┐      ┌──────────────┐
│ ConcreteA    │      │ ConcreteB    │
│ implements   │      │ implements   │
│ Interface    │      │ Interface    │
└──────────────┘      └──────────────┘
    │                         │
    └────────┬────────────────┘
             │
             ▼
       (Interface Type)
             │
             ▼
        Client Uses

Creation Decision Flow:

params = GetConfig()
        │
        ├─ Type == "A" ──> Create ConcreteA ──┐
        │                                     │
        ├─ Type == "B" ──> Create ConcreteB ──┼──> Return to Client
        │                                     │
        └─ Type == "C" ──> Create ConcreteC ──┘
```

## Real-World Examples

### 1. Payment Processor Factory

```go
processor := NewPaymentProcessor(paymentMethod)
// Returns PaymentProcessor interface
// Handles credit card, PayPal, Bitcoin, etc.
```

### 2. Database Connection Factory

```go
db := NewDatabaseConnection(config)
// Returns Database interface
// Handles PostgreSQL, MySQL, SQLite, MongoDB, etc.
```

### 3. Logger Factory

```go
logger := NewLogger(logLevel, format)
// Returns Logger interface
// Handles file, syslog, JSON, console formats
```

## Key Advantages

- **Decoupling**: Classes depend on abstractions, not implementations
- **Centralized creation**: All object creation logic in one place
- **Easy to extend**: Add new types without modifying existing code
- **Configuration**: Separate creation configuration from usage
- **Testability**: Easy to create mock implementations for testing
- **Maintainability**: Changes to creation logic affect only one place

## Key Gotchas

- **Overcomplexity**: Don't use factories for simple objects
- **Hidden dependencies**: Factory may hide complex initialization chains
- **Type switching**: Clients may need to know about types despite factory
- **Documentation**: Factory must clearly document what types it creates
- **Error handling**: Must properly handle invalid parameters
- **Performance**: Function call overhead (usually negligible)
