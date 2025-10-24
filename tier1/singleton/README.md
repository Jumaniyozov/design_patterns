# Singleton Pattern

## Overview

The Singleton Pattern is a creational design pattern that ensures a class has only one instance while providing a global point of access to it. In Go, this pattern is used to manage shared resources like database connections, configuration, and caches.

## Problem

Some objects need to be shared globally but should exist only once:

- **Resource management**: Database connections, thread pools, file handles are expensive
- **Shared state**: Configuration should be consistent across the application
- **Multiple instances problem**: Creating multiple instances wastes resources and causes inconsistency
- **Lazy initialization**: Some resources shouldn't be created until actually needed
- **Thread safety**: Multiple goroutines must safely access the singleton

### Real-World Context

Consider a database connection pool. You want exactly one instance shared across your entire application to avoid resource exhaustion. Creating multiple pools would waste connections and cause consistency issues. The singleton pattern ensures only one pool exists and provides safe access to it.

## Why Use This Pattern?

- **Resource efficiency**: Expensive resources created only once
- **Consistent state**: Single source of truth for shared data
- **Global access**: Easy to access from anywhere in the application
- **Lazy initialization**: Create resources only when needed
- **Thread safety**: Go's `sync.Once` makes safe initialization simple

## When to Use

- Shared resources (database connections, file handles, thread pools)
- Global configuration that shouldn't change during runtime
- Caches and registries
- Logging systems
- Session management
- Feature flags and configuration services

## When NOT to Use

- Objects that need multiple instances with different state
- Stateless utility objects (just use functions)
- When dependency injection is more appropriate
- Testing scenarios where you need different instances
- Highly scalable systems where singletons become bottlenecks

## Implementation Guidelines

1. **Private constructor**: Prevent direct instantiation
2. **Static creation**: Provide a way to get the single instance
3. **Lazy or eager initialization**: Decide when to create the instance
4. **Thread safety**: Use `sync.Once` for safe concurrent access
5. **Reset capability**: Consider providing a way to reset for testing

## Go Idioms

Go's approach to singletons differs from traditional languages:
- Use `sync.Once` for thread-safe lazy initialization
- Package-level variables initialized in `init()` for eager initialization
- Avoid global variables when possible; prefer dependency injection
- Constructor functions that manage instance state

## Visual Schema

```
┌────────────────────────────────────────────┐
│       Singleton Pattern Structure          │
└────────────────────────────────────────────┘

Approach 1: sync.Once (Lazy Initialization)

package db

var (
    instance *Database
    once     sync.Once
)

func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{...}
    })
    return instance
}

First Call:        Subsequent Calls:
  Client             Client
    │                  │
    ├─ GetInstance()   ├─ GetInstance()
    │    │             │    │
    │    ├─ once.Do()  │    └─ once already executed
    │    │  (execute)  │        │
    │    ├─ Create     │        ├─ Return cached instance
    │    │  instance   │        │
    │    └─ Cache it   │        └─ Fast return
    │                  │


Approach 2: init() (Eager Initialization)

package config

var (
    Config *Configuration
)

func init() {
    Config = &Configuration{...}
}


Connection Pool Lifecycle:

Application Start
        │
        ▼
┌──────────────────┐
│ GetDatabase()    │ (First call)
│ calls once.Do()  │
└────────┬─────────┘
         │
         ├─ Database not initialized
         │
         ├─ Create connection pool
         │
         ├─ Store in global variable
         │
         ├─ Return reference
         │
All subsequent calls return the same instance without recreating
```

## Real-World Examples

### 1. Database Connection Pool

```go
var (
    db *sql.DB
    once sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        var err error
        db, err = sql.Open("postgres", "...")
        // Handle error
    })
    return db
}
```

### 2. Global Configuration

```go
var Config *AppConfig

func init() {
    Config = &AppConfig{
        Port: 8080,
        Database: "localhost",
    }
}
```

### 3. Logger Instance

```go
var logger *Logger

func GetLogger() *Logger {
    once.Do(func() {
        logger = NewLogger("app.log")
    })
    return logger
}
```

## Key Advantages

- **Resource efficiency**: Expensive resources created only once
- **Consistent state**: Single point of truth for shared data
- **Lazy initialization**: Resources created only when needed
- **Global access**: Easy to use throughout application
- **Thread safety**: `sync.Once` ensures safe concurrent access
- **Memory efficient**: No duplicated resources

## Key Gotchas

- **Testing complexity**: Hard to test with singletons; provide reset methods
- **Hidden dependencies**: Global state makes dependencies implicit
- **Concurrency bottleneck**: Single instance may become a bottleneck
- **Mutable state**: Singleton state changes affect all users
- **Initialization errors**: Must handle initialization failures gracefully
- **Tight coupling**: Code couples to global state
- **Not always appropriate**: Often dependency injection is better design

## Best Practices

1. **Use sparingly**: Only for true shared resources
2. **Make thread-safe**: Always use `sync.Once` or init-based initialization
3. **Provide reset for testing**: Add `ResetForTesting()` method when needed
4. **Document clearly**: Make it obvious this is a singleton
5. **Consider DI first**: Dependency injection often leads to better design
