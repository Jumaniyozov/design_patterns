# Visitor Pattern

## Overview
The Visitor pattern separates algorithms from the objects they operate on, allowing new operations to be added without modifying the object structures. It's about moving operations out of element classes into visitor classes.

## Problem
You have a complex object structure (like a document tree, AST, or composite) and need to perform various operations on it. Adding each operation to element classes violates Single Responsibility Principle and makes classes bloated. Each new operation requires modifying all element classes.

## Why Use This Pattern?

**Benefits:**
- **Open/Closed Principle**: Add operations without modifying elements
- **Single Responsibility**: Separate algorithms from data structures
- **Gather Related Operations**: Group related behavior in visitor
- **Easy to Add Operations**: New visitor = new operation

## When to Use

Use Visitor when:
- Object structure is stable but operations change frequently
- Many distinct operations needed on object structure
- Object structure contains many classes with different interfaces
- Operations should be defined outside the classes

**Real-world scenarios:**
- Compiler AST processing (type checking, code generation)
- Document structure processing (export, rendering, validation)
- Scene graph operations (rendering, collision detection)
- File system traversal (search, backup, analysis)

## When NOT to Use

Avoid when:
- Object structure changes frequently (all visitors need updates)
- Few operations needed
- Operations are naturally part of element classes
- Visitor pattern complexity isn't justified

## Code Example Structure

Demonstrates:
- Document element visitor (export, word count)
- Shape visitor (area calculation, rendering)
- File system visitor (size calculation, search)
