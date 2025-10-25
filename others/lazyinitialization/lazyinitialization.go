// Package lazyinitialization demonstrates the Lazy Initialization pattern.
// It defers expensive object creation or computation until first use,
// improving startup time and resource utilization.
package lazyinitialization

import (
	"fmt"
	"sync"
	"time"
)

// ExpensiveResource represents a resource that's costly to create
type ExpensiveResource struct {
	data string
}

// NewExpensiveResource simulates expensive initialization
func NewExpensiveResource() *ExpensiveResource {
	fmt.Println("Creating expensive resource... (this takes time)")
	time.Sleep(500 * time.Millisecond) // Simulate expensive operation
	return &ExpensiveResource{
		data: "Expensive data loaded",
	}
}

// GetData retrieves resource data
func (e *ExpensiveResource) GetData() string {
	return e.data
}

// LazyResource lazily initializes a resource using sync.Once
type LazyResource struct {
	resource *ExpensiveResource
	once     sync.Once
}

// Get lazily initializes and returns the resource
func (l *LazyResource) Get() *ExpensiveResource {
	l.once.Do(func() {
		l.resource = NewExpensiveResource()
	})
	return l.resource
}

// Config demonstrates lazy configuration loading
type Config struct {
	settings map[string]string
	loaded   bool
	mu       sync.RWMutex
}

// NewConfig creates a config without loading
func NewConfig() *Config {
	return &Config{
		settings: nil,
		loaded:   false,
	}
}

// loadSettings simulates expensive config loading
func (c *Config) loadSettings() {
	fmt.Println("Loading configuration from file/database...")
	time.Sleep(200 * time.Millisecond) // Simulate I/O
	c.settings = map[string]string{
		"database": "postgresql://localhost:5432",
		"cache":    "redis://localhost:6379",
		"api_key":  "secret-key-12345",
	}
	c.loaded = true
}

// Get lazily loads config and returns a setting
func (c *Config) Get(key string) string {
	c.mu.RLock()
	if c.loaded {
		defer c.mu.RUnlock()
		return c.settings[key]
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if !c.loaded {
		c.loadSettings()
	}

	return c.settings[key]
}

// LazyMap demonstrates lazy initialization of map entries
type LazyMap struct {
	data    map[string]*ExpensiveResource
	mu      sync.RWMutex
	factory func(string) *ExpensiveResource
}

// NewLazyMap creates a lazy map with a factory function
func NewLazyMap(factory func(string) *ExpensiveResource) *LazyMap {
	return &LazyMap{
		data:    make(map[string]*ExpensiveResource),
		factory: factory,
	}
}

// Get lazily creates and caches the resource for the key
func (lm *LazyMap) Get(key string) *ExpensiveResource {
	lm.mu.RLock()
	if resource, exists := lm.data[key]; exists {
		lm.mu.RUnlock()
		return resource
	}
	lm.mu.RUnlock()

	lm.mu.Lock()
	defer lm.mu.Unlock()

	// Double-check
	if resource, exists := lm.data[key]; exists {
		return resource
	}

	// Create resource lazily
	resource := lm.factory(key)
	lm.data[key] = resource
	return resource
}

// ComputedValue demonstrates lazy evaluation with caching
type ComputedValue struct {
	compute func() interface{}
	value   interface{}
	cached  bool
	mu      sync.Mutex
}

// NewComputedValue creates a lazy computed value
func NewComputedValue(compute func() interface{}) *ComputedValue {
	return &ComputedValue{
		compute: compute,
		cached:  false,
	}
}

// Get computes and caches the value on first call
func (cv *ComputedValue) Get() interface{} {
	cv.mu.Lock()
	defer cv.mu.Unlock()

	if !cv.cached {
		fmt.Println("Computing value... (expensive operation)")
		cv.value = cv.compute()
		cv.cached = true
	}
	return cv.value
}

// Reset clears the cached value
func (cv *ComputedValue) Reset() {
	cv.mu.Lock()
	defer cv.mu.Unlock()
	cv.cached = false
	cv.value = nil
}

// Database demonstrates lazy connection
type Database struct {
	conn *DatabaseConnection
	once sync.Once
}

// DatabaseConnection represents a database connection
type DatabaseConnection struct {
	connectionString string
	connected        bool
}

// NewDatabaseConnection creates a database connection
func NewDatabaseConnection(connStr string) *DatabaseConnection {
	fmt.Println("Establishing database connection...")
	time.Sleep(300 * time.Millisecond) // Simulate connection time
	return &DatabaseConnection{
		connectionString: connStr,
		connected:        true,
	}
}

// Query executes a query
func (dc *DatabaseConnection) Query(sql string) string {
	if !dc.connected {
		return "Error: not connected"
	}
	return fmt.Sprintf("Executed: %s", sql)
}

// NewDatabase creates a database without connecting
func NewDatabase() *Database {
	return &Database{}
}

// GetConnection lazily establishes database connection
func (db *Database) GetConnection() *DatabaseConnection {
	db.once.Do(func() {
		db.conn = NewDatabaseConnection("postgresql://localhost:5432/mydb")
	})
	return db.conn
}

// Query executes a query with lazy connection
func (db *Database) Query(sql string) string {
	return db.GetConnection().Query(sql)
}
