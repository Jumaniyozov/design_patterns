package singleton

import (
	"fmt"
	"sync"
)

// Example1_ConfigManager demonstrates the singleton configuration manager.
func Example1_ConfigManager() {
	fmt.Println("\n=== Example 1: Configuration Manager ===")

	// Get the singleton instance
	config1 := GetConfigManager()
	fmt.Println("\nFirst instance - Reading default config:")
	if val, ok := config1.Get("app_name"); ok {
		fmt.Printf("  app_name: %s\n", val)
	}
	if val, ok := config1.Get("version"); ok {
		fmt.Printf("  version: %s\n", val)
	}

	// Set a new value
	config1.Set("api_key", "secret-key-123")
	fmt.Println("\nSet api_key on first instance")

	// Get the instance again from another part of the code
	config2 := GetConfigManager()
	fmt.Println("\nSecond instance - Verify it's the same instance:")

	// Verify both references point to the same instance
	if val, ok := config2.Get("api_key"); ok {
		fmt.Printf("  api_key from second instance: %s\n", val)
	}

	// Demonstrate they're the same instance
	fmt.Printf("\nAre they the same instance? %v\n", config1 == config2)

	// Show all settings
	fmt.Println("\nAll settings:")
	for k, v := range config2.GetAll() {
		fmt.Printf("  %s: %s\n", k, v)
	}
}

// Example2_DatabasePool demonstrates the singleton database connection pool.
func Example2_DatabasePool() {
	fmt.Println("\n=== Example 2: Database Connection Pool ===")

	// Get the singleton pool
	pool1 := GetDatabasePool()
	fmt.Printf("\nPool initialized with max size: %d\n", pool1.maxSize)

	// Acquire connections
	fmt.Println("\nAcquiring connections:")
	conn1 := pool1.AcquireConnection()
	conn2 := pool1.AcquireConnection()
	fmt.Printf("Acquired: %s, %s\n", conn1, conn2)

	// Get pool from another part of the code
	pool2 := GetDatabasePool()

	// Release connections back to pool
	fmt.Println("\nReleasing connections:")
	pool2.ReleaseConnection(conn1)
	pool2.ReleaseConnection(conn2)

	fmt.Printf("\nCurrent pool size: %d\n", pool2.GetPoolSize())

	// Acquire again - should reuse from pool
	fmt.Println("\nAcquiring connection again (should reuse from pool):")
	conn3 := pool1.AcquireConnection()
	fmt.Printf("Acquired: %s\n", conn3)

	fmt.Printf("\nAre pool references the same? %v\n", pool1 == pool2)
}

// Example3_Logger demonstrates the eagerly initialized singleton logger.
func Example3_Logger() {
	fmt.Println("\n=== Example 3: Logger Singleton ===")

	// Get logger instances from different parts of the application
	logger1 := GetLogger()
	logger2 := GetLogger()

	// Both should be the same instance
	fmt.Printf("\nAre logger instances the same? %v\n", logger1 == logger2)

	// Use logger
	fmt.Println("\nLogging with default prefix:")
	logger1.Log("Application started")
	logger1.Log("Processing request")

	// Change prefix from one instance
	fmt.Println("\nChanging prefix from logger1:")
	logger1.SetPrefix("[PRODUCTION]")

	// Log from second instance - should use updated prefix
	fmt.Println("Logging from logger2 (same instance):")
	logger2.Log("Configuration loaded")
	logger2.Log("System ready")
}

// Example4_ConcurrentAccess demonstrates thread-safe singleton initialization.
func Example4_ConcurrentAccess() {
	fmt.Println("\n=== Example 4: Concurrent Access ===")

	var wg sync.WaitGroup
	numGoroutines := 10

	fmt.Printf("\nLaunching %d goroutines to access ConfigManager...\n", numGoroutines)

	// Multiple goroutines trying to get the singleton concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			config := GetConfigManager()
			config.Set(fmt.Sprintf("key_%d", id), fmt.Sprintf("value_%d", id))
			fmt.Printf("Goroutine %d: got singleton instance\n", id)
		}(i)
	}

	wg.Wait()

	// Verify singleton instance has all the values
	config := GetConfigManager()
	fmt.Println("\nFinal configuration after concurrent access:")
	allSettings := config.GetAll()
	count := 0
	for k := range allSettings {
		if k[:4] == "key_" {
			count++
		}
	}
	fmt.Printf("Total keys set by goroutines: %d\n", count)
	fmt.Println("Singleton successfully handled concurrent access!")
}

// Example5_ComparisonOfPatterns demonstrates different singleton initialization patterns.
func Example5_ComparisonOfPatterns() {
	fmt.Println("\n=== Example 5: Comparison of Singleton Patterns ===")

	fmt.Println("\n1. Lazy Initialization with sync.Once (ConfigManager, DatabasePool):")
	fmt.Println("   - Instance created only when first accessed")
	fmt.Println("   - Thread-safe initialization guaranteed")
	fmt.Println("   - Best for expensive initialization")

	fmt.Println("\n2. Eager Initialization (Logger):")
	fmt.Println("   - Instance created at package load time")
	fmt.Println("   - Simple and fast access")
	fmt.Println("   - Best for lightweight objects always needed")

	fmt.Println("\n3. Trade-offs:")
	fmt.Println("   Lazy Init:")
	fmt.Println("     ✓ Saves resources if singleton never used")
	fmt.Println("     ✓ Can defer expensive initialization")
	fmt.Println("     ✗ Slightly more complex code")
	fmt.Println("   Eager Init:")
	fmt.Println("     ✓ Simpler code")
	fmt.Println("     ✓ Faster access (no sync overhead)")
	fmt.Println("     ✗ Always consumes resources even if unused")

	// Demonstrate both
	fmt.Println("\nDemonstration:")
	fmt.Println("Lazy - ConfigManager (watch for initialization message):")
	config := GetConfigManager()
	config.Set("demo", "lazy")

	fmt.Println("\nEager - Logger (no initialization message, already created):")
	logger := GetLogger()
	logger.Log("Eager initialization complete")
}
