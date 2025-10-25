# Event Sourcing Pattern

## Overview
Event Sourcing stores the state of an application as a sequence of events rather than just the current state. The current state is derived by replaying events from the beginning.

## Problem
Traditional CRUD loses history—you only see current state. Auditing is difficult, debugging past issues is impossible, and you can't recompute state with new business rules.

## Why Use This Pattern?
- **Complete Audit Trail**: Every state change recorded
- **Time Travel**: Reconstruct state at any point in time
- **Event-Driven Architecture**: Natural fit for event-driven systems
- **Debugging**: Replay events to reproduce bugs
- **Business Intelligence**: Analyze historical event patterns

## When to Use
- Audit requirements are critical
- Need to reconstruct past states
- Event-driven architecture
- Complex domain logic
- Multiple projections of same data

## Real-world scenarios
- Banking (transaction history)
- E-commerce (order lifecycle tracking)
- Version control systems
- Collaborative editing
- Gaming (replay functionality)

## When NOT to Use
- Simple CRUD applications
- Event history not valuable
- Query performance critical (need snapshots)
- Storage constraints

## Go Idioms
```go
type Event interface {
    EventType() string
    Timestamp() time.Time
}

type EventStore interface {
    AppendEvent(aggregateID string, event Event) error
    GetEvents(aggregateID string) ([]Event, error)
}
```
