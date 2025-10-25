// Package activeobject demonstrates the Active Object pattern.
// It decouples method execution from invocation, running methods
// asynchronously in their own goroutine with future-based results.
package activeobject

import (
	"fmt"
	"time"
)

// Future represents a future result
type Future[T any] struct {
	result chan T
}

// NewFuture creates a future
func NewFuture[T any]() *Future[T] {
	return &Future[T]{
		result: make(chan T, 1),
	}
}

// Get gets the result (blocks until available)
func (f *Future[T]) Get() T {
	return <-f.result
}

// TryGet tries to get result without blocking
func (f *Future[T]) TryGet() (T, bool) {
	select {
	case result := <-f.result:
		return result, true
	default:
		var zero T
		return zero, false
	}
}

// Set sets the result
func (f *Future[T]) Set(value T) {
	f.result <- value
}

// ActiveObject executes methods asynchronously
type ActiveObject struct {
	requests chan func()
	stop     chan struct{}
}

// NewActiveObject creates an active object
func NewActiveObject() *ActiveObject {
	ao := &ActiveObject{
		requests: make(chan func(), 100),
		stop:     make(chan struct{}),
	}
	go ao.run()
	return ao
}

func (ao *ActiveObject) run() {
	for {
		select {
		case req := <-ao.requests:
			req()
		case <-ao.stop:
			return
		}
	}
}

// Stop stops the active object
func (ao *ActiveObject) Stop() {
	close(ao.stop)
}

// Submit submits a task and returns a future
func (ao *ActiveObject) Submit(task func() interface{}) *Future[interface{}] {
	future := NewFuture[interface{}]()
	ao.requests <- func() {
		result := task()
		future.Set(result)
	}
	return future
}

// DataProcessor demonstrates active object pattern
type DataProcessor struct {
	activeObject *ActiveObject
}

// NewDataProcessor creates a data processor
func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		activeObject: NewActiveObject(),
	}
}

// Process processes data asynchronously
func (dp *DataProcessor) Process(data string) *Future[string] {
	future := NewFuture[string]()
	dp.activeObject.requests <- func() {
		// Simulate processing
		time.Sleep(100 * time.Millisecond)
		result := fmt.Sprintf("Processed: %s", data)
		future.Set(result)
	}
	return future
}

// Batch processes data asynchronously
func (dp *DataProcessor) Batch(items []string) *Future[[]string] {
	future := NewFuture[[]string]()
	dp.activeObject.requests <- func() {
		results := make([]string, len(items))
		for i, item := range items {
			time.Sleep(50 * time.Millisecond)
			results[i] = fmt.Sprintf("Batch processed: %s", item)
		}
		future.Set(results)
	}
	return future
}

// Stop stops the processor
func (dp *DataProcessor) Stop() {
	dp.activeObject.Stop()
}

// AsyncCalculator demonstrates async operations
type AsyncCalculator struct {
	activeObject *ActiveObject
}

// NewAsyncCalculator creates an async calculator
func NewAsyncCalculator() *AsyncCalculator {
	return &AsyncCalculator{
		activeObject: NewActiveObject(),
	}
}

// Add adds two numbers asynchronously
func (ac *AsyncCalculator) Add(a, b int) *Future[int] {
	future := NewFuture[int]()
	ac.activeObject.requests <- func() {
		time.Sleep(50 * time.Millisecond) // Simulate work
		future.Set(a + b)
	}
	return future
}

// Multiply multiplies two numbers asynchronously
func (ac *AsyncCalculator) Multiply(a, b int) *Future[int] {
	future := NewFuture[int]()
	ac.activeObject.requests <- func() {
		time.Sleep(50 * time.Millisecond)
		future.Set(a * b)
	}
	return future
}

// Stop stops the calculator
func (ac *AsyncCalculator) Stop() {
	ac.activeObject.Stop()
}

// FileWriter demonstrates async I/O
type FileWriter struct {
	activeObject *ActiveObject
}

// NewFileWriter creates a file writer
func NewFileWriter() *FileWriter {
	return &FileWriter{
		activeObject: NewActiveObject(),
	}
}

// Write writes data asynchronously
func (fw *FileWriter) Write(filename, content string) *Future[bool] {
	future := NewFuture[bool]()
	fw.activeObject.requests <- func() {
		// Simulate file write
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Writing to %s: %s\n", filename, content)
		future.Set(true)
	}
	return future
}

// Append appends data asynchronously
func (fw *FileWriter) Append(filename, content string) *Future[bool] {
	future := NewFuture[bool]()
	fw.activeObject.requests <- func() {
		time.Sleep(150 * time.Millisecond)
		fmt.Printf("Appending to %s: %s\n", filename, content)
		future.Set(true)
	}
	return future
}

// Stop stops the writer
func (fw *FileWriter) Stop() {
	fw.activeObject.Stop()
}
