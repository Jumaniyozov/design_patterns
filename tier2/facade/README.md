# Facade Pattern

## Overview

The Facade Pattern is a structural design pattern that provides a unified, simplified interface to a set of interfaces in a subsystem. It hides the complexity of the subsystem by providing a single entry point that orchestrates multiple components.

## Problem

As systems grow complex with multiple interconnected subsystems, you face several challenges:

- **Overwhelming complexity**: Clients must understand multiple classes and their interactions
- **Tight coupling**: Clients depend on many internal components
- **Difficult to learn**: New developers struggle to understand how to use the system
- **Fragile code**: Changes in subsystem details break multiple clients
- **Scattered logic**: Related operations spread across different classes
- **Initialization complexity**: Setting up multiple components becomes tedious

### Real-World Context

Consider a banking system with multiple subsystems: Account management, Transaction processing, Fraud detection, and Notification service. Instead of clients dealing with all these subsystems directly, a Facade provides a simple interface: `TransferMoney()` that internally orchestrates all necessary operations.

## Why Use This Pattern?

- **Simplification**: Complex subsystems appear simple to clients
- **Decoupling**: Clients depend on facade, not internal components
- **Organized structure**: Clear separation between interface and implementation
- **Easier integration**: New clients can use complex subsystems without deep knowledge
- **Change insulation**: Internal changes don't affect clients if facade stays stable

## When to Use

- Subsystems have become complex with multiple interconnected components
- Need to provide a simple interface to a complex set of classes
- Want to reduce dependencies between clients and implementation details
- Building layers in a large application (presentation, business, data layers)
- Integrating multiple third-party libraries into a cohesive interface
- Creating APIs that hide internal complexity

## When NOT to Use

- Simple systems that don't need simplification
- Hiding important details that clients should understand
- Facade adds unnecessary indirection without value
- All clients need direct access to subsystem details
- Creating a facade for a single class (unnecessary wrapper)

## Implementation Guidelines

1. **Identify subsystems**: Determine what components the facade will coordinate
2. **Common operations**: Define facade methods representing frequent client needs
3. **Delegation**: Facade delegates to appropriate subsystem components
4. **Hide details**: Keep subsystem classes internal; expose only through facade
5. **Simple interface**: Facade methods should have simple, intuitive signatures

## Go Idioms

Go's minimalist philosophy aligns well with facades:

- Simple interfaces focus on essential operations
- Package-level functions act as facades for package functionality
- Avoid over-abstraction; facades should simplify, not obscure
- Clear separation between public (facade) and unexported internals

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│      Facade Pattern Architecture             │
└──────────────────────────────────────────────┘

Without Facade (Complexity):

     Client
       │
    ┌──┴─────┬──────┬──────┐
    │        │      │      │
    ▼        ▼      ▼      ▼
  Class    Class  Class  Class
    A        B      C      D
    │        │      │      │
    └────────┴──────┴──────┘
     (Client must manage all)


With Facade (Simplified):

     Client
       │
       ▼
    ┌─────────────┐
    │   Facade    │
    │             │
    │ simpleOp()  │
    │ complexOp() │
    └──────┬──────┘
           │
    ┌──────┴────┬──────┬──────┐
    │           │      │      │
    ▼           ▼      ▼      ▼
  Class       Class  Class  Class
    A           B      C      D

Banking System Example:

┌─────────────────────────────────┐
│      BankingFacade              │
│                                 │
│ TransferMoney(from, to, amt)    │
│ OpenAccount(customer)           │
│ CloseAccount(account)           │
└───────────────┬─────────────────┘
                │
    ┌───────────┼────────┬──────────┐
    │           │        │          │
    ▼           ▼        ▼          ▼
┌────────┐┌────────┐ ┌────────┐┌────────┐
│Account ││Trans-  │ │ Fraud  ││Notif.. │
│Service ││action  │ │Detect  ││Service │
│        ││Service | │        │└────────┘
└────────┘└────────┘ └────────┘

Facade workflow:
1. Receive TransferMoney request
2. Validate accounts (Account Service)
3. Check fraud (Fraud Detection)
4. Process transaction (Transaction Service)
5. Notify parties (Notification Service)
6. Return result to client

Client sees one simple operation, not five subsystem calls.
```

## Real-World Examples

### 1. Banking System Facade

```go
type BankingFacade struct {
    accountService  *AccountService
    transactionSvc  *TransactionService
    fraudDetection  *FraudDetectionService
    notification    *NotificationService
}

func (bf *BankingFacade) TransferMoney(from, to string, amount float64) error {
    // Validates, checks fraud, processes, notifies
    // All coordinated transparently
}
```

### 2. Web Framework Facade

```go
type WebFramework struct {
    router    *Router
    database  *Database
    cache     *Cache
    validator *Validator
}

func (wf *WebFramework) HandleRequest(req *Request) *Response {
    // Routes, validates, queries DB, caches results
    // Client just calls one method
}
```

### 3. Reporting System Facade

```go
type ReportingFacade struct {
    dataSource *DataSource
    processor  *DataProcessor
    formatter  *ReportFormatter
    exporter   *Exporter
}

func (rf *ReportingFacade) GenerateReport(params ReportParams) (Report, error) {
    // Fetches data, processes, formats, returns
}
```

## Key Advantages

- **Simplification**: Complex systems appear simple through unified interface
- **Decoupling**: Clients depend on facade, not internal details
- **Change insulation**: Internal refactoring doesn't break clients
- **Easier learning curve**: New users start with simple facade interface
- **Clear boundaries**: Obvious separation between interface and implementation
- **Layering**: Enables creation of architectural layers

## Key Gotchas

- **Oversimplification**: Facade may hide important details clients need
- **Bloated interface**: Too many methods make facade complex (defeats purpose)
- **Unnecessary wrapper**: Using facade for single simple class adds no value
- **Feature loss**: Facade may not expose all subsystem capabilities
- **Performance**: Extra indirection can impact performance-critical code
- **Testing complexity**: Facades with many dependencies are harder to test
- **Outdated facade**: If subsystem changes, facade must be updated
- **Wrong level of abstraction**: Facade placed at wrong architectural level
