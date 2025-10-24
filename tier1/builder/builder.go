// Package builder demonstrates the Builder Pattern, a creational design pattern
// that separates the construction of complex objects from their representation.
//
// The Builder Pattern is particularly useful in Go where:
// - Constructor overloading is not supported
// - Objects have many optional fields
// - Complex initialization logic is required
// - You want to create immutable objects with validation
package builder

import (
	"errors"
	"fmt"
	"time"
)

// =============================================================================
// Example 1: Database Connection Pool Builder
// =============================================================================

// DBConnectionPool represents a configured database connection pool.
// This object has many optional configuration parameters, making it
// an ideal candidate for the Builder pattern.
type DBConnectionPool struct {
	host               string
	port               int
	database           string
	username           string
	password           string
	maxConnections     int
	minConnections     int
	connectionTimeout  time.Duration
	readTimeout        time.Duration
	writeTimeout       time.Duration
	sslMode            string
	retryAttempts      int
	enableCompression  bool
	enableQueryLogging bool
}

// DBConnectionPoolBuilder builds DBConnectionPool instances step by step.
// The builder is mutable during construction and should be used once.
type DBConnectionPoolBuilder struct {
	pool *DBConnectionPool
	errs []error
}

// NewDBConnectionPoolBuilder creates a new builder with sensible defaults.
func NewDBConnectionPoolBuilder() *DBConnectionPoolBuilder {
	return &DBConnectionPoolBuilder{
		pool: &DBConnectionPool{
			host:              "localhost",
			port:              5432,
			maxConnections:    10,
			minConnections:    1,
			connectionTimeout: 30 * time.Second,
			readTimeout:       5 * time.Second,
			writeTimeout:      5 * time.Second,
			sslMode:           "prefer",
			retryAttempts:     3,
		},
		errs: []error{},
	}
}

// WithHost sets the database host.
func (b *DBConnectionPoolBuilder) WithHost(host string) *DBConnectionPoolBuilder {
	if host == "" {
		b.errs = append(b.errs, errors.New("host cannot be empty"))
	}
	b.pool.host = host
	return b
}

// WithPort sets the database port.
func (b *DBConnectionPoolBuilder) WithPort(port int) *DBConnectionPoolBuilder {
	if port <= 0 || port > 65535 {
		b.errs = append(b.errs, fmt.Errorf("invalid port: %d", port))
	}
	b.pool.port = port
	return b
}

// WithDatabase sets the database name.
func (b *DBConnectionPoolBuilder) WithDatabase(database string) *DBConnectionPoolBuilder {
	if database == "" {
		b.errs = append(b.errs, errors.New("database name cannot be empty"))
	}
	b.pool.database = database
	return b
}

// WithCredentials sets both username and password.
func (b *DBConnectionPoolBuilder) WithCredentials(username, password string) *DBConnectionPoolBuilder {
	if username == "" {
		b.errs = append(b.errs, errors.New("username cannot be empty"))
	}
	b.pool.username = username
	b.pool.password = password
	return b
}

// WithMaxConnections sets the maximum number of connections in the pool.
func (b *DBConnectionPoolBuilder) WithMaxConnections(max int) *DBConnectionPoolBuilder {
	if max <= 0 {
		b.errs = append(b.errs, errors.New("max connections must be positive"))
	}
	b.pool.maxConnections = max
	return b
}

// WithMinConnections sets the minimum number of connections in the pool.
func (b *DBConnectionPoolBuilder) WithMinConnections(min int) *DBConnectionPoolBuilder {
	if min < 0 {
		b.errs = append(b.errs, errors.New("min connections cannot be negative"))
	}
	b.pool.minConnections = min
	return b
}

// WithConnectionTimeout sets the connection timeout.
func (b *DBConnectionPoolBuilder) WithConnectionTimeout(timeout time.Duration) *DBConnectionPoolBuilder {
	if timeout <= 0 {
		b.errs = append(b.errs, errors.New("connection timeout must be positive"))
	}
	b.pool.connectionTimeout = timeout
	return b
}

// WithReadTimeout sets the read timeout.
func (b *DBConnectionPoolBuilder) WithReadTimeout(timeout time.Duration) *DBConnectionPoolBuilder {
	if timeout <= 0 {
		b.errs = append(b.errs, errors.New("read timeout must be positive"))
	}
	b.pool.readTimeout = timeout
	return b
}

// WithWriteTimeout sets the write timeout.
func (b *DBConnectionPoolBuilder) WithWriteTimeout(timeout time.Duration) *DBConnectionPoolBuilder {
	if timeout <= 0 {
		b.errs = append(b.errs, errors.New("write timeout must be positive"))
	}
	b.pool.writeTimeout = timeout
	return b
}

// WithSSLMode sets the SSL mode (disable, prefer, require, verify-ca, verify-full).
func (b *DBConnectionPoolBuilder) WithSSLMode(mode string) *DBConnectionPoolBuilder {
	validModes := map[string]bool{
		"disable":     true,
		"prefer":      true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validModes[mode] {
		b.errs = append(b.errs, fmt.Errorf("invalid SSL mode: %s", mode))
	}
	b.pool.sslMode = mode
	return b
}

// WithRetryAttempts sets the number of retry attempts.
func (b *DBConnectionPoolBuilder) WithRetryAttempts(attempts int) *DBConnectionPoolBuilder {
	if attempts < 0 {
		b.errs = append(b.errs, errors.New("retry attempts cannot be negative"))
	}
	b.pool.retryAttempts = attempts
	return b
}

// WithCompression enables or disables compression.
func (b *DBConnectionPoolBuilder) WithCompression(enabled bool) *DBConnectionPoolBuilder {
	b.pool.enableCompression = enabled
	return b
}

// WithQueryLogging enables or disables query logging.
func (b *DBConnectionPoolBuilder) WithQueryLogging(enabled bool) *DBConnectionPoolBuilder {
	b.pool.enableQueryLogging = enabled
	return b
}

// Build validates all settings and returns the final DBConnectionPool.
// Returns an error if any validation failed during construction.
func (b *DBConnectionPoolBuilder) Build() (*DBConnectionPool, error) {
	// Check for accumulated errors
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("builder validation failed: %v", b.errs)
	}

	// Additional cross-field validation
	if b.pool.minConnections > b.pool.maxConnections {
		return nil, errors.New("min connections cannot exceed max connections")
	}

	// Required fields validation
	if b.pool.database == "" {
		return nil, errors.New("database name is required")
	}

	if b.pool.username == "" {
		return nil, errors.New("username is required")
	}

	return b.pool, nil
}

// String provides a readable representation of the connection pool config.
func (p *DBConnectionPool) String() string {
	return fmt.Sprintf(
		"DBConnectionPool{host=%s, port=%d, database=%s, maxConn=%d, minConn=%d, timeout=%s, ssl=%s}",
		p.host, p.port, p.database, p.maxConnections, p.minConnections,
		p.connectionTimeout, p.sslMode,
	)
}

// =============================================================================
// Example 2: HTTP Request Builder
// =============================================================================

// HTTPRequest represents a configured HTTP request.
// Demonstrates builder pattern for complex request configuration.
type HTTPRequest struct {
	method  string
	url     string
	headers map[string]string
	body    []byte
	timeout time.Duration
	retries int
}

// HTTPRequestBuilder builds HTTPRequest instances.
type HTTPRequestBuilder struct {
	request *HTTPRequest
}

// NewHTTPRequestBuilder creates a new HTTP request builder with defaults.
func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		request: &HTTPRequest{
			method:  "GET",
			headers: make(map[string]string),
			timeout: 30 * time.Second,
			retries: 0,
		},
	}
}

// WithMethod sets the HTTP method (GET, POST, PUT, DELETE, etc.).
func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	b.request.method = method
	return b
}

// WithURL sets the request URL.
func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.request.url = url
	return b
}

// WithHeader adds a header to the request.
func (b *HTTPRequestBuilder) WithHeader(key, value string) *HTTPRequestBuilder {
	b.request.headers[key] = value
	return b
}

// WithHeaders sets multiple headers at once.
func (b *HTTPRequestBuilder) WithHeaders(headers map[string]string) *HTTPRequestBuilder {
	for k, v := range headers {
		b.request.headers[k] = v
	}
	return b
}

// WithBody sets the request body.
func (b *HTTPRequestBuilder) WithBody(body []byte) *HTTPRequestBuilder {
	b.request.body = body
	return b
}

// WithTimeout sets the request timeout.
func (b *HTTPRequestBuilder) WithTimeout(timeout time.Duration) *HTTPRequestBuilder {
	b.request.timeout = timeout
	return b
}

// WithRetries sets the number of retry attempts.
func (b *HTTPRequestBuilder) WithRetries(retries int) *HTTPRequestBuilder {
	b.request.retries = retries
	return b
}

// Build validates and returns the final HTTPRequest.
func (b *HTTPRequestBuilder) Build() (*HTTPRequest, error) {
	if b.request.url == "" {
		return nil, errors.New("URL is required")
	}
	if b.request.method == "" {
		return nil, errors.New("method is required")
	}
	if b.request.timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}
	return b.request, nil
}

// String provides a readable representation of the HTTP request.
func (r *HTTPRequest) String() string {
	return fmt.Sprintf(
		"HTTPRequest{method=%s, url=%s, headers=%d, body=%d bytes, timeout=%s, retries=%d}",
		r.method, r.url, len(r.headers), len(r.body), r.timeout, r.retries,
	)
}

// =============================================================================
// Example 3: Application Configuration Builder
// =============================================================================

// AppConfig represents application configuration with many optional settings.
type AppConfig struct {
	serverPort     int
	serverHost     string
	databaseURL    string
	logLevel       string
	maxWorkers     int
	enableMetrics  bool
	enableProfiling bool
	enableTracing  bool
	shutdownTimeout time.Duration
	environment    string
}

// AppConfigBuilder builds AppConfig instances.
type AppConfigBuilder struct {
	config *AppConfig
}

// NewAppConfigBuilder creates a new application config builder with defaults.
func NewAppConfigBuilder() *AppConfigBuilder {
	return &AppConfigBuilder{
		config: &AppConfig{
			serverPort:      8080,
			serverHost:      "0.0.0.0",
			logLevel:        "info",
			maxWorkers:      4,
			enableMetrics:   false,
			enableProfiling: false,
			enableTracing:   false,
			shutdownTimeout: 30 * time.Second,
			environment:     "development",
		},
	}
}

// WithServerPort sets the server port.
func (b *AppConfigBuilder) WithServerPort(port int) *AppConfigBuilder {
	b.config.serverPort = port
	return b
}

// WithServerHost sets the server host.
func (b *AppConfigBuilder) WithServerHost(host string) *AppConfigBuilder {
	b.config.serverHost = host
	return b
}

// WithDatabaseURL sets the database connection URL.
func (b *AppConfigBuilder) WithDatabaseURL(url string) *AppConfigBuilder {
	b.config.databaseURL = url
	return b
}

// WithLogLevel sets the logging level (debug, info, warn, error).
func (b *AppConfigBuilder) WithLogLevel(level string) *AppConfigBuilder {
	b.config.logLevel = level
	return b
}

// WithMaxWorkers sets the maximum number of worker goroutines.
func (b *AppConfigBuilder) WithMaxWorkers(workers int) *AppConfigBuilder {
	b.config.maxWorkers = workers
	return b
}

// WithMetrics enables or disables metrics collection.
func (b *AppConfigBuilder) WithMetrics(enabled bool) *AppConfigBuilder {
	b.config.enableMetrics = enabled
	return b
}

// WithProfiling enables or disables profiling.
func (b *AppConfigBuilder) WithProfiling(enabled bool) *AppConfigBuilder {
	b.config.enableProfiling = enabled
	return b
}

// WithTracing enables or disables distributed tracing.
func (b *AppConfigBuilder) WithTracing(enabled bool) *AppConfigBuilder {
	b.config.enableTracing = enabled
	return b
}

// WithShutdownTimeout sets the graceful shutdown timeout.
func (b *AppConfigBuilder) WithShutdownTimeout(timeout time.Duration) *AppConfigBuilder {
	b.config.shutdownTimeout = timeout
	return b
}

// WithEnvironment sets the environment (development, staging, production).
func (b *AppConfigBuilder) WithEnvironment(env string) *AppConfigBuilder {
	b.config.environment = env
	return b
}

// Build validates and returns the final AppConfig.
func (b *AppConfigBuilder) Build() (*AppConfig, error) {
	if b.config.serverPort <= 0 || b.config.serverPort > 65535 {
		return nil, fmt.Errorf("invalid server port: %d", b.config.serverPort)
	}
	if b.config.maxWorkers <= 0 {
		return nil, errors.New("max workers must be positive")
	}
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[b.config.logLevel] {
		return nil, fmt.Errorf("invalid log level: %s", b.config.logLevel)
	}
	return b.config, nil
}

// String provides a readable representation of the application config.
func (c *AppConfig) String() string {
	return fmt.Sprintf(
		"AppConfig{port=%d, host=%s, logLevel=%s, workers=%d, metrics=%t, env=%s}",
		c.serverPort, c.serverHost, c.logLevel, c.maxWorkers, c.enableMetrics, c.environment,
	)
}
