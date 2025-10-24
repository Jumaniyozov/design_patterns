package templatemethod

import "fmt"

// RunAllExamples executes all Template Method Pattern examples in sequence.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║        TEMPLATE METHOD PATTERN - COMPREHENSIVE EXAMPLES        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples
	Example1_DataProcessing()
	Example2_ReportGeneration()
	Example3_GameAI()
	Example4_MultipleProcessors()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}
