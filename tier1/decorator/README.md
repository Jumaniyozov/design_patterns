# Decorator Pattern

## Overview

The Decorator Pattern is a structural design pattern that lets you attach new behaviors to objects by placing these objects inside special wrapper objects that contain the behaviors. It provides a flexible alternative to subclassing for extending functionality.

## Problem

When you need to add functionality to objects dynamically, you face several challenges:

- **Subclassing explosion**: Creating a subclass for each behavior combination leads to exponential class growth
- **Static composition**: Traditional inheritance is rigid and determined at compile time
- **Cross-cutting concerns**: Adding logging, caching, or validation to existing objects requires modification
- **Single Responsibility Violation**: Objects accumulate unrelated responsibilities

### Real-World Context

Consider a data processing pipeline where you need to add features like compression, encryption, logging, validation, and caching to your data. With inheritance, you'd need 2^5 = 32 different classes for all combinations. Decorators let you compose these features cleanly.

## Why Use This Pattern?

- **Flexibility**: Add behavior at runtime without modifying original objects
- **Composability**: Stack decorators in any order to create complex behavior
- **Single Responsibility**: Each decorator adds one specific concern
- **Avoids class explosion**: No need for exponential subclass combinations
- **Go Idioms**: Functions and closures make decorators particularly elegant in Go

## When to Use

- Adding responsibilities to objects dynamically
- Extending behavior through composition, not inheritance
- Multiple independent features that can be combined
- Middleware and cross-cutting concerns (logging, caching, validation)
- HTTP handlers and request processing chains
- Data processing pipelines

## When NOT to Use

- Simple single feature to add (direct composition may be cleaner)
- Performance is critical and wrapper overhead matters
- Few combinations of behaviors needed
- Deep decorator chains become hard to debug
- When interface changes are needed

## Implementation Guidelines

1. **Component interface**: Define the contract being decorated
2. **Concrete component**: The original object being enhanced
3. **Decorator base**: Wraps the component and satisfies the same interface
4. **Concrete decorators**: Add specific behaviors without modifying the component
5. **Ordering matters**: Stack decorators in the intended execution order

## Go Idioms

Go's first-class functions and closures make decorators particularly natural:

- Function types can be decorated like objects
- Closures capture state elegantly
- Middleware pattern is a functional decorator
- No need for base decorator classes; just compose functions

## Visual Schema

```go
┌────────────────────────────────────────────────┐
│        Decorator Pattern Structure             │
└────────────────────────────────────────────────┘

Component Interface:
  ┌──────────────────┐
  │   Operation()    │
  └──────────────────┘

Inheritance Approach (Explosion):
  ConcreteComponent
         │
    ┌────┴─────────┬────────────────┬─────────┐
    │              │                │         │
  WithA          WithB          WithA+B     WithC
    │              │                │         │
    ├──────┬───────┼────────────────┼─────────┤
   A+C    B+C    A+B+C             ...    [32 classes]

Decorator Approach (Composition):
  ConcreteComponent: Component Interface
        ↑
  DecoratorA: wraps Component
        ↑
  DecoratorB: wraps DecoratorA
        ↑
  DecoratorC: wraps DecoratorB

Client can compose any combination without subclassing!

Functional Decorator (Go Style):
  func MakeReader(r io.Reader) io.Reader {
    return &DecoratorA{
      reader: &DecoratorB{
        reader: &DecoratorC{
          reader: r,
        },
      },
    }
  }
```

## Real-World Examples

### 1. HTTP Handler Middleware

```go
// Middleware decorator chain
handler := loggingMiddleware(
    authMiddleware(
        metricsMiddleware(
            originalHandler,
        ),
    ),
)
```

### 2. Data Stream Processing

```go
stream := compressDecorator(
    encryptDecorator(
        validateDecorator(
            logDecorator(
                originalStream,
            ),
        ),
    ),
)
```

### 3. File Reader with Features

```go
reader := NewBufferedReader(
    NewEncryptedReader(
        NewValidatingReader(
            NewLoggingReader(
                originalFile,
            ),
        ),
    ),
)
```

## Key Advantages

- **Open/Closed Principle**: Open for extension, closed for modification
- **Runtime flexibility**: Add behavior dynamically as needed
- **Clean composition**: No side effects from multiple inheritance
- **Testability**: Decorators can be tested independently
- **Middleware pattern**: Natural fit for cross-cutting concerns

## Key Gotchas

- **Order matters**: Decorator order affects behavior and performance
- **Deep stacks**: Multiple layers can become hard to debug and understand
- **Interface compliance**: All decorators must satisfy the component interface
- **Type assertions**: May lose access to concrete type methods
- **Performance**: Each decorator adds a small overhead
- **Complexity**: Too many decorators can make code harder to follow
