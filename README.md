# `progress` [![GoDoc](https://godoc.org/github.com/machinebox/progress?status.png)](http://godoc.org/github.com/machinebox/progress) [![Build Status](https://travis-ci.org/machinebox/progress.svg?branch=master)](https://travis-ci.org/machinebox/progress) [![Go Report Card](https://goreportcard.com/badge/github.com/machinebox/progress)](https://goreportcard.com/report/github.com/machinebox/progress)

`io.Reader` and `io.Writer` with progress, remaining time estimation, and operation timing.

## Usage

```go
ctx := context.Background()

// get a reader and the total expected number of bytes
s := `Now that's what I call progress`
size := len(s)
r := progress.NewReader(strings.NewReader(s))

// Start a goroutine printing progress
go func() {
	ctx := context.Background()
	progressChan := progress.NewTicker(ctx, r, size, 1*time.Second)
	for p := range progressChan {
		fmt.Printf("\r%v remaining...", p.Remaining().Round(time.Second))
	}
	fmt.Println("\rdownload is completed")
}()

// use the Reader as normal
if _, err := io.Copy(dest, r); err != nil {
	log.Fatalln(err)
}
```

1. Wrap an `io.Reader` or `io.Writer` with `NewReader` and `NewWriter` respectively
1. Capture the total number of expected bytes
1. Use `progress.NewTicker` to get a channel on which progress updates will be sent
1. Start a Goroutine to periodically check the progress, and do something with it - like log it
1. Use the readers and writers as normal

## Operation Timing

Both `Reader` and `Writer` now provide average duration metrics for read and write operations:

```go
// Get the average time taken per Read operation
reader := progress.NewReader(r)
// ... perform some reads ...
avgReadTime := reader.AverageDuration()
fmt.Printf("Average read time: %v\n", avgReadTime)

// Get the average time taken per Write operation
writer := progress.NewWriter(w)
// ... perform some writes ...
avgWriteTime := writer.AverageWriteDuration()
fmt.Printf("Average write time: %v\n", avgWriteTime)

// Get the average time taken per byte written
avgByteTime := writer.AverageByteDuration()
fmt.Printf("Average time per byte: %v\n", avgByteTime)

// Get all stats at once (total duration, write count, and bytes written)
totalDuration, writeCount, bytesWritten := writer.Stats()
fmt.Printf("Total duration: %v, Write operations: %d, Bytes written: %d\n", 
           totalDuration, writeCount, bytesWritten)
```

The `TimedCounter` interface is implemented by both `Reader` and `Writer`:

```go
type TimedCounter interface {
	Counter
	// AverageDuration returns the average time taken per operation.
	AverageDuration() time.Duration
}
```

The `Writer` provides additional timing methods:

```go
// AverageByteDuration returns the average time taken per byte written
AverageByteDuration() time.Duration

// AverageWriteDuration returns the average time taken per write operation
AverageWriteDuration() time.Duration

// Stats returns total duration, write count, and bytes written in one call
Stats() (time.Duration, int64, int64)
```

See the [timing example](example/timing/main.go) for more details on how to use this feature.
