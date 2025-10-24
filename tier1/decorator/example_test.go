package decorator_test

import (
	"testing"

	"github.com/jumaniyozov/design_patterns/tier1/decorator"
)

// TestRunAllExamples demonstrates how to run all examples.
// This is useful for seeing all examples in action during testing.
func TestRunAllExamples(t *testing.T) {
	// This will output all examples to stdout
	// Run with: go test -v ./tier1/decorator -run TestRunAllExamples
	decorator.RunAllExamples()
}

// Example test that can be run with `go test`
func ExampleSimpleProcessor_Process() {
	processor := decorator.NewSimpleProcessor("ExampleProcessor")
	result, _ := processor.Process("Hello, Decorator Pattern!")
	// Output is deterministic so we can verify it
	_ = result
}
