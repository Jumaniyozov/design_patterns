# Iterator Pattern

## Overview
The Iterator pattern provides a way to access elements of a collection sequentially without exposing its underlying representation. It decouples iteration logic from the collection itself.

## Problem
In Go, the `range` keyword handles most iteration needs elegantly. However, Iterator pattern is still valuable for:
- Custom traversal algorithms (tree traversal, graph traversal)
- Lazy evaluation and infinite sequences
- Multiple concurrent iterations over same collection
- Hiding complex internal structure

## Why Use This Pattern?

**Benefits:**
- **Uniform Interface**: Single way to traverse different collections
- **Multiple Iterators**: Multiple simultaneous iterations over same collection
- **Encapsulation**: Hide collection's internal structure
- **Flexibility**: Support different traversal strategies

## When to Use

Use Iterator when:
- You need custom traversal logic (DFS, BFS, in-order, etc.)
- Collection structure is complex (tree, graph)
- You want lazy evaluation or infinite sequences
- Multiple concurrent iterations needed
- Client shouldn't depend on collection's internal structure

**Real-world scenarios:**
- Tree traversal (binary trees, N-ary trees)
- Graph algorithms (BFS, DFS)
- File system navigation
- Pagination through large datasets
- Stream processing

## When NOT to Use

Avoid this pattern when:
- Simple slice/map with `range` is sufficient
- Collection structure is simple
- Only one traversal strategy needed
- Go's native iteration is adequate

## Go Idioms

Go 1.23+ has range-over-function iterators. For earlier versions:
```go
type Iterator interface {
    HasNext() bool
    Next() interface{}
}
```

Or channel-based:
```go
func (c *Collection) Iterator() <-chan Item {
    ch := make(chan Item)
    go func() {
        defer close(ch)
        // Send items
    }()
    return ch
}
```

## Code Example Structure

Demonstrates:
- Tree traversal (in-order, pre-order, post-order)
- Channel-based iterators
- Lazy evaluation
- Custom collection iteration
