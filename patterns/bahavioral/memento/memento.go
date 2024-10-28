package main

import "fmt"

// Когда использовать: Когда нужно сохранить текущее состояние объекта и восстановить его позже.

// Снимок
type Memento struct {
	state string
}

// Создатель
type Originator struct {
	state string
}

func (o *Originator) createMemento() *Memento {
	return &Memento{state: o.state}
}

func (o *Originator) restore(m *Memento) {
	o.state = m.state
}

// Хранитель
type Caretaker struct {
	mementos []*Memento
}

func (c *Caretaker) addMemento(m *Memento) {
	c.mementos = append(c.mementos, m)
}

func (c *Caretaker) getMemento(index int) *Memento {
	return c.mementos[index]
}

func main() {
	originator := &Originator{}
	caretaker := &Caretaker{}

	originator.state = "Состояние1"
	caretaker.addMemento(originator.createMemento())

	originator.state = "Состояние2"
	caretaker.addMemento(originator.createMemento())

	originator.restore(caretaker.getMemento(0))
	fmt.Println("Восстановлено до:", originator.state)
}
