package chain_of_responsibility

import "fmt"

func Example1_RequestChain() {
	fmt.Println("\n=== Example 1: Request Chain ===")

	handler1 := &ConcreteHandler1{}
	handler2 := &ConcreteHandler2{}

	handler1.SetNext(handler2)

	handler1.Handle("Request1")
	fmt.Println()
	handler1.Handle("Request2")
}
