package adapter

import "fmt"

// RunAllExamples executes all Adapter Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Adapter pattern.
//
// Usage: Call from cmd/main.go or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║          ADAPTER PATTERN - COMPREHENSIVE EXAMPLES         ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples in sequence
	Example1_PaymentProcessorWithoutAdapter()
	Example2_PaymentProcessorWithAdapter()
	Example3_PolymorphicPaymentProcessing()
	Example4_RefundWithAdapters()
	Example5_DatabaseAdapters()
	Example6_PolymorphicDatabaseOperations()
	Example7_LoggerAdapters()
	Example8_PolymorphicLogging()
	Example9_AdapterBenefits()
	Example10_WhenNotToUseAdapter()

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    EXAMPLES COMPLETED SUCCESSFULLY        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	fmt.Println("\n📝 Key Takeaways from Adapter Pattern:")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("1. ✓ Converts incompatible interfaces to work together")
	fmt.Println("2. ✓ Decouples client code from external library interfaces")
	fmt.Println("3. ✓ Enables polymorphic use of different implementations")
	fmt.Println("4. ✓ Isolates code from external API changes")
	fmt.Println("5. ✓ Improves testability through stable interfaces")
	fmt.Println("6. ✓ Go's interfaces make adapters lightweight and elegant")
	fmt.Println("\n🧪 To run the test suite:")
	fmt.Println("────────────────────────")
	fmt.Println("   go test -v ./tier2/adapter/")
	fmt.Println("   go test -cover ./tier2/adapter/")
	fmt.Println("   go test -bench=. ./tier2/adapter/")
}
