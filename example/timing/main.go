package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/zing22845/progress"
)

func main() {
	fmt.Println("=== Writer Timing Example ===")
	writerExample()

	fmt.Println("\n=== Reader Timing Example ===")
	readerExample()
}

func writerExample() {
	var buf bytes.Buffer
	w := progress.NewWriter(&buf)

	// Write data in chunks of different sizes
	for i := 1; i <= 5; i++ {
		data := make([]byte, i*1000) // Chunks of increasing size
		start := time.Now()
		n, err := w.Write(data)
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			return
		}
		fmt.Printf("Wrote %d bytes in %v\n", n, time.Since(start))
	}

	// Get average duration
	fmt.Printf("Total bytes written: %d\n", w.N())
	fmt.Printf("Average write duration: %v\n", w.AverageDuration())
}

func readerExample() {
	// Create a large string to read from
	data := strings.Repeat("Hello, World! ", 1000)
	r := progress.NewReader(strings.NewReader(data))

	// Read data in small chunks
	buf := make([]byte, 128)
	readCount := 0

	for {
		start := time.Now()
		n, err := r.Read(buf)
		duration := time.Since(start)

		if n > 0 {
			readCount++
			fmt.Printf("Read #%d: %d bytes in %v\n", readCount, n, duration)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}
	}

	// Get average duration
	fmt.Printf("Total bytes read: %d\n", r.N())
	fmt.Printf("Number of read operations: %d\n", readCount)
	fmt.Printf("Average read duration: %v\n", r.AverageDuration())
}
