package templatemethod

import (
	"strings"
	"testing"
)

func TestCSVProcessor(t *testing.T) {
	input := []string{"test1", "test2"}
	processor := NewCSVProcessor(input)

	err := ProcessTemplate(processor)
	if err != nil {
		t.Errorf("ProcessTemplate failed: %v", err)
	}

	output := processor.GetOutput()
	if len(output) != 2 {
		t.Errorf("Expected 2 output items, got %d", len(output))
	}

	if output[0] != "TEST1" || output[1] != "TEST2" {
		t.Errorf("Expected uppercase output, got %v", output)
	}
}

func TestJSONProcessor(t *testing.T) {
	input := []string{"value1", "value2"}
	processor := NewJSONProcessor(input)

	err := ProcessTemplate(processor)
	if err != nil {
		t.Errorf("ProcessTemplate failed: %v", err)
	}

	output := processor.GetOutput()
	if len(output) != 2 {
		t.Errorf("Expected 2 output items, got %d", len(output))
	}

	if !strings.Contains(output[0], "value1") {
		t.Errorf("Expected JSON formatted output, got %s", output[0])
	}
}

func TestXMLProcessor(t *testing.T) {
	input := []string{"item1", "item2"}
	processor := NewXMLProcessor(input)

	err := ProcessTemplate(processor)
	if err != nil {
		t.Errorf("ProcessTemplate failed: %v", err)
	}

	output := processor.GetOutput()
	if len(output) != 2 {
		t.Errorf("Expected 2 output items, got %d", len(output))
	}

	if !strings.Contains(output[0], "<item>") {
		t.Errorf("Expected XML formatted output, got %s", output[0])
	}
}

func TestHTMLReportBuilder(t *testing.T) {
	builder := NewHTMLReportBuilder("Test Report")
	data := []string{"Item 1", "Item 2"}

	report := BuildReport(builder, data)

	if !strings.Contains(report, "<html>") {
		t.Error("Report should contain HTML tags")
	}

	if !strings.Contains(report, "Test Report") {
		t.Error("Report should contain title")
	}

	if !strings.Contains(report, "Item 1") {
		t.Error("Report should contain data items")
	}
}

func TestMarkdownReportBuilder(t *testing.T) {
	builder := NewMarkdownReportBuilder("Test Report")
	data := []string{"Item 1", "Item 2"}

	report := BuildReport(builder, data)

	if !strings.Contains(report, "# Test Report") {
		t.Error("Report should contain Markdown header")
	}

	if !strings.Contains(report, "- Item 1") {
		t.Error("Report should contain list items")
	}
}

func TestPlainTextReportBuilder(t *testing.T) {
	builder := NewPlainTextReportBuilder("Test Report")
	data := []string{"Item 1", "Item 2"}

	report := BuildReport(builder, data)

	if !strings.Contains(report, "Test Report") {
		t.Error("Report should contain title")
	}

	if !strings.Contains(report, "1. Item 1") {
		t.Error("Report should contain numbered items")
	}
}

func TestZombieAI(t *testing.T) {
	zombie := NewZombieAI()

	// Test without player nearby
	zombie.SetPlayerNear(false)
	err := AIUpdateCycle(zombie)
	if err != nil {
		t.Errorf("AIUpdateCycle failed: %v", err)
	}

	// Test with player nearby
	zombie.SetPlayerNear(true)
	err = AIUpdateCycle(zombie)
	if err != nil {
		t.Errorf("AIUpdateCycle failed: %v", err)
	}

	if zombie.position != "chasing" {
		t.Errorf("Expected zombie to be chasing, got %s", zombie.position)
	}
}

func TestGuardAI(t *testing.T) {
	guard := NewGuardAI()

	// Test normal patrol
	guard.SetAlertLevel(0)
	err := AIUpdateCycle(guard)
	if err != nil {
		t.Errorf("AIUpdateCycle failed: %v", err)
	}

	// Test with suspicious activity
	guard.SetAlertLevel(1)
	err = AIUpdateCycle(guard)
	if err != nil {
		t.Errorf("AIUpdateCycle failed: %v", err)
	}
}

// Benchmark tests
func BenchmarkProcessTemplate(b *testing.B) {
	input := []string{"test1", "test2", "test3", "test4", "test5"}
	processor := NewCSVProcessor(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTemplate(processor)
	}
}

func BenchmarkBuildReport(b *testing.B) {
	data := []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"}
	builder := NewHTMLReportBuilder("Test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildReport(builder, data)
	}
}
