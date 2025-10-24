# Adapter Pattern

## Overview

The Adapter Pattern is a structural design pattern that converts the interface of a class into another interface that clients expect. It lets classes with incompatible interfaces work together by providing a wrapper that translates between them.

## Problem

When integrating third-party libraries or legacy systems, you face interface mismatches:

- **Incompatible interfaces**: Existing code expects a different interface than what the library provides
- **Legacy system integration**: Old systems have outdated interfaces
- **Library changes**: Library updates change the interface, breaking dependent code
- **Multiple implementations**: Need to support multiple libraries with different interfaces
- **Unstable external APIs**: Want to insulate code from external API changes

### Real-World Context

Imagine you have a codebase that expects a `PaymentProcessor` interface with a specific signature, but you need to integrate with a third-party payment library that has a completely different interface. Instead of rewriting your code to match the library's interface, you create an adapter that translates between them.

## Why Use This Pattern?

- **Integration**: Makes incompatible interfaces work together
- **Decoupling**: Isolates your code from external library interfaces
- **Flexibility**: Support multiple implementations with different interfaces
- **Stability**: Changes to external libraries affect only the adapter
- **Reusability**: Existing code can use new libraries without modification

## When to Use

- Integrating third-party libraries with incompatible interfaces
- Supporting multiple implementations with different interfaces
- Isolating code from unstable external APIs
- Gradually migrating from legacy systems
- Creating a stable internal API that hides external complexity
- Supporting multiple database drivers or storage backends

## When NOT to Use

- Interfaces are already compatible (unnecessary indirection)
- Adapter becomes more complex than direct integration
- Only one implementation needed and it's unlikely to change
- Performance is critical and adapter overhead matters
- Adapter obscures rather than clarifies intent

## Implementation Guidelines

1. **Target interface**: Define what interface clients expect
2. **Adaptee**: The existing incompatible class or interface
3. **Adapter class**: Implements target interface and wraps adaptee
4. **Delegation**: Adapter translates calls to adaptee methods
5. **Data transformation**: Handle any necessary data format conversion

## Go Idioms

Go's interface system makes adapters elegant:

- Small, focused interfaces are easier to adapt to
- Implicit interface implementation simplifies adapter creation
- Function types can be adapted like structs
- Composition over inheritance makes adapters lightweight

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│        Adapter Pattern Architecture          │
└──────────────────────────────────────────────┘

Problem: Incompatible Interfaces

Client expects:        Library provides:
┌──────────────┐      ┌──────────────┐
│ PaymentProc  │      │ PaymentGW    │
│ Process(...) │      │ DoPayment()  │
└──────────────┘      │ GetStatus()  │
       ▲              └──────────────┘
       │                     │
       └─────── X ───────────┘
          (Incompatible!)


Solution: Adapter Bridge

Client Code
    │
    ▼
┌─────────────────────┐
│  Target Interface   │
│  (PaymentProcessor) │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────────────────┐
│  Adapter                        │
│  ┌─────────────────────────────┐│
│  │ func (a *Adapter)           ││
│  │ Process(...) error          ││
│  │  return a.gateway.DoPayment ││
│  └─────────────────────────────┘│
└──────────────┬──────────────────┘
               │
               ▼
        ┌──────────────┐
        │ External Lib │
        │ (Adaptee)    │
        └──────────────┘

Data Transformation:
┌──────────────────────────────────────┐
│  Client Format:                      │
│  Process(customerID, amount)         │
│                                      │
│  Adapter translates to:              │
│                                      │
│  Library Format:                     │
│  DoPayment(customer, value, fee)     │
└──────────────────────────────────────┘
```

## Real-World Examples

### 1. Database Driver Adapter

```go
// Client expects this interface
type Database interface {
    Query(sql string, args ...interface{}) (Result, error)
}

// Library provides different interface
type SQLiteDB struct { /*...*/ }
func (db *SQLiteDB) Execute(query string) (*Rows, error) { /*...*/ }

// Adapter bridges them
type SQLiteAdapter struct {
    db *SQLiteDB
}
func (a *SQLiteAdapter) Query(sql string, args ...interface{}) (Result, error) {
    return a.db.Execute(sql)
}
```

### 2. Message Queue Adapter

```go
// Client interface
type MessageQueue interface {
    Publish(message string) error
}

// Support multiple backends
type RabbitMQAdapter struct { /*...*/ }
type KafkaAdapter struct { /*...*/ }
```

### 3. Payment Gateway Adapter

```go
// Stable internal interface
type PaymentProcessor interface {
    Process(payment Payment) (receipt Receipt, err error)
}

// Different external libraries
type StripeAdapter struct { /*...*/ }
type PayPalAdapter struct { /*...*/ }
type SquareAdapter struct { /*...*/ }
```

## Key Advantages

- **Integration**: Enables use of incompatible external libraries
- **Decoupling**: Internal code independent of external interface changes
- **Flexibility**: Switch implementations without changing client code
- **Stability**: Provides stable internal API hiding external complexity
- **Testing**: Easy to mock adapted interfaces for testing
- **Single Responsibility**: Adapter focuses solely on translation

## Key Gotchas

- **Over-adaptation**: Creating adapters for compatible interfaces adds unnecessary complexity
- **Feature loss**: Adapter may not support all external library features
- **Performance overhead**: Each call goes through adapter translation layer
- **Maintenance burden**: Changes to external API require adapter updates
- **Leaky abstraction**: Adapter may expose underlying implementation details
- **Documentation**: Must document what adapter does and doesn't support
- **Complexity**: Too many adapters can obscure the actual implementation
