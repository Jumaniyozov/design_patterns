# Worker Pool Pattern

## Overview

The Worker Pool Pattern is a concurrency design pattern that maintains a fixed number of goroutines (workers) that process tasks from a shared work queue. Instead of creating a new goroutine for each task, a pool of reusable workers processes tasks from a queue, providing efficient resource management and rate limiting.

## Problem

Creating a new goroutine for each task creates several problems:

- **Resource exhaustion**: Creating too many goroutines consumes memory and context switches
- **Unbounded concurrency**: No control over maximum concurrent operations
- **Rate limiting**: No way to throttle concurrent work
- **Inefficiency**: Creating/destroying goroutines has overhead
- **Scalability limits**: Systems can only handle so many concurrent goroutines
- **Resource contention**: Unlimited goroutines compete for system resources

### Real-World Context

Imagine a web server handling incoming requests. Creating a new goroutine for each request seems natural, but under high load, thousands of goroutines would be created, exhausting memory. A worker pool maintains a fixed number of request handlers that process requests from a queue, ensuring controlled resource usage.

## Why Use This Pattern?

- **Resource control**: Limit concurrent operations to available resources
- **Efficiency**: Reuse goroutines instead of creating/destroying
- **Rate limiting**: Control throughput of work processing
- **Backpressure**: Queue provides natural backpressure on tasks
- **Scalability**: Handle unbounded task load with bounded resources
- **Performance**: Reduced memory usage and context switching

## When to Use

- Server handling many client requests
- Task processing systems with high volume
- Database connection management
- Rate limiting concurrent operations
- Controlling resource usage in concurrent systems
- Processing jobs from a queue
- Load balancing work across workers

## When NOT to Use

- Only a few tasks to process (overhead not justified)
- Tasks are very short (overhead exceeds benefit)
- Unlimited concurrency is acceptable
- Task arrival is unpredictable and bursty
- Must process tasks in strict order

## Implementation Guidelines

1. **Worker count**: Fixed number of goroutines processing tasks
2. **Job queue**: Channel holding tasks waiting to be processed
3. **Job definition**: Task that workers execute
4. **Dispatch**: Send jobs to job channel
5. **Result collection**: Gather results from workers
6. **Graceful shutdown**: Wait for in-flight tasks, close workers

## Go Idioms

Worker pools are idiomatic Go:

- Goroutines are the unit of concurrency
- Channels pass work between goroutines
- sync.WaitGroup coordinates goroutine shutdown
- Context can signal cancellation

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│      Worker Pool Pattern Architecture        │
└──────────────────────────────────────────────┘

Fixed Pool of Workers:

┌──────────────────────────────┐
│  Job Queue (Channel)         │
│  ┌──────────────────────────┐│
│  │ Job1 Job2 Job3 Job4 ...  ││
│  └──────────────────────────┘│
└──┬───┬───┬───┬───┬───────────┘
   │   │   │   │   │
   ▼   ▼   ▼   ▼   ▼
┌─────┐┌─────┐┌─────┐
│ W1  |│ W2  ││ W3  │
│     ││     ││     │
└─────┘└─────┘└─────┘
(Fixed number of workers)

Results
   │
   ▼
┌──────────────────────┐
│ Result Queue         │
│ (Aggregated output)  │
└──────────────────────┘


Comparison: Creating Goroutines vs Worker Pool

WITHOUT Worker Pool:
┌──────────────────┐
│ Client Requests  │ 1000 requests arrive
└────────┬─────────┘
         │
    ┌────┴────┬─────────┬────────┐
    │         │         │        │
    ▼         ▼         ▼        ▼
   G1        G2        G3  ... G1000
    │         │         │        │
    └─────────┴─────────┴────────┘

   Problems:
   - 1000 goroutines created
   - High memory usage
   - Too much context switching
   - Resource exhaustion

WITH Worker Pool:
┌──────────────────┐
│ Client Requests  │ 1000 requests arrive
└────────┬─────────┘
         │
    ┌────┴────────────┐
    │  Job Queue      │
    │ (Buffered ch)   │
    └────────┬────────┘
             │
    ┌────────┼────────┐
    │        │        │
    ▼        ▼        ▼
   W1       W2       W3  (Fixed: e.g., 10 workers)
    │        │        │
    └────────┴────────┘

   Benefits:
   - Only 10 goroutines
   - Low memory usage
   - Minimal context switching
   - Graceful rate limiting


Work Queue Buffering:

Job arrivals: ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔ (fast, many jobs)

            ┌────────────────┐
            │   Job Queue    │ Buffers excess jobs
            │ (size: 100)    │
            └────────────────┘
                    │
          ┌─────────┴─────────┐
          │        │          │
          ▼        ▼          ▼
          W1       W2    ...  W10

Processing rate: ▄▄▄ (slower, fewer workers)

Queue acts as buffer, preventing request loss
while workers process at sustainable rate.


Graceful Shutdown:

┌──────────────────────────────┐
│ New requests: BLOCKED        │
│ In-flight tasks: Complete    │
└──────────────────────────────┘
         │
         ▼
    ┌────────────┐
    │ Job Queue  │ ← No new jobs
    │ (Closed)   │
    └────────────┘
         │
    ┌────┴────┬────┐
    │         │    │
    ▼         ▼    ▼
   W1  ...   W10  (Processing remaining)
    │         │    │
    └─────────┴────┘
         │
         ▼
    All workers done
    Shutdown complete
```

## Real-World Examples

### 1. Basic Worker Pool

```go
type Job struct {
    id  int
    fn  func() error
}

type Result struct {
    job   Job
    err   error
}

type WorkerPool struct {
    jobs    chan Job
    results chan Result
    workers int
}

func NewWorkerPool(workers int, bufferSize int) *WorkerPool {
    return &WorkerPool{
        jobs:    make(chan Job, bufferSize),
        results: make(chan Result, bufferSize),
        workers: workers,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    for job := range wp.jobs {
        wp.results <- Result{
            job: job,
            err: job.fn(),
        }
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobs <- job
}

func (wp *WorkerPool) Stop() {
    close(wp.jobs)
}
```

### 2. HTTP Request Handler Pool

```go
type RequestHandler struct {
    jobs    chan *http.Request
    workers int
}

func NewRequestHandler(workers int) *RequestHandler {
    return &RequestHandler{
        jobs:    make(chan *http.Request, workers*2),
        workers: workers,
    }
}

func (rh *RequestHandler) Start() {
    for i := 0; i < rh.workers; i++ {
        go rh.worker()
    }
}

func (rh *RequestHandler) worker() {
    for req := range rh.jobs {
        rh.handleRequest(req)
    }
}

func (rh *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    select {
    case rh.jobs <- r:
        // Job queued
    case <-time.After(1 * time.Second):
        http.Error(w, "Queue full", http.StatusServiceUnavailable)
    }
}
```

### 3. Database Connection Pool with Workers

```go
type DBWorker struct {
    queries chan Query
    db      *sql.DB
    results chan QueryResult
}

func NewDBWorkerPool(poolSize int, workers int) *DBWorker {
    return &DBWorker{
        queries: make(chan Query, poolSize),
        results: make(chan QueryResult),
    }
}

func (dw *DBWorker) Start() {
    for i := 0; i < dw.workers; i++ {
        go func() {
            for query := range dw.queries {
                rows, err := dw.db.Query(query.sql)
                dw.results <- QueryResult{
                    query: query,
                    rows:  rows,
                    err:   err,
                }
            }
        }()
    }
}
```

## Key Advantages

- **Resource control**: Limit concurrent operations to available resources
- **Efficiency**: Reuse goroutines instead of creating/destroying
- **Rate limiting**: Control throughput of work processing
- **Backpressure**: Queue provides natural flow control
- **Scalability**: Handle unbounded task load with bounded resources
- **Performance**: Reduced memory usage and context switching
- **Graceful shutdown**: Can wait for in-flight tasks
- **Queue monitoring**: Monitor queue depth for system health

## Key Gotchas

- **Overhead**: For small number of fast tasks, overhead not justified
- **Deadlocks**: Improper synchronization can cause deadlocks
- **Buffer sizing**: Too large wastes memory, too small causes blocking
- **Goroutine leaks**: Must properly close channels and sync on shutdown
- **Error handling**: Errors in workers should not crash pool
- **Monitoring**: Need to monitor queue depth and worker health
- **Fair scheduling**: No guarantee of fair job distribution
- **Task ordering**: Tasks may not complete in submission order
