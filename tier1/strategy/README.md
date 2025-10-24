# Strategy Pattern

## Overview

The Strategy Pattern encapsulates a family of algorithms, making them interchangeable. It lets the algorithm vary independently from clients that use it. Instead of implementing multiple algorithms directly within a class or using conditional logic, you define each algorithm as a separate strategy and allow the context to select which strategy to use at runtime.

## Problem

Imagine you're building a payment processing system that needs to support multiple payment methods: credit cards, PayPal, cryptocurrency, and bank transfers. Each has completely different logic, validation rules, and processing flows.

Without the Strategy pattern, you might end up with massive conditional statements:

```go
func ProcessPayment(method string, amount float64, details map[string]string) error {
    if method == "credit_card" {
        // 50 lines of credit card validation and processing
    } else if method == "paypal" {
        // 40 lines of PayPal OAuth and confirmation logic
    } else if method == "crypto" {
        // 60 lines of blockchain monitoring and confirmation
    }
    // ... the nightmare continues
}
```

This approach violates **Single Responsibility Principle** (each payment method isn't cohesive), **Open/Closed Principle** (you must modify existing code to add new methods), and makes testing nearly impossible since all logic is tightly coupled.

## Why Use This Pattern?

1. **Eliminates Conditional Logic** - Replace if/switch statements with polymorphism
2. **Single Responsibility** - Each strategy handles one algorithm completely
3. **Open/Closed Principle** - Add new strategies without modifying existing code
4. **Runtime Flexibility** - Switch algorithms based on runtime conditions (user choice, configuration, data characteristics)
5. **Testability** - Test each algorithm independently without dependencies
6. **Code Reusability** - Strategies can be shared across different contexts
7. **Clear Intent** - Code explicitly shows that multiple approaches are available

## When to Use

✅ **Use Strategy when:**

- You have multiple ways to accomplish the same task
- You want to avoid conditional logic (if/switch statements)
- The algorithm might change at runtime based on user choice or system state
- You want to test each algorithm independently
- New algorithms will be added frequently (extensibility requirement)
- Different clients need different implementations of the same concept
- You're implementing variant behavior in a clean, decoupled way

❌ **When NOT to Use:**

- You only have one or two simple alternatives (overengineering)
- The logic is tightly coupled to specific data structures
- The overhead of an interface isn't justified by complexity (simple conditionals are fine)
- The algorithms are rarely changed or added
- All variants will only ever exist in one place in the codebase

## Implementation Guidelines

### 1. Define a Clear Interface

The strategy interface should represent the core operation(s) needed. Keep it focused:

```go
type Strategy interface {
    Execute(params interface{}) (result interface{}, err error)
}
```

For domain-specific patterns, be more explicit:

```go
type SortingStrategy interface {
    Compare(a, b interface{}) int
    Swap(slice interface{}, i, j int)
}
```

### 2. Implement Concrete Strategies

Each concrete strategy should be a complete, standalone implementation:

- Keep each strategy in its own type
- Don't share state between strategies unless necessary
- Make each strategy independently testable

### 3. Create a Context

The context holds a reference to the strategy and uses it:

- The context shouldn't know implementation details of strategies
- The context should be able to switch strategies if needed
- Encapsulate strategy selection logic if complex

### 4. Configure Strategy Selection

Choose how clients select strategies:

- **Constructor injection** (most common): Pass strategy when creating context
- **Setter method**: Allow changing strategy after creation
- **Factory function**: Centralize strategy selection logic
- **Configuration-driven**: Load strategy based on config files

### 5. Best Practices

- **Make strategies stateless when possible** - Easier to reuse and thread-safe
- **Use dependency injection** - Pass dependencies to strategy, don't have it access globals
- **Name strategies clearly** - `CreditCardPayment`, `PayPalPayment` are better than `Strategy1`, `Strategy2`
- **Document the interface contract** - What parameters, what guarantees, what can fail
- **Consider strategy composition** - Strategies can use other patterns (decorator, chain of responsibility)

## Visual Schemas

### Pattern Structure

The Strategy pattern consists of these key components and their relationships:

```go
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT CODE                             │
└────────────────────────┬────────────────────────────────────────┘
                         │ uses
                         ▼
        ┌────────────────────────────┐
        │   PaymentProcessor (Ctx)   │
        │                            │
        │  - strategy: Strategy      │
        │  + Process(amount)         │
        │  + SetStrategy(s)          │
        └───────────┬────────────────┘
                    │ delegates to
                    ▼
        ┌───────────────────────────────────────┐
        │    <<interface>>                      │
        │    PaymentStrategy                    │
        ├───────────────────────────────────────┤
        │ + Validate(details)                   │
        │ + Process(amount) -> txID             │
        │ + Refund(txID, amount)                │
        │ + GetName()                           │
        └───────────────────────────────────────┘
                    ▲
        ┌───────────┼────────────┬─────────────────┐
        │           │            │                 │
        │           │            │                 │
    ┌───┴───────┐ ┌─┴─────────┐ ┌┴──────────────┐ ...
    │ CreditCard│ │PayPal     │ │Cryptocurrency │
    │ Strategy  │ │ Strategy  │ │Strategy       │
    ├───────────┤ ├───────────┤ ├───────────────┤
    │ -cardNum  │ │ -email    │ │ -walletAddr   │
    │ -expiry   │ │           │ │ -cryptoType   │
    │ -cvv      │ │           │ │               │
    └───────────┘ └───────────┘ └───────────────┘
```

### Workflow: How Strategy Pattern Eliminates Conditionals

**WITHOUT Strategy Pattern (Conditional Nightmare):**

```go
Client Code
    │
    ├─ paymentMethod == "credit_card" ?
    │  └─ Run 50 lines of credit card logic
    │
    ├─ paymentMethod == "paypal" ?
    │  └─ Run 40 lines of PayPal logic
    │
    ├─ paymentMethod == "crypto" ?
    │  └─ Run 60 lines of crypto logic
    │
    └─ Add new method? Modify this code!
```

**WITH Strategy Pattern (Clean & Extensible):**

```go
Client Code
    │
    ├─ Select Strategy
    │  (CreditCardStrategy, PayPalStrategy, etc.)
    │
    ├─ Create Context with Strategy
    │  PaymentProcessor(strategy)
    │
    └─ Use Context
       processor.Process(amount)
       └─ Context delegates to Strategy
          └─ Strategy handles its own algorithm

    Add new method? Create new Strategy!
    No changes to existing code needed.
```

### Execution Flow

```go
process(amount)
    │
    ▼
┌─────────────────────────────────────────┐
│ PaymentProcessor.Process()              │
├─────────────────────────────────────────┤
│ 1. Validate strategy is set             │
│ 2. Call strategy.Validate(details)      │
│    ├─ CreditCard: Check card format     │
│    ├─ PayPal: Check email format        │
│    └─ Crypto: Check wallet length       │
│ 3. If valid, call strategy.Process()    │
│    ├─ CreditCard: Call payment gateway  │
│    ├─ PayPal: Initiate OAuth flow       │
│    └─ Crypto: Monitor blockchain        │
│ 4. Return transaction ID                │
└─────────────────────────────────────────┘
```

### Switching Strategies at Runtime

```go
Step 1: Initialize
┌─────────────────────────────────┐
│ processor := PaymentProcessor   │
│         (creditCardStrategy)    │
└─────────────────────────────────┘

Step 2: Process Payment 1
┌─────────────────────────────────┐
│ processor.Process($50.00)       │
│ Execution: CreditCardStrategy   │
│ Result: CC-transaction-ID       │
└─────────────────────────────────┘

Step 3: Change Strategy
┌─────────────────────────────────┐
│ processor.SetStrategy           │
│         (paypalStrategy)        │
└─────────────────────────────────┘

Step 4: Process Payment 2
┌─────────────────────────────────┐
│ processor.Process($75.00)       │
│ Execution: PayPalStrategy       │
│ Result: PP-transaction-ID       │
└─────────────────────────────────┘

Same processor, different algorithms!
```

## Go Idioms

Go's design makes the Strategy pattern incredibly natural:

1. **Implicit Interface Satisfaction** - Types automatically satisfy interfaces without explicit declaration
2. **Function Types as Strategies** - For simple cases, use `func` types directly
3. **Interface{} Less Common** - Go's explicit interfaces are cleaner than type assertions
4. **Composition Over Inheritance** - Go encourages composition which aligns perfectly with Strategy
5. **No Setter Boilerplate** - Simple, clean strategy assignment

Example of Go's simplicity:

```go
// Go makes this pattern trivial - no boilerplate needed
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Any type with a Read method is automatically a Reader
// File, Buffer, Network connection, etc. all work without declaring it
```

## Real-World Examples

### 1. Data Compression

Different compression algorithms (GZIP, ZSTD, Brotli) each have different speed/ratio tradeoffs. Users or systems choose based on requirements:

```go
type Compressor interface {
    Compress(data []byte) ([]byte, error)
    Decompress(data []byte) ([]byte, error)
}
```

### 2. Sorting and Filtering

Go's standard library extensively uses Strategy pattern. The `sort.Interface` defines the comparison strategy, allowing sorting any type:

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

### 3. HTTP Request Authentication

Different authentication methods (OAuth, JWT, API Key, Basic Auth) each have different validation logic:

```go
type AuthStrategy interface {
    Authenticate(r *http.Request) (userID string, err error)
}
```

### 4. Database Query Builders

Different SQL dialects (PostgreSQL, MySQL, SQLite) need different query generation strategies:

```go
type SQLDialect interface {
    GenerateQuery(q *Query) string
    EscapeString(s string) string
}
```

### 5. Cache Eviction Policies
Different caching strategies (LRU, LFU, FIFO, TTL) each implement different eviction logic:

```go
type EvictionStrategy interface {
    MarkAccess(key string)
    ShouldEvict() bool
    Evict() string
}
```

## Key Advantages

1. **Flexibility** - Change behavior at runtime without recompilation
2. **Testability** - Mock each strategy independently
3. **Maintainability** - Each strategy is isolated and focused
4. **Extensibility** - Add new strategies without touching existing code
5. **Code Organization** - Related logic is grouped together
6. **Reusability** - Same strategy can be used by different contexts
7. **Clear Architecture** - Explicit about multiple approaches existing

## Key Gotchas

### 1. Over-engineering

Don't use Strategy for simple conditional logic. A simple if/else is fine:

```go
// This is FINE - no need for Strategy
if isFast {
    quickSort(data)
} else {
    bubbleSort(data)
}
```

### 2. Sharing State Between Strategies

If strategies share state, be careful about thread safety:

```go
// BAD - strategies sharing mutable state
strategy.sharedCounter++  // Race condition in concurrent code

// GOOD - pass state as parameters or through dependency injection
strategy.Process(sharedState)
```

### 3. Complex Strategy Selection

If strategy selection logic is complex, it might hide a different pattern:

```go
// If you have this, consider Strategy + Factory pattern
if condition1 && condition2 && condition3 {
    use StrategyA
} else if condition4 || condition5 {
    use StrategyB
}
```

### 4. Performance Overhead

Interface dispatch in Go is very cheap, but be aware:

- Each call goes through the interface indirection
- This is negligible for most use cases but matters in tight loops

### 5. Forgetting the Context

The context is as important as the strategies. Make sure it:

- Properly handles all possible error states from strategies
- Provides clear feedback about which strategy is being used
- Validates strategy compatibility before using

## Common Implementation Patterns

### Simple Strategy Selection

```go
strategy := getStrategy(userPreference)
ctx := NewContext(strategy)
result := ctx.Execute()
```

### Strategy Factory

```go
ctx := NewContext(strategyFactory.Create("payment_method"))
```

### Strategy Composition

```go
// Strategy using another strategy internally
type CachingStrategy struct {
    underlying Strategy
    cache      Cache
}
```

### Function-based Strategies (Go-specific)

For simple cases, use function types:

```go
type Handler func(ctx context.Context) error

type Server struct {
    handlers map[string]Handler
}
```

This is often sufficient and simpler than full interface-based strategies.
