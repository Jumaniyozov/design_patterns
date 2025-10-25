# Dependency Injection Pattern

## Overview
Dependency Injection (DI) inverts control of dependency creation by providing (injecting) dependencies to an object rather than having it create them. This promotes loose coupling, testability, and flexibility.

## Problem
When objects create their own dependencies (e.g., `db := sql.Open(...)`), they become tightly coupled to specific implementations. This makes testing difficult (can't mock dependencies), violates Single Responsibility (object manages dependencies), and reduces flexibility (can't swap implementations).

## Why Use This Pattern?

**Benefits:**
- **Testability**: Easy to inject mocks/stubs for testing
- **Loose Coupling**: Components don't depend on concrete implementations
- **Flexibility**: Swap implementations without changing code
- **Single Responsibility**: Objects focus on their job, not dependency management
- **Configuration**: Centralized dependency wiring

## When to Use

Use DI when:
- You need to test with mocks
- Multiple implementations exist
- Dependencies should be configurable
- You want loose coupling between components

**Real-world scenarios:**
- Database connections
- HTTP clients
- Logging services
- Configuration providers
- External API clients

## When NOT to Use

Avoid when:
- Dependencies are stable and unchanging
- Application is very simple
- Overhead outweighs benefits

## Go Idioms

Go favors explicit DI through:
1. **Constructor Injection**: Pass dependencies to constructors
2. **Interface Injection**: Accept interfaces, return structs
3. **Struct Embedding**: Compose dependencies

```go
type Service struct {
    db     Database    // Interface
    logger Logger      // Interface
}

func NewService(db Database, logger Logger) *Service {
    return &Service{db: db, logger: logger}
}
```

## Code Example Structure

Demonstrates:
- Constructor injection
- Interface-based DI
- Service container pattern
- Wire-style DI framework
