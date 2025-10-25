# Builder Pattern

## Overview

The Builder pattern separates the construction of a complex object from its representation, allowing the same construction process to create different representations. It's especially useful when an object has many optional parameters or requires step-by-step construction.

## Problem

You're building a system that creates complex objects like HTTP requests, database queries, or UI components. These objects have:
- Many optional parameters (10+ fields)
- Complex initialization logic
- Order-dependent construction steps
- Validation requirements

Without the Builder pattern, you face:
```go
// Telescoping constructor anti-pattern
func NewHTTPRequest(url, method, body, contentType, authToken,
    userAgent, timeout, retryCount, followRedirects, ...20 more params) *Request {
    // Unmanageable!
}

// Or massive struct initialization
req := &Request{
    URL: "...",
    Method: "POST",
    Body: "...",
    Headers: map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer ...",
        // ... many more
    },
    // ... 20 more fields
}
```

You need a clean, readable way to construct complex objects step-by-step with optional parameters.

## Why Use This Pattern?

- **Readable Construction**: Fluent interface makes object creation self-documenting
- **Optional Parameters**: Handle many optional fields without telescoping constructors
- **Immutability**: Build mutable object, return immutable result
- **Validation**: Validate during construction before object creation
- **Step-by-Step**: Complex construction broken into manageable steps
- **Multiple Representations**: Same builder process, different products

## When to Use

- **Many Parameters**: Objects with 4+ optional parameters
- **Complex Construction**: Multi-step initialization process
- **Immutable Objects**: Build mutable, return immutable
- **Configuration Objects**: API clients, database connections, HTTP requests
- **Document Builders**: SQL queries, XML/JSON builders, HTML generation
- **Test Data Builders**: Creating test fixtures with many variations

## When NOT to Use

- **Simple Objects**: Few required parameters, no optional ones
- **No Variation**: Construction process never varies
- **Performance Critical**: Builder adds slight overhead
- **Overkill**: When simple struct initialization suffices

## Implementation Guidelines

### Basic Builder Pattern in Go

```go
// 1. Product - the complex object
type House struct {
    foundation string
    structure  string
    roof       string
    interior   string
}

// 2. Builder - constructs the product
type HouseBuilder struct {
    house House
}

func NewHouseBuilder() *HouseBuilder {
    return &HouseBuilder{}
}

func (b *HouseBuilder) Foundation(foundation string) *HouseBuilder {
    b.house.foundation = foundation
    return b // Return self for chaining
}

func (b *HouseBuilder) Structure(structure string) *HouseBuilder {
    b.house.structure = structure
    return b
}

func (b *HouseBuilder) Build() House {
    return b.house
}

// Usage
house := NewHouseBuilder().
    Foundation("concrete").
    Structure("wood").
    Build()
```

### Functional Options Pattern (Go Idiom)

```go
type Server struct {
    host string
    port int
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080} // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(WithHost("0.0.0.0"), WithPort(9000))
```

## Go Idioms

- **Method Chaining**: Return `*Builder` from each method
- **Functional Options**: Use for simple cases with optional params
- **Validation in Build()**: Validate before returning product
- **Director Optional**: Often omit director in Go, builder is enough
- **Pointer Receivers**: Use pointer receivers for builder methods

## Visual Schema

```
Client
  │
  │ creates
  ↓
┌─────────────────────────┐
│  Builder                │
│                         │
│  - product: Product     │
│                         │
│  + SetPartA()  ────────┼─→ builds product incrementally
│  + SetPartB()  ────────┼─→
│  + SetPartC()  ────────┼─→
│  + Build(): Product ───┼─→ returns completed product
└─────────────────────────┘
         │
         │ constructs
         ↓
    ┌──────────┐
    │  Product │
    └──────────┘
```

### Fluent Interface Flow

```
NewBuilder()
    │
    ├─→ SetField1(val) ──→ Builder
    │        │
    │        ├─→ SetField2(val) ──→ Builder
    │        │        │
    │        │        ├─→ SetField3(val) ──→ Builder
    │        │        │        │
    │        │        │        └─→ Build() ──→ Product
    │        │        │
    │        │        └─→ Build() ──→ Product
    │        │
    │        └─→ Build() ──→ Product (with defaults)
    │
    └─→ Build() ──→ Product (all defaults)

Method chaining enables flexible, readable construction
```

## Real-World Examples

### 1. HTTP Request Builder
```go
req := NewRequestBuilder().
    URL("https://api.example.com/users").
    Method("POST").
    Header("Content-Type", "application/json").
    Body(`{"name": "John"}`).
    Timeout(30 * time.Second).
    Build()
```

### 2. SQL Query Builder
```go
query := NewQueryBuilder().
    Select("id", "name", "email").
    From("users").
    Where("age > ?", 18).
    OrderBy("name ASC").
    Limit(10).
    Build()
```

### 3. Configuration Builder
```go
config := NewConfigBuilder().
    DatabaseURL("postgres://localhost/mydb").
    MaxConnections(100).
    Timeout(30 * time.Second).
    EnableLogging(true).
    Build()
```

## Key Advantages

✓ **Readability**: Self-documenting, fluent construction code
✓ **Flexibility**: Easy to add/modify optional parameters
✓ **Immutability**: Build complex objects that become immutable
✓ **Validation**: Centralized validation before object creation
✓ **Defaults**: Easy to provide sensible default values
✓ **Testability**: Easy to create test fixtures with variations

## Key Gotchas

⚠️ **Verbosity**: More code than simple struct initialization
⚠️ **Mutability**: Builder itself is mutable (by design)
⚠️ **Nil Pointers**: Builder methods should handle nil gracefully
⚠️ **Required Fields**: Must validate required fields in Build()
⚠️ **Thread Safety**: Builders are typically not thread-safe
⚠️ **Copies vs References**: Be careful with pointer fields

## Best Practices

1. **Return Self**: Builder methods return `*Builder` for chaining
2. **Validate in Build()**: Check required fields, business rules
3. **Provide Defaults**: Set sensible defaults in constructor
4. **Immutable Product**: Make the built object immutable
5. **Clear Naming**: Use verb names (Set, With, Add, Enable)
6. **Document Required Fields**: Clearly state what's mandatory
7. **Error Handling**: Return errors from Build() for validation failures
8. **Copy Don't Share**: Return copies of complex fields, not references
9. **One Builder Per Product**: Don't reuse builders across products
10. **Consider Functional Options**: For simpler cases, use functional options pattern

## Builder vs Functional Options

**Use Builder when:**
- Complex multi-step construction
- Need to build multiple similar objects
- Construction logic is complex
- Want explicit Build() step

**Use Functional Options when:**
- Simple configuration
- Most parameters have good defaults
- No complex validation
- Prefer concise code
