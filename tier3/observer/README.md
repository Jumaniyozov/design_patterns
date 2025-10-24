# Observer Pattern

## Overview

The Observer Pattern is a behavioral design pattern that defines a one-to-many dependency between objects such that when one object changes state, all its dependents are notified automatically. It creates a publish-subscribe relationship where subjects (publishers) notify observers (subscribers) of state changes.

## Problem

When multiple objects need to react to state changes in another object, direct coupling becomes problematic:

- **Tight coupling**: Observers hard-code dependencies on specific subjects
- **Unpredictable updates**: Unclear which objects depend on a subject
- **Cascading logic**: Changes trigger modifications in multiple places
- **Hard to extend**: Adding new observers requires modifying existing code
- **Scattered dependencies**: Observer relationships scattered throughout codebase
- **Difficult debugging**: Tracing cause-and-effect relationships is hard

### Real-World Context

Consider a stock market system where multiple traders want to be notified when stock prices change. Each trader (observer) registers with the stock price tracker (subject). When the price changes, the tracker automatically notifies all registered traders. Traders don't need to know about each other; they independently react to the price changes.

## Why Use This Pattern?

- **Loose coupling**: Subjects and observers depend on abstractions
- **Dynamic relationships**: Add/remove observers at runtime
- **Automatic notification**: No need to manually notify dependents
- **Scalability**: Easy to add new observers without changing subject
- **Channels and goroutines**: Natural fit for Go's concurrent patterns
- **Event-driven architecture**: Foundation for reactive systems

## When to Use

- Multiple objects need to react to state changes
- Number of observers is unknown or changes dynamically
- Event-driven systems and reactive programming
- Model-View architectures where views observe model changes
- Pub-Sub systems and message brokers
- Real-time notifications and updates
- Decoupling producers from consumers

## When NOT to Use

- Single observer (direct method call is simpler)
- Observer order matters critically (use chain of responsibility)
- Performance critical code where observer overhead matters
- Fixed, known set of observers
- Synchronous ordering is essential

## Implementation Guidelines

1. **Subject interface**: Methods to attach/detach/notify observers
2. **Concrete subject**: Maintains state and notifies observers of changes
3. **Observer interface**: Define the notification method
4. **Concrete observers**: Implement observer interface and react to notifications
5. **Notifications**: Subject calls observer methods when state changes

## Go Idioms

Go's channels and goroutines make observers particularly elegant:

- Channels can implement observer pattern naturally
- Goroutines enable asynchronous notifications
- Function types can be observers (simple callbacks)
- No need for complex observer interfaces

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│        Observer Pattern Architecture         │
└──────────────────────────────────────────────┘

Subject Interface:
  ┌──────────────────────────────────┐
  │ Subject                          │
  │ Attach(observer Observer)        │
  │ Detach(observer Observer)        │
  │ Notify()                         │
  │ GetState() State                 │
  └──────────────────────────────────┘
           ▲
           │
           ▼
  ┌──────────────────────────────────┐
  │ ConcreteSubject                  │
  │ ┌───────────────────────────────┐│
  │ │ state     State               ││
  │ │ observers []Observer          ││
  │ │                               ││
  │ │ SetState(s State)             ││
  │ │   state = s                   ││
  │ │   Notify()                    ││
  │ └───────────────────────────────┘│
  └──────────────────────────────────┘

Observer Interface:
  ┌──────────────────────────┐
  │ Observer                 │
  │ Update(subject Subject)  │
  └──────────────────────────┘
           ▲
           │
    ┌──────┴──────────┬──────────────┐
    │                 │              │
    ▼                 ▼              ▼
┌────────┐      ┌────────┐      ┌────────┐
│Obs1    │      │Obs2    │      │Obs3    │
│        │      │        │      │        │
│Update()│      │Update()│      │Update()|
└────────┘      └────────┘      └────────┘


Event Flow (Stock Price Example):

Stock Price (Subject):
  ┌──────────────────────────┐
  │ Price = $100             │
  │ Attached Observers:      │
  │ [Trader1, Trader2, ... ] |
  └────┬─────────────────────┘
       │
       └─ Price changes to $105
       │
       └─ Notify()
           │
        ┌──┴────────┬──────────┬────────┐
        │           │          │        │
        ▼           ▼          ▼        ▼
       Trader1   Trader2    Analyst   App
        │           │          │        │
        ├─ React: Update  →  React →  React display
        │  portfolio       │ reports    │ 
        │                  └────────────┘


Pub-Sub with Channels:

 Publisher               Channel               Subscribers
     │                      │                        │
     ├─ price change ──────>│                        │
     │                      ├─────> Trader1 (listen) |
     │                      │                        │
     │                      ├─────> Trader2 (listen) |
     │                      │                        │
     │                      ├─────> Analyst (listen) |
     │                      │
     └─ more changes ───────>
                            └─> All get notified
```

## Real-World Examples

### 1. Stock Price Observer

```go
type StockPrice interface {
    GetValue() float64
}

type Observer interface {
    Update(subject StockPrice)
}

type Stock struct {
    price     float64
    observers []Observer
}

func (s *Stock) SetPrice(p float64) {
    s.price = p
    s.NotifyObservers()
}

type Trader struct { /*...*/ }
func (t *Trader) Update(subject StockPrice) {
    // React to price change
}
```

### 2. Event Bus with Channels

```go
type EventBus struct {
    subscribers map[string][]chan interface{}
}

func (eb *EventBus) Subscribe(event string, handler chan interface{}) {
    eb.subscribers[event] = append(eb.subscribers[event], handler)
}

func (eb *EventBus) Publish(event string, data interface{}) {
    for _, handler := range eb.subscribers[event] {
        handler <- data
    }
}
```

### 3. UI Component Observer (MVC)

```go
type Model struct {
    data      string
    observers []View
}

type View interface {
    Update(model *Model)
}

func (m *Model) SetData(d string) {
    m.data = d
    for _, view := range m.observers {
        view.Update(m)
    }
}
```

## Key Advantages

- **Loose coupling**: Subjects and observers depend on abstractions
- **Dynamic relationships**: Add/remove observers at runtime
- **Automatic propagation**: No manual notification code needed
- **Scalability**: Easy to add new observers without changing subject
- **Separation of concerns**: Each observer handles its own logic
- **Reactive programming**: Foundation for event-driven systems
- **Channels and goroutines**: Natural fit for Go's concurrency model

## Key Gotchas

- **Memory leaks**: Forgotten observer references prevent garbage collection
- **Observer order**: Undefined execution order of observers
- **Circular updates**: Observer modifications trigger more updates
- **Performance**: Many observers can slow down updates
- **Complexity**: Can become hard to trace update chains
- **Synchronous overhead**: All observers notified before SetState returns
- **Testing difficulty**: Multiple observers make testing complex
- **Exception handling**: One observer exception shouldn't affect others
