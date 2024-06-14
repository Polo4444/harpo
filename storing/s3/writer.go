package storing_s3

import (
	"errors"
	"io"
	"sync"
)

// WriterAtFromWriter is a wrapper that allows using an io.Writer as an io.WriterAt.
type WriterAtFromWriter struct {
	writer io.Writer
	mu     sync.Mutex
	offset int64
}

// NewWriterAtFromWriter creates a new WriterAtFromWriter.
func NewWriterAtFromWriter(writer io.Writer) *WriterAtFromWriter {
	return &WriterAtFromWriter{
		writer: writer,
	}
}

// WriteAt writes len(p) bytes from p to the underlying writer at offset off.
func (w *WriterAtFromWriter) WriteAt(p []byte, off int64) (n int, err error) {
	if off != w.offset {
		return 0, errors.New("non-sequential write not supported")
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	n, err = w.writer.Write(p)
	w.offset += int64(n)
	return n, err
}
