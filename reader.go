package progress

import (
	"io"
	"sync"
	"time"
)

// Reader counts the bytes read through it.
type Reader struct {
	r io.Reader

	lock sync.RWMutex // protects n and err
	n    int64
	err  error

	// Time tracking fields
	totalDuration time.Duration
	readCount     int64
}

// NewReader makes a new Reader that counts the bytes
// read through it.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	start := time.Now()
	n, err = r.r.Read(p)
	duration := time.Since(start)

	r.lock.Lock()
	r.n += int64(n)
	r.err = err
	r.totalDuration += duration
	r.readCount++
	r.lock.Unlock()
	return
}

// N gets the number of bytes that have been read
// so far.
func (r *Reader) N() int64 {
	var n int64
	r.lock.RLock()
	n = r.n
	r.lock.RUnlock()
	return n
}

// Err gets the last error from the Reader.
func (r *Reader) Err() error {
	var err error
	r.lock.RLock()
	err = r.err
	r.lock.RUnlock()
	return err
}

// AverageByteDuration returns the average time taken per byte read.
// Returns 0 if no bytes have been read.
func (r *Reader) AverageByteDuration() time.Duration {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.n == 0 {
		return 0
	}

	return r.totalDuration / time.Duration(r.n)
}

// AverageOperationDuration returns the average time taken per read operation.
// Returns 0 if no reads have been performed.
func (r *Reader) AverageOperationDuration() time.Duration {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.readCount == 0 {
		return 0
	}

	return r.totalDuration / time.Duration(r.readCount)
}

// Stats returns both the total duration and the number of bytes read.
// This allows retrieving both values atomically with a single lock acquisition.
func (r *Reader) Stats() (totalDuration time.Duration, totalCount int64, totalBytes int64) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.totalDuration, r.readCount, r.n
}
