package bridge

import "fmt"

func Example1_RemoteControl() {
	fmt.Println("\n=== Example 1: Remote Control ===")
	tv := &TV{}
	remote := NewRemoteControl(tv)
	remote.TogglePower()
	remote.VolumeUp()

	fmt.Println("\nSwitching to radio:")
	radio := &Radio{}
	remote2 := NewRemoteControl(radio)
	remote2.TogglePower()
	remote2.VolumeUp()
}
