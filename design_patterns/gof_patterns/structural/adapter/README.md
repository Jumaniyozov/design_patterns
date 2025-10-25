# Adapter Pattern

## Overview

The Adapter pattern converts the interface of a class into another interface clients expect. It allows classes with incompatible interfaces to work together by wrapping one interface with another. In Go, this is typically implemented through composition and interface implementation.

## Problem

You're integrating with multiple third-party payment gateways: Stripe, PayPal, and Square. Each has a different API:

```go
// Stripe
stripe.ProcessPayment(amount, currency, cardToken)

// PayPal
paypal.MakePayment(price, currencyCode, paymentMethod)

// Square
square.Charge(moneyAmount, paymentSource, currencyType)
```

Your application expects a consistent interface:
```go
type PaymentProcessor interface {
    Pay(amount float64, currency string) error
}
```

Without adapters, you'd need conditional logic everywhere:
```go
if gateway == "stripe" {
    stripe.ProcessPayment(...)
} else if gateway == "paypal" {
    paypal.MakePayment(...)
} // Scattered throughout codebase
```

You need a way to make incompatible interfaces compatible without modifying the original code.

## Why Use This Pattern?

- **Interface Compatibility**: Makes incompatible interfaces work together
- **Code Reuse**: Use existing classes without modification
- **Decoupling**: Isolates client code from third-party implementations
- **Single Responsibility**: Adapter handles translation logic separately
- **Open/Closed Principle**: Add new adapters without changing existing code

## When to Use

- **Third-Party Integration**: Wrapping external libraries with inconsistent APIs
- **Legacy Code**: Making old code work with new interfaces
- **Multiple Implementations**: Standardizing diverse implementations
- **API Versioning**: Supporting multiple API versions simultaneously
- **Database Drivers**: Uniform interface across different databases
- **Cloud Providers**: Consistent API for AWS, Azure, GCP services

## When NOT to Use

- **Identical Interfaces**: When interfaces already match
- **Simple Wrappers**: When a simple function wrapper suffices
- **Overengineering**: For one-time, simple conversions
- **Performance Critical**: Adds indirection (usually negligible)

## Implementation Guidelines

### Basic Adapter Pattern in Go

```go
// Target interface (what your application expects)
type Target interface {
    Request() string
}

// Adaptee (incompatible third-party code)
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Specific request from adaptee"
}

// Adapter (makes Adaptee compatible with Target)
type Adapter struct {
    adaptee *Adaptee
}

func (a *Adapter) Request() string {
    // Translate the interface
    return a.adaptee.SpecificRequest()
}
```

### Object Adapter (Composition)
```go
type Adapter struct {
    adaptee Adaptee // Composition
}
```

### Class Adapter (Embedding - Go's inheritance alternative)
```go
type Adapter struct {
    *Adaptee // Embedding
}

func (a *Adapter) Request() string {
    return a.SpecificRequest() // Direct access
}
```

## Go Idioms

- **Prefer Composition**: Use struct fields over embedding for clarity
- **Interface-Based Design**: Define target interfaces in your package
- **Error Handling**: Adapt error types and handling strategies
- **Constructor Functions**: `NewXAdapter()` for initialization
- **Thin Adapters**: Keep translation logic simple; delegate to adaptee

## Visual Schema

```
Client
  │
  │ expects Target interface
  │
  ↓
┌─────────────────────────┐
│  Target Interface       │
│  ├─ Method1()           │
│  └─ Method2()           │
└─────────────────────────┘
          ↑
          │ implements
          │
┌─────────────────────────┐
│  Adapter                │
│                         │
│  Contains: Adaptee      │
│                         │
│  Method1() {            │
│    return adaptee       │
│      .DifferentMethod() │
│  }                      │
└─────────────────────────┘
          │
          │ uses
          ↓
┌─────────────────────────┐
│  Adaptee                │
│  (Third-party code)     │
│                         │
│  ├─ DifferentMethod()   │
│  └─ AnotherMethod()     │
└─────────────────────────┘
```

### Before vs After

```
BEFORE (Direct Dependencies):

Client Code
  ├─→ if Stripe: call stripe.ProcessPayment()
  ├─→ if PayPal: call paypal.MakePayment()
  └─→ if Square: call square.Charge()

  Problem: Tightly coupled, repetitive, hard to test


AFTER (With Adapters):

Client Code → PaymentProcessor Interface
                      ↑
                      │
              ┌───────┼───────┐
              │       │       │
          StripeAdapter  PayPalAdapter  SquareAdapter
              │       │       │
              ↓       ↓       ↓
           Stripe   PayPal  Square

  Solution: Loose coupling, uniform interface, testable
```

## Real-World Examples

### 1. Payment Gateway Adapters
```go
type PaymentProcessor interface {
    Pay(amount float64, currency string) error
}

type StripeAdapter struct {
    client *stripe.Client
}

func (s *StripeAdapter) Pay(amount float64, currency string) error {
    return s.client.ProcessPayment(int(amount*100), currency, token)
}
```

### 2. Database Driver Adapters
```go
type Database interface {
    Query(sql string) ([]Row, error)
}

type PostgresAdapter struct {
    conn *pgx.Conn
}

func (p *PostgresAdapter) Query(sql string) ([]Row, error) {
    rows, err := p.conn.Query(context.Background(), sql)
    // Adapt pgx.Rows to []Row
    return adaptRows(rows), err
}
```

### 3. Logging Library Adapters
```go
type Logger interface {
    Info(msg string)
    Error(msg string)
}

type ZapAdapter struct {
    logger *zap.Logger
}

func (z *ZapAdapter) Info(msg string) {
    z.logger.Info(msg)
}

func (z *ZapAdapter) Error(msg string) {
    z.logger.Error(msg)
}
```

## Key Advantages

✓ **Compatibility**: Integrate incompatible interfaces seamlessly
✓ **Reusability**: Reuse existing code without modification
✓ **Flexibility**: Swap implementations easily
✓ **Testability**: Mock adapters for unit testing
✓ **Separation of Concerns**: Translation logic isolated in adapters
✓ **Open/Closed**: Extend with new adapters without modifying existing code

## Key Gotchas

⚠️ **Complexity**: Adds extra layer of indirection
⚠️ **Performance**: Minimal overhead from additional method calls
⚠️ **Over-adaptation**: Don't adapt what already fits
⚠️ **Leaky Abstractions**: Ensure adapter fully hides adaptee details
⚠️ **Error Translation**: Properly map errors between systems
⚠️ **State Management**: Be careful with stateful adaptees

## Best Practices

1. **Define Clear Interfaces**: Your target interface should be well-designed
2. **Keep Adapters Simple**: Only translate interfaces, don't add business logic
3. **One Responsibility**: One adapter per adaptee type
4. **Document Mappings**: Clearly document how methods/parameters map
5. **Handle Errors**: Properly translate error types and messages
6. **Test Thoroughly**: Test adapter behavior, not just interface compliance
7. **Version Control**: Consider versioning when adapting external APIs
8. **Use Factories**: Create adapters through factory functions for consistency
