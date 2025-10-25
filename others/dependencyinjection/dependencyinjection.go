// Package dependencyinjection demonstrates the Dependency Injection pattern.
// It inverts control by providing dependencies to objects rather than having
// them create dependencies, promoting loose coupling and testability.
package dependencyinjection

import (
	"errors"
	"fmt"
)

// Database interface abstracts data storage
type Database interface {
	Query(query string) ([]map[string]interface{}, error)
	Execute(command string) error
}

// Logger interface abstracts logging
type Logger interface {
	Info(message string)
	Error(message string)
	Debug(message string)
}

// Cache interface abstracts caching
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

// Concrete implementations

// PostgreSQL is a concrete database implementation
type PostgreSQL struct {
	connectionString string
}

func NewPostgreSQL(connStr string) *PostgreSQL {
	return &PostgreSQL{connectionString: connStr}
}

func (p *PostgreSQL) Query(query string) ([]map[string]interface{}, error) {
	// Simulate database query
	return []map[string]interface{}{
		{"id": 1, "name": "User1"},
	}, nil
}

func (p *PostgreSQL) Execute(command string) error {
	// Simulate database command
	return nil
}

// MySQL is another concrete database implementation
type MySQL struct {
	host string
	port int
}

func NewMySQL(host string, port int) *MySQL {
	return &MySQL{host: host, port: port}
}

func (m *MySQL) Query(query string) ([]map[string]interface{}, error) {
	// Simulate MySQL query
	return []map[string]interface{}{
		{"id": 1, "name": "User1"},
	}, nil
}

func (m *MySQL) Execute(command string) error {
	// Simulate MySQL command
	return nil
}

// ConsoleLogger logs to console
type ConsoleLogger struct {
	prefix string
}

func NewConsoleLogger(prefix string) *ConsoleLogger {
	return &ConsoleLogger{prefix: prefix}
}

func (c *ConsoleLogger) Info(message string) {
	fmt.Printf("[%s INFO] %s\n", c.prefix, message)
}

func (c *ConsoleLogger) Error(message string) {
	fmt.Printf("[%s ERROR] %s\n", c.prefix, message)
}

func (c *ConsoleLogger) Debug(message string) {
	fmt.Printf("[%s DEBUG] %s\n", c.prefix, message)
}

// MemoryCache is an in-memory cache implementation
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

// UserService demonstrates constructor injection
type UserService struct {
	db     Database
	logger Logger
	cache  Cache
}

// NewUserService creates a user service with injected dependencies
func NewUserService(db Database, logger Logger, cache Cache) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
		cache:  cache,
	}
}

// GetUser retrieves a user
func (us *UserService) GetUser(id string) (map[string]interface{}, error) {
	us.logger.Debug(fmt.Sprintf("Getting user %s", id))

	// Check cache first
	if cached, exists := us.cache.Get(id); exists {
		us.logger.Info("Cache hit")
		return cached.(map[string]interface{}), nil
	}

	// Query database
	us.logger.Info("Cache miss, querying database")
	results, err := us.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %s", id))
	if err != nil {
		us.logger.Error(fmt.Sprintf("Database error: %v", err))
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("user not found")
	}

	// Cache result
	us.cache.Set(id, results[0])
	return results[0], nil
}

// CreateUser creates a new user
func (us *UserService) CreateUser(name string) error {
	us.logger.Info(fmt.Sprintf("Creating user: %s", name))
	return us.db.Execute(fmt.Sprintf("INSERT INTO users (name) VALUES ('%s')", name))
}

// EmailService demonstrates method injection
type EmailService struct {
	logger Logger
}

func NewEmailService(logger Logger) *EmailService {
	return &EmailService{logger: logger}
}

// SendWithProvider uses method-level dependency injection
func (es *EmailService) SendWithProvider(provider EmailProvider, to, subject, body string) error {
	es.logger.Info(fmt.Sprintf("Sending email to %s via %s", to, provider.Name()))
	return provider.Send(to, subject, body)
}

// EmailProvider interface for different email providers
type EmailProvider interface {
	Send(to, subject, body string) error
	Name() string
}

// SMTPProvider sends emails via SMTP
type SMTPProvider struct {
	server string
}

func NewSMTPProvider(server string) *SMTPProvider {
	return &SMTPProvider{server: server}
}

func (s *SMTPProvider) Send(to, subject, body string) error {
	// Simulate sending
	return nil
}

func (s *SMTPProvider) Name() string {
	return "SMTP"
}

// Container is a simple DI container
type Container struct {
	services map[string]interface{}
}

// NewContainer creates a DI container
func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

// Register registers a service
func (c *Container) Register(name string, service interface{}) {
	c.services[name] = service
}

// Resolve resolves a service
func (c *Container) Resolve(name string) (interface{}, error) {
	if service, exists := c.services[name]; exists {
		return service, nil
	}
	return nil, fmt.Errorf("service %s not found", name)
}

// BuildUserService is a factory that wires dependencies
func (c *Container) BuildUserService() (*UserService, error) {
	db, err := c.Resolve("database")
	if err != nil {
		return nil, err
	}

	logger, err := c.Resolve("logger")
	if err != nil {
		return nil, err
	}

	cache, err := c.Resolve("cache")
	if err != nil {
		return nil, err
	}

	return NewUserService(
		db.(Database),
		logger.(Logger),
		cache.(Cache),
	), nil
}
