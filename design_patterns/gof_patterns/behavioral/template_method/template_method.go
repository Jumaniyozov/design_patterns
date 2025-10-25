package template_method

import "fmt"

type DataProcessor interface {
	ReadData() string
	ProcessData(data string) string
	WriteData(data string)
}

type AbstractDataProcessor struct {
	processor DataProcessor
}

func (a *AbstractDataProcessor) Process() {
	data := a.processor.ReadData()
	processed := a.processor.ProcessData(data)
	a.processor.WriteData(processed)
}

type CSVProcessor struct{}

func (c *CSVProcessor) ReadData() string {
	fmt.Println("Reading CSV file")
	return "csv data"
}

func (c *CSVProcessor) ProcessData(data string) string {
	fmt.Println("Processing CSV data")
	return "processed " + data
}

func (c *CSVProcessor) WriteData(data string) {
	fmt.Println("Writing CSV:", data)
}

type XMLProcessor struct{}

func (x *XMLProcessor) ReadData() string {
	fmt.Println("Reading XML file")
	return "xml data"
}

func (x *XMLProcessor) ProcessData(data string) string {
	fmt.Println("Processing XML data")
	return "processed " + data
}

func (x *XMLProcessor) WriteData(data string) {
	fmt.Println("Writing XML:", data)
}
