# Pipeline Pattern

## Overview

The Pipeline Pattern is a concurrency design pattern that processes data through a series of processing stages connected by channels. Each stage runs concurrently, receiving input from the previous stage and sending output to the next stage, creating a data processing pipeline.

## Problem

When processing data through multiple sequential transformation steps, several challenges arise:

- **Sequential bottlenecks**: Each stage must complete before the next starts
- **Wasted resources**: Stages can't work on different data items simultaneously
- **Poor throughput**: Processing speed limited by slowest stage
- **Difficult composition**: Hard to build complex processing chains
- **State management**: Tracking data through multiple stages becomes messy
- **Resource underutilization**: Not taking advantage of concurrent processing

### Real-World Context

Consider an image processing pipeline: load → resize → apply filter → compress → save. With sequential processing, each image completes all stages before the next image starts. With pipeline pattern, you can have 5 images in different stages simultaneously, dramatically increasing throughput.

## Why Use This Pattern?

- **Concurrency**: Process multiple items at different stages simultaneously
- **Throughput**: Each stage works independently, improving overall speed
- **Composition**: Easy to build complex pipelines from simple stages
- **Resource utilization**: Better use of multi-core systems
- **Go idioms**: Perfect use of goroutines and channels
- **Scalability**: Scale individual stages independently

## When to Use

- Multi-stage data transformations
- Stream processing of large datasets
- ETL (Extract, Transform, Load) pipelines
- Image/video processing chains
- Request processing through multiple handlers
- Data aggregation from multiple sources
- Complex business logic pipelines

## When NOT to Use

- Single transformation step (overhead not justified)
- Stages must run sequentially (can't parallelize)
- Data must be processed in strict order (dependencies)
- Overhead exceeds benefits for small datasets
- Simple linear code is clearer

## Implementation Guidelines

1. **Stage definition**: Each stage is a goroutine with input/output channels
2. **Input channel**: Receives data from previous stage
3. **Processing**: Transform the data
4. **Output channel**: Sends results to next stage
5. **Closure**: Concurrent execution through goroutines
6. **Channel management**: Properly close channels when complete

## Go Idioms

Pipelines are a fundamental Go pattern:

- Goroutines are cheap (thousands feasible)
- Channels naturally represent data flow
- Range over channels consumes items naturally
- Composition of stages creates larger pipelines

## Visual Schema

```go
┌──────────────────────────────────────────────┐
│        Pipeline Pattern Architecture         │
└──────────────────────────────────────────────┘

Simple Pipeline:

Input ────> Stage1 ────> Stage2 ────> Stage3 ────> Output
  │           │            │           │
  │ receive   │ process    │ process   │ send
  │ send      │ receive    │ receive   │
  │           │ send       │ send      │

Time Diagram (3 items, 3 stages):

Input     Item1    Item2    Item3
  │         │        │        │
  ▼         ▼        ▼        ▼
Stage1: ┌─────┐                        (Item1)
        │ W1  │   ┌─────┐              (Item2)
        │     │   │ W2  │ ┌─────┐      (Item3)
        └─────┘   │     │ │ W3  │
                  └─────┘ │     │
                          └─────┘
                             ▼
                           Output

All stages working in parallel!
While Stage1 processes Item2,
Stage2 processes Item1,
Stage3 processes Item0


Data Flow Example (Image Processing):

Input Images
    │
    ▼
┌──────────────┐
│ Load Stage   │ Load image from disk
│ (goroutine)  │
└────────┬─────┘
         │ image channel
         ▼
    ┌──────────────┐
    │ Resize Stage │ Resize to 800x600
    │ (goroutine)  │
    └────────┬─────┘
             │ resized channel
             ▼
        ┌──────────────┐
        │ Filter Stage │ Apply blur/sharpen
        │ (goroutine)  │
        └────────┬─────┘
                 │ filtered channel
                 ▼
            ┌──────────────┐
            │ Compress     │ JPEG compression
            │ (goroutine)  │
            └────────┬─────┘
                     │ compressed channel
                     ▼
                Output Images


Pipeline Code Structure:

func Pipeline(input <-chan Data) <-chan Result {
    // Stage 1: Receive from input
    stage1 := make(chan Intermediate)
    go func() {
        for item := range input {
            stage1 <- process1(item)
        }
        close(stage1)
    }()

    // Stage 2: Receive from stage1, send to stage2
    stage2 := make(chan Result)
    go func() {
        for item := range stage1 {
            stage2 <- process2(item)
        }
        close(stage2)
    }()

    // Return final stage output
    return stage2
}

Multiple Worker Stages:

Generate ──┬──> Worker1 ┐
           ├──> Worker2 ├──> Merge ──> Output
           └──> Worker3 ┘

Allows scaling individual stages:
- Slow stage gets more workers
- Fast stage gets fewer workers
```

## Real-World Examples

### 1. Data Processing Pipeline

```go
func pipeline(input <-chan int) <-chan int {
    // Stage 1: Square numbers
    squared := make(chan int)
    go func() {
        for n := range input {
            squared <- n * n
        }
        close(squared)
    }()

    // Stage 2: Double squared values
    doubled := make(chan int)
    go func() {
        for n := range squared {
            doubled <- n * 2
        }
        close(doubled)
    }()

    return doubled
}

// Usage
input := make(chan int)
result := pipeline(input)

go func() {
    for i := 1; i <= 5; i++ {
        input <- i
    }
    close(input)
}()

for val := range result {
    fmt.Println(val) // 2, 8, 18, 32, 50
}
```

### 2. Image Processing Pipeline

```go
type Image struct {
    data []byte
    name string
}

func imageProcessingPipeline(input <-chan Image) <-chan Image {
    resized := resizeStage(input)
    filtered := filterStage(resized)
    return compressStage(filtered)
}

func resizeStage(input <-chan Image) <-chan Image {
    out := make(chan Image)
    go func() {
        for img := range input {
            out <- resizeImage(img)
        }
        close(out)
    }()
    return out
}

func filterStage(input <-chan Image) <-chan Image {
    out := make(chan Image)
    go func() {
        for img := range input {
            out <- applyFilter(img)
        }
        close(out)
    }()
    return out
}
```

### 3. ETL Pipeline with Multiple Workers

```go
func ETLPipeline(input <-chan Record, workers int) <-chan LoadedRecord {
    // Extract stage
    extracted := extractStage(input)

    // Transform stage with multiple workers
    transformed := transformStage(extracted, workers)

    // Load stage
    return loadStage(transformed)
}

func transformStage(input <-chan ExtractedRecord, workers int) <-chan TransformedRecord {
    out := make(chan TransformedRecord)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for record := range input {
                out <- transformRecord(record)
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

- **Concurrency**: Process multiple items at different stages simultaneously
- **Throughput**: Each stage works independently, improving speed
- **Composition**: Easy to build complex pipelines from simple stages
- **Scalability**: Scale individual stages independently by adding workers
- **Separation of concerns**: Each stage has single responsibility
- **Resource utilization**: Better use of multi-core systems
- **Go idiomatic**: Perfect use of goroutines and channels
- **Testability**: Each stage can be tested independently

## Key Gotchas

- **Channel management**: Must properly close channels to avoid leaks
- **Deadlocks**: Improper channel handling can cause deadlocks
- **Buffering decisions**: Too much buffering wastes memory, too little causes blocking
- **Error handling**: Errors in one stage can break the pipeline
- **Backpressure**: Slow consumer can back up entire pipeline
- **Goroutine leaks**: Must ensure all goroutines complete properly
- **Complexity**: Large pipelines become hard to debug
- **Ordering**: Results may not arrive in input order
- **Resource limits**: Can't create unlimited goroutines/channels
