package decorator

import (
	"fmt"
	"strings"
)

// RunAllExamples executes all Decorator Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Decorator pattern.
//
// Usage: Call from your own main package or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     DECORATOR PATTERN EXAMPLES                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples
	examples := []struct {
		name string
		fn   func()
	}{
		{"Example 1: Basic Decoration", Example1_BasicDecoration},
		{"Example 2: Stacking Multiple Decorators", Example2_StackingDecorators},
		{"Example 3: Data Pipeline with Compression", Example3_DataPipelineWithCompression},
		{"Example 4: Caching Decorator", Example4_CachingDecorator},
		{"Example 5: Encryption Decorator", Example5_EncryptionDecorator},
		{"Example 6: HTTP Middleware Chain", Example6_HTTPMiddlewareChain},
		{"Example 7: Decorator Order Matters", Example7_DecoratorOrderMatters},
		{"Example 8: Dynamic Decoration", Example8_DynamicDecoration},
	}

	for i, example := range examples {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("─", 76) + "\n")
		}
		example.fn()
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                      KEY TAKEAWAYS                                     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("1. ✓ Decorators add behavior without modifying the original object")
	fmt.Println("2. ✓ Multiple decorators can be stacked to compose complex behaviors")
	fmt.Println("3. ✓ Decorator order matters - it affects both behavior and performance")
	fmt.Println("4. ✓ Go's interfaces make decorators particularly elegant")
	fmt.Println("5. ✓ Functional decorators (middleware) are idiomatic in Go")
	fmt.Println("6. ✓ Use for cross-cutting concerns: logging, caching, validation, etc.")
	fmt.Println()
	fmt.Println("🧪 To run the test suite:")
	fmt.Println("────────────────────────")
	fmt.Println("   go test -v ./tier1/decorator/")
	fmt.Println("   go test -cover ./tier1/decorator/")
	fmt.Println("   go test -bench=. ./tier1/decorator/")
}
