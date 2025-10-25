package memento

import "fmt"

type Memento struct {
	state string
}

type Originator struct {
	state string
}

func (o *Originator) SetState(state string) {
	fmt.Printf("Setting state to: %s\n", state)
	o.state = state
}

func (o *Originator) GetState() string {
	return o.state
}

func (o *Originator) SaveToMemento() *Memento {
	fmt.Printf("Saving state: %s\n", o.state)
	return &Memento{state: o.state}
}

func (o *Originator) RestoreFromMemento(m *Memento) {
	o.state = m.state
	fmt.Printf("Restored state: %s\n", o.state)
}

type Caretaker struct {
	mementos []*Memento
}

func (c *Caretaker) Add(m *Memento) {
	c.mementos = append(c.mementos, m)
}

func (c *Caretaker) Get(index int) *Memento {
	return c.mementos[index]
}
