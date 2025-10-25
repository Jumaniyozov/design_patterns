// Package module demonstrates the Module pattern.
// In Go, the package system naturally implements the module pattern through
// exported (public) and unexported (private) identifiers.
package module

import "fmt"

// Public API - Exported types and functions

// Counter is a public counter with private state
type Counter struct {
	value int // unexported - private
}

// NewCounter creates a new counter (factory function)
func NewCounter() *Counter {
	return &Counter{value: 0}
}

// Increment increments the counter (exported - public)
func (c *Counter) Increment() {
	c.value++
}

// Value returns the current value (exported - public)
func (c *Counter) Value() int {
	return c.value
}

// Private implementation - unexported

// internalValidator is private to the package
type internalValidator struct {
	rules []validationRule
}

// validationRule is private
type validationRule func(string) bool

// newValidator is private constructor
func newValidator() *internalValidator {
	return &internalValidator{
		rules: []validationRule{
			isNotEmpty,
			hasMinLength,
		},
	}
}

// isNotEmpty is a private validation rule
func isNotEmpty(s string) bool {
	return len(s) > 0
}

// hasMinLength is a private validation rule
func hasMinLength(s string) bool {
	return len(s) >= 3
}

// validate is private validation method
func (iv *internalValidator) validate(s string) bool {
	for _, rule := range iv.rules {
		if !rule(s) {
			return false
		}
	}
	return true
}

// Public Validator wraps private implementation
type Validator struct {
	impl *internalValidator
}

// NewValidator creates a validator (public API)
func NewValidator() *Validator {
	return &Validator{
		impl: newValidator(),
	}
}

// Validate validates a string (public API)
func (v *Validator) Validate(s string) bool {
	return v.impl.validate(s)
}

// Configuration module example

// config holds private configuration
var config *configuration

// configuration is private
type configuration struct {
	settings map[string]string
}

// init initializes private config
func init() {
	config = &configuration{
		settings: map[string]string{
			"version": "1.0.0",
			"env":     "development",
		},
	}
}

// GetConfig returns a config value (public API)
func GetConfig(key string) string {
	return config.settings[key]
}

// SetConfig sets a config value (public API)
func SetConfig(key, value string) {
	config.settings[key] = value
}

// Logger module with private state

// logger is a private singleton instance
var logger *loggerImpl

// loggerImpl is private implementation
type loggerImpl struct {
	prefix string
	buffer []string
}

// initLogger initializes the private logger
func initLogger() {
	if logger == nil {
		logger = &loggerImpl{
			prefix: "[APP]",
			buffer: make([]string, 0),
		}
	}
}

// Log logs a message (public API)
func Log(message string) {
	initLogger()
	formatted := fmt.Sprintf("%s %s", logger.prefix, message)
	logger.buffer = append(logger.buffer, formatted)
	fmt.Println(formatted)
}

// GetLogs returns all logs (public API)
func GetLogs() []string {
	initLogger()
	return append([]string{}, logger.buffer...)
}

// SetLogPrefix sets the log prefix (public API)
func SetLogPrefix(prefix string) {
	initLogger()
	logger.prefix = prefix
}

// Math module demonstrating utility functions

// private constants
const (
	pi = 3.14159
	e  = 2.71828
)

// private helper
func square(x float64) float64 {
	return x * x
}

// CircleArea calculates circle area (public API)
func CircleArea(radius float64) float64 {
	return pi * square(radius)
}

// CircleCircumference calculates circumference (public API)
func CircleCircumference(radius float64) float64 {
	return 2 * pi * radius
}
