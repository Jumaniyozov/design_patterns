package singleton

import "fmt"

// RunAllExamples executes all Singleton Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Singleton pattern.
//
// Usage: Call from cmd/main.go or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         Singleton Pattern - Comprehensive Examples            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("The Singleton pattern ensures a class has only one instance")
	fmt.Println("while providing a global point of access to it.")
	fmt.Println()
	fmt.Println("This demonstration covers:")
	fmt.Println("  • Lazy initialization with sync.Once")
	fmt.Println("  • Eager initialization with init()")
	fmt.Println("  • Thread-safe concurrent access")
	fmt.Println("  • Reset capability for testing")
	fmt.Println("  • Real-world use cases")
	fmt.Println()

	// Run all examples
	Example1_LazyDatabaseSingleton()

	Example2_ConcurrentAccess()

	Example3_EagerConfiguration()

	Example4_LoggerWithReset()

	Example5_CacheSingleton()

	Example6_RealWorldDatabasePool()

	Example7_ConcurrentCacheAccess()

	// Summary
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Key Takeaways                               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✓ Use sync.Once for thread-safe lazy initialization")
	fmt.Println("✓ Use init() for eager initialization at package load time")
	fmt.Println("✓ Singletons provide global access to shared resources")
	fmt.Println("✓ Consider reset methods for testing scenarios")
	fmt.Println("✓ Thread safety is critical in concurrent environments")
	fmt.Println()
	fmt.Println("When to use:")
	fmt.Println("  • Database connection pools")
	fmt.Println("  • Application configuration")
	fmt.Println("  • Logging systems")
	fmt.Println("  • Caching layers")
	fmt.Println("  • Resource managers")
	fmt.Println()
	fmt.Println("When NOT to use:")
	fmt.Println("  • When dependency injection is more appropriate")
	fmt.Println("  • When you need multiple instances with different state")
	fmt.Println("  • For stateless utility objects (use functions instead)")
	fmt.Println("  • When singletons become performance bottlenecks")
	fmt.Println()
}