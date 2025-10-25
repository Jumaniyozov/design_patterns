// Package servicelocator demonstrates the Service Locator pattern.
// It provides a central registry for service lookup, though Dependency Injection
// is generally preferred in modern Go applications for better testability.
package servicelocator

import (
	"errors"
	"fmt"
	"sync"
)

// ServiceLocator is a central registry for services
type ServiceLocator struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewServiceLocator creates a service locator
func NewServiceLocator() *ServiceLocator {
	return &ServiceLocator{
		services: make(map[string]interface{}),
	}
}

// Register registers a service
func (sl *ServiceLocator) Register(name string, service interface{}) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.services[name] = service
}

// Get retrieves a service
func (sl *ServiceLocator) Get(name string) (interface{}, error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	service, exists := sl.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}

// Has checks if service exists
func (sl *ServiceLocator) Has(name string) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	_, exists := sl.services[name]
	return exists
}

// Remove removes a service
func (sl *ServiceLocator) Remove(name string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	delete(sl.services, name)
}

// Clear removes all services
func (sl *ServiceLocator) Clear() {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.services = make(map[string]interface{})
}

// Service interfaces

type Logger interface {
	Log(message string)
}

type Database interface {
	Query(sql string) ([]map[string]interface{}, error)
}

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

// Concrete implementations

type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(message string) {
	fmt.Printf("[LOG] %s\n", message)
}

type MockDatabase struct{}

func (m *MockDatabase) Query(sql string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"id": 1, "name": "test"},
	}, nil
}

type MemoryCache struct {
	data map[string]interface{}
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{data: make(map[string]interface{})}
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	val, exists := m.data[key]
	return val, exists
}

func (m *MemoryCache) Set(key string, value interface{}) {
	m.data[key] = value
}

// Global service locator instance (anti-pattern, but shown for completeness)
var globalLocator = NewServiceLocator()

// GetService gets a service from global locator
func GetService(name string) (interface{}, error) {
	return globalLocator.Get(name)
}

// RegisterService registers a service globally
func RegisterService(name string, service interface{}) {
	globalLocator.Register(name, service)
}

// Application using service locator
type Application struct {
	locator *ServiceLocator
}

// NewApplication creates an application
func NewApplication(locator *ServiceLocator) *Application {
	return &Application{locator: locator}
}

// DoWork performs work using located services
func (a *Application) DoWork() error {
	// Locate logger
	loggerService, err := a.locator.Get("logger")
	if err != nil {
		return err
	}
	logger := loggerService.(Logger)
	logger.Log("Starting work")

	// Locate database
	dbService, err := a.locator.Get("database")
	if err != nil {
		return err
	}
	db := dbService.(Database)
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		return err
	}

	logger.Log(fmt.Sprintf("Found %d results", len(results)))

	// Locate cache
	cacheService, err := a.locator.Get("cache")
	if err != nil {
		return err
	}
	cache := cacheService.(Cache)
	cache.Set("last_query", results)

	logger.Log("Work completed")
	return nil
}

// TypedServiceLocator provides type-safe service location
type TypedServiceLocator struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewTypedServiceLocator creates a typed service locator
func NewTypedServiceLocator() *TypedServiceLocator {
	return &TypedServiceLocator{
		services: make(map[string]interface{}),
	}
}

// RegisterLogger registers a logger
func (tsl *TypedServiceLocator) RegisterLogger(logger Logger) {
	tsl.mu.Lock()
	defer tsl.mu.Unlock()
	tsl.services["logger"] = logger
}

// GetLogger gets a logger
func (tsl *TypedServiceLocator) GetLogger() (Logger, error) {
	tsl.mu.RLock()
	defer tsl.mu.RUnlock()

	service, exists := tsl.services["logger"]
	if !exists {
		return nil, errors.New("logger not found")
	}
	return service.(Logger), nil
}

// RegisterDatabase registers a database
func (tsl *TypedServiceLocator) RegisterDatabase(db Database) {
	tsl.mu.Lock()
	defer tsl.mu.Unlock()
	tsl.services["database"] = db
}

// GetDatabase gets a database
func (tsl *TypedServiceLocator) GetDatabase() (Database, error) {
	tsl.mu.RLock()
	defer tsl.mu.RUnlock()

	service, exists := tsl.services["database"]
	if !exists {
		return nil, errors.New("database not found")
	}
	return service.(Database), nil
}
