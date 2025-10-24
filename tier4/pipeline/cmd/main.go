package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jumaniyozov/design_patterns/tier4/pipeline"
)

func main() {
	examples := map[string]func(){
		"1":  pipeline.Example1_BasicPipeline,
		"2":  pipeline.Example2_FluentPipeline,
		"3":  pipeline.Example3_FanOutFanIn,
		"4":  pipeline.Example4_ErrorHandling,
		"5":  pipeline.Example5_RealWorldLogProcessing,
		"6":  pipeline.Example6_DataBatching,
		"7":  pipeline.Example7_PipelineCancellation,
		"8":  pipeline.Example8_ComplexDataTransformation,
		"9":  pipeline.Example9_TeePipeline,
		"10": pipeline.Example10_AdvancedETL,
	}

	// If argument provided, run specific example
	if len(os.Args) > 1 {
		num := os.Args[1]
		if fn, ok := examples[num]; ok {
			fn()
			return
		}
		fmt.Printf("Unknown example: %s\n", num)
		printUsage()
		return
	}

	// Run all examples
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║         Pipeline Pattern - All Examples               ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()

	for i := 1; i <= 10; i++ {
		if fn, ok := examples[strconv.Itoa(i)]; ok {
			fn()
			fmt.Println("────────────────────────────────────────────────────────")
			fmt.Println()
		}
	}

	fmt.Println("✓ All examples completed successfully!")
}

func printUsage() {
	fmt.Println("Usage: go run main.go [example_number]")
	fmt.Println()
	fmt.Println("Available examples:")
	fmt.Println("  1  - Basic Pipeline")
	fmt.Println("  2  - Fluent Pipeline Builder")
	fmt.Println("  3  - Fan-Out/Fan-In Parallel Processing")
	fmt.Println("  4  - Error Handling")
	fmt.Println("  5  - Real-World Log Processing")
	fmt.Println("  6  - Data Batching")
	fmt.Println("  7  - Pipeline Cancellation")
	fmt.Println("  8  - Complex Data Transformation")
	fmt.Println("  9  - Tee - Split Pipeline")
	fmt.Println("  10 - Advanced ETL Pipeline")
	fmt.Println()
	fmt.Println("Run without arguments to execute all examples")
}
