package decorator

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Basic Component Tests
// ============================================================================

func TestSimpleProcessor(t *testing.T) {
	processor := NewSimpleProcessor("TestProcessor")
	result, err := processor.Process("test data")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "[TestProcessor] test data"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// ============================================================================
// Decorator Tests
// ============================================================================

func TestLoggingDecorator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLogs []string
	}{
		{
			name:  "logs processing",
			input: "test data",
			wantLogs: []string{
				"Processing started",
				"Processing completed",
			},
		},
		{
			name:  "logs input length",
			input: "hello",
			wantLogs: []string{
				"input length = 5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := NewSimpleProcessor("Test")
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)

			decorated := NewLoggingDecorator(base, logger)
			_, err := decorated.Process(tt.input)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			logOutput := logBuf.String()
			for _, wantLog := range tt.wantLogs {
				if !strings.Contains(logOutput, wantLog) {
					t.Errorf("Expected log to contain %q, got %q", wantLog, logOutput)
				}
			}
		})
	}
}

func TestValidationDecorator(t *testing.T) {
	tests := []struct {
		name      string
		config    ValidationConfig
		input     string
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid data",
			config: ValidationConfig{
				MinLength:  5,
				MaxLength:  20,
				AllowEmpty: false,
			},
			input:     "hello world",
			wantError: false,
		},
		{
			name: "empty data not allowed",
			config: ValidationConfig{
				MinLength:  1,
				MaxLength:  100,
				AllowEmpty: false,
			},
			input:     "",
			wantError: true,
			errorMsg:  "empty data not allowed",
		},
		{
			name: "data too short",
			config: ValidationConfig{
				MinLength:  10,
				MaxLength:  100,
				AllowEmpty: false,
			},
			input:     "short",
			wantError: true,
			errorMsg:  "less than minimum",
		},
		{
			name: "data too long",
			config: ValidationConfig{
				MinLength:  1,
				MaxLength:  5,
				AllowEmpty: false,
			},
			input:     "this is too long",
			wantError: true,
			errorMsg:  "exceeds maximum",
		},
		{
			name: "empty data allowed",
			config: ValidationConfig{
				MinLength:  0,
				MaxLength:  100,
				AllowEmpty: true,
			},
			input:     "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := NewSimpleProcessor("Test")
			decorated := NewValidationDecorator(base, tt.config)

			_, err := decorated.Process(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestCompressionDecorator(t *testing.T) {
	base := NewSimpleProcessor("Test")
	decorated := NewCompressionDecorator(base, 6)

	input := "This is some data that should be compressed successfully."
	result, err := decorated.Process(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Compressed result should be base64 encoded
	if !isBase64(result) {
		t.Error("Expected result to be base64 encoded")
	}
}

func TestEncryptionDecorator(t *testing.T) {
	tests := []struct {
		name      string
		keySize   int
		wantError bool
	}{
		{"AES-128", 16, false},
		{"AES-192", 24, false},
		{"AES-256", 32, false},
		{"Invalid key size", 15, true},
		{"Invalid key size", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := NewSimpleProcessor("Test")
			key := make([]byte, tt.keySize)

			decorated, err := NewEncryptionDecorator(base, key)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error for invalid key size, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Test processing
			result, err := decorated.Process("test data")
			if err != nil {
				t.Errorf("Expected no error during processing, got %v", err)
			}

			if result == "" {
				t.Error("Expected non-empty encrypted result")
			}

			// Each encryption should produce different output (due to random nonce)
			result2, _ := decorated.Process("test data")
			if result == result2 {
				t.Error("Expected different encrypted outputs for same input (nonce should be random)")
			}
		})
	}
}

func TestCachingDecorator(t *testing.T) {
	base := NewSimpleProcessor("Test")
	cache := NewCachingDecorator(base)

	// First call should be a cache miss
	result1, err := cache.Process("test data")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	hits, misses := cache.Stats()
	if hits != 0 || misses != 1 {
		t.Errorf("Expected 0 hits and 1 miss, got %d hits and %d misses", hits, misses)
	}

	// Second call with same data should be a cache hit
	result2, err := cache.Process("test data")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result1 != result2 {
		t.Error("Expected cached result to match original result")
	}

	hits, misses = cache.Stats()
	if hits != 1 || misses != 1 {
		t.Errorf("Expected 1 hit and 1 miss, got %d hits and %d misses", hits, misses)
	}

	// Call with different data should be a cache miss
	_, err = cache.Process("different data")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	hits, misses = cache.Stats()
	if hits != 1 || misses != 2 {
		t.Errorf("Expected 1 hit and 2 misses, got %d hits and %d misses", hits, misses)
	}

	// Clear cache
	cache.ClearCache()
	hits, misses = cache.Stats()
	if hits != 0 || misses != 0 {
		t.Errorf("Expected cache to be cleared, got %d hits and %d misses", hits, misses)
	}
}

// ============================================================================
// Decorator Composition Tests
// ============================================================================

func TestDecoratorStacking(t *testing.T) {
	base := NewSimpleProcessor("Test")

	// Stack multiple decorators
	decorated := NewValidationDecorator(
		base,
		ValidationConfig{MinLength: 5, MaxLength: 100, AllowEmpty: false},
	)

	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)
	decorated = NewLoggingDecorator(decorated, logger)

	cache := NewCachingDecorator(decorated)

	// Test valid data
	result, err := cache.Process("valid data")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Check that logging occurred
	if logBuf.Len() == 0 {
		t.Error("Expected logging output")
	}

	// Test invalid data
	_, err = cache.Process("bad")
	if err == nil {
		t.Error("Expected validation error for short data")
	}
}

func TestDecoratorOrderMatters(t *testing.T) {
	base := NewSimpleProcessor("Test")
	key := make([]byte, 32)

	// Order 1: Compress then Encrypt
	comp1 := NewCompressionDecorator(base, 6)
	enc1, _ := NewEncryptionDecorator(comp1, key)
	result1, _ := enc1.Process("This is some test data that will be processed.")
	len1 := len(result1)

	// Order 2: Encrypt then Compress
	enc2, _ := NewEncryptionDecorator(base, key)
	comp2 := NewCompressionDecorator(enc2, 6)
	result2, _ := comp2.Process("This is some test data that will be processed.")
	len2 := len(result2)

	// Compressing then encrypting should generally produce smaller output
	// than encrypting then compressing (since compression works better on plaintext)
	if len1 >= len2 {
		t.Logf("Warning: Expected compress-then-encrypt (%d bytes) to be smaller than encrypt-then-compress (%d bytes)", len1, len2)
		// Note: This might not always be true for very small data, but it's a general principle
	}
}

// ============================================================================
// HTTP Middleware Tests
// ============================================================================

func TestLoggingMiddleware(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	middleware := LoggingMiddleware(logger)
	decorated := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	decorated(w, req)

	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "Request started") {
		t.Error("Expected log to contain 'Request started'")
	}
	if !strings.Contains(logOutput, "Request completed") {
		t.Error("Expected log to contain 'Request completed'")
	}
}

func TestAuthMiddleware(t *testing.T) {
	validToken := "Bearer secret-token"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}

	middleware := AuthMiddleware(validToken)
	decorated := middleware(handler)

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{"valid token", validToken, http.StatusOK},
		{"invalid token", "Bearer wrong-token", http.StatusUnauthorized},
		{"no token", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			decorated(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestMetricsMiddleware(t *testing.T) {
	metrics := &RequestMetrics{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}

	middleware := MetricsMiddleware(metrics)
	decorated := middleware(handler)

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	decorated(w, req)

	if metrics.TotalRequests != 1 {
		t.Errorf("Expected 1 request, got %d", metrics.TotalRequests)
	}

	if metrics.TotalErrors != 0 {
		t.Errorf("Expected 0 errors, got %d", metrics.TotalErrors)
	}

	if metrics.TotalDuration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", metrics.TotalDuration)
	}

	// Make an error request
	errorHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	decorated = middleware(errorHandler)

	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	decorated(w, req)

	if metrics.TotalRequests != 2 {
		t.Errorf("Expected 2 requests, got %d", metrics.TotalRequests)
	}

	if metrics.TotalErrors != 1 {
		t.Errorf("Expected 1 error, got %d", metrics.TotalErrors)
	}
}

func TestChainMiddleware(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)
	metrics := &RequestMetrics{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	// Chain multiple middleware
	chained := ChainMiddleware(
		handler,
		LoggingMiddleware(logger),
		AuthMiddleware("Bearer token"),
		MetricsMiddleware(metrics),
	)

	// Test with valid auth
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	chained(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if metrics.TotalRequests != 1 {
		t.Errorf("Expected 1 request recorded, got %d", metrics.TotalRequests)
	}

	if !strings.Contains(logBuf.String(), "Request started") {
		t.Error("Expected logging to occur")
	}

	// Test with invalid auth (should be blocked by auth middleware)
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	chained(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkSimpleProcessor(b *testing.B) {
	processor := NewSimpleProcessor("Benchmark")
	data := "benchmark data"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.Process(data)
	}
}

func BenchmarkLoggingDecorator(b *testing.B) {
	base := NewSimpleProcessor("Benchmark")
	processor := NewLoggingDecorator(base, log.New(bytes.NewBuffer(nil), "", 0))
	data := "benchmark data"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.Process(data)
	}
}

func BenchmarkCachingDecorator_Hit(b *testing.B) {
	base := NewSimpleProcessor("Benchmark")
	cache := NewCachingDecorator(base)
	data := "benchmark data"

	// Prime the cache
	cache.Process(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Process(data)
	}
}

func BenchmarkCachingDecorator_Miss(b *testing.B) {
	base := NewSimpleProcessor("Benchmark")
	cache := NewCachingDecorator(base)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Process(string(rune(i)))
	}
}

func BenchmarkCompressionDecorator(b *testing.B) {
	base := NewSimpleProcessor("Benchmark")
	processor := NewCompressionDecorator(base, 6)
	data := "This is some benchmark data that will be compressed."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.Process(data)
	}
}

func BenchmarkEncryptionDecorator(b *testing.B) {
	base := NewSimpleProcessor("Benchmark")
	key := make([]byte, 32)
	processor, _ := NewEncryptionDecorator(base, key)
	data := "benchmark data for encryption"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.Process(data)
	}
}

func BenchmarkStackedDecorators(b *testing.B) {
	base := NewSimpleProcessor("Benchmark")
	processor := NewLoggingDecorator(
		NewValidationDecorator(
			NewCachingDecorator(base),
			ValidationConfig{MinLength: 1, MaxLength: 1000, AllowEmpty: false},
		),
		log.New(bytes.NewBuffer(nil), "", 0),
	)
	data := "benchmark data"

	// Prime cache
	processor.Process(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.Process(data)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func isBase64(s string) bool {
	// Simple check: base64 strings only contain A-Z, a-z, 0-9, +, /, and =
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '+' || c == '/' || c == '=') {
			return false
		}
	}
	return len(s) > 0
}
