# Null Object Pattern

## Overview
The Null Object pattern provides a default object that does nothing, eliminating the need for null checks. It implements the same interface as real objects but with neutral/no-op behavior.

## Problem
Checking for nil everywhere clutters code and causes potential nil pointer panics. You need a safe default that behaves correctly without special handling.

## Why Use This Pattern?
- **Eliminate Nil Checks**: No more `if obj != nil`
- **Simplify Code**: Clean, linear logic flow
- **Safe Defaults**: Prevent nil pointer panics
- **Polymorphism**: Null object behaves like real objects

## When to Use
- Optional dependencies
- Default behaviors
- Avoiding null checks
- Chain of responsibility termination

## Real-world scenarios
- Logging (NullLogger for disabled logging)
- Caching (NullCache for no-cache mode)
- Notifications (NullNotifier for silent mode)

## Go Idioms
```go
type Logger interface {
    Log(msg string)
}

type NullLogger struct{}
func (n *NullLogger) Log(msg string) {} // no-op
```
