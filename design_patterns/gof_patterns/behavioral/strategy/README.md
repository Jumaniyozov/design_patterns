# Strategy Pattern

## Overview

The Strategy pattern defines a family of algorithms, encapsulates each one, and makes them interchangeable. Strategy lets the algorithm vary independently from clients that use it. In Go, this is naturally implemented using interfaces and function types.

## Problem

You're building an e-commerce system that needs to calculate shipping costs. Different strategies apply:
- Standard shipping: $5.00 flat rate
- Express shipping: $15.00 flat rate
- International: Based on weight and destination
- Free shipping: For orders over $100

Without the Strategy pattern, you'd have conditional logic everywhere:
```go
func calculateShipping(order Order, method string) float64 {
    if method == "standard" {
        return 5.00
    } else if method == "express" {
        return 15.00
    } else if method == "international" {
        // complex calculation...
    } // ... becomes unmanageable
}
```

This violates:
- **Open/Closed Principle**: Adding new strategies requires modifying existing code
- **Single Responsibility**: Shipping calculation mixed with selection logic
- **Testability**: Hard to test individual strategies in isolation

You need a way to define algorithms as objects that can be selected and swapped at runtime.

## Why Use This Pattern?

- **Eliminate Conditionals**: Replace complex if-else/switch statements with polymorphism
- **Runtime Flexibility**: Change algorithm behavior at runtime
- **Testability**: Test each strategy independently
- **Open/Closed**: Add new strategies without modifying existing code
- **Encapsulation**: Each algorithm encapsulated in its own class/type

## When to Use

- **Multiple Algorithms**: Different ways to accomplish the same task
- **Conditional Complexity**: Large if-else or switch statements based on type
- **Runtime Selection**: Algorithm choice depends on runtime conditions
- **Sorting/Search Algorithms**: Different sorting strategies (quicksort, mergesort, etc.)
- **Compression Algorithms**: ZIP, GZIP, BZIP2, etc.
- **Payment Methods**: Credit card, PayPal, crypto, etc.
- **Pricing Strategies**: Regular, discount, wholesale, seasonal
- **Validation Rules**: Different validation strategies for different contexts

## When NOT to Use

- **Single Algorithm**: Only one way to do something
- **Simple Logic**: When if-else is clearer than abstraction
- **Static Behavior**: Algorithm never changes
- **Overengineering**: Adding unnecessary complexity

## Implementation Guidelines

### Basic Strategy Pattern in Go

```go
// 1. Define strategy interface
type CompressionStrategy interface {
    Compress(data []byte) []byte
}

// 2. Implement concrete strategies
type ZipCompression struct{}
func (z *ZipCompression) Compress(data []byte) []byte {
    // ZIP compression logic
    return compressedData
}

type GzipCompression struct{}
func (g *GzipCompression) Compress(data []byte) []byte {
    // GZIP compression logic
    return compressedData
}

// 3. Context uses strategy
type FileCompressor struct {
    strategy CompressionStrategy
}

func (f *FileCompressor) SetStrategy(s CompressionStrategy) {
    f.strategy = s
}

func (f *FileCompressor) Compress(data []byte) []byte {
    return f.strategy.Compress(data)
}
```

### Functional Strategy (Go Idiomatic)

```go
// Strategy as function type
type CompressionFunc func(data []byte) []byte

type FileCompressor struct {
    compressFunc CompressionFunc
}

func (f *FileCompressor) Compress(data []byte) []byte {
    return f.compressFunc(data)
}

// Strategies as functions
func zipCompress(data []byte) []byte { /* ... */ }
func gzipCompress(data []byte) []byte { /* ... */ }

// Usage
compressor := &FileCompressor{compressFunc: zipCompress}
```

## Go Idioms

- **Interface-Based**: Define small, focused strategy interfaces
- **Functional Strategies**: Use function types for simple strategies
- **Constructor Injection**: Pass strategy via `New()` functions
- **Setter Methods**: Allow runtime strategy changes via `SetStrategy()`
- **Nil Checks**: Handle nil strategies gracefully
- **Composition**: Embed strategy in context struct

## Visual Schema

```
Client
  │
  │ configures
  ↓
┌─────────────────────────┐
│  Context                │
│                         │
│  - strategy: Strategy   │
│                         │
│  + SetStrategy(s)       │
│  + ExecuteStrategy()  ──┼──→ strategy.Execute()
│                         │
└─────────────────────────┘
          │
          │ uses
          ↓
┌─────────────────────────┐
│  Strategy Interface     │
│                         │
│  + Execute()            │
└─────────────────────────┘
          ↑
          │ implement
          │
    ┌─────┴─────┬─────────┐
    │           │         │
┌───────┐  ┌───────┐  ┌───────┐
│ ConcreteStrategyA  ConcreteStrategyB  ConcreteStrategyC
│       │  │       │  │       │
│ Execute│  Execute│  Execute│
└───────┘  └───────┘  └───────┘
```

### Before vs After

```
BEFORE (Conditional Logic):

calculateCost(order, method):
  if method == "A":
    // algorithm A
  else if method == "B":
    // algorithm B
  else if method == "C":
    // algorithm C
  // Hard to maintain, test, extend


AFTER (Strategy Pattern):

Context.setStrategy(strategyA)
Context.execute()
      │
      ↓
  strategyA.execute()

Context.setStrategy(strategyB)
Context.execute()
      │
      ↓
  strategyB.execute()

// Easy to maintain, test, extend
```

## Real-World Examples

### 1. Sorting Strategies
```go
type SortStrategy interface {
    Sort(data []int) []int
}

type Sorter struct {
    strategy SortStrategy
}

func (s *Sorter) Sort(data []int) []int {
    return s.strategy.Sort(data)
}

// QuickSort, MergeSort, BubbleSort all implement SortStrategy
```

### 2. Pricing Strategies
```go
type PricingStrategy interface {
    CalculatePrice(basePrice float64) float64
}

type RegularPricing struct{}
func (r *RegularPricing) CalculatePrice(base float64) float64 {
    return base
}

type DiscountPricing struct{ discount float64 }
func (d *DiscountPricing) CalculatePrice(base float64) float64 {
    return base * (1 - d.discount)
}
```

### 3. Validation Strategies
```go
type ValidationStrategy interface {
    Validate(input string) error
}

type EmailValidator struct{}
type PhoneValidator struct{}
type CreditCardValidator struct{}
```

## Key Advantages

✓ **Flexibility**: Swap algorithms at runtime
✓ **Testability**: Test each strategy independently
✓ **Open/Closed**: Add new strategies without modifying context
✓ **Eliminate Conditionals**: Replace if-else with polymorphism
✓ **Encapsulation**: Algorithm details hidden in strategy
✓ **Reusability**: Strategies can be shared across contexts

## Key Gotchas

⚠️ **Strategy Proliferation**: Too many strategies can complicate codebase
⚠️ **Client Awareness**: Clients must know about different strategies
⚠️ **Overhead**: Extra objects/interfaces add slight complexity
⚠️ **Communication**: Context and strategy must share necessary data
⚠️ **Nil Strategies**: Must handle nil strategy gracefully
⚠️ **Thread Safety**: Ensure strategies are thread-safe if shared

## Best Practices

1. **Small Interfaces**: Keep strategy interfaces focused and minimal
2. **Immutable Strategies**: Make strategies stateless when possible
3. **Factory Pattern**: Use factories to create and select strategies
4. **Default Strategy**: Provide sensible default if no strategy set
5. **Configuration**: Load strategy selection from configuration
6. **Documentation**: Clearly document each strategy's behavior
7. **Naming**: Use descriptive names (StandardShipping, ExpressShipping)
8. **Thread Safety**: Ensure strategies can be used concurrently
9. **Error Handling**: Define consistent error handling across strategies
10. **Combine Patterns**: Use with Factory, Template Method, or Dependency Injection

## Pattern Variations

### Functional Strategy
Use function types instead of interfaces for simple cases:
```go
type ProcessFunc func(data []byte) []byte

type Processor struct {
    process ProcessFunc
}
```

### Strategy with State
Strategies can maintain internal state:
```go
type CacheStrategy interface {
    Store(key, value string)
    Get(key string) (string, bool)
}
```

### Strategy Registry
Register strategies by name:
```go
var strategies = map[string]Strategy{
    "fast": &FastStrategy{},
    "slow": &SlowStrategy{},
}

func GetStrategy(name string) Strategy {
    return strategies[name]
}
```
