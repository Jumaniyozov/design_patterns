package memento

import "fmt"

func Example1_TextEditor() {
	fmt.Println("\n=== Example 1: Text Editor Undo ===")

	originator := &Originator{}
	caretaker := &Caretaker{}

	originator.SetState("State 1")
	caretaker.Add(originator.SaveToMemento())

	originator.SetState("State 2")
	caretaker.Add(originator.SaveToMemento())

	originator.SetState("State 3")

	fmt.Println("\nRestoring to State 2:")
	originator.RestoreFromMemento(caretaker.Get(1))
}
