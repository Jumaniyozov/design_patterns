# Saga Pattern

## Overview
The Saga pattern manages distributed transactions across multiple services by breaking them into a series of local transactions, each with a compensating transaction to undo changes if the saga fails.

## Problem
Distributed transactions with 2PC (two-phase commit) don't scale and create tight coupling. You need a way to maintain data consistency across services without distributed locks.

## Why Use This Pattern?
- **Distributed Consistency**: Maintain consistency without 2PC
- **Scalability**: No distributed locks
- **Fault Tolerance**: Handle partial failures gracefully
- **Loose Coupling**: Services remain independent

## When to Use
- Microservices architecture
- Long-running business processes
- Need distributed transactions
- Services have independent databases

## Real-world scenarios
- Order processing (reserve inventory, charge payment, create order, ship)
- Travel booking (book flight, hotel, car; cancel all if any fails)
- E-commerce checkout (payment, inventory, shipping)

## Types
**Choreography**: Services listen to events and react
**Orchestration**: Central coordinator manages saga flow

## Go Idioms
```go
type Step interface {
    Execute(ctx context.Context) error
    Compensate(ctx context.Context) error
}

type Saga struct {
    steps []Step
}

func (s *Saga) Execute() error {
    for i, step := range s.steps {
        if err := step.Execute(); err != nil {
            // Compensate previous steps
            s.compensate(i)
            return err
        }
    }
}
```
