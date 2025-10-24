package singleton_test

import (
	"testing"

	"github.com/jumaniyozov/design_patterns/tier1/singleton"
)

// TestRunAllExamples executes all examples to verify they work correctly.
// Run with: go test -v -run TestRunAllExamples ./tier1/singleton/
func TestRunAllExamples(t *testing.T) {
	// This test runs all examples in sequence
	// It serves as both a test and a demonstration
	singleton.RunAllExamples()
}

// ExampleRunAllExamples demonstrates how to run all singleton examples.
// This appears in the package documentation.
func ExampleRunAllExamples() {
	singleton.RunAllExamples()
}
