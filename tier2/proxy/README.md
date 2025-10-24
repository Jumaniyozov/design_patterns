# Proxy Pattern

## Overview

The Proxy Pattern is a structural design pattern that provides a placeholder or surrogate for another object to control access to it. A proxy acts as an intermediary, intercepting calls to the real object and adding behavior like validation, caching, logging, or lazy initialization.

## Problem

When objects are expensive to create or access, or when you need to add behavior before accessing them, direct access becomes problematic:

- **Expensive resources**: Creating objects on demand wastes resources
- **No access control**: Can't restrict who accesses sensitive objects
- **Lack of logging**: Can't track access patterns or usage
- **No lazy loading**: Heavy objects load even if never used
- **Missing validation**: No way to validate before processing
- **Caching challenges**: No central point to cache results

### Real-World Context

Consider a document management system where loading a high-resolution image from disk is expensive and slow. Instead of loading images immediately, a proxy can defer loading until the image is actually accessed. The proxy acts as a lightweight placeholder, loading the real image only when needed.

## Why Use This Pattern?

- **Lazy initialization**: Defer expensive operations until actually needed
- **Access control**: Restrict access based on permissions
- **Caching**: Store results to avoid repeated expensive computations
- **Logging and monitoring**: Track access and usage patterns
- **Remote objects**: Handle objects on different machines transparently
- **Validation**: Check preconditions before delegating to real object

## When to Use

- Lazy initialization (expensive objects created only when used)
- Access control (restrict who can access objects)
- Caching expensive computations
- Logging, auditing, or monitoring object access
- Remote object access (RPC, distributed systems)
- Rate limiting or resource throttling
- Object pooling and resource management

## When NOT to Use

- Simple objects where proxy overhead isn't justified
- Direct access is clearer than proxy indirection
- Performance is critical and proxy overhead is unacceptable
- No real benefit from proxy behavior (logging not needed, not expensive, etc.)
- Just need a simple wrapper (decorator may be more appropriate)

## Implementation Guidelines

1. **Subject interface**: Define the interface that proxy and real object share
2. **Real subject**: The actual expensive or controlled object
3. **Proxy**: Implements subject interface and holds reference to real subject
4. **Control behavior**: Add logic for lazy loading, access control, caching, etc.
5. **Transparency**: Proxy should be transparent to clients

## Go Idioms

Go's interface system makes proxies elegant:

- Small focused interfaces are easier to proxy
- No inheritance needed; implicit interface satisfaction
- Function types can act as proxies
- Channels and goroutines enable remote proxy patterns

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│       Proxy Pattern Architecture             │
└──────────────────────────────────────────────┘

Subject Interface:
    ┌──────────────────┐
    │ Operation()      │
    └──────────────────┘
           ▲
           │
        ┌──┴──────────────┐
        │                 │
        ▼                 ▼
    ┌─────────┐      ┌──────────────┐
    │  Proxy  │      │ RealSubject  │
    │         │      │              │
    │ Check   │◄────►│ Expensive    │
    │ Log     │      │ Operations   │
    │ Cache   │      │              │
    └─────────┘      └──────────────┘
        ▲
        │
      Client

Proxy Control Flow:

Client calls: proxy.Operation()
    │
    ▼
Proxy intercepts:
├─ Check access? (Access Control Proxy)
├─ Check cache? (Caching Proxy)
├─ Log call? (Logging Proxy)
├─ Check rate limit? (Rate Limiting Proxy)
├─ Load if needed? (Lazy Loading Proxy)
    │
    ├─ YES: Call real object
    │   └─> RealSubject.Operation()
    │       │
    │       ▼
    │   (Expensive work)
    │       │
    │       ▼
    │   Return result
    │
    └─ NO: Return cached result
        │
        ▼
    Return to client


Lazy Loading Image Example:

WITHOUT Proxy:
  Load image immediately
    │
    ├─ Read from disk (slow!)
    ├─ Decode (slow!)
    ├─ Store in memory
    └─ Maybe never used!

WITH Proxy:
  Return proxy immediately
    │
    ▼
  Is image accessed? ────NO────> Keep as proxy
    │
    YES
    │
    ▼
  Load real image (first access only)
    │
    ▼
  Future accesses use loaded image
```

## Real-World Examples

### 1. Lazy Loading Image Proxy

```go
type Image interface {
    Display() string
}

type RealImage struct {
    filename string
    data     []byte
}

type ImageProxy struct {
    filename  string
    realImage *RealImage
}

func (p *ImageProxy) Display() string {
    if p.realImage == nil {
        p.realImage = LoadImageFromDisk(p.filename)
    }
    return p.realImage.Display()
}
```

### 2. Access Control Proxy

```go
type Database interface {
    Query(sql string) Result
}

type DBProxy struct {
    user   string
    db     *RealDatabase
}

func (p *DBProxy) Query(sql string) Result {
    if !p.canAccess(sql) {
        return Error("Access denied")
    }
    return p.db.Query(sql)
}
```

### 3. Caching Proxy

```go
type CachingProxy struct {
    cache       map[string]Result
    realService *ExpensiveService
}

func (p *CachingProxy) GetData(key string) Result {
    if cached, exists := p.cache[key]; exists {
        return cached
    }
    result := p.realService.GetData(key)
    p.cache[key] = result
    return result
}
```

## Key Advantages

- **Lazy initialization**: Create expensive objects only when needed
- **Access control**: Restrict who can access objects
- **Caching**: Avoid repeated expensive operations
- **Logging**: Track all access to objects
- **Transparency**: Proxy looks like real object to clients
- **Resource management**: Control object creation and lifecycle
- **Rate limiting**: Throttle access to expensive resources

## Key Gotchas

- **Extra indirection**: Proxy adds call overhead
- **Complexity**: Too much logic in proxy obscures intent
- **Caching invalidation**: Keeping cached data fresh is hard
- **Type switching**: Clients may need to know about proxy vs. real object
- **Performance**: Proxy overhead may exceed benefits
- **Synchronization**: Thread-safety becomes more complex
- **Debugging difficulty**: Extra layer makes debugging harder
- **Not transparent**: Proxy and real object may not be perfectly equivalent
