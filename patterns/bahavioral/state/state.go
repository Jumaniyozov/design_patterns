package main

import "fmt"

// Когда использовать: Когда поведение объекта должно меняться в зависимости от его состояния.

// Состояние
type State interface {
	doAction(*Context)
}

// Контекст
type Context struct {
	state State
}

func (c *Context) setState(s State) {
	c.state = s
}

func (c *Context) request() {
	c.state.doAction(c)
}

// Конкретные состояния
type StartState struct{}

func (s *StartState) doAction(c *Context) {
	fmt.Println("Плеер в состоянии старта")
	c.setState(s)
}

type StopState struct{}

func (s *StopState) doAction(c *Context) {
	fmt.Println("Плеер в состоянии остановки")
	c.setState(s)
}

func main() {
	context := &Context{}

	startState := &StartState{}
	startState.doAction(context)

	stopState := &StopState{}
	stopState.doAction(context)
}
