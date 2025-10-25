// Package singleton demonstrates the Singleton pattern in Go.
// The Singleton pattern ensures a class has only one instance and provides
// a global point of access to it.
package singleton

import (
	"fmt"
	"sync"
)

// ConfigManager represents a singleton configuration manager.
// It holds application-wide configuration that should only exist once.
type ConfigManager struct {
	settings map[string]string
	mu       sync.RWMutex
}

var (
	configInstance *ConfigManager
	configOnce     sync.Once
)

// GetConfigManager returns the singleton instance of ConfigManager.
// Uses sync.Once to ensure thread-safe lazy initialization.
func GetConfigManager() *ConfigManager {
	configOnce.Do(func() {
		fmt.Println("Initializing ConfigManager singleton...")
		configInstance = &ConfigManager{
			settings: make(map[string]string),
		}
		// Load default configuration
		configInstance.settings["app_name"] = "MyApp"
		configInstance.settings["version"] = "1.0.0"
		configInstance.settings["environment"] = "development"
	})
	return configInstance
}

// Get retrieves a configuration value by key.
func (c *ConfigManager) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.settings[key]
	return val, ok
}

// Set updates a configuration value.
func (c *ConfigManager) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.settings[key] = value
}

// GetAll returns all configuration settings.
func (c *ConfigManager) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent external modification
	result := make(map[string]string, len(c.settings))
	for k, v := range c.settings {
		result[k] = v
	}
	return result
}

// DatabasePool represents a singleton database connection pool.
// This is a common use case where you want exactly one pool managing connections.
type DatabasePool struct {
	connections []string // Simplified; would be []*sql.DB in real code
	maxSize     int
	mu          sync.Mutex
}

var (
	dbPoolInstance *DatabasePool
	dbPoolOnce     sync.Once
)

// GetDatabasePool returns the singleton database pool instance.
func GetDatabasePool() *DatabasePool {
	dbPoolOnce.Do(func() {
		fmt.Println("Initializing DatabasePool singleton...")
		dbPoolInstance = &DatabasePool{
			connections: make([]string, 0),
			maxSize:     10,
		}
	})
	return dbPoolInstance
}

// AcquireConnection simulates acquiring a connection from the pool.
func (dp *DatabasePool) AcquireConnection() string {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if len(dp.connections) > 0 {
		// Reuse existing connection
		conn := dp.connections[0]
		dp.connections = dp.connections[1:]
		return conn
	}

	// Create new connection
	connID := fmt.Sprintf("conn_%d", len(dp.connections)+1)
	fmt.Printf("Creating new connection: %s\n", connID)
	return connID
}

// ReleaseConnection returns a connection back to the pool.
func (dp *DatabasePool) ReleaseConnection(conn string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if len(dp.connections) < dp.maxSize {
		dp.connections = append(dp.connections, conn)
		fmt.Printf("Connection %s returned to pool\n", conn)
	} else {
		fmt.Printf("Pool full, discarding connection %s\n", conn)
	}
}

// GetPoolSize returns the current number of pooled connections.
func (dp *DatabasePool) GetPoolSize() int {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	return len(dp.connections)
}

// Logger represents a singleton logger instance.
// Demonstrates eager initialization via package-level variable.
type Logger struct {
	prefix string
	mu     sync.Mutex
}

// Eager initialization - instance created at package load time
var loggerInstance = &Logger{
	prefix: "[APP]",
}

// GetLogger returns the singleton logger instance.
// No sync.Once needed since it's eagerly initialized.
func GetLogger() *Logger {
	return loggerInstance
}

// Log writes a log message with the configured prefix.
func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("%s %s\n", l.prefix, message)
}

// SetPrefix updates the log prefix.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}
