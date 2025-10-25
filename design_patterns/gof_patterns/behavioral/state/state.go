package state

import "fmt"

type State interface {
	Handle(context *Context)
}

type Context struct {
	state State
}

func NewContext(state State) *Context {
	return &Context{state: state}
}

func (c *Context) SetState(state State) {
	c.state = state
}

func (c *Context) Request() {
	c.state.Handle(c)
}

type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle(context *Context) {
	fmt.Println("State A: Transitioning to State B")
	context.SetState(&ConcreteStateB{})
}

type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(context *Context) {
	fmt.Println("State B: Transitioning to State A")
	context.SetState(&ConcreteStateA{})
}
