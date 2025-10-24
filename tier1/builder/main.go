package builder

import "fmt"

// RunAllExamples executes all Builder Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Builder pattern.
//
// Usage: Call from cmd/main.go or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║          BUILDER PATTERN - COMPREHENSIVE EXAMPLES         ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples in sequence
	Example1_BasicDatabasePool()
	Example2_ProductionDatabasePool()
	Example3_ValidationErrors()
	Example4_CrossFieldValidation()
	Example5_HTTPRequestBuilder()
	Example6_HTTPRequestWithMultipleHeaders()
	Example7_ApplicationConfig()
	Example8_DevelopmentVsProduction()
	Example9_PartialConfiguration()
	Example10_BuilderReuse()

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    EXAMPLES COMPLETED SUCCESSFULLY        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	fmt.Println("\n📝 Key Takeaways from Builder Pattern:")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("1. ✓ Builder pattern creates complex objects step-by-step")
	fmt.Println("2. ✓ Fluent API (method chaining) improves readability")
	fmt.Println("3. ✓ Validation can happen at each step or at Build() time")
	fmt.Println("4. ✓ Sensible defaults reduce boilerplate for common cases")
	fmt.Println("5. ✓ Builders should be used once and then discarded")
	fmt.Println("6. ✓ Perfect for objects with many optional fields")
	fmt.Println("\n🧪 To run the test suite:")
	fmt.Println("────────────────────────")
	fmt.Println("   go test -v ./tier1/builder/")
	fmt.Println("   go test -cover ./tier1/builder/")
	fmt.Println("   go test -bench=. ./tier1/builder/")
}
