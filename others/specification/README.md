# Specification Pattern

## Overview
The Specification pattern encapsulates business rules as composable objects that can be combined using boolean logic (AND, OR, NOT). It separates query/validation logic from business objects.

## Problem
Business rules scattered throughout code are hard to maintain and reuse. Complex conditions become unreadable. You need composable, reusable business rules.

## Why Use This Pattern?
- **Composability**: Combine rules with AND/OR/NOT
- **Reusability**: Use same rules in queries and validation
- **Testability**: Test rules in isolation
- **Clarity**: Self-documenting business logic

## When to Use
- Complex selection criteria
- Business rule validation
- Query building
- Filtering collections

## Real-world scenarios
- Product filtering (price range AND category AND in-stock)
- User permissions (role-based access control)
- Data validation rules
- Search queries

## Code Example Structure
Demonstrates composable specifications for filtering and validation.
