// Package proxy demonstrates the Proxy Pattern, a structural design pattern
// that provides a placeholder or surrogate for another object to control access to it.
//
// The Proxy Pattern is essential for:
// - Lazy initialization of expensive objects
// - Access control and permission checking
// - Caching expensive computations
// - Logging and monitoring access patterns
// - Rate limiting and resource throttling
package proxy

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// Example 1: Lazy Loading Image Proxy
// =============================================================================

// Image is the subject interface that both RealImage and ImageProxy implement.
type Image interface {
	Display() string
	GetSize() int
}

// RealImage represents an expensive image loaded from disk.
// Loading and decoding images is a heavy operation.
type RealImage struct {
	filename string
	data     []byte
	loadTime time.Time
}

// NewRealImage simulates loading an image from disk (expensive operation).
func NewRealImage(filename string) *RealImage {
	fmt.Printf("📁 Loading image from disk: %s...\n", filename)
	time.Sleep(100 * time.Millisecond) // Simulate slow disk I/O

	// Simulate image data
	fakeData := make([]byte, 1024*100) // 100KB
	for i := range fakeData {
		fakeData[i] = byte(i % 256)
	}

	fmt.Printf("✓ Image loaded: %s (100KB)\n", filename)
	return &RealImage{
		filename: filename,
		data:     fakeData,
		loadTime: time.Now(),
	}
}

// Display shows the image.
func (r *RealImage) Display() string {
	return fmt.Sprintf("Displaying image: %s (loaded at %s)",
		r.filename, r.loadTime.Format("15:04:05"))
}

// GetSize returns the image size in bytes.
func (r *RealImage) GetSize() int {
	return len(r.data)
}

// ImageProxy provides lazy loading for images.
// The real image is only loaded when first accessed.
type ImageProxy struct {
	filename  string
	realImage *RealImage
	mu        sync.Mutex // Ensure thread-safe lazy initialization
}

// NewImageProxy creates a new image proxy (fast - doesn't load the image).
func NewImageProxy(filename string) *ImageProxy {
	fmt.Printf("📋 Created image proxy for: %s (not loaded yet)\n", filename)
	return &ImageProxy{
		filename: filename,
	}
}

// Display loads the image if needed (lazy loading) and displays it.
func (p *ImageProxy) Display() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.realImage == nil {
		fmt.Println("🔄 First access - loading real image...")
		p.realImage = NewRealImage(p.filename)
	} else {
		fmt.Println("⚡ Using already-loaded image (fast!)")
	}
	return p.realImage.Display()
}

// GetSize loads the image if needed and returns its size.
func (p *ImageProxy) GetSize() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.realImage == nil {
		p.realImage = NewRealImage(p.filename)
	}
	return p.realImage.GetSize()
}

// =============================================================================
// Example 2: Access Control Database Proxy
// =============================================================================

// Database is the subject interface for database operations.
type Database interface {
	Query(sql string) (string, error)
	Execute(sql string) (int, error)
}

// RealDatabase represents the actual database.
type RealDatabase struct {
	name string
}

// NewRealDatabase creates a new real database instance.
func NewRealDatabase(name string) *RealDatabase {
	return &RealDatabase{name: name}
}

// Query executes a SELECT query.
func (d *RealDatabase) Query(sql string) (string, error) {
	// Simulate query execution
	return fmt.Sprintf("Query result from %s: [row1, row2, row3]", d.name), nil
}

// Execute executes an INSERT/UPDATE/DELETE query.
func (d *RealDatabase) Execute(sql string) (int, error) {
	// Simulate execution
	return 3, nil // 3 rows affected
}

// User represents a user with permissions.
type User struct {
	Name        string
	Permissions []string
}

// HasPermission checks if user has a specific permission.
func (u *User) HasPermission(permission string) bool {
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// DatabaseProxy provides access control to the database.
type DatabaseProxy struct {
	user *User
	db   *RealDatabase
	mu   sync.Mutex
}

// NewDatabaseProxy creates a new database proxy with access control.
func NewDatabaseProxy(user *User, dbName string) *DatabaseProxy {
	return &DatabaseProxy{
		user: user,
		db:   NewRealDatabase(dbName),
	}
}

// Query checks permissions before allowing query execution.
func (p *DatabaseProxy) Query(sql string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("🔒 Checking permissions for user '%s' to run query...\n", p.user.Name)

	if !p.user.HasPermission("read") {
		return "", fmt.Errorf("access denied: user '%s' lacks 'read' permission", p.user.Name)
	}

	fmt.Printf("✓ Permission granted to user '%s'\n", p.user.Name)
	return p.db.Query(sql)
}

// Execute checks permissions before allowing modification queries.
func (p *DatabaseProxy) Execute(sql string) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("🔒 Checking permissions for user '%s' to execute command...\n", p.user.Name)

	if !p.user.HasPermission("write") {
		return 0, fmt.Errorf("access denied: user '%s' lacks 'write' permission", p.user.Name)
	}

	fmt.Printf("✓ Permission granted to user '%s'\n", p.user.Name)
	return p.db.Execute(sql)
}

// =============================================================================
// Example 3: Caching Proxy
// =============================================================================

// DataService is the subject interface for data retrieval.
type DataService interface {
	GetData(key string) (string, error)
	ComputeExpensive(input int) (int, error)
}

// ExpensiveDataService represents a service with expensive operations.
type ExpensiveDataService struct {
	name string
}

// NewExpensiveDataService creates a new expensive service.
func NewExpensiveDataService(name string) *ExpensiveDataService {
	return &ExpensiveDataService{name: name}
}

// GetData simulates expensive data retrieval.
func (s *ExpensiveDataService) GetData(key string) (string, error) {
	fmt.Printf("💰 Expensive operation: Fetching data for key '%s'...\n", key)
	time.Sleep(200 * time.Millisecond) // Simulate slow operation
	return fmt.Sprintf("Data for %s from %s", key, s.name), nil
}

// ComputeExpensive simulates expensive computation.
func (s *ExpensiveDataService) ComputeExpensive(input int) (int, error) {
	fmt.Printf("💰 Expensive computation: Processing %d...\n", input)
	time.Sleep(300 * time.Millisecond) // Simulate slow computation
	return input * input * input, nil // Cubic calculation
}

// CachingProxy caches results to avoid repeated expensive operations.
type CachingProxy struct {
	service      *ExpensiveDataService
	dataCache    map[string]string
	computeCache map[int]int
	mu           sync.RWMutex
}

// NewCachingProxy creates a new caching proxy.
func NewCachingProxy(serviceName string) *CachingProxy {
	return &CachingProxy{
		service:      NewExpensiveDataService(serviceName),
		dataCache:    make(map[string]string),
		computeCache: make(map[int]int),
	}
}

// GetData checks cache first, only calls real service on cache miss.
func (p *CachingProxy) GetData(key string) (string, error) {
	// Check cache (read lock)
	p.mu.RLock()
	if cached, exists := p.dataCache[key]; exists {
		p.mu.RUnlock()
		fmt.Printf("⚡ Cache HIT for key '%s' (no expensive operation!)\n", key)
		return cached, nil
	}
	p.mu.RUnlock()

	// Cache miss - acquire write lock
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check in case another goroutine cached it
	if cached, exists := p.dataCache[key]; exists {
		fmt.Printf("⚡ Cache HIT on retry for key '%s'\n", key)
		return cached, nil
	}

	fmt.Printf("❌ Cache MISS for key '%s'\n", key)
	result, err := p.service.GetData(key)
	if err != nil {
		return "", err
	}

	p.dataCache[key] = result
	fmt.Printf("📥 Cached result for key '%s'\n", key)
	return result, nil
}

// ComputeExpensive checks cache first for computation results.
func (p *CachingProxy) ComputeExpensive(input int) (int, error) {
	p.mu.RLock()
	if cached, exists := p.computeCache[input]; exists {
		p.mu.RUnlock()
		fmt.Printf("⚡ Cache HIT for computation input %d (instant!)\n", input)
		return cached, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("❌ Cache MISS for computation input %d\n", input)
	result, err := p.service.ComputeExpensive(input)
	if err != nil {
		return 0, err
	}

	p.computeCache[input] = result
	fmt.Printf("📥 Cached computation result for input %d\n", input)
	return result, nil
}

// =============================================================================
// Example 4: Logging Proxy
// =============================================================================

// PaymentService is the subject interface for payment operations.
type PaymentService interface {
	ProcessPayment(amount float64, customer string) (string, error)
}

// RealPaymentService processes actual payments.
type RealPaymentService struct{}

// NewRealPaymentService creates a new payment service.
func NewRealPaymentService() *RealPaymentService {
	return &RealPaymentService{}
}

// ProcessPayment processes the payment.
func (s *RealPaymentService) ProcessPayment(amount float64, customer string) (string, error) {
	if amount <= 0 {
		return "", errors.New("invalid amount")
	}
	txnID := fmt.Sprintf("TXN-%d", time.Now().Unix())
	return txnID, nil
}

// LoggingPaymentProxy logs all payment operations.
type LoggingPaymentProxy struct {
	service   *RealPaymentService
	mu        sync.Mutex
	callCount int
}

// NewLoggingPaymentProxy creates a new logging proxy.
func NewLoggingPaymentProxy() *LoggingPaymentProxy {
	return &LoggingPaymentProxy{
		service: NewRealPaymentService(),
	}
}

// ProcessPayment logs the payment operation and delegates to real service.
func (p *LoggingPaymentProxy) ProcessPayment(amount float64, customer string) (string, error) {
	p.mu.Lock()
	p.callCount++
	callNum := p.callCount
	p.mu.Unlock()

	start := time.Now()
	fmt.Printf("📝 [LOG #%d] Payment request: $%.2f for customer '%s' at %s\n",
		callNum, amount, customer, start.Format("15:04:05"))

	txnID, err := p.service.ProcessPayment(amount, customer)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("📝 [LOG #%d] Payment FAILED: %v (took %v)\n", callNum, err, duration)
		return "", err
	}

	fmt.Printf("📝 [LOG #%d] Payment SUCCESS: %s (took %v)\n", callNum, txnID, duration)
	return txnID, nil
}

// GetCallCount returns the number of payment calls made.
func (p *LoggingPaymentProxy) GetCallCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.callCount
}

// =============================================================================
// Example 5: Rate Limiting Proxy
// =============================================================================

// APIService is the subject interface for API operations.
type APIService interface {
	MakeRequest(endpoint string) (string, error)
}

// RealAPIService makes actual API requests.
type RealAPIService struct{}

// NewRealAPIService creates a new API service.
func NewRealAPIService() *RealAPIService {
	return &RealAPIService{}
}

// MakeRequest makes an API request.
func (s *RealAPIService) MakeRequest(endpoint string) (string, error) {
	return fmt.Sprintf("Response from %s", endpoint), nil
}

// RateLimitingProxy limits the rate of API calls.
type RateLimitingProxy struct {
	service       *RealAPIService
	maxRequests   int
	timeWindow    time.Duration
	requests      []time.Time
	mu            sync.Mutex
}

// NewRateLimitingProxy creates a new rate-limiting proxy.
func NewRateLimitingProxy(maxRequests int, window time.Duration) *RateLimitingProxy {
	return &RateLimitingProxy{
		service:     NewRealAPIService(),
		maxRequests: maxRequests,
		timeWindow:  window,
		requests:    make([]time.Time, 0),
	}
}

// MakeRequest enforces rate limiting before calling the real service.
func (p *RateLimitingProxy) MakeRequest(endpoint string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()

	// Remove old requests outside the time window
	cutoff := now.Add(-p.timeWindow)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range p.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	p.requests = validRequests

	// Check if rate limit exceeded
	if len(p.requests) >= p.maxRequests {
		fmt.Printf("🚫 Rate limit exceeded! (%d requests in %v window)\n",
			len(p.requests), p.timeWindow)
		return "", errors.New("rate limit exceeded")
	}

	// Record this request
	p.requests = append(p.requests, now)
	fmt.Printf("✓ Request allowed (%d/%d in current window)\n",
		len(p.requests), p.maxRequests)

	return p.service.MakeRequest(endpoint)
}
