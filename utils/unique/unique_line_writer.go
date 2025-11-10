// Package unique are utilities for unique line writer.
package unique

import "strings"

// LineWriter is a custom io.Writer that removes duplicate lines.
type LineWriter struct {
	builder *strings.Builder
	lines   map[string]bool
}

// NewLineWriter creates a new LineWriter.
func NewLineWriter() *LineWriter {
	return &LineWriter{
		builder: &strings.Builder{},
		lines:   make(map[string]bool),
	}
}

// Write appends data to the writer and removes duplicate lines.
func (w *LineWriter) Write(p []byte) (n int, err error) {
	// Split input into lines.
	lines := strings.SplitSeq(string(p), "\n")

	for line := range lines {
		if !w.lines[line] {
			// If the line is not a duplicate, write it and mark it as seen.
			w.builder.WriteString(line + "\n")
			w.lines[line] = true
		}
	}

	return len(p), nil
}

// String returns the  lines as a string.
func (w *LineWriter) String() string {
	return w.builder.String()
}

// Reset clears the buffer and the list of seen lines.
func (w *LineWriter) Reset() {
	w.builder.Reset()
	w.lines = make(map[string]bool)
}
