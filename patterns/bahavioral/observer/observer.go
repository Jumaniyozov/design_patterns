package main

import "fmt"

// Когда использовать: Когда изменение состояния одного объекта требует изменения других объектов.

// Наблюдатель
type Observer interface {
	update(string)
}

// Субъект
type Subject struct {
	observers []Observer
	state     string
}

func (s *Subject) attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *Subject) setState(state string) {
	s.state = state
	s.notifyAll()
}

func (s *Subject) notifyAll() {
	for _, observer := range s.observers {
		observer.update(s.state)
	}
}

// Конкретные наблюдатели
type ConcreteObserver struct {
	name string
}

func (c *ConcreteObserver) update(state string) {
	fmt.Printf("%s получил обновление состояния: %s\n", c.name, state)
}

func main() {
	subject := &Subject{}

	observer1 := &ConcreteObserver{name: "Наблюдатель1"}
	observer2 := &ConcreteObserver{name: "Наблюдатель2"}

	subject.attach(observer1)
	subject.attach(observer2)

	subject.setState("Новое состояние")
}
