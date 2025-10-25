package state

import "fmt"

func Example1_StateMachine() {
	fmt.Println("\n=== Example 1: State Machine ===")

	context := NewContext(&ConcreteStateA{})
	context.Request()
	context.Request()
	context.Request()
}
