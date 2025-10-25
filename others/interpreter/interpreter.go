// Package interpreter demonstrates the Interpreter pattern.
// It defines a grammar and interprets sentences in that language,
// useful for DSLs, expression evaluation, and parsers.
package interpreter

import "fmt"

// Expression interface for all expressions
type Expression interface {
	Interpret(context map[string]bool) bool
	String() string
}

// VariableExpression represents a variable
type VariableExpression struct {
	name string
}

func NewVariable(name string) *VariableExpression {
	return &VariableExpression{name: name}
}

func (v *VariableExpression) Interpret(context map[string]bool) bool {
	return context[v.name]
}

func (v *VariableExpression) String() string {
	return v.name
}

// AndExpression represents logical AND
type AndExpression struct {
	left, right Expression
}

func NewAnd(left, right Expression) *AndExpression {
	return &AndExpression{left: left, right: right}
}

func (a *AndExpression) Interpret(context map[string]bool) bool {
	return a.left.Interpret(context) && a.right.Interpret(context)
}

func (a *AndExpression) String() string {
	return fmt.Sprintf("(%s AND %s)", a.left.String(), a.right.String())
}

// OrExpression represents logical OR
type OrExpression struct {
	left, right Expression
}

func NewOr(left, right Expression) *OrExpression {
	return &OrExpression{left: left, right: right}
}

func (o *OrExpression) Interpret(context map[string]bool) bool {
	return o.left.Interpret(context) || o.right.Interpret(context)
}

func (o *OrExpression) String() string {
	return fmt.Sprintf("(%s OR %s)", o.left.String(), o.right.String())
}

// NotExpression represents logical NOT
type NotExpression struct {
	expr Expression
}

func NewNot(expr Expression) *NotExpression {
	return &NotExpression{expr: expr}
}

func (n *NotExpression) Interpret(context map[string]bool) bool {
	return !n.expr.Interpret(context)
}

func (n *NotExpression) String() string {
	return fmt.Sprintf("(NOT %s)", n.expr.String())
}

// Numeric expressions

// NumExpression interface for numeric expressions
type NumExpression interface {
	Evaluate() float64
	String() string
}

// Number represents a constant number
type Number struct {
	value float64
}

func NewNumber(value float64) *Number {
	return &Number{value: value}
}

func (n *Number) Evaluate() float64 {
	return n.value
}

func (n *Number) String() string {
	return fmt.Sprintf("%.2f", n.value)
}

// Add represents addition
type Add struct {
	left, right NumExpression
}

func NewAdd(left, right NumExpression) *Add {
	return &Add{left: left, right: right}
}

func (a *Add) Evaluate() float64 {
	return a.left.Evaluate() + a.right.Evaluate()
}

func (a *Add) String() string {
	return fmt.Sprintf("(%s + %s)", a.left.String(), a.right.String())
}

// Multiply represents multiplication
type Multiply struct {
	left, right NumExpression
}

func NewMultiply(left, right NumExpression) *Multiply {
	return &Multiply{left: left, right: right}
}

func (m *Multiply) Evaluate() float64 {
	return m.left.Evaluate() * m.right.Evaluate()
}

func (m *Multiply) String() string {
	return fmt.Sprintf("(%s * %s)", m.left.String(), m.right.String())
}
