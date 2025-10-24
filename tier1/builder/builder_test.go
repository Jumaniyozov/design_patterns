package builder

import (
	"testing"
	"time"
)

// TestDBConnectionPoolBuilder_BasicBuild tests basic database pool construction.
func TestDBConnectionPoolBuilder_BasicBuild(t *testing.T) {
	pool, err := NewDBConnectionPoolBuilder().
		WithDatabase("testdb").
		WithCredentials("user", "pass").
		Build()

	if err != nil {
		t.Fatalf("Expected successful build, got error: %v", err)
	}

	if pool.database != "testdb" {
		t.Errorf("Expected database 'testdb', got '%s'", pool.database)
	}

	if pool.username != "user" {
		t.Errorf("Expected username 'user', got '%s'", pool.username)
	}

	// Check defaults are applied
	if pool.host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", pool.host)
	}

	if pool.port != 5432 {
		t.Errorf("Expected default port 5432, got %d", pool.port)
	}
}

// TestDBConnectionPoolBuilder_FullConfiguration tests all builder options.
func TestDBConnectionPoolBuilder_FullConfiguration(t *testing.T) {
	pool, err := NewDBConnectionPoolBuilder().
		WithHost("db.example.com").
		WithPort(3306).
		WithDatabase("production").
		WithCredentials("admin", "secret").
		WithMaxConnections(100).
		WithMinConnections(10).
		WithConnectionTimeout(20 * time.Second).
		WithReadTimeout(5 * time.Second).
		WithWriteTimeout(5 * time.Second).
		WithSSLMode("require").
		WithRetryAttempts(5).
		WithCompression(true).
		WithQueryLogging(true).
		Build()

	if err != nil {
		t.Fatalf("Expected successful build, got error: %v", err)
	}

	if pool.host != "db.example.com" {
		t.Errorf("Expected host 'db.example.com', got '%s'", pool.host)
	}

	if pool.maxConnections != 100 {
		t.Errorf("Expected maxConnections 100, got %d", pool.maxConnections)
	}

	if pool.sslMode != "require" {
		t.Errorf("Expected sslMode 'require', got '%s'", pool.sslMode)
	}

	if !pool.enableCompression {
		t.Error("Expected compression to be enabled")
	}
}

// TestDBConnectionPoolBuilder_MissingRequiredFields tests validation of required fields.
func TestDBConnectionPoolBuilder_MissingRequiredFields(t *testing.T) {
	// Missing database
	_, err := NewDBConnectionPoolBuilder().
		WithCredentials("user", "pass").
		Build()

	if err == nil {
		t.Error("Expected error for missing database, got nil")
	}

	// Missing username
	_, err = NewDBConnectionPoolBuilder().
		WithDatabase("testdb").
		Build()

	if err == nil {
		t.Error("Expected error for missing username, got nil")
	}
}

// TestDBConnectionPoolBuilder_InvalidPort tests port validation.
func TestDBConnectionPoolBuilder_InvalidPort(t *testing.T) {
	_, err := NewDBConnectionPoolBuilder().
		WithDatabase("testdb").
		WithCredentials("user", "pass").
		WithPort(-1).
		Build()

	if err == nil {
		t.Error("Expected error for invalid port, got nil")
	}

	_, err = NewDBConnectionPoolBuilder().
		WithDatabase("testdb").
		WithCredentials("user", "pass").
		WithPort(70000).
		Build()

	if err == nil {
		t.Error("Expected error for port > 65535, got nil")
	}
}

// TestDBConnectionPoolBuilder_InvalidSSLMode tests SSL mode validation.
func TestDBConnectionPoolBuilder_InvalidSSLMode(t *testing.T) {
	_, err := NewDBConnectionPoolBuilder().
		WithDatabase("testdb").
		WithCredentials("user", "pass").
		WithSSLMode("invalid").
		Build()

	if err == nil {
		t.Error("Expected error for invalid SSL mode, got nil")
	}
}

// TestDBConnectionPoolBuilder_CrossFieldValidation tests min/max connections validation.
func TestDBConnectionPoolBuilder_CrossFieldValidation(t *testing.T) {
	_, err := NewDBConnectionPoolBuilder().
		WithDatabase("testdb").
		WithCredentials("user", "pass").
		WithMinConnections(50).
		WithMaxConnections(10).
		Build()

	if err == nil {
		t.Error("Expected error when min > max connections, got nil")
	}
}

// TestHTTPRequestBuilder_BasicBuild tests basic HTTP request construction.
func TestHTTPRequestBuilder_BasicBuild(t *testing.T) {
	request, err := NewHTTPRequestBuilder().
		WithURL("https://example.com").
		Build()

	if err != nil {
		t.Fatalf("Expected successful build, got error: %v", err)
	}

	if request.url != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got '%s'", request.url)
	}

	if request.method != "GET" {
		t.Errorf("Expected default method 'GET', got '%s'", request.method)
	}
}

// TestHTTPRequestBuilder_WithHeaders tests setting multiple headers.
func TestHTTPRequestBuilder_WithHeaders(t *testing.T) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}

	request, err := NewHTTPRequestBuilder().
		WithURL("https://api.example.com").
		WithHeaders(headers).
		Build()

	if err != nil {
		t.Fatalf("Expected successful build, got error: %v", err)
	}

	if len(request.headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(request.headers))
	}

	if request.headers["Content-Type"] != "application/json" {
		t.Error("Content-Type header not set correctly")
	}
}

// TestHTTPRequestBuilder_MissingURL tests URL validation.
func TestHTTPRequestBuilder_MissingURL(t *testing.T) {
	_, err := NewHTTPRequestBuilder().
		WithMethod("POST").
		Build()

	if err == nil {
		t.Error("Expected error for missing URL, got nil")
	}
}

// TestAppConfigBuilder_BasicBuild tests basic app config construction.
func TestAppConfigBuilder_BasicBuild(t *testing.T) {
	config, err := NewAppConfigBuilder().Build()

	if err != nil {
		t.Fatalf("Expected successful build with defaults, got error: %v", err)
	}

	if config.serverPort != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.serverPort)
	}

	if config.environment != "development" {
		t.Errorf("Expected default environment 'development', got '%s'", config.environment)
	}
}

// TestAppConfigBuilder_InvalidPort tests port validation.
func TestAppConfigBuilder_InvalidPort(t *testing.T) {
	_, err := NewAppConfigBuilder().
		WithServerPort(-1).
		Build()

	if err == nil {
		t.Error("Expected error for invalid port, got nil")
	}

	_, err = NewAppConfigBuilder().
		WithServerPort(70000).
		Build()

	if err == nil {
		t.Error("Expected error for port > 65535, got nil")
	}
}

// TestAppConfigBuilder_InvalidLogLevel tests log level validation.
func TestAppConfigBuilder_InvalidLogLevel(t *testing.T) {
	_, err := NewAppConfigBuilder().
		WithLogLevel("invalid").
		Build()

	if err == nil {
		t.Error("Expected error for invalid log level, got nil")
	}
}

// TestAppConfigBuilder_InvalidMaxWorkers tests max workers validation.
func TestAppConfigBuilder_InvalidMaxWorkers(t *testing.T) {
	_, err := NewAppConfigBuilder().
		WithMaxWorkers(0).
		Build()

	if err == nil {
		t.Error("Expected error for zero max workers, got nil")
	}

	_, err = NewAppConfigBuilder().
		WithMaxWorkers(-5).
		Build()

	if err == nil {
		t.Error("Expected error for negative max workers, got nil")
	}
}

// TestAppConfigBuilder_EnvironmentConfigurations tests different environment setups.
func TestAppConfigBuilder_EnvironmentConfigurations(t *testing.T) {
	// Development config
	devConfig, err := NewAppConfigBuilder().
		WithEnvironment("development").
		WithLogLevel("debug").
		WithProfiling(true).
		Build()

	if err != nil {
		t.Fatalf("Failed to build dev config: %v", err)
	}

	if devConfig.logLevel != "debug" {
		t.Error("Dev config should use debug log level")
	}

	// Production config
	prodConfig, err := NewAppConfigBuilder().
		WithEnvironment("production").
		WithLogLevel("error").
		WithMetrics(true).
		WithTracing(true).
		Build()

	if err != nil {
		t.Fatalf("Failed to build prod config: %v", err)
	}

	if !prodConfig.enableMetrics {
		t.Error("Prod config should have metrics enabled")
	}
}

// BenchmarkDBConnectionPoolBuilder benchmarks the builder performance.
func BenchmarkDBConnectionPoolBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewDBConnectionPoolBuilder().
			WithDatabase("benchdb").
			WithCredentials("user", "pass").
			WithMaxConnections(20).
			Build()
	}
}

// BenchmarkHTTPRequestBuilder benchmarks HTTP request builder performance.
func BenchmarkHTTPRequestBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewHTTPRequestBuilder().
			WithURL("https://example.com").
			WithMethod("POST").
			WithHeader("Content-Type", "application/json").
			Build()
	}
}
