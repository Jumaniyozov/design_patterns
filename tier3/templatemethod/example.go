package templatemethod

import "fmt"

// Example1_DataProcessing demonstrates the template method with data processors.
func Example1_DataProcessing() {
	fmt.Println("=== Example 1: Data Processing Pipeline ===")

	data := []string{"apple", "banana", "cherry", "date"}

	// Process with CSV
	fmt.Println("--- CSV Processing ---")
	csv := NewCSVProcessor(data)
	ProcessTemplate(csv)
	fmt.Printf("Output: %v\n\n", csv.GetOutput())

	// Process with JSON
	fmt.Println("--- JSON Processing ---")
	json := NewJSONProcessor(data)
	ProcessTemplate(json)
	fmt.Printf("Output: %v\n\n", json.GetOutput())

	// Process with XML
	fmt.Println("--- XML Processing ---")
	xml := NewXMLProcessor(data)
	ProcessTemplate(xml)
	fmt.Printf("Output: %v\n\n", xml.GetOutput())
}

// Example2_ReportGeneration demonstrates report building with different formats.
func Example2_ReportGeneration() {
	fmt.Println("=== Example 2: Report Generation ===")

	data := []string{
		"Q1 Revenue: $1.2M",
		"Q2 Revenue: $1.5M",
		"Q3 Revenue: $1.8M",
		"Q4 Revenue: $2.1M",
	}

	// HTML Report
	fmt.Println("--- HTML Report ---")
	htmlBuilder := NewHTMLReportBuilder("Quarterly Revenue Report")
	htmlReport := BuildReport(htmlBuilder, data)
	fmt.Println(htmlReport)

	// Markdown Report
	fmt.Println("--- Markdown Report ---")
	mdBuilder := NewMarkdownReportBuilder("Quarterly Revenue Report")
	mdReport := BuildReport(mdBuilder, data)
	fmt.Println(mdReport)

	// Plain Text Report
	fmt.Println("--- Plain Text Report ---")
	txtBuilder := NewPlainTextReportBuilder("Quarterly Revenue Report")
	txtReport := BuildReport(txtBuilder, data)
	fmt.Println(txtReport)
}

// Example3_GameAI demonstrates AI behavior with template method.
func Example3_GameAI() {
	fmt.Println("=== Example 3: Game AI Update Cycle ===")

	// Zombie AI - player not nearby
	fmt.Println("--- Zombie AI (No Player) ---")
	zombie := NewZombieAI()
	zombie.SetPlayerNear(false)
	AIUpdateCycle(zombie)
	fmt.Println()

	// Zombie AI - player nearby
	fmt.Println("--- Zombie AI (Player Nearby) ---")
	zombie.SetPlayerNear(true)
	AIUpdateCycle(zombie)
	fmt.Println()

	// Guard AI - normal patrol
	fmt.Println("--- Guard AI (Normal Patrol) ---")
	guard := NewGuardAI()
	guard.SetAlertLevel(0)
	AIUpdateCycle(guard)
	fmt.Println()

	// Guard AI - suspicious activity
	fmt.Println("--- Guard AI (Suspicious Activity) ---")
	guard.SetAlertLevel(1)
	AIUpdateCycle(guard)
	fmt.Println()
}

// Example4_MultipleProcessors demonstrates processing the same data differently.
func Example4_MultipleProcessors() {
	fmt.Println("=== Example 4: Same Data, Different Processors ===")

	input := []string{"user1", "user2", "user3"}

	processors := []struct {
		name      string
		processor DataProcessor
	}{
		{"CSV", NewCSVProcessor(input)},
		{"JSON", NewJSONProcessor(input)},
		{"XML", NewXMLProcessor(input)},
	}

	for _, p := range processors {
		fmt.Printf("--- %s Processor ---\n", p.name)
		ProcessTemplate(p.processor)
		fmt.Println()
	}
}
