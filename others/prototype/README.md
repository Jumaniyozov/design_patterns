# Prototype Pattern

## Overview
The Prototype pattern creates new objects by cloning existing objects (prototypes) rather than constructing them from scratch. This is particularly useful when object creation is expensive or complex, and you need multiple similar instances.

## Problem
Imagine you're building a game where you need to spawn thousands of enemy characters. Each enemy has complex initialization—loading 3D models, setting up AI behaviors, configuring animation states, etc. Creating each enemy from scratch would be prohibitively slow.

Or consider a document editor where users can duplicate complex shapes, charts, or formatted text blocks. Each object might have dozens of properties and nested structures. Manually copying all properties is error-prone and couples your code to the concrete classes.

## Why Use This Pattern?

**Benefits:**
- **Performance**: Cloning can be much faster than constructing from scratch
- **Simplification**: Avoid complex initialization logic repetition
- **Runtime Configuration**: Create objects with specific configurations at runtime
- **Reduced Subclassing**: Avoid creating factory hierarchies for each variation
- **Independence**: Client code doesn't need to know concrete classes

## When to Use

Use Prototype when:
- Object creation is expensive (database queries, network calls, complex calculations)
- You need objects configured in specific ways that are easier to clone than construct
- You want to avoid building parallel class hierarchies of factories
- The system should be independent of how objects are created and represented
- Classes to instantiate are specified at runtime

**Real-world scenarios:**
- Game object spawning (enemies, projectiles, particles)
- Document editors (clone shapes, text blocks, images)
- Configuration management (clone preset configurations)
- Testing (create test fixtures from prototypes)
- GUI builders (duplicate UI components with their states)

## When NOT to Use

Avoid this pattern when:
- Objects are simple and cheap to construct
- Deep vs shallow copy complexity outweighs benefits
- Objects contain resources that shouldn't be duplicated (file handles, connections)
- The language/framework doesn't support efficient cloning
- Circular references make cloning complex

## Implementation Guidelines

**Key Components:**
1. **Prototype Interface**: Declares clone method
2. **Concrete Prototypes**: Implement clone to return a copy of themselves
3. **Client**: Creates new objects by cloning prototypes

**Go-Specific Considerations:**
- Go doesn't have built-in deep copy—implement carefully
- Use copy constructors or explicit Clone() methods
- Consider shallow vs deep copy based on your needs
- For deep copies of complex structures, consider serialization/deserialization
- Be careful with pointers, slices, maps—they need deep copying

**Shallow vs Deep Copy:**
```
Shallow Copy: Copy struct values, but pointers point to same objects
Deep Copy: Recursively copy all referenced objects
```

## Visual Schema

```
Client
   |
   | calls Clone()
   |
   v
Prototype (interface)
   |
   +-- Clone() Prototype
   |
   +------------------+------------------+
   |                  |                  |
ConcretePrototypeA  ConcretePrototypeB  ConcretePrototypeC
   |                  |                  |
Each implements     Clone() to        return copy
   own Clone()        of itself         of itself

Prototype Registry (optional):
  Map[string]Prototype
  - stores pre-configured prototypes
  - clients request clones by key
```

**Flow Diagram:**
```
1. Initialize prototypes with desired configurations
2. Store prototypes in registry (optional)
3. Client requests clone: registry.Get("enemyType").Clone()
4. Prototype creates and returns copy of itself
5. Client customizes the clone if needed
6. Repeat for as many clones as needed
```

## Real-World Examples

### 1. Game Enemy Spawning
Pre-configure different enemy types (weak, normal, boss) with different stats, then clone them rapidly during gameplay.

### 2. Database Connection Pooling
Clone a prototype connection configuration to create multiple connections with the same settings.

### 3. Document Templates
Provide template documents (invoice, report, letter) that users clone and customize rather than build from scratch.

## Key Advantages

- **Performance Boost**: Cloning faster than construction for complex objects
- **Simplicity**: Avoid complex initialization code duplication
- **Dynamic Configuration**: Add/remove prototypes at runtime
- **Hide Complexity**: Client doesn't need to know construction details

## Key Gotchas

- **Deep vs Shallow Copy**: Shallow copies share references—can cause unexpected mutations
- **Circular References**: Objects referencing each other make cloning complex
- **Resource Management**: Some resources (files, locks) shouldn't be cloned
- **Implementation Burden**: Must implement Clone for every prototype
- **Copy Constructor Maintenance**: Changes to struct require Clone() updates

## Go Idioms

In Go, implementing Prototype requires explicit work:
- No built-in clone mechanism like some languages
- Manual deep copy for pointer fields, slices, maps
- Can use encoding/gob or similar for deep copy serialization
- Copy constructors are idiomatic: `NewFromExisting(existing *Type)`

```go
// Explicit clone method
func (p *Person) Clone() *Person {
    return &Person{
        Name: p.Name,
        Age:  p.Age,
        Address: &Address{ // Deep copy
            Street: p.Address.Street,
            City:   p.Address.City,
        },
    }
}
```

## Code Example Structure

The implementation in `prototype.go` demonstrates:
- Prototype interface with Clone method
- Multiple concrete prototypes (game characters)
- Deep vs shallow copy examples
- Prototype registry for managing prototypes
- Performance comparison: construction vs cloning
