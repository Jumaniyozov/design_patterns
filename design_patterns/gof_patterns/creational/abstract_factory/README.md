# Abstract Factory Pattern

## Overview

The Abstract Factory pattern provides an interface for creating families of related or dependent objects without specifying their concrete classes. It's a factory of factories, creating groups of related objects that work together.

## Problem

You're building a cross-platform UI library that needs to create components for different operating systems: Windows, Mac, and Linux. Each platform has its own:
- Button style and behavior
- Checkbox appearance
- Text field rendering
- Scroll bar implementation

Without Abstract Factory:
```go
if platform == "windows" {
    button = &WindowsButton{}
    checkbox = &WindowsCheckbox{}
    // ... must remember to use Windows components together
} else if platform == "mac" {
    button = &MacButton{}
    checkbox = &MacCheckbox{}
    // ... risk of mixing incompatible components
}
```

You need to ensure related components are created together and are compatible with each other.

## Why Use This Pattern?

- **Consistency**: Ensures related objects are used together
- **Isolation**: Client code independent of concrete classes
- **Easy Switching**: Change entire product family at once
- **Guaranteed Compatibility**: Products from same factory work together
- **Single Responsibility**: Product creation logic in one place

## When to Use

- **Product Families**: Need to create families of related objects
- **Platform Independence**: Cross-platform applications (Windows, Mac, Linux)
- **Theme Systems**: UI themes (dark mode, light mode, high contrast)
- **Database Abstraction**: Different database vendors (MySQL, PostgreSQL, MongoDB)
- **Document Types**: Different document formats (PDF, HTML, Markdown)
- **Cloud Providers**: Abstract AWS, Azure, GCP services

## When NOT to Use

- **Single Product**: Only creating one type of object
- **No Families**: Products don't need to be related
- **Simple Factory Sufficient**: When Factory Method pattern is enough
- **Over-engineering**: Adding unnecessary complexity

## Implementation Guidelines

### Basic Pattern

```go
// Abstract products
type Button interface {
    Render()
}

type Checkbox interface {
    Render()
}

// Abstract factory
type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}

// Concrete factory for Windows
type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button {
    return &WindowsButton{}
}

func (f *WindowsFactory) CreateCheckbox() Checkbox {
    return &WindowsCheckbox{}
}

// Client code
func RenderUI(factory GUIFactory) {
    button := factory.CreateButton()
    checkbox := factory.CreateCheckbox()
    button.Render()
    checkbox.Render()
}
```

## Visual Schema

```
Client
  │
  │ uses
  ↓
AbstractFactory
  ├─ CreateProductA()
  └─ CreateProductB()
        ↑
        │ implements
        │
    ┌───┴────┬────────┐
    │        │        │
ConcreteFactory1  ConcreteFactory2  ConcreteFactory3
    │        │        │
    ├─→ ProductA1   ProductA2   ProductA3
    └─→ ProductB1   ProductB2   ProductB3

Each factory creates a family of compatible products
```

## Real-World Examples

### 1. Cross-Platform UI
```go
type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
    CreateTextField() TextField
}

windowsFactory := &WindowsGUIFactory{}
macFactory := &MacGUIFactory{}
```

### 2. Database Abstraction
```go
type DatabaseFactory interface {
    CreateConnection() Connection
    CreateQuery() QueryBuilder
    CreateTransaction() Transaction
}
```

### 3. Document Generation
```go
type DocumentFactory interface {
    CreateDocument() Document
    CreateParagraph() Paragraph
    CreateHeading() Heading
}
```

## Key Advantages

✓ **Isolation**: Concrete classes isolated from client code
✓ **Consistency**: Ensures compatible products used together
✓ **Easy Switching**: Change entire family by changing factory
✓ **Open/Closed**: Easy to introduce new product families
✓ **Single Responsibility**: Creation logic centralized

## Key Gotchas

⚠️ **Complexity**: More complex than Factory Method
⚠️ **Rigid Structure**: Hard to add new product types
⚠️ **Overkill**: Too complex for simple scenarios
⚠️ **Interface Changes**: Adding product type requires changing all factories

## Best Practices

1. **Return Interfaces**: Factories return product interfaces, not concrete types
2. **Consistent Families**: Ensure products in a family work together
3. **Configuration**: Use configuration to select factory at runtime
4. **Combine with Singleton**: Factory itself can be singleton
5. **Document Compatibility**: Clearly document which products work together
