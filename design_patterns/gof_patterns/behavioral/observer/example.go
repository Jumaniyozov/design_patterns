package observer

import "fmt"

func Example1_NewsPublisher() {
	fmt.Println("\n=== Example 1: News Publisher ===")

	publisher := NewConcreteSubject()
	subscriber1 := NewConcreteObserver("Subscriber 1")
	subscriber2 := NewConcreteObserver("Subscriber 2")

	publisher.Attach(subscriber1)
	publisher.Attach(subscriber2)

	publisher.SetState("Breaking News!")
	publisher.SetState("Weather Update")
}
