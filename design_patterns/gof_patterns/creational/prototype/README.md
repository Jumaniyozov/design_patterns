# Prototype Pattern

## Overview

The Prototype pattern creates new objects by cloning existing instances rather than creating from scratch. It specifies the kinds of objects to create using a prototypical instance and creates new objects by copying this prototype.

## Problem

You're building a document editor where users create complex objects like diagrams, charts, and formatted text blocks. Creating these from scratch involves:
- Complex initialization
- Many configuration steps
- Expensive resource loading
- Identical or similar setup

Without Prototype:
```go
// Recreating complex objects repeatedly
doc1 := CreateComplexDocument() // 50 lines of setup
doc2 := CreateComplexDocument() // Same 50 lines again
// Tedious, error-prone, inefficient
```

You need a way to create similar objects quickly by copying existing ones.

## Why Use This Pattern?

- **Performance**: Cloning faster than creating from scratch
- **Reduced Complexity**: Avoid complex initialization logic
- **Dynamic Creation**: Create objects at runtime without knowing their class
- **Preserve State**: Clone objects with current state
- **Avoid Subclassing**: Alternative to Factory patterns

## When to Use

- **Expensive Creation**: Object initialization is costly
- **Similar Objects**: Need many objects with similar configuration
- **Runtime Prototypes**: Object types determined at runtime
- **State Preservation**: Clone objects with specific state
- **Avoid Class Explosion**: Alternative to creating many subclasses

## When NOT to Use

- **Simple Objects**: Creation is trivial
- **Deep Copy Complexity**: Complex object graphs hard to clone
- **Unique Objects**: Each instance significantly different

## Implementation Guidelines

```go
// Prototype interface
type Prototype interface {
    Clone() Prototype
}

// Concrete prototype
type ConcretePrototype struct {
    field1 string
    field2 int
}

func (p *ConcretePrototype) Clone() Prototype {
    return &ConcretePrototype{
        field1: p.field1,
        field2: p.field2,
    }
}
```

## Real-World Examples

### 1. Document Cloning
```go
template := LoadComplexDocument()
doc1 := template.Clone()
doc2 := template.Clone()
```

### 2. Game Object Spawning
```go
enemyPrototype := &Enemy{health: 100, speed: 5}
enemy1 := enemyPrototype.Clone()
enemy2 := enemyPrototype.Clone()
```

## Key Advantages

✓ **Performance**: Fast object creation through cloning
✓ **Simplicity**: Avoid complex initialization code
✓ **Flexibility**: Add/remove prototypes at runtime
✓ **Preserve Configuration**: Clone pre-configured objects

## Key Gotchas

⚠️ **Deep vs Shallow Copy**: Be careful with nested objects
⚠️ **Circular References**: Can cause infinite loops
⚠️ **Mutable State**: Cloned objects share mutable state if not deep copied

## Best Practices

1. **Implement Clone() Method**: All prototypes should have Clone()
2. **Deep Copy When Needed**: Use deep copy for nested structures
3. **Copy Constructors**: Alternative to Clone() method
4. **Prototype Registry**: Store and manage prototypes centrally
5. **Document Copy Semantics**: Clearly specify shallow vs deep copy
