package strategy

import "fmt"

// RunAllExamples executes all Strategy Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Strategy pattern.
//
// Usage: Call from cmd/main.go or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          STRATEGY PATTERN - COMPREHENSIVE EXAMPLES             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	// Run all examples in sequence
	Example1_BasicPaymentProcessing()
	Example2_SwitchingStrategiesAtRuntime()
	Example3_ErrorHandlingPerStrategy()
	Example4_MultiplePaymentsWithDifferentMethods()
	Example5_DynamicStrategySelection()
	Example6_StrategyComparison()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	fmt.Println("\n📝 Key Takeaways from Strategy Pattern:")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("1. ✓ Eliminates conditional logic (if/switch statements)")
	fmt.Println("2. ✓ Allows runtime algorithm switching without code changes")
	fmt.Println("3. ✓ Makes testing individual algorithms independent and simple")
	fmt.Println("4. ✓ Enables adding new strategies without modifying existing code")
	fmt.Println("5. ✓ Improves code organization and maintainability")
	fmt.Println("\n🧪 To run the test suite:")
	fmt.Println("────────────────────────")
	fmt.Println("   go test -v ./tier1/strategy/")
	fmt.Println("   go test -cover ./tier1/strategy/")
	fmt.Println("   go test -bench=. ./tier1/strategy/")
}
