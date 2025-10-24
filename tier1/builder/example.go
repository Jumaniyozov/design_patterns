package builder

import (
	"fmt"
	"time"
)

// Example1_BasicDatabasePool demonstrates basic usage of the database pool builder
// with minimal configuration.
func Example1_BasicDatabasePool() {
	fmt.Println("=== Example 1: Basic Database Connection Pool ===")
	fmt.Println()

	// Simple configuration with just required fields
	pool, err := NewDBConnectionPoolBuilder().
		WithDatabase("myapp").
		WithCredentials("admin", "secret123").
		Build()

	if err != nil {
		fmt.Printf("Error building pool: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", pool)
	fmt.Println("\nThis demonstrates the builder with minimal configuration,")
	fmt.Println("using sensible defaults for all optional fields.")
}

// Example2_ProductionDatabasePool demonstrates a production-grade database pool
// configuration with all optional settings specified.
func Example2_ProductionDatabasePool() {
	fmt.Println("\n=== Example 2: Production Database Pool ===")
	fmt.Println()

	// Full production configuration with all optional fields
	pool, err := NewDBConnectionPoolBuilder().
		WithHost("db.production.com").
		WithPort(5432).
		WithDatabase("production_db").
		WithCredentials("prod_user", "secure_password").
		WithMaxConnections(50).
		WithMinConnections(5).
		WithConnectionTimeout(10 * time.Second).
		WithReadTimeout(30 * time.Second).
		WithWriteTimeout(30 * time.Second).
		WithSSLMode("verify-full").
		WithRetryAttempts(5).
		WithCompression(true).
		WithQueryLogging(true).
		Build()

	if err != nil {
		fmt.Printf("Error building pool: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", pool)
	fmt.Println("\nThis demonstrates a fully configured production database pool")
	fmt.Println("with security, performance, and monitoring features enabled.")
}

// Example3_ValidationErrors demonstrates how the builder handles validation errors.
func Example3_ValidationErrors() {
	fmt.Println("\n=== Example 3: Builder Validation Errors ===")
	fmt.Println()

	// Attempt to build with invalid configuration
	pool, err := NewDBConnectionPoolBuilder().
		WithPort(-1).                   // Invalid port
		WithMaxConnections(-5).         // Invalid max connections
		WithSSLMode("invalid_mode").    // Invalid SSL mode
		WithConnectionTimeout(-1).      // Invalid timeout
		Build()

	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	} else {
		fmt.Printf("Unexpected success: %s\n", pool)
	}

	fmt.Println("\nThis demonstrates how the builder accumulates and reports")
	fmt.Println("all validation errors, helping developers catch issues early.")
}

// Example4_CrossFieldValidation demonstrates validation that depends on multiple fields.
func Example4_CrossFieldValidation() {
	fmt.Println("\n=== Example 4: Cross-Field Validation ===")
	fmt.Println()

	// Min connections greater than max connections
	pool, err := NewDBConnectionPoolBuilder().
		WithDatabase("test_db").
		WithCredentials("user", "pass").
		WithMinConnections(20).
		WithMaxConnections(10). // Less than min!
		Build()

	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	} else {
		fmt.Printf("Unexpected success: %s\n", pool)
	}

	fmt.Println("\nThis demonstrates cross-field validation performed in Build(),")
	fmt.Println("ensuring logical consistency across all configuration settings.")
}

// Example5_HTTPRequestBuilder demonstrates building HTTP requests.
func Example5_HTTPRequestBuilder() {
	fmt.Println("\n=== Example 5: HTTP Request Builder ===")
	fmt.Println()

	// Build a POST request with headers and body
	request, err := NewHTTPRequestBuilder().
		WithMethod("POST").
		WithURL("https://api.example.com/users").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer token123").
		WithBody([]byte(`{"name":"John","email":"john@example.com"}`)).
		WithTimeout(15 * time.Second).
		WithRetries(3).
		Build()

	if err != nil {
		fmt.Printf("Error building request: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", request)
	fmt.Println("\nThis demonstrates building a complex HTTP request with")
	fmt.Println("method, headers, body, timeout, and retry configuration.")
}

// Example6_HTTPRequestWithMultipleHeaders demonstrates setting multiple headers at once.
func Example6_HTTPRequestWithMultipleHeaders() {
	fmt.Println("\n=== Example 6: HTTP Request with Multiple Headers ===")
	fmt.Println()

	headers := map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"Authorization":   "Bearer token123",
		"X-Request-ID":    "req-12345",
		"X-Correlation-ID": "corr-67890",
	}

	request, err := NewHTTPRequestBuilder().
		WithMethod("GET").
		WithURL("https://api.example.com/data").
		WithHeaders(headers).
		WithTimeout(10 * time.Second).
		Build()

	if err != nil {
		fmt.Printf("Error building request: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", request)
	fmt.Println("\nThis demonstrates batch setting of multiple headers,")
	fmt.Println("showing flexibility in the builder API design.")
}

// Example7_ApplicationConfig demonstrates building application configuration.
func Example7_ApplicationConfig() {
	fmt.Println("\n=== Example 7: Application Configuration Builder ===")
	fmt.Println()

	// Build production application configuration
	config, err := NewAppConfigBuilder().
		WithServerPort(8080).
		WithServerHost("0.0.0.0").
		WithDatabaseURL("postgres://localhost:5432/myapp").
		WithLogLevel("info").
		WithMaxWorkers(8).
		WithMetrics(true).
		WithProfiling(false).
		WithTracing(true).
		WithShutdownTimeout(45 * time.Second).
		WithEnvironment("production").
		Build()

	if err != nil {
		fmt.Printf("Error building config: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", config)
	fmt.Println("\nThis demonstrates building application configuration")
	fmt.Println("with observability features and graceful shutdown settings.")
}

// Example8_DevelopmentVsProduction demonstrates different configurations
// for different environments.
func Example8_DevelopmentVsProduction() {
	fmt.Println("\n=== Example 8: Development vs Production Config ===")
	fmt.Println()

	// Development configuration
	devConfig, _ := NewAppConfigBuilder().
		WithServerPort(3000).
		WithLogLevel("debug").
		WithMaxWorkers(2).
		WithProfiling(true).
		WithEnvironment("development").
		Build()

	fmt.Printf("Development: %s\n", devConfig)

	// Production configuration
	prodConfig, _ := NewAppConfigBuilder().
		WithServerPort(8080).
		WithLogLevel("warn").
		WithMaxWorkers(16).
		WithMetrics(true).
		WithTracing(true).
		WithEnvironment("production").
		Build()

	fmt.Printf("Production:  %s\n", prodConfig)

	fmt.Println("\nThis demonstrates how the same builder pattern can create")
	fmt.Println("different configurations optimized for different environments.")
}

// Example9_PartialConfiguration demonstrates that you only need to specify
// the fields that differ from defaults.
func Example9_PartialConfiguration() {
	fmt.Println("\n=== Example 9: Partial Configuration (Defaults) ===")
	fmt.Println()

	// Only override specific fields, use defaults for everything else
	config, err := NewAppConfigBuilder().
		WithServerPort(9090).
		WithLogLevel("debug").
		Build()

	if err != nil {
		fmt.Printf("Error building config: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", config)
	fmt.Println("\nThis demonstrates the power of sensible defaults:")
	fmt.Println("you only specify what you need to change, keeping code clean.")
}

// Example10_BuilderReuse demonstrates why builders should be used once.
func Example10_BuilderReuse() {
	fmt.Println("\n=== Example 10: Builder Reuse (Anti-pattern) ===")
	fmt.Println()

	// Create a builder
	builder := NewAppConfigBuilder().
		WithServerPort(8080).
		WithLogLevel("info")

	// Build once - OK
	config1, _ := builder.Build()
	fmt.Printf("First build:  %s\n", config1)

	// Modify and build again - This is an anti-pattern!
	// The builder may have accumulated state from previous build
	builder.WithServerPort(9090)
	config2, _ := builder.Build()
	fmt.Printf("Second build: %s\n", config2)

	fmt.Println("\n⚠️  WARNING: Builders should be used once and discarded.")
	fmt.Println("Reusing builders can lead to unexpected behavior.")
	fmt.Println("Best practice: Create a new builder for each object.")
}
