# Memento Pattern

## Overview
The Memento pattern captures and externalizes an object's internal state without violating encapsulation, allowing the object to be restored to this state later. It's the foundation of undo/redo functionality.

## Problem
You need to save and restore object state (for undo/redo, checkpoints, snapshots) without exposing internal implementation details. Direct access to fields violates encapsulation, while serialization may expose sensitive data.

## Why Use This Pattern?

**Benefits:**
- **Encapsulation**: Preserves object's encapsulation boundary
- **State Management**: Clean snapshot and restore mechanism
- **Undo/Redo**: Foundation for undo/redo systems
- **Checkpointing**: Save states at specific points

## When to Use

Use Memento when:
- You need to save/restore object state
- Direct field access would violate encapsulation
- Implementing undo/redo functionality
- Creating checkpoints or snapshots

**Real-world scenarios:**
- Text editors (undo/redo)
- Game save systems
- Database transactions (rollback)
- Configuration snapshots

## When NOT to Use

Avoid when:
- State is simple and public
- Memory overhead is too high
- Serialization is sufficient

## Code Example Structure

Demonstrates:
- Text editor with undo/redo
- Game state snapshots
- Configuration rollback
