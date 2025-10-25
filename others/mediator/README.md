# Mediator Pattern

## Overview
The Mediator pattern defines an object that encapsulates how a set of objects interact. It promotes loose coupling by keeping objects from referring to each other explicitly, allowing their interaction to vary independently.

## Problem
Imagine a complex UI dialog with multiple components (buttons, checkboxes, text fields) that need to interact. Without Mediator, each component directly references others, creating tight coupling and a maintenance nightmare. Changes to one component ripple through all related components.

## Why Use This Pattern?

**Benefits:**
- **Reduced Coupling**: Components don't reference each other directly
- **Centralized Control**: All interaction logic in one place
- **Easier Maintenance**: Changes affect only mediator
- **Reusability**: Components more reusable independently

## When to Use

Use Mediator when:
- Objects communicate in complex but well-defined ways
- Reusing objects is difficult due to many dependencies
- Behavior distributed among classes should be customizable without subclassing

**Real-world scenarios:**
- UI dialog coordination
- Chat room systems
- Air traffic control
- Workflow engines

## When NOT to Use

Avoid when:
- Communication is simple
- Few objects interact
- Mediator becomes too complex (God object)

## Code Example Structure

Demonstrates:
- Chat room mediator
- UI dialog coordination
- Air traffic control system
