package progress

import (
	"io"
	"sync"
	"time"
)

// Writer counts the bytes written through it.
type Writer struct {
	w io.Writer

	lock sync.RWMutex // protects n and err
	n    int64
	err  error

	// Time tracking fields
	totalDuration time.Duration
	writeCount    int64
}

// NewWriter gets a Writer that counts the number
// of bytes written.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	start := time.Now()
	n, err = w.w.Write(p)
	duration := time.Since(start)

	w.lock.Lock()
	w.n += int64(n)
	w.err = err
	w.totalDuration += duration
	w.writeCount++
	w.lock.Unlock()
	return
}

// N gets the number of bytes that have been written
// so far.
func (w *Writer) N() int64 {
	var n int64
	w.lock.RLock()
	n = w.n
	w.lock.RUnlock()
	return n
}

// Err gets the last error from the Writer.
func (w *Writer) Err() error {
	var err error
	w.lock.RLock()
	err = w.err
	w.lock.RUnlock()
	return err
}

// AverageDuration returns the average time taken per Write operation.
// Returns 0 if no writes have been performed.
func (w *Writer) AverageDuration() time.Duration {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.writeCount == 0 {
		return 0
	}

	return time.Duration(w.totalDuration.Nanoseconds() / w.writeCount)
}
