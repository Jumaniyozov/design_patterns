# Composite Pattern

## Overview

The Composite Pattern is a structural design pattern that lets you compose objects into tree structures to represent part-whole hierarchies. It allows clients to treat individual objects and compositions of objects uniformly.

## Problem

When working with hierarchical or tree-structured data, you face complexity in handling both individual items and groups:

- **Type switching**: Code must check whether an object is a leaf or composite
- **Recursive traversal**: Difficult to implement operations that work on trees
- **Inconsistent handling**: Leaf and composite objects are treated differently
- **Complexity duplication**: Similar logic needed for both leaves and composites
- **Tree navigation**: Traversing and manipulating hierarchies becomes tedious

### Real-World Context

Consider a file system with files and directories. A file is a leaf node (can't contain other files), while a directory is a composite node (can contain files and other directories). You want to perform operations like calculating total size or listing contents uniformly on both types.

## Why Use This Pattern?

- **Uniform interface**: Handle leaves and composites through the same interface
- **Recursive elegance**: Operations on trees naturally use recursion
- **Easy traversal**: Navigate hierarchies without type checking
- **Flexible composition**: Build complex structures from simple parts
- **Client simplicity**: Clients don't know if they're dealing with leaf or composite

## When to Use

- Tree or hierarchical structures (files, organizational charts, UI components)
- Operations that work uniformly on individual items and collections
- Building complex objects from simpler components
- Part-whole hierarchies where parts can contain other parts
- Recursive data structures (ASTs, DOM trees, menu systems)

## When NOT to Use

- Flat structures where composition adds unnecessary complexity
- Performance critical code where tree traversal overhead matters
- Simple parent-child relationships (just use direct pointers)
- When leaves and composites need very different behavior
- Type information must be preserved for different handling

## Implementation Guidelines

1. **Component interface**: Define operations for both leaf and composite
2. **Leaf objects**: Implement component interface (no children)
3. **Composite objects**: Implement component interface and contain children
4. **Child management**: Add/remove/get children operations in composite
5. **Recursive operations**: Default implementations work on both types

## Go Idioms

Go's interface-based design makes composites natural:

- Small focused interfaces describe what operations are available
- Implicit satisfaction means any type implementing the interface works
- Composition over inheritance is idiomatic
- No inheritance hierarchy needed

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│       Composite Pattern Structure            │
└──────────────────────────────────────────────┘

Component Interface:
    ┌─────────────────────┐
    │ Operation()         │
    │ GetName()           │
    └─────────────────────┘
           ▲
           │
    ┌──────┴──────────┐
    │                 │
    ▼                 ▼
┌──────────┐    ┌──────────────┐
│   Leaf   │    │  Composite   │
│          │    │              │
│ value    │    │ children []  │
│          │    │              │
└──────────┘    │ Add/Remove   │
                │ Iterate      │
                └──────────────┘
                       │
                    ┌──┴──┐
                    │     │
                 Leaf  Composite
                         │
                      ┌──┴──┐
                      │     │
                    Leaf  Leaf

File System Example:

                    /root/
                      │
        ┌─────────────┼──────────────┐
        │             │              │
     file1        dir1/          file2
        │             │
        │     ┌───────┴────────┐
        │     │                │
        │   file3           dir2/
        │                      │
        │                    file4

Operations work the same on all nodes:
- GetSize(): Returns file size OR sum of children sizes
- ListContents(): Returns name OR lists all children
- Delete(): Removes file OR directory and contents
```

## Real-World Examples

### 1. File System

```go
type FileNode interface {
    GetName() string
    GetSize() int64
    Accept(visitor Visitor)
}

type File struct {
    name string
    size int64
}

type Directory struct {
    name     string
    children []FileNode
}

// Both implement FileNode
// Size: file returns its size, directory sums children
```

### 2. UI Component Tree

```go
type UIComponent interface {
    Render() string
    GetWidth() int
    GetHeight() int
}

type Button struct { /*...*/ }      // Leaf
type Panel struct {                 // Composite
    components []UIComponent
}
```

### 3. Organization Hierarchy

```go
type Member interface {
    GetName() string
    GetSalary() float64
    GetReports() []Member  // Empty for individuals
}

type Employee struct {
    name   string
    salary float64
}

type Manager struct {
    name    string
    salary  float64
    reports []Member
}
```

## Key Advantages

- **Simplicity**: Treat leaves and composites uniformly
- **Flexibility**: Build complex structures from simple parts
- **Easy to extend**: Add new component types without changing tree logic
- **Recursive elegance**: Natural recursive algorithms
- **Client code simplicity**: No type checking needed
- **Real-world mapping**: Natural representation of hierarchies

## Key Gotchas

- **Circular references**: Need to prevent cycles in tree structure
- **Performance**: Deep trees can cause stack overflow or slow traversal
- **Complexity**: Can be overkill for simple parent-child relationships
- **Type information loss**: Must support different behaviors for leaves vs. composites
- **Memory overhead**: Each object adds memory (especially with many small leaves)
- **Modification during iteration**: Changing children while traversing causes issues
- **Null safety**: Handling empty composites and null children requires care
