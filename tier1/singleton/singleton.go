// Package singleton demonstrates the Singleton pattern in Go.
//
// The Singleton pattern ensures a class has only one instance while providing
// a global point of access to it. In Go, this is commonly achieved using sync.Once
// for lazy initialization or init() for eager initialization.
//
// This package demonstrates three different approaches:
// 1. Lazy initialization with sync.Once (thread-safe)
// 2. Eager initialization with init()
// 3. Singleton with reset capability for testing
package singleton

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Approach 1: Lazy Initialization with sync.Once
// ============================================================================

// Database represents a database connection pool.
// This is the singleton instance that will be shared across the application.
type Database struct {
	connectionString string
	maxConnections   int
	activeConnections int
	mu               sync.Mutex
}

var (
	// dbInstance holds the single database instance
	dbInstance *Database
	// dbOnce ensures the database is initialized only once
	dbOnce sync.Once
)

// GetDatabase returns the singleton database instance.
// The instance is created lazily on the first call using sync.Once,
// which guarantees thread-safe initialization even with concurrent calls.
func GetDatabase() *Database {
	dbOnce.Do(func() {
		fmt.Println("Initializing database connection pool (lazy)...")
		dbInstance = &Database{
			connectionString: "postgresql://localhost:5432/mydb",
			maxConnections:   100,
			activeConnections: 0,
		}
		// Simulate expensive initialization
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Database connection pool initialized")
	})
	return dbInstance
}

// Connect simulates acquiring a connection from the pool
func (db *Database) Connect() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeConnections >= db.maxConnections {
		return fmt.Errorf("max connections reached: %d", db.maxConnections)
	}

	db.activeConnections++
	return nil
}

// Disconnect simulates releasing a connection back to the pool
func (db *Database) Disconnect() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeConnections <= 0 {
		return fmt.Errorf("no active connections to release")
	}

	db.activeConnections--
	return nil
}

// Stats returns the current connection statistics
func (db *Database) Stats() string {
	db.mu.Lock()
	defer db.mu.Unlock()
	return fmt.Sprintf("Active: %d/%d connections", db.activeConnections, db.maxConnections)
}

// ============================================================================
// Approach 2: Eager Initialization with init()
// ============================================================================

// Configuration represents application configuration.
// It's initialized eagerly when the package is loaded.
type Configuration struct {
	AppName     string
	Port        int
	Environment string
	FeatureFlags map[string]bool
}

var (
	// Config is the global configuration instance
	// Initialized immediately when the package is loaded
	Config *Configuration
)

func init() {
	fmt.Println("Initializing configuration (eager)...")
	Config = &Configuration{
		AppName:     "MyApplication",
		Port:        8080,
		Environment: "development",
		FeatureFlags: map[string]bool{
			"feature_x": true,
			"feature_y": false,
			"beta_mode": true,
		},
	}
	fmt.Println("Configuration initialized")
}

// IsFeatureEnabled checks if a feature flag is enabled
func (c *Configuration) IsFeatureEnabled(feature string) bool {
	if c.FeatureFlags == nil {
		return false
	}
	return c.FeatureFlags[feature]
}

// GetServerAddress returns the full server address
func (c *Configuration) GetServerAddress() string {
	return fmt.Sprintf("http://localhost:%d", c.Port)
}

// ============================================================================
// Approach 3: Singleton with Reset Capability (for testing)
// ============================================================================

// Logger represents a singleton logger with reset capability for testing
type Logger struct {
	logFile string
	level   string
	mu      sync.Mutex
	logs    []string
}

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
	loggerMu       sync.Mutex // Protects reset operations
)

// GetLogger returns the singleton logger instance.
// Supports reset for testing scenarios.
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		fmt.Println("Initializing logger...")
		loggerInstance = &Logger{
			logFile: "app.log",
			level:   "INFO",
			logs:    make([]string, 0),
		}
	})
	return loggerInstance
}

// Log adds a log entry
func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), l.level, message)
	l.logs = append(l.logs, entry)
	fmt.Println(entry)
}

// GetLogs returns all logged messages
func (l *Logger) GetLogs() []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Return a copy to prevent external modification
	logsCopy := make([]string, len(l.logs))
	copy(logsCopy, l.logs)
	return logsCopy
}

// SetLevel changes the log level
func (l *Logger) SetLevel(level string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// ResetForTesting resets the logger singleton.
// WARNING: This should only be used in test code.
// It allows creating a fresh logger instance for each test.
func ResetForTesting() {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	loggerInstance = nil
	loggerOnce = sync.Once{}
}

// ============================================================================
// Approach 4: Thread-Safe Cache Singleton
// ============================================================================

// Cache represents a thread-safe in-memory cache singleton
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

var (
	cacheInstance *Cache
	cacheOnce     sync.Once
)

// GetCache returns the singleton cache instance
func GetCache() *Cache {
	cacheOnce.Do(func() {
		fmt.Println("Initializing cache...")
		cacheInstance = &Cache{
			data: make(map[string]interface{}),
		}
	})
	return cacheInstance
}

// Set stores a value in the cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Size returns the number of items in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]interface{})
}