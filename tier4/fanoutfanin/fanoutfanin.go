// Package fanoutfanin implements the Fan-Out/Fan-In concurrency pattern.
//
// The Fan-Out/Fan-In pattern distributes work across multiple goroutines (fan-out)
// and then combines their results into a single channel (fan-in). This pattern
// is essential for parallel processing and maximizing throughput.
//
// Key characteristics:
// - Fan-Out: Distribute work from one channel to multiple workers
// - Fan-In: Merge results from multiple channels into one
// - Parallel Execution: Workers process independently
// - Result Aggregation: Combine outputs while preserving data
//
// Use cases:
// - Parallel data processing
// - Load distribution across workers
// - Aggregating results from multiple sources
// - Maximizing CPU utilization
package fanoutfanin

import (
	"context"
	"sync"
)

// WorkerFunc defines a function that processes input and produces output.
type WorkerFunc[In, Out any] func(In) Out

// FanOut distributes work from input channel to multiple worker goroutines.
// Each worker reads from the shared input channel and writes to its own output channel.
func FanOut[In, Out any](ctx context.Context, input <-chan In, numWorkers int, worker WorkerFunc[In, Out]) []<-chan Out {
	outputs := make([]<-chan Out, numWorkers)

	for i := 0; i < numWorkers; i++ {
		output := make(chan Out)
		outputs[i] = output

		go func(out chan<- Out) {
			defer close(out)
			for item := range input {
				select {
				case out <- worker(item):
				case <-ctx.Done():
					return
				}
			}
		}(output)
	}

	return outputs
}

// FanIn merges multiple input channels into a single output channel.
// Reads from all input channels concurrently and forwards to output.
func FanIn[T any](ctx context.Context, inputs ...<-chan T) <-chan T {
	output := make(chan T)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, input := range inputs {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for item := range ch {
				select {
				case output <- item:
				case <-ctx.Done():
					return
				}
			}
		}(input)
	}

	// Close output channel when all inputs are drained
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// FanOutFanIn combines fan-out and fan-in into a single operation.
// This is a convenience function for the common pattern of distributing work
// and then collecting results.
func FanOutFanIn[In, Out any](ctx context.Context, input <-chan In, numWorkers int, worker WorkerFunc[In, Out]) <-chan Out {
	workerOutputs := FanOut(ctx, input, numWorkers, worker)
	return FanIn(ctx, workerOutputs...)
}

// OrderedFanOutFanIn maintains input order in the output.
// This is more expensive but guarantees that output order matches input order.
func OrderedFanOutFanIn[In, Out any](ctx context.Context, input <-chan In, numWorkers int, worker WorkerFunc[In, Out]) <-chan Out {
	type indexedItem struct {
		index int
		value Out
	}

	output := make(chan Out)

	go func() {
		defer close(output)

		// Create indexed input channel
		indexedInput := make(chan struct {
			index int
			value In
		})

		// Index input items
		go func() {
			defer close(indexedInput)
			index := 0
			for item := range input {
				select {
				case indexedInput <- struct {
					index int
					value In
				}{index: index, value: item}:
					index++
				case <-ctx.Done():
					return
				}
			}
		}()

		// Fan-out to workers with indexed items
		indexedOutputs := make([]<-chan indexedItem, numWorkers)
		for i := 0; i < numWorkers; i++ {
			indexedOutput := make(chan indexedItem)
			indexedOutputs[i] = indexedOutput

			go func(out chan<- indexedItem) {
				defer close(out)
				for item := range indexedInput {
					result := worker(item.value)
					select {
					case out <- indexedItem{index: item.index, value: result}:
					case <-ctx.Done():
						return
					}
				}
			}(indexedOutput)
		}

		// Fan-in with ordering
		merged := FanIn(ctx, indexedOutputs...)

		// Collect and reorder results
		results := make(map[int]Out)
		nextIndex := 0

		for item := range merged {
			results[item.index] = item.value

			// Output results in order
			for {
				if val, ok := results[nextIndex]; ok {
					select {
					case output <- val:
						delete(results, nextIndex)
						nextIndex++
					case <-ctx.Done():
						return
					}
				} else {
					break
				}
			}
		}
	}()

	return output
}

// WorkPool represents a pool of workers that can process jobs.
type WorkPool[In, Out any] struct {
	ctx        context.Context
	numWorkers int
	worker     WorkerFunc[In, Out]
	jobs       chan In
	results    chan Out
	wg         sync.WaitGroup
}

// NewWorkPool creates a new work pool with the specified number of workers.
func NewWorkPool[In, Out any](ctx context.Context, numWorkers int, worker WorkerFunc[In, Out]) *WorkPool[In, Out] {
	return &WorkPool[In, Out]{
		ctx:        ctx,
		numWorkers: numWorkers,
		worker:     worker,
		jobs:       make(chan In),
		results:    make(chan Out),
	}
}

// Start begins processing jobs with the worker pool.
func (wp *WorkPool[In, Out]) Start() {
	// Start workers
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobs {
				result := wp.worker(job)
				select {
				case wp.results <- result:
				case <-wp.ctx.Done():
					return
				}
			}
		}()
	}

	// Close results channel when all workers complete
	go func() {
		wp.wg.Wait()
		close(wp.results)
	}()
}

// Submit adds a job to the work pool.
func (wp *WorkPool[In, Out]) Submit(job In) bool {
	select {
	case wp.jobs <- job:
		return true
	case <-wp.ctx.Done():
		return false
	}
}

// Close closes the jobs channel and waits for workers to complete.
func (wp *WorkPool[In, Out]) Close() {
	close(wp.jobs)
}

// Results returns the results channel for reading worker outputs.
func (wp *WorkPool[In, Out]) Results() <-chan Out {
	return wp.results
}

// Generator creates a channel and populates it with values from a slice.
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

// Collect drains a channel and returns all values as a slice.
func Collect[T any](ctx context.Context, input <-chan T) []T {
	var results []T
	for item := range input {
		select {
		case <-ctx.Done():
			return results
		default:
			results = append(results, item)
		}
	}
	return results
}

// ParallelMap applies a function to each element in parallel using fan-out/fan-in.
func ParallelMap[In, Out any](ctx context.Context, input []In, numWorkers int, mapper WorkerFunc[In, Out]) []Out {
	inputChan := Generator(ctx, input...)
	outputChan := FanOutFanIn(ctx, inputChan, numWorkers, mapper)
	return Collect(ctx, outputChan)
}

// OrderedParallelMap applies a function in parallel while preserving input order.
func OrderedParallelMap[In, Out any](ctx context.Context, input []In, numWorkers int, mapper WorkerFunc[In, Out]) []Out {
	inputChan := Generator(ctx, input...)
	outputChan := OrderedFanOutFanIn(ctx, inputChan, numWorkers, mapper)
	return Collect(ctx, outputChan)
}

// BatchFanOut distributes batches of work to workers.
type BatchFanOut[In, Out any] struct {
	ctx        context.Context
	numWorkers int
	batchSize  int
	worker     func([]In) []Out
}

// NewBatchFanOut creates a new batch fan-out processor.
func NewBatchFanOut[In, Out any](ctx context.Context, numWorkers int, batchSize int, worker func([]In) []Out) *BatchFanOut[In, Out] {
	return &BatchFanOut[In, Out]{
		ctx:        ctx,
		numWorkers: numWorkers,
		batchSize:  batchSize,
		worker:     worker,
	}
}

// Process processes input in batches across multiple workers.
func (bf *BatchFanOut[In, Out]) Process(input <-chan In) <-chan Out {
	// Create batches
	batches := make(chan []In)
	go func() {
		defer close(batches)
		batch := make([]In, 0, bf.batchSize)

		for item := range input {
			batch = append(batch, item)
			if len(batch) >= bf.batchSize {
				select {
				case batches <- batch:
					batch = make([]In, 0, bf.batchSize)
				case <-bf.ctx.Done():
					return
				}
			}
		}

		// Send remaining batch
		if len(batch) > 0 {
			select {
			case batches <- batch:
			case <-bf.ctx.Done():
				return
			}
		}
	}()

	// Fan-out to workers
	workerOutputs := make([]<-chan Out, bf.numWorkers)
	for i := 0; i < bf.numWorkers; i++ {
		output := make(chan Out)
		workerOutputs[i] = output

		go func(out chan<- Out) {
			defer close(out)
			for batch := range batches {
				results := bf.worker(batch)
				for _, result := range results {
					select {
					case out <- result:
					case <-bf.ctx.Done():
						return
					}
				}
			}
		}(output)
	}

	// Fan-in results
	return FanIn(bf.ctx, workerOutputs...)
}
