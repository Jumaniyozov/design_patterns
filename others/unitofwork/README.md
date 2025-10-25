# Unit of Work Pattern

## Overview
The Unit of Work pattern maintains a list of objects affected by a business transaction and coordinates writing out changes. It tracks changes during a transaction and commits them as a single unit.

## Problem
Multiple repository operations should succeed or fail together (atomicity). Without coordination, partial updates lead to inconsistent state. You need transactional boundaries.

## Why Use This Pattern?
- **Atomicity**: All changes succeed or all fail
- **Consistency**: Maintain data integrity
- **Performance**: Batch database operations
- **Change Tracking**: Know what changed during transaction

## When to Use
- Multiple related changes must be atomic
- Need transaction management
- Complex business operations spanning multiple entities
- Optimizing database round-trips

## Real-world scenarios
- E-commerce orders (order + line items + inventory + payment)
- User registration (user + profile + preferences)
- Financial transactions (debit + credit + audit log)

## Go Idioms
```go
type UnitOfWork interface {
    RegisterNew(entity interface{})
    RegisterDirty(entity interface{})
    RegisterDeleted(entity interface{})
    Commit() error
    Rollback()
}
```
