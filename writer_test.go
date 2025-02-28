package progress

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestNewWriter(t *testing.T) {
	is := is.New(t)

	// check Writer interfaces
	var (
		_ io.Writer    = (*Writer)(nil)
		_ Counter      = (*Writer)(nil)
		_ TimedCounter = (*Writer)(nil)
	)

	var buf bytes.Buffer
	w := NewWriter(&buf)

	n, err := w.Write([]byte("1"))
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(w.N(), int64(1)) // r.N()

	n, err = w.Write([]byte("1"))
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(w.N(), int64(2)) // r.N()

	n, err = w.Write([]byte("123"))
	is.NoErr(err)
	is.Equal(n, 3)            // n
	is.Equal(w.N(), int64(5)) // r.N()

	// Test AverageDuration
	avgDuration := w.AverageByteDuration()
	is.True(avgDuration >= 0)          // Average duration should be non-negative
	is.True(avgDuration < time.Second) // Average duration should be reasonable
}

// TestAverageDuration tests the AverageDuration method specifically
func TestWriterAverageDuration(t *testing.T) {
	is := is.New(t)

	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Initial average should be 0 when no writes have been performed
	is.Equal(w.AverageByteDuration(), time.Duration(0))

	// Perform a few writes
	for i := 0; i < 5; i++ {
		_, err := w.Write([]byte("test"))
		is.NoErr(err)
	}

	// Average duration should be non-zero after writes
	avgDuration := w.AverageByteDuration()
	is.True(avgDuration > 0)
	is.True(avgDuration < time.Second) // Sanity check for reasonable duration
}
