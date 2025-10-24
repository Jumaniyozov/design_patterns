package pipeline

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Example1_BasicPipeline demonstrates a simple linear pipeline.
func Example1_BasicPipeline() {
	fmt.Println("=== Example 1: Basic Pipeline ===")
	fmt.Println("Square numbers, then double them, then filter evens")
	fmt.Println()

	ctx := context.Background()

	// Generate numbers 1-5
	numbers := Generator(ctx, 1, 2, 3, 4, 5)

	// Stage 1: Square each number
	squared := Map(ctx, numbers, func(n int) int {
		result := n * n
		fmt.Printf("Square: %d -> %d\n", n, result)
		return result
	})

	// Stage 2: Double each squared number
	doubled := Map(ctx, squared, func(n int) int {
		result := n * 2
		fmt.Printf("Double: %d -> %d\n", n, result)
		return result
	})

	// Stage 3: Filter for even numbers (all will be even in this case)
	filtered := Filter(ctx, doubled, func(n int) bool {
		isEven := n%2 == 0
		fmt.Printf("Filter: %d -> %v\n", n, isEven)
		return isEven
	})

	// Collect results
	fmt.Println("\nFinal results:")
	for result := range filtered {
		fmt.Printf("  %d\n", result)
	}
	fmt.Println()
}

// Example2_FluentPipeline demonstrates using the fluent Pipeline builder.
func Example2_FluentPipeline() {
	fmt.Println("=== Example 2: Fluent Pipeline Builder ===")
	fmt.Println("Using method chaining for cleaner syntax")
	fmt.Println()

	ctx := context.Background()

	// Create source channel
	source := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Build pipeline with method chaining
	results := NewPipeline(ctx, source).
		Map(func(n int) int {
			return n * n // square
		}).
		Filter(func(n int) bool {
			return n > 20 // keep only values > 20
		}).
		Take(3). // take first 3 results
		Collect()

	fmt.Println("Input: 1-10")
	fmt.Println("Square, filter (>20), take(3):")
	for _, v := range results {
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()
}

// Example3_FanOutFanIn demonstrates parallel processing with multiple workers.
func Example3_FanOutFanIn() {
	fmt.Println("=== Example 3: Fan-Out/Fan-In Parallel Processing ===")
	fmt.Println("Using 3 workers to process data in parallel")
	fmt.Println()

	ctx := context.Background()

	// Generate work items
	numbers := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	// Fan-out: Create 3 workers that process numbers in parallel
	workerCount := 3
	expensiveWork := func(n int) int {
		// Simulate expensive computation
		time.Sleep(100 * time.Millisecond)
		result := n * n
		fmt.Printf("Worker processed: %d -> %d\n", n, result)
		return result
	}

	workerChannels := FanOut(ctx, numbers, workerCount, expensiveWork)

	// Fan-in: Merge results from all workers
	results := FanIn(ctx, workerChannels...)

	// Collect results
	fmt.Println("\nResults (order may vary due to parallel processing):")
	var collected []int
	for result := range results {
		collected = append(collected, result)
	}

	fmt.Printf("  %v\n", collected)
	fmt.Println()
}

// Example4_ErrorHandling demonstrates error handling in pipelines.
func Example4_ErrorHandling() {
	fmt.Println("=== Example 4: Error Handling ===")
	fmt.Println("Processing strings to integers with error handling")
	fmt.Println()

	ctx := context.Background()

	// Generate input with some invalid values
	inputs := Generator(ctx, "1", "2", "invalid", "4", "bad", "6")

	// Process with error handling
	results := MapWithError(ctx, inputs, func(s string) (int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("failed to parse '%s': %w", s, err)
		}
		return n, nil
	})

	// Separate values from errors
	values, errors := CollectErrors(ctx, results)

	// Process values and errors separately
	var successCount, errorCount int

	// Collect in goroutines to avoid deadlock
	done := make(chan bool, 2)

	go func() {
		for v := range values {
			fmt.Printf("Success: %d\n", v)
			successCount++
		}
		done <- true
	}()

	go func() {
		for err := range errors {
			fmt.Printf("Error: %v\n", err)
			errorCount++
		}
		done <- true
	}()

	// Wait for both to complete
	<-done
	<-done

	fmt.Printf("\nSummary: %d successful, %d errors\n", successCount, errorCount)
	fmt.Println()
}

// Example5_RealWorldLogProcessing demonstrates a realistic log processing pipeline.
func Example5_RealWorldLogProcessing() {
	fmt.Println("=== Example 5: Real-World Log Processing ===")
	fmt.Println("Parse -> Filter -> Transform -> Aggregate")
	fmt.Println()

	ctx := context.Background()

	// Simulate log entries
	logEntries := []string{
		"2024-01-15 10:00:00 INFO User login successful",
		"2024-01-15 10:00:05 ERROR Database connection failed",
		"2024-01-15 10:00:10 INFO Request processed in 45ms",
		"2024-01-15 10:00:15 ERROR Invalid authentication token",
		"2024-01-15 10:00:20 WARN High memory usage detected",
		"2024-01-15 10:00:25 ERROR Service timeout after 30s",
		"2024-01-15 10:00:30 INFO Cache refreshed",
	}

	// Stage 1: Generate log entries
	logs := Generator(ctx, logEntries...)

	// Stage 2: Parse log entries
	type LogEntry struct {
		Timestamp string
		Level     string
		Message   string
	}

	parsed := Map(ctx, logs, func(log string) LogEntry {
		parts := strings.SplitN(log, " ", 4)
		if len(parts) < 4 {
			return LogEntry{}
		}
		return LogEntry{
			Timestamp: parts[0] + " " + parts[1],
			Level:     parts[2],
			Message:   parts[3],
		}
	})

	// Stage 3: Filter for ERROR level only
	errors := Filter(ctx, parsed, func(entry LogEntry) bool {
		return entry.Level == "ERROR"
	})

	// Stage 4: Transform to alert format
	alerts := Map(ctx, errors, func(entry LogEntry) string {
		return fmt.Sprintf("🚨 ALERT [%s]: %s", entry.Timestamp, entry.Message)
	})

	// Stage 5: Output alerts
	fmt.Println("Processing logs...")
	fmt.Println("\nGenerated Alerts:")
	alertCount := 0
	for alert := range alerts {
		fmt.Println(alert)
		alertCount++
	}

	fmt.Printf("\nTotal alerts generated: %d\n", alertCount)
	fmt.Println()
}

// Example6_DataBatching demonstrates batching for efficiency.
func Example6_DataBatching() {
	fmt.Println("=== Example 6: Batching for Efficient Processing ===")
	fmt.Println("Group individual items into batches of 3")
	fmt.Println()

	ctx := context.Background()

	// Generate numbers
	numbers := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Batch into groups of 3
	batched := Batch(ctx, numbers, 3)

	// Process batches
	fmt.Println("Batched results:")
	batchNum := 1
	for batch := range batched {
		fmt.Printf("Batch %d: %v\n", batchNum, batch)
		batchNum++
	}
	fmt.Println()
}

// Example7_PipelineCancellation demonstrates cancellation with context.
func Example7_PipelineCancellation() {
	fmt.Println("=== Example 7: Pipeline Cancellation ===")
	fmt.Println("Cancel pipeline after 200ms")
	fmt.Println()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Infinite generator (would run forever without cancellation)
	infinite := make(chan int)
	go func() {
		defer close(infinite)
		for i := 0; ; i++ {
			select {
			case infinite <- i:
				time.Sleep(50 * time.Millisecond) // Simulate slow generation
			case <-ctx.Done():
				return
			}
		}
	}()

	// Process with pipeline
	results := Map(ctx, infinite, func(n int) int {
		return n * 2
	})

	// Consume until cancelled
	fmt.Println("Processing (will stop after 200ms):")
	for v := range results {
		fmt.Printf("  Received: %d\n", v)
	}

	fmt.Println("Pipeline cancelled successfully")
	fmt.Println()
}

// Example8_ComplexDataTransformation demonstrates a multi-stage transformation pipeline.
func Example8_ComplexDataTransformation() {
	fmt.Println("=== Example 8: Complex Data Transformation Pipeline ===")
	fmt.Println("Simulating image processing: Load -> Resize -> Filter -> Compress")
	fmt.Println()

	ctx := context.Background()

	type Image struct {
		Name   string
		Width  int
		Height int
		Size   int // in KB
	}

	// Stage 1: Generate images
	images := []Image{
		{Name: "photo1.jpg", Width: 4000, Height: 3000, Size: 5000},
		{Name: "photo2.jpg", Width: 3840, Height: 2160, Size: 4200},
		{Name: "photo3.jpg", Width: 1920, Height: 1080, Size: 2500},
	}
	source := Generator(ctx, images...)

	// Stage 2: Resize large images
	resized := Map(ctx, source, func(img Image) Image {
		if img.Width > 1920 {
			ratio := 1920.0 / float64(img.Width)
			img.Width = 1920
			img.Height = int(float64(img.Height) * ratio)
			img.Size = int(float64(img.Size) * ratio * ratio)
			fmt.Printf("Resized: %s -> %dx%d (%dKB)\n", img.Name, img.Width, img.Height, img.Size)
		} else {
			fmt.Printf("Skipped resize: %s already optimal\n", img.Name)
		}
		return img
	})

	// Stage 3: Apply filter (simulate processing time)
	filtered := Map(ctx, resized, func(img Image) Image {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Applied filter: %s\n", img.Name)
		return img
	})

	// Stage 4: Compress
	compressed := Map(ctx, filtered, func(img Image) Image {
		img.Size = int(float64(img.Size) * 0.7) // 30% compression
		fmt.Printf("Compressed: %s -> %dKB\n", img.Name, img.Size)
		return img
	})

	// Stage 5: Collect final results
	fmt.Println("\nFinal processed images:")
	for img := range compressed {
		fmt.Printf("  ✓ %s (%dx%d, %dKB)\n", img.Name, img.Width, img.Height, img.Size)
	}
	fmt.Println()
}

// Example9_TeePipeline demonstrates splitting a pipeline into two branches.
func Example9_TeePipeline() {
	fmt.Println("=== Example 9: Tee - Split Pipeline into Two Branches ===")
	fmt.Println()

	ctx := context.Background()

	// Generate numbers
	numbers := Generator(ctx, 1, 2, 3, 4, 5)

	// Split into two branches
	branch1, branch2 := Tee(ctx, numbers)

	// Branch 1: Calculate squares
	squares := Map(ctx, branch1, func(n int) int {
		return n * n
	})

	// Branch 2: Calculate cubes
	cubes := Map(ctx, branch2, func(n int) int {
		return n * n * n
	})

	// Process both branches concurrently
	done := make(chan bool, 2)

	go func() {
		fmt.Println("Branch 1 (squares):")
		for v := range squares {
			fmt.Printf("  %d\n", v)
		}
		done <- true
	}()

	go func() {
		fmt.Println("Branch 2 (cubes):")
		for v := range cubes {
			fmt.Printf("  %d\n", v)
		}
		done <- true
	}()

	<-done
	<-done
	fmt.Println()
}

// Example10_AdvancedETL demonstrates a realistic ETL pipeline with all stages.
func Example10_AdvancedETL() {
	fmt.Println("=== Example 10: Advanced ETL Pipeline ===")
	fmt.Println("Extract -> Validate -> Transform -> Enrich -> Load")
	fmt.Println()

	ctx := context.Background()

	type RawRecord struct {
		ID    int
		Email string
		Age   int
	}

	type ProcessedRecord struct {
		ID       int
		Email    string
		Age      int
		Domain   string
		AgeGroup string
		Valid    bool
	}

	// Stage 1: Extract - simulate reading from database
	rawData := []RawRecord{
		{ID: 1, Email: "user1@example.com", Age: 25},
		{ID: 2, Email: "user2@test.org", Age: 35},
		{ID: 3, Email: "invalid-email", Age: 45},
		{ID: 4, Email: "user4@company.net", Age: 17},
		{ID: 5, Email: "user5@example.com", Age: 65},
	}
	extracted := Generator(ctx, rawData...)

	// Stage 2: Validate
	validated := Map(ctx, extracted, func(r RawRecord) ProcessedRecord {
		valid := strings.Contains(r.Email, "@") && r.Age >= 18
		return ProcessedRecord{
			ID:    r.ID,
			Email: r.Email,
			Age:   r.Age,
			Valid: valid,
		}
	})

	// Stage 3: Filter invalid records
	validOnly := Filter(ctx, validated, func(r ProcessedRecord) bool {
		if !r.Valid {
			fmt.Printf("Rejected: ID=%d (invalid data)\n", r.ID)
		}
		return r.Valid
	})

	// Stage 4: Enrich with additional data
	enriched := Map(ctx, validOnly, func(r ProcessedRecord) ProcessedRecord {
		// Extract domain
		parts := strings.Split(r.Email, "@")
		if len(parts) == 2 {
			r.Domain = parts[1]
		}

		// Determine age group
		switch {
		case r.Age < 30:
			r.AgeGroup = "Young Adult"
		case r.Age < 50:
			r.AgeGroup = "Middle Age"
		default:
			r.AgeGroup = "Senior"
		}

		return r
	})

	// Stage 5: Load - simulate saving to data warehouse
	fmt.Println("Loading records to warehouse:")
	loadCount := 0
	for record := range enriched {
		// Simulate random processing delay
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		fmt.Printf("  ✓ Loaded ID=%d, Email=%s, Domain=%s, AgeGroup=%s\n",
			record.ID, record.Email, record.Domain, record.AgeGroup)
		loadCount++
	}

	fmt.Printf("\nETL Complete: %d records loaded\n", loadCount)
	fmt.Println()
}
