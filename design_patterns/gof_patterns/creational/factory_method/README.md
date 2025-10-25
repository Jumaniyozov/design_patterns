# Factory Method Pattern

## Overview

The Factory Method pattern defines an interface for creating objects, but lets subclasses or implementations decide which class to instantiate. In Go, this is achieved through functions that return interface types, allowing the factory to determine the concrete type based on input parameters.

## Problem

You're building a notification system that sends messages through different channels: Email, SMS, and Push notifications. Each channel has different:
- Initialization requirements (SMTP servers, SMS gateways, push services)
- Message formatting rules
- Delivery mechanisms
- Error handling

Hardcoding the creation logic spreads `if-else` or `switch` statements throughout your code:
```go
var notifier Notifier
if channel == "email" {
    notifier = &EmailNotifier{smtpServer: "..."}
} else if channel == "sms" {
    notifier = &SMSNotifier{gateway: "..."}
} // ... repeated everywhere
```

You need a centralized way to create objects that encapsulates the instantiation logic and allows for easy extension.

## Why Use This Pattern?

- **Encapsulation**: Object creation logic is centralized in one place
- **Flexibility**: Easy to add new types without modifying existing code
- **Decoupling**: Client code depends on interfaces, not concrete types
- **Testability**: Easy to inject mock implementations for testing
- **Single Responsibility**: Creation logic separated from business logic

## When to Use

- **Multiple Implementations**: When you have several types implementing the same interface
- **Runtime Decisions**: Object type determined by runtime parameters (config, user input, environment)
- **Plugin Systems**: Loading different implementations based on configuration
- **Database Drivers**: Creating connections to different database types
- **File Parsers**: Instantiating parsers based on file extension (JSON, XML, YAML)
- **Payment Gateways**: Creating payment processors for different providers

## When NOT to Use

- **Single Implementation**: When there's only one concrete type
- **Simple Construction**: When object creation is trivial (just `&Type{}`)
- **No Variation**: When the creation process never changes
- **Overengineering**: Adding unnecessary abstraction for simple cases

## Implementation Guidelines

### Basic Pattern in Go

```go
// 1. Define the interface
type Product interface {
    Use() string
}

// 2. Implement concrete types
type ConcreteProductA struct{}
func (p *ConcreteProductA) Use() string { return "Product A" }

type ConcreteProductB struct{}
func (p *ConcreteProductB) Use() string { return "Product B" }

// 3. Factory function
func NewProduct(productType string) Product {
    switch productType {
    case "A":
        return &ConcreteProductA{}
    case "B":
        return &ConcreteProductB{}
    default:
        return nil
    }
}
```

### Advanced: Registry Pattern

```go
type ProductFactory func() Product

var registry = make(map[string]ProductFactory)

func Register(name string, factory ProductFactory) {
    registry[name] = factory
}

func Create(name string) Product {
    if factory, ok := registry[name]; ok {
        return factory()
    }
    return nil
}
```

## Go Idioms

- **Constructor Functions**: Use `New...()` naming convention for factory functions
- **Return Interfaces**: Factory functions return interface types, not concrete types
- **Variadic Options**: Combine with functional options pattern for flexible construction
- **Error Handling**: Return `(Product, error)` for robust error handling
- **Package Design**: Keep factory and interfaces in one package, implementations can be separate

## Visual Schema

```
Client Code
    │
    │ calls NewNotifier("email", config)
    │
    ↓
┌─────────────────────────────┐
│   Factory Function          │
│                             │
│  func NewNotifier(          │
│      type string,           │
│      config Config          │
│  ) Notifier                 │
│                             │
│  switch type {              │
│    case "email":            │
│      return &EmailNotifier  │
│    case "sms":              │
│      return &SMSNotifier    │
│    case "push":             │
│      return &PushNotifier   │
│  }                          │
└─────────────────────────────┘
    │
    ├─→ EmailNotifier ──┐
    ├─→ SMSNotifier    ─┼→ All implement Notifier interface
    └─→ PushNotifier   ─┘

                        ↓
                  Notifier Interface
                  ├─ Send(msg string) error
                  └─ GetStatus() string
```

### Before vs After

```
BEFORE (Scattered Creation Logic):

┌──────────────────┐         ┌──────────────────┐
│   ServiceA       │         │   ServiceB       │
│                  │         │                  │
│  if type == "X"  │         │  if type == "X"  │
│    create X      │         │    create X      │
│  else            │         │  else            │
│    create Y      │         │    create Y      │
└──────────────────┘         └──────────────────┘
  Duplication!                  Duplication!

AFTER (Centralized Factory):

┌──────────────────┐         ┌──────────────────┐
│   ServiceA       │         │   ServiceB       │
│                  │         │                  │
│  obj := Factory  │         │  obj := Factory  │
│    .Create(type) │         │    .Create(type) │
└──────────────────┘         └──────────────────┘
         │                            │
         └────────────┬───────────────┘
                      ↓
              ┌──────────────────┐
              │  Factory         │
              │  (Single Source) │
              └──────────────────┘
```

## Real-World Examples

### 1. Database Connection Factory
```go
func NewDatabase(driver string, connString string) (Database, error) {
    switch driver {
    case "postgres":
        return NewPostgresDB(connString)
    case "mysql":
        return NewMySQLDB(connString)
    case "mongodb":
        return NewMongoDB(connString)
    default:
        return nil, fmt.Errorf("unknown driver: %s", driver)
    }
}
```

### 2. Document Parser Factory
```go
func NewParser(filename string) (Parser, error) {
    ext := filepath.Ext(filename)
    switch ext {
    case ".json":
        return &JSONParser{}, nil
    case ".xml":
        return &XMLParser{}, nil
    case ".yaml", ".yml":
        return &YAMLParser{}, nil
    default:
        return nil, fmt.Errorf("unsupported format: %s", ext)
    }
}
```

### 3. HTTP Client Factory
```go
func NewHTTPClient(env string) HTTPClient {
    switch env {
    case "production":
        return &ProductionClient{timeout: 30 * time.Second}
    case "development":
        return &DevelopmentClient{timeout: 5 * time.Minute, verbose: true}
    case "test":
        return &MockClient{}
    default:
        return &DefaultClient{}
    }
}
```

## Key Advantages

✓ **Centralized Creation Logic**: All instantiation in one place
✓ **Easy Extension**: Add new types by extending factory, not modifying clients
✓ **Type Safety**: Compiler ensures all products implement the interface
✓ **Dependency Inversion**: Clients depend on abstractions (interfaces)
✓ **Testability**: Easy to swap implementations for testing

## Key Gotchas

⚠️ **God Factory**: Don't create one factory for unrelated types
⚠️ **Tight Coupling to Strings**: String-based switching is fragile; consider type parameters
⚠️ **Nil Returns**: Always check for nil or use error returns
⚠️ **Over-abstraction**: Don't use for simple object creation
⚠️ **Registration Timing**: With registry pattern, ensure registrations happen before first use
⚠️ **Thread Safety**: If factory maintains state, ensure thread-safe access

## Best Practices

1. **Return Interfaces**: `func NewX() Interface`, not `func NewX() *ConcreteType`
2. **Error Handling**: Return errors for invalid inputs: `(Product, error)`
3. **Validate Early**: Check parameters before attempting creation
4. **Document Types**: Clearly document supported types and their parameters
5. **Use Enums**: Instead of strings, use constants or enums for type selection
6. **Naming Convention**: Prefix with `New` (NewParser, NewDatabase, NewClient)
7. **Single Responsibility**: One factory per interface/product family
