# Abstract Factory Pattern

## Overview
The Abstract Factory pattern provides an interface for creating families of related or dependent objects without specifying their concrete classes. It's a factory of factories—a super-factory that creates other factories.

## Problem
Imagine you're building a cross-platform UI library that needs to render buttons, checkboxes, and windows differently on Windows, macOS, and Linux. Each platform has its own look and feel, but your application code shouldn't be littered with platform-specific conditionals. You need a way to create entire families of related UI components that work together cohesively.

Without this pattern, you might end up with code like:
```go
if platform == "windows" {
    button = new WindowsButton()
    checkbox = new WindowsCheckbox()
} else if platform == "macos" {
    button = new MacButton()
    checkbox = new MacCheckbox()
}
```

This becomes unmaintainable as you add more platforms and components.

## Why Use This Pattern?

**Benefits:**
- **Consistency**: Ensures that products from the same family are used together (no mixing Windows buttons with Mac checkboxes)
- **Isolation**: Isolates concrete classes from client code
- **Easy to extend**: Adding new product families (platforms) is straightforward
- **Single Responsibility**: Product creation code is in one place

## When to Use

Use Abstract Factory when:
- Your system needs to work with multiple families of related products
- You want to enforce that products from a family are used together
- You want to provide a library of products and reveal only interfaces, not implementations
- You're building plugins or themes where families of objects need to work together
- You need to configure systems with one of multiple families of products

**Real-world scenarios:**
- Cross-platform UI frameworks (Windows/Mac/Linux widgets)
- Database drivers (MySQL/PostgreSQL/SQLite connections, statements, results)
- Cloud providers (AWS/Azure/GCP services: storage, compute, networking)
- Document converters (PDF/Word/HTML: parsers, formatters, writers)
- Game engines (Medieval/SciFi/Fantasy: characters, weapons, environments)

## When NOT to Use

Avoid this pattern when:
- You only have one family of products (use simple Factory instead)
- Product families don't need to work together consistently
- The number of product families is fixed and unlikely to grow
- Adding new product types (not families) is more common than adding families

## Implementation Guidelines

**Key Components:**
1. **Abstract Factory Interface**: Declares creation methods for each product type
2. **Concrete Factories**: Implement creation methods for a specific product family
3. **Abstract Product Interfaces**: Declare interfaces for each product type
4. **Concrete Products**: Implement products for specific families
5. **Client**: Uses only abstract interfaces

**Go-Specific Considerations:**
- Use interfaces for both factories and products
- Factory functions can return interfaces, hiding concrete types
- Consider using functional options for factory configuration
- No need for abstract classes—interfaces are sufficient

## Visual Schema

```
Client Code
    |
    | uses
    v
AbstractFactory (interface)
    |
    +-- CreateProductA() ProductA
    +-- CreateProductB() ProductB
    |
    +----------------+------------------+
    |                |                  |
ConcreteFactory1  ConcreteFactory2  ConcreteFactory3
    |                |                  |
    creates          creates            creates
    |                |                  |
    v                v                  v
ProductA1        ProductA2          ProductA3
ProductB1        ProductB2          ProductB3

All ProductA* implement ProductA interface
All ProductB* implement ProductB interface
```

**Flow Diagram:**
```
1. Client requests factory (e.g., UIFactory for "macos")
2. Factory system returns MacUIFactory (concrete implementation)
3. Client calls factory.CreateButton() → receives MacButton (as Button interface)
4. Client calls factory.CreateCheckbox() → receives MacCheckbox (as Checkbox interface)
5. All products from same family work together consistently
```

## Real-World Examples

### 1. Database Abstraction Layer
Different database drivers (MySQL, PostgreSQL) provide connections, statements, and result sets that must work together.

### 2. Cloud Provider SDKs
AWS, Azure, and GCP each provide families of services (storage, compute, messaging) that need consistent APIs within a provider.

### 3. UI Theme Systems
Applications that support multiple themes (Light/Dark/HighContrast) where all UI components must match the selected theme.

## Key Advantages

- **Product Family Consistency**: Guaranteed that products work together
- **Loose Coupling**: Client code doesn't depend on concrete classes
- **Open/Closed Principle**: Easy to add new product families without modifying existing code
- **Single Point of Control**: Product creation logic centralized in factories

## Key Gotchas

- **Complexity**: Adds many new interfaces and classes
- **Rigidity**: Adding new product types requires changing all factory interfaces
- **Overhead**: May be overkill for simple scenarios with few product variations
- **Interface Explosion**: Large product families create many interfaces
- **Testing**: Need to test each factory implementation separately

## Go Idioms

Go's interface satisfaction is implicit, making Abstract Factory particularly clean:
- No need for explicit inheritance hierarchies
- Small, focused interfaces align with Go philosophy
- Factory functions can return interface types
- Easy to mock for testing

## Code Example Structure

The implementation in `abstractfactory.go` demonstrates:
- UI component factory for multiple platforms (Windows, Mac)
- Product interfaces (Button, Checkbox)
- Concrete products for each platform
- Factory interface and concrete factory implementations
- Client code that works with any factory
