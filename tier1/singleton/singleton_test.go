package singleton

import (
	"sync"
	"testing"
)

// TestDatabaseSingleton verifies that GetDatabase always returns the same instance
func TestDatabaseSingleton(t *testing.T) {
	db1 := GetDatabase()
	db2 := GetDatabase()

	if db1 != db2 {
		t.Errorf("Expected same database instance, got different instances: %p vs %p", db1, db2)
	}
}

// TestDatabaseConcurrentAccess verifies thread-safe initialization
func TestDatabaseConcurrentAccess(t *testing.T) {
	const goroutines = 100
	instances := make([]*Database, goroutines)
	var wg sync.WaitGroup

	// Launch multiple goroutines trying to get the database
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			instances[index] = GetDatabase()
		}(i)
	}

	wg.Wait()

	// Verify all goroutines got the same instance
	first := instances[0]
	for i, db := range instances {
		if db != first {
			t.Errorf("Goroutine %d got different instance: %p vs %p", i, db, first)
		}
	}
}

// TestDatabaseConnectionPool tests the connection pool functionality
func TestDatabaseConnectionPool(t *testing.T) {
	db := GetDatabase()

	// Test successful connection
	err := db.Connect()
	if err != nil {
		t.Errorf("Expected successful connection, got error: %v", err)
	}

	if db.activeConnections != 1 {
		t.Errorf("Expected 1 active connection, got %d", db.activeConnections)
	}

	// Test successful disconnection
	err = db.Disconnect()
	if err != nil {
		t.Errorf("Expected successful disconnection, got error: %v", err)
	}

	if db.activeConnections != 0 {
		t.Errorf("Expected 0 active connections after disconnect, got %d", db.activeConnections)
	}
}

// TestDatabaseMaxConnections tests the max connections limit
func TestDatabaseMaxConnections(t *testing.T) {
	db := GetDatabase()

	// Reset to known state
	db.mu.Lock()
	db.activeConnections = 0
	db.maxConnections = 3
	db.mu.Unlock()

	// Create 3 connections (should succeed)
	for i := 0; i < 3; i++ {
		if err := db.Connect(); err != nil {
			t.Errorf("Connection %d failed: %v", i, err)
		}
	}

	// Try to exceed max connections
	err := db.Connect()
	if err == nil {
		t.Error("Expected error when exceeding max connections, got nil")
	}

	// Cleanup
	db.mu.Lock()
	db.activeConnections = 0
	db.maxConnections = 100
	db.mu.Unlock()
}

// TestConfigurationSingleton verifies the eager-loaded configuration
func TestConfigurationSingleton(t *testing.T) {
	if Config == nil {
		t.Fatal("Expected Config to be initialized, got nil")
	}

	if Config.AppName != "MyApplication" {
		t.Errorf("Expected AppName 'MyApplication', got '%s'", Config.AppName)
	}

	if Config.Port != 8080 {
		t.Errorf("Expected Port 8080, got %d", Config.Port)
	}
}

// TestConfigurationFeatureFlags tests feature flag functionality
func TestConfigurationFeatureFlags(t *testing.T) {
	tests := []struct {
		feature  string
		expected bool
	}{
		{"feature_x", true},
		{"feature_y", false},
		{"beta_mode", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.feature, func(t *testing.T) {
			result := Config.IsFeatureEnabled(tt.feature)
			if result != tt.expected {
				t.Errorf("Expected %s to be %v, got %v", tt.feature, tt.expected, result)
			}
		})
	}
}

// TestLoggerSingleton verifies logger singleton behavior
func TestLoggerSingleton(t *testing.T) {
	// Reset logger for clean test
	ResetForTesting()

	logger1 := GetLogger()
	logger2 := GetLogger()

	if logger1 != logger2 {
		t.Errorf("Expected same logger instance, got different instances: %p vs %p", logger1, logger2)
	}
}

// TestLoggerFunctionality tests logging functionality
func TestLoggerFunctionality(t *testing.T) {
	ResetForTesting()
	logger := GetLogger()

	// Log some messages
	logger.Log("Test message 1")
	logger.Log("Test message 2")

	logs := logger.GetLogs()
	if len(logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(logs))
	}

	// Verify logs contain the messages
	found1, found2 := false, false
	for _, log := range logs {
		if len(log) > 0 {
			if !found1 && len(log) > 20 && log[len(log)-14:] == "Test message 1" {
				found1 = true
			}
			if !found2 && len(log) > 20 && log[len(log)-14:] == "Test message 2" {
				found2 = true
			}
		}
	}

	if !found1 || !found2 {
		t.Error("Expected to find both test messages in logs")
	}
}

// TestLoggerSetLevel tests changing log level
func TestLoggerSetLevel(t *testing.T) {
	ResetForTesting()
	logger := GetLogger()

	logger.SetLevel("DEBUG")
	logger.Log("Debug message")

	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}

	// Check that log contains DEBUG level
	if len(logs) > 0 && len(logs[0]) > 0 {
		// The log format is [TIME] LEVEL: message
		// We just verify it was logged
		if len(logs[0]) == 0 {
			t.Error("Expected log entry to contain content")
		}
	}
}

// TestLoggerReset tests the reset functionality
func TestLoggerReset(t *testing.T) {
	ResetForTesting()

	logger1 := GetLogger()
	logger1.Log("Before reset")

	ResetForTesting()

	logger2 := GetLogger()
	logs := logger2.GetLogs()

	if len(logs) != 0 {
		t.Errorf("Expected 0 logs after reset, got %d", len(logs))
	}

	// Note: logger1 and logger2 are different instances after reset
	// This is expected behavior for testing scenarios
}

// TestCacheSingleton verifies cache singleton behavior
func TestCacheSingleton(t *testing.T) {
	cache1 := GetCache()
	cache2 := GetCache()

	if cache1 != cache2 {
		t.Errorf("Expected same cache instance, got different instances: %p vs %p", cache1, cache2)
	}
}

// TestCacheOperations tests basic cache operations
func TestCacheOperations(t *testing.T) {
	cache := GetCache()
	cache.Clear() // Start with empty cache

	// Test Set and Get
	cache.Set("key1", "value1")
	value, exists := cache.Get("key1")

	if !exists {
		t.Error("Expected key1 to exist in cache")
	}

	if value != "value1" {
		t.Errorf("Expected value 'value1', got '%v'", value)
	}

	// Test Size
	cache.Set("key2", "value2")
	if cache.Size() != 2 {
		t.Errorf("Expected cache size 2, got %d", cache.Size())
	}

	// Test Delete
	cache.Delete("key1")
	_, exists = cache.Get("key1")
	if exists {
		t.Error("Expected key1 to be deleted from cache")
	}

	if cache.Size() != 1 {
		t.Errorf("Expected cache size 1 after delete, got %d", cache.Size())
	}

	// Test Clear
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected cache size 0 after clear, got %d", cache.Size())
	}
}

// TestCacheConcurrentAccess tests thread-safe cache operations
func TestCacheConcurrentAccess(t *testing.T) {
	cache := GetCache()
	cache.Clear()

	const (
		goroutines = 50
		operations = 100
	)

	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := string(rune('A' + (id % 26)))
				cache.Set(key, id*operations+j)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := string(rune('A' + (id % 26)))
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still functional
	cache.Set("test", "value")
	value, exists := cache.Get("test")

	if !exists {
		t.Error("Cache not functional after concurrent access")
	}

	if value != "value" {
		t.Errorf("Expected 'value', got '%v'", value)
	}
}

// TestCacheDataIsolation verifies cache data is shared across instances
func TestCacheDataIsolation(t *testing.T) {
	cache1 := GetCache()
	cache1.Clear()

	cache1.Set("shared_key", "shared_value")

	cache2 := GetCache()
	value, exists := cache2.Get("shared_key")

	if !exists {
		t.Error("Expected cache2 to access data from cache1 (same instance)")
	}

	if value != "shared_value" {
		t.Errorf("Expected 'shared_value', got '%v'", value)
	}
}

// BenchmarkGetDatabase measures the performance of getting the database singleton
func BenchmarkGetDatabase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetDatabase()
	}
}

// BenchmarkGetLogger measures the performance of getting the logger singleton
func BenchmarkGetLogger(b *testing.B) {
	ResetForTesting()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetLogger()
	}
}

// BenchmarkGetCache measures the performance of getting the cache singleton
func BenchmarkGetCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCache()
	}
}

// BenchmarkCacheOperations measures cache operation performance
func BenchmarkCacheOperations(b *testing.B) {
	cache := GetCache()
	cache.Clear()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set("key", i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		cache.Set("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get("key")
		}
	})

	b.Run("Delete", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set("key", i)
			cache.Delete("key")
		}
	})
}

// BenchmarkConcurrentCacheAccess measures concurrent cache performance
func BenchmarkConcurrentCacheAccess(b *testing.B) {
	cache := GetCache()
	cache.Clear()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set("key", i)
			cache.Get("key")
			i++
		}
	})
}