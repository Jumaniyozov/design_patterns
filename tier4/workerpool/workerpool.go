// Package workerpool implements the Worker Pool concurrency pattern.
//
// The Worker Pool pattern manages a fixed number of goroutines (workers) that
// process jobs from a shared queue. This pattern is essential for controlling
// resource usage and preventing goroutine explosion.
//
// Key characteristics:
// - Fixed number of workers
// - Shared job queue
// - Controlled concurrency
// - Resource management
// - Graceful shutdown
package workerpool

import (
	"context"
	"sync"
)

// Job represents a unit of work to be processed.
type Job[In, Out any] struct {
	ID    string
	Input In
}

// Result represents the output of processing a job.
type Result[Out any] struct {
	JobID string
	Value Out
	Error error
}

// Worker processes jobs from the job queue.
type Worker[In, Out any] struct {
	id         int
	jobQueue   <-chan Job[In, Out]
	resultChan chan<- Result[Out]
	processor  func(In) (Out, error)
	ctx        context.Context
	wg         *sync.WaitGroup
}

// Start begins processing jobs.
func (w *Worker[In, Out]) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case job, ok := <-w.jobQueue:
				if !ok {
					return // Job queue closed
				}
				value, err := w.processor(job.Input)
				result := Result[Out]{
					JobID: job.ID,
					Value: value,
					Error: err,
				}
				select {
				case w.resultChan <- result:
				case <-w.ctx.Done():
					return
				}
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

// WorkerPool manages a pool of workers.
type WorkerPool[In, Out any] struct {
	workers    []*Worker[In, Out]
	jobQueue   chan Job[In, Out]
	resultChan chan Result[Out]
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	size       int
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool[In, Out any](ctx context.Context, size int, queueSize int, processor func(In) (Out, error)) *WorkerPool[In, Out] {
	poolCtx, cancel := context.WithCancel(ctx)

	pool := &WorkerPool[In, Out]{
		workers:    make([]*Worker[In, Out], size),
		jobQueue:   make(chan Job[In, Out], queueSize),
		resultChan: make(chan Result[Out], queueSize),
		ctx:        poolCtx,
		cancel:     cancel,
		size:       size,
	}

	// Create workers
	for i := 0; i < size; i++ {
		pool.workers[i] = &Worker[In, Out]{
			id:         i,
			jobQueue:   pool.jobQueue,
			resultChan: pool.resultChan,
			processor:  processor,
			ctx:        poolCtx,
			wg:         &pool.wg,
		}
	}

	return pool
}

// Start starts all workers in the pool.
func (wp *WorkerPool[In, Out]) Start() {
	for _, worker := range wp.workers {
		worker.Start()
	}

	// Close result channel when all workers complete
	go func() {
		wp.wg.Wait()
		close(wp.resultChan)
	}()
}

// Submit submits a job to the worker pool.
func (wp *WorkerPool[In, Out]) Submit(job Job[In, Out]) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	}
}

// Results returns the results channel.
func (wp *WorkerPool[In, Out]) Results() <-chan Result[Out] {
	return wp.resultChan
}

// Shutdown gracefully shuts down the worker pool.
func (wp *WorkerPool[In, Out]) Shutdown() {
	close(wp.jobQueue)
	wp.wg.Wait()
}

// Stop immediately stops all workers.
func (wp *WorkerPool[In, Out]) Stop() {
	wp.cancel()
	close(wp.jobQueue)
}

// Size returns the number of workers in the pool.
func (wp *WorkerPool[In, Out]) Size() int {
	return wp.size
}

// SimpleWorkerPool provides a simpler interface for basic use cases.
type SimpleWorkerPool struct {
	workers    int
	jobQueue   chan func()
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewSimpleWorkerPool creates a simple worker pool.
func NewSimpleWorkerPool(ctx context.Context, workers int) *SimpleWorkerPool {
	poolCtx, cancel := context.WithCancel(ctx)
	return &SimpleWorkerPool{
		workers:  workers,
		jobQueue: make(chan func(), workers*2),
		ctx:      poolCtx,
		cancel:   cancel,
	}
}

// Start starts the simple worker pool.
func (sp *SimpleWorkerPool) Start() {
	for i := 0; i < sp.workers; i++ {
		sp.wg.Add(1)
		go func() {
			defer sp.wg.Done()
			for {
				select {
				case job, ok := <-sp.jobQueue:
					if !ok {
						return
					}
					job()
				case <-sp.ctx.Done():
					return
				}
			}
		}()
	}
}

// Submit submits a job function to the pool.
func (sp *SimpleWorkerPool) Submit(job func()) error {
	select {
	case sp.jobQueue <- job:
		return nil
	case <-sp.ctx.Done():
		return sp.ctx.Err()
	}
}

// Shutdown gracefully shuts down the pool.
func (sp *SimpleWorkerPool) Shutdown() {
	close(sp.jobQueue)
	sp.wg.Wait()
}

// Stop immediately stops the pool.
func (sp *SimpleWorkerPool) Stop() {
	sp.cancel()
	close(sp.jobQueue)
}
