# Flyweight Pattern

## Overview
The Flyweight pattern minimizes memory usage by sharing common data among multiple objects. It's useful when you need to create a large number of similar objects that share substantial state, allowing you to support more objects by storing shared state externally.

## Problem
Imagine a text editor that needs to represent every character on screen as an object with font, size, color, and position. A document with 100,000 characters would create 100,000 objects, most sharing the same font and formatting. This wastes enormous memory.

Or consider a game with thousands of trees—each tree object stores species, texture, mesh, but also position and health. Most trees of the same species share identical species/texture/mesh data.

## Why Use This Pattern?

**Benefits:**
- **Memory Efficiency**: Dramatically reduce memory usage for large numbers of similar objects
- **Performance**: Fewer objects mean less garbage collection pressure
- **Scalability**: Support massive numbers of objects that would otherwise exhaust memory

## When to Use

Use Flyweight when:
- Application uses a large number of objects
- Storage costs are high due to sheer quantity
- Most object state can be made extrinsic (moved outside)
- Many groups of objects share intrinsic state
- Application doesn't depend on object identity (shared objects are acceptable)

**Real-world scenarios:**
- Text editors (character glyphs with shared formatting)
- Game engines (particles, trees, enemies with shared resources)
- Map applications (map tiles, markers with shared icons)
- GUI toolkits (shared fonts, icons, cursors)

## When NOT to Use

Avoid this pattern when:
- Objects don't share significant state
- Few objects are created
- Memory usage isn't a concern
- Objects need to be independently mutable
- Complexity outweighs memory savings

## Implementation Guidelines

**Key Components:**
1. **Flyweight Interface**: Operations that use extrinsic state
2. **Concrete Flyweight**: Stores intrinsic (shared) state
3. **Flyweight Factory**: Creates and manages flyweight objects
4. **Client**: Maintains extrinsic state and passes it to flyweights

**Go-Specific Considerations:**
- Use maps for flyweight factories/pools
- Sync.Map for concurrent access
- Careful about when to share vs. copy
- Consider using pointers to shared data

**Intrinsic vs Extrinsic State:**
- **Intrinsic**: Shared state, stored in flyweight (e.g., tree species, font name)
- **Extrinsic**: Unique state, stored by client (e.g., tree position, character location)

## Visual Schema

```
FlyweightFactory
    |
    | maintains pool of flyweights
    |
    +-- map[key]Flyweight
            |
            +-- Flyweight (shared intrinsic state)
                    ^
                    |
                    | uses
                    |
Client (stores extrinsic state)
    - position
    - health
    - uses flyweight for shared data

Memory Comparison:
Without Flyweight: N objects × (intrinsic + extrinsic) = Large
With Flyweight: (unique intrinsic objects) + (N × extrinsic) = Small
```

## Real-World Examples

### 1. Text Formatting
Characters share Font objects rather than each character storing font name, size, style.

### 2. Game Particle Systems
Thousands of particles share texture/sprite data, only positions differ.

### 3. Icon Caching
Desktop environments share icon bitmaps across all file entries showing same file type.

## Key Advantages

- **Massive Memory Savings**: Orders of magnitude reduction
- **Better Performance**: Less memory allocation and GC pressure
- **More Objects**: Support numbers that would otherwise be impossible

## Key Gotchas

- **Complexity**: Separating intrinsic/extrinsic state adds complexity
- **Runtime Cost**: Extra indirection to access shared state
- **Thread Safety**: Shared state requires synchronization
- **Immutability**: Shared state should be immutable to avoid issues

## Go Idioms

```go
type TreeType struct { // Flyweight - intrinsic state
    name    string
    texture string
}

type Tree struct { // Context - extrinsic state
    x, y     float64
    treeType *TreeType // Reference to flyweight
}
```

## Code Example Structure

The implementation demonstrates:
- Character rendering with shared glyph data
- Tree forest with shared tree species
- Flyweight factory managing shared objects
- Memory usage comparison
