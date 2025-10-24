# Design Patterns in Go

A comprehensive, tier-based curriculum for mastering software design patterns using Go. This repository is designed for senior software engineers who want to deepen their architectural knowledge and learn how classic design patterns translate to Go's unique idioms and concurrency model.

## What Are Design Patterns?

Design patterns are reusable, proven solutions to common software design problems. They represent best practices evolved over decades of software engineering, codifying expert knowledge into recognizable structures that solve recurring architectural challenges.

**Design patterns are NOT:**
- Copy-paste code templates
- Silver bullets that solve all problems
- Mandatory rules you must always follow
- Language-specific implementations

**Design patterns ARE:**
- Shared vocabulary for discussing software architecture
- Proven approaches to common structural and behavioral problems
- Guidelines that can be adapted to your specific context
- Tools that help you write maintainable, flexible, and testable code

Think of design patterns as architectural blueprints. An architect doesn't reinvent how to design a load-bearing wall for every building—they apply proven structural principles. Similarly, design patterns provide architectural principles for software systems.

## Why Use Design Patterns?

### 1. **Communication and Shared Vocabulary**
When you say "let's use a Factory pattern here," other engineers immediately understand your intent. Patterns create a common language that transcends code, enabling clearer communication in design discussions and code reviews.

### 2. **Proven Solutions to Known Problems**
Rather than reinventing solutions, patterns let you leverage decades of collective experience. You avoid common pitfalls because others have already discovered and documented them.

### 3. **Code Maintainability**
Patterns make code more recognizable. When a future engineer (including future you) sees a familiar pattern, they can quickly understand the design intent and modify the code confidently.

### 4. **Flexibility and Extensibility**
Good patterns make it easier to change behavior without rewriting large portions of code. They promote loose coupling and high cohesion, key principles of maintainable software.

### 5. **Testing and Mocking**
Many patterns (especially Strategy, Factory, and Decorator) naturally create interfaces that make unit testing and dependency injection straightforward.

## When to Use Design Patterns?

### Use Patterns When:

✅ **You have a specific, identifiable problem** - Don't apply patterns just because they're "best practices." Apply them when they solve a real problem you're facing.

✅ **The codebase is growing in complexity** - Patterns help manage complexity as systems scale. A 100-line script doesn't need design patterns; a 100,000-line service probably does.

✅ **You need flexibility for future changes** - When you know certain parts of the system will need to vary or extend (new payment methods, different notification channels, etc.), patterns provide clean extension points.

✅ **Multiple implementations exist or will exist** - Strategy pattern makes sense when you have multiple sorting algorithms; Adapter makes sense when integrating multiple third-party services.

✅ **You're working in a team** - Patterns provide shared understanding and make code reviews more productive.

### Avoid Patterns When:

❌ **The problem is simple** - Don't use Builder pattern for a struct with two fields. Don't use Strategy for a single implementation. Over-engineering creates unnecessary complexity.

❌ **You're pattern hunting** - Don't force patterns into code. If you're asking "where can I apply this pattern?" instead of "what pattern solves this problem?", you're approaching it backwards.

❌ **Performance is critical and the pattern adds overhead** - Some patterns (like Proxy, Decorator) add indirection. In hot paths with microsecond requirements, simpler approaches may be necessary.

❌ **The pattern doesn't fit the language idioms** - Go is not Java or C++. Some classic patterns don't translate well. For example, Go doesn't have traditional inheritance, so classical Template Method needs adaptation.

## Categories of Design Patterns

Design patterns fall into four main categories, each addressing different types of architectural challenges:

### 1. Creational Patterns (Object Creation)

These patterns abstract the instantiation process, making systems independent of how objects are created, composed, and represented.

**Common Problems They Solve:**
- Complex object construction with many parameters
- Controlling how many instances exist (singletons)
- Creating objects without specifying exact classes
- Building complex objects step-by-step

**Patterns in This Repository:**
- **Factory Pattern** - Centralized object creation logic
- **Builder Pattern** - Step-by-step construction of complex objects
- **Singleton Pattern** - Ensure only one instance exists
- **Options Pattern** - Flexible, backward-compatible constructors (Go-specific)

**When to Use:** When object creation is complex, requires configuration, or needs centralized control.

### 2. Structural Patterns (Object Composition)

These patterns deal with how classes and objects are composed to form larger structures, ensuring flexibility and efficiency.

**Common Problems They Solve:**
- Adapting incompatible interfaces
- Adding functionality without modifying existing code
- Simplifying complex subsystems
- Managing object hierarchies

**Patterns in This Repository:**
- **Adapter Pattern** - Make incompatible interfaces work together
- **Decorator Pattern** - Add responsibilities dynamically
- **Proxy Pattern** - Control access to objects
- **Composite Pattern** - Treat individual objects and compositions uniformly
- **Facade Pattern** - Provide simplified interface to complex subsystems

**When to Use:** When you need to compose objects in flexible ways, integrate third-party code, or manage complexity through layering.

### 3. Behavioral Patterns (Object Interaction)

These patterns focus on communication between objects, defining how objects interact and distribute responsibility.

**Common Problems They Solve:**
- Decoupling senders and receivers
- Managing object state transitions
- Defining family of algorithms
- Processing requests through a chain

**Patterns in This Repository:**
- **Strategy Pattern** - Encapsulate interchangeable algorithms
- **Observer Pattern** - Notify dependents of state changes
- **Command Pattern** - Encapsulate requests as objects
- **Chain of Responsibility** - Pass requests along a handler chain
- **Template Method Pattern** - Define algorithm skeleton, vary steps
- **State Pattern** - Alter behavior when internal state changes

**When to Use:** When you need flexible communication patterns, event-driven architectures, or behavior that changes based on state.

### 4. Concurrency Patterns (Go-Specific Excellence)

While not part of the original Gang of Four patterns, concurrency patterns are essential in Go. They leverage goroutines and channels to build scalable, concurrent systems.

**Common Problems They Solve:**
- Processing data streams concurrently
- Distributing work across multiple workers
- Managing resource pools
- Preventing cascading failures in distributed systems

**Patterns in This Repository:**
- **Pipeline Pattern** - Process data through stages
- **Fan-Out/Fan-In Pattern** - Distribute work, aggregate results
- **Worker Pool Pattern** - Limit concurrent workers
- **Circuit Breaker Pattern** - Fail fast and recover gracefully

**When to Use:** When building concurrent systems, processing large datasets, calling external services, or optimizing for throughput.

## How to Use This Repository

### Learning Path: Tier-Based Progression

This repository organizes 19 essential patterns into four progressive tiers. **Start with Tier 1** and build mastery before advancing.

```
┌─────────────────────────────────────────────────────────┐
│ Tier 1: Essential Go Patterns                           │
│ Strategy, Factory, Builder, Decorator, Singleton        │
│ Foundation - Use these daily in production              │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│ Tier 2: Structural Patterns                             │
│ Adapter, Proxy, Composite, Facade                       │
│ Organize larger codebases, manage complexity            │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│ Tier 3: Behavioral Patterns                             │
│ Observer, Command, Chain of Responsibility,             │
│ Template Method, State                                  │
│ Advanced object communication and interaction           │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│ Tier 4: Concurrency & Advanced Patterns                 │
│ Pipeline, Fan-Out/Fan-In, Worker Pool,                  │
│ Circuit Breaker, Options                                │
│ Go-specific excellence for high-performance systems     │
└─────────────────────────────────────────────────────────┘
```

### Each Pattern Includes:

Each pattern directory contains:

- **README.md** - Comprehensive explanation with:
  - Problem statement and real-world context
  - Why, when, and when NOT to use the pattern
  - Visual schemas and diagrams
  - Go-specific idioms and best practices
  - Real-world examples and use cases

- **pattern.go** - Core implementation showing:
  - Interface definitions
  - Concrete implementations (2-3 examples)
  - Context/processor that uses the pattern

- **example.go** - Runnable examples demonstrating:
  - Basic usage
  - Real-world scenarios
  - How to switch between implementations

- **main_test.go** - Comprehensive tests including:
  - Integration tests
  - Strategy switching tests
  - Validation and edge cases
  - Benchmarks where applicable

### Running Examples

Navigate to any pattern directory and run:

```bash
# Run all examples for a pattern
go run tier1/strategy/*.go

# Run tests
go test ./tier1/strategy

# Run specific test
go test ./tier1/strategy -run TestStrategyPattern

# View coverage
go test -cover ./tier1/strategy
```

### Build and Test Everything

```bash
# Build all patterns
go build ./...

# Test all patterns
go test ./...

# View coverage across all patterns
go test -cover ./...

# Format code
go fmt ./...

# Lint code
go vet ./...
```

## Recommended Learning Approach

### 1. **Start Small, Build Projects**
Don't just read the patterns—apply them. After completing Tier 1, build 2-3 small projects using those patterns:
- A simple web service with Strategy for different authentication methods
- A configuration loader using Factory and Builder
- A middleware pipeline using Decorator

### 2. **Focus on Problems, Not Patterns**
Ask yourself: "What problem am I solving?" rather than "Which pattern should I use?" The pattern should emerge naturally from the problem.

### 3. **Understand Trade-offs**
Every pattern has costs (complexity, indirection, learning curve). The README for each pattern includes "When NOT to Use" sections—read them carefully.

### 4. **Connect Patterns Together**
As you progress, you'll see patterns compose naturally:
- Factory + Adapter for creating adapters
- Observer + Command for event-driven systems
- Decorator + Strategy for flexible behavior chains

### 5. **Apply Go Idioms**
Go is not Java or C++. These patterns are adapted to Go's:
- Interface-based composition (not inheritance)
- Explicit error handling
- Goroutines and channels for concurrency
- Simplicity and clarity over complexity

## General Best Practices

### ✨ Do's

- **Keep It Simple** - Use the simplest solution that solves your problem
- **Prefer Composition Over Inheritance** - Go's interfaces make composition natural
- **Design for Testability** - Patterns should make testing easier, not harder
- **Document Intent** - Explain WHY you chose a pattern in comments
- **Iterate** - Start simple, add patterns as complexity demands
- **Follow Go Conventions** - Accept interfaces, return concrete types

### ⚠️ Don'ts

- **Don't Over-Engineer** - Premature abstraction is costly
- **Don't Force Patterns** - If it doesn't fit naturally, don't use it
- **Don't Ignore Context** - What works for a microservice might not work for a CLI tool
- **Don't Forget YAGNI** - "You Aren't Gonna Need It" - add patterns when you need them
- **Don't Sacrifice Readability** - If a pattern makes code harder to understand, reconsider

## Pattern Quick Reference

| Pattern | Category | Primary Use Case | Key Benefit |
|---------|----------|------------------|-------------|
| Strategy | Behavioral | Interchangeable algorithms | Flexibility in behavior |
| Factory | Creational | Complex object creation | Centralized construction |
| Builder | Creational | Objects with many options | Readable initialization |
| Decorator | Structural | Add behavior dynamically | Composable functionality |
| Singleton | Creational | Single shared instance | Resource management |
| Adapter | Structural | Interface compatibility | Integration flexibility |
| Proxy | Structural | Control object access | Lazy loading, caching |
| Composite | Structural | Tree structures | Uniform treatment |
| Facade | Structural | Simplify subsystems | Reduced complexity |
| Observer | Behavioral | Event notification | Decoupled communication |
| Command | Behavioral | Encapsulate requests | Undo/redo, queuing |
| Chain of Responsibility | Behavioral | Request processing chain | Flexible pipelines |
| Template Method | Behavioral | Algorithm skeleton | Reuse with variation |
| State | Behavioral | Behavior varies by state | Clean state management |
| Pipeline | Concurrency | Data stream processing | Concurrent stages |
| Fan-Out/Fan-In | Concurrency | Parallel processing | Maximize throughput |
| Worker Pool | Concurrency | Limited concurrency | Resource control |
| Circuit Breaker | Concurrency | Failure isolation | System resilience |
| Options | Creational | Flexible constructors | Clean APIs |

## Project Structure

```
design_patterns/
├── README.md              # This file
├── CLAUDE.md             # Development guidelines
├── go.mod                # Go module definition
├── tier1/                # Essential Go Patterns
│   ├── strategy/
│   ├── factory/
│   ├── builder/
│   ├── decorator/
│   └── singleton/
├── tier2/                # Structural Patterns
│   ├── adapter/
│   ├── proxy/
│   ├── composite/
│   └── facade/
├── tier3/                # Behavioral Patterns
│   ├── observer/
│   ├── command/
│   ├── chainofresponsibility/
│   ├── templatemethod/
│   └── state/
└── tier4/                # Concurrency and Advanced
    ├── pipeline/
    ├── fanout/
    ├── workerpool/
    ├── circuitbreaker/
    └── options/
```

## Requirements

- **Go 1.25** or later
- Basic understanding of Go syntax and interfaces
- Familiarity with software engineering fundamentals

## Getting Started

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd design_patterns
   ```

2. **Start with Tier 1, Pattern 1: Strategy**
   ```bash
   cd tier1/strategy
   cat README.md  # Read the comprehensive explanation
   go run *.go    # Run the examples
   go test        # Run tests
   ```

3. **Work through each pattern systematically**
   - Read the README thoroughly
   - Study the implementation
   - Run and modify examples
   - Complete exercises (if provided)
   - Build small projects applying the pattern

4. **Progress to next tier only after mastery**
   - Can you explain when to use each pattern?
   - Can you implement it from scratch?
   - Can you identify it in real-world code?

## Contributing

This is a learning repository. If you find errors, have improvements, or want to add examples:
- Ensure additions follow Go conventions
- Include comprehensive documentation
- Provide runnable examples
- Add tests for new code

## Resources and Further Reading

- **Design Patterns: Elements of Reusable Object-Oriented Software** (Gang of Four) - The canonical reference
- **Head First Design Patterns** - Accessible, visual introduction
- **Go Concurrency Patterns** (Rob Pike) - Essential for Tier 4
- **Effective Go** - Go-specific idioms and best practices
- **The Go Blog: Advanced Go Concurrency Patterns** - Deep dive into Go patterns

## Philosophy

> "Design patterns should not be goals in themselves. They are tools—use them when they solve problems, ignore them when they don't. The best code is often the simplest code that solves the problem at hand."

Good architecture emerges from understanding principles (SOLID, DRY, YAGNI) and applying patterns judiciously when they provide clear value. This repository aims to build that judgment through comprehensive explanation and practical application.

---

**Happy Learning!** Start with `tier1/strategy/` and build your architectural expertise one pattern at a time.
