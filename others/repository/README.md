# Repository Pattern

## Overview
The Repository pattern mediates between the domain and data mapping layers, providing a collection-like interface for accessing domain objects. It encapsulates data access logic and promotes separation of concerns.

## Problem
Direct database access in business logic creates tight coupling, makes testing difficult, and scatters data access code. You need an abstraction that hides data access details.

## Why Use This Pattern?
- **Separation of Concerns**: Isolate data access from business logic
- **Testability**: Easy to mock repositories
- **Centralized Data Access**: All queries in one place
- **Domain-Focused**: Work with domain objects, not SQL/tables

## When to Use
- Domain-driven design
- Complex data access logic
- Multiple data sources
- Need to swap data sources

## Real-world scenarios
- User management (UserRepository)
- Product catalog (ProductRepository)
- Order processing (OrderRepository)
- Multi-database applications

## Go Idioms
```go
type UserRepository interface {
    FindByID(id string) (*User, error)
    FindAll() ([]*User, error)
    Save(user *User) error
    Delete(id string) error
}
```
