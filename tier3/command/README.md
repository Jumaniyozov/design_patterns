# Command Pattern

## Overview

The Command Pattern is a behavioral design pattern that encapsulates a request as an object, allowing you to parameterize clients with different requests, queue requests, log requests, and support undoable operations.

## Problem

When you need to decouple the object that invokes an operation from the object that performs the operation, direct method calls become problematic:

- **Tight coupling**: Invoker depends directly on receiver classes
- **Undo/Redo**: Hard to implement without storing operation history
- **Request queuing**: Difficult to queue and execute requests later
- **Logging and auditing**: No central place to log operations
- **Delayed execution**: Can't easily defer operation execution
- **Conditional logic**: Method calls scattered throughout code

### Real-World Context

Consider a text editor with undo/redo functionality. Each user action (delete text, insert character, format) should be an undoable command. Instead of having the UI directly call the text buffer, commands encapsulate each operation. This allows maintaining a history for undo/redo.

## Why Use This Pattern?

- **Decoupling**: Invoker doesn't depend on concrete receiver classes
- **Undo/Redo**: Easy to implement with command history
- **Request queuing**: Queue commands for later execution
- **Logging and auditing**: Central logging of all operations
- **Composability**: Combine commands into macros or sequences
- **Go concurrency**: Commands work naturally with goroutines

## When to Use

- Implementing undo/redo functionality
- Queuing operations for later execution
- Scheduling and executing tasks
- Building command-driven interfaces (REPL, CLI)
- Transaction management
- Macro or script recording
- Request logging and audit trails
- Delayed or asynchronous operations

## When NOT to Use

- Simple synchronous operations where direct calls suffice
- No need for undo/redo or delayed execution
- Simple object with few operations
- Performance critical code where command overhead matters
- Clear control flow is more important than decoupling

## Implementation Guidelines

1. **Command interface**: Defines Execute/Undo operations
2. **Concrete commands**: Implements command for specific operations
3. **Receiver**: The object that actually performs the work
4. **Invoker**: Triggers commands without knowing their details
5. **History**: Store commands for undo/redo functionality

## Go Idioms

Go's goroutines and channels make commands natural for async operations:

- Commands fit well in channels for async execution
- Closures can capture context for command execution
- Function types can represent simple commands
- Goroutines enable concurrent command processing

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│        Command Pattern Architecture          │
└──────────────────────────────────────────────┘

Command Interface:
  ┌────────────────────┐
  │ Command            │
  │ Execute()          │
  │ Undo()             │
  └────────────────────┘
           ▲
           │
    ┌──────┴──────────────┬──────────────┐
    │                     │              │
    ▼                     ▼              ▼
┌──────────┐      ┌──────────┐      ┌───────┐
│InsertCmd │      │DeleteCmd │      │CopyCmd│
│          │      │          │      │       │
│Receiver: │      │Receiver: │      │Recv:  │
│TextBuffer│      │TextBuffer│      │Buffer │
└──────────┘      └──────────┘      └───────┘

Execution Flow:

┌────────────┐
│   Client   │
└────┬───────┘
     │
     ▼
┌──────────────┐
│  Invoker     │ (e.g., Editor)
│              │
│ Execute(cmd) │
│ Undo()       │
└────┬─────────┘
     │
     ▼
┌───────────────┐     ┌───────────┐
│   Command     │────→│ Receiver  │
│ (Encapsulates │     │ (Actually │
│  operation    │     │ does work)│
│  and state)   │     └───────────┘
└───────────────┘

Undo/Redo with History:

┌─────────────────────────────────┐
│  Invoker (Editor)               │
│                                 │
│ ┌──────────────────────────────┐│
│ │ Execute(command)             ││
│ │ ├─ command.Execute()         ││
│ │ └─ history.push(command)     ││
│ │                              ││
│ │ Undo()                       ││
│ │ ├─ cmd = history.pop()       ││
│ │ └─ cmd.Undo()                ││
│ │                              ││
│ │ Redo()                       ││
│ │ ├─ cmd = undone.pop()        ││
│ │ └─ cmd.Execute()             ││
│ └──────────────────────────────┘│
│                                 │
│ ┌──────────────────────────────┐│
│ │ History Stack                ││
│ │ [Insert, Delete, Format,... ]││
│ └──────────────────────────────┘│
└─────────────────────────────────┘

Text Editor Command Chain:
1. User types 'A'
   └─ InsertCommand("A").Execute()

2. User deletes
   └─ DeleteCommand(pos).Execute()

3. User presses Ctrl+Z
   └─ DeleteCommand.Undo()

4. User presses Ctrl+Y
   └─ DeleteCommand.Execute()
```

## Real-World Examples

### 1. Text Editor Undo/Redo

```go
type Command interface {
    Execute() error
    Undo() error
}

type InsertCommand struct {
    buffer *TextBuffer
    text   string
    pos    int
}

type TextEditor struct {
    buffer   *TextBuffer
    history  []Command
    undoPos  int
}

func (e *TextEditor) Execute(cmd Command) error {
    cmd.Execute()
    e.history = append(e.history, cmd)
}

func (e *TextEditor) Undo() error {
    e.history[len(e.history)-1].Undo()
}
```

### 2. Job Queue with Commands

```go
type JobCommand interface {
    Execute() error
}

type ProcessDataJob struct { /*...*/ }
type SendEmailJob struct { /*...*/ }

type JobQueue struct {
    jobs chan JobCommand
}

func (q *JobQueue) Submit(job JobCommand) {
    q.jobs <- job
}
```

### 3. API Batch Operations

```go
type APICommand interface {
    Execute(api *API) error
}

type CreateUserCommand struct { /*...*/ }
type UpdateUserCommand struct { /*...*/ }
type DeleteUserCommand struct { /*...*/ }

// Queue and execute later
batch := []APICommand{create, update, delete}
for _, cmd := range batch {
    cmd.Execute(api)
}
```

## Key Advantages

- **Decoupling**: Invokers don't depend on concrete receiver classes
- **Undo/Redo**: Easy to implement operation history
- **Queuing**: Commands can be queued and executed later
- **Asynchronous execution**: Commands work well with goroutines and channels
- **Logging and auditing**: Central place to log all operations
- **Macros and sequences**: Combine commands into complex operations
- **Transactional**: Pair commands with rollback for transactions

## Key Gotchas

- **Memory overhead**: Storing commands for undo uses memory
- **State management**: Commands must preserve state for undo
- **Execution order**: Order of command execution may matter
- **Complexity**: Can be overkill for simple operations
- **Thread safety**: History must be thread-safe with goroutines
- **Large undo history**: May consume significant memory over time
- **Concurrent modifications**: Must handle race conditions carefully
