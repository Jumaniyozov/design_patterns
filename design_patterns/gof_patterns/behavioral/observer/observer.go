package observer

import "fmt"

type Observer interface {
	Update(message string)
}

type Subject interface {
	Attach(Observer)
	Detach(Observer)
	Notify()
}

type ConcreteSubject struct {
	observers []Observer
	state     string
}

func NewConcreteSubject() *ConcreteSubject {
	return &ConcreteSubject{observers: make([]Observer, 0)}
}

func (s *ConcreteSubject) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *ConcreteSubject) Detach(o Observer) {
	for i, obs := range s.observers {
		if obs == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *ConcreteSubject) Notify() {
	for _, o := range s.observers {
		o.Update(s.state)
	}
}

func (s *ConcreteSubject) SetState(state string) {
	s.state = state
	s.Notify()
}

type ConcreteObserver struct {
	name string
}

func NewConcreteObserver(name string) *ConcreteObserver {
	return &ConcreteObserver{name: name}
}

func (o *ConcreteObserver) Update(message string) {
	fmt.Printf("[%s] Received update: %s\n", o.name, message)
}
