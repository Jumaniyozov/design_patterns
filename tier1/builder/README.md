# Builder Pattern

## Overview

The Builder Pattern is a creational design pattern that separates the construction of a complex object from its representation. It allows you to build objects step-by-step, making it ideal for objects with many optional fields or complex initialization logic.

## Problem

When creating objects with many optional parameters or complex initialization logic, you face several challenges:

- **Constructor overloading hell**: Go doesn't support method overloading, making it difficult to provide multiple ways to create an object
- **Optional parameters**: Passing many nil or zero values makes code unreadable
- **Partial initialization**: You may want to specify only some fields, leaving others with defaults
- **Complex validation**: Building an object may require multiple steps or validation between steps

### Real-World Context

Consider creating a database connection pool with many optional settings: connection timeout, read timeout, write timeout, max idle connections, max open connections, SSL configuration, and retry policy. Creating this with a single constructor function becomes unwieldy and error-prone.

## Why Use This Pattern?

- **Readability**: Fluent API makes construction code self-documenting
- **Flexibility**: Easy to add new optional fields without breaking existing code
- **Validation**: Can validate at each step or at the end
- **Immutability**: Often used to create immutable objects by validating before final construction
- **Go Idioms**: Aligns with Go's functional options approach and builder conventions

## When to Use

- Objects with many optional fields (3 or more)
- Complex initialization logic or multi-step setup
- Need for readable, maintainable construction code
- Immutable objects that require validation
- Creating objects with optional configuration

## When NOT to Use

- Simple objects with 1-2 required fields
- Rare object creation where simplicity outweighs benefits
- Performance-critical code where builder overhead is unacceptable
- When the Options Pattern (functional options) is more suitable

## Implementation Guidelines

1. **Builder struct**: Hold mutable state during construction
2. **Setter methods**: Return the builder to enable method chaining
3. **Build method**: Perform final validation and return the constructed object
4. **Separation**: Keep the builder logic separate from the main object
5. **Validation**: Consider when to validate (eagerly or at build time)

## Go Idioms

Go developers often prefer the **Functional Options** pattern over traditional builders. However, builders remain useful for:

- Complex multi-step initialization
- Immutable value objects with many fields
- Clear visual separation between building and using

## Visual Schema

```go
┌─────────────────────────────────────────────┐
│         Builder Pattern Flow                │
└─────────────────────────────────────────────┘

  Client                Builder              Product
    │                     │                    │
    ├─ New Builder ──────>│                    │
    │                     │                    │
    ├─ WithField1 ───────>│ (store field1)     │
    │                     │                    │
    ├─ WithField2 ───────>│ (store field2)     │
    │                     │                    │
    ├─ WithField3 ───────>│ (store field3)     │
    │                     │                    │
    ├─── Build ──────────>│ (validate all)     │
    │                     │                    │
    │                     ├─ return Product ──>│
    │<─────────────────────────────────────────┤
    │           (Complex Object)               │

Problem vs Solution:

BEFORE (Hard to read):
  conn := NewConnection(timeout, readTimeout, writeTimeout,
                        maxIdle, maxOpen, useSSL, retries)

AFTER (Clear & Maintainable):
  conn := NewConnectionBuilder().
    WithTimeout(30 * time.Second).
    WithMaxConnections(10).
    WithSSL(true).
    Build()
```

## Real-World Examples

### 1. Database Connection Pool Builder

```go
pool := NewDatabasePoolBuilder().
    WithHost("localhost").
    WithPort(5432).
    WithMaxConnections(20).
    WithConnectionTimeout(30 * time.Second).
    WithSSLMode("require").
    Build()
```

### 2. HTTP Request Builder

```go
request := NewHTTPRequestBuilder().
    WithMethod("POST").
    WithURL("https://api.example.com/data").
    WithHeader("Authorization", "Bearer token").
    WithBody(payload).
    WithTimeout(10 * time.Second).
    Build()
```

### 3. Configuration Object Builder

```go
config := NewAppConfigBuilder().
    WithServerPort(8080).
    WithDatabase("postgres://localhost/mydb").
    WithLogLevel("info").
    WithMaxWorkers(4).
    WithEnableMetrics(true).
    Build()
```

## Key Advantages

- **Fluent API**: Method chaining creates readable, expressive code
- **Flexibility**: Easily make fields optional without breaking changes
- **Validation**: Centralized validation logic at build time
- **Immutability**: Construct immutable objects safely
- **Clarity**: Intent is clear when reading construction code

## Key Gotchas

- **Mutable builder**: Builder is mutable during construction; use once and discard
- **Nil safety**: Always handle nil checks in Build() method
- **Chaining complexity**: Long chains can become harder to format and read
- **Performance**: Creating builder object adds overhead (usually negligible)
- **Thread safety**: Builders are not thread-safe; don't share across goroutines
