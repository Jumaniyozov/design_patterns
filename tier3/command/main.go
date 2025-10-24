package command

import "fmt"

// RunAllExamples executes all Command Pattern examples in sequence.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          COMMAND PATTERN - COMPREHENSIVE EXAMPLES              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples
	Example1_TextEditorBasic()
	Example2_TextEditorAdvanced()
	Example3_MacroCommands()
	Example4_BankingTransactions()
	Example5_AsyncCommandQueue()
	Example6_FunctionCommands()
	Example7_ComplexMacro()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}
