# Bridge Pattern

## Overview
The Bridge pattern decouples an abstraction from its implementation so that both can vary independently. It splits a large class or set of closely related classes into two separate hierarchies—abstraction and implementation—which can be developed independently.

## Problem
Imagine you're building a graphics application that needs to render shapes (Circle, Rectangle) on different platforms (Windows, Linux, macOS). Without Bridge, you'd need separate classes for each combination: WindowsCircle, LinuxCircle, MacCircle, WindowsRectangle, LinuxRectangle, MacRectangle—creating an explosion of classes.

Each time you add a new shape OR a new platform, the number of classes multiplies exponentially. This violates the Single Responsibility Principle and makes the codebase unmaintainable.

## Why Use This Pattern?

**Benefits:**
- **Reduced Class Explosion**: N shapes × M platforms = N+M classes (not N×M)
- **Independent Extension**: Add shapes or platforms without affecting each other
- **Runtime Binding**: Switch implementations at runtime
- **Better Organization**: Separate abstraction concerns from platform details

## When to Use

Use Bridge when:
- You want to avoid permanent binding between abstraction and implementation
- Both abstractions and implementations should be extensible via subclassing
- Changes in implementation shouldn't impact clients
- You have a proliferation of classes from combining two dimensions of variation
- You want to share an implementation among multiple objects

**Real-world scenarios:**
- GUI frameworks (shapes × rendering engines)
- Database drivers (queries × database types)
- Messaging systems (message types × transport protocols)
- Device drivers (device types × operating systems)
- Payment processing (payment methods × payment gateways)

## When NOT to Use

Avoid this pattern when:
- You only have one implementation
- Abstraction and implementation rarely change
- The added complexity outweighs benefits
- Simple inheritance is sufficient

## Implementation Guidelines

**Key Components:**
1. **Abstraction**: Defines high-level control interface, maintains reference to Implementation
2. **Refined Abstraction**: Extends Abstraction with variants
3. **Implementation**: Defines interface for implementation classes
4. **Concrete Implementation**: Provides specific implementations

**Go-Specific Considerations:**
- Use interfaces for both Abstraction and Implementation
- Composition over inheritance (embed interfaces)
- Constructor functions initialize abstraction with implementation

## Visual Schema

```
Abstraction Hierarchy          Implementation Hierarchy
--------------------          ------------------------
     Abstraction                  Implementation
         |                             |
 has-a implementation          +-------+-------+
         |                     |               |
    +----+----+          ConcreteImplA   ConcreteImplB
    |         |
RefinedA  RefinedB

Client creates:
  abstraction = new RefinedA(new ConcreteImplA())
  abstraction = new RefinedB(new ConcreteImplB())
```

**Flow Diagram:**
```
1. Client creates implementation (e.g., WindowsRenderer)
2. Client creates abstraction with implementation (e.g., Circle with WindowsRenderer)
3. Client calls abstraction methods (circle.Draw())
4. Abstraction delegates to implementation (renderer.RenderCircle())
5. Implementation executes platform-specific code
```

## Real-World Examples

### 1. Remote Control and Devices
Remote control (abstraction) operates different devices (TV, Radio) through a common interface. You can add new remotes or devices independently.

### 2. Database Abstraction Layer
Query builder (abstraction) works with different database drivers (MySQL, PostgreSQL) through a common connection interface.

### 3. Notification System
Notification types (Email, SMS, Push) sent through different channels (SMTP, Twilio, FCM).

## Key Advantages

- **Decoupling**: Abstraction and implementation vary independently
- **Improved Extensibility**: Add new abstractions or implementations easily
- **Hide Implementation Details**: Clients work only with abstraction
- **Open/Closed Principle**: Open for extension, closed for modification

## Key Gotchas

- **Increased Complexity**: More interfaces and indirection
- **Overhead**: May be overkill for simple scenarios
- **Design Upfront**: Requires identifying proper abstraction/implementation split

## Go Idioms

Bridge aligns well with Go's composition model:
```go
type Shape struct {
    renderer Renderer // Bridge to implementation
}

type Circle struct {
    Shape
    radius float64
}
```

## Code Example Structure

The implementation in `bridge.go` demonstrates:
- Shape abstraction (Circle, Rectangle)
- Rendering implementation (Vector, Raster)
- Independent extension of both hierarchies
- Runtime switching of implementations
