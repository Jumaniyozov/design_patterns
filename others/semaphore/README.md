# Semaphore Pattern

## Overview
The Semaphore pattern limits concurrent access to a resource by maintaining a count of available permits. It's essential for rate limiting and resource management in concurrent systems.

## Problem
You need to limit the number of concurrent operations—database connections, API calls, goroutines. Without limits, resource exhaustion and cascading failures occur.

## Why Use This Pattern?
- **Resource Control**: Limit concurrent access
- **Rate Limiting**: Control request rates
- **Backpressure**: Prevent system overload
- **Graceful Degradation**: Manage peak loads

## When to Use
- Limiting concurrent operations
- Rate limiting API calls
- Connection pool management
- Throttling goroutines

## Real-world scenarios
- HTTP client connection limits
- Database connection pooling
- Worker goroutine limits
- API rate limiting

## Go Idioms
Use buffered channels as semaphores:
```go
sem := make(chan struct{}, maxConcurrent)

sem <- struct{}{} // Acquire
defer func() { <-sem }() // Release
```

Or use `golang.org/x/sync/semaphore` package.
