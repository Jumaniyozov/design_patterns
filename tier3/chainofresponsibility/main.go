package chainofresponsibility

import "fmt"

// RunAllExamples executes all Chain of Responsibility Pattern examples in sequence.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║    CHAIN OF RESPONSIBILITY - COMPREHENSIVE EXAMPLES            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples
	Example1_HTTPMiddlewareChain()
	Example2_RateLimiting()
	Example3_ExpenseApproval()
	Example4_ValidationChain()
	Example5_ConfigurableChain()
	Example6_ShortCircuiting()
	Example7_MultipleHandlerTypes()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}
