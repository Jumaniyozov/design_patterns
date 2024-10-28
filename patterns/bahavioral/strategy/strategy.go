package main

import "fmt"

// Когда использовать: Когда есть несколько алгоритмов для выполнения операции, и нужно выбрать один из них во время выполнения.

// Стратегия
type Strategy interface {
	execute(int, int) int
}

// Контекст
type Context struct {
	strategy Strategy
}

func (c *Context) setStrategy(s Strategy) {
	c.strategy = s
}

func (c *Context) executeStrategy(a, b int) int {
	return c.strategy.execute(a, b)
}

// Конкретные стратегии
type Addition struct{}

func (a *Addition) execute(num1, num2 int) int {
	return num1 + num2
}

type Subtraction struct{}

func (s *Subtraction) execute(num1, num2 int) int {
	return num1 - num2
}

func main() {
	context := &Context{}

	context.setStrategy(&Addition{})
	fmt.Println("10 + 5 =", context.executeStrategy(10, 5))

	context.setStrategy(&Subtraction{})
	fmt.Println("10 - 5 =", context.executeStrategy(10, 5))
}
