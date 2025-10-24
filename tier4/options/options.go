// Package options implements the Functional Options pattern.
//
// The Functional Options pattern is a Go-specific pattern for configuring
// objects with optional parameters. It uses variadic functions to provide a
// clean, extensible API without requiring multiple constructor variations.
//
// Key characteristics:
// - Self-documenting API
// - Extensible without breaking changes
// - Type-safe
// - Optional parameters with defaults
// - No need for builder pattern complexity
//
// This pattern is idiomatic in Go and used extensively in the standard library
// (e.g., grpc.Dial options).
package options

import (
	"crypto/tls"
	"time"
)

// Option is a function that configures an object.
// This is the core type for the functional options pattern.
type Option[T any] func(*T)

// Apply applies multiple options to a value.
func Apply[T any](target *T, options ...Option[T]) {
	for _, option := range options {
		option(target)
	}
}

// Server demonstrates the options pattern with a server configuration.
type Server struct {
	Host            string
	Port            int
	Timeout         time.Duration
	MaxConnections  int
	TLSConfig       *tls.Config
	EnableLogging   bool
	EnableMetrics   bool
	ShutdownTimeout time.Duration
	MiddlewareChain []string
}

// ServerOption is an option for configuring a Server.
type ServerOption func(*Server)

// NewServer creates a new server with default values and applies options.
func NewServer(options ...ServerOption) *Server {
	// Set defaults
	server := &Server{
		Host:            "localhost",
		Port:            8080,
		Timeout:         30 * time.Second,
		MaxConnections:  100,
		EnableLogging:   false,
		EnableMetrics:   false,
		ShutdownTimeout: 10 * time.Second,
		MiddlewareChain: []string{},
	}

	// Apply options
	for _, option := range options {
		option(server)
	}

	return server
}

// WithHost sets the server host.
func WithHost(host string) ServerOption {
	return func(s *Server) {
		s.Host = host
	}
}

// WithPort sets the server port.
func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.Port = port
	}
}

// WithTimeout sets the server timeout.
func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

// WithMaxConnections sets the maximum number of connections.
func WithMaxConnections(max int) ServerOption {
	return func(s *Server) {
		s.MaxConnections = max
	}
}

// WithTLS enables TLS with the given configuration.
func WithTLS(config *tls.Config) ServerOption {
	return func(s *Server) {
		s.TLSConfig = config
	}
}

// WithLogging enables logging.
func WithLogging() ServerOption {
	return func(s *Server) {
		s.EnableLogging = true
	}
}

// WithMetrics enables metrics collection.
func WithMetrics() ServerOption {
	return func(s *Server) {
		s.EnableMetrics = true
	}
}

// WithShutdownTimeout sets the graceful shutdown timeout.
func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.ShutdownTimeout = timeout
	}
}

// WithMiddleware adds middleware to the chain.
func WithMiddleware(middleware string) ServerOption {
	return func(s *Server) {
		s.MiddlewareChain = append(s.MiddlewareChain, middleware)
	}
}

// DatabaseConfig demonstrates options with validation.
type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxConnections  int
	ConnectTimeout  time.Duration
	QueryTimeout    time.Duration
	SSLMode         string
	RetryAttempts   int
	RetryDelay      time.Duration
}

// DatabaseOption is an option for configuring a database connection.
type DatabaseOption func(*DatabaseConfig) error

// NewDatabaseConfig creates a new database configuration with validation.
func NewDatabaseConfig(options ...DatabaseOption) (*DatabaseConfig, error) {
	config := &DatabaseConfig{
		Host:           "localhost",
		Port:           5432,
		MaxConnections: 10,
		ConnectTimeout: 5 * time.Second,
		QueryTimeout:   30 * time.Second,
		SSLMode:        "require",
		RetryAttempts:  3,
		RetryDelay:     time.Second,
	}

	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// WithDBHost sets the database host.
func WithDBHost(host string) DatabaseOption {
	return func(c *DatabaseConfig) error {
		c.Host = host
		return nil
	}
}

// WithDBPort sets the database port.
func WithDBPort(port int) DatabaseOption {
	return func(c *DatabaseConfig) error {
		if port <= 0 || port > 65535 {
			return nil // Could return error for validation
		}
		c.Port = port
		return nil
	}
}

// WithDBCredentials sets username and password.
func WithDBCredentials(username, password string) DatabaseOption {
	return func(c *DatabaseConfig) error {
		c.Username = username
		c.Password = password
		return nil
	}
}

// WithDatabase sets the database name.
func WithDatabase(name string) DatabaseOption {
	return func(c *DatabaseConfig) error {
		c.Database = name
		return nil
	}
}

// WithDBMaxConnections sets the connection pool size.
func WithDBMaxConnections(max int) DatabaseOption {
	return func(c *DatabaseConfig) error {
		if max < 1 {
			max = 1
		}
		c.MaxConnections = max
		return nil
	}
}

// WithDBTimeouts sets connection and query timeouts.
func WithDBTimeouts(connect, query time.Duration) DatabaseOption {
	return func(c *DatabaseConfig) error {
		c.ConnectTimeout = connect
		c.QueryTimeout = query
		return nil
	}
}

// WithSSLMode sets the SSL mode.
func WithSSLMode(mode string) DatabaseOption {
	return func(c *DatabaseConfig) error {
		c.SSLMode = mode
		return nil
	}
}

// WithRetry configures retry behavior.
func WithRetry(attempts int, delay time.Duration) DatabaseOption {
	return func(c *DatabaseConfig) error {
		c.RetryAttempts = attempts
		c.RetryDelay = delay
		return nil
	}
}

// Client demonstrates options with builder-like behavior.
type Client struct {
	BaseURL     string
	Timeout     time.Duration
	MaxRetries  int
	Headers     map[string]string
	BearerToken string
	UserAgent   string
}

// ClientOption is an option for configuring a client.
type ClientOption func(*Client)

// NewClient creates a new HTTP client with options.
func NewClient(baseURL string, options ...ClientOption) *Client {
	client := &Client{
		BaseURL:    baseURL,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Headers:    make(map[string]string),
		UserAgent:  "Go-Client/1.0",
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithClientTimeout sets the client timeout.
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(retries int) ClientOption {
	return func(c *Client) {
		c.MaxRetries = retries
	}
}

// WithHeader adds a header to the client.
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.Headers[key] = value
	}
}

// WithBearerToken sets the bearer token for authentication.
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.BearerToken = token
		c.Headers["Authorization"] = "Bearer " + token
	}
}

// WithUserAgent sets the user agent.
func WithUserAgent(agent string) ClientOption {
	return func(c *Client) {
		c.UserAgent = agent
	}
}

// Logger demonstrates options with different configuration levels.
type Logger struct {
	Level      string
	Output     string
	Format     string
	TimeFormat string
	Prefix     string
	EnableCaller bool
}

// LoggerOption is an option for configuring a logger.
type LoggerOption func(*Logger)

// NewLogger creates a new logger with options.
func NewLogger(options ...LoggerOption) *Logger {
	logger := &Logger{
		Level:      "info",
		Output:     "stdout",
		Format:     "json",
		TimeFormat: time.RFC3339,
		EnableCaller: false,
	}

	for _, option := range options {
		option(logger)
	}

	return logger
}

// WithLevel sets the log level.
func WithLevel(level string) LoggerOption {
	return func(l *Logger) {
		l.Level = level
	}
}

// WithOutput sets the output destination.
func WithOutput(output string) LoggerOption {
	return func(l *Logger) {
		l.Output = output
	}
}

// WithFormat sets the log format.
func WithFormat(format string) LoggerOption {
	return func(l *Logger) {
		l.Format = format
	}
}

// WithTimeFormat sets the time format.
func WithTimeFormat(format string) LoggerOption {
	return func(l *Logger) {
		l.TimeFormat = format
	}
}

// WithPrefix sets a log prefix.
func WithPrefix(prefix string) LoggerOption {
	return func(l *Logger) {
		l.Prefix = prefix
	}
}

// WithCaller enables caller information in logs.
func WithCaller() LoggerOption {
	return func(l *Logger) {
		l.EnableCaller = true
	}
}

// Combining options example
func DevelopmentLogger() LoggerOption {
	return func(l *Logger) {
		WithLevel("debug")(l)
		WithFormat("text")(l)
		WithCaller()(l)
	}
}

func ProductionLogger() LoggerOption {
	return func(l *Logger) {
		WithLevel("info")(l)
		WithFormat("json")(l)
		WithOutput("file")(l)
	}
}
