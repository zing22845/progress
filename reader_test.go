package progress

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestNewReader(t *testing.T) {
	is := is.New(t)

	// check Reader interfaces
	var (
		_ io.Reader    = (*Reader)(nil)
		_ Counter      = (*Reader)(nil)
		_ TimedCounter = (*Reader)(nil)
	)

	s := `Now that's what I call progress`
	r := NewReader(strings.NewReader(s))

	buf := make([]byte, 1)
	n, err := r.Read(buf)
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(r.N(), int64(1)) // r.N()

	n, err = r.Read(buf)
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(r.N(), int64(2)) // r.N()

	// read to the end
	b, err := io.ReadAll(r)
	is.NoErr(err)
	is.Equal(len(b), 29)       // len(b)
	is.Equal(r.N(), int64(31)) // r.N()

	// Test AverageDuration
	avgDuration := r.AverageByteDuration()
	is.True(avgDuration >= 0)          // Average duration should be non-negative
	is.True(avgDuration < time.Second) // Average duration should be reasonable
}

// TestReaderAverageDuration tests the AverageDuration method specifically
func TestReaderAverageDuration(t *testing.T) {
	is := is.New(t)

	// No reads performed yet
	r := NewReader(strings.NewReader("test data"))
	is.Equal(r.AverageByteDuration(), time.Duration(0))

	// Perform a few reads
	buf := make([]byte, 2)
	for i := 0; i < 3; i++ {
		_, err := r.Read(buf)
		is.NoErr(err)
	}

	// Average duration should be non-zero after reads
	avgDuration := r.AverageByteDuration()
	is.True(avgDuration > 0)
	is.True(avgDuration < time.Second) // Sanity check for reasonable duration
}
