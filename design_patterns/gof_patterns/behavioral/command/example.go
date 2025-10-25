package command

import "fmt"

func Example1_RemoteControl() {
	fmt.Println("\n=== Example 1: Remote Control ===")

	light := &Light{}
	onCmd := &LightOnCommand{light: light}
	offCmd := &LightOffCommand{light: light}

	remote := &RemoteControl{history: make([]Command, 0)}

	remote.SetCommand(onCmd)
	remote.PressButton()

	remote.SetCommand(offCmd)
	remote.PressButton()

	fmt.Println("\nUndo last command:")
	remote.PressUndo()
}
