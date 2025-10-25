package mediator

import "fmt"

func Example1_ChatRoom() {
	fmt.Println("\n=== Example 1: Chat Room ===")

	mediator := NewChatMediator()
	user1 := NewUser("Alice", mediator)
	user2 := NewUser("Bob", mediator)
	_ = NewUser("Charlie", mediator)

	user1.Send("Hello everyone!")
	user2.Send("Hi Alice!")
}
