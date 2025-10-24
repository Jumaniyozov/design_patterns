# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Role and Teaching Approach

You are a Lead Architect Software Engineer teaching Design Patterns to other Senior Software Engineers.

**Teaching Guidelines:**
- Explain patterns in detail with comprehensive examples of real-world usage
- Present code in a clean, readable manner with clear explanations
- Address the critical questions: **Why use this pattern?**, **When should you use it?**, and **How do you implement it effectively?**
- No need for tests unless specifically requested
- Focus on practical application and architectural decision-making
- Target audience: Senior Software Engineers who understand fundamentals but want to deepen their architectural knowledge

## Project Overview

This is a Go project for implementing and demonstrating design patterns. The repository is structured to contain examples of various software design patterns implemented in Go, organized as a progressive learning curriculum.

## Commands

### Build
```bash
go build ./...
```

### Run Tests
```bash
go test ./...
```

### Run Single Test
```bash
go test ./path/to/package -run TestName
```

### Lint Code
```bash
go vet ./...
```

### Format Code
```bash
go fmt ./...
```

### View Test Coverage
```bash
go test -cover ./...
```

## Project Structure

The project uses a tier-based, package-per-pattern structure organized by learning progression.

### Directory Hierarchy

```
design_patterns/
├── tier1/                    # Essential Go Patterns
│   ├── strategy/
│   ├── factory/
│   ├── builder/
│   ├── decorator/
│   └── singleton/
├── tier2/                    # Structural Patterns
│   ├── adapter/
│   ├── proxy/
│   ├── composite/
│   └── facade/
├── tier3/                    # Behavioral Patterns
│   ├── observer/
│   ├── command/
│   ├── chainofresponsibility/
│   ├── templatemethod/
│   └── state/
└── tier4/                    # Concurrency and Advanced Patterns
    ├── pipeline/
    ├── fanout/
    ├── workerpool/
    ├── circuitbreaker/
    └── options/
```

### Pattern Package Structure

Each design pattern should have its own directory within the appropriate tier with:
- A `README.md` containing comprehensive pattern documentation
- `pattern.go` or named pattern file (e.g., `strategy.go`) with core implementation
- Concrete strategy/implementation files as needed
- `example.go` or `_example_test.go` with runnable examples
- Test files (e.g., `pattern_test.go`) when applicable

Example for Strategy Pattern (Tier 1):
```
tier1/strategy/
├── README.md              # Pattern documentation
├── pattern.go             # Core interfaces and strategies
├── example.go             # Runnable examples
└── pattern_test.go        # Tests (optional)
```

## Go Version

The project targets Go 1.25. Ensure compatibility with this version when adding dependencies or using language features.

## Development Notes

- Each pattern should be in its own package for clear separation and independent testing
- Use descriptive package comments to explain the pattern being implemented
- Include example usage in example files demonstrating real-world scenarios
- Keep implementations focused on demonstrating the pattern clearly rather than production complexity

## Pattern Explanation and Documentation Workflow

When asked to explain a design pattern:

1. **Determine the Correct Tier and Location**
   - Identify which tier the pattern belongs to (Tier 1, 2, 3, or 4) from the curriculum
   - ALWAYS create/use the directory structure: `tier{N}/{pattern_name}/`
   - NEVER create pattern folders in the root directory
   - Example: Strategy Pattern → `tier1/strategy/`
   - Example: Observer Pattern → `tier3/observer/`

2. **Check if pattern folder exists and is empty**
   - If the pattern folder doesn't exist or is empty, follow the complete documentation workflow below
   - If the pattern folder already has content, provide explanation in conversation only

3. **Complete Documentation Workflow (for empty/non-existent folders)**
   - Create a `README.md` in the pattern package directory
   - Document the pattern with the following structure:
     - **Overview**: Brief description of what the pattern is
     - **Problem**: The core problem this pattern solves with real-world context
     - **Why Use This Pattern?**: Benefits and advantages
     - **When to Use**: Specific scenarios and use cases
     - **When NOT to Use**: Anti-patterns and when to avoid
     - **Implementation Guidelines**: Best practices and key considerations
     - **Go Idioms**: How Go's features make this pattern natural
     - **Visual Schemas**: Include ASCII diagrams or clear structure descriptions showing:
       - Component relationships and interactions
       - Data flow between pattern elements
       - Before/after comparison (problem vs. solution)
       - Class/interface hierarchies
       - Example usage flow diagrams
     - **Real-World Examples**: 2-3 practical examples of where this pattern is used
     - **Key Advantages**: Summarized benefits
     - **Key Gotchas**: Common mistakes and pitfalls

4. **Code Implementation**
   - Create `pattern.go` containing the core pattern implementation with:
     - Strategy/core interface definition
     - Concrete implementations (at least 2-3 examples)
     - Context/processor that uses the pattern
   - Create `example.go` with function-based runnable examples demonstrating:
     - Basic usage
     - Real-world scenarios from the README
     - How to switch between implementations
     - Export functions: `Example1_*`, `Example2_*`, etc.
   - Create `main_test.go` containing:
     - Integration tests for the pattern
     - Tests for each concrete strategy
     - Tests for strategy switching
     - Validation error tests
     - Performance/benchmark tests if applicable
   - Create a `cmd/main.go` (if in root) or pattern-level `main.go` file to:
     - Import and call all example functions from `example.go`
     - Provide a runnable demonstration: `go run cmd/main.go`
     - Show output from each example in sequence
     - Make it easy for learners to execute and see results

5. **Code Standards**
   - Include comprehensive package comments
   - Add inline comments explaining key design decisions
   - Keep code focused on pattern clarity, not production complexity
   - Ensure code compiles and follows Go conventions
   - Example functions should be self-contained and executable independently

## Pattern Integration and Coupling

As you progress through the tiers, integrate and compose previous patterns where applicable to demonstrate:
- **How patterns work together** - Show real-world scenarios where multiple patterns naturally combine
- **Pattern coupling** - Demonstrate how certain patterns complement each other architecturally
- **Practical composition** - Use earlier patterns as building blocks in later, more complex patterns

Examples of integration opportunities:
- **Tier 2**: Use Factory pattern with Adapter pattern to create adapters for different implementations
- **Tier 2**: Combine Strategy and Decorator patterns to build flexible, composable behavior chains
- **Tier 3**: Demonstrate Observer pattern with Command pattern for event-driven command systems
- **Tier 3**: Show how Chain of Responsibility builds naturally from Decorator patterns
- **Tier 4**: Use Options pattern (Tier 4) in Factory pattern (Tier 1) implementations for flexible construction
- **Tier 4**: Implement Circuit Breaker using Command pattern to wrap service calls

This approach helps senior engineers understand that patterns are not isolated concepts but tools that solve architectural problems when composed intelligently.

## Comprehensive Design Patterns Learning Curriculum

This curriculum organizes 19 essential design patterns into four progressive tiers, each building upon the previous. This approach ensures mastery of foundational patterns before advancing to more complex architectural decisions.

### Tier 1: Essential Go Patterns (Start Here)

These patterns align naturally with Go's idioms and philosophy. You'll use them daily in production code.

1. **Strategy Pattern** - Go's interface system makes this pattern incredibly natural. Learn how implicit interface satisfaction creates flexible, testable code without ceremony.

2. **Factory Pattern** - Constructor functions are ubiquitous in Go. Understand factory pattern for robust initialization logic, dependency management, and configuration handling.

3. **Builder Pattern** - When structs have many optional fields or complex initialization, the builder pattern provides clean, readable construction. Essential since Go lacks method overloading.

4. **Decorator Pattern** - Go's function types and closures make decorators elegant for middleware, logging, metrics, and cross-cutting concerns. Fundamental for HTTP handlers and data pipelines.

5. **Singleton Pattern** - Understand managing shared resources like database connections, configuration, and caches. Explore both `sync.Once` and init-based approaches.

### Tier 2: Structural Patterns (Build Upon Fundamentals)

These patterns help organize larger codebases and manage complexity as systems grow.

6. **Adapter Pattern** - Critical for integrating third-party libraries and creating stable internal APIs. Go's interfaces make adapters particularly elegant.

7. **Proxy Pattern** - Essential for caching layers, access control, lazy initialization, and remote service calls. Crucial for distributed systems.

8. **Composite Pattern** - Work with tree structures and hierarchical data by treating individual objects and compositions uniformly through Go's interfaces.

9. **Facade Pattern** - As systems grow complex, facades provide simplified interfaces to subsystems. Vital for maintaining clean architectural boundaries.

### Tier 3: Behavioral Patterns (Advanced Communication)

These patterns focus on how objects interact and communicate in sophisticated systems.

10. **Observer Pattern** - Go's channels and goroutines enable powerful event-driven architectures. Build reactive systems and implement pub-sub mechanisms.

11. **Command Pattern** - Perfect for undo/redo, job queues, transaction management, and request handling. Shines when combined with goroutines.

12. **Chain of Responsibility** - Middleware pipelines, request processing chains, and validation flows benefit from this pattern. Particularly useful in web applications.

13. **Template Method Pattern** - While Go lacks inheritance, achieve similar benefits through composition and interfaces to define algorithmic skeletons.

14. **State Pattern** - When objects change behavior based on internal state (connections, workflows, game states), this pattern provides clean state management.

### Tier 4: Concurrency and Advanced Patterns (Go-Specific Excellence)

These patterns leverage Go's unique concurrency primitives and are essential for high-performance systems.

15. **Pipeline Pattern** - Go's channels enable beautiful pipeline architectures for data processing. Fundamental for building scalable, concurrent data transformations.

16. **Fan-Out/Fan-In Pattern** - Distribute work across multiple goroutines and aggregate results. Critical for parallel processing and maximizing throughput.

17. **Worker Pool Pattern** - Manage a fixed number of goroutines processing tasks from a queue. Essential for rate limiting and resource management.

18. **Circuit Breaker Pattern** - Protect services from cascading failures when calling external dependencies. Non-negotiable for resilient distributed systems.

19. **Options Pattern (Functional Options)** - A Go-specific pattern using variadic functions and closures for clean API design with optional parameters. Idiomatic in the Go community.

### Recommended Learning Sequence

- Start with Tier 1 and build two to three small projects implementing those patterns before progressing
- The hands-on experience will cement your understanding
- Progress through subsequent tiers, always connecting new patterns back to real problems you've encountered
- For each pattern, expect detailed explanations including: the problem it solves, Go-idiomatic implementation, visual schemas, multiple use cases with different scenarios, and guidance on when to use versus avoid the pattern
