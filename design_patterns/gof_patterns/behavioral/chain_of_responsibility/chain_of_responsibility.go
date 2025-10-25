package chain_of_responsibility

import "fmt"

type Handler interface {
	SetNext(Handler)
	Handle(request string)
}

type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BaseHandler) Handle(request string) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

type ConcreteHandler1 struct {
	BaseHandler
}

func (h *ConcreteHandler1) Handle(request string) {
	if request == "Request1" {
		fmt.Println("Handler1: Processing", request)
	} else {
		fmt.Println("Handler1: Passing to next handler")
		h.BaseHandler.Handle(request)
	}
}

type ConcreteHandler2 struct {
	BaseHandler
}

func (h *ConcreteHandler2) Handle(request string) {
	if request == "Request2" {
		fmt.Println("Handler2: Processing", request)
	} else {
		fmt.Println("Handler2: Passing to next handler")
		h.BaseHandler.Handle(request)
	}
}
