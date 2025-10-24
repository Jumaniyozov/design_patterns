// Package pipeline implements the Pipeline concurrency pattern.
//
// The Pipeline pattern processes data through a series of stages connected by channels.
// Each stage is a goroutine or group of goroutines that:
//  1. Receives values from upstream via inbound channels
//  2. Performs some function on that data
//  3. Sends values downstream via outbound channels
//
// Key characteristics:
// - Each stage runs concurrently in its own goroutine
// - Stages are connected by channels that provide natural backpressure
// - The stage that creates a channel is responsible for closing it
// - Stages use range loops to receive values until the channel is closed
package pipeline

import (
	"context"
	"fmt"
	"sync"
)

// Stage represents a processing stage in a pipeline.
// It receives input from an inbound channel, processes it, and sends output to an outbound channel.
type Stage[In, Out any] func(ctx context.Context, in <-chan In) <-chan Out

// Generator creates values and sends them to a channel.
// This is typically the first stage of a pipeline.
func Generator[T any](ctx context.Context, values ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range values {
			select {
			case out <- v:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Map transforms each input value using the provided function.
// This is a basic transformation stage.
func Map[In, Out any](ctx context.Context, in <-chan In, fn func(In) Out) <-chan Out {
	out := make(chan Out)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- fn(v):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Filter passes through only values that satisfy the predicate.
func Filter[T any](ctx context.Context, in <-chan T, predicate func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range in {
			if predicate(v) {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// FanOut creates multiple workers that all read from the same input channel.
// This is useful for parallelizing CPU-intensive operations.
func FanOut[T any](ctx context.Context, in <-chan T, workerCount int, worker func(T) T) []<-chan T {
	channels := make([]<-chan T, workerCount)

	for i := 0; i < workerCount; i++ {
		out := make(chan T)
		channels[i] = out

		go func(out chan<- T) {
			defer close(out)
			for v := range in {
				select {
				case out <- worker(v):
				case <-ctx.Done():
					return
				}
			}
		}(out)
	}

	return channels
}

// FanIn merges multiple channels into a single channel.
// This is the counterpart to FanOut, combining results from multiple workers.
func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	// Close output channel when all input channels are drained
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Tee splits input into two output channels.
// Both outputs receive all input values.
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)

		for v := range in {
			// Send to both channels - we need to ensure both get the value
			// Create two copies of the value
			val1, val2 := v, v

			// Send to first channel
			select {
			case out1 <- val1:
			case <-ctx.Done():
				return
			}

			// Send to second channel
			select {
			case out2 <- val2:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out1, out2
}

// Batch groups input values into batches of the specified size.
func Batch[T any](ctx context.Context, in <-chan T, size int) <-chan []T {
	out := make(chan []T)

	go func() {
		defer close(out)
		batch := make([]T, 0, size)

		for v := range in {
			batch = append(batch, v)

			if len(batch) >= size {
				select {
				case out <- batch:
					batch = make([]T, 0, size)
				case <-ctx.Done():
					return
				}
			}
		}

		// Send remaining batch if not empty
		if len(batch) > 0 {
			select {
			case out <- batch:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Result wraps a value with a potential error for error handling in pipelines.
type Result[T any] struct {
	Value T
	Error error
}

// MapWithError transforms input values and handles errors.
func MapWithError[In, Out any](ctx context.Context, in <-chan In, fn func(In) (Out, error)) <-chan Result[Out] {
	out := make(chan Result[Out])

	go func() {
		defer close(out)
		for v := range in {
			result, err := fn(v)
			select {
			case out <- Result[Out]{Value: result, Error: err}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// CollectErrors separates successful values from errors.
func CollectErrors[T any](ctx context.Context, in <-chan Result[T]) (<-chan T, <-chan error) {
	values := make(chan T)
	errors := make(chan error)

	go func() {
		defer close(values)
		defer close(errors)

		for result := range in {
			if result.Error != nil {
				select {
				case errors <- result.Error:
				case <-ctx.Done():
					return
				}
			} else {
				select {
				case values <- result.Value:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return values, errors
}

// Buffer creates a buffered channel stage to smooth out bursty traffic.
func Buffer[T any](ctx context.Context, in <-chan T, size int) <-chan T {
	out := make(chan T, size)

	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- v:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Take limits the pipeline to the first n values.
func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		count := 0
		for v := range in {
			if count >= n {
				return
			}
			select {
			case out <- v:
				count++
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Sink consumes all values from a channel and calls the provided function.
// This is typically the final stage of a pipeline.
func Sink[T any](ctx context.Context, in <-chan T, fn func(T)) {
	for v := range in {
		select {
		case <-ctx.Done():
			return
		default:
			fn(v)
		}
	}
}

// Pipeline is a builder for constructing pipelines fluently.
type Pipeline[T any] struct {
	ctx context.Context
	out <-chan T
}

// NewPipeline creates a new pipeline with the given context and source channel.
func NewPipeline[T any](ctx context.Context, source <-chan T) *Pipeline[T] {
	return &Pipeline[T]{
		ctx: ctx,
		out: source,
	}
}

// Map applies a transformation function to each value in the pipeline.
func (p *Pipeline[T]) Map(fn func(T) T) *Pipeline[T] {
	return &Pipeline[T]{
		ctx: p.ctx,
		out: Map(p.ctx, p.out, fn),
	}
}

// Filter applies a predicate to filter values in the pipeline.
func (p *Pipeline[T]) Filter(predicate func(T) bool) *Pipeline[T] {
	return &Pipeline[T]{
		ctx: p.ctx,
		out: Filter(p.ctx, p.out, predicate),
	}
}

// Take limits the pipeline to the first n values.
func (p *Pipeline[T]) Take(n int) *Pipeline[T] {
	return &Pipeline[T]{
		ctx: p.ctx,
		out: Take(p.ctx, p.out, n),
	}
}

// Buffer adds a buffer to smooth out traffic.
func (p *Pipeline[T]) Buffer(size int) *Pipeline[T] {
	return &Pipeline[T]{
		ctx: p.ctx,
		out: Buffer(p.ctx, p.out, size),
	}
}

// Out returns the output channel of the pipeline.
func (p *Pipeline[T]) Out() <-chan T {
	return p.out
}

// Collect gathers all values from the pipeline into a slice.
func (p *Pipeline[T]) Collect() []T {
	var results []T
	for v := range p.out {
		select {
		case <-p.ctx.Done():
			return results
		default:
			results = append(results, v)
		}
	}
	return results
}

// ForEach applies a function to each value in the pipeline.
func (p *Pipeline[T]) ForEach(fn func(T)) {
	Sink(p.ctx, p.out, fn)
}

// String formats a Result for display
func (r Result[T]) String() string {
	if r.Error != nil {
		return fmt.Sprintf("Error: %v", r.Error)
	}
	return fmt.Sprintf("Value: %v", r.Value)
}
