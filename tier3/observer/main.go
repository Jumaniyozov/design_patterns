package observer

import "fmt"

// RunAllExamples executes all Observer Pattern examples in sequence.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          OBSERVER PATTERN - COMPREHENSIVE EXAMPLES             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples
	Example1_BasicStockObserver()
	Example2_AlertSystem()
	Example3_DynamicObservers()
	Example4_ChannelBasedObserver()
	Example5_MultipleSubjects()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}
