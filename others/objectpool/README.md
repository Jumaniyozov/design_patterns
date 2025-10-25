# Object Pool Pattern

## Overview
The Object Pool pattern manages a pool of reusable objects, reducing the overhead of creating and destroying expensive objects. Instead of creating new objects, clients acquire from the pool and return them when done.

## Problem
Creating certain objects is expensive—database connections, thread pools, large buffers, network connections. Creating them for every request wastes resources and degrades performance. You need a way to reuse these objects efficiently.

## Why Use This Pattern?

**Benefits:**
- **Performance**: Reuse expensive objects instead of recreating
- **Resource Management**: Control maximum number of objects
- **Memory Efficiency**: Reduce allocation/GC pressure
- **Predictable Performance**: Avoid creation spikes

## When to Use

Use Object Pool when:
- Object creation is expensive (time or resources)
- You need many short-lived instances of expensive objects
- Number of objects can be limited
- Objects can be safely reused after reset

**Real-world scenarios:**
- Database connection pools
- Thread pools / goroutine workers
- Buffer pools (bytes.Buffer, []byte)
- HTTP client pools
- Socket connection pools

## When NOT to Use

Avoid when:
- Object creation is cheap
- Objects can't be safely reused
- Pool management overhead exceeds benefits
- Objects are stateful and hard to reset

## Go Idioms

Go's standard library uses pools extensively:
- `sync.Pool` for temporary object caching
- Buffer pools in encoding packages
- Connection pools in database/sql

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}
```

## Code Example Structure

Demonstrates:
- Database connection pool
- Worker goroutine pool
- Buffer pool using sync.Pool
- Custom object pool with lifecycle management
