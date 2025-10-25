# CQRS Pattern (Command Query Responsibility Segregation)

## Overview
CQRS separates read operations (queries) from write operations (commands) using different models. Commands modify state, queries read state, and they use separate optimized data stores.

## Problem
Traditional CRUD uses the same model for reads and writes, leading to conflicts: writes need validation/business logic, reads need joins/aggregations. One model can't optimize for both.

## Why Use This Pattern?
- **Optimization**: Separate read/write optimization
- **Scalability**: Scale reads and writes independently
- **Simplicity**: Simpler models for each concern
- **Flexibility**: Different storage for reads vs writes
- **Performance**: Denormalized read models

## When to Use
- Read/write patterns are very different
- Need independent scaling
- Complex business logic on writes
- Multiple read projections needed
- High-performance requirements

## Real-world scenarios
- E-commerce (complex writes, fast reads)
- Financial systems (audited writes, fast queries)
- Social media (event-driven writes, optimized feeds)
- Analytics (write-heavy ingestion, read-heavy queries)

## When NOT to Use
- Simple CRUD applications
- Read/write patterns are similar
- Added complexity not justified
- Team unfamiliar with pattern

## Go Idioms
```go
// Command side
type CreateOrderCommand struct {
    OrderID string
    Items   []Item
}

// Query side
type OrderQuery struct {
    OrderID string
}

type CommandBus interface {
    Send(cmd Command) error
}

type QueryBus interface {
    Execute(query Query) (interface{}, error)
}
```
