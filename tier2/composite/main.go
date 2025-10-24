package composite

import "fmt"

// RunAllExamples executes all Composite Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Composite pattern.
//
// Usage: Call from an external main program or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         Composite Pattern - Comprehensive Examples            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("The Composite pattern composes objects into tree structures")
	fmt.Println("to represent part-whole hierarchies. It lets clients treat")
	fmt.Println("individual objects and compositions uniformly.")
	fmt.Println()
	fmt.Println("This demonstration covers:")
	fmt.Println("  • File systems (files and directories)")
	fmt.Println("  • Organization charts (employees and managers)")
	fmt.Println("  • UI component trees (simple and container components)")
	fmt.Println("  • Menu systems (items and submenus)")
	fmt.Println("  • Complex document structures")
	fmt.Println()

	// Run all examples
	Example1_FileSystem()
	Example2_OrganizationChart()
	Example3_UIComponents()
	Example4_MenuSystem()
	Example5_RealWorld()

	// Summary
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Key Takeaways                               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✓ Uniform interface for leaves and composites")
	fmt.Println("✓ Recursive operations on tree structures")
	fmt.Println("✓ Easy to add new component types")
	fmt.Println("✓ Client simplicity - no type checking needed")
	fmt.Println("✓ Natural representation of hierarchical data")
	fmt.Println()
	fmt.Println("When to use:")
	fmt.Println("  • Tree or hierarchical structures (files, org charts, UI)")
	fmt.Println("  • Operations that work uniformly on items and collections")
	fmt.Println("  • Building complex objects from simpler components")
	fmt.Println("  • Part-whole hierarchies where parts contain other parts")
	fmt.Println("  • Recursive data structures (ASTs, DOM trees, menus)")
	fmt.Println()
	fmt.Println("When NOT to use:")
	fmt.Println("  • Flat structures (adds unnecessary complexity)")
	fmt.Println("  • Performance-critical code (tree traversal overhead)")
	fmt.Println("  • Simple parent-child relationships")
	fmt.Println("  • When leaves and composites need very different behavior")
	fmt.Println()
}
