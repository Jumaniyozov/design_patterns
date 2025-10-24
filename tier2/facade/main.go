package facade

import "fmt"

// RunAllExamples executes all Facade Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Facade pattern.
//
// Usage: Call from an external main program or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Facade Pattern - Comprehensive Examples              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("The Facade pattern provides a simplified, unified interface to")
	fmt.Println("a complex subsystem. It hides complexity by creating a single")
	fmt.Println("entry point that orchestrates multiple components.")
	fmt.Println()
	fmt.Println("This demonstration covers:")
	fmt.Println("  • Banking system (account, transaction, fraud, notification)")
	fmt.Println("  • Smart home automation (lighting, security, climate, entertainment)")
	fmt.Println("  • E-commerce order processing (inventory, payment, shipping, email)")
	fmt.Println("  • Comparison: with vs without facade")
	fmt.Println("  • Real-world video encoding system")
	fmt.Println("  • Pattern composition and integration")
	fmt.Println()

	// Run all examples
	Example1_BankingSystem()
	Example2_SmartHome()
	Example3_ECommerce()
	Example4_ComparisonWithoutFacade()
	Example5_RealWorld()
	Example6_FacadePatternIntegration()

	// Summary
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Key Takeaways                               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✓ Simplifies complex subsystems with unified interface")
	fmt.Println("✓ Decouples clients from implementation details")
	fmt.Println("✓ Provides clear architectural boundaries")
	fmt.Println("✓ Makes systems easier to use and understand")
	fmt.Println("✓ Changes to subsystems don't affect clients if facade stays stable")
	fmt.Println()
	fmt.Println("When to use:")
	fmt.Println("  • Complex subsystems with multiple interconnected components")
	fmt.Println("  • Need to simplify interface to a set of classes")
	fmt.Println("  • Want to reduce dependencies on internal details")
	fmt.Println("  • Building layered architecture (presentation, business, data)")
	fmt.Println("  • Integrating third-party libraries with cohesive interface")
	fmt.Println()
	fmt.Println("When NOT to use:")
	fmt.Println("  • Simple systems that don't need simplification")
	fmt.Println("  • Hiding important details clients should understand")
	fmt.Println("  • Facade adds unnecessary indirection without value")
	fmt.Println("  • All clients need direct access to subsystem details")
	fmt.Println()
	fmt.Println("Real-world applications:")
	fmt.Println("  • Banking APIs (abstracting complex financial systems)")
	fmt.Println("  • Smart home automation (scene-based control)")
	fmt.Println("  • E-commerce platforms (order processing workflow)")
	fmt.Println("  • Video processing (encoding, transcoding, delivery)")
	fmt.Println("  • Database access layers (hiding query complexity)")
	fmt.Println()
}
