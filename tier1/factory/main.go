package factory

import "fmt"

// RunAllExamples executes all Factory Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Factory pattern.
//
// Usage: Call from cmd/main.go or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          FACTORY PATTERN - COMPREHENSIVE EXAMPLES              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples in sequence
	Example1_BasicFactoryUsage()
	Example2_MultiplePaymentMethods()
	Example3_SpecializedFactories()
	Example4_ErrorHandling()
	Example5_RealWorldScenario()
	Example6_FactoryWithDependencyInjection()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   EXAMPLES COMPLETED SUCCESSFULLY              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	fmt.Println("\n📝 Key Takeaways from Factory Pattern:")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("1. ✓ Decouples client code from concrete implementations")
	fmt.Println("2. ✓ Centralizes object creation logic in one place")
	fmt.Println("3. ✓ Makes extending with new types easy (Open/Closed Principle)")
	fmt.Println("4. ✓ Handles complex initialization and validation logic")
	fmt.Println("5. ✓ Improves testability through dependency abstraction")
	fmt.Println("6. ✓ Idiomatic in Go through 'New*' constructor functions")
	fmt.Println("\n🧪 To run the test suite:")
	fmt.Println("────────────────────────")
	fmt.Println("   go test -v ./tier1/factory/")
	fmt.Println("   go test -cover ./tier1/factory/")
	fmt.Println("   go test -bench=. ./tier1/factory/")
}
