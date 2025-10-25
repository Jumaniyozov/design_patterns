// Package nullobject demonstrates the Null Object pattern.
// It provides default objects with neutral behavior, eliminating nil checks
// and preventing nil pointer errors.
package nullobject

import "fmt"

// Logger interface
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

// ConsoleLogger logs to console
type ConsoleLogger struct{}

func (c *ConsoleLogger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (c *ConsoleLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}

func (c *ConsoleLogger) Debug(msg string) {
	fmt.Printf("[DEBUG] %s\n", msg)
}

// NullLogger is a no-op logger (null object)
type NullLogger struct{}

func (n *NullLogger) Info(msg string)  {}
func (n *NullLogger) Error(msg string) {}
func (n *NullLogger) Debug(msg string) {}

// Cache interface
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}

// MemoryCache is a real cache implementation
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

func (m *MemoryCache) Delete(key string) {
	delete(m.data, key)
}

// NullCache is a no-op cache (null object)
type NullCache struct{}

func (n *NullCache) Get(key string) (interface{}, bool) {
	return nil, false // Always cache miss
}

func (n *NullCache) Set(key string, value interface{}) {}
func (n *NullCache) Delete(key string)                 {}

// Service demonstrates using null objects
type Service struct {
	logger Logger
	cache  Cache
}

// NewService creates a service with logger and cache
// If nil is passed, null objects are used
func NewService(logger Logger, cache Cache) *Service {
	if logger == nil {
		logger = &NullLogger{}
	}
	if cache == nil {
		cache = &NullCache{}
	}
	return &Service{
		logger: logger,
		cache:  cache,
	}
}

// DoWork performs work with logging and caching
func (s *Service) DoWork(id string) string {
	s.logger.Info(fmt.Sprintf("Starting work for %s", id))

	// Check cache
	if cached, exists := s.cache.Get(id); exists {
		s.logger.Debug("Cache hit")
		return cached.(string)
	}

	// Do work
	result := fmt.Sprintf("Result for %s", id)

	// Cache result
	s.cache.Set(id, result)
	s.logger.Info("Work completed")

	return result
}

// Notifier interface
type Notifier interface {
	Notify(message string)
}

// EmailNotifier sends email notifications
type EmailNotifier struct{}

func (e *EmailNotifier) Notify(message string) {
	fmt.Printf("Email sent: %s\n", message)
}

// NullNotifier is a no-op notifier (null object)
type NullNotifier struct{}

func (n *NullNotifier) Notify(message string) {}

// Customer with optional notifier
type Customer struct {
	name     string
	notifier Notifier
}

// NewCustomer creates a customer with optional notifier
func NewCustomer(name string, notifier Notifier) *Customer {
	if notifier == nil {
		notifier = &NullNotifier{}
	}
	return &Customer{
		name:     name,
		notifier: notifier,
	}
}

// Purchase makes a purchase and notifies
func (c *Customer) Purchase(item string) {
	// No nil check needed!
	c.notifier.Notify(fmt.Sprintf("%s purchased %s", c.name, item))
}
