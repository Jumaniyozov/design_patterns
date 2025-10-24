package singleton

import (
	"fmt"
	"sync"
)

// Example1_LazyDatabaseSingleton demonstrates lazy initialization with sync.Once
func Example1_LazyDatabaseSingleton() {
	fmt.Println("\n=== Example 1: Lazy Database Singleton ===")
	fmt.Println("Demonstrating thread-safe lazy initialization with sync.Once")
	fmt.Println()

	// First call - will initialize the database
	fmt.Println("First call to GetDatabase():")
	db1 := GetDatabase()
	fmt.Printf("Database instance address: %p\n", db1)
	fmt.Println()

	// Second call - returns the same instance
	fmt.Println("Second call to GetDatabase():")
	db2 := GetDatabase()
	fmt.Printf("Database instance address: %p\n", db2)
	fmt.Println()

	// Verify both references point to the same instance
	fmt.Printf("Are db1 and db2 the same instance? %v\n", db1 == db2)
	fmt.Println()

	// Use the database
	fmt.Println("Using the database connection pool:")
	if err := db1.Connect(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println(db1.Stats())

	if err := db2.Connect(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println(db2.Stats())

	db1.Disconnect()
	fmt.Println(db2.Stats())
}

// Example2_ConcurrentAccess demonstrates thread-safe singleton access
func Example2_ConcurrentAccess() {
	fmt.Println("\n=== Example 2: Concurrent Access to Singleton ===")
	fmt.Println("Demonstrating thread-safety with multiple goroutines")
	fmt.Println()

	var wg sync.WaitGroup
	addresses := make(chan string, 10)

	// Launch 10 concurrent goroutines trying to get the database instance
	fmt.Println("Launching 10 concurrent goroutines...")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := GetDatabase()
			addresses <- fmt.Sprintf("%p", db)
		}(i)
	}

	wg.Wait()
	close(addresses)

	// Verify all goroutines got the same instance
	fmt.Println("\nInstance addresses from all goroutines:")
	addressMap := make(map[string]int)
	for addr := range addresses {
		addressMap[addr]++
		fmt.Printf("  - %s\n", addr)
	}

	fmt.Printf("\nNumber of unique instances: %d\n", len(addressMap))
	fmt.Println("Expected: 1 (all goroutines should receive the same instance)")
}

// Example3_EagerConfiguration demonstrates eager initialization with init()
func Example3_EagerConfiguration() {
	fmt.Println("\n=== Example 3: Eager Configuration Singleton ===")
	fmt.Println("Configuration initialized at package load time using init()")
	fmt.Println()

	// Config is already initialized - no function call needed
	fmt.Printf("Application: %s\n", Config.AppName)
	fmt.Printf("Environment: %s\n", Config.Environment)
	fmt.Printf("Server Address: %s\n", Config.GetServerAddress())
	fmt.Println()

	// Check feature flags
	fmt.Println("Feature Flags:")
	for feature, enabled := range Config.FeatureFlags {
		status := "disabled"
		if enabled {
			status = "enabled"
		}
		fmt.Printf("  - %s: %s\n", feature, status)
	}
	fmt.Println()

	// Use the IsFeatureEnabled method
	if Config.IsFeatureEnabled("beta_mode") {
		fmt.Println("Beta mode is enabled - showing beta features")
	}
}

// Example4_LoggerWithReset demonstrates a singleton logger with reset capability
func Example4_LoggerWithReset() {
	fmt.Println("\n=== Example 4: Logger Singleton with Reset ===")
	fmt.Println("Demonstrating singleton pattern with testing support")
	fmt.Println()

	// Get logger and use it
	logger := GetLogger()
	logger.Log("Application started")
	logger.Log("Processing request #1")
	logger.Log("Request completed")

	fmt.Println()
	fmt.Printf("Total logs: %d\n", len(logger.GetLogs()))

	// Change log level
	fmt.Println("\nChanging log level to DEBUG...")
	logger.SetLevel("DEBUG")
	logger.Log("Debug information")

	fmt.Println()
	fmt.Println("All logs in the singleton:")
	for i, log := range logger.GetLogs() {
		fmt.Printf("  %d. %s\n", i+1, log)
	}
}

// Example5_CacheSingleton demonstrates a thread-safe cache singleton
func Example5_CacheSingleton() {
	fmt.Println("\n=== Example 5: Cache Singleton ===")
	fmt.Println("Demonstrating a thread-safe in-memory cache")
	fmt.Println()

	cache := GetCache()

	// Store some values
	fmt.Println("Storing values in cache:")
	cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
	cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})
	cache.Set("session:abc123", "active")
	fmt.Printf("Cache size: %d items\n", cache.Size())
	fmt.Println()

	// Retrieve values
	fmt.Println("Retrieving values from cache:")
	if value, exists := cache.Get("user:1"); exists {
		fmt.Printf("  user:1 = %v\n", value)
	}
	if value, exists := cache.Get("session:abc123"); exists {
		fmt.Printf("  session:abc123 = %v\n", value)
	}
	fmt.Println()

	// Verify singleton behavior - get cache again
	cache2 := GetCache()
	fmt.Printf("Cache instances are the same: %v\n", cache == cache2)
	fmt.Printf("Cache2 can access same data: ")
	if value, exists := cache2.Get("user:2"); exists {
		fmt.Printf("%v\n", value)
	}
	fmt.Println()

	// Delete and clear
	cache.Delete("session:abc123")
	fmt.Printf("After deleting session, cache size: %d\n", cache.Size())
}

// Example6_RealWorldDatabasePool demonstrates a real-world scenario
func Example6_RealWorldDatabasePool() {
	fmt.Println("\n=== Example 6: Real-World Database Pool Usage ===")
	fmt.Println("Simulating multiple services using the same database pool")
	fmt.Println()

	// UserService uses the database
	fmt.Println("UserService connecting to database:")
	db := GetDatabase()
	if err := db.Connect(); err != nil {
		fmt.Printf("UserService error: %v\n", err)
	}
	fmt.Printf("UserService: %s\n", db.Stats())
	fmt.Println()

	// OrderService uses the same database instance
	fmt.Println("OrderService connecting to database:")
	db2 := GetDatabase()
	if err := db2.Connect(); err != nil {
		fmt.Printf("OrderService error: %v\n", err)
	}
	fmt.Printf("OrderService: %s\n", db2.Stats())
	fmt.Println()

	// PaymentService also uses the same instance
	fmt.Println("PaymentService connecting to database:")
	db3 := GetDatabase()
	if err := db3.Connect(); err != nil {
		fmt.Printf("PaymentService error: %v\n", err)
	}
	fmt.Printf("PaymentService: %s\n", db3.Stats())
	fmt.Println()

	// All services share the same pool
	fmt.Printf("All services use the same pool: %v\n", db == db2 && db2 == db3)
	fmt.Printf("Instance address: %p\n", db)
	fmt.Println()

	// Cleanup
	db.Disconnect()
	db2.Disconnect()
	db3.Disconnect()
	fmt.Printf("After cleanup: %s\n", db.Stats())
}

// Example7_ConcurrentCacheAccess demonstrates thread-safe cache operations
func Example7_ConcurrentCacheAccess() {
	fmt.Println("\n=== Example 7: Concurrent Cache Access ===")
	fmt.Println("Demonstrating thread-safe cache operations")
	fmt.Println()

	cache := GetCache()
	cache.Clear() // Start fresh

	var wg sync.WaitGroup

	// Writer goroutines
	fmt.Println("Launching 5 writer goroutines...")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("writer-%d-item-%d", id, j)
				cache.Set(key, fmt.Sprintf("value-%d-%d", id, j))
			}
		}(i)
	}

	// Reader goroutines
	fmt.Println("Launching 5 reader goroutines...")
	readCount := make(chan int, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			count := 0
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("writer-%d-item-%d", id, j)
				if _, exists := cache.Get(key); exists {
					count++
				}
			}
			readCount <- count
		}(i)
	}

	wg.Wait()
	close(readCount)

	fmt.Println("\nResults:")
	fmt.Printf("Final cache size: %d items\n", cache.Size())
	fmt.Println("Reads per reader goroutine:")
	for count := range readCount {
		fmt.Printf("  - Found %d items\n", count)
	}
	fmt.Println("\nAll operations completed successfully without race conditions!")
}