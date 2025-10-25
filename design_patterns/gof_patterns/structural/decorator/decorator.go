// Package decorator demonstrates the Decorator pattern in Go.
package decorator

import "fmt"

// Coffee interface
type Coffee interface {
	Cost() float64
	Description() string
}

// SimpleCoffee is the base component
type SimpleCoffee struct{}

func (c *SimpleCoffee) Cost() float64 {
	return 2.0
}

func (c *SimpleCoffee) Description() string {
	return "Simple Coffee"
}

// MilkDecorator adds milk
type MilkDecorator struct {
	coffee Coffee
}

func NewMilkDecorator(c Coffee) Coffee {
	return &MilkDecorator{coffee: c}
}

func (m *MilkDecorator) Cost() float64 {
	return m.coffee.Cost() + 0.5
}

func (m *MilkDecorator) Description() string {
	return m.coffee.Description() + ", Milk"
}

// SugarDecorator adds sugar
type SugarDecorator struct {
	coffee Coffee
}

func NewSugarDecorator(c Coffee) Coffee {
	return &SugarDecorator{coffee: c}
}

func (s *SugarDecorator) Cost() float64 {
	return s.coffee.Cost() + 0.25
}

func (s *SugarDecorator) Description() string {
	return s.coffee.Description() + ", Sugar"
}

// WhippedCreamDecorator adds whipped cream
type WhippedCreamDecorator struct {
	coffee Coffee
}

func NewWhippedCreamDecorator(c Coffee) Coffee {
	return &WhippedCreamDecorator{coffee: c}
}

func (w *WhippedCreamDecorator) Cost() float64 {
	return w.coffee.Cost() + 0.75
}

func (w *WhippedCreamDecorator) Description() string {
	return w.coffee.Description() + ", Whipped Cream"
}

// DataSource interface for IO example
type DataSource interface {
	WriteData(data string) error
	ReadData() (string, error)
}

// FileDataSource is the base component
type FileDataSource struct {
	filename string
	data     string
}

func NewFileDataSource(filename string) DataSource {
	return &FileDataSource{filename: filename}
}

func (f *FileDataSource) WriteData(data string) error {
	f.data = data
	fmt.Printf("[File] Writing to %s: %s\n", f.filename, data)
	return nil
}

func (f *FileDataSource) ReadData() (string, error) {
	fmt.Printf("[File] Reading from %s\n", f.filename)
	return f.data, nil
}

// EncryptionDecorator adds encryption
type EncryptionDecorator struct {
	wrapped DataSource
}

func NewEncryptionDecorator(ds DataSource) DataSource {
	return &EncryptionDecorator{wrapped: ds}
}

func (e *EncryptionDecorator) WriteData(data string) error {
	encrypted := "[ENCRYPTED]" + data
	fmt.Println("[Encryption] Encrypting data")
	return e.wrapped.WriteData(encrypted)
}

func (e *EncryptionDecorator) ReadData() (string, error) {
	data, err := e.wrapped.ReadData()
	if err != nil {
		return "", err
	}
	fmt.Println("[Encryption] Decrypting data")
	return data[len("[ENCRYPTED]"):], nil
}

// CompressionDecorator adds compression
type CompressionDecorator struct {
	wrapped DataSource
}

func NewCompressionDecorator(ds DataSource) DataSource {
	return &CompressionDecorator{wrapped: ds}
}

func (c *CompressionDecorator) WriteData(data string) error {
	compressed := "[COMPRESSED]" + data
	fmt.Println("[Compression] Compressing data")
	return c.wrapped.WriteData(compressed)
}

func (c *CompressionDecorator) ReadData() (string, error) {
	data, err := c.wrapped.ReadData()
	if err != nil {
		return "", err
	}
	fmt.Println("[Compression] Decompressing data")
	return data[len("[COMPRESSED]"):], nil
}
