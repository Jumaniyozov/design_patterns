# Interpreter Pattern

## Overview
The Interpreter pattern defines a grammatical representation for a language and an interpreter to interpret sentences in that language. It's used for parsing and evaluating structured languages or expressions.

## Problem
You need to evaluate expressions, parse configuration files, or process domain-specific languages (DSLs). While you could use regex or string parsing, complex grammars require structured interpretation with proper grammar rules.

## Why Use This Pattern?

**Benefits:**
- **Extensibility**: Easy to add new grammar rules
- **Modularity**: Each rule is a separate class
- **Flexibility**: Compose expressions dynamically

## When to Use

Use Interpreter when:
- Grammar is simple and stable
- Efficiency isn't critical
- You're building a DSL or query language
- Expression evaluation is needed

**Real-world scenarios:**
- Mathematical expression evaluators
- SQL query parsers
- Configuration languages
- Regular expression engines
- Boolean logic evaluators

## When NOT to Use

Avoid when:
- Grammar is complex (use parser generators instead)
- Performance is critical
- Grammar changes frequently

## Code Example Structure

Demonstrates:
- Boolean expression interpreter
- Mathematical expression evaluator
- Simple query language
