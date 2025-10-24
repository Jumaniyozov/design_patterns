// Package templatemethod implements the Template Method pattern.
//
// The Template Method pattern defines the skeleton of an algorithm in a method,
// deferring some steps to subclasses or implementations. It lets implementations
// redefine certain steps of an algorithm without changing the algorithm's structure.
//
// In Go, we use composition and interfaces instead of inheritance, making the
// pattern more flexible and idiomatic.
//
// Key components:
// - Algorithm interface: Defines customizable steps
// - Template executor: Defines the algorithm structure
// - Concrete implementations: Provide specific step implementations
package templatemethod

import (
	"fmt"
	"strings"
)

// DataProcessor defines the interface for data processing steps.
type DataProcessor interface {
	ReadData() ([]string, error)
	ProcessData(data []string) ([]string, error)
	WriteData(data []string) error
}

// ProcessTemplate defines the template method that orchestrates the algorithm.
func ProcessTemplate(processor DataProcessor) error {
	fmt.Println("[Template] Starting data processing pipeline...")

	// Step 1: Read data
	data, err := processor.ReadData()
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}
	fmt.Printf("[Template] Read %d records\n", len(data))

	// Step 2: Process data
	processed, err := processor.ProcessData(data)
	if err != nil {
		return fmt.Errorf("process failed: %w", err)
	}
	fmt.Printf("[Template] Processed %d records\n", len(processed))

	// Step 3: Write data
	err = processor.WriteData(processed)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	fmt.Println("[Template] Pipeline completed successfully")
	return nil
}

// CSVProcessor implements data processing for CSV files.
type CSVProcessor struct {
	inputData  []string
	outputData []string
}

// NewCSVProcessor creates a new CSV processor.
func NewCSVProcessor(input []string) *CSVProcessor {
	return &CSVProcessor{
		inputData: input,
	}
}

// ReadData reads CSV data.
func (p *CSVProcessor) ReadData() ([]string, error) {
	fmt.Println("[CSV] Reading CSV data...")
	return p.inputData, nil
}

// ProcessData processes CSV data (uppercase conversion).
func (p *CSVProcessor) ProcessData(data []string) ([]string, error) {
	fmt.Println("[CSV] Converting to uppercase...")
	result := make([]string, len(data))
	for i, line := range data {
		result[i] = strings.ToUpper(line)
	}
	return result, nil
}

// WriteData writes CSV data.
func (p *CSVProcessor) WriteData(data []string) error {
	fmt.Println("[CSV] Writing CSV output...")
	p.outputData = data
	return nil
}

// GetOutput returns the processed output.
func (p *CSVProcessor) GetOutput() []string {
	return p.outputData
}

// JSONProcessor implements data processing for JSON files.
type JSONProcessor struct {
	inputData  []string
	outputData []string
}

// NewJSONProcessor creates a new JSON processor.
func NewJSONProcessor(input []string) *JSONProcessor {
	return &JSONProcessor{
		inputData: input,
	}
}

// ReadData reads JSON data.
func (p *JSONProcessor) ReadData() ([]string, error) {
	fmt.Println("[JSON] Parsing JSON data...")
	return p.inputData, nil
}

// ProcessData processes JSON data (add JSON formatting).
func (p *JSONProcessor) ProcessData(data []string) ([]string, error) {
	fmt.Println("[JSON] Adding JSON formatting...")
	result := make([]string, len(data))
	for i, item := range data {
		result[i] = fmt.Sprintf(`{"value": "%s"}`, item)
	}
	return result, nil
}

// WriteData writes JSON data.
func (p *JSONProcessor) WriteData(data []string) error {
	fmt.Println("[JSON] Writing JSON output...")
	p.outputData = data
	return nil
}

// GetOutput returns the processed output.
func (p *JSONProcessor) GetOutput() []string {
	return p.outputData
}

// XMLProcessor implements data processing for XML files.
type XMLProcessor struct {
	inputData  []string
	outputData []string
}

// NewXMLProcessor creates a new XML processor.
func NewXMLProcessor(input []string) *XMLProcessor {
	return &XMLProcessor{
		inputData: input,
	}
}

// ReadData reads XML data.
func (p *XMLProcessor) ReadData() ([]string, error) {
	fmt.Println("[XML] Parsing XML data...")
	return p.inputData, nil
}

// ProcessData processes XML data (add XML tags).
func (p *XMLProcessor) ProcessData(data []string) ([]string, error) {
	fmt.Println("[XML] Adding XML formatting...")
	result := make([]string, len(data))
	for i, item := range data {
		result[i] = fmt.Sprintf("<item>%s</item>", item)
	}
	return result, nil
}

// WriteData writes XML data.
func (p *XMLProcessor) WriteData(data []string) error {
	fmt.Println("[XML] Writing XML output...")
	p.outputData = data
	return nil
}

// GetOutput returns the processed output.
func (p *XMLProcessor) GetOutput() []string {
	return p.outputData
}

// ReportBuilder defines the interface for building reports.
type ReportBuilder interface {
	GenerateHeader() string
	GenerateBody(data []string) string
	GenerateFooter() string
}

// BuildReport is the template method for report generation.
func BuildReport(builder ReportBuilder, data []string) string {
	fmt.Println("[Template] Building report...")

	report := ""

	// Step 1: Generate header
	report += builder.GenerateHeader()

	// Step 2: Generate body
	report += builder.GenerateBody(data)

	// Step 3: Generate footer
	report += builder.GenerateFooter()

	fmt.Println("[Template] Report built successfully")
	return report
}

// HTMLReportBuilder builds HTML reports.
type HTMLReportBuilder struct {
	title string
}

// NewHTMLReportBuilder creates a new HTML report builder.
func NewHTMLReportBuilder(title string) *HTMLReportBuilder {
	return &HTMLReportBuilder{title: title}
}

// GenerateHeader generates HTML header.
func (b *HTMLReportBuilder) GenerateHeader() string {
	return fmt.Sprintf("<html><head><title>%s</title></head><body>\n", b.title)
}

// GenerateBody generates HTML body.
func (b *HTMLReportBuilder) GenerateBody(data []string) string {
	body := "<ul>\n"
	for _, item := range data {
		body += fmt.Sprintf("  <li>%s</li>\n", item)
	}
	body += "</ul>\n"
	return body
}

// GenerateFooter generates HTML footer.
func (b *HTMLReportBuilder) GenerateFooter() string {
	return "</body></html>\n"
}

// MarkdownReportBuilder builds Markdown reports.
type MarkdownReportBuilder struct {
	title string
}

// NewMarkdownReportBuilder creates a new Markdown report builder.
func NewMarkdownReportBuilder(title string) *MarkdownReportBuilder {
	return &MarkdownReportBuilder{title: title}
}

// GenerateHeader generates Markdown header.
func (b *MarkdownReportBuilder) GenerateHeader() string {
	return fmt.Sprintf("# %s\n\n", b.title)
}

// GenerateBody generates Markdown body.
func (b *MarkdownReportBuilder) GenerateBody(data []string) string {
	body := ""
	for _, item := range data {
		body += fmt.Sprintf("- %s\n", item)
	}
	body += "\n"
	return body
}

// GenerateFooter generates Markdown footer.
func (b *MarkdownReportBuilder) GenerateFooter() string {
	return "---\n*Generated report*\n"
}

// PlainTextReportBuilder builds plain text reports.
type PlainTextReportBuilder struct {
	title string
}

// NewPlainTextReportBuilder creates a new plain text report builder.
func NewPlainTextReportBuilder(title string) *PlainTextReportBuilder {
	return &PlainTextReportBuilder{title: title}
}

// GenerateHeader generates plain text header.
func (b *PlainTextReportBuilder) GenerateHeader() string {
	return fmt.Sprintf("%s\n%s\n\n", b.title, strings.Repeat("=", len(b.title)))
}

// GenerateBody generates plain text body.
func (b *PlainTextReportBuilder) GenerateBody(data []string) string {
	body := ""
	for i, item := range data {
		body += fmt.Sprintf("%d. %s\n", i+1, item)
	}
	body += "\n"
	return body
}

// GenerateFooter generates plain text footer.
func (b *PlainTextReportBuilder) GenerateFooter() string {
	return "--- End of Report ---\n"
}

// GameAI defines the interface for game AI behavior.
type GameAI interface {
	Sense() []string
	Think(observations []string) string
	Act(decision string) error
}

// AIUpdateCycle is the template method for AI updates.
func AIUpdateCycle(ai GameAI) error {
	fmt.Println("[AI] Starting AI update cycle...")

	// Step 1: Sense environment
	observations := ai.Sense()

	// Step 2: Think/decide based on observations
	decision := ai.Think(observations)

	// Step 3: Act on decision
	err := ai.Act(decision)
	if err != nil {
		return err
	}

	fmt.Println("[AI] AI cycle completed")
	return nil
}

// ZombieAI implements AI for zombie characters.
type ZombieAI struct {
	position   string
	playerNear bool
}

// NewZombieAI creates a new zombie AI.
func NewZombieAI() *ZombieAI {
	return &ZombieAI{position: "idle"}
}

// Sense senses the environment.
func (z *ZombieAI) Sense() []string {
	fmt.Println("[Zombie] Sensing environment...")
	observations := []string{}

	if z.playerNear {
		observations = append(observations, "player-nearby")
	}

	return observations
}

// Think makes decisions based on observations.
func (z *ZombieAI) Think(observations []string) string {
	fmt.Println("[Zombie] Processing observations...")

	for _, obs := range observations {
		if obs == "player-nearby" {
			return "chase-player"
		}
	}

	return "wander"
}

// Act executes the decision.
func (z *ZombieAI) Act(decision string) error {
	fmt.Printf("[Zombie] Executing action: %s\n", decision)

	switch decision {
	case "chase-player":
		z.position = "chasing"
	case "wander":
		z.position = "wandering"
	}

	return nil
}

// SetPlayerNear sets whether player is nearby.
func (z *ZombieAI) SetPlayerNear(near bool) {
	z.playerNear = near
}

// GuardAI implements AI for guard characters.
type GuardAI struct {
	position      string
	patrolRoute   []string
	currentIndex  int
	alertLevel    int
}

// NewGuardAI creates a new guard AI.
func NewGuardAI() *GuardAI {
	return &GuardAI{
		position:    "patrol",
		patrolRoute: []string{"point-A", "point-B", "point-C"},
		alertLevel:  0,
	}
}

// Sense senses the environment.
func (g *GuardAI) Sense() []string {
	fmt.Println("[Guard] Scanning area...")
	observations := []string{}

	if g.alertLevel > 0 {
		observations = append(observations, "suspicious-activity")
	}

	return observations
}

// Think makes decisions based on observations.
func (g *GuardAI) Think(observations []string) string {
	fmt.Println("[Guard] Analyzing situation...")

	for _, obs := range observations {
		if obs == "suspicious-activity" {
			return "investigate"
		}
	}

	return "patrol"
}

// Act executes the decision.
func (g *GuardAI) Act(decision string) error {
	fmt.Printf("[Guard] Executing action: %s\n", decision)

	switch decision {
	case "patrol":
		currentPoint := g.patrolRoute[g.currentIndex]
		fmt.Printf("[Guard] Patrolling to %s\n", currentPoint)
		g.currentIndex = (g.currentIndex + 1) % len(g.patrolRoute)
	case "investigate":
		fmt.Println("[Guard] Investigating suspicious activity")
	}

	return nil
}

// SetAlertLevel sets the alert level.
func (g *GuardAI) SetAlertLevel(level int) {
	g.alertLevel = level
}
