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

// AverageByteDuration returns the average time taken per byte written.
// Returns 0 if no bytes have been written.
func (w *Writer) AverageByteDuration() time.Duration {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.n == 0 {
		return 0
	}

	return w.totalDuration / time.Duration(w.n)
}

// AverageOperationDuration returns the average time taken per write operation.
// Returns 0 if no writes have been performed.
func (w *Writer) AverageOperationDuration() time.Duration {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.writeCount == 0 {
		return 0
	}

	return w.totalDuration / time.Duration(w.writeCount)
}

// Stats returns both the total duration and the number of bytes written.
// This allows retrieving both values atomically with a single lock acquisition.
func (w *Writer) Stats() (totalDuration time.Duration, totalCount int64, totalBytes int64) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.totalDuration, w.writeCount, w.n
}
