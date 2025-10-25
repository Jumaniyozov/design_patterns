package facade

import "fmt"

func Example1_ComputerFacade() {
	fmt.Println("\n=== Example 1: Computer Facade ===")
	computer := NewComputerFacade()
	computer.Start()
}
