// Package decorator demonstrates the Decorator Pattern in Go.
//
// The Decorator Pattern allows behavior to be added to individual objects,
// either statically or dynamically, without affecting the behavior of other
// objects from the same class. It's particularly elegant in Go thanks to
// first-class functions and interface composition.
package decorator

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ============================================================================
// Object-Oriented Decorator Pattern
// ============================================================================

// DataProcessor defines the component interface that both concrete components
// and decorators must implement. This is the contract that allows decorators
// to wrap other processors transparently.
type DataProcessor interface {
	Process(data string) (string, error)
}

// ============================================================================
// Concrete Component
// ============================================================================

// SimpleProcessor is the concrete component - the base object we want to decorate.
// It provides the core functionality without any additional features.
type SimpleProcessor struct {
	name string
}

// NewSimpleProcessor creates a new simple processor with the given name.
func NewSimpleProcessor(name string) *SimpleProcessor {
	return &SimpleProcessor{name: name}
}

// Process implements the DataProcessor interface with basic processing.
func (sp *SimpleProcessor) Process(data string) (string, error) {
	// Base processing: just format the data
	return fmt.Sprintf("[%s] %s", sp.name, data), nil
}

// ============================================================================
// Concrete Decorators
// ============================================================================

// LoggingDecorator adds logging functionality to any DataProcessor.
// It demonstrates how decorators wrap components and add behavior
// before and after delegating to the wrapped component.
type LoggingDecorator struct {
	processor DataProcessor
	logger    *log.Logger
}

// NewLoggingDecorator wraps a processor with logging capabilities.
func NewLoggingDecorator(processor DataProcessor, logger *log.Logger) DataProcessor {
	if logger == nil {
		logger = log.Default()
	}
	return &LoggingDecorator{
		processor: processor,
		logger:    logger,
	}
}

// Process adds logging before and after processing.
func (ld *LoggingDecorator) Process(data string) (string, error) {
	ld.logger.Printf("Processing started: input length = %d", len(data))
	start := time.Now()

	result, err := ld.processor.Process(data)

	duration := time.Since(start)
	if err != nil {
		ld.logger.Printf("Processing failed after %v: %v", duration, err)
		return "", err
	}

	ld.logger.Printf("Processing completed in %v: output length = %d", duration, len(result))
	return result, nil
}

// ============================================================================

// ValidationDecorator adds validation functionality.
// It demonstrates pre-processing validation and error handling.
type ValidationDecorator struct {
	processor  DataProcessor
	minLength  int
	maxLength  int
	allowEmpty bool
}

// ValidationConfig holds configuration for the validation decorator.
type ValidationConfig struct {
	MinLength  int
	MaxLength  int
	AllowEmpty bool
}

// NewValidationDecorator wraps a processor with validation capabilities.
func NewValidationDecorator(processor DataProcessor, config ValidationConfig) DataProcessor {
	return &ValidationDecorator{
		processor:  processor,
		minLength:  config.MinLength,
		maxLength:  config.MaxLength,
		allowEmpty: config.AllowEmpty,
	}
}

// Process adds validation before processing.
func (vd *ValidationDecorator) Process(data string) (string, error) {
	// Pre-processing validation
	if !vd.allowEmpty && len(data) == 0 {
		return "", fmt.Errorf("validation error: empty data not allowed")
	}

	if vd.minLength > 0 && len(data) < vd.minLength {
		return "", fmt.Errorf("validation error: data length %d is less than minimum %d", len(data), vd.minLength)
	}

	if vd.maxLength > 0 && len(data) > vd.maxLength {
		return "", fmt.Errorf("validation error: data length %d exceeds maximum %d", len(data), vd.maxLength)
	}

	// Delegate to wrapped processor
	return vd.processor.Process(data)
}

// ============================================================================

// CompressionDecorator adds compression functionality.
// It demonstrates post-processing transformation of data.
type CompressionDecorator struct {
	processor DataProcessor
	level     int
}

// NewCompressionDecorator wraps a processor with compression capabilities.
// level is the gzip compression level (1-9, or -1 for default).
func NewCompressionDecorator(processor DataProcessor, level int) DataProcessor {
	return &CompressionDecorator{
		processor: processor,
		level:     level,
	}
}

// Process adds compression after processing.
func (cd *CompressionDecorator) Process(data string) (string, error) {
	result, err := cd.processor.Process(data)
	if err != nil {
		return "", err
	}

	// Post-processing compression
	compressed, err := cd.compress(result)
	if err != nil {
		return "", fmt.Errorf("compression error: %w", err)
	}

	return compressed, nil
}

func (cd *CompressionDecorator) compress(data string) (string, error) {
	var buf strings.Builder
	gzWriter, err := gzip.NewWriterLevel(&buf, cd.level)
	if err != nil {
		return "", err
	}

	if _, err := gzWriter.Write([]byte(data)); err != nil {
		return "", err
	}

	if err := gzWriter.Close(); err != nil {
		return "", err
	}

	// Return base64 encoded for safe string representation
	return base64.StdEncoding.EncodeToString([]byte(buf.String())), nil
}

// ============================================================================

// EncryptionDecorator adds encryption functionality.
// This demonstrates how decorators can add security features transparently.
type EncryptionDecorator struct {
	processor DataProcessor
	key       []byte
}

// NewEncryptionDecorator wraps a processor with AES encryption.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func NewEncryptionDecorator(processor DataProcessor, key []byte) (DataProcessor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}

	return &EncryptionDecorator{
		processor: processor,
		key:       key,
	}, nil
}

// Process adds encryption after processing.
func (ed *EncryptionDecorator) Process(data string) (string, error) {
	result, err := ed.processor.Process(data)
	if err != nil {
		return "", err
	}

	// Post-processing encryption
	encrypted, err := ed.encrypt(result)
	if err != nil {
		return "", fmt.Errorf("encryption error: %w", err)
	}

	return encrypted, nil
}

func (ed *EncryptionDecorator) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(ed.key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and append nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return base64 encoded for safe string representation
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// ============================================================================

// CachingDecorator adds caching functionality.
// This demonstrates stateful decorators that maintain their own data.
type CachingDecorator struct {
	processor DataProcessor
	cache     map[string]string
	hits      int
	misses    int
}

// NewCachingDecorator wraps a processor with caching capabilities.
func NewCachingDecorator(processor DataProcessor) *CachingDecorator {
	return &CachingDecorator{
		processor: processor,
		cache:     make(map[string]string),
	}
}

// Process checks cache before delegating to the wrapped processor.
func (cd *CachingDecorator) Process(data string) (string, error) {
	// Check cache first
	if result, found := cd.cache[data]; found {
		cd.hits++
		return result, nil
	}

	// Cache miss - delegate to wrapped processor
	cd.misses++
	result, err := cd.processor.Process(data)
	if err != nil {
		return "", err
	}

	// Store in cache
	cd.cache[data] = result
	return result, nil
}

// Stats returns cache statistics.
func (cd *CachingDecorator) Stats() (hits, misses int) {
	return cd.hits, cd.misses
}

// ClearCache clears the cache and resets statistics.
func (cd *CachingDecorator) ClearCache() {
	cd.cache = make(map[string]string)
	cd.hits = 0
	cd.misses = 0
}

// ============================================================================
// Functional Decorator Pattern (Go Idiomatic)
// ============================================================================

// HTTPHandler is a standard HTTP handler function type.
type HTTPHandler func(w http.ResponseWriter, r *http.Request)

// HTTPMiddleware is a decorator for HTTP handlers.
// This demonstrates the functional decorator pattern which is idiomatic in Go.
type HTTPMiddleware func(HTTPHandler) HTTPHandler

// LoggingMiddleware creates a logging decorator for HTTP handlers.
func LoggingMiddleware(logger *log.Logger) HTTPMiddleware {
	return func(next HTTPHandler) HTTPHandler {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Printf("Request started: %s %s", r.Method, r.URL.Path)

			// Call the wrapped handler
			next(w, r)

			duration := time.Since(start)
			logger.Printf("Request completed: %s %s in %v", r.Method, r.URL.Path, duration)
		}
	}
}

// AuthMiddleware creates an authentication decorator for HTTP handlers.
func AuthMiddleware(validToken string) HTTPMiddleware {
	return func(next HTTPHandler) HTTPHandler {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != validToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Authentication passed - call next handler
			next(w, r)
		}
	}
}

// MetricsMiddleware creates a metrics collection decorator.
func MetricsMiddleware(metrics *RequestMetrics) HTTPMiddleware {
	return func(next HTTPHandler) HTTPHandler {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			metrics.IncrementRequests()

			// Wrap ResponseWriter to capture status code
			wrappedWriter := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next(wrappedWriter, r)

			duration := time.Since(start)
			metrics.RecordDuration(duration)

			if wrappedWriter.statusCode >= 400 {
				metrics.IncrementErrors()
			}
		}
	}
}

// RequestMetrics holds metrics for HTTP requests.
type RequestMetrics struct {
	TotalRequests int
	TotalErrors   int
	TotalDuration time.Duration
}

// IncrementRequests increments the request counter.
func (rm *RequestMetrics) IncrementRequests() {
	rm.TotalRequests++
}

// IncrementErrors increments the error counter.
func (rm *RequestMetrics) IncrementErrors() {
	rm.TotalErrors++
}

// RecordDuration adds to the total duration.
func (rm *RequestMetrics) RecordDuration(d time.Duration) {
	rm.TotalDuration += d
}

// AverageDuration returns the average request duration.
func (rm *RequestMetrics) AverageDuration() time.Duration {
	if rm.TotalRequests == 0 {
		return 0
	}
	return rm.TotalDuration / time.Duration(rm.TotalRequests)
}

// statusRecorder wraps http.ResponseWriter to capture the status code.
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

// ChainMiddleware chains multiple middleware decorators together.
// This is a utility function that demonstrates how to compose decorators.
func ChainMiddleware(handler HTTPHandler, middlewares ...HTTPMiddleware) HTTPHandler {
	// Apply middlewares in reverse order so they execute in the order specified
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
