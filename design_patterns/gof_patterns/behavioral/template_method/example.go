package template_method

import "fmt"

func Example1_DataProcessors() {
	fmt.Println("\n=== Example 1: Data Processors ===")

	csvProc := &AbstractDataProcessor{processor: &CSVProcessor{}}
	csvProc.Process()

	fmt.Println()
	xmlProc := &AbstractDataProcessor{processor: &XMLProcessor{}}
	xmlProc.Process()
}
