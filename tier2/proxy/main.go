package proxy

import "fmt"

// RunAllExamples executes all Proxy Pattern examples in sequence.
// This demonstrates all use cases and benefits of the Proxy pattern.
//
// Usage: Call from cmd/main.go or in tests to see all examples in action.
func RunAllExamples() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║           PROXY PATTERN - COMPREHENSIVE EXAMPLES          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Run all examples in sequence
	Example1_WithoutLazyLoading()
	Example2_WithLazyLoading()
	Example3_MultipleAccesses()
	Example4_AccessControlAuthorized()
	Example5_AccessControlDenied()
	Example6_CachingProxyFirstAccess()
	Example7_CachingProxySubsequentAccess()
	Example8_CachingProxyComputations()
	Example9_LoggingProxy()
	Example10_RateLimiting()
	Example11_ProxyBenefits()
	Example12_WhenNotToUseProxy()
	Example13_ProxyVsDecorator()

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    EXAMPLES COMPLETED SUCCESSFULLY        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	fmt.Println("\n📝 Key Takeaways from Proxy Pattern:")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("1. ✓ Controls access to expensive or sensitive objects")
	fmt.Println("2. ✓ Lazy loading defers creation until needed")
	fmt.Println("3. ✓ Caching dramatically improves performance")
	fmt.Println("4. ✓ Access control enforces permissions")
	fmt.Println("5. ✓ Logging and rate limiting are natural proxy uses")
	fmt.Println("6. ✓ Proxy maintains same interface as real object")
	fmt.Println("\n🧪 To run the test suite:")
	fmt.Println("────────────────────────")
	fmt.Println("   go test -v ./tier2/proxy/")
	fmt.Println("   go test -cover ./tier2/proxy/")
	fmt.Println("   go test -bench=. ./tier2/proxy/")
}
