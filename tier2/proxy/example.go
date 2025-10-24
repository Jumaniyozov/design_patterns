package proxy

import (
	"fmt"
	"time"
)

// Example1_WithoutLazyLoading demonstrates the problem without lazy loading.
func Example1_WithoutLazyLoading() {
	fmt.Println("=== Example 1: Problem - Loading Without Lazy Proxy ===")
	fmt.Println()

	fmt.Println("❌ Loading 5 images immediately (even if never used):")
	start := time.Now()

	images := []*RealImage{
		NewRealImage("photo1.jpg"),
		NewRealImage("photo2.jpg"),
		NewRealImage("photo3.jpg"),
		NewRealImage("photo4.jpg"),
		NewRealImage("photo5.jpg"),
	}

	duration := time.Since(start)
	fmt.Printf("\nTotal load time: %v\n", duration)
	fmt.Printf("Memory used: %d bytes\n\n", len(images)*100*1024)
	fmt.Println("💡 All images loaded upfront, wasting time and memory!")
	fmt.Println("   What if we only need to display one image?")
}

// Example2_WithLazyLoading demonstrates the solution with lazy loading proxy.
func Example2_WithLazyLoading() {
	fmt.Println("\n=== Example 2: Solution - Lazy Loading with Proxy ===")
	fmt.Println()

	fmt.Println("✓ Creating 5 image proxies (instant!):")
	start := time.Now()

	images := []Image{
		NewImageProxy("photo1.jpg"),
		NewImageProxy("photo2.jpg"),
		NewImageProxy("photo3.jpg"),
		NewImageProxy("photo4.jpg"),
		NewImageProxy("photo5.jpg"),
	}

	creationTime := time.Since(start)
	fmt.Printf("\nProxy creation time: %v\n", creationTime)
	fmt.Println()

	// Only display one image
	fmt.Println("Now displaying just ONE image:")
	images[0].Display()

	fmt.Println()
	fmt.Println("💡 Only one image was loaded - saved 400KB memory and 400ms time!")
	fmt.Println("   Other images remain as lightweight proxies until accessed.")
}

// Example3_MultipleAccesses demonstrates that proxy caches the loaded image.
func Example3_MultipleAccesses() {
	fmt.Println("\n=== Example 3: Proxy Caches Loaded Image ===")
	fmt.Println()

	proxy := NewImageProxy("large_photo.jpg")

	fmt.Println("First access (loads image):")
	proxy.Display()

	fmt.Println("\nSecond access (uses cached image):")
	proxy.Display()

	fmt.Println("\nThird access (still cached):")
	proxy.Display()

	fmt.Println()
	fmt.Println("💡 Image loaded only once! Subsequent accesses are instant.")
}

// Example4_AccessControlAuthorized demonstrates authorized database access.
func Example4_AccessControlAuthorized() {
	fmt.Println("\n=== Example 4: Access Control - Authorized User ===")
	fmt.Println()

	// User with full permissions
	admin := &User{
		Name:        "admin",
		Permissions: []string{"read", "write", "delete"},
	}

	db := NewDatabaseProxy(admin, "production_db")

	fmt.Println("Admin attempting to query:")
	result, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}

	fmt.Println()

	fmt.Println("Admin attempting to execute:")
	rows, err := db.Execute("UPDATE users SET active = true")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Rows affected: %d\n", rows)
	}

	fmt.Println()
	fmt.Println("💡 Admin has permissions - all operations succeed.")
}

// Example5_AccessControlDenied demonstrates denied database access.
func Example5_AccessControlDenied() {
	fmt.Println("\n=== Example 5: Access Control - Unauthorized User ===")
	fmt.Println()

	// User with read-only permissions
	guest := &User{
		Name:        "guest",
		Permissions: []string{"read"},
	}

	db := NewDatabaseProxy(guest, "production_db")

	fmt.Println("Guest attempting to query (has read permission):")
	result, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✓ Result: %s\n", result)
	}

	fmt.Println()

	fmt.Println("Guest attempting to execute (lacks write permission):")
	_, err = db.Execute("DELETE FROM users WHERE id = 1")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✓ Success (unexpected!)")
	}

	fmt.Println()
	fmt.Println("💡 Access control proxy protected the database from")
	fmt.Println("   unauthorized modifications!")
}

// Example6_CachingProxyFirstAccess demonstrates cache miss scenario.
func Example6_CachingProxyFirstAccess() {
	fmt.Println("\n=== Example 6: Caching Proxy - First Access (Cache Miss) ===")
	fmt.Println()

	proxy := NewCachingProxy("DataService")

	fmt.Println("First request for 'user_profile':")
	start := time.Now()
	data, _ := proxy.GetData("user_profile")
	duration := time.Since(start)

	fmt.Printf("Result: %s\n", data)
	fmt.Printf("Time taken: %v\n", duration)

	fmt.Println()
	fmt.Println("💡 Cache miss - expensive operation executed.")
}

// Example7_CachingProxySubsequentAccess demonstrates cache hit scenario.
func Example7_CachingProxySubsequentAccess() {
	fmt.Println("\n=== Example 7: Caching Proxy - Subsequent Access (Cache Hit) ===")
	fmt.Println()

	proxy := NewCachingProxy("DataService")

	// Prime the cache
	fmt.Println("First request (priming cache):")
	proxy.GetData("user_profile")

	fmt.Println()

	// Hit the cache
	fmt.Println("Second request for same key:")
	start := time.Now()
	data, _ := proxy.GetData("user_profile")
	duration := time.Since(start)

	fmt.Printf("Result: %s\n", data)
	fmt.Printf("Time taken: %v (instant!)\n", duration)

	fmt.Println()
	fmt.Println("💡 Cache hit - no expensive operation, instant result!")
}

// Example8_CachingProxyComputations demonstrates caching expensive computations.
func Example8_CachingProxyComputations() {
	fmt.Println("\n=== Example 8: Caching Expensive Computations ===")
	fmt.Println()

	proxy := NewCachingProxy("ComputeService")

	inputs := []int{5, 10, 5, 10, 15, 5}

	fmt.Println("Processing sequence:", inputs)
	fmt.Println()

	for i, input := range inputs {
		fmt.Printf("Request %d: Computing cube of %d\n", i+1, input)
		result, _ := proxy.ComputeExpensive(input)
		fmt.Printf("Result: %d\n\n", result)
	}

	fmt.Println("💡 Notice: Inputs 5 and 10 were computed only once,")
	fmt.Println("   then retrieved from cache on subsequent requests!")
}

// Example9_LoggingProxy demonstrates logging all operations.
func Example9_LoggingProxy() {
	fmt.Println("\n=== Example 9: Logging Proxy ===")
	fmt.Println()

	proxy := NewLoggingPaymentProxy()

	fmt.Println("Processing multiple payments:")
	fmt.Println()

	// Successful payment
	proxy.ProcessPayment(99.99, "Alice")
	fmt.Println()

	// Another successful payment
	proxy.ProcessPayment(149.50, "Bob")
	fmt.Println()

	// Failed payment
	proxy.ProcessPayment(-50.00, "Charlie")
	fmt.Println()

	fmt.Printf("Total payment calls logged: %d\n", proxy.GetCallCount())

	fmt.Println()
	fmt.Println("💡 All operations automatically logged with timestamps,")
	fmt.Println("   duration, and outcome - perfect for auditing!")
}

// Example10_RateLimiting demonstrates rate limiting.
func Example10_RateLimiting() {
	fmt.Println("\n=== Example 10: Rate Limiting Proxy ===")
	fmt.Println()

	// Allow 3 requests per 2 seconds
	proxy := NewRateLimitingProxy(3, 2*time.Second)

	fmt.Println("Attempting 5 requests (limit: 3 per 2 seconds):")
	fmt.Println()

	endpoints := []string{"/api/users", "/api/posts", "/api/comments", "/api/likes", "/api/shares"}

	for i, endpoint := range endpoints {
		fmt.Printf("Request %d to %s:\n", i+1, endpoint)
		response, err := proxy.MakeRequest(endpoint)
		if err != nil {
			fmt.Printf("  ❌ %v\n", err)
		} else {
			fmt.Printf("  ✓ %s\n", response)
		}
		fmt.Println()
	}

	fmt.Println("💡 Rate limiting proxy protected the API from overload!")
	fmt.Println("   First 3 requests succeeded, remaining were throttled.")

	fmt.Println("\nWaiting 2 seconds for window to reset...")
	time.Sleep(2 * time.Second)

	fmt.Println("Trying again after window reset:")
	response, err := proxy.MakeRequest("/api/profile")
	if err != nil {
		fmt.Printf("❌ %v\n", err)
	} else {
		fmt.Printf("✓ %s\n", response)
	}
}

// Example11_ProxyBenefits demonstrates key benefits of proxy pattern.
func Example11_ProxyBenefits() {
	fmt.Println("\n=== Example 11: Key Benefits of Proxy Pattern ===")
	fmt.Println()

	fmt.Println("1. LAZY INITIALIZATION")
	fmt.Println("   ✓ Defer expensive object creation until needed")
	fmt.Println("   ✓ Save memory and startup time")
	fmt.Println()

	fmt.Println("2. ACCESS CONTROL")
	fmt.Println("   ✓ Enforce permissions before granting access")
	fmt.Println("   ✓ Protect sensitive resources")
	fmt.Println()

	fmt.Println("3. CACHING")
	fmt.Println("   ✓ Store expensive computation results")
	fmt.Println("   ✓ Dramatically improve performance")
	fmt.Println()

	fmt.Println("4. LOGGING & MONITORING")
	fmt.Println("   ✓ Track all access to objects")
	fmt.Println("   ✓ Audit operations automatically")
	fmt.Println()

	fmt.Println("5. RATE LIMITING")
	fmt.Println("   ✓ Protect resources from overload")
	fmt.Println("   ✓ Enforce usage quotas")
	fmt.Println()

	fmt.Println("6. TRANSPARENCY")
	fmt.Println("   ✓ Proxy and real object share same interface")
	fmt.Println("   ✓ Clients don't need to know about proxy")
}

// Example12_WhenNotToUseProxy demonstrates anti-patterns.
func Example12_WhenNotToUseProxy() {
	fmt.Println("\n=== Example 12: When NOT to Use Proxy Pattern ===")
	fmt.Println()

	fmt.Println("❌ AVOID proxies when:")
	fmt.Println()
	fmt.Println("1. Object creation is cheap")
	fmt.Println("   → Proxy overhead exceeds any benefit")
	fmt.Println()
	fmt.Println("2. Direct access is clearer")
	fmt.Println("   → Extra indirection obscures intent")
	fmt.Println()
	fmt.Println("3. Performance is critical")
	fmt.Println("   → Proxy adds measurable overhead")
	fmt.Println()
	fmt.Println("4. No access control, caching, or logging needed")
	fmt.Println("   → Proxy provides no value")
	fmt.Println()
	fmt.Println("5. Debugging is difficult")
	fmt.Println("   → Extra layer makes issues hard to trace")
	fmt.Println()
	fmt.Println("💡 Use proxies strategically when they solve real problems,")
	fmt.Println("   not reflexively for every object access.")
}

// Example13_ProxyVsDecorator demonstrates the difference.
func Example13_ProxyVsDecorator() {
	fmt.Println("\n=== Example 13: Proxy vs Decorator Pattern ===")
	fmt.Println()

	fmt.Println("PROXY Pattern:")
	fmt.Println("─────────────")
	fmt.Println("• Controls ACCESS to an object")
	fmt.Println("• Usually manages object lifecycle (lazy loading)")
	fmt.Println("• May not create real object until needed")
	fmt.Println("• Examples: caching, access control, lazy loading")
	fmt.Println()

	fmt.Println("DECORATOR Pattern:")
	fmt.Println("──────────────────")
	fmt.Println("• Adds BEHAVIOR to an object")
	fmt.Println("• Always wraps an existing object")
	fmt.Println("• Object exists when decorator is created")
	fmt.Println("• Examples: adding features, enhancing functionality")
	fmt.Println()

	fmt.Println("💡 Proxy = Control access / Decorator = Add behavior")
}
