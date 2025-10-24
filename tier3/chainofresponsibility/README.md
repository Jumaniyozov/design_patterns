# Chain of Responsibility Pattern

## Overview

The Chain of Responsibility Pattern is a behavioral design pattern that lets you pass requests along a chain of handlers. Each handler decides whether to process the request or pass it to the next handler in the chain. This allows multiple objects to handle a request without coupling the sender to any specific handler.

## Problem

When handling requests with multiple potential processors or when multiple validators need to run sequentially, direct routing becomes problematic:

- **Hard to determine handler**: No clear way to route requests to the right handler
- **Tight coupling**: Sender must know about all potential handlers
- **Fixed order issues**: Handler order may need to change dynamically
- **Complex conditional logic**: Long if-else chains to determine handler
- **Difficult to extend**: Adding new handlers requires modifying existing code
- **Multiple validations**: Running multiple validators in sequence is messy

### Real-World Context

Consider an HTTP request processing pipeline where a request must pass through authentication, authorization, validation, and business logic handlers in sequence. Each handler can either process the request or pass it to the next one. If authentication fails, the chain stops; otherwise it continues.

## Why Use This Pattern?

- **Decoupling**: Sender doesn't know which handler will process the request
- **Dynamic chains**: Handler order can be configured at runtime
- **Easy to extend**: Add new handlers without modifying existing ones
- **Flexible processing**: Multiple handlers can process the same request
- **Middleware pattern**: Natural fit for request processing pipelines
- **Separation of concerns**: Each handler has one responsibility

## When to Use

- Request handling pipelines (HTTP middleware, event processing)
- Validation chains where multiple validators run in sequence
- Approval workflows (document approval, expense requests)
- Error handling strategies that try different approaches
- Logging and monitoring with multiple handlers
- Event handling systems with multiple subscribers

## When NOT to Use

- Single handler can process request (no need for chain)
- Order of handlers doesn't matter (use observer instead)
- All handlers must process (not a chain, but broadcast)
- Simple if-else logic (chains add unnecessary complexity)
- Performance critical where chain overhead matters

## Implementation Guidelines

1. **Handler interface**: Define how handlers are called
2. **Concrete handlers**: Implement handler interface and maintain reference to next
3. **Chain construction**: Build the chain by linking handlers
4. **Request processing**: Each handler decides to process or pass
5. **Next handler**: Each handler knows how to call the next one

## Go Idioms

Go's function types and interfaces make chains elegant:

- Function types can be chained using closures
- Middleware pattern is idiomatic for HTTP
- Interfaces focus on simple operations
- Composition naturally builds chains

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│   Chain of Responsibility Pattern            │
└──────────────────────────────────────────────┘

Handler Interface:
  ┌────────────────────────┐
  │ Handler                │
  │ Handle(request)        │
  │ SetNext(next Handler)  │
  └────────────────────────┘
           ▲
           │
    ┌──────┴──────────────────┬──────────────┐
    │                         │              │
    ▼                         ▼              ▼
┌─────────────┐      ┌──────────────┐  ┌────────┐
│Handler A    │      │Handler B     │  │Handler │
│             │      │              │  │   C    │
│ Handle(){   │      │Handle(){     │  │ {      │
│  process()  │      │ process()    │  │ final  │
│  next()     │      │ next()       │  │ handle │
│ }           │      │ }            │  │ }      │
└──────┬──────┘      └──────┬───────┘  └────────┘
       │                    │
       └────next ─────► Handle B
                           │
                           └────next ─────► Handler C


HTTP Middleware Chain Example:

┌──────────┐
│ Request  │
└─────┬────┘
      │
      ▼
┌───────────────┐
│Authentication │ Check credentials
│ Middleware    │ Pass or block
└─────┬─────────┘
      │
      ▼
┌───────────────┐
│Authorization  │ Check permissions
│ Middleware    │ Pass or block
└─────┬─────────┘
      │
      ▼
┌───────────────┐
│Validation     │ Validate data
│ Middleware    │ Pass or reject
└─────┬─────────┘
      │
      ▼
┌───────────────┐
│Business Logic │ Process request
│ Handler       │ Return result
└─────┬─────────┘
      │
      ▼
┌──────────┐
│ Response │
└──────────┘

Order matters! Each handler receives
result of previous handler, can intercept
or modify before passing forward.

Building the Chain:

auth := NewAuthHandler()
authz := NewAuthorizationHandler()
validate := NewValidationHandler()
business := NewBusinessHandler()

auth.SetNext(authz)
authz.SetNext(validate)
validate.SetNext(business)

// Start processing
request.Handle(auth)
```

## Real-World Examples

### 1. HTTP Middleware Chain

```go
type Handler interface {
    Handle(w http.ResponseWriter, r *http.Request)
    SetNext(next Handler) Handler
}

type AuthMiddleware struct { /*...*/ }
type LoggingMiddleware struct { /*...*/ }
type ValidationMiddleware struct { /*...*/ }

// Chain them together
auth.SetNext(logging)
logging.SetNext(validation)
```

### 2. Approval Workflow

```go
type ApprovalHandler interface {
    Approve(request ApprovalRequest) bool
    SetNext(next ApprovalHandler) ApprovalHandler
}

type ManagerApproval struct { /*...*/ }
type DirectorApproval struct { /*...*/ }
type CFOApproval struct { /*...*/ }

// Chain based on amount
if amount < 1000 { handler = ManagerApproval }
else if amount < 10000 { handler = DirectorApproval }
else { handler = CFOApproval }
```

### 3. Validation Chain

```go
type Validator interface {
    Validate(data Data) error
    SetNext(next Validator) Validator
}

type EmailValidator struct { /*...*/ }
type LengthValidator struct { /*...*/ }
type FormatValidator struct { /*...*/ }

// All validators run in sequence
email.SetNext(length)
length.SetNext(format)
```

## Key Advantages

- **Loose coupling**: Sender doesn't know about handlers
- **Dynamic chains**: Change handler order at runtime
- **Easy extension**: Add new handlers without changing existing code
- **Separation of concerns**: Each handler has one responsibility
- **Flexible processing**: Can process, modify, or reject requests
- **Middleware pattern**: Natural fit for request processing

## Key Gotchas

- **Handler order matters**: Different order produces different results
- **Complexity**: Chains can become hard to debug and understand
- **Performance**: Each handler adds processing overhead
- **Silent failures**: Request can be lost if no handler processes it
- **Circular chains**: Must prevent cycles in chain construction
- **Thread safety**: Chain itself must be thread-safe if shared
- **Debugging difficulty**: Tracing through chain is harder than direct calls
- **Request modification**: Handlers may unintentionally affect other handlers
