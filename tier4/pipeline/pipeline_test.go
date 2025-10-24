package pipeline

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	ctx := context.Background()
	values := []int{1, 2, 3, 4, 5}
	ch := Generator(ctx, values...)

	var results []int
	for v := range ch {
		results = append(results, v)
	}

	if len(results) != len(values) {
		t.Errorf("Expected %d values, got %d", len(values), len(results))
	}

	for i, v := range results {
		if v != values[i] {
			t.Errorf("Expected value %d at index %d, got %d", values[i], i, v)
		}
	}
}

func TestMap(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3, 4, 5)

	double := func(n int) int { return n * 2 }
	output := Map(ctx, input, double)

	expected := []int{2, 4, 6, 8, 10}
	var results []int
	for v := range output {
		results = append(results, v)
	}

	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFilter(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	isEven := func(n int) bool { return n%2 == 0 }
	output := Filter(ctx, input, isEven)

	var results []int
	for v := range output {
		results = append(results, v)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 even numbers, got %d", len(results))
	}

	for _, v := range results {
		if v%2 != 0 {
			t.Errorf("Expected even number, got odd: %d", v)
		}
	}
}

func TestFanOutFanIn(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3, 4, 5, 6)

	// Fan out to 3 workers
	square := func(n int) int { return n * n }
	workers := FanOut(ctx, input, 3, square)

	// Fan in results
	merged := FanIn(ctx, workers...)

	// Collect results (order may vary)
	var results []int
	for v := range merged {
		results = append(results, v)
	}

	// Sort for comparison
	sort.Ints(results)

	expected := []int{1, 4, 9, 16, 25, 36}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestTee(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3)

	out1, out2 := Tee(ctx, input)

	var results1, results2 []int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range out1 {
			results1 = append(results1, v)
		}
	}()

	go func() {
		defer wg.Done()
		for v := range out2 {
			results2 = append(results2, v)
		}
	}()

	wg.Wait()

	if len(results1) != 3 || len(results2) != 3 {
		t.Errorf("Expected 3 values in each branch, got %d and %d", len(results1), len(results2))
	}

	expected := []int{1, 2, 3}
	for i := range expected {
		if results1[i] != expected[i] || results2[i] != expected[i] {
			t.Errorf("Mismatch at index %d", i)
		}
	}
}

func TestBatch(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3, 4, 5, 6, 7)

	batched := Batch(ctx, input, 3)

	var batches [][]int
	for batch := range batched {
		batches = append(batches, batch)
	}

	if len(batches) != 3 {
		t.Errorf("Expected 3 batches, got %d", len(batches))
	}

	// First two batches should have 3 items
	if len(batches[0]) != 3 || len(batches[1]) != 3 {
		t.Errorf("Expected first two batches to have 3 items")
	}

	// Last batch should have 1 item
	if len(batches[2]) != 1 {
		t.Errorf("Expected last batch to have 1 item, got %d", len(batches[2]))
	}
}

func TestMapWithError(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 0, 4) // 0 will cause an error

	divide10 := func(n int) (int, error) {
		if n == 0 {
			return 0, errors.New("division by zero")
		}
		return 10 / n, nil
	}

	results := MapWithError(ctx, input, divide10)

	successCount := 0
	errorCount := 0

	for result := range results {
		if result.Error != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	if successCount != 3 {
		t.Errorf("Expected 3 successful operations, got %d", successCount)
	}

	if errorCount != 1 {
		t.Errorf("Expected 1 error, got %d", errorCount)
	}
}

func TestCollectErrors(t *testing.T) {
	ctx := context.Background()

	// Create results channel with mixed success/error
	results := make(chan Result[int])
	go func() {
		defer close(results)
		results <- Result[int]{Value: 10, Error: nil}
		results <- Result[int]{Value: 0, Error: errors.New("error1")}
		results <- Result[int]{Value: 20, Error: nil}
		results <- Result[int]{Value: 0, Error: errors.New("error2")}
	}()

	values, errs := CollectErrors(ctx, results)

	var successValues []int
	var errorMsgs []error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range values {
			successValues = append(successValues, v)
		}
	}()

	go func() {
		defer wg.Done()
		for e := range errs {
			errorMsgs = append(errorMsgs, e)
		}
	}()

	wg.Wait()

	if len(successValues) != 2 {
		t.Errorf("Expected 2 success values, got %d", len(successValues))
	}

	if len(errorMsgs) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errorMsgs))
	}
}

func TestBuffer(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3, 4, 5)

	buffered := Buffer(ctx, input, 10)

	var results []int
	for v := range buffered {
		results = append(results, v)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 values, got %d", len(results))
	}
}

func TestTake(t *testing.T) {
	ctx := context.Background()
	input := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	limited := Take(ctx, input, 3)

	var results []int
	for v := range limited {
		results = append(results, v)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 values, got %d", len(results))
	}
}

func TestPipelineBuilder(t *testing.T) {
	ctx := context.Background()
	source := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	results := NewPipeline(ctx, source).
		Map(func(n int) int { return n * n }).
		Filter(func(n int) bool { return n > 20 }).
		Take(3).
		Collect()

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Results should be 25, 36, 49 (squares > 20, first 3)
	expected := []int{25, 36, 49}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create an infinite generator
	infinite := make(chan int)
	go func() {
		defer close(infinite)
		for i := 0; ; i++ {
			select {
			case infinite <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Process through pipeline
	processed := Map(ctx, infinite, func(n int) int { return n * 2 })

	// Take a few values then cancel
	count := 0
	for v := range processed {
		count++
		if count >= 3 {
			cancel()
			break
		}
		_ = v // Use the value
	}

	// Give goroutines time to shutdown
	time.Sleep(50 * time.Millisecond)

	if count != 3 {
		t.Errorf("Expected to receive 3 values before cancellation, got %d", count)
	}
}

func TestConcurrentSafety(t *testing.T) {
	ctx := context.Background()

	// Generate many values
	values := make([]int, 100)
	for i := 0; i < 100; i++ {
		values[i] = i
	}
	input := Generator(ctx, values...)

	// Process through multiple stages with fan-out
	workers := FanOut(ctx, input, 5, func(n int) int { return n * 2 })
	merged := FanIn(ctx, workers...)

	// Collect all results
	results := make(map[int]bool)
	for v := range merged {
		if results[v] {
			t.Errorf("Duplicate value received: %d", v)
		}
		results[v] = true
	}

	if len(results) != 100 {
		t.Errorf("Expected 100 unique results, got %d", len(results))
	}
}

func TestEmptyPipeline(t *testing.T) {
	ctx := context.Background()
	empty := Generator[int](ctx) // No values

	results := NewPipeline(ctx, empty).
		Map(func(n int) int { return n * 2 }).
		Filter(func(n int) bool { return n > 0 }).
		Collect()

	if len(results) != 0 {
		t.Errorf("Expected 0 results from empty pipeline, got %d", len(results))
	}
}

// Benchmark tests
func BenchmarkSimplePipeline(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		values := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			values[j] = j
		}

		input := Generator(ctx, values...)
		squared := Map(ctx, input, func(n int) int { return n * n })
		filtered := Filter(ctx, squared, func(n int) bool { return n%2 == 0 })

		// Consume results
		for range filtered {
		}
	}
}

func BenchmarkFanOutFanIn(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		values := make([]int, 100)
		for j := 0; j < 100; j++ {
			values[j] = j
		}

		input := Generator(ctx, values...)
		workers := FanOut(ctx, input, 4, func(n int) int { return n * n })
		merged := FanIn(ctx, workers...)

		// Consume results
		for range merged {
		}
	}
}

func BenchmarkPipelineBuilder(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		values := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			values[j] = j
		}

		source := Generator(ctx, values...)
		_ = NewPipeline(ctx, source).
			Map(func(n int) int { return n * n }).
			Filter(func(n int) bool { return n < 500000 }).
			Take(100).
			Collect()
	}
}
