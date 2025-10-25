# Intercepting Filter Pattern

## Overview
The Intercepting Filter pattern preprocesses and postprocesses requests through a chain of filters. Each filter can transform, validate, or log requests before they reach the target handler.

## Problem
Cross-cutting concerns (authentication, logging, compression) shouldn't clutter business logic. You need a way to handle these concerns uniformly across all requests.

## Why Use This Pattern?
- **Separation of Concerns**: Extract cross-cutting logic
- **Flexibility**: Add/remove filters dynamically
- **Reusability**: Filters work across multiple handlers
- **Clean Code**: Business logic stays focused

## When to Use
- HTTP middleware
- Request preprocessing
- Response postprocessing
- Cross-cutting concerns

## Real-world scenarios
- Authentication/Authorization
- Logging and monitoring
- Compression/Decompression
- Rate limiting

## Go Idioms
Go's HTTP middleware pattern is Intercepting Filter:
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```
