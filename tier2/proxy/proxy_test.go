package proxy

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// TestImageProxy_LazyLoading tests that image is not loaded until accessed.
func TestImageProxy_LazyLoading(t *testing.T) {
	proxy := NewImageProxy("test.jpg")

	// At this point, real image should NOT be loaded
	if proxy.realImage != nil {
		t.Error("Expected realImage to be nil before first access")
	}

	// First access should load the image
	proxy.Display()

	// Now real image should be loaded
	if proxy.realImage == nil {
		t.Error("Expected realImage to be loaded after first access")
	}
}

// TestImageProxy_MultipleAccesses tests that image is loaded only once.
func TestImageProxy_MultipleAccesses(t *testing.T) {
	proxy := NewImageProxy("test.jpg")

	// First access
	proxy.Display()
	firstLoad := proxy.realImage

	// Second access
	proxy.Display()
	secondLoad := proxy.realImage

	// Should be the same instance
	if firstLoad != secondLoad {
		t.Error("Expected same realImage instance on multiple accesses")
	}
}

// TestImageProxy_GetSize tests lazy loading through GetSize.
func TestImageProxy_GetSize(t *testing.T) {
	proxy := NewImageProxy("test.jpg")

	if proxy.realImage != nil {
		t.Error("Expected realImage to be nil before access")
	}

	size := proxy.GetSize()

	if size != 100*1024 {
		t.Errorf("Expected size 102400, got %d", size)
	}

	if proxy.realImage == nil {
		t.Error("Expected realImage to be loaded after GetSize")
	}
}

// TestImageProxy_ThreadSafe tests concurrent access is safe.
func TestImageProxy_ThreadSafe(t *testing.T) {
	proxy := NewImageProxy("concurrent.jpg")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			proxy.Display()
		}()
	}

	wg.Wait()

	// Should only have one real image loaded
	if proxy.realImage == nil {
		t.Error("Expected realImage to be loaded")
	}
}

// TestDatabaseProxy_AuthorizedAccess tests access with proper permissions.
func TestDatabaseProxy_AuthorizedAccess(t *testing.T) {
	user := &User{
		Name:        "admin",
		Permissions: []string{"read", "write"},
	}

	db := NewDatabaseProxy(user, "testdb")

	// Test Query (requires read permission)
	_, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Expected successful query, got error: %v", err)
	}

	// Test Execute (requires write permission)
	rows, err := db.Execute("UPDATE users SET active = true")
	if err != nil {
		t.Errorf("Expected successful execute, got error: %v", err)
	}
	if rows != 3 {
		t.Errorf("Expected 3 rows affected, got %d", rows)
	}
}

// TestDatabaseProxy_UnauthorizedQuery tests denied query access.
func TestDatabaseProxy_UnauthorizedQuery(t *testing.T) {
	user := &User{
		Name:        "guest",
		Permissions: []string{}, // No permissions
	}

	db := NewDatabaseProxy(user, "testdb")

	_, err := db.Query("SELECT * FROM users")
	if err == nil {
		t.Error("Expected error for unauthorized query, got nil")
	}

	if !strings.Contains(err.Error(), "access denied") {
		t.Errorf("Expected 'access denied' error, got: %v", err)
	}
}

// TestDatabaseProxy_UnauthorizedExecute tests denied execute access.
func TestDatabaseProxy_UnauthorizedExecute(t *testing.T) {
	user := &User{
		Name:        "reader",
		Permissions: []string{"read"}, // Only read permission
	}

	db := NewDatabaseProxy(user, "testdb")

	_, err := db.Execute("DELETE FROM users")
	if err == nil {
		t.Error("Expected error for unauthorized execute, got nil")
	}

	if !strings.Contains(err.Error(), "access denied") {
		t.Errorf("Expected 'access denied' error, got: %v", err)
	}
}

// TestCachingProxy_CacheMiss tests first access (cache miss).
func TestCachingProxy_CacheMiss(t *testing.T) {
	proxy := NewCachingProxy("test_service")

	start := time.Now()
	result, err := proxy.GetData("test_key")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected successful request, got error: %v", err)
	}

	if !strings.Contains(result, "test_key") {
		t.Errorf("Expected result to contain 'test_key', got: %s", result)
	}

	// Should take at least 200ms (simulated delay)
	if duration < 150*time.Millisecond {
		t.Errorf("Expected delay >= 150ms for cache miss, got %v", duration)
	}
}

// TestCachingProxy_CacheHit tests second access (cache hit).
func TestCachingProxy_CacheHit(t *testing.T) {
	proxy := NewCachingProxy("test_service")

	// Prime the cache
	proxy.GetData("test_key")

	// Second access should be instant (from cache)
	start := time.Now()
	result, err := proxy.GetData("test_key")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected successful request, got error: %v", err)
	}

	if !strings.Contains(result, "test_key") {
		t.Errorf("Expected result to contain 'test_key', got: %s", result)
	}

	// Should be instant (< 10ms)
	if duration > 10*time.Millisecond {
		t.Errorf("Expected cache hit < 10ms, got %v", duration)
	}
}

// TestCachingProxy_ComputeExpensive tests caching of computations.
func TestCachingProxy_ComputeExpensive(t *testing.T) {
	proxy := NewCachingProxy("compute_service")

	// First computation (cache miss)
	start := time.Now()
	result1, _ := proxy.ComputeExpensive(5)
	duration1 := time.Since(start)

	if result1 != 125 { // 5^3 = 125
		t.Errorf("Expected 125, got %d", result1)
	}

	if duration1 < 250*time.Millisecond {
		t.Errorf("Expected delay >= 250ms for cache miss, got %v", duration1)
	}

	// Second computation with same input (cache hit)
	start = time.Now()
	result2, _ := proxy.ComputeExpensive(5)
	duration2 := time.Since(start)

	if result2 != 125 {
		t.Errorf("Expected 125, got %d", result2)
	}

	if duration2 > 10*time.Millisecond {
		t.Errorf("Expected cache hit < 10ms, got %v", duration2)
	}
}

// TestCachingProxy_DifferentKeys tests that different keys don't share cache.
func TestCachingProxy_DifferentKeys(t *testing.T) {
	proxy := NewCachingProxy("test_service")

	result1, _ := proxy.GetData("key1")
	result2, _ := proxy.GetData("key2")

	if result1 == result2 {
		t.Error("Expected different results for different keys")
	}

	if !strings.Contains(result1, "key1") {
		t.Errorf("Expected result1 to contain 'key1', got: %s", result1)
	}

	if !strings.Contains(result2, "key2") {
		t.Errorf("Expected result2 to contain 'key2', got: %s", result2)
	}
}

// TestLoggingPaymentProxy_Success tests logging of successful payments.
func TestLoggingPaymentProxy_Success(t *testing.T) {
	proxy := NewLoggingPaymentProxy()

	txnID, err := proxy.ProcessPayment(100.00, "Alice")

	if err != nil {
		t.Fatalf("Expected successful payment, got error: %v", err)
	}

	if txnID == "" {
		t.Error("Expected transaction ID, got empty string")
	}

	if !strings.HasPrefix(txnID, "TXN-") {
		t.Errorf("Expected transaction ID to start with 'TXN-', got: %s", txnID)
	}

	if proxy.GetCallCount() != 1 {
		t.Errorf("Expected call count 1, got %d", proxy.GetCallCount())
	}
}

// TestLoggingPaymentProxy_Failure tests logging of failed payments.
func TestLoggingPaymentProxy_Failure(t *testing.T) {
	proxy := NewLoggingPaymentProxy()

	_, err := proxy.ProcessPayment(-50.00, "Bob")

	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}

	if !strings.Contains(err.Error(), "invalid amount") {
		t.Errorf("Expected 'invalid amount' error, got: %v", err)
	}

	if proxy.GetCallCount() != 1 {
		t.Errorf("Expected call count 1, got %d", proxy.GetCallCount())
	}
}

// TestLoggingPaymentProxy_MultipleCallsCount tests call counting.
func TestLoggingPaymentProxy_MultipleCallsCount(t *testing.T) {
	proxy := NewLoggingPaymentProxy()

	proxy.ProcessPayment(50.00, "Alice")
	proxy.ProcessPayment(75.00, "Bob")
	proxy.ProcessPayment(100.00, "Charlie")

	if proxy.GetCallCount() != 3 {
		t.Errorf("Expected call count 3, got %d", proxy.GetCallCount())
	}
}

// TestRateLimitingProxy_AllowedRequests tests requests within limit.
func TestRateLimitingProxy_AllowedRequests(t *testing.T) {
	proxy := NewRateLimitingProxy(3, 1*time.Second)

	// All 3 requests should succeed
	for i := 0; i < 3; i++ {
		_, err := proxy.MakeRequest("/api/test")
		if err != nil {
			t.Errorf("Request %d should succeed, got error: %v", i+1, err)
		}
	}
}

// TestRateLimitingProxy_ExceededLimit tests requests exceeding limit.
func TestRateLimitingProxy_ExceededLimit(t *testing.T) {
	proxy := NewRateLimitingProxy(2, 1*time.Second)

	// First 2 should succeed
	proxy.MakeRequest("/api/test")
	proxy.MakeRequest("/api/test")

	// Third should fail
	_, err := proxy.MakeRequest("/api/test")
	if err == nil {
		t.Error("Expected error for rate limit exceeded, got nil")
	}

	if !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("Expected 'rate limit exceeded' error, got: %v", err)
	}
}

// TestRateLimitingProxy_WindowReset tests that window resets after time.
func TestRateLimitingProxy_WindowReset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping time-based test in short mode")
	}

	proxy := NewRateLimitingProxy(2, 500*time.Millisecond)

	// Fill the limit
	proxy.MakeRequest("/api/test")
	proxy.MakeRequest("/api/test")

	// Third should fail
	_, err := proxy.MakeRequest("/api/test")
	if err == nil {
		t.Error("Expected error for rate limit exceeded")
	}

	// Wait for window to reset
	time.Sleep(600 * time.Millisecond)

	// Should succeed after reset
	_, err = proxy.MakeRequest("/api/test")
	if err != nil {
		t.Errorf("Expected success after window reset, got error: %v", err)
	}
}

// BenchmarkImageProxy_LazyLoading benchmarks lazy loading performance.
func BenchmarkImageProxy_LazyLoading(b *testing.B) {
	for i := 0; i < b.N; i++ {
		proxy := NewImageProxy("bench.jpg")
		proxy.Display()
	}
}

// BenchmarkCachingProxy_CacheHit benchmarks cache hit performance.
func BenchmarkCachingProxy_CacheHit(b *testing.B) {
	proxy := NewCachingProxy("bench_service")
	proxy.GetData("test_key") // Prime cache

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proxy.GetData("test_key")
	}
}

// BenchmarkLoggingProxy benchmarks logging overhead.
func BenchmarkLoggingProxy(b *testing.B) {
	proxy := NewLoggingPaymentProxy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proxy.ProcessPayment(100.00, "bench_user")
	}
}
