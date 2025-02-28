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

// AverageDuration returns the average time taken per Read operation.
// Returns 0 if no reads have been performed.
func (r *Reader) AverageDuration() time.Duration {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.readCount == 0 {
		return 0
	}

	return time.Duration(r.totalDuration.Nanoseconds() / r.readCount)
}
