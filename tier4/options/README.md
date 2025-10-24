# Options Pattern (Functional Options)

## Overview

The Options Pattern, also known as Functional Options, is a Go-idiomatic design pattern that provides a clean and flexible way to handle optional parameters and configuration. It uses variadic function parameters and closures to set options on a struct without needing method overloading or complex constructors.

## Problem

When creating objects with many optional parameters, traditional approaches become problematic:

- **No method overloading**: Go doesn't support overloading, making multiple constructors messy
- **Positional parameters**: Hard to remember parameter order
- **Nil parameters**: Ugly to pass many nil or zero values
- **Future extensibility**: Adding new optional parameters breaks existing code
- **Configuration clarity**: Not obvious which parameters are optional
- **Builder verbosity**: Builder pattern adds boilerplate for simple cases

### Real-World Context

Creating an HTTP server with optional configuration: timeout, TLS, middleware, logging, etc. Traditional approaches require either a massive constructor or a builder. Options pattern provides a clean, extensible way to specify only needed options.

## Why Use This Pattern?

- **Clean API**: Clear, readable configuration at call site
- **Extensible**: Add new options without breaking existing code
- **Go idiomatic**: Widely used in Go community and standard library
- **Type-safe**: Compile-time checking of option validity
- **Performant**: No reflection or runtime overhead
- **Simple implementation**: Elegant and straightforward

## When to Use

- Functions with many optional parameters
- Creating configurable objects (servers, clients, pools)
- Library APIs that need flexibility
- Deprecating parameters without breaking changes
- Complex initialization logic with variations
- Any public API with optional settings

## When NOT to Use

- Simple objects with few required parameters
- All parameters are required (just use function params)
- Simplicity is more important than extensibility
- Performance critical where function calls matter

## Implementation Guidelines

1. **Option type**: Define a type for option functions
2. **Getter function**: Main constructor accepts options
3. **Option functions**: Each returns a function that modifies the struct
4. **Variadic parameter**: Constructor uses ...Option
5. **Apply options**: Loop through options and apply each

## Go Idioms

The Options pattern is idiomatic Go, used extensively in the standard library:

```go
// From flag package
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet

// More modern Go style with options:
func NewServer(opts ...Option) *Server

// Used with option functions:
server := NewServer(
    WithPort(8080),
    WithTimeout(10*time.Second),
)
```

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│    Options Pattern Architecture              │
└──────────────────────────────────────────────┘

Pattern Structure:

type Option func(*Server)

Define options:
func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

Constructor:
func NewServer(opts ...Option) *Server {
    s := &Server{
        // defaults
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

Usage:
server := NewServer(
    WithPort(8080),
    WithTimeout(10*time.Second),
)

Comparison:

BEFORE (Without Options):
┌──────────────────────────────────────────┐
│ Ugly: Too many parameters                │
│ s := NewServer(8080, true, 10, 100, nil) │
│                ↑    ↑   ↑   ↑    ↑       │
│            What are these?!?             │
└──────────────────────────────────────────┘

AFTER (With Options):
┌──────────────────────────────────────────┐
│ Clear: Intent obvious from parameter     │
│ s := NewServer(                          │
│     WithPort(8080),                      │
│     WithSSL(true),                       │
│     WithTimeout(10*time.Second),         │
│     WithMaxConnections(100),             │
│ )                                        │
└──────────────────────────────────────────┘

Option Function Flow:

Client Code
    │
    ▼
NewServer(WithPort(8080), WithTimeout(5s), ...)
    │
    ├─ Create default server struct
    │
    ├─ For each option:
    │  ├─ Call WithPort(8080)
    │  │  └─ Returns func(s *Server) { s.port = 8080 }
    │  │
    │  ├─ Apply function to server
    │  │  └─ Executes the closure
    │  │
    │  ├─ Call WithTimeout(5s)
    │  │  └─ Returns func(s *Server) { s.timeout = 5s }
    │  │
    │  └─ Apply function to server
    │
    ▼
Return fully configured server

Stack of Modifiers:

Default Server
    │
    ├─ WithPort(8080) ────> Server.port = 8080
    │
    ├─ WithTimeout(5s) ───> Server.timeout = 5s
    │
    ├─ WithSSL(true) ─────> Server.ssl = true
    │
    ├─ WithLogger(log) ───> Server.logger = log
    │
    ▼
Fully Configured Server
```

## Real-World Examples

### 1. HTTP Server with Options

```go
type Server struct {
    port         int
    timeout      time.Duration
    ssl          bool
    maxConnections int
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func WithSSL(useSSL bool) Option {
    return func(s *Server) {
        s.ssl = useSSL
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        port:    8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(
    WithPort(9000),
    WithTimeout(60 * time.Second),
    WithSSL(true),
)
```

### 2. Database Connection Pool

```go
type ConnectionPool struct {
    host         string
    port         int
    maxConns     int
    idleTimeout  time.Duration
    retryPolicy  RetryPolicy
}

type Option func(*ConnectionPool)

func WithMaxConnections(max int) Option {
    return func(cp *ConnectionPool) {
        cp.maxConns = max
    }
}

func WithIdleTimeout(timeout time.Duration) Option {
    return func(cp *ConnectionPool) {
        cp.idleTimeout = timeout
    }
}

func NewConnectionPool(host string, port int, opts ...Option) *ConnectionPool {
    pool := &ConnectionPool{
        host:        host,
        port:        port,
        maxConns:    10,
        idleTimeout: 5 * time.Minute,
    }
    for _, opt := range opts {
        opt(pool)
    }
    return pool
}
```

### 3. Logger with Options

```go
type Logger struct {
    level  LogLevel
    format LogFormat
    output io.Writer
}

type Option func(*Logger)

func WithLevel(level LogLevel) Option {
    return func(l *Logger) { l.level = level }
}

func WithFormat(format LogFormat) Option {
    return func(l *Logger) { l.format = format }
}

func WithOutput(writer io.Writer) Option {
    return func(l *Logger) { l.output = writer }
}

func NewLogger(opts ...Option) *Logger {
    l := &Logger{
        level:  InfoLevel,
        format: JSONFormat,
        output: os.Stderr,
    }
    for _, opt := range opts {
        opt(l)
    }
    return l
}
```

## Key Advantages

- **Clean API**: Clear, readable configuration at call site
- **Extensible**: Add new options without breaking existing code
- **Type-safe**: Compile-time checking of options
- **Go idiomatic**: Matches Go community conventions
- **Performant**: No reflection, just function calls
- **Flexible**: Only specify needed options
- **Simple defaults**: Sensible defaults with optional overrides
- **No method overloading**: Works around Go limitation naturally

## Key Gotchas

- **Order matters**: Some options may conflict if specified in certain order
- **Complexity hiding**: Implementation details hidden in closures
- **Harder to debug**: Stack traces show internal option functions
- **Learning curve**: New Go developers may not be familiar with pattern
- **No validation**: Options should validate their inputs
- **Option conflicts**: No automatic detection of incompatible options
- **Documentation**: Must clearly document all available options
- **Testing**: Need to test all option combinations
