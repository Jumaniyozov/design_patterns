# Strangler Fig Pattern

## Overview
The Strangler Fig pattern incrementally migrates a legacy system by gradually replacing specific functionality with new services. Like the strangler fig tree that grows around a host tree, the new system gradually replaces the old.

## Problem
Big-bang rewrites are risky and often fail. You need to migrate a legacy system gradually while maintaining business continuity, minimizing risk, and delivering value incrementally.

## Why Use This Pattern?
- **Incremental Migration**: Replace piece by piece
- **Reduced Risk**: Small, reversible changes
- **Continuous Delivery**: Deliver value during migration
- **Parallel Running**: Old and new systems coexist
- **Rollback Capability**: Easy to revert changes

## When to Use
- Legacy system modernization
- Monolith to microservices migration
- Technology stack upgrade
- Gradual platform migration

## Real-world scenarios
- Migrating from monolith to microservices
- Replacing legacy mainframe systems
- Moving from on-premise to cloud
- Framework/language upgrades

## Implementation Strategy
1. **Identify boundaries**: Find seams in the legacy system
2. **Add facade/proxy**: Route requests to old or new system
3. **Implement new feature**: Build replacement service
4. **Switch routing**: Gradually redirect traffic to new service
5. **Monitor**: Ensure new service works correctly
6. **Decommission**: Remove old code when fully replaced

## Go Idioms
```go
type Proxy struct {
    legacy LegacyService
    modern ModernService
    router RoutingStrategy
}

func (p *Proxy) HandleRequest(req Request) Response {
    if p.router.ShouldUseModern(req) {
        return p.modern.Handle(req)
    }
    return p.legacy.Handle(req)
}
```
