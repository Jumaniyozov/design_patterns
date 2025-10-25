// Package semaphore demonstrates the Semaphore pattern.
// It limits concurrent access to resources using permits,
// essential for rate limiting and resource management.
package semaphore

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Semaphore using buffered channel
type Semaphore struct {
	permits chan struct{}
}

// NewSemaphore creates a semaphore with n permits
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: make(chan struct{}, n),
	}
}

// Acquire acquires a permit (blocks if none available)
func (s *Semaphore) Acquire() {
	s.permits <- struct{}{}
}

// TryAcquire tries to acquire without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.permits <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithTimeout tries to acquire with timeout
func (s *Semaphore) AcquireWithTimeout(timeout time.Duration) bool {
	select {
	case s.permits <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Release releases a permit
func (s *Semaphore) Release() {
	<-s.permits
}

// AvailablePermits returns number of available permits
func (s *Semaphore) AvailablePermits() int {
	return cap(s.permits) - len(s.permits)
}

// WeightedSemaphore allows acquiring multiple permits
type WeightedSemaphore struct {
	current int
	max     int
	mu      sync.Mutex
	cond    *sync.Cond
}

// NewWeightedSemaphore creates a weighted semaphore
func NewWeightedSemaphore(max int) *WeightedSemaphore {
	ws := &WeightedSemaphore{
		current: 0,
		max:     max,
	}
	ws.cond = sync.NewCond(&ws.mu)
	return ws
}

// Acquire acquires n permits
func (ws *WeightedSemaphore) Acquire(n int) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for ws.current+n > ws.max {
		ws.cond.Wait()
	}
	ws.current += n
}

// Release releases n permits
func (ws *WeightedSemaphore) Release(n int) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.current -= n
	if ws.current < 0 {
		ws.current = 0
	}
	ws.cond.Broadcast()
}

// Available returns available permits
func (ws *WeightedSemaphore) Available() int {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.max - ws.current
}

// RateLimiter limits operations per time window
type RateLimiter struct {
	rate     int
	interval time.Duration
	tokens   chan struct{}
	stop     chan struct{}
}

// NewRateLimiter creates a rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		rate:     rate,
		interval: interval,
		tokens:   make(chan struct{}, rate),
		stop:     make(chan struct{}),
	}

	// Fill initial tokens
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	// Start token refill goroutine
	go rl.refill()

	return rl
}

func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(rl.interval / time.Duration(rl.rate))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Bucket full
			}
		case <-rl.stop:
			return
		}
	}
}

// Allow checks if operation is allowed
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait waits until operation is allowed
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stop)
}

// ConcurrencyLimiter limits concurrent operations
type ConcurrencyLimiter struct {
	sem *Semaphore
}

// NewConcurrencyLimiter creates a concurrency limiter
func NewConcurrencyLimiter(maxConcurrent int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		sem: NewSemaphore(maxConcurrent),
	}
}

// Execute executes a function with concurrency limit
func (cl *ConcurrencyLimiter) Execute(ctx context.Context, fn func() error) error {
	// Try to acquire
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	cl.sem.Acquire()
	defer cl.sem.Release()

	return fn()
}

// ExecuteAll executes all functions with concurrency limit
func (cl *ConcurrencyLimiter) ExecuteAll(ctx context.Context, tasks []func() error) []error {
	results := make([]error, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t func() error) {
			defer wg.Done()
			results[index] = cl.Execute(ctx, t)
		}(i, task)
	}

	wg.Wait()
	return results
}

// Example: HTTP Client Pool with semaphore
type HTTPClientPool struct {
	sem        *Semaphore
	maxClients int
}

// NewHTTPClientPool creates an HTTP client pool
func NewHTTPClientPool(maxClients int) *HTTPClientPool {
	return &HTTPClientPool{
		sem:        NewSemaphore(maxClients),
		maxClients: maxClients,
	}
}

// Do performs an HTTP request with concurrency limit
func (p *HTTPClientPool) Do(url string) (string, error) {
	p.sem.Acquire()
	defer p.sem.Release()

	// Simulate HTTP request
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("Response from %s", url), nil
}
