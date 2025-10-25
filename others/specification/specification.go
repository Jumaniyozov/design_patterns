// Package specification demonstrates the Specification pattern.
// It encapsulates business rules as composable objects that can be combined
// using boolean logic, making complex criteria more manageable.
package specification

// Specification interface for generic specifications
type Specification[T any] interface {
	IsSatisfiedBy(candidate T) bool
	And(other Specification[T]) Specification[T]
	Or(other Specification[T]) Specification[T]
	Not() Specification[T]
}

// BaseSpecification provides common specification operations
type BaseSpecification[T any] struct {
	predicate func(T) bool
}

// NewSpecification creates a specification from a predicate
func NewSpecification[T any](predicate func(T) bool) *BaseSpecification[T] {
	return &BaseSpecification[T]{predicate: predicate}
}

func (s *BaseSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.predicate(candidate)
}

func (s *BaseSpecification[T]) And(other Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: s, right: other}
}

func (s *BaseSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: s, right: other}
}

func (s *BaseSpecification[T]) Not() Specification[T] {
	return &NotSpecification[T]{spec: s}
}

// Composite specifications

type AndSpecification[T any] struct {
	left, right Specification[T]
}

func (a *AndSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return a.left.IsSatisfiedBy(candidate) && a.right.IsSatisfiedBy(candidate)
}

func (a *AndSpecification[T]) And(other Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: a, right: other}
}

func (a *AndSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: a, right: other}
}

func (a *AndSpecification[T]) Not() Specification[T] {
	return &NotSpecification[T]{spec: a}
}

type OrSpecification[T any] struct {
	left, right Specification[T]
}

func (o *OrSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return o.left.IsSatisfiedBy(candidate) || o.right.IsSatisfiedBy(candidate)
}

func (o *OrSpecification[T]) And(other Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: o, right: other}
}

func (o *OrSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: o, right: other}
}

func (o *OrSpecification[T]) Not() Specification[T] {
	return &NotSpecification[T]{spec: o}
}

type NotSpecification[T any] struct {
	spec Specification[T]
}

func (n *NotSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return !n.spec.IsSatisfiedBy(candidate)
}

func (n *NotSpecification[T]) And(other Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: n, right: other}
}

func (n *NotSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: n, right: other}
}

func (n *NotSpecification[T]) Not() Specification[T] {
	return n.spec
}

// Example: Product specifications

type Product struct {
	Name     string
	Price    float64
	Category string
	InStock  bool
}

// Concrete specifications for products

func PriceGreaterThan(minPrice float64) Specification[*Product] {
	return NewSpecification(func(p *Product) bool {
		return p.Price > minPrice
	})
}

func PriceLessThan(maxPrice float64) Specification[*Product] {
	return NewSpecification(func(p *Product) bool {
		return p.Price < maxPrice
	})
}

func InCategory(category string) Specification[*Product] {
	return NewSpecification(func(p *Product) bool {
		return p.Category == category
	})
}

func IsInStock() Specification[*Product] {
	return NewSpecification(func(p *Product) bool {
		return p.InStock
	})
}

// Filter filters a slice based on specification
func Filter[T any](items []T, spec Specification[T]) []T {
	result := make([]T, 0)
	for _, item := range items {
		if spec.IsSatisfiedBy(item) {
			result = append(result, item)
		}
	}
	return result
}

// User specification example

type User struct {
	Name   string
	Age    int
	Active bool
	Role   string
}

func AgeGreaterThan(minAge int) Specification[*User] {
	return NewSpecification(func(u *User) bool {
		return u.Age > minAge
	})
}

func IsActive() Specification[*User] {
	return NewSpecification(func(u *User) bool {
		return u.Active
	})
}

func HasRole(role string) Specification[*User] {
	return NewSpecification(func(u *User) bool {
		return u.Role == role
	})
}
