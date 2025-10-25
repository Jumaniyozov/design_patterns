# Active Object Pattern

## Overview
The Active Object pattern decouples method execution from method invocation by running methods asynchronously in their own thread/goroutine. Clients get immediate control back with a future/promise for the result.

## Problem
Blocking operations hurt responsiveness. You need to execute long-running operations asynchronously while maintaining a clean, synchronous-looking API.

## Why Use This Pattern?
- **Asynchronous Execution**: Non-blocking method calls
- **Concurrency**: Methods execute in separate goroutine
- **Simplicity**: Clients see simple synchronous API
- **Thread Safety**: Internal synchronization handled

## When to Use
- Long-running operations
- Background processing
- Non-blocking APIs
- Actor-like patterns

## Real-world scenarios
- Async I/O operations
- Background data processing
- Message queue consumers
- Event handlers

## Go Idioms
Use goroutines and channels:
```go
type ActiveObject struct {
    requests chan func()
}

func (a *ActiveObject) DoWork(data string) <-chan Result {
    result := make(chan Result, 1)
    a.requests <- func() {
        result <- doActualWork(data)
    }
    return result
}
```
