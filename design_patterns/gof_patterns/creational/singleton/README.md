# Singleton Pattern

## Overview

The Singleton pattern ensures a class has only one instance and provides a global point of access to it. In Go, this is achieved through package-level variables, initialization functions, or `sync.Once` for thread-safe lazy initialization.

## Problem

Imagine you need to manage a shared resource like a database connection pool, application configuration, or logging system. Creating multiple instances would:
- Waste memory and resources
- Lead to inconsistent state across the application
- Cause race conditions in concurrent environments
- Make debugging and testing difficult

You need exactly ONE instance shared across your entire application, accessible globally, and initialized safely even in concurrent scenarios.

## Why Use This Pattern?

- **Controlled Access**: Single instance ensures controlled access to shared resources
- **Reduced Memory Footprint**: Prevents unnecessary instantiation of expensive objects
- **Global State Management**: Provides a well-known access point for shared state
- **Lazy Initialization**: Resources can be initialized only when first needed
- **Thread Safety**: Go's `sync.Once` ensures safe initialization in concurrent environments

## When to Use

- **Database Connection Pools**: One pool managing all database connections
- **Configuration Objects**: Single source of truth for application config
- **Logger Instances**: Centralized logging with consistent formatting
- **Cache Managers**: Shared in-memory cache across the application
- **Hardware Interface Managers**: Printer spooler, file system access, device drivers

## When NOT to Use

- **Testability Concerns**: Singletons introduce global state making unit testing harder
- **Hidden Dependencies**: Makes dependencies less explicit in function signatures
- **Tight Coupling**: Creates implicit dependencies throughout codebase
- **Simple Use Cases**: When dependency injection or simple package variables suffice
- **Microservices**: Each service instance should manage its own resources

## Implementation Guidelines

In Go, there are three common approaches:

### 1. Package-Level Variable with init()
```go
var instance *Singleton

func init() {
    instance = &Singleton{}
}

func GetInstance() *Singleton {
    return instance
}
```

### 2. sync.Once (Thread-Safe Lazy Loading)
```go
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

### 3. Package-Level Variable (Eager Initialization)
```go
var instance = &Singleton{}

func GetInstance() *Singleton {
    return instance
}
```

## Go Idioms

- **Package Privacy**: Use unexported struct types with exported accessor functions
- **sync.Once**: Preferred for thread-safe lazy initialization
- **No Constructors**: Go doesn't have constructors; use init() or package-level initialization
- **Interface Satisfaction**: Make singleton satisfy interfaces for better testability

## Visual Schema

```
Application Components
    │
    ├─→ Component A ──┐
    │                 │
    ├─→ Component B ──┼─→ GetInstance() ──→ [Singleton Instance]
    │                 │          ↑                    │
    ├─→ Component C ──┘          │                    │
    │                             │                    │
    └─→ Component D ─────────────┘                    │
                                                       │
                                          All components share
                                          the SAME instance
```

### Thread-Safe Initialization Flow

```
Thread 1                Thread 2                Thread 3
    │                       │                       │
    ├─→ GetInstance()       ├─→ GetInstance()       ├─→ GetInstance()
    │        │              │        │              │        │
    │        ↓              │        ↓              │        ↓
    │    sync.Once.Do()     │    sync.Once.Do()     │    sync.Once.Do()
    │        │              │        │              │        │
    │        ↓              │        ↓              │        ↓
    │   [First Call]        │   [Waits...]          │   [Waits...]
    │   Initialize          │        │              │        │
    │        │              │        ↓              │        ↓
    │        └──────────────┴──→ Return Instance ←──┘
    │                              (Same Instance)
```

## Real-World Examples

### 1. Database Connection Pool
```go
type DBPool struct {
    connections []*sql.DB
    maxConn     int
}

var (
    dbInstance *DBPool
    dbOnce     sync.Once
)

func GetDBPool() *DBPool {
    dbOnce.Do(func() {
        dbInstance = &DBPool{
            connections: make([]*sql.DB, 0),
            maxConn:     10,
        }
    })
    return dbInstance
}
```

### 2. Application Configuration
```go
type Config struct {
    APIKey      string
    DatabaseURL string
    Port        int
}

var config = loadConfigFromFile()

func GetConfig() *Config {
    return config
}
```

### 3. Logger Instance
```go
type Logger struct {
    writer io.Writer
    level  string
}

var logger = &Logger{
    writer: os.Stdout,
    level:  "INFO",
}

func GetLogger() *Logger {
    return logger
}
```

## Key Advantages

✓ **Guaranteed Single Instance**: Only one instance exists throughout application lifecycle
✓ **Thread-Safe Access**: sync.Once ensures safe concurrent initialization
✓ **Lazy Initialization**: Resources created only when needed
✓ **Global Access Point**: Easy access from anywhere in codebase
✓ **Resource Efficiency**: Prevents redundant instantiation of expensive objects

## Key Gotchas

⚠️ **Testing Challenges**: Global state makes unit tests dependent and hard to isolate
⚠️ **Hidden Dependencies**: Not clear from function signatures what depends on singleton
⚠️ **Tight Coupling**: Components become tightly coupled to singleton implementation
⚠️ **Concurrency Issues**: If not using sync.Once, race conditions can create multiple instances
⚠️ **Initialization Order**: Package-level initialization order can be tricky
⚠️ **Overuse**: Easy to abuse; consider dependency injection for better testability

## Best Practices

1. **Use sync.Once** for thread-safe lazy initialization
2. **Export Interfaces**, not concrete types, for better testability
3. **Consider Dependency Injection** as an alternative for better testing
4. **Document Why** the singleton is necessary
5. **Provide Reset** methods for testing scenarios (carefully!)
6. **Keep It Simple**: Don't overcomplicate with unnecessary features
