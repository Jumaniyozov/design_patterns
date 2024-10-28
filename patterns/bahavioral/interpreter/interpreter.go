package main

import (
	"fmt"
	"strings"
)

//Когда использовать: Когда нужно определить грамматику простого языка и интерпретировать предложения в этом языке.

// Интерфейс выражения
type Expression interface {
	interpret(context string) bool
}

// Конкретные выражения
type TerminalExpression struct {
	data string
}

func (t *TerminalExpression) interpret(context string) bool {
	return strings.Contains(context, t.data)
}

func getMaleExpression() Expression {
	john := &TerminalExpression{data: "Джон"}
	mike := &TerminalExpression{data: "Майк"}
	return &OrExpression{expr1: john, expr2: mike}
}

// Комбинированные выражения
type OrExpression struct {
	expr1, expr2 Expression
}

func (o *OrExpression) interpret(context string) bool {
	return o.expr1.interpret(context) || o.expr2.interpret(context)
}

func main() {
	isMale := getMaleExpression()
	fmt.Println("Джон мужчина?", isMale.interpret("Джон"))
	fmt.Println("Джулия мужчина?", isMale.interpret("Джулия"))
}
