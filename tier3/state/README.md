# State Pattern

## Overview

The State Pattern is a behavioral design pattern that allows an object to alter its behavior when its internal state changes. The object appears to change its class when its behavior changes based on its state.

## Problem

When objects have behavior that varies dramatically based on their internal state, code becomes difficult to manage:

- **Large conditional statements**: Extensive if-else or switch statements for state handling
- **State-specific logic scattered**: Each state's logic spreads across multiple methods
- **State transitions hard to trace**: Unclear what states exist and how to transition
- **Violation of Single Responsibility**: Methods handle multiple, unrelated states
- **Error-prone transitions**: Missing or invalid state transitions possible
- **Difficult to extend**: Adding new states requires modifying many methods

### Real-World Context

Consider a TCP connection with states: Established, Closed, Listen, etc. Each state has different behavior for sending/receiving data. With the State Pattern, each state is a separate class, making the code cleaner and state transitions explicit.

## Why Use This Pattern?

- **Clarity**: Each state encapsulated in its own class
- **Easy to extend**: Add new states without modifying existing ones
- **Explicit transitions**: Clear, documented state transitions
- **Single responsibility**: Each state handles its own logic
- **Eliminates conditionals**: No more massive if-else chains
- **Testability**: Each state can be tested independently

## When to Use

- Objects with behavior that varies significantly by state
- Objects with many state-dependent conditional statements
- Complex state machines with explicit transitions
- Workflow or lifecycle management
- Game character behaviors (idle, walking, attacking, dead)
- Connection/session management
- Protocol implementations with defined states

## When NOT to Use

- Simple objects with few states
- States don't significantly affect behavior
- Straightforward if-else logic is clearer
- State transitions are rare
- Simplicity is more important than extensibility

## Implementation Guidelines

1. **State interface**: Defines operations available in any state
2. **Concrete states**: Implement interface for each specific state
3. **Context**: Maintains reference to current state and delegates operations
4. **State transitions**: Determine when and how to change states
5. **Behavior encapsulation**: Each state knows its valid transitions

## Go Idioms

Go's interface system makes state machines elegant:

- State interface can be simple (focused operations)
- Implicit satisfaction removes boilerplate
- Function types can represent simple state behaviors
- Composition over inheritance fits naturally

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│         State Pattern Architecture           │
└──────────────────────────────────────────────┘

Context with State:
  ┌──────────────────────────────┐
  │ Context (e.g., Connection)   │
  │                              │
  │ currentState: State          │
  │                              │
  │ Open() ─────┐                │
  │ Close() ────┼─>              │
  │ Send() ─────┤  delegate to   │
  │ Receive() ──┤  currentState  │
  │             │                │
  └──────────────────────────────┘
           │
           │ references
           │
    ┌──────┴──────────────────────┐
    │                             │
    ▼                             ▼
┌─────────┐                  ┌─────────┐
│ Closed  │                  │Established
│  State  │                  │  State
│         │                  │
│Open() ──┼──────┐           │Send()  ──┼──> Send data
│Close()  │      │           │Receive() ┼──> Read data
│Send() ──┼──► Transition    │Close() ──┼──┐
│         │  to Connected    │          │  │
└─────────┘                  └──────────┘  │
                                           │
                              ┌────────────┘
                              │
                              ▼
                          ┌──────────┐
                          │ Closed   │
                          │ State    │
                          └──────────┘

TCP Connection State Machine:

             ┌──────────┐
             │  CLOSED  │
             └────┬─────┘
                  │ listen()
                  ▼
             ┌──────────┐
             │  LISTEN  │
             └────┬─────┘
                  │ accept()
                  ▼
          ┌───────────────┐
          │  ESTABLISHED  │◄──────┐
          └───────┬───────┘       │
                  │               │
         ┌────────┴────────┐      │
         │                 │      │
    Send/Recv         Close()     │
         │                 │      │
         ▼                 ▼      │
      ┌─────────┐    ┌──────────┐ │
      │ SENDING │    │ CLOSING  │ │
      └─────────┘    └────┬─────┘ │
                          │       │
                      Timeout     │
                          │       │
                          └───────┘


Bad Code (Without State Pattern):


func (c *Connection) Send(data []byte) error {
    if c.state == "CLOSED" {
        return errors.New("closed")
    }
    if c.state == "ESTABLISHED" {
        // Send logic
    }
    if c.state == "LISTENING" {
        return errors.New("can't send while listening")
    }
    // ... more conditions
}

Good Code (With State Pattern):


func (c *Connection) Send(data []byte) error {
    return c.currentState.Send(c, data)
}

// Each state handles its logic:
func (es *EstablishedState) Send(c *Connection, data []byte) error {
    // Send logic
}

func (cs *ClosedState) Send(c *Connection, data []byte) error {
    return errors.New("closed")
}
```

## Real-World Examples

### 1. TCP Connection State Machine

```go
type State interface {
    Open(ctx *Connection) error
    Close(ctx *Connection) error
    Send(ctx *Connection, data []byte) error
}

type ClosedState struct{}
func (s *ClosedState) Open(ctx *Connection) error {
    ctx.state = &EstablishedState{}
    return nil
}

type EstablishedState struct{}
func (s *EstablishedState) Send(ctx *Connection, data []byte) error {
    // Send data
    return nil
}
```

### 2. Document Workflow States

```go
type DocumentState interface {
    Publish(doc *Document) error
    Reject(doc *Document) error
}

type DraftState struct{}
func (s *DraftState) Publish(doc *Document) error {
    doc.state = &ReviewState{}
    return nil
}

type ReviewState struct{}
func (s *ReviewState) Reject(doc *Document) error {
    doc.state = &DraftState{}
    return nil
}
```

### 3. Game Character States

```go
type CharacterState interface {
    Update(char *Character)
    HandleInput(char *Character, input Input)
}

type IdleState struct{}
func (s *IdleState) HandleInput(char *Character, input Input) {
    if input == MoveInput {
        char.state = &WalkingState{}
    }
}

type DeadState struct{}
func (s *DeadState) HandleInput(char *Character, input Input) {
    // Dead characters ignore input
}
```

## Key Advantages

- **Clarity**: Each state encapsulated in its own class
- **Maintainability**: No massive conditional statements
- **Easy to extend**: Add new states without modifying existing code
- **Single responsibility**: Each state handles its own behavior
- **Explicit transitions**: Clear state transition logic
- **Testability**: Test each state independently
- **Eliminates duplication**: State-specific logic not scattered

## Key Gotchas

- **Over-engineering**: Overkill for simple state logic
- **Context dependencies**: States may need many context references
- **Memory overhead**: Each state instance adds memory
- **Transition complexity**: Complex transition logic can be hard to follow
- **State leakage**: States may expose implementation details
- **Invalid transitions**: Must prevent invalid state transitions
- **Circular dependencies**: States may reference each other
- **Thread safety**: State changes must be thread-safe
