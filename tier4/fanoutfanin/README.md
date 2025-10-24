# Fan-Out/Fan-In Pattern

## Overview

The Fan-Out/Fan-In Pattern is a concurrency design pattern that distributes a task across multiple goroutines (fan-out) and then collects their results (fan-in). It enables parallel processing of independent work and aggregation of results for maximum throughput.

## Problem

When you have a task that can be split into independent subtasks, sequential processing wastes computational resources:

- **Underutilized resources**: Single goroutine doesn't take advantage of multiple cores
- **Slow processing**: Processing items one at a time is inefficient
- **Wasted CPU**: Modern systems have many cores that sit idle
- **Poor throughput**: Sequential processing can't match parallel performance
- **Scalability issues**: Vertical scaling by adding cores doesn't help without parallelization
- **Result aggregation**: No clear way to collect results from parallel work

### Real-World Context

Imagine processing a batch of user records where each record needs validation, enrichment, and storage. Processing them sequentially takes N seconds. With fan-out/fan-in, you can process 10 records in parallel, completing in roughly N/10 seconds. You fan-out 10 goroutines to process different records, then fan-in to collect all results.

## Why Use This Pattern?

- **Parallelism**: Process independent work in parallel
- **Throughput**: Dramatically increase processing speed
- **Resource utilization**: Take full advantage of multi-core systems
- **Go concurrency**: Perfect fit for goroutines and channels
- **Scalability**: Easy to adjust parallelism level
- **Simple coordination**: Channels handle synchronization automatically

## When to Use

- Independent tasks that can be processed in parallel
- Batch processing of items
- Parallel API calls to multiple services
- Map-reduce operations
- Processing pipelines with parallel stages
- Aggregating results from multiple sources
- Load distribution across workers

## When NOT to Use

- Tasks are dependent or require sequential ordering
- Coordination overhead exceeds parallelism benefit
- Single-threaded performance sufficient
- Tasks are very fast (context switching overhead matters)
- Result ordering must be preserved (use ordered channels)
- Memory constraints prevent creating many goroutines

## Implementation Guidelines

1. **Fan-out**: Create goroutines to process items in parallel
2. **Work distribution**: Send items to workers via channel
3. **Worker processing**: Each goroutine processes its item
4. **Result collection**: Gather results from worker goroutines
5. **Fan-in**: Combine results from parallel workers

## Go Idioms

Fan-out/fan-in is idiomatic in Go:

- Goroutines are cheap (thousands feasible)
- Channels naturally distribute work and collect results
- sync.WaitGroup coordinates goroutines
- Patterns compose naturally with pipelines

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│    Fan-Out/Fan-In Pattern Architecture       │
└──────────────────────────────────────────────┘

Fan-Out: Distribute work

     Input Channel
            │
    ┌───────┴────────────┬──────────┐
    │                    │          │
    ▼                    ▼          ▼
┌─────────┐          ┌─────────┐ ┌─────────┐
│Worker 1 │          │Worker 2 │ │Worker 3 │
│Process  │          │Process  │ │Process  │
│Item A   │          │Item B   │ │Item C   │
└────┬────┘          └────┬────┘ └────┬────┘
     │                    │           │
     └────────┬───────────┴───────────┘
              │
           Fan-In: Collect results
              │
        ┌─────┴──────┐
        │            │
        ▼            ▼
   Result A    Result B
   Result C

Parallel Processing Example:

Input Items: [Item1, Item2, Item3, Item4, Item5]

Sequential (5 units of work, 5 time units):
Time: 1───2───3───4───5
      │   │   │   │   │
   Item1 Item2 Item3 Item4 Item5
   Total: 5 time units

Parallel with 3 Workers:
Worker1: Item1──Item4
Worker2: Item2──Item5
Worker3: Item3
Time: 1───2
Total: 2 time units (much faster!)

Process Flow:

┌──────────────┐
│ Input List   │ [1, 2, 3, 4, 5, 6, 7, 8]
└──────┬───────┘
       │
       ▼
┌─────────────────┐
│ Create Workers  │ N goroutines
└──────┬──────────┘
       │
       ▼
    ┌──┴──────────────────┐
    │                     │
    ▼                     ▼
Distribute Items      Worker Pool
    │                     │
    │    Item 1 ─────────>│
    │    Item 2 ─────────>│ Processing
    │    Item 3 ─────────>│ in parallel
    │    Item 4 ─────────>│
    │    Item 5 ─────────>│
    │    ...              │
    │                     │
    └─────────────────────┘
                │
                ▼
         Result Channel
                │
      ┌─────────┴─────────┐
      ▼                   ▼   
    Collect Results    Aggregate
        Return              
    
```

## Real-World Examples

### 1. Batch Item Processing

```go
func ProcessItems(items []Item) []Result {
    workerCount := 10
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))

    // Fan-out: Start workers
    for i := 0; i < workerCount; i++ {
        go func() {
            for item := range jobs {
                results <- processItem(item)
            }
        }()
    }

    // Distribute items
    go func() {
        for _, item := range items {
            jobs <- item
        }
        close(jobs)
    }()

    // Fan-in: Collect results
    var processed []Result
    for i := 0; i < len(items); i++ {
        processed = append(processed, <-results)
    }
    return processed
}
```

### 2. Parallel API Calls

```go
func FetchUserData(userIDs []string) []UserData {
    results := make(chan UserData, len(userIDs))
    var wg sync.WaitGroup

    // Fan-out: Call API for each user in parallel
    for _, id := range userIDs {
        wg.Add(1)
        go func(userID string) {
            defer wg.Done()
            data := fetchFromAPI(userID)
            results <- data
        }(id)
    }

    // Fan-in: Collect results
    go func() {
        wg.Wait()
        close(results)
    }()

    var users []UserData
    for user := range results {
        users = append(users, user)
    }
    return users
}
```

### 3. Data Processing Pipeline

```go
func ProcessLargeDataset(data []DataItem) []ProcessedItem {
    workers := 8

    // Fan-out stage
    itemChan := distribute(data, workers)

    // Processing (workers process items)
    resultChan := process(itemChan, workers)

    // Fan-in stage
    return collect(resultChan, len(data))
}

func distribute(items []DataItem, workers int) <-chan DataItem {
    out := make(chan DataItem)
    go func() {
        for _, item := range items {
            out <- item
        }
        close(out)
    }()
    return out
}

func process(in <-chan DataItem, workers int) <-chan ProcessedItem {
    out := make(chan ProcessedItem)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range in {
                out <- processItem(item)
            }
        }()
    }

    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

## Key Advantages

- **Parallelism**: Process work in parallel on multiple cores
- **Throughput**: Dramatically increase processing speed
- **Resource utilization**: Full use of multi-core systems
- **Scalability**: Easy to adjust worker count
- **Go idiomatic**: Natural fit for goroutines and channels
- **Simplicity**: Pattern is straightforward to understand
- **Flexibility**: Works with different data types and processing

## Key Gotchas

- **Context switching**: Too many goroutines can hurt performance
- **Channel overhead**: May not be worth it for very fast operations
- **Memory usage**: Buffered channels hold results in memory
- **Result ordering**: Results arrive in unpredictable order
- **Worker coordination**: Proper synchronization is critical
- **Resource limits**: Can't create unlimited goroutines
- **Error handling**: Need to handle errors from individual workers
- **Backpressure**: Unbounded work can exhaust memory
- **Goroutine leaks**: Must ensure all goroutines complete properly
