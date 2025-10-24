package state

import "fmt"

// RunAllExamples executes all State Pattern examples in sequence.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║            STATE PATTERN - COMPREHENSIVE EXAMPLES              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples
	Example1_TCPConnection()
	Example2_DocumentWorkflow()
	Example3_PlayerCharacter()
	Example4_StateTransitions()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}
