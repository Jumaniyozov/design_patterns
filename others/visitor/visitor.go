// Package visitor demonstrates the Visitor pattern.
// It separates algorithms from object structures, allowing new operations
// without modifying the classes of elements being operated on.
package visitor

import (
	"fmt"
	"strings"
)

// ShapeVisitor defines operations on shapes
type ShapeVisitor interface {
	VisitCircle(c *Circle) string
	VisitRectangle(r *Rectangle) string
	VisitTriangle(t *Triangle) string
}

// Shape interface accepts visitors
type Shape interface {
	Accept(v ShapeVisitor) string
}

// Circle shape
type Circle struct {
	Radius float64
}

func (c *Circle) Accept(v ShapeVisitor) string {
	return v.VisitCircle(c)
}

// Rectangle shape
type Rectangle struct {
	Width, Height float64
}

func (r *Rectangle) Accept(v ShapeVisitor) string {
	return v.VisitRectangle(r)
}

// Triangle shape
type Triangle struct {
	Base, Height float64
}

func (t *Triangle) Accept(v ShapeVisitor) string {
	return v.VisitTriangle(t)
}

// AreaCalculator visitor calculates areas
type AreaCalculator struct{}

func (a *AreaCalculator) VisitCircle(c *Circle) string {
	area := 3.14159 * c.Radius * c.Radius
	return fmt.Sprintf("Circle area: %.2f", area)
}

func (a *AreaCalculator) VisitRectangle(r *Rectangle) string {
	area := r.Width * r.Height
	return fmt.Sprintf("Rectangle area: %.2f", area)
}

func (a *AreaCalculator) VisitTriangle(t *Triangle) string {
	area := 0.5 * t.Base * t.Height
	return fmt.Sprintf("Triangle area: %.2f", area)
}

// PerimeterCalculator visitor calculates perimeters
type PerimeterCalculator struct{}

func (p *PerimeterCalculator) VisitCircle(c *Circle) string {
	perimeter := 2 * 3.14159 * c.Radius
	return fmt.Sprintf("Circle perimeter: %.2f", perimeter)
}

func (p *PerimeterCalculator) VisitRectangle(r *Rectangle) string {
	perimeter := 2 * (r.Width + r.Height)
	return fmt.Sprintf("Rectangle perimeter: %.2f", perimeter)
}

func (p *PerimeterCalculator) VisitTriangle(t *Triangle) string {
	// Assuming right triangle for simplicity
	hypotenuse := 0.0
	if t.Base > 0 && t.Height > 0 {
		hypotenuse = ((t.Base * t.Base) + (t.Height * t.Height))
	}
	perimeter := t.Base + t.Height + hypotenuse
	return fmt.Sprintf("Triangle perimeter: %.2f", perimeter)
}

// Document element visitor example

// DocumentVisitor defines operations on document elements
type DocumentVisitor interface {
	VisitParagraph(p *Paragraph) string
	VisitImage(i *Image) string
	VisitTable(t *Table) string
}

// DocumentElement accepts visitors
type DocumentElement interface {
	Accept(v DocumentVisitor) string
}

// Paragraph element
type Paragraph struct {
	Text string
}

func (p *Paragraph) Accept(v DocumentVisitor) string {
	return v.VisitParagraph(p)
}

// Image element
type Image struct {
	URL    string
	Alt    string
	Width  int
	Height int
}

func (i *Image) Accept(v DocumentVisitor) string {
	return v.VisitImage(i)
}

// Table element
type Table struct {
	Rows    int
	Columns int
	Data    [][]string
}

func (t *Table) Accept(v DocumentVisitor) string {
	return v.VisitTable(t)
}

// HTMLExporter exports to HTML
type HTMLExporter struct{}

func (h *HTMLExporter) VisitParagraph(p *Paragraph) string {
	return fmt.Sprintf("<p>%s</p>", p.Text)
}

func (h *HTMLExporter) VisitImage(i *Image) string {
	return fmt.Sprintf(`<img src="%s" alt="%s" width="%d" height="%d"/>`,
		i.URL, i.Alt, i.Width, i.Height)
}

func (h *HTMLExporter) VisitTable(t *Table) string {
	html := "<table>\n"
	for _, row := range t.Data {
		html += "  <tr>\n"
		for _, cell := range row {
			html += fmt.Sprintf("    <td>%s</td>\n", cell)
		}
		html += "  </tr>\n"
	}
	html += "</table>"
	return html
}

// MarkdownExporter exports to Markdown
type MarkdownExporter struct{}

func (m *MarkdownExporter) VisitParagraph(p *Paragraph) string {
	return p.Text + "\n"
}

func (m *MarkdownExporter) VisitImage(i *Image) string {
	return fmt.Sprintf("![%s](%s)", i.Alt, i.URL)
}

func (m *MarkdownExporter) VisitTable(t *Table) string {
	if len(t.Data) == 0 {
		return ""
	}

	md := ""
	// Header row
	md += "| " + strings.Join(t.Data[0], " | ") + " |\n"
	// Separator
	md += "|" + strings.Repeat(" --- |", t.Columns) + "\n"
	// Data rows
	for i := 1; i < len(t.Data); i++ {
		md += "| " + strings.Join(t.Data[i], " | ") + " |\n"
	}
	return md
}

// WordCounter counts words in document
type WordCounter struct {
	count int
}

func (w *WordCounter) VisitParagraph(p *Paragraph) string {
	words := strings.Fields(p.Text)
	w.count += len(words)
	return fmt.Sprintf("Paragraph: %d words", len(words))
}

func (w *WordCounter) VisitImage(i *Image) string {
	return "Image: 0 words"
}

func (w *WordCounter) VisitTable(t *Table) string {
	count := 0
	for _, row := range t.Data {
		for _, cell := range row {
			count += len(strings.Fields(cell))
		}
	}
	w.count += count
	return fmt.Sprintf("Table: %d words", count)
}

func (w *WordCounter) GetTotalCount() int {
	return w.count
}
