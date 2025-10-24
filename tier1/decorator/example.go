package decorator

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// Example1_BasicDecoration demonstrates the simplest use case: wrapping a
// component with a single decorator.
func Example1_BasicDecoration() {
	fmt.Println("=== Example 1: Basic Decoration ===")

	// Create base processor
	processor := NewSimpleProcessor("BaseProcessor")

	// Process without decoration
	result, _ := processor.Process("Hello, World!")
	fmt.Println("Without decoration:", result)

	// Add logging decoration
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)
	decoratedProcessor := NewLoggingDecorator(processor, logger)

	// Process with decoration
	result, _ = decoratedProcessor.Process("Hello, World!")
	fmt.Println("With logging decoration:", result)
	fmt.Println("\nLog output:")
	fmt.Print(logBuf.String())
	fmt.Println()
}

// Example2_StackingDecorators demonstrates how multiple decorators can be
// stacked to combine behaviors. This is the real power of the decorator pattern.
func Example2_StackingDecorators() {
	fmt.Println("=== Example 2: Stacking Multiple Decorators ===")

	// Create base processor
	base := NewSimpleProcessor("SecureProcessor")

	// Stack decorators: Logging -> Validation -> Base
	// Note: Order matters! Logging is outermost, so it runs first
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "DECORATOR: ", 0)

	processor := NewLoggingDecorator(
		NewValidationDecorator(
			base,
			ValidationConfig{
				MinLength:  5,
				MaxLength:  100,
				AllowEmpty: false,
			},
		),
		logger,
	)

	// Test with valid data
	fmt.Println("Processing valid data:")
	result, err := processor.Process("This is valid data")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// Test with invalid data
	fmt.Println("\nProcessing invalid data (too short):")
	result, err = processor.Process("Hi")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	fmt.Println("\nLog output:")
	fmt.Print(logBuf.String())
	fmt.Println()
}

// Example3_DataPipelineWithCompression demonstrates a real-world scenario:
// building a data processing pipeline with validation, compression, and logging.
func Example3_DataPipelineWithCompression() {
	fmt.Println("=== Example 3: Data Pipeline with Compression ===")

	// Build a complete data processing pipeline
	// Order: Logging -> Compression -> Validation -> Base
	base := NewSimpleProcessor("DataPipeline")

	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "PIPELINE: ", 0)

	pipeline := NewLoggingDecorator(
		NewCompressionDecorator(
			NewValidationDecorator(
				base,
				ValidationConfig{
					MinLength:  1,
					MaxLength:  1000,
					AllowEmpty: false,
				},
			),
			6, // Compression level
		),
		logger,
	)

	// Process data through the pipeline
	data := "This is some data that will be validated, processed, and compressed."
	fmt.Println("Input data:", data)

	result, err := pipeline.Process(data)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Output (compressed): %d bytes\n", len(result))
		fmt.Printf("Original length: %d bytes\n", len(data))
	}

	fmt.Println("\nPipeline logs:")
	fmt.Print(logBuf.String())
	fmt.Println()
}

// Example4_CachingDecorator demonstrates a stateful decorator that maintains
// a cache to avoid reprocessing the same data.
func Example4_CachingDecorator() {
	fmt.Println("=== Example 4: Caching Decorator ===")

	// Create a processor with caching
	base := NewSimpleProcessor("ExpensiveProcessor")
	cache := NewCachingDecorator(base)

	// Process the same data multiple times
	testData := []string{
		"First request",
		"Second request",
		"First request",  // Cache hit
		"Third request",
		"First request",  // Cache hit
		"Second request", // Cache hit
	}

	for i, data := range testData {
		result, _ := cache.Process(data)
		fmt.Printf("%d. Processing: %s -> %s\n", i+1, data, result)
	}

	// Show cache statistics
	hits, misses := cache.Stats()
	fmt.Printf("\nCache Statistics:\n")
	fmt.Printf("  Hits: %d\n", hits)
	fmt.Printf("  Misses: %d\n", misses)
	fmt.Printf("  Hit Rate: %.2f%%\n", float64(hits)/float64(hits+misses)*100)
	fmt.Println()
}

// Example5_EncryptionDecorator demonstrates adding encryption to a data processor.
func Example5_EncryptionDecorator() {
	fmt.Println("=== Example 5: Encryption Decorator ===")

	// Create a secure processing pipeline with encryption
	base := NewSimpleProcessor("SecureProcessor")

	// AES-256 key (32 bytes)
	key := []byte("12345678901234567890123456789012")

	encryptedProcessor, err := NewEncryptionDecorator(base, key)
	if err != nil {
		fmt.Println("Error creating encryption decorator:", err)
		return
	}

	// Process sensitive data
	sensitiveData := "Secret message: Project Phoenix is a go!"
	fmt.Println("Original data:", sensitiveData)

	result, err := encryptedProcessor.Process(sensitiveData)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Encrypted result (%d bytes): %s...\n", len(result), result[:50])
	}

	// Demonstrate combining encryption with other decorators
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "SECURE: ", 0)

	secureLogger := NewLoggingDecorator(encryptedProcessor, logger)
	result, _ = secureLogger.Process(sensitiveData)

	fmt.Println("\nWith logging:")
	fmt.Print(logBuf.String())
	fmt.Println()
}

// Example6_HTTPMiddlewareChain demonstrates the functional decorator pattern
// commonly used for HTTP middleware in Go.
func Example6_HTTPMiddlewareChain() {
	fmt.Println("=== Example 6: HTTP Middleware Chain (Functional Decorators) ===")

	// Create a simple handler
	baseHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the base handler!")
	}

	// Create middleware decorators
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "HTTP: ", 0)
	metrics := &RequestMetrics{}

	// Chain middleware together
	// Order: Logging -> Auth -> Metrics -> BaseHandler
	handler := ChainMiddleware(
		baseHandler,
		LoggingMiddleware(logger),
		AuthMiddleware("Bearer secret-token"),
		MetricsMiddleware(metrics),
	)

	// Test with valid authentication
	fmt.Println("Request 1: Valid authentication")
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	w := httptest.NewRecorder()
	handler(w, req)
	fmt.Println("Response:", w.Code, "-", w.Body.String())

	// Test with invalid authentication
	fmt.Println("\nRequest 2: Invalid authentication")
	req = httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	w = httptest.NewRecorder()
	handler(w, req)
	fmt.Println("Response:", w.Code)

	// Show metrics
	fmt.Println("\nMetrics:")
	fmt.Printf("  Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("  Total Errors: %d\n", metrics.TotalErrors)
	fmt.Printf("  Average Duration: %v\n", metrics.AverageDuration())

	fmt.Println("\nMiddleware logs:")
	fmt.Print(logBuf.String())
	fmt.Println()
}

// Example7_DecoratorOrderMatters demonstrates that the order in which
// decorators are applied affects the behavior and performance of the system.
func Example7_DecoratorOrderMatters() {
	fmt.Println("=== Example 7: Decorator Order Matters ===")

	base := NewSimpleProcessor("OrderDemo")
	key := []byte("12345678901234567890123456789012")

	// Scenario 1: Compress then Encrypt
	// This is typically BETTER - compression works on plaintext
	fmt.Println("Scenario 1: Validation -> Compression -> Encryption")
	encryptor1, _ := NewEncryptionDecorator(
		NewCompressionDecorator(
			NewValidationDecorator(
				base,
				ValidationConfig{MinLength: 1, MaxLength: 1000, AllowEmpty: false},
			),
			6,
		),
		key,
	)

	data := "This is some test data that will be compressed and encrypted."
	result1, _ := encryptor1.Process(data)
	fmt.Printf("Result length: %d bytes\n", len(result1))

	// Scenario 2: Encrypt then Compress
	// This is typically WORSE - compression doesn't work well on encrypted data
	fmt.Println("\nScenario 2: Validation -> Encryption -> Compression")
	compressor2, _ := NewEncryptionDecorator(base, key)
	compressor2 = NewCompressionDecorator(
		NewValidationDecorator(
			compressor2,
			ValidationConfig{MinLength: 1, MaxLength: 1000, AllowEmpty: false},
		),
		6,
	)

	result2, _ := compressor2.Process(data)
	fmt.Printf("Result length: %d bytes\n", len(result2))

	fmt.Printf("\nDifference: %d bytes (Scenario 2 is larger because compression is less effective on encrypted data)\n", len(result2)-len(result1))
	fmt.Println()
}

// Example8_DynamicDecoration demonstrates adding decorators at runtime based
// on configuration or user preferences.
func Example8_DynamicDecoration() {
	fmt.Println("=== Example 8: Dynamic Decoration ===")

	// Configuration that determines which decorators to apply
	type Config struct {
		EnableLogging     bool
		EnableValidation  bool
		EnableCompression bool
		EnableCaching     bool
	}

	// Function to build a processor based on configuration
	buildProcessor := func(cfg Config) DataProcessor {
		var processor DataProcessor = NewSimpleProcessor("DynamicProcessor")

		// Apply decorators based on configuration
		if cfg.EnableValidation {
			processor = NewValidationDecorator(processor, ValidationConfig{
				MinLength:  1,
				MaxLength:  200,
				AllowEmpty: false,
			})
			fmt.Println("✓ Validation enabled")
		}

		if cfg.EnableCompression {
			processor = NewCompressionDecorator(processor, 6)
			fmt.Println("✓ Compression enabled")
		}

		if cfg.EnableCaching {
			processor = NewCachingDecorator(processor)
			fmt.Println("✓ Caching enabled")
		}

		if cfg.EnableLogging {
			processor = NewLoggingDecorator(processor, log.New(io.Discard, "", 0))
			fmt.Println("✓ Logging enabled")
		}

		return processor
	}

	// Example 1: Minimal configuration
	fmt.Println("Configuration 1: Minimal (no decorators)")
	processor1 := buildProcessor(Config{})
	result, _ := processor1.Process("Test data")
	fmt.Println("Result:", result)

	// Example 2: Full configuration
	fmt.Println("\nConfiguration 2: Full (all decorators)")
	processor2 := buildProcessor(Config{
		EnableLogging:     true,
		EnableValidation:  true,
		EnableCompression: true,
		EnableCaching:     true,
	})
	result, _ = processor2.Process("Test data")
	fmt.Printf("Result: %d bytes (compressed)\n", len(result))

	fmt.Println()
}
