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

The project follows the classic Gang of Four (GOF) categorization, organizing all 23 design patterns by their fundamental purpose: Creational, Structural, and Behavioral.

### Directory Hierarchy

```
design_patterns/
└── gof_patterns/                    # All Gang of Four Patterns
    ├── creational/                  # Object Creation Patterns
    │   ├── singleton/               # Ensure one instance
    │   ├── factory_method/          # Create objects via factory
    │   ├── abstract_factory/        # Families of related objects
    │   ├── builder/                 # Construct complex objects step-by-step
    │   └── prototype/               # Clone existing objects
    │
    ├── structural/                  # Object Composition & Relationships
    │   ├── adapter/                 # Convert interfaces
    │   ├── bridge/                  # Separate abstraction from implementation
    │   ├── composite/               # Tree structures of objects
    │   ├── decorator/               # Add responsibilities dynamically
    │   ├── facade/                  # Simplified interface to subsystem
    │   ├── flyweight/               # Share objects to save memory
    │   └── proxy/                   # Surrogate or placeholder
    │
    └── behavioral/                  # Object Interaction & Responsibility
        ├── chain_of_responsibility/ # Pass requests along chain
        ├── command/                 # Encapsulate requests as objects
        ├── interpreter/             # Language grammar interpreter
        ├── iterator/                # Sequential access to elements
        ├── mediator/                # Centralize complex communications
        ├── memento/                 # Capture and restore object state
        ├── observer/                # Notify dependents of changes
        ├── state/                   # Change behavior when state changes
        ├── strategy/                # Encapsulate interchangeable algorithms
        ├── template_method/         # Define algorithm skeleton
        └── visitor/                 # Operations on object structure elements
```

### Pattern Package Structure

Each design pattern has its own directory with a standard structure:

**Required Files:**
- `README.md` - Comprehensive pattern documentation
- `{pattern_name}.go` - Core implementation (e.g., `singleton.go`, `adapter.go`, `strategy.go`)
- `example.go` - Runnable examples demonstrating usage

**Optional Files:**
- `{pattern_name}_test.go` - Unit tests (when explicitly requested)

Example for Singleton Pattern:
```
creational/singleton/
├── README.md              # Pattern documentation
├── singleton.go           # Core implementation
└── example.go             # Runnable examples
```

Example for Strategy Pattern:
```
behavioral/strategy/
├── README.md              # Pattern documentation
├── strategy.go            # Core interfaces and implementations
└── example.go             # Runnable examples
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

1. **Determine the Correct Category and Location**
   - Identify which GOF category the pattern belongs to: **Creational**, **Structural**, or **Behavioral**
   - ALWAYS use the directory structure: `design_patterns/gof_patterns/{category}/{pattern_name}/`
   - NEVER create pattern folders in the root directory or outside the GOF structure
   - Examples:
     - Singleton → `design_patterns/gof_patterns/creational/singleton/`
     - Adapter → `design_patterns/gof_patterns/structural/adapter/`
     - Strategy → `design_patterns/gof_patterns/behavioral/strategy/`
     - Observer → `design_patterns/gof_patterns/behavioral/observer/`

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
   - Create `{pattern_name}.go` containing the core pattern implementation with:
     - Core interface definition
     - Concrete implementations (at least 2-3 examples)
     - Context/processor that uses the pattern (if applicable)
   - Create `example.go` with function-based runnable examples demonstrating:
     - Basic usage
     - Real-world scenarios from the README
     - How to switch between implementations
     - Export functions: `Example1_*`, `Example2_*`, etc.
   - Tests are OPTIONAL unless explicitly requested:
     - Create `{pattern_name}_test.go` only when asked
     - Include unit tests for the pattern
     - Tests for each concrete implementation
     - Tests for edge cases and error handling
     - Performance/benchmark tests if applicable

5. **Code Standards**
   - Include comprehensive package comments
   - Add inline comments explaining key design decisions
   - Keep code focused on pattern clarity, not production complexity
   - Ensure code compiles and follows Go conventions
   - Example functions should be self-contained and executable independently

## Pattern Integration and Coupling

Patterns are rarely used in isolation. Demonstrate how patterns naturally compose and complement each other:

### Creational + Structural Integration
- **Factory + Adapter**: Use Factory Method to create appropriate adapters for different third-party services
- **Builder + Composite**: Build complex composite structures step-by-step
- **Abstract Factory + Bridge**: Create families of related objects that work with different implementations
- **Singleton + Facade**: Singleton ensures one facade instance managing subsystem access

### Creational + Behavioral Integration
- **Factory + Strategy**: Factory creates and selects appropriate strategy implementations
- **Builder + Command**: Build complex commands with multiple parameters
- **Prototype + Memento**: Clone objects to save/restore state efficiently

### Structural + Behavioral Integration
- **Decorator + Strategy**: Wrap strategies with decorators to add cross-cutting concerns (logging, caching)
- **Composite + Visitor**: Apply operations across composite structures
- **Proxy + Command**: Proxy can queue commands for remote execution
- **Adapter + Template Method**: Adapt third-party code into template method steps

### Multi-Pattern Compositions
- **Observer + Command + Memento**: Event-driven systems with undo/redo capabilities
- **Chain of Responsibility + Command**: Process command objects through handler chains
- **Strategy + Factory + Singleton**: Factory selects strategies, Singleton manages factory instance
- **Facade + Adapter + Proxy**: Simplified interface (Facade) to adapted (Adapter) remote services (Proxy)

This approach helps senior engineers understand that patterns are not isolated concepts but tools that solve architectural problems when composed intelligently.

## Gang of Four (GOF) Design Patterns Curriculum

This curriculum covers all 23 classic GOF design patterns, organized by their fundamental purpose. Each pattern is implemented in Go with idiomatic examples and comprehensive documentation.

### Creational Patterns (5 patterns)

Creational patterns abstract the object instantiation process, making systems independent of how objects are created, composed, and represented.

1. **Singleton** - Ensures a class has only one instance and provides a global access point. Essential for managing shared resources (database pools, configuration, loggers). In Go, use `sync.Once` for thread-safe lazy initialization or package-level variables for eager initialization.

2. **Factory Method** - Defines an interface for creating objects but lets functions decide which concrete type to instantiate. Ubiquitous in Go through constructor functions (`New...()`). Critical for runtime type selection, plugin systems, and decoupling client code from concrete types.

3. **Abstract Factory** - Provides an interface for creating families of related or dependent objects without specifying concrete classes. Use when you need to create coordinated sets of objects (UI themes, cross-platform libraries, database drivers for different vendors).

4. **Builder** - Separates complex object construction from representation, allowing the same construction process to create different representations. Essential in Go since there's no method overloading. Perfect for objects with many optional parameters or complex initialization sequences.

5. **Prototype** - Creates new objects by cloning existing instances rather than creating from scratch. Useful when object creation is expensive or complex. In Go, implement through Clone() methods or copy constructors.

### Structural Patterns (7 patterns)

Structural patterns explain how to compose classes and objects to form larger structures while keeping them flexible and efficient.

6. **Adapter** - Converts one interface into another that clients expect, making incompatible interfaces work together. Critical for third-party library integration, legacy code modernization, and maintaining stable internal APIs. Go's interfaces make this pattern particularly elegant.

7. **Bridge** - Separates abstraction from implementation so they can vary independently. Use when you want to avoid permanent binding between abstraction and implementation, or when both need to be extended through subclassing.

8. **Composite** - Composes objects into tree structures to represent part-whole hierarchies. Lets clients treat individual objects and compositions uniformly. Perfect for file systems, UI component trees, organizational hierarchies, and any recursive data structures.

9. **Decorator** - Attaches additional responsibilities to objects dynamically without affecting other objects. Go's function types and closures make decorators elegant. Fundamental for middleware, logging, metrics, caching, and HTTP handler chains.

10. **Facade** - Provides a simplified, unified interface to a complex subsystem. Essential as systems grow to maintain clean architectural boundaries, reduce coupling, and hide complexity from clients.

11. **Flyweight** - Uses sharing to support large numbers of fine-grained objects efficiently by sharing common state. Useful for rendering systems, text editors (sharing character objects), game development (reusing terrain textures).

12. **Proxy** - Provides a surrogate or placeholder for another object to control access to it. Critical for virtual proxies (lazy loading), protection proxies (access control), remote proxies (RPC), and smart references (caching, reference counting).

### Behavioral Patterns (11 patterns)

Behavioral patterns characterize how classes and objects interact and distribute responsibility, focusing on communication between objects.

13. **Chain of Responsibility** - Passes requests along a chain of handlers where each handler decides either to process the request or pass it to the next handler. Excellent for middleware pipelines, validation chains, event bubbling, and request processing in web frameworks.

14. **Command** - Encapsulates a request as an object, allowing you to parameterize clients with different requests, queue requests, and support undoable operations. Perfect for job queues, transaction systems, macro recording, and undo/redo functionality.

15. **Interpreter** - Defines a grammatical representation for a language and provides an interpreter to process sentences in that language. Use for domain-specific languages (DSLs), query languages, configuration parsers, and expression evaluators.

16. **Iterator** - Provides a way to access elements of a collection sequentially without exposing its underlying representation. In Go, this is often implemented through channels, range-based iteration, or explicit iterator types.

17. **Mediator** - Defines an object that encapsulates how a set of objects interact, promoting loose coupling by keeping objects from referring to each other explicitly. Useful for complex UI interactions, chat room servers, air traffic control systems.

18. **Memento** - Captures and externalizes an object's internal state so it can be restored later, without violating encapsulation. Essential for undo/redo, snapshots, transaction rollback, and game save states.

19. **Observer** - Defines a one-to-many dependency between objects so when one object changes state, all dependents are notified. Go's channels and goroutines make this particularly powerful. Critical for event systems, reactive programming, pub-sub architectures, and model-view separation.

20. **State** - Allows an object to alter its behavior when its internal state changes, appearing to change its class. Perfect for state machines, connection management (connected/disconnected), order processing workflows, and game character states.

21. **Strategy** - Defines a family of algorithms, encapsulates each one, and makes them interchangeable. Go's interface system makes this pattern incredibly natural. Use for sorting algorithms, compression algorithms, payment processing, validation rules, and pricing strategies.

22. **Template Method** - Defines the skeleton of an algorithm in a base operation, deferring some steps to subclasses. While Go lacks inheritance, achieve similar benefits through composition and interfaces. Useful for frameworks, data processing pipelines, and algorithmic skeletons.

23. **Visitor** - Represents an operation to be performed on elements of an object structure, letting you define new operations without changing the classes of elements. Useful for AST traversal, reporting operations, serialization, and operations on composite structures.

### Recommended Learning Sequence

**Phase 1 - Foundations (Start Here):**
- Singleton, Factory Method, Strategy, Adapter
- These are the most commonly used patterns and form the foundation

**Phase 2 - Essential Extensions:**
- Builder, Decorator, Observer, Command
- Build on Phase 1 knowledge with more sophisticated patterns

**Phase 3 - Structural Sophistication:**
- Composite, Facade, Proxy, Bridge
- Master structural organization and relationships

**Phase 4 - Behavioral Mastery:**
- State, Chain of Responsibility, Template Method, Mediator
- Advanced communication and responsibility distribution

**Phase 5 - Complete the Collection:**
- Abstract Factory, Prototype, Flyweight, Iterator, Memento, Interpreter, Visitor
- Specialized patterns for specific use cases

### Learning Approach

- **Understand the Problem First**: Each pattern solves a specific problem. Understand the problem before learning the solution.
- **Go-Idiomatic Implementation**: Focus on how Go's features (interfaces, embedding, closures, channels) make patterns more elegant.
- **Real-World Context**: Every pattern includes practical examples from actual software systems.
- **When NOT to Use**: Understanding anti-patterns is as important as understanding patterns.
- **Pattern Composition**: Learn how patterns work together to solve complex architectural challenges.
