package main

import "fmt"

// Когда использовать: Когда нужно передавать запрос по цепочке обработчиков до тех пор, пока один из них не обработает запрос.

// Интерфейс обработчика
type Handler interface {
	setNext(Handler)
	execute(*Request)
}

// Запрос
type Request struct {
	amount int
}

// Конкретные обработчики
type ConcreteHandler1 struct {
	next Handler
}

func (h *ConcreteHandler1) setNext(next Handler) {
	h.next = next
}

func (h *ConcreteHandler1) execute(r *Request) {
	if r.amount < 10 {
		fmt.Println("ConcreteHandler1 обработал запрос")
	} else if h.next != nil {
		h.next.execute(r)
	}
}

type ConcreteHandler2 struct {
	next Handler
}

func (h *ConcreteHandler2) setNext(next Handler) {
	h.next = next
}

func (h *ConcreteHandler2) execute(r *Request) {
	if r.amount >= 10 {
		fmt.Println("ConcreteHandler2 обработал запрос")
	} else if h.next != nil {
		h.next.execute(r)
	}
}

func main() {
	handler1 := &ConcreteHandler1{}
	handler2 := &ConcreteHandler2{}
	handler1.setNext(handler2)

	request := &Request{amount: 15}
	handler1.execute(request)
}
