# Template Method Pattern

## Overview

The Template Method Pattern is a behavioral design pattern that defines the skeleton of an algorithm in a method, deferring some steps to subclasses. It lets subclasses redefine certain steps of an algorithm without changing the algorithm's structure.

## Problem

When you have multiple classes that perform similar algorithms with slight variations, code duplication becomes a problem:

- **Code duplication**: The same algorithm structure repeated in multiple classes
- **Difficult maintenance**: Changes to shared logic require updates in multiple places
- **Hard to extend**: Adding variations requires careful code duplication
- **Unclear structure**: Algorithm intent obscured by implementation details
- **Scattered logic**: Related algorithmic steps separated across classes
- **Inheritance hierarchy**: Leads to deep inheritance hierarchies

### Real-World Context

Consider different document formats (PDF, Word, HTML) that all need to be opened, processed, and saved in the same order, but with format-specific implementations. The Template Method defines the overall process (open → validate → process → save) while allowing each format to implement format-specific details.

## Why Use This Pattern?

- **Code reuse**: Common algorithm structure in one place
- **Inversion of control**: Subclasses hook into algorithm at specific points
- **Consistency**: Ensures all variants follow the same process
- **Easy to extend**: Add new variants without duplicating algorithm
- **Clarity**: Algorithm intent clear from main method
- **Go idioms**: Composition and interfaces fit better than inheritance

## When to Use

- Multiple classes perform similar algorithms with variations
- Want to avoid duplication of common algorithm steps
- Defining hooks where subclasses can customize behavior
- Creating frameworks or libraries with customization points
- Processing pipelines with format-specific steps
- Workflows with standard sequence but variable implementations

## When NOT to Use

- Simple algorithms without much variation
- Only one implementation needed
- Subclasses need completely different behavior
- Algorithm steps have no logical ordering
- Composition/strategy is more appropriate
- Deep inheritance hierarchies to avoid

## Implementation Guidelines

1. **Template method**: Define algorithm skeleton in base class
2. **Abstract operations**: Steps that subclasses must implement
3. **Concrete subclasses**: Implement specific steps
4. **Hooks**: Optional override points for customization
5. **Control flow**: Template method controls order of operations

## Go Idioms

Go's composition over inheritance approach:

- Use interfaces instead of base classes
- Embed structs to share behavior
- Function types represent customization points
- Avoid deep inheritance; prefer composition

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│    Template Method Pattern Architecture      │
└──────────────────────────────────────────────┘

Base Class (Algorithm Structure):
  ┌────────────────────────────────┐
  │ BaseAlgorithm                  │
  │                                │
  │ TemplateMethod() {             │
  │   Step1()                      │
  │   Step2()                      │
  │   Step3()                      │
  │   Step4()                      │
  │ }                              │
  │                                │
  │ Step1() { ... }  ◄─ Concrete   │
  │                                │
  │ abstract Step2()   ◄─ Abstract │
  │ abstract Step3()   ◄─ Abstract │
  │                                │
  │ Step4() { ... }  ◄─ Hook       │
  └────────┬───────────────────────┘
           │
       inherits
           │
    ┌──────┴────────┬──────────────┐
    │               │              │
    ▼               ▼              ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│Concrete1 │ │Concrete2 │ │Concrete3 │
│          │ │          │ │          │
│Step2()   │ │Step2()   │ │Step2()   │
│Step3()   │ │Step3()   │ │Step3()   │
│Step4()   │ │Step4()   │ │Step4()   │
└──────────┘ └──────────┘ └──────────┘


Document Processing Example:

Template Method:
  ┌────────────────────────────────────┐
  │ Process() {                        │
  │   Open()          ◄──┐             │
  │   Validate()      ◄──┤ Defined     │
  │   Transform()     ◄──┤ in child    │
  │   Save()          ◄──┘ classes     │
  │ }                                  │
  └────────────────────────────────────┘

PDFDocument:           WordDocument:       HTMLDocument:
│                      │                   │
├─ Open: PDF parser    ├─ Open: Word API   ├─ Open: HTML parser
├─ Validate: Check     ├─ Validate: Check  ├─ Validate: Check
│  signatures          │  compatibility    │  well-formed
├─ Transform: Render   ├─ Transform:       ├─ Transform:
│  to bitmap           │  to text          │  to DOM tree
└─ Save: PDF stream    └─ Save: DOC        └─ Save: HTML

All follow same Process() structure
but implement steps differently.

Without Template Method (Duplication):

func ProcessPDF(filename string) error {
    pdf := OpenPDF(filename)
    ValidatePDF(pdf)
    TransformPDF(pdf)
    SavePDF(pdf)
}

func ProcessWord(filename string) error {
    word := OpenWord(filename)
    ValidateWord(word)      ◄── Duplicated
    TransformWord(word)     ◄── structure
    SaveWord(word)
}

With Template Method (DRY):

func (d *Document) Process() error {
    d.Open()
    d.Validate()
    d.Transform()
    d.Save()
}

// Each subclass just implements steps
type PDFDocument struct{}
func (p *PDFDocument) Open() { /*...*/ }
func (p *PDFDocument) Transform() { /*...*/ }
```

## Real-World Examples

### 1. Document Processing

```go
type Document interface {
    Open(filename string) error
    Validate() error
    Transform() error
    Save(filename string) error
    Process(filename string) error
}

type BaseDocument struct {
    content []byte
}

func (d *BaseDocument) Process(filename string) error {
    d.Open(filename)
    d.Validate()
    d.Transform()
    d.Save(filename)
    return nil
}

type PDFDocument struct{ BaseDocument }
func (p *PDFDocument) Transform() error { /*...*/ }
```

### 2. Data Processing Pipeline

```go
type DataProcessor interface {
    Load() error
    Validate() error
    Process() error
    Export() error
    Execute() error
}

type CSVProcessor struct{}
func (c *CSVProcessor) Execute() error {
    c.Load()
    c.Validate()
    c.Process()
    c.Export()
}

type JSONProcessor struct{}
func (j *JSONProcessor) Execute() error {
    j.Load()
    j.Validate()
    j.Process()
    j.Export()
}
```

### 3. Game Entity AI

```go
type GameEntity interface {
    Think() error
    Move() error
    Act() error
    Update() error
}

type BaseEntity struct{}
func (e *BaseEntity) Update() error {
    e.Think()
    e.Move()
    e.Act()
    return nil
}

type Zombie struct{ BaseEntity }
func (z *Zombie) Think() error { /*zombie logic*/ }

type Guard struct{ BaseEntity }
func (g *Guard) Think() error { /*guard logic*/ }
```

## Key Advantages

- **Code reuse**: Common algorithm structure defined once
- **Consistency**: All variants follow the same process
- **Easy to extend**: Add new variants without duplicating algorithm
- **Clarity**: Algorithm intent clear from template method
- **Inversion of control**: Subclasses don't control flow; template does
- **Maintainability**: Changes to common logic in one place
- **Flexibility**: Subclasses can customize specific steps

## Key Gotchas

- **Deep inheritance**: Can lead to complex inheritance hierarchies
- **Coupling**: Subclasses tightly coupled to base class
- **Limited flexibility**: Algorithm structure is fixed
- **Over-engineering**: Overkill for simple algorithms
- **Violation of Liskov**: Subclasses may break algorithm assumptions
- **Hard to refactor**: Changing algorithm structure affects all subclasses
- **Go philosophy**: Go prefers composition over inheritance
- **Limited testing**: Hard to test algorithm and variations independently
