# Anti-Corruption Layer Pattern

## Overview
The Anti-Corruption Layer (ACL) acts as a translation layer between two subsystems that don't share the same semantics. It prevents concepts from one bounded context from leaking into another, maintaining clean domain models.

## Problem
When integrating with external systems or legacy code, their models and concepts can "corrupt" your clean domain model. Direct integration creates tight coupling and forces your domain to conform to external constraints.

## Why Use This Pattern?
- **Domain Purity**: Keep your domain model clean
- **Isolation**: Protect from external changes
- **Translation**: Convert between different models
- **Independence**: Evolve systems independently
- **Testing**: Mock external dependencies easily

## When to Use
- Integrating with legacy systems
- Connecting different bounded contexts
- Third-party API integration
- Protecting domain model from external influence

## Real-world scenarios
- Modern system integrating with legacy mainframe
- Microservice calling external payment gateway
- Domain-driven design bounded context integration
- Wrapping poorly designed third-party APIs

## Components
1. **Facade**: Simplifies access to external system
2. **Adapter**: Translates between models
3. **Translator**: Converts data structures
4. **Service**: Encapsulates external operations

## Go Idioms
```go
// External system has different model
type ExternalOrder struct {
    OrderNum string
    CustID   int
    Items    string // comma-separated
}

// Your clean domain model
type Order struct {
    ID       string
    Customer Customer
    Items    []LineItem
}

// ACL translates between them
type OrderACL struct {
    external ExternalOrderService
}

func (acl *OrderACL) GetOrder(id string) (*Order, error) {
    extOrder := acl.external.GetOrder(id)
    return acl.translate(extOrder), nil
}
```
