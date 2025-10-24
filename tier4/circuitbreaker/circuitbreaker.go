// Package circuitbreaker implements the Circuit Breaker pattern.
//
// The Circuit Breaker pattern prevents cascading failures in distributed systems
// by wrapping service calls and monitoring for failures. When failures exceed a
// threshold, the circuit "opens" and subsequent calls fail fast without attempting
// the operation.
//
// States:
// - Closed: Normal operation, requests pass through
// - Open: Too many failures, requests fail immediately
// - Half-Open: Testing if service recovered, limited requests allowed
//
// Key characteristics:
// - Failure detection and fast-fail
// - Automatic recovery attempts
// - Prevents resource exhaustion
// - Protects downstream services
package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrCircuitOpen is returned when the circuit is open.
	ErrCircuitOpen = errors.New("circuit breaker is open")

	// ErrTooManyRequests is returned when too many requests are made in half-open state.
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// State represents the circuit breaker state.
type State int

const (
	// StateClosed means requests pass through normally.
	StateClosed State = iota

	// StateOpen means circuit is open, requests fail fast.
	StateOpen

	// StateHalfOpen means testing if service recovered.
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// Config holds circuit breaker configuration.
type Config struct {
	// MaxRequests is the maximum number of requests allowed in half-open state.
	MaxRequests uint32

	// Interval is the cyclic period in closed state to clear internal counts.
	// If 0, never clears.
	Interval time.Duration

	// Timeout is how long to wait in open state before transitioning to half-open.
	Timeout time.Duration

	// ReadyToTrip is called when a request fails in closed state.
	// If it returns true, the circuit breaker trips to open state.
	ReadyToTrip func(counts Counts) bool

	// OnStateChange is called whenever the state changes.
	OnStateChange func(from State, to State)
}

// Counts holds statistics about requests.
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu           sync.RWMutex
	config       Config
	state        State
	counts       Counts
	expiry       time.Time
	halfOpenReqs uint32
}

// New creates a new circuit breaker.
func New(config Config) *CircuitBreaker {
	cb := &CircuitBreaker{
		config: config,
		state:  StateClosed,
		counts: Counts{},
	}

	if cb.config.MaxRequests == 0 {
		cb.config.MaxRequests = 1
	}

	if cb.config.Timeout == 0 {
		cb.config.Timeout = 60 * time.Second
	}

	if cb.config.ReadyToTrip == nil {
		cb.config.ReadyToTrip = func(counts Counts) bool {
			return counts.ConsecutiveFailures > 5
		}
	}

	return cb
}

// Execute runs the given function if the circuit breaker allows it.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	generation, err := cb.beforeRequest()
	if err != nil {
		return err
	}

	err = fn()
	cb.afterRequest(generation, err == nil)

	return err
}

// ExecuteWithContext runs the given function with context support.
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	generation, err := cb.beforeRequest()
	if err != nil {
		return err
	}

	err = fn(ctx)
	cb.afterRequest(generation, err == nil)

	return err
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Counts returns the current counts.
func (cb *CircuitBreaker) Counts() Counts {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.counts
}

// beforeRequest checks if the request should be allowed.
func (cb *CircuitBreaker) beforeRequest() (uint32, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	state := cb.state

	// Check if we should transition from open to half-open
	if state == StateOpen && cb.expiry.Before(now) {
		cb.setState(StateHalfOpen)
		return 0, nil
	}

	// Handle each state
	switch state {
	case StateClosed:
		// Check if interval expired and we should reset counts
		if cb.config.Interval > 0 && cb.expiry.Before(now) {
			cb.counts = Counts{}
			cb.expiry = now.Add(cb.config.Interval)
		}
		return 0, nil

	case StateOpen:
		return 0, ErrCircuitOpen

	case StateHalfOpen:
		// Check if too many requests in half-open state
		if cb.halfOpenReqs >= cb.config.MaxRequests {
			return 0, ErrTooManyRequests
		}
		cb.halfOpenReqs++
		return cb.halfOpenReqs, nil

	default:
		return 0, errors.New("unknown circuit breaker state")
	}
}

// afterRequest records the result of a request.
func (cb *CircuitBreaker) afterRequest(generation uint32, success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Ignore if generation doesn't match (stale request)
	if cb.state == StateHalfOpen && generation != cb.halfOpenReqs {
		return
	}

	cb.counts.Requests++

	if success {
		cb.onSuccess()
	} else {
		cb.onFailure()
	}
}

// onSuccess handles a successful request.
func (cb *CircuitBreaker) onSuccess() {
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0

	// Transition from half-open to closed after enough successes
	if cb.state == StateHalfOpen && cb.counts.ConsecutiveSuccesses >= cb.config.MaxRequests {
		cb.setState(StateClosed)
	}
}

// onFailure handles a failed request.
func (cb *CircuitBreaker) onFailure() {
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0

	// Check if we should trip the circuit
	switch cb.state {
	case StateClosed:
		if cb.config.ReadyToTrip(cb.counts) {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		// Any failure in half-open state reopens the circuit
		cb.setState(StateOpen)
	}
}

// setState transitions to a new state.
func (cb *CircuitBreaker) setState(newState State) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.counts = Counts{}

	switch newState {
	case StateClosed:
		if cb.config.Interval > 0 {
			cb.expiry = time.Now().Add(cb.config.Interval)
		}
	case StateOpen:
		cb.expiry = time.Now().Add(cb.config.Timeout)
	case StateHalfOpen:
		cb.halfOpenReqs = 0
	}

	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(oldState, newState)
	}
}

// Reset manually resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.setState(StateClosed)
}
