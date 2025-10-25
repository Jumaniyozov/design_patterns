# Decorator Pattern

## Overview

The Decorator pattern attaches additional responsibilities to an object dynamically without affecting other objects of the same class. Decorators provide a flexible alternative to subclassing for extending functionality.

## Problem

You're building a coffee shop system where customers can customize drinks with various add-ons (milk, sugar, whipped cream, caramel). Without decorators, you'd need a class for every combination (CoffeeWithMilk, CoffeeWithMilkAndSugar, etc.), leading to class explosion.

## Why Use This Pattern?

- **Flexible Extension**: Add responsibilities dynamically
- **Open/Closed**: Extend without modifying existing code
- **Composition Over Inheritance**: More flexible than subclassing
- **Single Responsibility**: Each decorator has one responsibility
- **Runtime Configuration**: Combine decorators at runtime

## When to Use

- **Add Responsibilities Dynamically**: Attach features at runtime
- **Avoid Class Explosion**: Too many subclass combinations
- **Cross-Cutting Concerns**: Logging, caching, authentication
- **Middleware**: HTTP middleware, data processing pipelines
- **UI Components**: Adding borders, scrollbars, tooltips

## Key Advantages

✓ **More flexible than inheritance**
✓ **Avoid class explosion**  
✓ **Single responsibility per decorator**
✓ **Composable at runtime**

## Best Practices

1. Keep decorators focused on single concern
2. Maintain same interface as component
3. Use composition, not inheritance
4. Allow stacking multiple decorators
