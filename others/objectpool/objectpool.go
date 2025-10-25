// Package objectpool demonstrates the Object Pool pattern.
// It manages a pool of reusable objects to avoid expensive creation/destruction,
// critical for performance-sensitive applications.
package objectpool

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Connection represents an expensive resource
type Connection struct {
	ID        int
	connected bool
	lastUsed  time.Time
}

// Connect simulates expensive connection setup
func (c *Connection) Connect() error {
	time.Sleep(100 * time.Millisecond) // Simulate expensive operation
	c.connected = true
	c.lastUsed = time.Now()
	return nil
}

// Close closes the connection
func (c *Connection) Close() {
	c.connected = false
}

// Execute simulates query execution
func (c *Connection) Execute(query string) string {
	if !c.connected {
		return "Error: not connected"
	}
	c.lastUsed = time.Now()
	return fmt.Sprintf("Connection %d executed: %s", c.ID, query)
}

// Reset resets connection state for reuse
func (c *Connection) Reset() {
	c.lastUsed = time.Now()
}

// ConnectionPool manages a pool of connections
type ConnectionPool struct {
	available chan *Connection
	inUse     map[int]*Connection
	mu        sync.Mutex
	maxSize   int
	nextID    int
}

// NewConnectionPool creates a connection pool
func NewConnectionPool(size int) *ConnectionPool {
	return &ConnectionPool{
		available: make(chan *Connection, size),
		inUse:     make(map[int]*Connection),
		maxSize:   size,
		nextID:    1,
	}
}

// Acquire gets a connection from the pool
func (cp *ConnectionPool) Acquire() (*Connection, error) {
	select {
	case conn := <-cp.available:
		// Reuse existing connection
		cp.mu.Lock()
		cp.inUse[conn.ID] = conn
		cp.mu.Unlock()
		conn.Reset()
		return conn, nil
	default:
		// Create new connection if pool not full
		cp.mu.Lock()
		defer cp.mu.Unlock()

		if len(cp.inUse)+len(cp.available) >= cp.maxSize {
			return nil, errors.New("pool exhausted")
		}

		conn := &Connection{ID: cp.nextID}
		cp.nextID++
		if err := conn.Connect(); err != nil {
			return nil, err
		}
		cp.inUse[conn.ID] = conn
		return conn, nil
	}
}

// Release returns connection to pool
func (cp *ConnectionPool) Release(conn *Connection) {
	cp.mu.Lock()
	delete(cp.inUse, conn.ID)
	cp.mu.Unlock()

	select {
	case cp.available <- conn:
		// Successfully returned to pool
	default:
		// Pool full, close connection
		conn.Close()
	}
}

// Stats returns pool statistics
func (cp *ConnectionPool) Stats() string {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	return fmt.Sprintf("Pool: %d available, %d in use, max: %d",
		len(cp.available), len(cp.inUse), cp.maxSize)
}

// WorkerPool demonstrates goroutine pool pattern
type WorkerPool struct {
	tasks   chan func()
	workers int
	wg      sync.WaitGroup
}

// NewWorkerPool creates a worker pool
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		tasks:   make(chan func(), 100),
		workers: workers,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes tasks
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for task := range wp.tasks {
		task()
	}
}

// Submit submits a task to the pool
func (wp *WorkerPool) Submit(task func()) {
	wp.tasks <- task
}

// Shutdown stops the worker pool
func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}

// BufferPool demonstrates sync.Pool usage
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a buffer pool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 1024)
			},
		},
	}
}

// Get acquires a buffer from pool
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put returns buffer to pool
func (bp *BufferPool) Put(buf []byte) {
	// Clear buffer before returning
	for i := range buf {
		buf[i] = 0
	}
	bp.pool.Put(buf)
}

// ObjectPool is a generic object pool implementation
type ObjectPool struct {
	pool    chan interface{}
	factory func() interface{}
	reset   func(interface{})
}

// NewObjectPool creates a generic object pool
func NewObjectPool(size int, factory func() interface{}, reset func(interface{})) *ObjectPool {
	return &ObjectPool{
		pool:    make(chan interface{}, size),
		factory: factory,
		reset:   reset,
	}
}

// Get acquires an object
func (op *ObjectPool) Get() interface{} {
	select {
	case obj := <-op.pool:
		return obj
	default:
		return op.factory()
	}
}

// Put returns an object to pool
func (op *ObjectPool) Put(obj interface{}) {
	if op.reset != nil {
		op.reset(obj)
	}
	select {
	case op.pool <- obj:
	default:
		// Pool full, discard
	}
}
