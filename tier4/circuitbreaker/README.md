# Circuit Breaker Pattern

## Overview

The Circuit Breaker Pattern is a behavioral design pattern used to prevent an application from making requests to a service that is failing or unresponsive. It acts like an electrical circuit breaker, stopping requests when the service is down and allowing recovery time before retrying.

## Problem

When calling external services, failures can cascade and create problems:

- **Cascading failures**: Calling a failing service repeatedly wastes resources and slows response
- **Resource exhaustion**: Too many stuck requests consume connection pools and memory
- **No recovery time**: Service never gets a chance to recover from temporary issues
- **Poor user experience**: Users wait for timeouts on failed requests
- **No feedback**: System doesn't know when to start retrying a failed service
- **Silent failures**: Failures silently compound across the system

### Real-World Context

Imagine a microservices system where Service A calls Service B, which calls Service C. If Service C fails, Service B keeps retrying, consuming resources. Service A then waits for Service B, and the failure cascades upward. A circuit breaker at each service boundary prevents this cascade by immediately failing requests to a downed service.

## Why Use This Pattern?

- **Failure isolation**: Prevent cascading failures across services
- **Fast fail**: Stop trying a downed service immediately
- **Resource protection**: Prevent connection pool and thread exhaustion
- **Recovery time**: Give failing services time to recover
- **Graceful degradation**: Show meaningful errors instead of timeouts
- **Essential for microservices**: Required for resilient distributed systems

## When to Use

- Calling external services (HTTP, databases, APIs)
- Microservices communication
- Protecting against cascading failures
- Services with occasional failures or slow response
- Need for resilience and automatic recovery
- Distributed systems requiring fault tolerance
- Any network-dependent operation

## When NOT to Use

- Local function calls (no failure mode)
- Operations that must always complete (no fallback possible)
- Circuit breaker overhead exceeds benefit
- Simple retry is sufficient
- Service is always reliable (no failures)

## Implementation Guidelines

1. **States**: Track Closed (normal), Open (failing), Half-Open (testing)
2. **Threshold**: Count consecutive failures before opening circuit
3. **Timeout**: Wait period before attempting recovery
4. **Reset logic**: How to transition between states
5. **Monitoring**: Track circuit state and metrics

## Go Idioms

Go's simplicity and concurrency features fit circuit breaker naturally:

- Goroutines handle concurrent requests
- Channels can signal state changes
- sync.Mutex protects shared state
- Interfaces define service contracts

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│     Circuit Breaker Pattern States           │
└──────────────────────────────────────────────┘

Three States:

CLOSED (Normal Operation):
  ┌─────────────────────┐
  │ CLOSED              │ Request flows normally
  │ (Circuit Normal)    │ Failures counted
  │                     │
  │ Success: stay       │
  │ Failures: count     │
  │ failCount++         │
  │                     │
  │ If failCount >      │
  │ threshold ──────┐   │
  └─────────────────┼───┘
                    │
                    ▼
               ┌─────────────┐
               │ OPEN        │ Stop all requests
               │ (Circuit    │ Service failing
               │  OPEN)      │
               │             │
               │ Block all   │ Fail immediately
               │ requests    │
               │             │
               │ Wait time   │
               │ passes ──┐  │
               └──────────┼──┘
                          │
                          ▼
                    ┌──────────────┐
                    │ HALF-OPEN    │ Try request
                    │ (Testing)    │ Testing recovery
                    │              │
                    │ Success:     │
                    │ Reset ────┐  │
                    │           │  │
                    │ Failure:  │  │
                    │ Open ──┐  │  │
                    └────┬───┼──┘  │
                         │   │     │
                         │   └─────┼──> CLOSED
                         │         │
                         └─────────┼──> OPEN

Request Flow:

Client Request
    │
    ▼
Check Circuit State
    │
    ├─ CLOSED: Send request
    │          │
    │          ├─ Success: Return
    │          │
    │          └─ Failure: Increment counter
    │                     │
    │                     ├─ Count < threshold: Continue
    │                     │
    │                     └─ Count >= threshold: Open circuit
    │
    ├─ OPEN: Return error immediately (fast fail!)
    │        No request sent to service
    │        Wait for timeout
    │
    └─ HALF-OPEN: Send 1 test request
                   │
                   ├─ Success: Close circuit, resume normal
                   │
                   └─ Failure: Open circuit, wait longer

Metrics and Monitoring:

┌────────────────────────────────┐
│ Circuit Breaker State          │
│ Service: PaymentAPI            │
│                                │
│ Current State: OPEN            │
│ Failures: 15/10                │
│ Success Rate: 42%              │
│ Last failure: 2s ago           │
│ Time until retry: 8s           │
│ Total requests: 1000           │
│ Successful: 420                │
│ Failed: 580                    │
└────────────────────────────────┘
```

## Real-World Examples

### 1. HTTP Service Circuit Breaker

```go
type CircuitBreaker struct {
    state       string // CLOSED, OPEN, HALF_OPEN
    failures    int
    threshold   int
    timeout     time.Duration
    lastFailure time.Time
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == "OPEN" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "HALF_OPEN"
        } else {
            return errors.New("circuit open")
        }
    }

    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.threshold {
            cb.state = "OPEN"
        }
        return err
    }

    cb.failures = 0
    cb.state = "CLOSED"
    return nil
}
```

### 2. Database Connection Circuit Breaker

```go
type DBCircuitBreaker struct {
    db          *sql.DB
    breaker     *CircuitBreaker
}

func (db *DBCircuitBreaker) Query(query string) (*sql.Rows, error) {
    return db.breaker.Call(func() error {
        rows, err := db.db.Query(query)
        return err
    })
}
```

### 3. API Client with Circuit Breaker

```go
type APIClient struct {
    httpClient  *http.Client
    breaker     *CircuitBreaker
}

func (client *APIClient) GetUser(id string) (*User, error) {
    var user User
    err := client.breaker.Call(func() error {
        resp, err := client.httpClient.Get("/users/" + id)
        // ... parse response
        return err
    })
    return &user, err
}
```

## Key Advantages

- **Prevents cascading failures**: Stops propagation of service failures
- **Fast failure**: Returns error immediately when circuit open
- **Resource protection**: Prevents connection pool exhaustion
- **Graceful degradation**: Can provide fallback responses
- **Automatic recovery**: Attempts to recover when service recovers
- **Monitoring**: Clear view of service health
- **User experience**: Users get immediate feedback instead of hanging

## Key Gotchas

- **State synchronization**: Distributed systems need shared circuit breaker state
- **Timeout tuning**: Too short reopens too often, too long prevents recovery
- **Threshold tuning**: Too low triggers too often, too high delays response
- **Partial failures**: Some requests may succeed even when circuit should be open
- **Fallback handling**: Need meaningful error messages or fallback responses
- **Half-open storms**: Multiple clients in half-open can overwhelm recovering service
- **Metrics needed**: Must monitor circuit state to diagnose issues
- **Configuration**: Different services may need different thresholds/timeouts
