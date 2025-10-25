# Service Locator Pattern

## Overview
The Service Locator pattern provides a central registry to locate services/dependencies. It acts as a global access point for obtaining service instances.

## Problem
Hard-coding service instantiation creates tight coupling. Passing dependencies everywhere is cumbersome. You need centralized service management and lookup.

## Why Use This Pattern?
- **Centralized Registry**: One place to manage services
- **Decoupling**: Clients don't create dependencies
- **Runtime Configuration**: Swap implementations dynamically
- **Convenience**: Easy access to services

## When to Use
- Plugin architectures
- Service-oriented applications
- Legacy code integration
- Framework-level service management

## When NOT to Use
- **Prefer Dependency Injection**: Service Locator is considered an anti-pattern by many
- Makes dependencies implicit (harder to understand)
- Testing becomes more difficult
- Hidden dependencies violate explicit is better than implicit

## Note
While useful in some scenarios, modern Go applications typically prefer explicit dependency injection over Service Locator for better testability and clarity.

## Go Idioms
```go
type ServiceLocator struct {
    services map[string]interface{}
}

func (sl *ServiceLocator) Register(name string, service interface{})
func (sl *ServiceLocator) Get(name string) (interface{}, error)
```
